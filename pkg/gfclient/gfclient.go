package gfClient

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"maps"
	"math/rand"
	"net/http"
	"strings"

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

type AuthErrorResponse struct {
	Message    string
	errorTypes []string
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
	errEmptyClientVersion = errors.New("server didn't send a client version")
	errTokenNotSent       = errors.New("server didn't send a token")
)

func headerAuthorization(bearer string) http.Header {
	return http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", bearer)},
	}
}

func headerJsonContentType() http.Header {
	return http.Header{
		"Content-Type": {"application/json"},
	}
}

func headerOrigin() http.Header {
	return http.Header{
		"Origin": {"spark://www.gameforge.com"},
	}
}

const (
	_clientVersionEndpoint = "http://dl.tnt.gameforge.com/tnt/final-ms3/clientversioninfo.json"
	_gameforgeSparkUrl     = "https://spark.gameforge.com"
	_apiV1BaseUrl          = _gameforgeSparkUrl + "/api/v1"
	_authSessionsEndpoint  = _apiV1BaseUrl + "/auth/sessions"
	_accountsEndpoint      = _apiV1BaseUrl + "/user/accounts"
	_iovationEndpoint      = _apiV1BaseUrl + "/auth/iovation"
	_codesEndpoint         = _apiV1BaseUrl + "/auth/thin/codes"
)

const (
	_apiV2BaseUrl                  = _gameforgeSparkUrl + "/api/v2"
	_authProvidersSessionsEndpoint = _apiV2BaseUrl + "/authProviders/credentials/sessions"
)

func New(gfUserAgent, installationId string) *Client {
	return &Client{gfUserAgent: gfUserAgent,
		installationId: installationId,
		httpClient:     new(http.Client),
		gfHeaders: map[string][]string{
			"tnt-installation-id": {installationId},
			"User-Agent":          {gfUserAgent},
		},
	}
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

	header := headerOrigin()
	maps.Copy(header, headerJsonContentType())

	httpResp, err := c.makeRequest(http.MethodPost, _authProvidersSessionsEndpoint, http.StatusCreated, bytes.NewBuffer(body), header)
	if err != nil {
		errResp := AuthErrorResponse{}
		err = json.NewDecoder(httpResp.Body).Decode(&errResp)
		if err != nil {
			return
		}
		return "", fmt.Errorf("login failed: %s", errResp.Message)
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
	_, err := c.makeRequest(http.MethodDelete, _authSessionsEndpoint, http.StatusAccepted, nil, headerAuthorization(bearer))
	return err
}

func (c *Client) GetGameAccounts(bearer string) ([]GameAccount, error) {
	header := headerOrigin()
	maps.Copy(header, headerAuthorization(bearer))

	httpResp, err := c.makeRequest(http.MethodGet, _accountsEndpoint, http.StatusOK, nil, header)
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
	maps.Copy(header, headerOrigin())
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

	ua, err := c.getCefUserAgent(accountId)
	if err != nil {
		return "", err
	}

	header := http.Header{
		"User-Agent": {ua},
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

func (c *Client) getCefUserAgent(accountId string) (string, error) {
	if c.cefUserAgent != "" {
		return c.cefUserAgent, nil
	}

	v, err := c.getVersion()
	if err != nil {
		return "", err
	}
	return c.createCefUserAgent(accountId, v), nil
}

func (c *Client) calcCefUserAgentChecksum(accountId, clientVersion string) string {
	const (
		certSha256 = "99025da70af1ef39d2acd049018887ef5140daebc6f11d80461bcf8d02f2d36b"
		certSha1   = "d68f9401f15791cc396d4d6af3b977bc58ad0002"
	)

	hashChain := func(hashers []hash.Hash, certHash string) string {
		var (
			sb     strings.Builder
			inputs = []string{"C" + clientVersion, c.installationId, accountId}
		)

		sb.WriteString(certHash)
		for i, data := range inputs {
			h := hashers[i]
			h.Write([]byte(data))
			e := hex.EncodeToString(h.Sum(nil))
			sb.WriteString(e)
		}

		return fmt.Sprintf("%x", sha256.Sum256([]byte(sb.String())))
	}

	var firstDigit byte
	// accountId contains only ASCII characters
	for _, r := range c.installationId {
		if r >= '0' && r <= '9' {
			firstDigit = byte(r)
			break
		}
	}

	var (
		evenHashers = []hash.Hash{sha1.New(), sha256.New(), sha1.New()}
		oddHashers  = []hash.Hash{sha256.New(), sha1.New(), sha256.New()}
	)

	if firstDigit == 0 || (firstDigit-'0')%2 == 0 {
		h := hashChain(evenHashers, certSha256)
		return strings.Clone(h[:8])

	}

	h := hashChain(oddHashers, certSha1)
	return strings.Clone(h[len(h)-8:])
}

func (c *Client) createCefUserAgent(accountId, clientVersion string) string {
	h := c.calcCefUserAgentChecksum(accountId, clientVersion)
	return fmt.Sprintf("Chrome/C%s (%s%s)", clientVersion, accountId[:2], h)
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
