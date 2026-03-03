package gfclient

import (
	"net/http"
)

const (
	_headerAuthorization  = "Authorization"
	_headerOrigin         = "Origin"
	_originSparkGameforge = "spark://www.gameforge.com"
)

func setAuthorizationHeader(headers http.Header, bearer string) {
	headers.Set(_headerAuthorization, "Bearer "+bearer)
}

func setOriginHeader(headers http.Header) {
	headers.Set(_headerOrigin, _originSparkGameforge)
}
