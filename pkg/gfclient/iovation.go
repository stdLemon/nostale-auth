package gfclient

import (
	"errors"
	"net/http"

	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
)

const iovationEndpoint = apiV1BaseURL + "/auth/iovation"

type IovationResponse struct {
	Status string
}

type IovationRequest struct {
	AccountID string `json:"accountId"`
	Blackbox  string `json:"blackbox"`
	Type      string `json:"type"`
}

func (c *Client) Iovation(bearer string, manager identitymgr.Manager, accountID string) error {
	blackbox, err := manager.NewBlackbox(nil)
	if err != nil {
		return err
	}

	body := IovationRequest{
		AccountID: accountID,
		Blackbox:  blackbox.String(),
		Type:      "play_now",
	}

	headers := make(http.Header, 2)
	setOriginHeader(headers)
	setAuthorizationHeader(headers, bearer)

	resp := new(IovationResponse)
	httpResp, err := c.restClient.NewRequest().
		SetHeaders(headers).
		SetBody(body).
		SetResult(resp).
		Post(iovationEndpoint)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return errors.New(httpResp.Status())
	}

	return nil
}
