package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CurrentTokenPrice(mint string) float64 {
	url := fmt.Sprintf("%sdefi/price?address=%s", BIRD_EYE_API_URL, mint)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-chain", "solana")

	res, _ := http.DefaultClient.Do(req)

	type Response struct {
		Data struct {
			Value float64 `json:"value"`
		} `json:"data"`
	}

	var response Response

	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
	}

	return response.Data.Value
}

func tokenPriceDaysAgo(mint string, numDays int) float64 {
	currentTime := time.Now().Unix()
	sevenDaysAgo := currentTime - int64(numDays*24*60*60)
	url := fmt.Sprintf("%sdefi/price?address=%s&address_type=token&type=1D&time_from=%d&time_to=%d", BIRD_EYE_API_URL, mint, sevenDaysAgo, currentTime)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-chain", "solana")

	res, _ := http.DefaultClient.Do(req)

	type Response struct {
		Data struct {
			Items []struct {
				Value float64 `json:"value"`
			} `json:"items"`
		} `json:"data"`
	}

	var response Response

	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
	}

	return response.Data.Items[0].Value
}

// func CalculatePNL(buyPrice float64, amount float64, mint string, time int) {

// 	currPrice := CurrentTokenPrice(mint)

// 	buy := buyPrice * amount
// 	curr := currPrice * amount

// 	pnl := (curr - buy) / amount

// }

// func CalculateROI(buyPrice float64, amount float64, mint string) {
// 	currTokenPrice := CurrentTokenPrice(mint)
// 	buy := buyPrice * amount
// 	curr := currTokenPrice * amount

// 	roi := ((curr - buy) / buy) * 100

// 	fmt.Println(roi)
// }

// // func CalculateHoldTime(mint string)

// func CalculateWinRate(buyPrice float64, amount float64, mint string) {
// 	currentTokenPrice := CurrentTokenPrice(mint)
// 	buy := buyPrice * amount
// 	curr := currentTokenPrice * amount

// 	winRate := (buy - curr) / curr

// 	fmt.Println(winRate)
// }
