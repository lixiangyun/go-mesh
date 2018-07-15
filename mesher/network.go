package main

import (
	"net"
	"strings"

	"github.com/lixiangyun/go-mesh/mesher/log"
)

var gLocalIp []string

func init() {

	IpAddr := make([]string, 0)

	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		log.Println(log.WARNING, "Get local IP addr failed!!!")
		return
	}

	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() && ipnet.IP.IsGlobalUnicast() {
				IpAddr = append(IpAddr, ipnet.IP.String())
			}
		}
	}

	gLocalIp = IpAddr
}

func SetNetWork(addr string) {
	if addr != "" {
		gLocalIp = strings.Split(addr, ",")
	}
	log.Printf(log.INFO, "NetWork %+v\r\n", gLocalIp)
}
