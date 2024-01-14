//go:build integration

package gfClient

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
)

type GfAccountData struct {
	Email    string
	Password string
	Locale   string
	Name     string
}

func TestCodes(t *testing.T) {
	content, err := ioutil.ReadFile("test/account.json")
	if err != nil {
		t.Fatal(err)
	}

	manager, err := identitymgr.New("test/identity.json")
	if err != nil {
		t.Fatal(err)
	}

	accountData := new(GfAccountData)
	json.Unmarshal(content, accountData)

	identity := manager.Get()

	client := New(
		identity.Fingerprint.UserAgent,
		"Chrome/C2.5.0.1857 (49a6fac338)",
		identity.InstallationId,
	)

	bearer, err := client.Login(accountData.Email, accountData.Password, accountData.Locale, manager)
	if err != nil {
		t.Fatal(err)
	}

	if bearer == "" {
		t.Fatal("bearer can't be empty")
	}

	accountList, err := client.GetGameAccounts(bearer)
	if err != nil {
		t.Fatal(err)
	}

	account, err := FindGameAccount(accountData.Name, accountList)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Iovation(bearer, manager, account.Id)
	if err != nil {
		t.Fatal(err)
	}

	code, err := client.Codes(bearer, manager, account.Id, account.GameId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code", code)
	manager.Save()
}
