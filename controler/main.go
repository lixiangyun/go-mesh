package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	LISTEN_ADDR string
	CFG_FILE    string
	h           bool
)

func init() {
	flag.StringVar(&LISTEN_ADDR, "p", "127.0.0.1:301", "server listen address.")
	flag.StringVar(&CFG_FILE, "cfg", "config.json", "server proxy cfg file.")
	flag.BoolVar(&h, "h", false, "this help.")
}

func main() {

	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	err := ProxyCfgLoadFromFile(CFG_FILE)
	if err != nil {
		log.Printf("load cfg from(%s) failed! (%s)\r\n", CFG_FILE, err.Error())
		return
	}

	ServerCheckTimeout()

	mux := http.NewServeMux()

	mux.HandleFunc("/proxy/cfg", ProxyCfgHandler)
	mux.HandleFunc("/server/query", ServerQueryHandler)
	mux.HandleFunc("/server/register", ServerRegisterHandler)

	http.ListenAndServe(LISTEN_ADDR, mux)
}
