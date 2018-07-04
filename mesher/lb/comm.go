package lb

type LBPOLICY_TYPE string

const (
	LBPOLICY_NONE LBPOLICY_TYPE = ""
	LBPOLICY_RR                 = "roundrobin"
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
	} else { //默认采用无lb策略，即始终返回第一个
		return NewLBNONE(list)
	}
}
