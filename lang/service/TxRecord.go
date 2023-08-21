package service

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/hash"
)

const (
	FATAL   byte = 30
	WARNING byte = 20
	INFO    byte = 10
	NONE    byte = 0
)

type TxRecord struct {
	Txid    int64
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

	HttpcCount int32
	HttpcTime  int32

	Active       bool
	StepsDataPos int64
	Cipher       int32

	IpAddr    int32
	WClientId int64
	UserAgent int32
	Referer   int32
	Status    int32

	Mtid    int64
	Mdepth  int32
	Mcaller int64

	McallerPcode int64
	McallerOkind int32
	McallerOid   int32
	McallerSpec  int32
	McallerUrl   int32
	MthisSpec    int32

	HttpMethod byte

	Domain int32
	Fields *value.MapValue

	Login int32

	ErrorLevel byte

	Oid   int32
	Okind int32
	Onode int32

	Uuid string

	DbcTime int32

	//20210628 -
	Apdex byte

	McallerStepId int64
	OriginUrl     string

	StepSplitCount int
}

func NewTxRecord() *TxRecord {
	p := new(TxRecord)
	return p
}
func (this *TxRecord) setUuid(uuid string) int64 {
	this.Uuid = uuid
	if uuid == "" {
		this.Mtid = 0
	} else {
		this.Mtid = hash.Hash64Str(uuid)
	}
	return this.Mtid
}
func (this *TxRecord) Write(dout *io.DataOutputX) {

	//version
	dout.WriteByte(10)

	o := io.NewDataOutputX()

	o.WriteLong(this.Txid)
	o.WriteDecimal(this.EndTime)
	o.WriteDecimal(int64(this.Service))
	o.WriteDecimal(int64(this.Elapsed))
	o.WriteDecimal(this.Error)
	o.WriteDecimal(int64(this.CpuTime))
	o.WriteDecimal(this.Malloc)

	o.WriteDecimal(int64(this.SqlCount))
	o.WriteDecimal(int64(this.SqlTime))
	o.WriteDecimal(int64(this.SqlFetchCount))
	o.WriteDecimal(int64(this.SqlFetchTime))

	o.WriteDecimal(int64(this.HttpcCount))
	o.WriteDecimal(int64(this.HttpcTime))

	o.WriteBool(this.Active)
	o.WriteDecimal(this.StepsDataPos)
	o.WriteDecimal(int64(this.Cipher))

	o.WriteInt(this.IpAddr)
	o.WriteDecimal(int64(this.WClientId))
	o.WriteDecimal(int64(this.UserAgent))
	o.WriteDecimal(int64(this.Referer))
	o.WriteDecimal(int64(this.Status))

	if this.Mtid != 0 {
		o.WriteByte(1)
		o.WriteDecimal(this.Mtid)
		o.WriteDecimal(int64(this.Mdepth))
		o.WriteDecimal(this.Mcaller)
	} else {
		o.WriteByte(0)
	}
	if this.McallerPcode != 0 {
		o.WriteByte(6)
		o.WriteDecimal(this.McallerPcode)
		o.WriteDecimal(int64(this.McallerOkind))
		o.WriteDecimal(int64(this.McallerOid))
		o.WriteDecimal(int64(this.McallerSpec))
		o.WriteDecimal(int64(this.McallerUrl))
		o.WriteDecimal(int64(this.MthisSpec)) // 보통 0, 서버에서 셋팅
	} else {
		o.WriteByte(0)
	}

	o.WriteByte(this.HttpMethod)
	o.WriteDecimal(int64(this.Domain))

	if this.Fields == nil {
		o.WriteByte(0)
	} else {
		sz := this.Fields.Size()
		o.WriteByte(byte(sz))
		keys := this.Fields.Keys()
		for keys.HasMoreElements() {
			key := keys.NextString()
			tmp := this.Fields.Get(key)
			if tmp != nil {
				if v, ok := tmp.(value.Value); ok {
					o.WriteText(key)
					value.WriteValue(o, v)
				}
			} else {
				// empty string
				o.WriteText(key)
				value.WriteValue(o, value.NewTextValue(""))
			}
		}
	}
	o.WriteDecimal(int64(this.Login))

	o.WriteByte(this.ErrorLevel)

	// 2018.12.28
	o.WriteDecimal(int64(this.Oid))
	o.WriteDecimal(int64(this.Okind))
	o.WriteDecimal(int64(this.Onode))

	// To-DO
	// 2019.04.04
	o.WriteText(this.Uuid)
	//
	// 2019.07.09
	o.WriteDecimal(int64(this.DbcTime))

	// 2021.06.28 , java //2021.05.13
	o.WriteByte(this.Apdex)

	// 2023.07.17 , java //2021.12.10
	o.WriteDecimal(this.McallerStepId)
	o.WriteText(this.OriginUrl)

	o.WriteDecimal(int64(this.StepSplitCount))

	////////////// BLOB ///////////////
	dout.WriteBlob(o.ToByteArray())
}

func (this *TxRecord) Read(din *io.DataInputX) *TxRecord {
	ver := din.ReadByte()

	if ver < 10 {
		panic("not supported version TxRecord")
		//throw new RuntimeException("not supported version TxRecord");
	}
	in := io.NewDataInputX(din.ReadBlob())

	this.Txid = in.ReadLong()
	this.EndTime = in.ReadDecimal()
	this.Service = int32(in.ReadDecimal())
	this.Elapsed = int32(in.ReadDecimal())
	this.Error = in.ReadDecimal()
	this.CpuTime = int32(in.ReadDecimal())
	this.Malloc = in.ReadDecimal()

	this.SqlCount = int32(in.ReadDecimal())
	this.SqlTime = int32(in.ReadDecimal())
	this.SqlFetchCount = int32(in.ReadDecimal())
	this.SqlFetchTime = int32(in.ReadDecimal())

	this.HttpcCount = int32(in.ReadDecimal())
	this.HttpcTime = int32(in.ReadDecimal())

	this.Active = in.ReadBool()
	this.StepsDataPos = in.ReadDecimal()
	this.Cipher = int32(in.ReadDecimal())

	this.IpAddr = in.ReadInt()
	this.WClientId = in.ReadDecimal()
	this.UserAgent = int32(in.ReadDecimal())
	this.Referer = int32(in.ReadDecimal())
	this.Status = int32(in.ReadDecimal())

	if in.ReadByte() > 0 {
		this.Mtid = in.ReadDecimal()
		this.Mdepth = int32(in.ReadDecimal())
		this.Mcaller = in.ReadDecimal()
	}
	switch in.ReadByte() {
	case 1:
		this.McallerPcode = in.ReadDecimal()
		break
	case 3:
		this.McallerPcode = in.ReadDecimal()
		this.McallerSpec = int32(in.ReadDecimal())
		this.McallerUrl = int32(in.ReadDecimal())
		break
	case 4:
		this.McallerPcode = in.ReadDecimal()
		this.McallerSpec = int32(in.ReadDecimal())
		this.McallerUrl = int32(in.ReadDecimal())
		this.MthisSpec = int32(in.ReadDecimal())
		break
	case 5:
		this.McallerPcode = in.ReadDecimal()
		this.McallerOid = int32(in.ReadDecimal())
		this.McallerSpec = int32(in.ReadDecimal())
		this.McallerUrl = int32(in.ReadDecimal())
		this.MthisSpec = int32(in.ReadDecimal())
		break
	case 6:
		this.McallerPcode = in.ReadDecimal()
		this.McallerOkind = int32(in.ReadDecimal())
		this.McallerOid = int32(in.ReadDecimal())
		this.McallerSpec = int32(in.ReadDecimal())
		this.McallerUrl = int32(in.ReadDecimal())
		this.MthisSpec = int32(in.ReadDecimal())
		break

	}
	this.HttpMethod = in.ReadByte()
	this.Domain = int32(in.ReadDecimal())

	n := int(in.ReadByte())
	if n > 0 {
		this.Fields = value.NewMapValue()
		for i := 0; i < n; i++ {
			key := in.ReadText()
			v := value.ReadValue(in)
			this.Fields.Put(key, v)
		}
	}
	//if in.Available() > 0  {
	this.Login = int32(in.ReadDecimal())

	//if in.Available() > 0  {
	this.ErrorLevel = in.ReadByte()
	if this.ErrorLevel == 0 {
		if this.Error != 0 {
			this.ErrorLevel = WARNING
		}
	}

	//if in.Available() > 0  {
	this.Oid = int32(in.ReadDecimal())
	this.Okind = int32(in.ReadDecimal())
	this.Onode = int32(in.ReadDecimal())

	// TODO
	//if in.Available() > 0  {
	this.Uuid = in.ReadText()

	//if in.Available() > 0  {
	this.DbcTime = int32(in.ReadDecimal())

	//if in.Available() > 0  {
	this.Apdex = in.ReadByte()

	//if (in.available() > 0) {
	this.McallerStepId = in.ReadDecimal()
	this.OriginUrl = in.ReadText()

	//if (in.available() > 0) {
	this.StepSplitCount = int(in.ReadDecimal())

	return this
}

func (this *TxRecord) ToBytes() []byte {
	o := io.NewDataOutputX()
	this.Write(o)
	return o.ToByteArray()
}

func (this *TxRecord) ToObject(b []byte) *TxRecord {

	this.Read(io.NewDataInputX(b))
	return this
}
