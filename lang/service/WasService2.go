package service

import (
	"github.com/whatap/golib/io"
)

type WasService2 struct {
	WasService
}

func NewWasService2() *WasService2 {
	p := new(WasService2)
	return p
}

func (this *WasService2) GetServiceType() byte {
	return SERVICE_WAS_2
}

func (this *WasService2) Write(out *io.DataOutputX) {
	this.WasService.Write(out)
}

func (this *WasService2) Read(in *io.DataInputX) {
	this.WasService.Read(in)
}
