package gfclient_auth

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/google/uuid"
)

const RANDOM_TEXT_FILENAME = "./test_data/random_text.txt"

func TestEncryptBlackbox(t *testing.T) {
	content, err := ioutil.ReadFile(RANDOM_TEXT_FILENAME)
	if err != nil {
		t.Fatal(err)
	}

	blackbox := string(content)
	gs_id := uuid.New().String()
	account_id := uuid.New().String()
	encrypted_blackbox := encryptBlackbox(blackbox, gs_id, account_id)

	var out bytes.Buffer
	cmd := exec.Command("./gfclient_poc/encrypt_blackbox.js", RANDOM_TEXT_FILENAME, gs_id, account_id)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	expected := out.Bytes()

	if bytes.Compare(encrypted_blackbox, expected) != 0 {
		t.Fatalf("blackbox was not encrypted correctly \ngot: %v\nexpected: %v\n", encrypted_blackbox, expected)
	}
}
