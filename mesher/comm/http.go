package comm

import (
	"net"
	"net/http"
	"time"
)

var defaultTransport http.RoundTripper = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConnsPerHost:   10,
	MaxIdleConns:          10,
	DisableKeepAlives:     true,
	IdleConnTimeout:       30 * time.Second,
	TLSHandshakeTimeout:   30 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var HttpClient *http.Client

// init HTTPClient
func init() {
	HttpClient = &http.Client{
		Transport: defaultTransport,
		Timeout:   10 * time.Second,
	}
}
