package gfclient_auth

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/url"
	"strings"

	"github.com/fatih/structs"
)

func createBlackbox(fingerprint *Fingerprint) (string, error) {
	fields := structs.Fields(fingerprint)
	values := make([]interface{}, len(fields))
	for i := range fields {
		values[i] = fields[i].Value()
	}

	json_array, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	uri_encoded := url.QueryEscape(string(json_array))
	uri_encoded = strings.ReplaceAll(uri_encoded, "+", "%20")
	uri_encoded = strings.ReplaceAll(uri_encoded, "%29", ")")
	uri_encoded = strings.ReplaceAll(uri_encoded, "%28", "(")

	blackbox := make([]byte, len(uri_encoded))
	blackbox[0] = uri_encoded[0]
	for i := 1; i < len(uri_encoded); i++ {
		blackbox[i] = blackbox[i-1] + uri_encoded[i]
	}

	return "tra:" + base64.RawURLEncoding.EncodeToString(blackbox), nil
}

func NewBlackbox(identity_manager IdentityManager, request *Request) (string, error) {
	fingerprint, err := createFingerprint(identity_manager)
	if err != nil {
		return "", err
	}

	fingerprint.Request = request
	return createBlackbox(&fingerprint)
}

func NewEncryptedBlackbox(identity_manager IdentityManager, gs_id, account_id, installation string) ([]byte, error) {
	delim_index := strings.LastIndexByte(gs_id, '-')
	session := gs_id[:delim_index]

	feature := float64(rand.Intn(0xFFFFFFFE-1) + 1)
	request := Request{Features: []float64{feature}, Installation: installation, Session: session}

	blackbox, err := NewBlackbox(identity_manager, &request)

	if err != nil {
		return nil, err
	}

	return encryptBlackbox(blackbox, gs_id, account_id), nil
}
