package blackbox

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/fatih/structs"
)

type Blackbox string

func New(fingerprint *Fingerprint) (Blackbox, error) {
	var (
		fields = structs.Fields(fingerprint)
		values = make([]interface{}, len(fields))
	)

	for i := range fields {
		values[i] = fields[i].Value()
	}

	valuesJson, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	encodedValues := url.QueryEscape(string(valuesJson))
	encodedValues = strings.NewReplacer("+", "%20", "%29", ")", "%28", "(").Replace(encodedValues)
	var (
		encodedLen = len(encodedValues)
		blackbox   = make([]byte, encodedLen)
	)

	blackbox[0] = encodedValues[0]
	for i := 1; i < encodedLen; i++ {
		blackbox[i] = blackbox[i-1] + encodedValues[i]
	}

	return Blackbox("tra:" + base64.RawURLEncoding.EncodeToString(blackbox)), nil
}

func (b Blackbox) Encrypt(gsId, accountId string) []byte {
	var (
		encrypted = xor([]byte(b), createKey(gsId, accountId))
		encoded   = make([]byte, base64.StdEncoding.EncodedLen(len(encrypted)))
	)

	base64.StdEncoding.Encode(encoded, encrypted)
	return encoded
}

func (b Blackbox) String() string {
	return string(b)
}
