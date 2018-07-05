package main

import (
	"flag"
	"log"
	"net"
	"sync"
	"time"

	"github.com/lixiangyun/go-mesh/mesher/stat"
)

var (
	SERVER_NAME    string
	SERVER_VERSION string
	LISTEN_ADDRESS string

	h bool
)

var gStat *stat.Stat

func init() {
	flag.StringVar(&SERVER_NAME, "n", "demotcp", "set the service name.")
	flag.StringVar(&SERVER_VERSION, "v", "1.0.0", "set the service version.")
	flag.StringVar(&LISTEN_ADDRESS, "p", "127.0.0.1:10001", "set the service listen addr.")

	flag.BoolVar(&h, "h", false, "this help.")

	gStat = stat.NewStat(10)
}

func process(conn net.Conn) {

	var recvbuf [65535]byte
	var sendbuf [65535]byte

	log.Printf("new connect remoteaddr %s, localaddr %s\r\n",
		conn.RemoteAddr().String(),
		conn.LocalAddr().String())

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
			gStat.Recv(cnt)
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
			gStat.Send(cnt)
		}
	}()

	wait.Wait()
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

	var delaytime time.Duration

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err.Error())

			if delaytime == 0 {
				delaytime = 5 * time.Millisecond
			} else if delaytime < 1*time.Second {
				delaytime = delaytime * 2
			}
			time.Sleep(delaytime)

			continue
		}
		go process(conn)
	}
}
