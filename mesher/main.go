package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	BIN_NAME = "mesher"
	BIN_VER  = "0.0.1"
)

var (
	SERVER_NAME    string
	SERVER_VERSION string
	CONTROLER_ADDR string
	NETWORK_IP     string

	h bool
)

func init() {
	flag.StringVar(&SERVER_NAME, "n", "demo", "set the service name for mesher proxy.")
	flag.StringVar(&SERVER_VERSION, "v", "1.1.1", "set the service version for mesher proxy.")
	flag.StringVar(&CONTROLER_ADDR, "c", "127.0.0.1:301", "set the mesher control service addr.")
	flag.StringVar(&NETWORK_IP, "b", "", "set the mesher bind network addr.")

	flag.BoolVar(&h, "h", false, "this help.")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `mesher version: mesher/0.0.1
Usage: mesher [-h] [-n servicename] [-v serviceversion] [-c ip:port] [-b ip1,ip2,ip3...]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	log.SetPrefix(fmt.Sprintf("[%s %s] ", BIN_NAME, BIN_VER))
	SetNetWork(NETWORK_IP)
	MesherStart(SERVER_NAME, SERVER_VERSION, CONTROLER_ADDR)
}
