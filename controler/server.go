package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"

	"github.com/lixiangyun/go-mesh/mesher/api"
)

var gSvcAll map[api.SvcType]*api.SvcBase
var gInstanceID uint32

func init() {

	gSvcAll = make(map[api.SvcType]*api.SvcBase, 0)
}

func NewServerBase(svc api.SvcType) *api.SvcBase {
	svcbase := new(api.SvcBase)
	svcbase.Server = svc
	svcbase.Instance = make([]api.SvcInstance, 0)
	return svcbase
}

func UUID() string {
	return fmt.Sprintf("%08x%08x", rand.Uint32(), atomic.AddUint32(&gInstanceID, 1))
}

func ServerAddInstance(s *api.SvcBase, inst api.SvcInstance) api.SvcInstance {

	if inst.ID == "" {
		inst.ID = UUID()
		s.Instance = append(s.Instance, inst)
		log.Printf("new instance (%s) success!\r\n", inst.ID)
		return inst
	}

	for idx, _ := range s.Instance {

		if s.Instance[idx].ID != inst.ID {
			continue
		}

		if api.InstanceCompare(s.Instance[idx], inst) {
			log.Printf("heartbeat instance (%s) success!\r\n", inst.ID)
		} else {
			s.Instance[idx].Array = make([]api.EndPoint, len(inst.Array))
			copy(s.Instance[idx].Array, inst.Array)
			log.Printf("update instance (%s) success!\r\n", inst.ID)
		}
		return inst
	}

	s.Instance = append(s.Instance, inst)
	log.Printf("add instance (%s) success!\r\n", inst.ID)

	return inst
}

func ServerRegisterHandler(rw http.ResponseWriter, req *http.Request) {

	servername := req.Header.Get("X-Server-Name")
	serverversion := req.Header.Get("X-Server-Version")

	if servername == "" || serverversion == "" {
		err := fmt.Sprintf("have not found \"X-Server-Name\" or \"X-Server-Version\" in request header!\r\n")
		http.Error(rw, err, http.StatusBadRequest)
		log.Println(err)
		return
	}

	svc := api.SvcType{Name: servername, Version: serverversion}

	if req.Method != "POST" {
		err := fmt.Sprintf("method(%s) is invalid!\r\n", req.Method)
		http.Error(rw, err, http.StatusNotFound)
		log.Println(err)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	var instance api.SvcInstance

	err = json.Unmarshal(body, &instance)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	svcbase, b := gSvcAll[svc]
	if b == false {
		svcbase = NewServerBase(svc)
		gSvcAll[svc] = svcbase
	}

	instance = ServerAddInstance(svcbase, instance)

	body, err = json.Marshal(&instance)
	if err != nil {
		log.Println(err.Error())
	}

	rw.WriteHeader(http.StatusOK)

	cnt, err := rw.Write(body)
	if err != nil {
		log.Println(err.Error())
	} else if cnt != len(body) {
		log.Println("write to body not finish!")
	}

	log.Printf("server (%s %s) register/update/heartbeat success!\r\n", svc.Name, svc.Version)
}

func ServerQueryHandler(rw http.ResponseWriter, req *http.Request) {

	servername := req.Header.Get("X-Server-Name")
	serverversion := req.Header.Get("X-Server-Version")

	if servername == "" || serverversion == "" {
		err := fmt.Sprintf("have not found \"X-Server-Name\" or \"X-Server-Version\" in request header!\r\n")
		http.Error(rw, err, http.StatusBadRequest)
		log.Println(err)
		return
	}
	svc := api.SvcType{Name: servername, Version: serverversion}

	if req.Method != "GET" {
		err := fmt.Sprintf("method(%s) is invalid!\r\n", req.Method)
		http.Error(rw, err, http.StatusNotFound)
		log.Println(err)
		return
	}

	svcInstance, b := gSvcAll[svc]
	if b == false {
		err := fmt.Sprintf("can not found (%s %s) svc instance on db base!\r\n", svc.Name, svc.Version)
		http.Error(rw, err, http.StatusNotFound)
		log.Println(err)
		return
	}

	body, err := json.Marshal(&svcInstance)
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

	log.Printf("server (%s %s) query success!\r\n", svc.Name, svc.Version)
}
