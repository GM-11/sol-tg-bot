package wallet

type Criterias struct {
	MinAmount int `json:"minAmount"`
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

type TokenHoldersResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
}

type TransactionsResponse struct {
	Result struct {
		Data []struct {
			Signature string `json:"signature"`
		} `json:"data"`
	} `json:"result"`
}

type Result struct {
	Tokens []struct {
		Owner  string `json:"owner"`
		Amount int    `json:"amount"`
	} `json:"token_accounts"`
}
