package service

import (
	"github.com/whatap/golib/io"
)

type WasService struct {
	AbstractService

	IpAddr    int32
	WClientId int64
	UserAgent int32
	Referer   int32
	Status    int32

	Mtid    int64
	Mdepth  int32
	Mcaller int64
}

func NewWasService() *WasService {
	p := new(WasService)
	return p
}

func (this *WasService) GetServiceType() byte {
	return SERVICE_WAS
}

func (this *WasService) Write(out *io.DataOutputX) {
	this.AbstractService.Write(out)

	out.WriteInt(this.IpAddr)
	out.WriteDecimal(0)
	out.WriteDecimal(this.WClientId)
	out.WriteDecimal(int64(this.UserAgent))
	out.WriteDecimal(int64(this.Referer))
	out.WriteDecimal(int64(this.Status))

	out.WriteDecimal(this.Mtid)
	out.WriteDecimal(int64(this.Mdepth))
	out.WriteDecimal(this.Mcaller)
}

func (this *WasService) Read(in *io.DataInputX) {
	this.AbstractService.Read(in)

	this.IpAddr = in.ReadInt()
	in.ReadDecimal()
	this.WClientId = in.ReadDecimal()
	this.UserAgent = int32(in.ReadDecimal())
	this.Referer = int32(in.ReadDecimal())
	this.Status = int32(in.ReadDecimal())

	this.Mtid = in.ReadDecimal()
	this.Mdepth = int32(in.ReadDecimal())
	this.Mcaller = in.ReadDecimal()

}
