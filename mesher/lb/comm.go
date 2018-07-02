package lb

type LBPOLICY_TYPE string

const (
	LBPOLICY_RR   LBPOLICY_TYPE = "roundrobin"
	LBPOLICY_RAND               = "random"
)

type LBE interface {
	Select() interface{}
	ReFlash([]interface{})
}

func NewLB(policy LBPOLICY_TYPE, list []interface{}) LBE {

	if policy == LBPOLICY_RR {
		return NewLBRR(list)
	} else if policy == LBPOLICY_RAND {
		return NewLBRandmod(list)
	}

	return nil
}
