package blackbox

import (
	"encoding/base32"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	resp, err := http.Get(SERVER_FILE_GAME1_FILE)
	if err != nil {
		return "", err
	}

	date, err := time.Parse(time.RFC1123, resp.Header["Date"][0])
	if err != nil {
		return "", err
	}

	return date.Format(time.RFC3339), nil
}

func CreateVector(content string, time int64) string {
	return fmt.Sprintf("%v %v", content, time)
}

func GenerateVector() string {
	content := randomString(VECTOR_CONTENT_LENGTH)
	time := time.Now().UnixMilli()

	return CreateVector(content, time)
}

func UnpackVector(vector string) (string, string) {
	i := strings.LastIndexByte(vector, ' ')
	content := vector[:i]
	time := vector[i:]
	return content, time
}

func UpdateVector(vector *string) {
	content, _ := UnpackVector(*vector)

	nContent := content[1:] + string(randomAscii())
	nTime := time.Now().UnixMilli()
	*vector = CreateVector(nContent, nTime)
}

func GenerateUuid() string {
	str := randomString(UUID_LENGTH)
	return strings.ToLower(base32.StdEncoding.EncodeToString([]byte(str)))[:UUID_LENGTH]
}
