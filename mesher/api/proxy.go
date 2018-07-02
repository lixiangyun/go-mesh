package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/lixiangyun/go-mesh/mesher/lb"
)

type PROXY_TYPE string

const (
	PROXY_IN  PROXY_TYPE = "in"
	PROXY_OUT            = "out"
)

type PROTOCOL_TYPE string

const (
	PROTOCOL_TCP  PROTOCOL_TYPE = "tcp"
	PROTOCOL_HTTP               = "http"
)

type Server struct {
	Svc SvcType          `json:"svctype"`
	Lb  lb.LBPOLICY_TYPE `json:"policy"`
}

type ProxyCfg struct {
	Protocal  PROTOCOL_TYPE    `json:"protocal"`
	Addr      string           `json:"listen"`
	ProxyType PROXY_TYPE       `json:"type"`
	Service   []Server         `json:"cluster"`
	LB        lb.LBPOLICY_TYPE `json:"policy"`
}

type ProxyCfgAll struct {
	ProxyList []ProxyCfg `json:"proxy"`
}

func GetProxyCfg(addr string, svctype SvcType) ([]ProxyCfg, error) {

	transport := http.DefaultTransport

	url := "http://" + addr + "/proxycfg"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Server-Name", svctype.Name)
	req.Header.Add("X-Server-Version", svctype.Version)

	rsp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("get proxy cfg failed! " + rsp.Status)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	proxycfg := new(ProxyCfgAll)
	proxycfg.ProxyList = make([]ProxyCfg, 0)

	err = json.Unmarshal(body, proxycfg)
	if err != nil {
		return nil, err
	}

	return proxycfg.ProxyList, nil
}
