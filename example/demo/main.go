package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/lixiangyun/go-mesh/mesher/api"
	"github.com/lixiangyun/go-mesh/mesher/stat"
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

	conn.SetDeadline(time.Now().Add(1 * time.Second))
	_, err = conn.Write([]byte(body))
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	conn.SetDeadline(time.Now().Add(1 * time.Second))
	cnt, err := conn.Read(recvbuf[:])
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	return string(recvbuf[:cnt])
}

func TcpBenchMark(addr string, duration int) {
	var recvbuf [65535]byte
	var sendbuf [65535]byte

	log.Println("start tcp bench mark test....")

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	s := stat.NewStat(5)
	defer s.Delete()

	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		defer wait.Done()
		defer conn.Close()
		for {
			cnt, err := conn.Read(recvbuf[:])
			if err != nil {
				log.Println(err.Error())
				return
			}
			s.Recv(cnt)
		}
	}()

	go func() {
		defer wait.Done()
		defer conn.Close()
		for {
			cnt, err := conn.Write(sendbuf[:])
			if err != nil {
				log.Println(err.Error())
				return
			}
			s.Send(cnt)
		}
	}()

	stop := make(chan struct{}, 1)
	go func() {
		wait.Wait()
		stop <- struct{}{}
	}()

	go func() {
		time.Sleep(time.Duration(duration) * time.Second)
		conn.Close()
		wait.Wait()
		stop <- struct{}{}
	}()

	<-stop
}

func HttpBenchMark(addr string, duration int) {

	log.Println("start http bench mark test....")

	s := stat.NewStat(5)
	defer s.Delete()

	var stop bool

	var wait sync.WaitGroup
	wait.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wait.Done()

			for {
				if stop {
					return
				}

				req := []byte("helloworld!")
				s.Send(len(req))

				rsp, err := HttpRequest(addr, "/123", req)
				if err != nil {
					log.Println(err.Error())
					return
				}
				s.Recv(len(rsp))
			}
		}()
	}

	time.Sleep(time.Duration(duration) * time.Second)

	stop = true

	wait.Wait()
}

func HttpRequest(addr string, url string, body []byte) ([]byte, error) {

	path := "http://" + addr + url

	request, err := http.NewRequest("GET", path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	rsp, err := api.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	log.Printf("%s %s start success!\r\n", SERVER_NAME, SERVER_VERSION)

	for {
		time.Sleep(5 * time.Second)
		TcpBenchMark("127.0.0.1:1000", 60)

		time.Sleep(5 * time.Second)
		HttpBenchMark("127.0.0.1:2000", 60)
	}

	/*
		var cnt int

		for {
			req := fmt.Sprintf("hello world! (tcp:%d)", cnt)
			rsp := TcpSend("127.0.0.1:1000", req)

			log.Printf("TCP_TEST:\r\nREQ:%s\r\nRSP:%s\r\n", req, rsp)

			req = fmt.Sprintf("hello world! (http:%d)", cnt)
			rsp = HttpRequest("127.0.0.1:2000", "/abc", req)
			log.Printf("HTTP_TEST:\r\nREQ:%s\r\nRSP:%s\r\n", req, rsp)

			time.Sleep(1 * time.Second)
		}*/
}
