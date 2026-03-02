package gfclient

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"maps"
	"math/rand"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
	"github.com/stdLemon/nostale-auth/pkg/restclient"
)

type Client struct {
	gfUserAgent    string
	cefUserAgent   string
	installationId string
	httpClient     restclient.Client
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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Locale   string `json:"locale"`
	Blackbox string `json:"blackbox"`
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
		httpClient: restclient.New(map[string][]string{
			"tnt-installation-id": {installationId},
			"User-Agent":          {gfUserAgent},
		}),
	}
}

func (c *Client) Login(email, password, locale string, manager identitymgr.Manager) (string, error) {
	blackbox, err := manager.NewBlackbox(nil)
	if err != nil {
		return "", err
	}

	header := headerOrigin()
	maps.Copy(header, headerJsonContentType())

	authResp := new(AuthResponse)
	errResp := new(AuthErrorResponse)
	_, err = c.httpClient.NewRequest().
		SetHeaders(header).
		SetBody(loginRequest{
			Email:    email,
			Password: password,
			Locale:   locale,
			Blackbox: blackbox.String(),
		}).
		SetResult(authResp).
		SetError(errResp).
		Post(_authProvidersSessionsEndpoint)

	if err != nil {
		return "", fmt.Errorf("login failed: %w response: %q", err, errResp.Message)
	}

	if authResp.Token == "" {
		return "", errTokenNotSent
	}

	return authResp.Token, nil
}

func (c *Client) Logout(bearer string) error {
	_, err := c.httpClient.NewRequest().
		SetHeaders(headerAuthorization(bearer)).
		Delete(_authSessionsEndpoint)

	if err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	return nil
}

func (c *Client) GetGameAccounts(bearer string) ([]GameAccount, error) {
	header := headerOrigin()
	maps.Copy(header, headerAuthorization(bearer))

	resp := make(map[string]GameAccount)
	_, err := c.httpClient.NewRequest().
		SetHeaders(header).
		SetResult(&resp).
		Get(_accountsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("getting game accounts failed: %w", err)
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

	body := map[string]string{
		"accountId": accountId,
		"blackbox":  blackbox.String(),
		"type":      "play_now",
	}

	header := headerJsonContentType()
	maps.Copy(header, headerOrigin())
	maps.Copy(header, headerAuthorization(bearer))

	resp := new(IovationResponse)
	httpResp, err := c.httpClient.NewRequest().
		SetHeaders(header).
		SetBody(body).
		SetResult(resp).
		Post(_iovationEndpoint)
	if err != nil {
		return fmt.Errorf("iovation failed: %w", err)
	}

	if resp.Status != "ok" {
		return errors.New(httpResp.Status())
	}

	return nil
}

func (c *Client) Codes(bearer string, manager identitymgr.Manager, accountId, gameId string) (string, error) {
	gsId := generateGsid()
	encBlackbox, err := manager.NewEncryptedBlackbox(gsId, accountId)
	if err != nil {
		return "", err
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

	body := map[string]string{
		"blackbox":              string(encBlackbox),
		"gameId":                gameId,
		"gsid":                  gsId,
		"platformGameAccountId": accountId,
	}

	resp := new(CodesResponse)
	_, err = c.httpClient.NewRequest().
		SetHeaders(header).
		SetBody(body).
		SetResult(resp).
		Post(_codesEndpoint)
	if err != nil {
		return "", fmt.Errorf("getting codes failed: %w", err)
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

	if (firstDigit-'0')%2 == 0 {
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
	respJson := make(map[string]any)
	_, err := c.httpClient.NewRequest().
		SetResult(&respJson).
		Get(_clientVersionEndpoint)
	if err != nil {
		return "", fmt.Errorf("getting client version failed: %w", err)
	}

	clientVersion, ok := respJson["version"].(string)
	if !ok {
		return "", errEmptyClientVersion
	}

	return clientVersion, nil
}

func generateGsid() string {
	session := uuid.New().String()
	num := rand.Intn(9999) + 1
	return fmt.Sprintf("%s-%4d", session, num)
}
