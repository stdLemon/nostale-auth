package gfClient

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type GfAccountData struct {
	Email    string
	Password string
	Locale   string
	Name     string
}

func TestCalcCefUserAgentChecksum(t *testing.T) {
	tests := []struct {
		name             string
		installationId   string
		expectedChecksum string
	}{
		{
			name:             "odd installation id",
			installationId:   "c37f161c-7201-48a5-9a27-89c20a2a243d",
			expectedChecksum: "14fef9b2",
		},
		{
			name:             "even installation id",
			installationId:   "c47f161c-7201-48a5-9a27-89c20a2a243d",
			expectedChecksum: "74646fd4",
		},
	}

	const (
		accountId     = "f3779a30-fe2e-4eb9-b757-24897f4cd7d1"
		clientVersion = "2.8.0.1876"
	)

	c := Client{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.installationId = tt.installationId
			h := c.calcCefUserAgentChecksum(accountId, clientVersion)
			assert.Equal(t, tt.expectedChecksum, h)
		})
	}
}

func TestCodes(t *testing.T) {
	if os.Getenv("integration") != "true" {
		t.Skip("skipping integration test")
	}
	content, err := os.ReadFile("test/account.json")
	require.NoError(t, err)

	manager, err := identitymgr.New("test/identity.json")
	require.NoError(t, err)

	accountData := new(GfAccountData)
	require.NoError(t, json.Unmarshal(content, accountData))

	identity := manager.Get()
	client := New(
		identity.Fingerprint.UserAgent,
		identity.InstallationId,
	)

	bearer, err := client.Login(accountData.Email, accountData.Password, accountData.Locale, manager)
	require.NoError(t, err)
	require.NotEmpty(t, bearer, "bearer can't be empty")

	accountList, err := client.GetGameAccounts(bearer)
	require.NoError(t, err)

	account, ok := FindGameAccount(accountData.Name, accountList)
	require.True(t, ok, "account with name %s not found", accountData.Name)

	err = client.Iovation(bearer, manager, account.Id)
	require.NoError(t, err)

	code, err := client.Codes(bearer, manager, account.Id, account.GameId)
	require.NoError(t, err)

	require.NoError(t, client.Logout(bearer))
	manager.Save()

	t.Log("code", code)
}
