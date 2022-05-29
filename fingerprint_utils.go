package gfclient_auth

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
	builder := strings.Builder{}
	for i := uint(0); i < n; i++ {
		builder.WriteByte(randomAscii())
	}

	return builder.String()
}

func getServerDate() (string, error) {
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

func createVector(content string, time int64) string {
	return fmt.Sprintf("%v %v", content, time)
}

func generateVector() string {
	content := randomString(VECTOR_CONTENT_LENGTH)
	time := time.Now().UnixMilli()

	return createVector(content, time)
}

func unpackVector(vector string) (string, string) {
	delim_index := strings.LastIndexByte(vector, ' ')
	content := vector[:delim_index]
	time := vector[delim_index:]
	return content, time
}

func updateVector(vector *string) {
	content, _ := unpackVector(*vector)

	new_content := content[1:] + string(randomAscii())
	new_time := time.Now().UnixMilli()
	*vector = createVector(new_content, new_time)
}

func generateUuid() string {
	str := randomString(UUID_LENGTH)
	return strings.ToLower(base32.StdEncoding.EncodeToString([]byte(str)))[:UUID_LENGTH]
}
