package gfclient_auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
)

type GfClient struct {
	gf_user_agent   string
	cef_user_agent  string
	installation_id string
	bearer          string
	gf_headers      http.Header
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

func NewGfClient(gf_user_agent, cef_user_agent, installation_id string) GfClient {
	return GfClient{gf_user_agent: gf_user_agent,
		cef_user_agent:  cef_user_agent,
		installation_id: installation_id,
		gf_headers: map[string][]string{
			"TNT-Installation-Id": {installation_id},
			"Origin":              {"spark://www.gameforge.com"},
			"User-Agent":          {gf_user_agent},
		},
	}
}

func (client *GfClient) Auth(email, password, locale string) error {
	const url string = "https://spark.gameforge.com/api/v1/auth/sessions"

	body, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
		"locale":   locale,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header = map[string][]string{
		"Content-Type": {"application/json", "charset=UTF-8"},
	}
	client.addDefaultHeaders(request.Header)

	http_client := new(http.Client)
	http_response, err := http_client.Do(request)
	if err != nil {
		return err
	}

	switch http_response.StatusCode {
	case http.StatusForbidden:
		return errors.New("invalid account data")
	case http.StatusConflict:
		return errors.New("captcha")
	}

	var parsed_response AuthResponse
	err = json.NewDecoder(http_response.Body).Decode(&parsed_response)
	if err != nil {
		return err
	}

	if parsed_response.Token == "" {
		return errors.New("server did not send token")
	}

	client.bearer = parsed_response.Token
	return nil
}

func (client *GfClient) GetGameAccounts() ([]GameAccount, error) {
	const url string = "https://spark.gameforge.com/api/v1/user/accounts"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", client.bearer)},
	}
	client.addDefaultHeaders(request.Header)

	http_client := new(http.Client)
	resp, err := http_client.Do(request)
	if err != nil {
		return nil, err

	}

	parsed_response := make(map[string]GameAccount)
	err = json.NewDecoder(resp.Body).Decode(&parsed_response)
	if err != nil {
		return nil, err

	}

	game_accounts := make([]GameAccount, 0, len(parsed_response))
	for _, v := range parsed_response {
		game_accounts = append(game_accounts, v)
	}

	return game_accounts, nil
}

func (client *GfClient) Iovation(identity_manager IdentityManager, account_id string) error {
	const url string = "https://spark.gameforge.com/api/v1/auth/iovation"

	blackbox, err := NewBlackbox(identity_manager, nil)
	if err != nil {
		return nil
	}

	body, err := json.Marshal(map[string]string{
		"accountId": account_id,
		"blackbox":  blackbox,
		"type":      "play_now",
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header = map[string][]string{
		"Content-Type":  {"application/json", "charset=UTF-8"},
		"Authorization": {fmt.Sprintf("Bearer %s", client.bearer)},
	}
	client.addDefaultHeaders(request.Header)

	http_client := new(http.Client)
	resp, err := http_client.Do(request)
	if err != nil {
		return err
	}

	parsed_response := new(IovationResponse)
	err = json.NewDecoder(resp.Body).Decode(parsed_response)
	defer resp.Body.Close()
	if err != nil {
		return err

	}

	if parsed_response.Status != "ok" {
		return errors.New(resp.Status)
	}

	return nil
}

func (client *GfClient) Codes(identity_manager IdentityManager, account_id, game_id string) (string, error) {
	const url string = "https://spark.gameforge.com/api/v1/auth/thin/codes"

	gs_id := generateGsid()
	encrypted_blackbox, err := NewEncryptedBlackbox(identity_manager, gs_id, account_id, client.installation_id)
	if err != nil {
		return "", nil
	}

	body, err := json.Marshal(map[string]string{
		"blackbox":              string(encrypted_blackbox),
		"gameId":                game_id,
		"gsid":                  gs_id,
		"platformGameAccountId": account_id,
	})
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err

	}

	request.Header = map[string][]string{
		"tnt-installation-id": {client.installation_id},
		"Authorization":       {fmt.Sprintf("Bearer %s", client.bearer)},
		"User-Agent":          {client.cef_user_agent},
		"Content-Type":        {"application/json", "charset=UTF-8"},
	}

	http_client := new(http.Client)
	resp, err := http_client.Do(request)
	if err != nil {
		return "", err
	}

	parsed_response := new(CodesResponse)
	err = json.NewDecoder(resp.Body).Decode(parsed_response)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New(fmt.Sprintf("expected 201 Created status, got %s", resp.Status))
	}

	if parsed_response.Message != "" {
		return "", errors.New(parsed_response.Message)
	}

	if parsed_response.Code == "" {
		return "", errors.New("server did not send code")
	}

	return parsed_response.Code, nil
}

func FindGameAccount(name string, accounts []GameAccount) (GameAccount, error) {
	for _, acc := range accounts {
		if acc.DisplayName == name {
			return acc, nil
		}
	}
	return GameAccount{}, errors.New(fmt.Sprintf("account with name %s was not found", name))
}

func (client *GfClient) addDefaultHeaders(headers http.Header) {
	for k, v := range client.gf_headers {
		headers[k] = v
	}
}

func generateGsid() string {
	session := uuid.New().String()
	num := rand.Intn(9999-1) + 1
	return fmt.Sprintf("%s-%4d", session, num)
}
