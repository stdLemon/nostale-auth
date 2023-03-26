package gfClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
)

type Client struct {
	gfUserAgent    string
	cefUserAgent   string
	installationId string
	gfHeaders      http.Header
	httpClient     *http.Client
}

type GameAccount struct {
	Id          string
	DisplayName string
	GameId      string
}

type AuthResponse struct {
	Token string
}

type IovationResponse struct {
	Status string
}

type CodesResponse struct {
	Message string
	Code    string
}

func New(gfUserAgent, cefUserAgent, installationId string) *Client {
	return &Client{gfUserAgent: gfUserAgent,
		cefUserAgent:   cefUserAgent,
		installationId: installationId,
		httpClient:     new(http.Client),
		gfHeaders: map[string][]string{
			"TNT-Installation-Id": {installationId},
			"Origin":              {"spark://www.gameforge.com"},
			"User-Agent":          {gfUserAgent},
		},
	}
}

func (c *Client) Login(email, password, locale string, manager identitymgr.Manager) (bearer string, err error) {
	const url string = "https://spark.gameforge.com/api/v1/auth/sessions"

	blackbox, err := manager.NewBlackbox(nil)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
		"locale":   locale,
		"blackbox": blackbox.String(),
	})
	if err != nil {
		return
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	request.Header = map[string][]string{
		"Content-Type": {"application/json", "charset=UTF-8"},
	}
	c.addDefaultHeaders(request.Header)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return
	}

	switch resp.StatusCode {
	case http.StatusForbidden:
		err = errors.New("invalid account data")
		return
	case http.StatusConflict:
		err = errors.New("captcha")
		return
	}

	if err = checkStatusCode(http.StatusCreated, resp.StatusCode); err != nil {
		return
	}

	authResp := AuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return
	}

	if authResp.Token == "" {
		err = errors.New("server did not send token")
		return
	}

	bearer = authResp.Token
	return
}

func (c *Client) Logout(bearer string) error {
	const url string = "https://spark.gameforge.com/api/v1/auth/sessions"

	r, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	r.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", bearer)},
	}
	c.addDefaultHeaders(r.Header)
	httpResp, err := c.httpClient.Do(r)
	if err != nil {
		return err

	}
	return checkStatusCode(http.StatusAccepted, httpResp.StatusCode)
}

func (c *Client) GetGameAccounts(bearer string) ([]GameAccount, error) {
	const url string = "https://spark.gameforge.com/api/v1/user/accounts"

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	r.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", bearer)},
	}
	c.addDefaultHeaders(r.Header)

	httpResp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err

	}

	if err = checkStatusCode(http.StatusOK, httpResp.StatusCode); err != nil {
		return nil, err
	}

	resp := make(map[string]GameAccount)
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err

	}

	accs := make([]GameAccount, 0, len(resp))
	for _, v := range resp {
		accs = append(accs, v)
	}

	return accs, nil
}

func (c *Client) Iovation(bearer string, manager identitymgr.Manager, accountId string) error {
	const url string = "https://spark.gameforge.com/api/v1/auth/iovation"

	blackbox, err := manager.NewBlackbox(nil)
	if err != nil {
		return err
	}

	body, err := json.Marshal(map[string]string{
		"accountId": accountId,
		"blackbox":  blackbox.String(),
		"type":      "play_now",
	})
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	r.Header = map[string][]string{
		"Content-Type":  {"application/json", "charset=UTF-8"},
		"Authorization": {fmt.Sprintf("Bearer %s", bearer)},
	}
	c.addDefaultHeaders(r.Header)

	httpResp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}

	if err = checkStatusCode(http.StatusOK, httpResp.StatusCode); err != nil {
		return err
	}

	resp := new(IovationResponse)
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	if err != nil {
		return err

	}

	if resp.Status != "ok" {
		return errors.New(httpResp.Status)
	}

	return nil
}

func (c *Client) Codes(bearer string, manager identitymgr.Manager, accountId, gameId string) (string, error) {
	const url string = "https://spark.gameforge.com/api/v1/auth/thin/codes"

	gsId := generateGsid()
	encBlackbox, err := manager.NewEncryptedBlackbox(gsId, accountId)
	if err != nil {
		return "", nil
	}

	body, err := json.Marshal(map[string]string{
		"blackbox":              string(encBlackbox),
		"gameId":                gameId,
		"gsid":                  gsId,
		"platformGameAccountId": accountId,
	})
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err

	}

	r.Header = map[string][]string{
		"tnt-installation-id": {c.installationId},
		"Authorization":       {fmt.Sprintf("Bearer %s", bearer)},
		"User-Agent":          {c.cefUserAgent},
		"Content-Type":        {"application/json", "charset=UTF-8"},
	}

	httpResp, err := c.httpClient.Do(r)
	if err != nil {
		return "", err
	}

	resp := new(CodesResponse)
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	if err != nil {
		return "", err
	}

	if err = checkStatusCode(http.StatusCreated, httpResp.StatusCode); err != nil {
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

func FindGameAccount(name string, accounts []GameAccount) (GameAccount, error) {
	for _, acc := range accounts {
		if acc.DisplayName == name {
			return acc, nil
		}
	}
	return GameAccount{}, fmt.Errorf("account with name %s was not found", name)
}

func (c *Client) addDefaultHeaders(headers http.Header) {
	for k, v := range c.gfHeaders {
		headers[k] = v
	}
}

func checkStatusCode(expected, returned int) error {
	if expected != returned {
		return fmt.Errorf("got unexpected status code expected %s, got %s",
			http.StatusText(expected),
			http.StatusText(returned),
		)
	}
	return nil
}

func generateGsid() string {
	session := uuid.New().String()
	num := rand.Intn(9999-1) + 1
	return fmt.Sprintf("%s-%4d", session, num)
}
