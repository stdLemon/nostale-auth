package blackbox

import (
	"testing"
	"time"
)

func TestGenerateVector(t *testing.T) {
	vec := GenerateVector()
	content, _ := UnpackVector(vec)
	len := len(content)

	if len != VECTOR_CONTENT_LENGTH {
		t.Fatal("content length must be equal to", VECTOR_CONTENT_LENGTH, "got", len)
	}

	t.Log(vec)
}

func TestUpdateVector(t *testing.T) {
	vec := GenerateVector()
	content1, time1 := UnpackVector(vec)
	time.Sleep(1 * time.Millisecond)
	UpdateVector(&vec)
	content2, time2 := UnpackVector(vec)

	if len(content2) != VECTOR_CONTENT_LENGTH {
		t.Error("content length must be equal to", VECTOR_CONTENT_LENGTH, "got", len(vec))
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
	uuid := GenerateUuid()
	len := len(uuid)

	if len != UUID_LENGTH {
		t.Error("length must be equal to", UUID_LENGTH, "got", len)
	}

	for _, c := range uuid {
		if ((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) == false {
			t.Fatal("only lowercase alphanumeric characters are allowed")
		}
	}

	t.Log(uuid)
}

func TestGetServerDate(t *testing.T) {
	s, err := GetServerDate()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(s)
}
