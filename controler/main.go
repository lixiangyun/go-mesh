package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	LOG "github.com/lixiangyun/go-mesh/mesher/log"

	"github.com/lixiangyun/go-mesh/controler/etcd"
)

var (
	LISTEN_ADDR  string
	CFG_FILE     string
	LOG_FILE     string
	ETCD_CLUSTER string
	h            bool
)

func init() {
	flag.StringVar(&LISTEN_ADDR, "p", "127.0.0.1:301", "server listen address.")
	flag.StringVar(&CFG_FILE, "cfg", "config.json", "server proxy cfg file.")
	flag.StringVar(&LOG_FILE, "log", "controler.log", "the controler record log file.")
	flag.StringVar(&ETCD_CLUSTER, "etcd", "127.0.0.1:2379", "the etcd cluster address list.[IP1,IP2,IP3...]")
	flag.BoolVar(&h, "h", false, "this help.")
}

func main() {

	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	LOG.SetLogFile(LOG_FILE)

	err := ProxyCfgLoadFromFile(CFG_FILE)
	if err != nil {
		log.Printf("load cfg from(%s) failed! (%s)\r\n", CFG_FILE, err.Error())
		return
	}

	endpoints := strings.Split(ETCD_CLUSTER, ",")
	log.Println("connect etcd cluster : ", endpoints)

	err = etcd.ServiceConnect(endpoints)
	if err != nil {
		log.Printf("connect etcd cluster failed!(%s)\r\n", err.Error())
		return
	}

	ServerCheckTimeout()

	mux := http.NewServeMux()

	mux.HandleFunc("/proxy/cfg", ProxyCfgHandler)
	mux.HandleFunc("/server/query", ServerQueryHandler)
	mux.HandleFunc("/server/register", ServerRegisterHandler)

	http.ListenAndServe(LISTEN_ADDR, mux)
}
