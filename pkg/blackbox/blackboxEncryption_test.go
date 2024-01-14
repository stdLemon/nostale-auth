package blackbox

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const RANDOM_TEXT_FILENAME = "test/random_text.txt"

func TestEncryptBlackbox(t *testing.T) {
	content, err := os.ReadFile(RANDOM_TEXT_FILENAME)
	require.NoError(t, err)

	var (
		expectedEncryptedBlackbox bytes.Buffer
		blackbox                  = Blackbox(content)
		gsId                      = uuid.New().String()
		accountId                 = uuid.New().String()
		encryptedBlackbox         = blackbox.Encrypt(gsId, accountId)
		cmd                       = exec.Command("gfclient_poc/encrypt_blackbox.js", RANDOM_TEXT_FILENAME, gsId, accountId)
	)

	cmd.Stdout = &expectedEncryptedBlackbox
	err = cmd.Run()
	require.NoError(t, err)

	assert.Equal(t, expectedEncryptedBlackbox.Bytes(), encryptedBlackbox)
}
