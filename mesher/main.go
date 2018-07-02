package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	SERVER_NAME    string
	SERVER_VERSION string
	CONTROLER_ADDR string

	h bool
)

func init() {
	flag.StringVar(&SERVER_NAME, "n", "demo", "set the service name for mesher proxy.")
	flag.StringVar(&SERVER_VERSION, "v", "v0.1.0", "set the service version for mesher proxy.")
	flag.StringVar(&CONTROLER_ADDR, "c", "127.0.0.1:3001", "set the mesher control service addr.")

	flag.BoolVar(&h, "h", false, "this help.")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `mesher version: mesher/0.0.1
Usage: mesher [-h] [-n servicename] [-v serviceversion] [-c ip:port]

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

	MesherStart(SERVER_NAME, SERVER_VERSION, CONTROLER_ADDR)
}
