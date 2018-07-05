package proxy

import (
	"log"
)

type PROTOCOL_TYPE string

const (
	PROTOCOL_TCP  PROTOCOL_TYPE = "tcp"
	PROTOCOL_HTTP               = "http"
)

type SELECT_ADDR func() string

type PROXY interface {
	Close()
}

func NewProxy(protocal PROTOCOL_TYPE, addr string, fun SELECT_ADDR) PROXY {
	if protocal == PROTOCOL_TCP {
		return NewTcpProxy(addr, fun)
	} else if protocal == PROTOCOL_HTTP {
		return NewHttpProxy(addr, fun)
	} else {
		log.Printf("protocal %s not support.\r\n", protocal)
		return nil
	}
}
