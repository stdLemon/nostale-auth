package gfClient

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
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

var (
	errInvalidAccountData = errors.New("invalid account data")
	errEmptyClientVersion = errors.New("server didn't send a client version")
	errCaptchaRequired    = errors.New("captcha is required")
	errTokenNotSent       = errors.New("server didn't send a token")
)

func headerAuthorization(bearer string) http.Header {
	return http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", bearer)},
	}
}

func headerJsonContentType() http.Header {
	return http.Header{
		"Content-Type": {"application/json", "charset=UTF-8"},
	}
}

const (
	_clientVersionEndpoint = "http://dl.tnt.gameforge.com/tnt/final-ms3/clientversioninfo.json"
	_gameforgeSparkUrl     = "https://spark.gameforge.com"
	_apiV1BaseUrl          = _gameforgeSparkUrl + "/api/v1"
	_sessionsEndpoint      = _apiV1BaseUrl + "/auth/sessions"
	_accountsEndpoint      = _apiV1BaseUrl + "/user/accounts"
	_iovationEndpoint      = _apiV1BaseUrl + "/auth/iovation"
	_codesEndpoint         = _apiV1BaseUrl + "/auth/thin/codes"
)

func New(gfUserAgent, installationId string) *Client {
	return &Client{gfUserAgent: gfUserAgent,
		installationId: installationId,
		httpClient:     new(http.Client),
		gfHeaders: map[string][]string{
			"Origin":              {_gameforgeSparkUrl},
			"tnt-installation-id": {installationId},
			"User-Agent":          {gfUserAgent},
		},
	}
}

func (c *Client) Init() error {
	userAgent, err := c.generateCefUserAgent()
	if err != nil {
		return nil
	}
	c.cefUserAgent = userAgent
	return nil
}

func (c *Client) Login(email, password, locale string, manager identitymgr.Manager) (bearer string, err error) {
	blackbox, err := manager.NewBlackbox(nil)
	if err != nil {
		return
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

	httpResp, err := c.makeRequest(http.MethodPost, _sessionsEndpoint, http.StatusCreated, bytes.NewBuffer(body), headerJsonContentType())
	if err != nil {
		switch httpResp.StatusCode {
		case http.StatusForbidden:
			err = errInvalidAccountData
			return
		case http.StatusConflict:
			err = errCaptchaRequired
			return
		default:
			return
		}

	}

	authResp := AuthResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&authResp)
	if err != nil {
		return
	}

	if authResp.Token == "" {
		err = errTokenNotSent
		return
	}

	return authResp.Token, nil
}

func (c *Client) Logout(bearer string) error {
	_, err := c.makeRequest(http.MethodDelete, _sessionsEndpoint, http.StatusAccepted, nil, headerAuthorization(bearer))
	return err
}

func (c *Client) GetGameAccounts(bearer string) ([]GameAccount, error) {
	httpResp, err := c.makeRequest(http.MethodGet, _accountsEndpoint, http.StatusOK, nil, headerAuthorization(bearer))
	if err != nil {
		return nil, err
	}

	resp := make(map[string]GameAccount)
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	accs := make([]GameAccount, 0, len(resp))
	for _, v := range resp {
		accs = append(accs, v)
	}

	return accs, nil
}

func (c *Client) Iovation(bearer string, manager identitymgr.Manager, accountId string) error {
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

	header := headerJsonContentType()
	maps.Copy(header, headerAuthorization(bearer))

	httpResp, err := c.makeRequest(http.MethodPost, _iovationEndpoint, http.StatusOK, bytes.NewBuffer(body), header)
	if err != nil {
		return err
	}

	resp := new(IovationResponse)
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return err

	}

	if resp.Status != "ok" {
		return errors.New(httpResp.Status)
	}

	return nil
}

func (c *Client) Codes(bearer string, manager identitymgr.Manager, accountId, gameId string) (string, error) {
	gsId := generateGsid()
	encBlackbox, err := manager.NewEncryptedBlackbox(gsId, accountId)
	if err != nil {
		return "", nil
	}

	header := http.Header{
		"User-Agent": {c.cefUserAgent},
	}
	maps.Copy(header, headerJsonContentType())
	maps.Copy(header, headerAuthorization(bearer))

	body, err := json.Marshal(map[string]string{
		"blackbox":              string(encBlackbox),
		"gameId":                gameId,
		"gsid":                  gsId,
		"platformGameAccountId": accountId,
	})
	if err != nil {
		return "", err
	}

	httpResp, err := c.makeRequest(http.MethodPost, _codesEndpoint, http.StatusCreated, bytes.NewBuffer(body), header)
	if err != nil {
		return "", err
	}
	resp := new(CodesResponse)
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
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

func FindGameAccount(name string, accounts []GameAccount) (GameAccount, bool) {
	for _, acc := range accounts {
		if acc.DisplayName == name {
			return acc, true
		}
	}
	return GameAccount{}, false
}

// lazy user agenet generation
// other than the launcher implementation
func (c Client) generateCefUserAgent() (string, error) {
	clientVersion, err := c.getVersion()
	if err != nil {
		return "", err
	}

	instalationHash := fmt.Sprintf("%x", sha256.Sum256([]byte(c.installationId)))
	return fmt.Sprintf("Chrome/C%s (%s%s)", clientVersion, c.installationId[:2], instalationHash[:8]), nil
}

func (c Client) getVersion() (string, error) {
	httpResp, err := c.makeRequest(http.MethodGet, _clientVersionEndpoint, http.StatusOK, nil, nil)
	if err != nil {
		return "", err
	}

	respJson := make(map[string]interface{})
	if err := json.NewDecoder(httpResp.Body).Decode(&respJson); err != nil {
		return "", err
	}

	clientVersion, ok := respJson["version"]
	if !ok {
		return "", errEmptyClientVersion
	}

	return clientVersion.(string), nil
}

func (c Client) makeRequest(method, endpoint string, expectedStatusCode int, body io.Reader, header http.Header) (*http.Response, error) {
	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	request.Header = c.gfHeaders.Clone()
	maps.Copy(request.Header, header)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err

	}

	if expectedStatusCode != response.StatusCode {
		return response, fmt.Errorf("got unexpected status code expected %s, got %s",
			http.StatusText(expectedStatusCode),
			http.StatusText(response.StatusCode),
		)
	}

	return response, nil
}

func generateGsid() string {
	session := uuid.New().String()
	num := rand.Intn(9999-1) + 1
	return fmt.Sprintf("%s-%4d", session, num)
}
