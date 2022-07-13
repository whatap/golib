package service

import (
	"github.com/whatap/golib/io"
)

const (
	SERVICE_WAS   = 1
	SERVICE_APP   = 2
	SERVICE_WAS_2 = 3
)

type Service interface {
	GetServiceType() byte
	Write(out *io.DataOutputX)
	Read(in *io.DataInputX)
}

type AbstractService struct {
	Seq     int64
	EndTime int64
	Service int32
	Elapsed int32
	Error   int64
	CpuTime int32
	Malloc  int64

	SqlCount      int32
	SqlTime       int32
	SqlFetchCount int32
	SqlFetchTime  int32

	HttpcCount     int32
	HttpcTime      int32
	Active         bool
	Steps_data_pos int64

	Mtid    int64
	Mdepth  int32
	Mcaller int64
}

func CreateService(t byte) Service {
	switch t {
	case SERVICE_WAS:
		return NewWasService()
	case SERVICE_APP:
		return NewAppService()
	case SERVICE_WAS_2:
		return NewWasService2()
	}
	return nil
}
func (this *AbstractService) Write(dout *io.DataOutputX) {
	dout.WriteLong(this.Seq)
	dout.WriteDecimal(this.EndTime)
	dout.WriteDecimal(int64(this.Service))

	dout.WriteDecimal(int64(this.Elapsed))
	dout.WriteDecimal(this.Error)
	dout.WriteDecimal(int64(this.CpuTime))

	dout.WriteDecimal(int64(this.SqlCount))
	dout.WriteDecimal(int64(this.SqlTime))
	dout.WriteDecimal(int64(this.SqlFetchCount))
	dout.WriteDecimal(int64(this.SqlFetchTime))

	dout.WriteDecimal(this.Malloc)

	dout.WriteDecimal(int64(this.HttpcCount))
	dout.WriteDecimal(int64(this.HttpcTime))
	dout.WriteBool(this.Active)
	dout.WriteDecimal(this.Steps_data_pos)

}

func (this *AbstractService) Read(din *io.DataInputX) {
	this.Seq = din.ReadLong()
	this.EndTime = din.ReadDecimal()
	this.Service = int32(din.ReadDecimal())

	this.Elapsed = int32(din.ReadDecimal())
	this.Error = din.ReadDecimal()
	this.CpuTime = int32(din.ReadDecimal())

	this.SqlCount = int32(din.ReadDecimal())
	this.SqlTime = int32(din.ReadDecimal())
	this.SqlFetchCount = int32(din.ReadDecimal())
	this.SqlFetchTime = int32(din.ReadDecimal())

	this.Malloc = din.ReadDecimal()

	this.HttpcCount = int32(din.ReadDecimal())
	this.HttpcTime = int32(din.ReadDecimal())
	this.Active = din.ReadBool()
	this.Steps_data_pos = din.ReadDecimal()
}

func ToBytes(s Service, dout *io.DataOutputX) {
	dout.WriteByte(s.GetServiceType())
	s.Write(dout)
}

func ToObject(din *io.DataInputX) Service {
	service := CreateService(din.ReadByte())
	service.Read(din)
	return service
}
