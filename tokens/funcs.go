package tokens

import (
	"bot/m/v2/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func getLatestTokenProfiles() []string {
	res, err := http.Get(fmt.Sprintf("%s/token-profiles/latest/v1", utils.DEX_SCREENER_API_BASE))

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var responseJSON []TokenProfile

	var tokenAddresses []string
	if err := json.NewDecoder(res.Body).Decode(&responseJSON); err != nil {
		panic(err)
	}

	for _, token := range responseJSON {
		if len(token.TokenAddress) == 44 {
			tokenAddresses = append(tokenAddresses, token.TokenAddress)
		}
	}
	return tokenAddresses[0:10]
}

func getTokenBoosts() []string {
	res, err := http.Get(fmt.Sprintf("%s/token-boosts/latest/v1", utils.DEX_SCREENER_API_BASE))

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var responseJSON []TokenProfile

	var tokenAddresses []string
	if err := json.NewDecoder(res.Body).Decode(&responseJSON); err != nil {
		panic(err)
	}

	for _, token := range responseJSON {
		if len(token.TokenAddress) == 44 {
			tokenAddresses = append(tokenAddresses, token.TokenAddress)
		}
	}
	return tokenAddresses
}

func getTopTokenBoosts() []string {
	res, err := http.Get(fmt.Sprintf("%s/token-boosts/top/v1", utils.DEX_SCREENER_API_BASE))

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var responseJSON []TokenProfile

	var tokenAddresses []string
	if err := json.NewDecoder(res.Body).Decode(&responseJSON); err != nil {
		panic(err)
	}

	for _, token := range responseJSON {
		tokenAddresses = append(tokenAddresses, token.TokenAddress)
	}
	return tokenAddresses[0:10]
}

func GetTrendingTokens() []string {
	latestTokenProfiles := getLatestTokenProfiles()
	topTokenBoosts := getTopTokenBoosts()

	combination := []string{}

	combination = append(combination, latestTokenProfiles...)
	combination = append(combination, topTokenBoosts...)
	return combination
}
