package gfclient

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
)

const codesEndpoint = apiV1BaseURL + "/auth/thin/codes"

type CodesResponse struct {
	Message string
	Code    string
}

type CodesRequest struct {
	Blackbox              string `json:"blackbox"`
	GameID                string `json:"gameId"`
	GSID                  string `json:"gsid"`
	PlatformGameAccountID string `json:"platformGameAccountId"`
}

func (c *Client) Codes(bearer string, manager identitymgr.Manager, accountID, gameID string) (string, error) {
	gsID := generateGsid()
	encBlackbox, err := manager.NewEncryptedBlackbox(gsID, accountID)
	if err != nil {
		return "", err
	}

	ua, err := c.getCefUserAgent(accountID)
	if err != nil {
		return "", err
	}

	headers := make(http.Header, 2)
	headers.Set("User-Agent", ua)
	setAuthorizationHeader(headers, bearer)

	body := CodesRequest{
		Blackbox:              string(encBlackbox),
		GameID:                gameID,
		GSID:                  gsID,
		PlatformGameAccountID: accountID,
	}

	resp := new(CodesResponse)
	_, err = c.restClient.NewRequest().
		SetHeaders(headers).
		SetBody(body).
		SetResult(resp).
		Post(codesEndpoint)
	if err != nil {
		return "", err
	}

	if resp.Message != "" {
		return "", errors.New(resp.Message)
	}

	if resp.Code == "" {
		return "", errors.New("server did not send code")
	}

	return resp.Code, nil
}

func generateGsid() string {
	session := uuid.New().String()
	num := rand.Intn(9999) + 1
	return fmt.Sprintf("%s-%04d", session, num)
}
