package gfclient_auth

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

func xor(data []byte, key []byte) []byte {
	key_length := len(key)
	result := make([]byte, len(data))

	for i := range data {
		wrapping_i := i % key_length
		result[i] = data[i] ^ key[wrapping_i] ^ key[key_length-wrapping_i-1]
	}

	return result
}

func createKey(gs_id, account_id string) []byte {
	v := fmt.Sprintf("%s-%s", gs_id, account_id)
	hash := sha512.Sum512([]byte(v))

	return []byte(fmt.Sprintf("%x", hash))
}

func encryptBlackbox(blackbox string, gs_id, account_id string) []byte {
	encrypted := xor([]byte(blackbox), createKey(gs_id, account_id))
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(encrypted)))

	base64.StdEncoding.Encode(encoded, encrypted)
	return encoded
}
