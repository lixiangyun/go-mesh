package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	SERVER_NAME    string
	SERVER_VERSION string

	h bool
)

func init() {
	flag.StringVar(&SERVER_NAME, "n", "demo", "set the service name.")
	flag.StringVar(&SERVER_VERSION, "v", "1.1.1", "set the service version.")

	flag.BoolVar(&h, "h", false, "this help.")
}

func TcpSend(addr string, body string) string {

	var recvbuf [1024]byte

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	defer conn.Close()

	_, err = conn.Write([]byte(body))
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	cnt, err := conn.Read(recvbuf[:])
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	return string(recvbuf[:cnt])
}

func HttpRequest(addr string, url string, body string) string {

	path := "http://" + addr + url

	transport := http.DefaultTransport

	request, err := http.NewRequest("GET", path, strings.NewReader(body))
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	rsp, err := transport.RoundTrip(request)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	body2, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	return string(body2)
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	log.Printf("%s %s start success!\r\n", SERVER_NAME, SERVER_VERSION)

	var cnt int

	for {
		req := fmt.Sprintf("hello world! (tcp:%d)", cnt)
		rsp := TcpSend("127.0.0.1:1000", req)

		log.Printf("TCP_TEST:\r\nREQ:%s\r\nRSP:%s\r\n", req, rsp)

		req = fmt.Sprintf("hello world! (http:%d)", cnt)
		rsp = HttpRequest("127.0.0.1:2000", "/abc", req)
		log.Printf("HTTP_TEST:\r\nREQ:%s\r\nRSP:%s\r\n", req, rsp)

		time.Sleep(1 * time.Second)
	}
}
