package pack

import ()

type WEBSOCKET struct {
	Count int32
	In    int64
	Out   int64
}

func NewWEBSOCKET() *WEBSOCKET {
	p := new(WEBSOCKET)
	return p
}
