package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lixiangyun/go-mesh/mesher/api"
)

var gProxyCfgMap map[api.SvcType]*api.ProxyCfg

func init() {
	gProxyCfgMap = make(map[api.SvcType]*api.ProxyCfg, 100)
}

func ProxyCfgHandler(rw http.ResponseWriter, req *http.Request) {

	servername := req.Header.Get("X-Server-Name")
	serverversion := req.Header.Get("X-Server-Version")

	if servername == "" || serverversion == "" {
		err := fmt.Sprintf("have not found server name & version in request header!\r\n")
		http.Error(rw, err, http.StatusBadRequest)
		log.Println(err)
		return
	}
	svc := api.SvcType{Name: servername, Version: serverversion}

	if req.Method == "GET" {

		proxycfg, b := gProxyCfgMap[svc]
		if b == false {
			err := fmt.Sprintf("can not found (%v) proxy cfg on db base!\r\n", svc)
			http.Error(rw, err, http.StatusNotFound)
			log.Println(err)
			return
		}

		body, err := json.Marshal(proxycfg)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		rw.WriteHeader(http.StatusOK)

		cnt, err := rw.Write(body)
		if err != nil {
			log.Println(err.Error())
		} else if cnt != len(body) {
			log.Println("write to body not finish!")
		}

	} else if req.Method == "POST" {

		proxycfg := &api.ProxyCfg{In: make([]api.InStream, 0), Out: make([]api.OutStream, 0)}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}

		err = json.Unmarshal(body, proxycfg)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}

		gProxyCfgMap[svc] = proxycfg

		rw.WriteHeader(http.StatusOK)

	} else {
		err := fmt.Sprintf("method(%s) is invalid!\r\n", req.Method)
		http.Error(rw, err, http.StatusNotFound)
		log.Println(err)
	}
}
