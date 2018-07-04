package lb

import (
	"sync/atomic"
)

type LBRR struct {
	Index int32
	Array []interface{}
}

func NewLBRR(list []interface{}) LBE {
	rr := new(LBRR)
	rr.Index = 0
	rr.Array = make([]interface{}, len(list))
	copy(rr.Array, list)
	return rr
}

func (rr *LBRR) Select() interface{} {
	if len(rr.Array) == 0 {
		return nil
	}
	before := atomic.AddInt32(&rr.Index, 1)
	before = before % int32(len(rr.Array))
	return rr.Array[before]
}

func (rr *LBRR) ReFlash(list []interface{}) {
	rr.Array = make([]interface{}, len(list))
	copy(rr.Array, list)
}
