package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/lixiangyun/go-mesh/mesher/api"
)

const (
	INSTANCE_TIMEOUT = 3
)

type ServerInstance struct {
	Svc       api.SvcType
	Instances map[string]api.SvcInstance
}

type ServerBaseDB struct {
	sync.Mutex

	SvcList map[api.SvcType]*ServerInstance
}

var gSvcDB ServerBaseDB

func init() {
	gSvcDB.SvcList = make(map[api.SvcType]*ServerInstance, 0)
}

func ServerCheckTimeout() {
	go func() {

		for {
			gSvcDB.Lock()

			instdel := make([]api.SvcInstance, 0)
			newtime := time.Now()

			for _, svc := range gSvcDB.SvcList {
				for _, inst := range svc.Instances {
					subtime := newtime
					if subtime.Sub(inst.Time) >= time.Second*INSTANCE_TIMEOUT {
						instdel = append(instdel, inst)
					}
				}
				for _, inst := range instdel {
					delete(svc.Instances, inst.ID)
				}
			}

			gSvcDB.Unlock()

			time.Sleep(1 * time.Second)
		}
	}()
}

func (svcbase *ServerInstance) ServerAddInstance(inst api.SvcInstance) api.SvcInstance {

	inst.Time = time.Now()

	if inst.ID == "" {
		for {
			inst.ID = UUID()
			_, b := svcbase.Instances[inst.ID]
			if b == false {
				svcbase.Instances[inst.ID] = inst
				log.Printf("new instance (%v) success!\r\n", inst)
				return inst
			}
		}
	}

	instold, b := svcbase.Instances[inst.ID]
	if b == false {
		log.Printf("add instance (%v) success!\r\n", inst)
	} else {
		if api.InstanceCompare(instold, inst) {
			log.Printf("heartbeat instance (%v) success!\r\n", inst)
		} else {
			log.Printf("update instance (%v) success!\r\n", inst)
		}
	}

	svcbase.Instances[inst.ID] = inst

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

	gSvcDB.Lock()
	defer gSvcDB.Unlock()

	svcbase, b := gSvcDB.SvcList[svc]
	if b == false {
		svcbase = &ServerInstance{Svc: svc}
		svcbase.Instances = make(map[string]api.SvcInstance, 0)
		gSvcDB.SvcList[svc] = svcbase
	}

	instance = svcbase.ServerAddInstance(instance)

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

	gSvcDB.Lock()
	defer gSvcDB.Unlock()

	svcbase, b := gSvcDB.SvcList[svc]
	if b == false {
		err := fmt.Sprintf("can not found (%s %s) svc instance on db base!\r\n", svc.Name, svc.Version)
		http.Error(rw, err, http.StatusNotFound)
		log.Println(err)
		return
	}

	Instances := make([]api.SvcInstance, len(svcbase.Instances))
	var idx int
	for _, inst := range svcbase.Instances {
		Instances[idx] = inst
		idx++
	}

	svcqurey := &api.SvcBase{Server: svc, Instance: Instances}

	body, err := json.Marshal(&svcqurey)
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
