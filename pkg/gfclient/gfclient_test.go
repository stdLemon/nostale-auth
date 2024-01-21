//go:build integration

package gfClient

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
	"github.com/stretchr/testify/require"
)

type GfAccountData struct {
	Email    string
	Password string
	Locale   string
	Name     string
}

func TestCodes(t *testing.T) {
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
	require.NoError(t, client.Init())

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
