package blackbox

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBlackbox(t *testing.T) {
	for _, filename := range []string{"test/fingerprint1.json", "test/fingerprint2.json"} {
		content, err := os.ReadFile(filename)
		require.NoError(t, err)

		finterprint := new(Fingerprint)
		require.NoError(t, json.Unmarshal(content, finterprint))

		blackbox, err := New(finterprint)
		require.NoError(t, err)

		var expectedBlackbox bytes.Buffer
		cmd := exec.Command("gfclient_poc/create_blackbox.js", filename)
		cmd.Stdout = &expectedBlackbox
		require.NoError(t, cmd.Run())

		assert.Equal(t, expectedBlackbox.String(), blackbox.String())
	}
}
