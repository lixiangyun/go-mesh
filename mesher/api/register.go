package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type SvcType struct {
	Name    string `json:"servername"`
	Version string `json:"version"`
}

type EndPoint struct {
	Addr string `json:"address"`
	Type int    `json:"type"`
}

type SvcInstance struct {
	ID    string     `json:"instanceid"`
	Array []EndPoint `json:"endpoint"`
	Time  time.Time  `json:"-"`
}

type SvcBase struct {
	Server   SvcType       `json:"server"`
	Instance []SvcInstance `json:"instance"`
}

func InstanceToAddr(instances []SvcInstance, addrtype int) []string {
	addrs := make([]string, 0)
	for _, inst := range instances {
		for _, edp := range inst.Array {
			if edp.Type == addrtype {
				addrs = append(addrs, edp.Addr)
			}
		}
	}
	return addrs
}

func InstanceCompare(a, b SvcInstance) bool {

	if a.ID != b.ID {
		return false
	}

	if len(a.Array) != len(b.Array) {
		return false
	}

	for _, av := range a.Array {
		var bfound bool
		for _, bv := range b.Array {
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

func InstanceArrayCompare(a, b []SvcInstance) bool {

	if len(a) != len(b) {
		return false
	}

	for _, av := range a {
		var bfound bool
		for _, bv := range b {
			if InstanceCompare(av, bv) {
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

func ServerRegister(addr string, svctype SvcType, inst *SvcInstance) error {

	body, err := json.Marshal(inst)
	if err != nil {
		return err
	}

	path := "http://" + addr + "/server/register"
	req, err := http.NewRequest("POST", path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("X-Server-Name", svctype.Name)
	req.Header.Add("X-Server-Version", svctype.Version)

	rsp, err := HttpClient(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		errstr := fmt.Sprintf("register service(%s.%s) failed! (ret=%s)",
			svctype.Name, svctype.Version, rsp.Status)
		return errors.New(errstr)
	}

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, inst)
	if err != nil {
		return err
	}

	return nil
}

func ServerQuery(addr string, svctype SvcType) ([]SvcInstance, error) {

	path := "http://" + addr + "/server/query"
	req, err := http.NewRequest("GET", path, nil)
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
		errstr := fmt.Sprintf("service(%v) does not exist! (ret=%s)", svctype, rsp.Status)
		return nil, errors.New(errstr)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	svcquery := &SvcBase{Instance: make([]SvcInstance, 0)}

	err = json.Unmarshal(body, &svcquery)
	if err != nil {
		return nil, err
	}

	return svcquery.Instance, nil
}
