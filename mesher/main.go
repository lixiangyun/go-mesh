package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lixiangyun/go-mesh/mesher/log"
)

const (
	BIN_NAME = "mesher"
	BIN_VER  = "0.1.0"
)

var (
	SERVER_NAME    string
	SERVER_VERSION string
	CONTROLER_ADDR string
	NETWORK_IP     string
	LOG_FILE       string

	h bool
)

func init() {
	flag.StringVar(&SERVER_NAME, "name", "", "set the service name for mesher proxy.")
	flag.StringVar(&SERVER_VERSION, "ver", "", "set the service version for mesher proxy.")
	flag.StringVar(&CONTROLER_ADDR, "control", "127.0.0.1:301", "set the mesher control service addr.")
	flag.StringVar(&NETWORK_IP, "bind", "", "set the mesher bind network addr.")
	flag.StringVar(&LOG_FILE, "log", "", "the controler record log file.")

	flag.BoolVar(&h, "h", false, "this help.")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `mesher version: `+BIN_VER+`

Usage: mesher [-h] [-name servicename] [-ver serviceversion] [-control ip:port] [-bind ip1,ip2,ip3...]

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

	if LOG_FILE == "" {
		LOG_FILE = fmt.Sprintf("mesher_%s_%s.log", SERVER_NAME, SERVER_VERSION)
	}
	log.SetLogFile(LOG_FILE)

	SetNetWork(NETWORK_IP)
	MesherStart(SERVER_NAME, SERVER_VERSION, CONTROLER_ADDR)
}
