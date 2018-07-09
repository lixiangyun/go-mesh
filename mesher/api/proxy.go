package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/lixiangyun/go-mesh/mesher/lb"
	"github.com/lixiangyun/go-mesh/mesher/proxy"
)

type Server struct {
	Svc SvcType          `json:"type"`
	LB  lb.LBPOLICY_TYPE `json:"policy,omitempty"`
}

type OutStream struct {
	Protocal proxy.PROTOCOL_TYPE `json:"protocal"`
	Addr     string              `json:"listen"`
	Svc      []Server            `json:"service"`
	LB       lb.LBPOLICY_TYPE    `json:"policy,omitempty"`
}

type InStream struct {
	Protocal proxy.PROTOCOL_TYPE `json:"protocal"`
	Addr     string              `json:"listen"`
	Local    []string            `json:"local"`
	LB       lb.LBPOLICY_TYPE    `json:"policy,omitempty"`
}

type ProxyCfg struct {
	Out []OutStream `json:"out,omitempty"`
	In  []InStream  `json:"in,omitempty"`
}

func OutStreamCompare(a, b OutStream) bool {
	if a.Protocal != b.Protocal || a.Addr != b.Addr || a.LB != b.LB {
		return false
	}

	if len(a.Svc) != len(b.Svc) {
		return false
	}

	for _, av := range a.Svc {
		var bfound bool
		for _, bv := range b.Svc {
			if av == bv {
				bfound = true
				break
			}
		}
		if bfound == false {
			return false
		}
	}

	return true
}

func InStreamCompare(a, b InStream) bool {

	if a.Protocal != b.Protocal || a.Addr != b.Addr || a.LB != b.LB {
		return false
	}

	if len(a.Local) != len(b.Local) {
		return false
	}

	for _, av := range a.Local {
		var bfound bool
		for _, bv := range b.Local {
			if av == bv {
				bfound = true
				break
			}
		}
		if bfound == false {
			return false
		}
	}

	return true
}

func ProxyCfgCompare(a, b ProxyCfg) bool {

	if len(a.In) != len(b.In) || len(a.Out) != len(b.Out) {
		return false
	}

	for _, av := range a.In {
		var bfound bool
		for _, bv := range b.In {
			if InStreamCompare(av, bv) {
				bfound = true
				break
			}
		}
		if bfound == false {
			return false
		}
	}

	for _, av := range a.Out {
		var bfound bool
		for _, bv := range b.Out {
			if OutStreamCompare(av, bv) {
				bfound = true
				break
			}
		}
		if bfound == false {
			return false
		}
	}

	return true
}

func LoadProxyCfg(addr string, svctype SvcType) (*ProxyCfg, error) {

	url := "http://" + addr + "/proxy/cfg"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Server-Name", svctype.Name)
	req.Header.Add("X-Server-Version", svctype.Version)

	rsp, err := HttpClient(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("get proxy cfg failed! " + rsp.Status)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var proxycfg ProxyCfg

	err = json.Unmarshal(body, &proxycfg)
	if err != nil {
		return nil, err
	}

	return &proxycfg, nil
}
