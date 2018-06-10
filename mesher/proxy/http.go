package proxy

import (
	"net/http"
)

type HttpProxy struct {
}

// HTTPHandler handles http connections.
func (proxy *HttpProxy) HTTPHandler(rw http.ResponseWriter, req *http.Request) {

	log.Infof("%v is sending request %v %v \n", proxy.User, req.Method, req.URL.Host)
	RmProxyHeaders(req)

	resp, err := proxy.Tr.RoundTrip(req)
	if err != nil {
		log.Errorf("%v", err)
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	ClearHeaders(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode) //写入响应状态

	nr, err := ioCopy(rw, resp.Body)
	if err != nil && err != io.EOF {
		log.Errorf("%v got an error when copy remote response to client. %v\n", proxy.User, err)
		return
	}
	log.Infof("%v copied %v bytes from %v.\n", proxy.User, nr, req.URL.Host)
}

// CopyHeaders copy headers from source to destination.
// Nothing would be returned.
func CopyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

// ClearHeaders clear headers.
func ClearHeaders(headers http.Header) {
	for key := range headers {
		headers.Del(key)
	}
}

// RmProxyHeaders remove Hop-by-hop headers.
func RmProxyHeaders(req *http.Request) {
	req.RequestURI = ""
	req.Header.Del("Proxy-Connection")
	req.Header.Del("Connection")
	req.Header.Del("Keep-Alive")
	req.Header.Del("Proxy-Authenticate")
	req.Header.Del("Proxy-Authorization")
	req.Header.Del("TE")
	req.Header.Del("Trailers")
	req.Header.Del("Transfer-Encoding")
	req.Header.Del("Upgrade")
}

func (*proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	RmProxyHeaders(req)

	resp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Errorf("%v", err)
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	// step 1
	outReq := new(http.Request)
	*outReq = *req // this only does shallow copies of maps

	// step 3
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func NewHttpProxy(addr string) error {

	http.ListenAndServe(addr)

	http.ListenAndServe(*addr, proxy)

}

func HttpFilter(req *http.Request) {

}
