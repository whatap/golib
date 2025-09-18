package udp

import (
	"fmt"
	"github.com/whatap/golib/util/stringutil"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/paramtext"
)

type UdpTxSqlPack struct {
	AbstractPack
	Dbc string
	Sql string
	//error
	ErrorType    string
	ErrorMessage string

	Stack string
	Fetch int32
}

func NewUdpTxSqlPack() *UdpTxSqlPack {
	p := new(UdpTxSqlPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}

func NewUdpTxSqlPackVer(ver int32) *UdpTxSqlPack {
	p := new(UdpTxSqlPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}

func (this *UdpTxSqlPack) GetPackType() uint8 {
	return TX_SQL
}

func (this *UdpTxSqlPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",dbc=", this.Dbc, ",sql=", this.Sql, ",desc=")
}

func (this *UdpTxSqlPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.Dbc = ""
	this.Sql = ""
	//error
	this.ErrorType = ""
	this.ErrorMessage = ""

	this.Stack = ""
	this.Fetch = 0
}

func (this *UdpTxSqlPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Dbc)
	dout.WriteTextShortLength(this.Sql)

	if this.Ver > 60000 {
		// Node.js
		dout.WriteTextShortLength(string(this.Fetch))
	} else if this.Ver > 50000 {
		// Golang
		dout.WriteTextShortLength(this.ErrorType)
		dout.WriteTextShortLength(this.ErrorMessage)
		dout.WriteTextShortLength(this.Stack)
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
		dout.WriteTextShortLength(this.ErrorType)
		dout.WriteTextShortLength(this.ErrorMessage)
		dout.WriteTextShortLength(this.Stack)
	} else if this.Ver > 20000 {
		// Python
		if this.Ver >= 20102 {
			dout.WriteTextShortLength(string(this.Fetch))
		}
	} else {
		// PHP
		if this.Ver >= 10105 {
			dout.WriteTextShortLength(this.ErrorType)
			dout.WriteTextShortLength(this.ErrorMessage)
			dout.WriteTextShortLength(this.Stack)
		}
	}
}

func (this *UdpTxSqlPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Dbc = din.ReadTextShortLength()
	this.Sql = din.ReadTextShortLength()

	if this.Ver > 60000 {
		// Node.js
		this.Fetch = stringutil.ParseInt32(din.ReadTextShortLength())
	} else if this.Ver > 50000 {
		// Golang
		this.ErrorType = din.ReadTextShortLength()
		this.ErrorMessage = din.ReadTextShortLength()
		this.Stack = din.ReadTextShortLength()
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
		this.ErrorType = din.ReadTextShortLength()
		this.ErrorMessage = din.ReadTextShortLength()
		this.Stack = din.ReadTextShortLength()
	} else if this.Ver > 20000 {
		// Python
		if this.Ver >= 20102 {
			this.Fetch = stringutil.ParseInt32(din.ReadTextShortLength())
		}
	} else {
		// PHP
		if this.Ver >= 10105 {
			this.ErrorType = din.ReadTextShortLength()
			this.ErrorMessage = din.ReadTextShortLength()
			this.Stack = din.ReadTextShortLength()
		}
	}
}

func (this *UdpTxSqlPack) Process() {

	if this.Ver > 60000 {
		// Node.js
	} else if this.Ver > 50000 {
		// Golang
		if this.Dbc != "" {
			p := paramtext.NewParamKVSeperate(this.Dbc, " ", "=")
			this.Dbc = p.ToStringStr("password", "#")
			p = paramtext.NewParamKVSeperate(this.Dbc, ";", "=")
			this.Dbc = p.ToStringStr("password", "#")
		}
		if len(this.Sql) >= UDP_PACKET_SQL_MAX_SIZE {
			this.Sql = "[QUERY TOO LONG]\r\n" + this.Sql
		}
	} else if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
	} else if this.Ver > 20000 {
		// Python
	} else {
		// PHP
		if this.Dbc != "" {
			p := paramtext.NewParamKVSeperate(this.Dbc, " ", "=")
			this.Dbc = p.ToStringStr("password", "#")
			p = paramtext.NewParamKVSeperate(this.Dbc, ";", "=")
			this.Dbc = p.ToStringStr("password", "#")
		}
		if len(this.Sql) >= UDP_PACKET_SQL_MAX_SIZE {
			this.Sql = "[QUERY TOO LONG]\r\n" + this.Sql
		}
	}
}
