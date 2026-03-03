package gfclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type GfAccountData struct {
	Email    string
	Password string
	Locale   string
	Name     string
}

func TestCalcCefUserAgentChecksum(t *testing.T) {
	tests := []struct {
		name             string
		installationId   string
		expectedChecksum string
	}{
		{
			name:             "odd installation id",
			installationId:   "c37f161c-7201-48a5-9a27-89c20a2a243d",
			expectedChecksum: "14fef9b2",
		},
		{
			name:             "even installation id",
			installationId:   "c47f161c-7201-48a5-9a27-89c20a2a243d",
			expectedChecksum: "74646fd4",
		},
	}

	const (
		accountId     = "f3779a30-fe2e-4eb9-b757-24897f4cd7d1"
		clientVersion = "2.8.0.1876"
	)

	c := Client{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.installationId = tt.installationId
			h := c.calcCefUserAgentChecksum(accountId, clientVersion)
			assert.Equal(t, tt.expectedChecksum, h)
		})
	}
}
