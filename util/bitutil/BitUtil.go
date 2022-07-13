package bitutil
/*
  TODO : 검증 테스트해야함.
*/
func Composite64(hkey, wkey int32) int64 {
	return (int64(hkey) << 32) | (int64(wkey) & 0xffffffff)
}
func Composite32(hkey, wkey int16) int32 {
	return (int32(hkey) << 16) | (int32(wkey) & 0xffff)
}
func Composite16(hkey, wkey byte) int16 {
	return (int16(hkey) << 8) | (int16(wkey) & int16(0xff))
}

func SetHigh64(src int64, hkey int32) int64 {
	return (src & 0x00000000ffffffff) | (int64(hkey) << 32)
}
func SetLow64(src int64, wkey int32) int64 {
	var x uint64 = 0xffffffff00000000
	return (src & int64(x)) | (int64(wkey) & 0xffffffff)
}
func GetHigh64(key int64) int32 {
	var x uint32 = 0xffffffff
	return int32(key>>32) & int32(x)
}
func GetLow64(key int64) int32 {
	var x uint32 = 0xffffffff
	return int32(key) & int32(x)
}
func GetHigh32(key int32) int16 {
	var x uint16 = 0xffff
	return int16(key >> 16) & int16(x)
}
func GetLow32(key int32) int16 {
	return int16(key & 0xffff)
}

func GetHigh16(key int16) byte {
	return byte((key >> 8) & int16(0xff))
}
func GetLow16(key int16) byte {
	return byte(key & int16(0xff))
}
