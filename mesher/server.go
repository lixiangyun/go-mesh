package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lixiangyun/go-mesh/mesher/api"
	"github.com/lixiangyun/go-mesh/mesher/lb"
	"github.com/lixiangyun/go-mesh/mesher/proxy"
)

type SvcCluster struct {
	Svc      api.SvcType
	Instance []api.SvcInstance
	Addr     []string
	Lbe      lb.LBE
}

type OutChannel struct {
	Proxy proxy.PROXY
	Lbe   lb.LBE
	Svc   []SvcCluster
}

type InChannel struct {
	Proxy proxy.PROXY
	Lbe   lb.LBE
	Addr  []string
}

type ProxyMap struct {
	Cfg api.ProxyCfg

	InChan  map[string]*InChannel
	OutChan map[string]*OutChannel

	Svc      api.SvcType     // 本服务类型信息
	Instance api.SvcInstance // 本服务实例信息

	DepInstance map[api.SvcType][]api.SvcInstance // 依赖的服务实例信息缓存
}

var gProxyMap ProxyMap

var gControlerAddr string

func init() {

	gProxyMap.InChan = make(map[string]*InChannel, 0)
	gProxyMap.OutChan = make(map[string]*OutChannel, 0)

	gProxyMap.DepInstance = make(map[api.SvcType][]api.SvcInstance, 0)
}

func (s *SvcCluster) SelectAddr() string {
	addr := s.Lbe.Select()
	if addr == nil {
		return ""
	}
	return *(addr.(*string))
}

func (in *InChannel) SelectAddr() string {
	addr := in.Lbe.Select()
	if addr == nil {
		return ""
	}
	return *(addr.(*string))
}

func (out *OutChannel) SelectAddr() string {
	selector := out.Lbe.Select()
	svc := selector.(*SvcCluster)
	if svc == nil {
		return ""
	}
	return svc.SelectAddr()
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
func getInstreamEndpoint(inlist []api.InStream) []api.EndPoint {

	endpoints := make([]api.EndPoint, 0)

	for _, in := range inlist {

		endpointtmp := parseEndpoint(in.Addr)
		if len(endpointtmp) != 0 {
			for _, addr := range endpointtmp {
				endpoints = append(endpoints, api.EndPoint{Addr: addr})
			}
		} else {
			endpoints = append(endpoints, api.EndPoint{Addr: in.Addr})
		}
	}

	return endpoints
}

func NewInChannel(in api.InStream) *InChannel {

	channel := new(InChannel)
	channel.Addr = in.Local

	selectlist := make([]interface{}, len(channel.Addr))

	for i := 0; i < len(channel.Addr); i++ {
		selectlist[i] = &channel.Addr[i]
	}

	channel.Lbe = lb.NewLB(in.LB, selectlist)
	channel.Proxy = proxy.NewProxy(in.Protocal, in.Addr, channel.SelectAddr)

	if channel.Proxy == nil {
		return nil
	}

	return channel
}

func NewSvcCluster(svccfg api.Server) *SvcCluster {
	svcCluster := new(SvcCluster)
	svcCluster.Svc = svccfg.Svc

	instances, err := api.ServerQuery(gControlerAddr, svccfg.Svc)
	if err != nil {
		log.Printf("server(%v) query failed! (%s)", svccfg.Svc, err.Error())
	}

	svcCluster.Instance = instances
	svcCluster.Addr = api.InstanceToAddr(instances, 0)

	selectlist := make([]interface{}, len(svcCluster.Addr))

	for i := 0; i < len(svcCluster.Addr); i++ {
		selectlist[i] = &svcCluster.Addr[i]
	}

	svcCluster.Lbe = lb.NewLB(svccfg.LB, selectlist)
	return svcCluster
}

func NewOutChannel(in api.OutStream) *OutChannel {

	channel := new(OutChannel)
	channel.Svc = make([]SvcCluster, 0)

	// 服务集群信息
	for _, svccfg := range in.Svc {
		svccluster := NewSvcCluster(svccfg)
		if svccluster == nil {
			continue
		}
		channel.Svc = append(channel.Svc, *svccluster)
	}

	selectlist := make([]interface{}, len(channel.Svc))
	for i := 0; i < len(channel.Svc); i++ {
		selectlist[i] = &channel.Svc[i]
	}

	channel.Lbe = lb.NewLB(in.LB, selectlist)
	channel.Proxy = proxy.NewProxy(in.Protocal, in.Addr, channel.SelectAddr)

	if channel.Proxy == nil {
		return nil
	}

	return channel
}

func FlashSvcCluster(svcCluster *SvcCluster) {

	instances, err := api.ServerQuery(gControlerAddr, svcCluster.Svc)
	if err != nil {
		log.Printf("server(%v) query failed! (%s)", svcCluster.Svc, err.Error())
		return
	}

	if api.InstanceArrayCompare(svcCluster.Instance, instances) {
		return
	}

	svcCluster.Instance = instances
	svcCluster.Addr = api.InstanceToAddr(instances, 0)

	selectlist := make([]interface{}, len(svcCluster.Addr))

	for i := 0; i < len(svcCluster.Addr); i++ {
		selectlist[i] = &svcCluster.Addr[i]
	}

	svcCluster.Lbe.ReFlash(selectlist)
}

func FlashOutChannel(channel *OutChannel) {

	// 更新服务集群的实例信息
	for idx, _ := range channel.Svc {
		FlashSvcCluster(&channel.Svc[idx])
	}
}

func UpdateProxyChanel(proxycfg *api.ProxyCfg) error {

	for _, InStream := range proxycfg.In {
		channel, b := gProxyMap.InChan[InStream.Addr]
		if b == true {
			continue
		}
		channel = NewInChannel(InStream)
		if channel != nil {
			gProxyMap.InChan[InStream.Addr] = channel
		}
	}

	for _, OutStream := range proxycfg.Out {
		channel, b := gProxyMap.OutChan[OutStream.Addr]
		if b == true {
			FlashOutChannel(channel)
			continue
		}
		channel = NewOutChannel(OutStream)
		if channel != nil {
			gProxyMap.OutChan[OutStream.Addr] = channel
		}
	}

	gProxyMap.Cfg = *proxycfg

	//log.Println("update proxy success!")

	return nil
}

func MesherStart(name, version string, addr string) {

	var errcnt int

	gControlerAddr = addr

	svctype := api.SvcType{Name: name, Version: version}
	gProxyMap.Svc = svctype

	for {

		if errcnt > 0 {
			time.Sleep(3 * time.Second)
		} else if errcnt > 10 {
			log.Println("retry same times and continue failed! exit mesher proess!")
			return
		}

		proxycfg, err := api.LoadProxyCfg(gControlerAddr, svctype)
		if err != nil {
			log.Println("load proxy cfg from controler failed! ", err.Error())
			errcnt++
			continue
		}

		endpoint := getInstreamEndpoint(proxycfg.In)

		instance := &api.SvcInstance{ID: gProxyMap.Instance.ID, Array: endpoint}

		err = api.ServerRegister(gControlerAddr, svctype, instance)
		if err != nil {
			log.Println("server register failed! ", err.Error())
			errcnt++
			continue
		}

		if !api.InstanceCompare(gProxyMap.Instance, *instance) {
			log.Println("get instance ", instance.ID)
			gProxyMap.Instance = *instance
		}

		err = UpdateProxyChanel(proxycfg)
		if err != nil {
			log.Println(err.Error())
			errcnt++
			continue
		}

		errcnt = 0
		time.Sleep(100 * time.Second)
	}
}
