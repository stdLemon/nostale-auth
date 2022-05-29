package gfclient_auth

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type GfAccountData struct {
	Email    string
	Password string
	Locale   string
	Name     string
}

func TestCodes(t *testing.T) {
	content, err := ioutil.ReadFile("account.json")
	if err != nil {
		t.Fatal(err)
	}

	identity_manager, err := NewIdentityManager("identity.json")
	if err != nil {
		t.Fatal(err)
	}

	account_data := new(GfAccountData)
	json.Unmarshal(content, account_data)

	identity := identity_manager.Get()

	gfclient := NewGfClient(
		identity.Fingerprint.UserAgent,
		"Chrome/C2.2.23.1813 (49c0acbee1)",
		identity.Installation_id,
	)

	err = gfclient.Auth(account_data.Email, account_data.Password, account_data.Locale)
	if err != nil {
		t.Fatal(err)
	}

	game_account_list, err := gfclient.GetGameAccounts()
	if err != nil {
		t.Fatal(err)
	}

	game_account, err := FindGameAccount(account_data.Name, game_account_list)
	if err != nil {
		t.Fatal(err)
	}

	err = gfclient.Iovation(identity_manager, game_account.Id)
	if err != nil {
		t.Fatal(err)
	}

	code, err := gfclient.Codes(identity_manager, game_account.Id, game_account.GameId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code", code)
	identity_manager.Save()
}
