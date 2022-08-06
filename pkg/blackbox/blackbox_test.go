package blackbox

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
)

const FINGERPRINT1_FILENAME = "test/fingerprint1.json"
const FINGERPRINT2_FILENAME = "test/fingerprint2.json"

func testBlackbox(t *testing.T, fingerprintFilename string) {
	content, err := os.ReadFile(fingerprintFilename)
	if err != nil {
		t.Fatal(err)
	}

	finterprint := new(Fingerprint)
	err = json.Unmarshal(content, finterprint)
	if err != nil {
		t.Fatal(err)
	}

	blackbox, err := New(finterprint)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := exec.Command("gfclient_poc/create_blackbox.js", fingerprintFilename)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	expectedBlackbox := out.String()
	if blackbox.String() != expectedBlackbox {
		t.Fatalf("blackbox was not encoded correctly \ngot: %v\nexpected: %v\n", blackbox, expectedBlackbox)
	}
}

func TestCreateBlackbox(t *testing.T) {
	testBlackbox(t, FINGERPRINT1_FILENAME)
	testBlackbox(t, FINGERPRINT2_FILENAME)
}
