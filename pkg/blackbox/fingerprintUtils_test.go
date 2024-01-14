package blackbox

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVector(t *testing.T) {
	content, _ := UnpackVector(RandomVector())
	assert.Len(t, content, _vectorContentLength)
}

func TestUpdateVector(t *testing.T) {
	var (
		oldTime             = time.Now()
		expectedTime        = oldTime.Add(time.Millisecond)
		oldVec              = CreateVector("abc", oldTime)
		gotContent, gotTime = UnpackVector(updateVector(oldVec, 'x', expectedTime))
	)

	assert.Len(t, gotContent, 3)
	assert.Equal(t, fmt.Sprintf("%d", expectedTime.UnixMilli()), gotTime)
	assert.Equal(t, "bcx", gotContent)
}

func TestGenerateUuid(t *testing.T) {
	uuid := GenerateUuid()

	assert.Len(t, uuid, _uuidLength)
	for _, c := range uuid {
		if ((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) == false {
			t.Fatal("only lowercase alphanumeric characters are allowed")
		}
	}
}

func TestGetServerDate(t *testing.T) {
	date, err := GetServerDate()
	assert.NoError(t, err)
	assert.NotEmpty(t, date)
}
