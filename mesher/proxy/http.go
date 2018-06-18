package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

type HttpProxy struct {
	Server string
	Addr   string
	Svc    *http.Server
}

func (h *HttpProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	//fmt.Println(req.URL.Path)

	// step 1
	req.Host = h.Server
	req.RequestURI = "http://" + h.Server + "/" + req.URL.Path

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
	resp.Body.Close()
}

func NewHttpProcy(addr string, servername string) *HttpProxy {
	proxy := new(HttpProxy)

	proxy.Addr = addr
	proxy.Server = servername

	lis, err := net.Listen("tcp", proxy.Addr)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	proxy.Svc = &http.Server{Handler: proxy}

	go proxy.Svc.Serve(lis)

	return proxy
}

func (h *HttpProxy) Close() {
	h.Svc.Close()
}

func HttpFilter(req *http.Request) {

}