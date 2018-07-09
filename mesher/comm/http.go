package comm

import (
	"net"
	"net/http"
	"time"
)

func newTransport() http.RoundTripper {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConnsPerHost:   10,
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func NewHttpClient() *http.Client {
	return &http.Client{
		Transport: newTransport(),
		Timeout:   10 * time.Second,
	}
}
