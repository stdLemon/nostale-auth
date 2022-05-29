package gfclient_auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"testing"
)

const FINGERPRINT1_FILENAME = "./test_data/fingerprint1.json"
const FINGERPRINT2_FILENAME = "./test_data/fingerprint2.json"

func testBlackbox(t *testing.T, fingerprintFilename string) {
	content, err := ioutil.ReadFile(fingerprintFilename)
	if err != nil {
		t.Fatal(err)
	}

	finterprint := new(Fingerprint)
	err = json.Unmarshal(content, finterprint)
	if err != nil {
		t.Fatal(err)
	}

	blackbox, err := createBlackbox(finterprint)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := exec.Command("./gfclient_poc/create_blackbox.js", fingerprintFilename)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	expectedBlackbox := out.String()
	if blackbox != expectedBlackbox {
		t.Fatalf("blackbox was not encoded correctly \ngot: %v\nexpected: %v\n", blackbox, expectedBlackbox)
	}
}

func TestCreateBlackbox(t *testing.T) {
	testBlackbox(t, FINGERPRINT1_FILENAME)
	testBlackbox(t, FINGERPRINT2_FILENAME)
}
