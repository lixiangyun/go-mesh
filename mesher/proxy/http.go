package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

type HttpProxy struct {
	Fun  SELECT_ADDR
	Addr string
	Svc  *http.Server
}

func (h *HttpProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	redirect := h.Fun()

	// step 1
	req.Host = redirect
	req.RequestURI = "http://" + redirect + "/" + req.URL.Path

	req.URL, _ = url.Parse(req.RequestURI)

	// step 2
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {

		fmt.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}
	defer resp.Body.Close()

	// step 3
	for key, value := range resp.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

func NewHttpProxy(addr string, fun SELECT_ADDR) *HttpProxy {
	proxy := new(HttpProxy)

	proxy.Addr = addr
	proxy.Fun = fun

	lis, err := net.Listen("tcp", proxy.Addr)
	if err != nil {
		log.Println("http listen failed!", err.Error())
		return nil
	}

	log.Printf("Http Proxy Listen : %s\r\n", addr)

	proxy.Svc = &http.Server{Handler: proxy}

	go proxy.Svc.Serve(lis)

	return proxy
}

func (h *HttpProxy) Close() {
	h.Svc.Close()
}
