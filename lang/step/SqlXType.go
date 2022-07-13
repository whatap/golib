package step

const (
	SQL_XTYPE_STMT     = 0x0
	SQL_XTYPE_PREPARED = 0x1

	SQL_XTYPE_METHOD_KNOWN   = 0x00
	SQL_XTYPE_METHOD_EXECUTE = 0x10
	SQL_XTYPE_METHOD_UPDATE  = 0x20
	SQL_XTYPE_METHOD_QUERY   = 0x30
)

func GetSqlClassType(xtype byte) byte {
	return xtype & 0x0f
}

func GetSqlMethodType(xtype byte) byte {
	return xtype & 0xf0
}

func GetSqlClassString(xtype byte) string {
	switch GetSqlClassType(xtype) {
	case SQL_XTYPE_STMT:
		return "Statement"
	case SQL_XTYPE_PREPARED:
		return " PreparedStatement"
	}
	return "Unknown"
}
