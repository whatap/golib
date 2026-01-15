package pack

import (
	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/pack/open"
)

const (
	// 256
	PACK_PARAMETER = 0x0100
	// 512
	//PACK_COUNTER       	= 0x0200
	// 513
	PACK_COUNTER_1 = 0x0201
	// 768
	PACK_PROFILE            = 0x0300
	PACK_BIZMON_HELP        = 0x0301
	PACK_PROFILE_STEP_SPLIT = 0x302
	// 1025
	PACK_ACTIVESTACK_1 = 0x0401
	// 1792
	PACK_TEXT    = 0x0700
	PACK_HV_TEXT = 0x0701
	// 2048
	//PACK_ERROR_SNAP    	= 0x0800
	// 2049
	PACK_ERROR_SNAP_1 = 0x0801
	// 3840
	PACK_REALTIME_USER   = 0x0f00
	PACK_REALTIME_USER_1 = 0x0f01

	// 2304
	PACK_STAT_SERVICE   = 0x0900
	PACK_STAT_SERVICE_1 = 0x0901
	// 2320
	PACK_STAT_GENERAL   = 0x0910
	PACK_STAT_GENERAL_1 = 0x0911

	// 2560
	PACK_STAT_SQL = 0x0a00
	// 2816
	PACK_STAT_HTTPC = 0x0b00
	// 3072
	PACK_STAT_ERROR = 0x0c00
	// 3584
	PACK_STAT_METHOD = 0x0e00
	// 4096
	PACK_STAT_TOP_SERVICE = 0x1000
	// 4352
	PACK_STAT_REMOTE_IP = 0x1100
	// 4608
	PACK_STAT_USER_AGENT   = 0x1200
	PACK_STAT_USER_AGENT_1 = 0x1201

	// 5120
	PACK_EVENT = 0x1400
	// 5376
	//PACK_HITMAP    		= 0x1500
	// 5377
	PACK_HITMAP_1 = 0x1501
	HITVIEW       = 0x1506
	// 5632
	PACK_EXTENSION = 0x1600
	TAG_COUNT      = 0x1601
	TAG_LOG        = 0x1602

	// To avoid cycle import errors, move open.Packenum
	// PACK_OPEN_MX_PACK      = 0x1603
	// PACK_OPEN_MX_HELP_PACK = 0x1604

	// 5888
	PACK_COMPOSITE = 0x1700

	// 5889
	PACK_BSM_RECORD = 0x1701
	// 5890
	PACK_AP_NUT = 0x1702
	// 5891
	PACK_ADDIN_COUNT = 0x1703

	PACK_KUBE_MASTER_COUNT = 0x1704
	PACK_KUBE_MASTER_STAT  = 0x1705
	PACK_KUBE_NODE         = 0x1706
	PACK_WEB_CHECK_COUNT   = 0x1707
	PACK_LOGSINK           = 0x170a
	PACK_ZIP               = 0x170b
	PACK_AGENT_MAPPING     = 0x170c
	PACK_LOGSINK_ZIP       = 0x170d
	PACK_AGENT_PROPERTY    = 0x170e

	//	// 12288
	//	PACK_SM_BASE       = 0x3000
	//	// 12289
	//	PACK_SM_DISK_QUATA = 0x3001
	//	// 12290
	//	PACK_SM_NET_PERF   = 0x3002
	//	// 12300
	//	PACK_SM_PROC_PERF  = 0x3003
	//	// 12301
	//	PACK_SM_PORT_PERF  = 0x3004
	//	// 12302
	//	PACK_SM_LOG_EVENT  = 0x3005

	PACK_SM_BASE       = 0x3008
	PACK_SM_DISK_QUATA = 0x3001
	PACK_SM_NET_PERF   = 0x3002
	PACK_SM_PROC_PERF  = 0x3003
	PACK_SM_PORT_PERF  = 0x3004
	PACK_SM_LOG_EVENT  = 0x3005
	PACK_SM_DOWN_CHECK = 0x3006
	PACK_SM_PROC_GROUP = 0x3007
	PACK_SM_BASE_1     = 0x3008
	PACK_SM_META       = 0x3009

	PACK_SM_EXTENSION = 0x1600

	PACK_SM_NUT  = 0x3010
	PACK_SM_ATTR = 0x3011
	PACK_SM_PING = 0x3012

	ADDIN_VTYPE_CONTAINER = "container"

	PACK_SERVERINFO = 0x6500
)

type Pack interface {
	GetPackType() int16
	Write(out *io.DataOutputX)
	Read(in *io.DataInputX)

	// OID, PCODE 설정을 위한 함수
	SetOID(oid int32)
	SetPCODE(pcode int64)
	GetPCODE() int64
	SetOKIND(okind int32)
	SetONODE(onode int32)

	// Time interface 추가
	GetTime() int64
	SetTime(t int64)
}

func CreatePack(t int16) Pack {
	switch t {
	case PACK_PARAMETER:
		return NewParamPack()
	case PACK_COUNTER_1:
		return NewCounterPack1()
	case PACK_PROFILE:
		return NewProfilePack()
	case PACK_ACTIVESTACK_1:
		return NewActiveStackPack()
	case PACK_TEXT:
		return NewTextPack()
	case PACK_ERROR_SNAP_1:
		return NewErrorSnapPack1()
	case PACK_REALTIME_USER:
		return NewRealtimeUserPack()

	case PACK_STAT_SERVICE:
		return NewStatServicePack()
	case PACK_STAT_GENERAL:
		return NewStatGeneralPack()
	case PACK_STAT_SQL:
		return NewStatSqlPack()
	case PACK_STAT_HTTPC:
		return NewStatHttpcPack()
	case PACK_STAT_ERROR:
		return NewStatErrorPack()
	// case PACK_STAT_METHOD:
	// 	return NewStatMethodPack()
	// case PACK_STAT_TOP_SERVICE
	// return NewStatTopServicePack()
	case PACK_STAT_REMOTE_IP:
		return NewStatRemoteIpPack()
	case PACK_STAT_USER_AGENT:
		return NewStatUserAgentPack()
	case PACK_EVENT:
		return NewEventPack()
	case PACK_HITMAP_1:
		return NewHitMapPack1()
	case PACK_EXTENSION:
		return NewExtensionPack()
	case TAG_COUNT:
		return NewTagCountPack()
	case TAG_LOG:
		return NewTagLogPack()
	case open.PACK_OPEN_MX_HELP_PACK:
		return open.NewOpenMxHelpPack()
	case open.PACK_OPEN_MX_PACK:
		return open.NewOpenMxPack()
	case PACK_COMPOSITE:
		return NewCompositePack()
	// case PACK_ADDIN_COUNT:
	// 	return NewAddinCountPack()
	case PACK_LOGSINK:
		return NewLogSinkPack()
	case PACK_ZIP:
		return NewZipPack()
	case PACK_LOGSINK_ZIP:
		return NewLogSinkZipPack()
	case PACK_AGENT_PROPERTY:
		return NewAgentPropertyPack()

	case PACK_SERVERINFO:
		return NewServerInfoPack()
	}

	return nil
}

func WritePack(out *io.DataOutputX, p Pack) *io.DataOutputX {
	out.WriteShort(int16(p.GetPackType()))
	p.Write(out)
	return out
}
func ReadPack(in *io.DataInputX) Pack {
	t := in.ReadShort()
	v := CreatePack(t)
	v.Read(in)
	return v
}

func ToBytesPack(p Pack) []byte {
	out := io.NewDataOutputX()
	WritePack(out, p)
	return out.ToByteArray()
}
func ToPack(b []byte) Pack {
	in := io.NewDataInputX(b)
	return ReadPack(in)
}
func ToBytesPackECB(p Pack, fmtLen int) []byte {
	out := io.NewDataOutputX()
	WritePack(out, p)
	remainder := out.Size() % fmtLen
	if remainder != 0 {
		b := make([]byte, fmtLen-remainder)
		out.Write(b, 0, len(b))
	}
	return out.ToByteArray()
}

func GetPackTypeString(t int16) string {
	switch t {
	case PACK_PARAMETER:
		return "ParamPack"
	case PACK_COUNTER_1:
		return "CounterPack1"
	case PACK_PROFILE:
		return "ProfilePack"
	case PACK_ACTIVESTACK_1:
		return "ActiveStackPack"
	case PACK_TEXT:
		return "TextPack"
	case PACK_ERROR_SNAP_1:
		return "ErrorSnapPack1"
	case PACK_REALTIME_USER:
		return "RealtimeUserPack"

	case PACK_STAT_SERVICE:
		return "StatServicePack"
	case PACK_STAT_GENERAL:
		return "StatGeneralPack"
	case PACK_STAT_SQL:
		return "StatSqlPack"
	case PACK_STAT_HTTPC:
		return "StatHttpcPack"
	case PACK_STAT_ERROR:
		return "StatErrorPack"
	//	case PACK_STAT_METHOD:
	//	return NewStatMethodPack()
	//	case PACK_STAT_TOP_SERVICE
	//	return NewStatTopServicePack()
	case PACK_STAT_REMOTE_IP:
		return "StatRemoteIpPack"
	case PACK_STAT_USER_AGENT:
		return "StatUserAgentPack"
	case PACK_EVENT:
		return "EventPack"
	case PACK_HITMAP_1:
		return "HitMapPack1"
	case PACK_EXTENSION:
		return "ExtensionPack"
	case TAG_COUNT:
		return "TagCountPack"
		//		case TAG_LOG:
		//		return NewTagLogPack()
	case open.PACK_OPEN_MX_HELP_PACK:
		return "OpenMxHelpPack"
	case open.PACK_OPEN_MX_PACK:
		return "OpenMxPack"
	case PACK_COMPOSITE:
		return "CompositePack"
		//	case PACK_BSM_RECORD:
		//		return NewBsmRecordPack()
		//	case PACK_AP_NUT:
		//	return NewApNutPack()
	case PACK_ADDIN_COUNT:
		return "AddinCountPack"
	case PACK_LOGSINK:
		return "LogSinkPack"
	case PACK_ZIP:
		return "ZipPack"
	case PACK_LOGSINK_ZIP:
		return "LogsinkZipPack"
	case PACK_AGENT_PROPERTY:
		return "AgentPropertyPack"

	}
	return "Unknown"
}
