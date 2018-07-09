package api

import (
	"net/http"
	"sync"

	"github.com/lixiangyun/go-mesh/mesher/comm"
)

var gHttpClient *http.Client

var gLock sync.Mutex

func init() {
	gHttpClient = comm.NewHttpClient()
}

func HttpClient(req *http.Request) (*http.Response, error) {
	gLock.Lock()
	defer gLock.Unlock()

	return gHttpClient.Do(req)
}
