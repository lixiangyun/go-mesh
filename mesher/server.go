package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/lixiangyun/go-mesh/mesher/api"
	"github.com/lixiangyun/go-mesh/mesher/lb"
	//	"github.com/lixiangyun/go-mesh/mesher/proxy"
)

type Server struct {
	Type     api.SvcType
	Instance []api.SvcInstance
	Lbe      lb.LBE
}

type ProxyChannel struct {
	Before api.ProxyCfg
	Addr   string
	Lbe    lb.LBE
	Svc    []Server
	Stop   chan struct{}
}

type ProxyMap struct {
	sync.RWMutex
	Run map[string]ProxyChannel
}

var gLocalIp []string

var gProxyMap ProxyMap

func init() {
	gLocalIp = getLocalIp()
	gProxyMap.Run = make(map[string]ProxyChannel, 10)
}

func getLocalIp() []string {

	IpAddr := make([]string, 0)

	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		log.Println("Get local IP addr failed!!!")
		return IpAddr
	}

	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				IpAddr = append(IpAddr, ipnet.IP.String())
			}
		}
	}

	return IpAddr
}

// 如果监听地址为 :8080 ，即没有指定IP地址。则添加本地地址做为IP地址。
func parseEndpoint(addr string) []string {

	endpoints := make([]string, 0)
	ipaddr := strings.Split(addr, ":")

	if len(ipaddr) == 2 && ipaddr[0] == "" {
		for _, v := range gLocalIp {
			endpoint := fmt.Sprintf("%s:%s", v, ipaddr[1])
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

// 进入的流量地址，属于该微服务的endpoint地址

func getInstreamEndpoint(proxycfg []api.ProxyCfg) []api.EndPoint {
	endpoints := make([]api.EndPoint, 0)

	for _, proxy := range proxycfg {
		if proxy.ProxyType != api.PROXY_IN {
			continue
		}

		endpointtmp := parseEndpoint(proxy.Addr)
		if len(endpointtmp) != 0 {
			for _, addr := range endpointtmp {
				endpoints = append(endpoints, api.EndPoint{Addr: addr})
			}
		} else {
			endpoints = append(endpoints, api.EndPoint{Addr: proxy.Addr})
		}
	}

	return endpoints
}

func NewProxyChanel(proxy ProxyChannel) error {

	return nil
}

func DelProxyChanel(proxy ProxyChannel) error {
	proxy.Stop <- struct{}{}
	delete(gProxyMap.Run, proxy.Addr)
	return nil
}

func IsSameProxyCfg(a, b api.ProxyCfg) bool {

	if a.Addr != b.Addr || a.LB != b.LB {
		return false
	}

	if a.Protocal != b.Protocal || a.ProxyType != b.ProxyType {
		return false
	}

	if len(a.Service) != len(b.Service) {
		return false
	}

	for idx, av := range a.Service {
		if b.Service[idx] != av {
			return false
		}
	}

	return true
}

func UpdateProxyChanel(proxycfg []api.ProxyCfg) error {

	var err error

	for _, proxy := range proxycfg {

		channel, b := gProxyMap.Run[proxy.Addr]
		if b == false {

			channel = ProxyChannel{Before: proxy}
			err = NewProxyChanel(channel)
			if err != nil {
				return err
			}

			continue
		}

		if IsSameProxyCfg(channel.Before, proxy) {
			log.Println("proxy cfg not change!")
			continue
		}

		// 简单处理，后续可以考虑如何优化升级
		DelProxyChanel(channel)

		channel.Before = proxy
		NewProxyChanel(channel)
	}

	return err
}

func MesherStart(name, version string, addr string) {

	var errcnt int

	svctype := api.SvcType{Name: name, Version: version}

	for {

		if errcnt > 0 {
			time.Sleep(3 * time.Second)
		} else if errcnt > 10 {
			log.Println("retry same times and continue failed! exit mesher proess!")
			return
		}

		proxycfg, err := api.GetProxyCfg(addr, svctype)
		if err != nil {
			log.Println(err.Error())
			errcnt++
			continue
		}

		endpoint := getInstreamEndpoint(proxycfg)

		instance := api.SvcInstance{Edp: endpoint}

		err = api.ServerRegister(addr, svctype, instance)
		if err != nil {
			log.Println(err.Error())
			errcnt++
			continue
		}

		err = UpdateProxyChanel(proxycfg)
		if err != nil {
			log.Println(err.Error())
			errcnt++
			continue
		}

		errcnt = 0
		time.Sleep(1 * time.Second)
	}
}
