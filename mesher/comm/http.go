package comm

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type httpClientMap struct {
	sync.Mutex
	clients map[string]*http.Client
}

var gHttpClientAll httpClientMap

func init() {
	gHttpClientAll.clients = make(map[string]*http.Client, 100)
}

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

func newhttpclient() *http.Client {
	return &http.Client{
		Transport: newTransport(),
		Timeout:   10 * time.Second,
	}
}

func HttpClient(addr string) *http.Client {
	gHttpClientAll.Lock()
	client, b := gHttpClientAll.clients[addr]
	if b == false {
		client = newhttpclient()
		gHttpClientAll.clients[addr] = client
	}
	gHttpClientAll.Unlock()

	return client
}
