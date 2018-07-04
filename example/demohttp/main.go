package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	SERVER_NAME    string
	SERVER_VERSION string
	LISTEN_ADDRESS string

	h bool
)

type DemoHttp struct{}

func (*DemoHttp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	body := fmt.Sprintf("[%s %s]Received request [%s %s %s]\n", SERVER_NAME, SERVER_VERSION, req.Method, req.Host, req.RemoteAddr)

	log.Println(body)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(body))
}

func init() {
	flag.StringVar(&SERVER_NAME, "n", "demohttp", "set the service name.")
	flag.StringVar(&SERVER_VERSION, "v", "1.0.0", "set the service version.")
	flag.StringVar(&LISTEN_ADDRESS, "p", "127.0.0.1:8001", "set the service listen addr.")

	flag.BoolVar(&h, "h", false, "this help.")
}

func main() {

	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	log.Printf("%s %s start success!\r\n", SERVER_NAME, SERVER_VERSION)

	err := http.ListenAndServe(LISTEN_ADDRESS, &DemoHttp{})
	if err != nil {
		log.Println(err.Error())
	}
}
