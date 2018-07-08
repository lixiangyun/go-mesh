package comm

import (
	"net"
	"net/http"
	"sync"
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
	IdleConnTimeout:       30 * time.Second,
	TLSHandshakeTimeout:   30 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var gLock sync.Mutex

var gHttpClientMap map[string]*http.Client

func init() {
	gHttpClientMap = make(map[string]*http.Client, 100)
}

func newhttpclient() *http.Client {
	return &http.Client{
		Transport: defaultTransport,
		Timeout:   10 * time.Second,
	}
}

func HttpClient(addr string) *http.Client {
	gLock.Lock()
	client, b := gHttpClientMap[addr]
	if b == false {
		client = newhttpclient()
		gHttpClientMap[addr] = client
	}
	defer gLock.Unlock()

	return client
}
