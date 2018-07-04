package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	SERVER_NAME    string
	SERVER_VERSION string
	LISTEN_ADDRESS string

	h bool
)

func init() {
	flag.StringVar(&SERVER_NAME, "n", "demotcp", "set the service name.")
	flag.StringVar(&SERVER_VERSION, "v", "1.0.0", "set the service version.")
	flag.StringVar(&LISTEN_ADDRESS, "p", "127.0.0.1:10001", "set the service listen addr.")

	flag.BoolVar(&h, "h", false, "this help.")
}

func process(conn net.Conn) {
	var recvbuf [1024]byte

	defer conn.Close()

	for {
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		cnt, err := conn.Read(recvbuf[:])
		if err != nil {
			log.Println(err.Error())
			break
		}

		sendbuf := fmt.Sprintf("send from [%s %s] body [%v]\r\n",
			SERVER_NAME, SERVER_VERSION, recvbuf[:cnt])

		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		_, err = conn.Write([]byte(sendbuf))
		if err != nil {
			log.Println(err.Error())
			break
		}
	}
}

func main() {

	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	log.Printf("%s %s start success!\r\n", SERVER_NAME, SERVER_VERSION)

	lis, err := net.Listen("tcp", LISTEN_ADDRESS)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		go process(conn)
	}
}
