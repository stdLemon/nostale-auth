package gfclient

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strings"
)

const clientVersionEndpoint = "http://dl.tnt.gameforge.com/tnt/final-ms3/clientversioninfo.json"

type clientVersionInfo struct {
	Version                        string                `json:"version"`
	MinimumVersionForDelayedUpdate string                `json:"minimumVersionForDelayedUpdate"`
	Force                          bool                  `json:"force"`
	Full                           clientVersionFull     `json:"full"`
	Patch                          clientVersionPatch    `json:"patch"`
	GSL                            clientVersionGSL      `json:"gsl"`
	Frontend                       clientVersionFrontend `json:"frontend"`
}

type clientVersionFull struct {
	URL  string `json:"url"`
	Hash string `json:"hash"`
}

type clientVersionPatch struct {
	URL        string `json:"url"`
	Hash       string `json:"hash"`
	MinVersion string `json:"minVersion"`
}

type clientVersionGSL struct {
	URL       string `json:"url"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type clientVersionFrontend struct {
	Hash                       string `json:"hash"`
	URL                        string `json:"url"`
	Version                    string `json:"version"`
	MinVersionForDelayedUpdate string `json:"minVersionForDelayedUpdate"`
}

func (c *Client) getCefUserAgent(accountID string) (string, error) {
	if ua, ok := c.cefUserAgents[accountID]; ok {
		return ua, nil
	}

	if c.cefUserAgents == nil {
		c.cefUserAgents = make(map[string]string)
	}

	version, err := c.getVersion()
	if err != nil {
		return "", err
	}

	ua := c.createCefUserAgent(accountID, version)

	c.cefUserAgents[accountID] = ua
	return ua, nil
}

func (c *Client) calcCefUserAgentChecksum(accountID, clientVersion string) string {
	const (
		certSHA256 = "99025da70af1ef39d2acd049018887ef5140daebc6f11d80461bcf8d02f2d36b"
		certSHA1   = "d68f9401f15791cc396d4d6af3b977bc58ad0002"
	)

	hashChain := func(hashers []hash.Hash, certHash string) string {
		var (
			sb     strings.Builder
			inputs = []string{"C" + clientVersion, c.installationId, accountID}
		)

		sb.WriteString(certHash)
		for i, data := range inputs {
			h := hashers[i]
			h.Write([]byte(data))
			e := hex.EncodeToString(h.Sum(nil))
			sb.WriteString(e)
		}

		return fmt.Sprintf("%x", sha256.Sum256([]byte(sb.String())))
	}

	var firstDigit byte
	for _, r := range c.installationId {
		if r >= '0' && r <= '9' {
			firstDigit = byte(r)
			break
		}
	}

	var (
		evenHashers = []hash.Hash{sha1.New(), sha256.New(), sha1.New()}
		oddHashers  = []hash.Hash{sha256.New(), sha1.New(), sha256.New()}
	)

	if (firstDigit-'0')%2 == 0 {
		h := hashChain(evenHashers, certSHA256)
		return strings.Clone(h[:8])
	}

	h := hashChain(oddHashers, certSHA1)
	return strings.Clone(h[len(h)-8:])
}

func (c *Client) createCefUserAgent(accountID, clientVersion string) string {
	h := c.calcCefUserAgentChecksum(accountID, clientVersion)
	return fmt.Sprintf("Chrome/C%s (%s%s)", clientVersion, accountID[:2], h)
}

func (c Client) getVersion() (string, error) {
	resp := new(clientVersionInfo)
	_, err := c.restClient.NewRequest().
		SetResult(resp).
		Get(clientVersionEndpoint)
	if err != nil {
		return "", err
	}

	if resp.Version == "" {
		return "", errors.New("server didn't send a client version")
	}

	return resp.Version, nil
}
