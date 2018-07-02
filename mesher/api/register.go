package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	ID  string     `json:"instanceid"`
	Edp []EndPoint `json:"endpoint"`
}

type SvcQuery struct {
	Server   SvcType       `json:"server"`
	Instance []SvcInstance `json:"instance"`
}

func ServerRegister(addr string, svctype SvcType, inst SvcInstance) error {

	body, err := json.Marshal(inst)
	if err != nil {
		return err
	}

	transport := http.DefaultTransport

	path := "http://" + addr + "/server/register"
	req, err := http.NewRequest("POST", path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("X-Server-Name", svctype.Name)
	req.Header.Add("X-Server-Version", svctype.Version)

	rsp, err := transport.RoundTrip(req)
	if rsp.StatusCode != http.StatusOK {

		errstr := fmt.Sprintf("register service(%s.%s) failed! (ret=%s)",
			svctype.Name, svctype.Version, rsp.Status)
		return errors.New(errstr)
	}

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &inst.ID)
	if err != nil {
		return err
	}

	return nil
}

func ServerQuery(addr string, svctype SvcType) ([]SvcInstance, error) {

	transport := http.DefaultTransport

	path := "http://" + addr + "/server/query"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Server-Name", svctype.Name)
	req.Header.Add("X-Server-Version", svctype.Version)

	rsp, err := transport.RoundTrip(req)
	if rsp.StatusCode != http.StatusOK {

		errstr := fmt.Sprintf("register service(%s.%s) failed! (ret=%s)",
			svctype.Name, svctype.Version, rsp.Status)

		return nil, errors.New(errstr)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	svcquery := new(SvcQuery)

	err = json.Unmarshal(body, &svcquery)
	if err != nil {
		return nil, err
	}

	return svcquery.Instance, nil
}
