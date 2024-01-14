package blackbox

import (
	"encoding/base32"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func randomAscii() byte {
	s := rand.Intn('~'-' ') + ' '
	return byte(s)
}

func randomString(n uint) string {
	b := strings.Builder{}
	for i := uint(0); i < n; i++ {
		b.WriteByte(randomAscii())
	}

	return b.String()
}

func GetServerDate() (string, error) {
	resp, err := http.Get(_gfAuthScriptUrl)
	if err != nil {
		return "", err
	}

	date, err := time.Parse(time.RFC1123, resp.Header["Date"][0])
	if err != nil {
		return "", err
	}

	return date.Format(time.RFC3339), nil
}

func CreateVector(content string, at time.Time) string {
	return fmt.Sprintf("%v %v", content, at.UnixMilli())
}

func RandomVector() string {
	return CreateVector(randomString(_vectorContentLength), time.Now())
}

func UnpackVector(vector string) (string, string) {
	i := strings.LastIndexByte(vector, ' ')
	return vector[:i], vector[i+1:]
}

func UpdateVector(vector string) string {
	return updateVector(vector, randomAscii(), time.Now())
}

func updateVector(vector string, newCharacter byte, at time.Time) string {
	var (
		content, _ = UnpackVector(vector)
		newContent = content[1:] + string(newCharacter)
	)
	return CreateVector(newContent, at)
}

func GenerateUuid() string {
	str := randomString(_uuidLength)
	return strings.ToLower(base32.StdEncoding.EncodeToString([]byte(str)))[:_uuidLength]
}
