package tokens

type TokenProfile struct {
	TokenAddress string `json:"tokenAddress"`
}

type TrendingTokenData struct {
	Data struct {
		Tokens []struct {
			Address string `json:"address"`
		} `json:"tokens"`
	} `json:"data"`
}

type Params struct {
	Limit int    `json:"limit"`
	Mint  string `json:"mint"`
}

type RequestBody struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Response struct {
	JsonRpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
}

type Result struct {
	Tokens []struct {
		Owner  string `json:"owner"`
		Amount int    `json:"amount"`
	} `json:"token_accounts"`
}
