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
	fields := structs.Fields(fingerprint)
	values := make([]interface{}, len(fields))
	for i := range fields {
		values[i] = fields[i].Value()
	}

	j, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	uri := url.QueryEscape(string(j))
	r := strings.NewReplacer("+", "%20", "%29", ")", "%28", "(")
	uri = r.Replace(uri)

	blackbox := make([]byte, len(uri))
	blackbox[0] = uri[0]
	for i := 1; i < len(uri); i++ {
		blackbox[i] = blackbox[i-1] + uri[i]
	}

	return Blackbox("tra:" + base64.RawURLEncoding.EncodeToString(blackbox)), nil
}

func (b Blackbox) Encrypt(gsId, accountId string) []byte {
	encrypted := xor([]byte(b), createKey(gsId, accountId))
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(encrypted)))

	base64.StdEncoding.Encode(encoded, encrypted)
	return encoded
}

func (b Blackbox) String() string {
	return string(b)
}
