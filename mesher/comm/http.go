package comm

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type httpClientMap struct {
	sync.RWMutex
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

func newhttpclient(addr string) *http.Client {
	client, b := gHttpClientAll.clients[addr]
	if b == false {
		client = &http.Client{
			Transport: newTransport(),
			Timeout:   10 * time.Second,
		}
		gHttpClientAll.clients[addr] = client
	}
	return client
}

func HttpClient(addr string) *http.Client {
	gHttpClientAll.RLock()
	client, b := gHttpClientAll.clients[addr]
	gHttpClientAll.RUnlock()

	if b == false {
		gHttpClientAll.Lock()
		client = newhttpclient(addr)
		gHttpClientAll.Unlock()
	}

	return client
}
