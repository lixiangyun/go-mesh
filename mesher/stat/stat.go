package stat

import (
	"fmt"
	"log"
	"time"
)

type Item struct {
	sendSize int
	recvSize int
	sendCnt  int
	recvCnt  int
}

type Stat struct {
	now      Item
	old      Item
	interval int
	stop     chan struct{}
}

func (now *Item) Sub(old Item) {
	now.sendCnt -= old.sendCnt
	now.sendSize -= old.sendSize
	now.recvCnt -= old.recvCnt
	now.recvSize -= old.recvSize
}

func (now *Item) Div(interval int) {
	now.sendCnt = now.sendCnt / interval
	now.sendSize = now.sendSize / interval
	now.recvCnt = now.recvCnt / interval
	now.recvSize = now.recvSize / interval
}

func calcUnit(cnt int) string {
	if cnt < 1024 {
		return fmt.Sprintf("%d", cnt)
	} else if cnt < 1024*1024 {
		return fmt.Sprintf("%.3fk", float32(cnt)/1024)
	} else if cnt < 1024*1024*1024 {
		return fmt.Sprintf("%.3fM", float32(cnt)/(1024*1024))
	} else {
		return fmt.Sprintf("%.3fG", float32(cnt)/(1024*1024*1024))
	}
}

func (now *Item) Format() string {

	str := fmt.Sprintf("sendCnt : %s ,Size %s\r\n",
		calcUnit(now.sendCnt), calcUnit(now.sendSize))

	str += fmt.Sprintf("recvCnt : %s ,Size %s\r\n",
		calcUnit(now.recvCnt), calcUnit(now.recvSize))

	return str
}

func (s *Stat) display() {
	timer := time.NewTimer(time.Duration(s.interval) * time.Second)
	for {
		select {
		case <-timer.C:
			{
				now := s.now
				old := s.old

				now.Sub(old)
				now.Div(s.interval)
				str := now.Format()

				log.Println("stat result : ", str)

				s.old = s.now

				timer.Reset(time.Duration(s.interval) * time.Second)
			}
		case <-s.stop:
			{
				timer.Stop()
				return
			}
		}
	}
}

func NewStat(interval int) *Stat {
	s := new(Stat)
	s.interval = interval
	s.stop = make(chan struct{}, 1)
	go s.display()
	return s
}

func (s *Stat) Send(size int) {
	s.now.sendCnt++
	s.now.sendSize += size
}

func (s *Stat) Recv(size int) {
	s.now.recvCnt++
	s.now.recvSize += size
}

func (s *Stat) Delete() {
	s.stop <- struct{}{}
}
