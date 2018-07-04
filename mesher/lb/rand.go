package lb

import (
	"math/rand"
)

type LBRandmod struct {
	Array []interface{}
}

func NewLBRandmod(list []interface{}) LBE {
	rr := new(LBRandmod)
	rr.Array = make([]interface{}, len(list))
	copy(rr.Array, list)
	return rr
}

func (rr *LBRandmod) Select() interface{} {
	if len(rr.Array) == 0 {
		return nil
	}
	Index := rand.Int31() % int32(len(rr.Array))
	return rr.Array[Index]
}

func (rr *LBRandmod) ReFlash(list []interface{}) {
	rr.Array = make([]interface{}, len(list))
	copy(rr.Array, list)
}
