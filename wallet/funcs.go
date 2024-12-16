package wallet

import (
	"bot/m/v2/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetTokenHolders(tokenAddress string, minBalance float64, minPNL float64, minROI float64, minHoldTime int64) []string {

	helius := fmt.Sprintf("https://mainnet.helius-rpc.com/?api-key=%s", utils.HELIUS_API_KEY)

	params := Params{
		Limit: 20,
		Mint:  tokenAddress,
	}

	requestBody := RequestBody{
		Jsonrpc: "2.0",
		Id:      "helius-test",
		Method:  "getTokenAccounts",
		Params:  params,
	}

	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	res, err := http.Post(helius, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	var response TokenHoldersResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		panic(err)
	}
	holders := make([]string, 0)
	for _, token := range response.Result.Tokens {
		fmt.Printf("\t Processing owner: %s\n", token.Owner)
		sigs := GetWalletTransactionsSignatures(token.Owner, tokenAddress)

		if GetWalletTransactions(token.Owner, sigs, tokenAddress, minPNL, minROI, minHoldTime) && token.Amount > int(minBalance) {
			fmt.Printf("\t\t Found holder: %s\n", token.Owner)
			holders = append(holders, token.Owner)
		}

	}
	return holders
}

func GetWalletTransactionsSignatures(walletAddress string, tokenMint string) []string {
	url := fmt.Sprintf("%s%s/transactions?mint=%s&limit=4", utils.SOLANA_FM_API_URL, walletAddress, tokenMint)

	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	signatures := make([]string, 0)

	var response TransactionsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	for _, transaction := range response.Result.Data {

		signatures = append(signatures, transaction.Signature)
	}

	return signatures
}

func GetWalletTransactions(walletAddress string, signatures []string, mint string, minReqPNL float64, minReqROIance float64, minReqHoldTime int64) bool {
	helius := fmt.Sprintf("https://api.helius.xyz/v0/transactions?api-key=%s", utils.HELIUS_API_KEY)

	buys := make([]float64, 0)
	sells := make([]float64, 0)
	buyTimeStamp := make([]int, 0)

	var pnl float64 = 0
	var roi float64 = 0
	type RequestBody struct {
		Transactions []string `json:"transactions"`
	}

	for _, signature := range signatures {
		requestBody := RequestBody{
			Transactions: []string{signature},
		}
		jsonValue, err := json.Marshal(requestBody)

		if err != nil {
			panic(err)
		}
		res, err := http.Post(helius, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			panic(err)
		}

		defer res.Body.Close()

		type HeliusResponse struct {
			Type           string `json:"type"`
			Timestamp      int    `json:"timestamp"`
			TokenTransfers []struct {
				Mint            string  `json:"mint"`
				FromUserAccount string  `json:"fromUserAccount"`
				ToUserAccount   string  `json:"toUserAccount"`
				TokenAmount     float64 `json:"tokenAmount"`
			} `json:"tokenTransfers"`
		}
		var transactions []HeliusResponse
		err = json.NewDecoder(res.Body).Decode(&transactions)
		if err != nil {
			return false
		}

		if len(transactions) == 0 {
			return false
		}

		for _, res := range transactions {
			if len(res.TokenTransfers) == 0 {
				return false
			}

			index := -1

			for i, token := range res.TokenTransfers {
				if token.Mint == mint {
					index = i
				}
			}

			if index == -1 {
				return false
			}

			for _, token := range res.TokenTransfers {
				if token.ToUserAccount == walletAddress && token.Mint == mint {
					price := getPriceAtTime(mint, res.Timestamp)
					buys = append(buys, token.TokenAmount*price)

					roi = roi - token.TokenAmount*price
					buyTimeStamp = append(buyTimeStamp, res.Timestamp)
				} else if token.FromUserAccount == walletAddress && token.Mint == mint {
					price := getPriceAtTime(mint, res.Timestamp)
					sells = append(sells, token.TokenAmount*price)

					if len(buys) == 0 {
						return false
					}

					buy := buys[len(buys)-1]
					_pnl := (token.TokenAmount*price - buy) / token.TokenAmount
					if _pnl > 0 {
						roi = roi + _pnl
					}
					pnl = pnl + _pnl
				} else {
					continue
				}
			}
		}
	}
	roi = roi / sum(buys)
	holdTime := time.Now().Unix() - min(buyTimeStamp)

	if roi < minReqROIance || holdTime < minReqHoldTime || pnl < minReqPNL {
		return false
	}

	return true

}

func getPriceAtTime(mint string, timestamp int) float64 {
	priceUrl := fmt.Sprintf("%sdefi/history_price?address_type=token&type=1m&address=%s&time_from=%d&time_to=%d", utils.BIRD_EYE_API_URL, mint, timestamp-60, timestamp+600)

	client := &http.Client{}
	req, err := http.NewRequest("GET", priceUrl, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-chain", "solana")
	req.Header.Add("X-API-KEY", "a43ab680e43d4459b82f433976b0e9bd")

	r, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	type BirdEyeResponse struct {
		Data struct {
			Items []struct {
				Value float64 `json:"value"`
			} `json:"items"`
		} `json:"data"`
	}

	var response BirdEyeResponse
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	if len(response.Data.Items) == 0 {
		return 0
	}
	return response.Data.Items[0].Value
}

func sum(list []float64) float64 {
	var s float64
	for _, l := range list {
		s = s + l
	}
	return s
}

func min(list []int) int64 {
	var m int64 = time.Now().Unix()
	for _, l := range list {
		if int64(l) < m {
			m = int64(l)
		}
	}
	return m
}
