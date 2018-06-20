package main

import (
	"github.com/lixiangyun/go-mesh/mesher/proxy"
)




func main() {

	httpproxy := proxy.NewHttpProxy(":808", "127.0.0.1:809")

	tcpproxy := proxy.NewTcpProxy(":123", "127.0.0.1:124")

	
	

}
