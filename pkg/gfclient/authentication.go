package gfclient

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
)

const (
	authSessionsEndpoint          = apiV1BaseURL + "/auth/sessions"
	authProvidersSessionsEndpoint = apiV2BaseURL + "/authProviders/credentials/sessions"
)
const (
	errorTypeCredentialsInvalid = "CREDENTIALS_INVALID"
	errorTypeChallengeRequired  = "CHALLENGE_REQUIRED"
)

const (
	headerGFInstallationID = "gf-installation-id"
)

type AuthErrorResponse struct {
	Message    string
	ErrorTypes []string
}

type AuthResponse struct {
	Token  string
	UserId string
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Locale   string `json:"locale"`
	Blackbox string `json:"blackbox"`
}

func (c *Client) handleLoginError(err error, errResp *AuthErrorResponse) (string, error) {
	switch {
	case slices.Contains(errResp.ErrorTypes, errorTypeCredentialsInvalid):
		return "", errors.New(errResp.Message)
	case slices.Contains(errResp.ErrorTypes, errorTypeChallengeRequired):
		return "", errors.New("challenge required, but challenge solving is not implemented yet")

	default:
		return "", fmt.Errorf("login failed: %w", err)
	}
}

func (c *Client) Login(email, password, locale string, manager identitymgr.Manager) (string, error) {
	blackbox, err := manager.NewBlackbox(nil)
	if err != nil {
		return "", err
	}

	var (
		authResp = new(AuthResponse)
		errResp  = new(AuthErrorResponse)
		headers  = make(http.Header, 2)
	)

	setOriginHeader(headers)
	headers.Set(headerGFInstallationID, c.installationId)

	_, err = c.restClient.NewRequest().
		SetHeaders(headers).
		SetBody(loginRequest{
			Email:    email,
			Password: password,
			Locale:   locale,
			Blackbox: blackbox.String(),
		}).
		SetResult(authResp).
		SetError(errResp).
		Post(authProvidersSessionsEndpoint)
	if err != nil {
		return c.handleLoginError(err, errResp)
	}

	if authResp.Token == "" {
		return "", errors.New("server didn't send a token")
	}

	return authResp.Token, nil
}

func (c *Client) Logout(bearer string) error {
	headers := make(http.Header, 1)
	setAuthorizationHeader(headers, bearer)

	_, err := c.restClient.NewRequest().
		SetHeaders(headers).
		Delete(authSessionsEndpoint)
	return err
}
