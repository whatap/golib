package pack

import (
	"fmt"

	"github.com/whatap/golib/io"
)

type AbstractPack struct {
	Pcode int64
	Oid   int32
	Okind int32
	Onode int32
	Time  int64
}

func (this *AbstractPack) Write(dout *io.DataOutputX) {
	if (this.Okind | this.Onode) == 0 {
		dout.WriteDecimal(this.Pcode)
		dout.WriteInt(this.Oid)
		dout.WriteLong(this.Time)
	} else {
		dout.WriteByte(9)
		dout.WriteDecimal(this.Pcode)
		dout.WriteInt(this.Oid)
		dout.WriteInt(this.Okind)
		dout.WriteInt(this.Onode)
		dout.WriteLong(this.Time)
	}
}
func (this *AbstractPack) Read(din *io.DataInputX) {
	ver := din.ReadByte()
	if ver <= 8 {
		this.Pcode = din.ReadDecimalLen(int(ver))
		this.Oid = din.ReadInt()
		this.Time = din.ReadLong()
		return
	}
	this.Pcode = din.ReadDecimal()
	this.Oid = din.ReadInt()
	this.Okind = din.ReadInt()
	this.Onode = din.ReadInt()
	this.Time = din.ReadLong()
}

// oid 설정   pack interface
func (this *AbstractPack) SetOID(oid int32) {
	this.Oid = oid
}

// pcode 설정   pack interface
func (this *AbstractPack) SetPCODE(pcode int64) {
	this.Pcode = pcode
}

// pcode 설정   pack interface
func (this *AbstractPack) GetPCODE() int64 {
	return this.Pcode
}

// oid 설정   pack interface
func (this *AbstractPack) SetOKIND(okind int32) {
	this.Okind = okind
}

// pcode 설정   pack interface
func (this *AbstractPack) SetONODE(onode int32) {
	this.Onode = onode
}

func (this *AbstractPack) ToString() string {
	return fmt.Sprintln("\nPcode=", this.Pcode, ",Oid=", this.Oid, ",Okind=", this.Okind, ",ONode=", this.Onode, ",Time=", this.Time)
}
