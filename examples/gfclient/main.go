package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/stdLemon/nostale-auth/pkg/gfclient"
	"github.com/stdLemon/nostale-auth/pkg/identitymgr"
)

type GfAccountData struct {
	Email    string
	Password string
	Locale   string
	Name     string
}

func main() {
	content, err := os.ReadFile("account.json")
	if err != nil {
		panic(err)
	}

	manager, err := identitymgr.New("identity.json")
	if err != nil {
		panic(err)
	}

	accountData := new(GfAccountData)
	if err := json.Unmarshal(content, accountData); err != nil {
		panic(err)
	}

	identity := manager.Get()
	client := gfclient.New(
		identity.Fingerprint.UserAgent,
		identity.InstallationId,
	)

	bearer, err := client.Login(accountData.Email, accountData.Password, accountData.Locale, manager)
	if err != nil {
		panic(err)
	}
	if bearer == "" {
		panic("bearer can't be empty")
	}

	accountList, err := client.GetGameAccounts(bearer)
	if err != nil {
		panic(err)
	}

	account, ok := gfclient.FindGameAccount(accountData.Name, accountList)
	if !ok {
		panic(fmt.Sprintf("account with name %s not found", accountData.Name))
	}

	if err := client.Iovation(bearer, manager, account.Id); err != nil {
		panic(err)
	}

	code, err := client.Codes(bearer, manager, account.Id, account.GameId)
	if err != nil {
		panic(err)
	}

	if err := client.Logout(bearer); err != nil {
		panic(err)
	}
	if err := manager.Save(); err != nil {
		panic(err)
	}

	fmt.Println("code", code)
}
