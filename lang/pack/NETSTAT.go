package pack

import ()

type NETSTAT struct {
	Est  int32
	FinW int32
	TimW int32
	CloW int32
}

func NewNETSTAT() *NETSTAT {
	p := new(NETSTAT)
	return p
}
