package blackbox

import (
	"crypto/sha512"
	"fmt"
)

func xor(data []byte, key []byte) []byte {
	var (
		l      = len(key)
		result = make([]byte, len(data))
	)

	for i := range data {
		iMod := i % l
		result[i] = data[i] ^ key[iMod] ^ key[l-iMod-1]
	}

	return result
}

func createKey(gsId, accountId string) []byte {
	var (
		key     = fmt.Sprintf("%s-%s", gsId, accountId)
		keyHash = sha512.Sum512([]byte(key))
	)

	return []byte(fmt.Sprintf("%x", keyHash))
}
