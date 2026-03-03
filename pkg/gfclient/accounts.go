package gfclient

import (
	"fmt"
	"net/http"
)

const accountsEndpoint = apiV1BaseURL + "/user/accounts"

type GameAccount struct {
	Id          string
	DisplayName string
	GameId      string
}

func (c *Client) GetGameAccounts(bearer string) ([]GameAccount, error) {
	headers := make(http.Header, 2)
	setOriginHeader(headers)
	setAuthorizationHeader(headers, bearer)

	resp := make(map[string]GameAccount)
	_, err := c.restClient.NewRequest().
		SetHeaders(headers).
		SetResult(&resp).
		Get(accountsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get game accounts: %w", err)
	}

	accounts := make([]GameAccount, 0, len(resp))
	for _, account := range resp {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func FindGameAccount(name string, accounts []GameAccount) (GameAccount, bool) {
	for _, account := range accounts {
		if account.DisplayName == name {
			return account, true
		}
	}

	return GameAccount{}, false
}
