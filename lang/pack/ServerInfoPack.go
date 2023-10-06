package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/dateutil"
)

const RUMCTL string = "rumctl"

type ServerInfoPack struct {
	AbstractPack

	Version    int32
	Host       string
	Port       int64
	ServerName string
	UpTime     int64
	Attr       *value.MapValue
	KeepTime   int64
}

func NewServerInfoPack() *ServerInfoPack {
	p := new(ServerInfoPack)
	p.Attr = value.NewMapValue()
	p.KeepTime = 10000
	p.UpTime = dateutil.Now()
	return p
}

func (this *ServerInfoPack) GetPackType() int16 {
	return PACK_SERVERINFO
}

func (this *ServerInfoPack) Write(dout *io.DataOutputX) {
	dout.WriteInt3(this.Version) //this.Version)
	dout.WriteDecimal(this.Port) //this.Port)
	dout.WriteDecimal(this.UpTime)
	dout.WriteDecimal(this.KeepTime)
	dout.WriteText(this.ServerName) //this.ServerName)
	this.Attr.Write(dout)
}

func (this *ServerInfoPack) Read(din *io.DataInputX) {
	this.Version = din.ReadInt3() //this.Version)
	this.Port = din.ReadDecimal() //this.Port)
	this.UpTime = din.ReadDecimal()
	this.KeepTime = din.ReadDecimal()
	this.ServerName = din.ReadText() //this.ServerName)
	this.Attr = value.ReadMapValue(din)
}
