package lb

type LBNONE struct {
	Array []interface{}
}

func NewLBNONE(list []interface{}) LBE {
	rr := new(LBNONE)
	rr.Array = make([]interface{}, len(list))
	copy(rr.Array, list)
	return rr
}

func (rr *LBNONE) Select() interface{} {
	if len(rr.Array) == 0 {
		return nil
	}
	return rr.Array[0]
}

func (rr *LBNONE) ReFlash(list []interface{}) {
	rr.Array = make([]interface{}, len(list))
	copy(rr.Array, list)
}
