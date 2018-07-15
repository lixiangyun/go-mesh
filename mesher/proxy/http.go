package proxy

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/lixiangyun/go-mesh/mesher/comm"
	"github.com/lixiangyun/go-mesh/mesher/log"
)

type HttpProxy struct {
	Fun  SELECT_ADDR
	Addr string
	Svc  *http.Server

	GoCnt int
	Que   chan *HttpRequest
	Wait  sync.WaitGroup
	Stop  chan struct{}
}

type HttpRsponse struct {
	status int
	header http.Header
	body   []byte

	err error
}

type HttpRequest struct {
	addr   string
	url    string
	method string
	header http.Header
	body   []byte
	rsp    chan *HttpRsponse
}

func (h *HttpProxy) Process() {
	defer h.Wait.Done()

	httpclient := comm.NewHttpClient()

	for {
		select {
		case proxyreq := <-h.Que:
			{
				proxyrsp := new(HttpRsponse)

				request, err := http.NewRequest(proxyreq.method,
					proxyreq.url,
					bytes.NewBuffer(proxyreq.body))
				if err != nil {
					proxyrsp.err = err
					proxyrsp.status = http.StatusInternalServerError

					proxyreq.rsp <- proxyrsp
					continue
				}

				for key, value := range proxyreq.header {
					for _, v := range value {
						request.Header.Add(key, v)
					}
				}

				resp, err := httpclient.Do(request)
				if err != nil {
					proxyrsp.err = err
					proxyrsp.status = http.StatusInternalServerError

					proxyreq.rsp <- proxyrsp
					continue
				} else {
					proxyrsp.status = resp.StatusCode
					proxyrsp.header = resp.Header
				}

				proxyrsp.body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					proxyrsp.err = err
					proxyrsp.status = http.StatusInternalServerError

					proxyreq.rsp <- proxyrsp
					continue
				}
				resp.Body.Close()

				proxyreq.rsp <- proxyrsp
			}
		case <-h.Stop:
			{
				return
			}
		}
	}
}

func (h *HttpProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	var err error

	defer req.Body.Close()

	redirect := h.Fun()

	// step 1
	proxyreq := new(HttpRequest)
	proxyreq.addr = redirect
	proxyreq.url = "http://" + redirect + "/" + req.URL.Path
	proxyreq.method = req.Method
	proxyreq.header = req.Header
	proxyreq.rsp = make(chan *HttpRsponse, 1)

	proxyreq.body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(log.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Que <- proxyreq
	proxyrsp := <-proxyreq.rsp

	// step 2
	if proxyrsp.err != nil {
		log.Println(log.ERROR, proxyrsp.err.Error())
		http.Error(rw, proxyrsp.err.Error(), http.StatusInternalServerError)
		return
	}

	// step 3
	for key, value := range proxyrsp.header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(proxyrsp.status)
	rw.Write(proxyrsp.body)
}

func NewHttpProxy(addr string, fun SELECT_ADDR) *HttpProxy {
	proxy := new(HttpProxy)

	proxy.Addr = addr
	proxy.Fun = fun

	lis, err := net.Listen("tcp", proxy.Addr)
	if err != nil {
		log.Println(log.ERROR, "http listen failed!", err.Error())
		return nil
	}

	log.Printf(log.INFO, "Http Proxy Listen %s\r\n", addr)

	proxy.Svc = &http.Server{Handler: proxy}

	proxy.GoCnt = 10
	proxy.Que = make(chan *HttpRequest, 100)
	proxy.Stop = make(chan struct{}, proxy.GoCnt)

	proxy.Wait.Add(proxy.GoCnt)
	for i := 0; i < proxy.GoCnt; i++ {
		go proxy.Process()
	}

	go proxy.Svc.Serve(lis)

	return proxy
}

func (h *HttpProxy) Close() {

	log.Println(log.INFO, "Http Proxy Shut Down!", h.Addr)

	h.Svc.Close()
	for i := 0; i < h.GoCnt; i++ {
		h.Stop <- struct{}{}
	}
	h.Wait.Wait()
}
