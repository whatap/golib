package udp

import (
	"sync"

	"github.com/whatap/golib/io"
)

const (
	UDP_PACK_VERSION = 50100

	TX_BLANK uint8 = 0
	TX_START uint8 = 1

	TX_DB_CONN   uint8 = 2
	TX_DB_FETCH  uint8 = 3
	TX_SQL       uint8 = 4
	TX_SQL_START uint8 = 5
	TX_SQL_END   uint8 = 6

	TX_HTTPC       uint8 = 7
	TX_HTTPC_START uint8 = 8
	TX_HTTPC_END   uint8 = 9

	TX_ERROR  uint8 = 10
	TX_MSG    uint8 = 11
	TX_METHOD uint8 = 12

	// secure msg
	TX_SECURE_MSG uint8 = 13

	// sql & param
	TX_SQL_PARAM uint8 = 14

	TX_PARAM     uint8 = 30
	ACTIVE_STACK uint8 = 40
	ACTIVE_STATS uint8 = 41
	DBCONN_POOL  uint8 = 42

	// golang config
	CONFIG_INFO uint8 = 230

	// relay pack
	RELAY_PACK uint8 = 244

	TX_START_END            uint8 = 254
	TX_END                  uint8 = 255
	UDP_PACKET_SQL_MAX_SIZE       = 32768
)

const (
	PACKET_DB_MAX_SIZE           = 4 * 1024  // max size of sql
	PACKET_SQL_MAX_SIZE          = 32 * 1024 // max size of sql
	PACKET_HTTPC_MAX_SIZE        = 32 * 1024 // max size of sql
	PACKET_MESSAGE_MAX_SIZE      = 32 * 1024 // max size of message
	PACKET_METHOD_STACK_MAX_SIZE = 32 * 1024 // max size of message

	COMPILE_FILE_MAX_SIZE = 2 * 1024 // max size of filename

	HTTP_HOST_MAX_SIZE   = 2 * 1024 // max size of host
	HTTP_URI_MAX_SIZE    = 2 * 1024 // max size of uri
	HTTP_METHOD_MAX_SIZE = 256      // max size of method
	HTTP_IP_MAX_SIZE     = 256      // max size of ip(request_addr)
	HTTP_UA_MAX_SIZE     = 2 * 1024 // max size of user agent
	HTTP_REF_MAX_SIZE    = 2 * 1024 // max size of referer
	HTTP_USERID_MAX_SIZE = 2 * 1024 // max size of userid

	HTTP_PARAM_MAX_COUNT      = 20
	HTTP_PARAM_KEY_MAX_SIZE   = 255 // = 을 빼고 255 byte
	HTTP_PARAM_VALUE_MAX_SIZE = 256

	HTTP_HEADER_MAX_COUNT      = 20
	HTTP_HEADER_KEY_MAX_SIZE   = 255 // = 을 빼고 255 byte
	HTTP_HEADER_VALUE_MAX_SIZE = 256

	SQL_PARAM_MAX_COUNT      = 20
	SQL_PARAM_VALUE_MAX_SIZE = 256

	STEP_ERROR_MESSAGE_MAX_SIZE = 4 * 1024
)

type UdpPack interface {
	GetPackType() uint8
	Write(out *io.DataOutputX)
	Read(in *io.DataInputX)

	SetVersion(ver int32)
	GetVersion() int32

	SetFlush(flush bool)
	IsFlush() bool

	// Processing data
	Process()

	Clear()
}

var udpStartPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxStartPack()
	},
}
var udpStartEndPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxStartEndPack()
	},
}
var udpEndPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxEndPack()
	},
}
var udpSqlPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxSqlPack()
	},
}
var udpSqlParamPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxSqlParamPack()
	},
}
var udpHttpcPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxHttpcPack()
	},
}
var udpErrorPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxErrorPack()
	},
}
var udpMessagePool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxMessagePack()
	},
}
var udpSecureMessagePool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxSecureMessagePack()
	},
}
var udpMethodPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxMethodPack()
	},
}
var udpDbcPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxDbcPack()
	},
}
var udpRelayPool = sync.Pool{
	New: func() interface{} {
		return NewUdpRelayPack()
	},
}
var udpActiveStackPool = sync.Pool{
	New: func() interface{} {
		return NewUdpActiveStackPack()
	},
}
var udpTxParamPool = sync.Pool{
	New: func() interface{} {
		return NewUdpTxParamPack()
	},
}
var udpActiveStatsPool = sync.Pool{
	New: func() interface{} {
		return NewUdpActiveStatsPack()
	},
}
var udpDBConPool = sync.Pool{
	New: func() interface{} {
		return NewUdpDBConPoolPack()
	},
}
var udpConfigPool = sync.Pool{
	New: func() interface{} {
		return NewUdpConfigPack()
	},
}

func CreatePack(t uint8, ver int32) UdpPack {
	switch t {
	case TX_START:
		p := udpStartPool.Get().(*UdpTxStartPack)
		p.Ver = ver
		return p
		//return NewUdpTxStartPackVer(ver)
	case TX_START_END:
		p := udpStartEndPool.Get().(*UdpTxStartEndPack)
		p.Ver = ver
		return p
		//return NewUdpTxStartEndPackVer(ver)
	case TX_END:
		p := udpEndPool.Get().(*UdpTxEndPack)
		p.Ver = ver
		return p
		//return NewUdpTxEndPackVer(ver)
	case TX_SQL:
		p := udpSqlPool.Get().(*UdpTxSqlPack)
		p.Ver = ver
		return p
	//return NewUdpTxSqlPackVer(ver)
	case TX_SQL_PARAM:
		p := udpSqlParamPool.Get().(*UdpTxSqlParamPack)
		p.Ver = ver
		return p
		//return NewUdpTxSqlParamPackVer(ver)
	case TX_HTTPC:
		p := udpHttpcPool.Get().(*UdpTxHttpcPack)
		p.Ver = ver
		return p
		//return NewUdpTxHttpcPackVer(ver)
	case TX_ERROR:
		p := udpErrorPool.Get().(*UdpTxErrorPack)
		p.Ver = ver
		return p
		//return NewUdpTxErrorPackVer(ver)
	case TX_MSG:
		p := udpMessagePool.Get().(*UdpTxMessagePack)
		p.Ver = ver
		return p
		//return NewUdpTxMessagePackVer(ver)
	case TX_SECURE_MSG:
		p := udpSecureMessagePool.Get().(*UdpTxSecureMessagePack)
		p.Ver = ver
		return p
		//return NewUdpTxSecureMessagePackVer(ver)
	case TX_METHOD:
		p := udpMethodPool.Get().(*UdpTxMethodPack)
		p.Ver = ver
		return p
		//return NewUdpTxMethodPackVer(ver)
	case TX_DB_CONN:
		p := udpDbcPool.Get().(*UdpTxDbcPack)
		p.Ver = ver
		return p
		//return NewUdpTxDbcPackVer(ver)
	case RELAY_PACK:
		p := udpRelayPool.Get().(*UdpRelayPack)
		p.Ver = ver
		return p
		//return NewUdpRelayPackVer(ver)
	case ACTIVE_STACK:
		p := udpActiveStackPool.Get().(*UdpActiveStackPack)
		p.Ver = ver
		return p
		//return NewUdpActiveStackPackVer(ver)
	case TX_PARAM:
		p := udpTxParamPool.Get().(*UdpTxParamPack)
		p.Ver = ver
		return p
		//return NewUdpTxParamPackVer(ver)
	case ACTIVE_STATS:
		p := udpActiveStatsPool.Get().(*UdpActiveStatsPack)
		p.Ver = ver
		return p
		//return NewUdpActiveStatsPackVer(ver)
	case DBCONN_POOL:
		p := udpDBConPool.Get().(*UdpDBConPoolPack)
		p.Ver = ver
		return p
		//return NewUdpDBConPoolPackVer(ver)
	case CONFIG_INFO:
		p := udpConfigPool.Get().(*UdpConfigPack)
		p.Ver = ver
		return p
		//return NewUdpConfigPackVer(ver)
	}
	return nil
}
func ClosePack(p UdpPack) {
	p.Clear()
	switch p.GetPackType() {
	case TX_START:
		udpStartPool.Put(p)
	case TX_START_END:
		udpStartEndPool.Put(p)
	case TX_END:
		udpEndPool.Put(p)
	case TX_SQL:
		udpSqlPool.Put(p)
	case TX_SQL_PARAM:
		udpSqlParamPool.Put(p)
	case TX_HTTPC:
		udpHttpcPool.Put(p)
	case TX_ERROR:
		udpErrorPool.Put(p)
	case TX_MSG:
		udpMessagePool.Put(p)
	case TX_SECURE_MSG:
		udpSecureMessagePool.Put(p)
	case TX_METHOD:
		udpMethodPool.Put(p)
	case TX_DB_CONN:
		udpDbcPool.Put(p)
	case RELAY_PACK:
		udpRelayPool.Put(p)
	case ACTIVE_STACK:
		udpActiveStackPool.Put(p)
	case TX_PARAM:
		udpTxParamPool.Put(p)
	case ACTIVE_STATS:
		udpActiveStatsPool.Put(p)
	case DBCONN_POOL:
		udpDBConPool.Put(p)
	case CONFIG_INFO:
		udpConfigPool.Put(p)
	}
}

func WritePack(out *io.DataOutputX, p UdpPack) *io.DataOutputX {
	p.Write(out)
	return out
}

func ReadPack(t uint8, ver int32, in *io.DataInputX) UdpPack {
	v := CreatePack(uint8(t), ver)
	v.Read(in)
	v.Process()
	return v
}

func ToBytesPack(p UdpPack) []byte {
	out := io.NewDataOutputX()
	WritePack(out, p)
	return out.ToByteArray()
}

func ToPack(t uint8, ver int32, b []byte) UdpPack {
	in := io.NewDataInputX(b)
	return ReadPack(t, ver, in)
}
