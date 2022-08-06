package blackbox

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/google/uuid"
)

const RANDOM_TEXT_FILENAME = "test/random_text.txt"

func TestEncryptBlackbox(t *testing.T) {
	c, err := os.ReadFile(RANDOM_TEXT_FILENAME)
	if err != nil {
		t.Fatal(err)
	}

	blackbox := Blackbox(c)
	gsId := uuid.New().String()
	accountId := uuid.New().String()
	encBlackbox := blackbox.Encrypt(gsId, accountId)

	var out bytes.Buffer
	cmd := exec.Command("gfclient_poc/encrypt_blackbox.js", RANDOM_TEXT_FILENAME, gsId, accountId)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	expected := out.Bytes()

	if !bytes.Equal(encBlackbox, expected) {
		t.Fatalf("blackbox was not encrypted correctly \ngot: %v\nexpected: %v\n", encBlackbox, expected)
	}
}
