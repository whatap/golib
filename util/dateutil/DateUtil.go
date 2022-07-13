package dateutil

import (
	//"log"
	"time"
)

var helper = getDateTimeHelper("")

func DateTime(time int64) string {
	return helper.datetime(time)
}
func TimeStamp(time int64) string {
	return helper.timestamp(time)
}
func TimeStampNow() string {
	return helper.timestamp(Now())
}
func WeekDay(time int64) string {
	return helper.weekday(time)
}
func GetDateUnitNow() int64 {
	return helper.getDateUnit(Now())
}

func GetDateUnit(time int64) int64 {
	return helper.getDateUnit(time)
}
func Ymdhms(time int64) string {
	return helper.ymdhms(time)
}
func YYYYMMDD(time int64) string {
	return helper.yyyymmdd(time)
}
func HHMMSS(time int64) string {
	return helper.hhmmss(time)
}
func HHMM(time int64) string {
	return helper.hhmm(time)
}
func YmdNow() string {
	return helper.yyyymmdd(Now())
}

var delta int64 = 0

func SystemNow() int64 {
	return (time.Now().UnixNano() / 1000000)
}
func Now() int64 {
	t := SystemNow()
	return t + delta
}
func SetDelta(t int64) {
	delta = t
}

func SetServerTime(serverTime int64, syncfactor float64) int64 {
	now := SystemNow()
	delta = serverTime - now
	if delta != 0 {
		delta = int64(float64(delta) * syncfactor)
	}
	return delta
}
func GetDelta() int64 {
	return delta
}

//
//	func timestamp(time int64 ) string{
//		return helper.timestamp(time);
//	}
//
//	func yyyymmdd(time int64 ) string{
//		return helper.yyyymmddStr(time);
//	}

func GetFiveMinUnit(time int64) int64 {
	return helper.getFiveMinUnit(time)
}

func GetMinUnit(time int64) int64 {
	return helper.getMinUnit(time)
}

func GetYmdTime(yyyyMMdd string) int64 {
	return helper.getYmdTime(yyyyMMdd)
}
