package gfclient

import (
	"net/http"

	"github.com/stdLemon/nostale-auth/pkg/restclient"
)

type Client struct {
	gfUserAgent    string
	cefUserAgents  map[string]string
	installationId string
	restClient     restclient.Client
}

func New(gfUserAgent, installationId string) *Client {
	return &Client{
		gfUserAgent:    gfUserAgent,
		cefUserAgents:  make(map[string]string),
		installationId: installationId,
		restClient: restclient.New(http.Header{
			"tnt-installation-id": {installationId},
			"User-Agent":          {gfUserAgent},
		}),
	}
}
