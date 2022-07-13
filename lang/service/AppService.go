package service

import (
	"github.com/whatap/golib/io"
)

type AppService struct {
	AbstractService
}

func NewAppService() *AppService {
	p := new(AppService)
	return p
}

func (this *AppService) GetServiceType() byte {
	return SERVICE_APP
}

func (this *AppService) Write(out *io.DataOutputX) {
	this.AbstractService.Write(out)
}

func (this *AppService) Read(in *io.DataInputX) {
	this.AbstractService.Read(in)
}
