package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	LISTEN_ADDRESS = "127.0.0.1:8001"
)

type DemoHttp struct{}

func (*DemoHttp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	body := fmt.Sprintf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	log.Println(body)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(body))
}

func main() {
	err := http.ListenAndServe(LISTEN_ADDRESS, &DemoHttp{})
	if err != nil {
		log.Println(err.Error())
	}
}
