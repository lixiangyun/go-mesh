package proxy

import (
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type connect struct {
	sessionID uint32
	conn1     net.Conn
	conn2     net.Conn
}

type TcpProxy struct {
	ListenAddr string
	SelectAddr SELECT_ADDR

	listen    net.Listener
	sessionID uint32
	connbuf   map[uint32]*connect
	sync.RWMutex

	stop chan struct{}
}

func NewTcpProxy(addr string, fun SELECT_ADDR) *TcpProxy {

	proxy := &TcpProxy{ListenAddr: addr, SelectAddr: fun}

	proxy.connbuf = make(map[uint32]*connect)
	proxy.stop = make(chan struct{}, 1)

	listen, err := net.Listen("tcp", proxy.ListenAddr)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	log.Printf("Tcp Proxy Listen : %s\r\n", addr)

	proxy.listen = listen

	go proxy.start()

	return proxy
}

func writeFull(conn net.Conn, buf []byte) error {
	totallen := len(buf)
	sendcnt := 0

	for {
		cnt, err := conn.Write(buf[sendcnt:])
		if err != nil {
			return err
		}
		if cnt+sendcnt >= totallen {
			return nil
		}
		sendcnt += cnt
	}
}

// tcp通道互通
func tcpChannel(localconn net.Conn, remoteconn net.Conn, wait *sync.WaitGroup) {

	defer wait.Done()
	defer localconn.Close()
	defer remoteconn.Close()

	buf := make([]byte, 65535)

	for {
		cnt, err := localconn.Read(buf[0:])
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
			}
			break
		}

		err = writeFull(remoteconn, buf[0:cnt])
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
			}
			break
		}
	}
}

// tcp代理处理
func (t *TcpProxy) tcpProxyProcess(wait *sync.WaitGroup, conn *connect) {

	syncSem := new(sync.WaitGroup)

	defer wait.Done()

	log.Println("new connect. ", conn.conn1.RemoteAddr(), " -> ", conn.conn2.RemoteAddr())

	syncSem.Add(2)

	go tcpChannel(conn.conn1, conn.conn2, syncSem)
	go tcpChannel(conn.conn2, conn.conn1, syncSem)

	syncSem.Wait()

	t.Lock()
	t.connbuf[conn.sessionID] = nil
	t.Unlock()

	log.Println("close connect. ", conn.conn1.RemoteAddr(), " -> ", conn.conn2.RemoteAddr())
}

// 正向tcp代理启动和处理入口
func (t *TcpProxy) start() {

	var wait sync.WaitGroup

	for {

		session := atomic.AddUint32(&t.sessionID, 1)

		localconn, err := t.listen.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		remoteconn, err := net.Dial("tcp", t.SelectAddr())
		if err != nil {
			log.Println(err.Error())
			localconn.Close()
			continue
		}

		newconn := &connect{conn1: localconn, conn2: remoteconn, sessionID: session}

		t.Lock()
		t.connbuf[session] = newconn
		t.Unlock()

		wait.Add(1)

		go t.tcpProxyProcess(&wait, newconn)
	}

	wait.Wait()

	t.stop <- struct{}{}
	t.listen = nil
}

func (t *TcpProxy) Close() {
	t.listen.Close()

	for _, v := range t.connbuf {
		if v != nil {
			v.conn1.Close()
			v.conn2.Close()
		}
	}

	<-t.stop
}
