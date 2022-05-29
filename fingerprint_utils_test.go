package gfclient_auth

import (
	"testing"
	"time"
)

func TestGenerateVector(t *testing.T) {
	vector := generateVector()
	content, _ := unpackVector(vector)
	content_len := len(content)

	if content_len != VECTOR_CONTENT_LENGTH {
		t.Fatal("content length must be equal to", VECTOR_CONTENT_LENGTH, "got", content_len)
	}

	t.Log(vector)
}

func TestUpdateVector(t *testing.T) {
	vector := generateVector()
	content1, time1 := unpackVector(vector)
	time.Sleep(1 * time.Millisecond)
	updateVector(&vector)
	content2, time2 := unpackVector(vector)

	if len(content2) != VECTOR_CONTENT_LENGTH {
		t.Error("content length must be equal to", VECTOR_CONTENT_LENGTH, "got", len(vector))
	}

	if time1 == time2 {
		t.Error("time was not updated", time1, time2)
	}

	if content1[1] != content2[0] {
		t.Error("first character was not shifted")
	}

	if content1[len(content1)-1] == content2[len(content2)-1] {
		t.Error("last character was not updated")
	}
}

func TestGenerateUuid(t *testing.T) {
	uuid := generateUuid()
	uuid_len := len(uuid)

	if uuid_len != UUID_LENGTH {
		t.Error("length must be equal to", UUID_LENGTH, "got", uuid_len)
	}

	for _, c := range uuid {
		if ((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) == false {
			t.Fatal("only lowercase alphanumeric characters are allowed")
		}
	}

	t.Log(uuid)
}

func TestGetServerDate(t *testing.T) {
	s, err := getServerDate()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(s)
}
