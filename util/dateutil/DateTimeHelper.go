package dateutil

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	MILLIS_PER_SECOND      = 1000
	MILLIS_PER_MINUTE      = 60 * MILLIS_PER_SECOND
	MILLIS_PER_FIVE_MINUTE = 5 * 60 * MILLIS_PER_SECOND
	MILLIS_PER_TEN_MINUTE  = 10 * MILLIS_PER_MINUTE
	MILLIS_PER_HOUR        = 60 * MILLIS_PER_MINUTE
	MILLIS_PER_DAY         = 24 * MILLIS_PER_HOUR
)

type Day struct {
	yyyy int
	mm   int
	dd   int
	date string
	wday string
	time int64
}
type DateTimeHelper struct {
	BASE_TIME int64
	table     [][][]*Day
	dateTable []*Day
	LAST_DATE int
}

var wday = []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
var mdayLen = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func newDateTimeHelper(location string) *DateTimeHelper {
	d := new(DateTimeHelper)
	d.BASE_TIME = 0
	d.table = make([][][]*Day, 100)
	d.dateTable = make([]*Day, 40000)
	d.LAST_DATE = 0

	// BASE_TIME 생성 location 이 없을 경우 time.UTC 로 설정.
	if location != "" {
		loc, _ := time.LoadLocation(location)
		//t, _ := time.ParseInLocation(time.RFC3339, "2001-01-01T00:00:00Z", loc)
		t := time.Date(2000, time.January, 1, 0, 0, 0, 0, loc)
		d.BASE_TIME = t.Unix() * 1000
	} else {
		//t, _ := time.Parse(time.RFC3339, "2001-01-01T00:00:00Z")
		t := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		d.BASE_TIME = t.Unix() * 1000
	}
	d.open()

	return d
}

var lock sync.Mutex
var _table = make(map[string]*DateTimeHelper)

func getDateTimeHelper(timezone string) *DateTimeHelper {
	lock.Lock()
	defer lock.Unlock()

	helper := _table[timezone]
	if helper == nil {
		helper = newDateTimeHelper(timezone)
		_table[timezone] = helper
	}
	return helper
}

func (this *DateTimeHelper) open() {
	mtime := this.BASE_TIME
	seq := 0
	wdayIdx := 5 // 20000101:Saturday
	for year := 0; year < 100; year += 1 {
		this.table[year] = make([][]*Day, 12)
		for mm := 0; mm < 12; mm++ {
			monLen := mdayLen[mm]
			if mm == 1 && isYun(year) {
				monLen++
			}
			this.table[year][mm] = make([]*Day, monLen)
			for dd := 0; dd < monLen; dd++ {
				yyyyMMdd := fmt.Sprintf("%d%02d%02d", (year + 2000), mm+1, dd+1)
				this.dateTable[seq] = new(Day)
				this.dateTable[seq].date = yyyyMMdd
				this.dateTable[seq].yyyy = year + 2000
				this.dateTable[seq].mm = mm + 1
				this.dateTable[seq].dd = dd + 1
				this.dateTable[seq].wday = wday[wdayIdx]
				this.dateTable[seq].time = mtime
				this.table[year][mm][dd] = this.dateTable[seq]

				if wdayIdx == 6 {
					wdayIdx = 0
				} else {
					wdayIdx += 1
				}
				seq += 1
				mtime += MILLIS_PER_DAY
			}
		}
	}
}
func isYun(year int) bool {
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	} else {
		return false
	}
}
func (this *DateTimeHelper) getYmdTime(yyyyMMdd string) int64 {
	if yyyyMMdd == "" {
		return 0
	}
	year, _ := strconv.Atoi(yyyyMMdd[0:4])
	mm, _ := strconv.Atoi(yyyyMMdd[4:6])
	dd, _ := strconv.Atoi(yyyyMMdd[6:8])
	if year >= 2100 {
		return this.table[99][11][30].time + MILLIS_PER_DAY
	}
	if year < 2000 {
		return this.BASE_TIME
	}
	return this.table[year-2000][mm-1][dd-1].time
}
func (this *DateTimeHelper) getWeekDay(yyyyMMdd string) string {
	if yyyyMMdd == "" {
		return ""
	}
	year, _ := strconv.Atoi(yyyyMMdd[0:4])
	mm, _ := strconv.Atoi(yyyyMMdd[4:6])
	dd, _ := strconv.Atoi(yyyyMMdd[6:8])
	if year >= 2100 {
		return this.table[99][11][30].wday
	}
	if year < 2000 {
		return this.table[0][0][0].wday
	}
	return this.table[year-2000][mm-1][dd-1].wday
}
func (this *DateTimeHelper) hhmmTo(date string) int64 {
	if date == "" {
		return 0
	}
	h, _ := strconv.Atoi(date[0:2])
	m, _ := strconv.Atoi(date[2:4])
	return int64(h*MILLIS_PER_HOUR + m*MILLIS_PER_MINUTE)
}
func (this *DateTimeHelper) yyyymmdd(time int64) string {
	idx := int((time - this.BASE_TIME) / MILLIS_PER_DAY)
	if idx < 0 {
		idx = 0
	}
	return this.dateTable[idx].date
}
func (this *DateTimeHelper) weekday(time int64) string {
	idx := int((time - this.BASE_TIME) / MILLIS_PER_DAY)
	if idx < 0 {
		idx = 0
	}
	return this.dateTable[idx].wday
}
func (this *DateTimeHelper) datetime(time int64) string {
	if time < this.BASE_TIME {
		return "20000101 00:00:00"
	}
	idx := (int)((time - this.BASE_TIME) / MILLIS_PER_DAY)
	dtime := (int)((time - this.BASE_TIME) % MILLIS_PER_DAY)
	hh := (int)(dtime / MILLIS_PER_HOUR)
	dtime = (int)(dtime % MILLIS_PER_HOUR)
	mm := (int)(dtime / MILLIS_PER_MINUTE)
	dtime = (int)(dtime % MILLIS_PER_MINUTE)
	ss := (int)(dtime / MILLIS_PER_SECOND)

	var buffer bytes.Buffer
	buffer.WriteString(this.dateTable[idx].date)
	buffer.WriteString(" ")
	buffer.WriteString(mk2(hh))
	buffer.WriteString(":")
	buffer.WriteString(mk2(mm))
	buffer.WriteString(":")
	buffer.WriteString(mk2(ss))
	return buffer.String()
}
func (this *DateTimeHelper) timestamp(time int64) string {
	if time < this.BASE_TIME {
		return "20000101 00:00:00"
	}
	idx := (int)((time - this.BASE_TIME) / MILLIS_PER_DAY)
	dtime := (int)((time - this.BASE_TIME) % MILLIS_PER_DAY)
	hh := (int)(dtime / MILLIS_PER_HOUR)
	dtime = (int)(dtime % MILLIS_PER_HOUR)
	mm := (int)(dtime / MILLIS_PER_MINUTE)
	dtime = (int)(dtime % MILLIS_PER_MINUTE)
	ss := (int)(dtime / MILLIS_PER_SECOND)
	sss := (int)(dtime % 1000)

	var buffer bytes.Buffer
	buffer.WriteString(this.dateTable[idx].date)
	buffer.WriteString(" ")
	buffer.WriteString(mk2(hh))
	buffer.WriteString(":")
	buffer.WriteString(mk2(mm))
	buffer.WriteString(":")
	buffer.WriteString(mk2(ss))
	buffer.WriteString(".")
	buffer.WriteString(mk2(sss))
	return buffer.String()
}
func (this *DateTimeHelper) getDateMillis(time int64) int {
	if time < this.BASE_TIME {
		return 0
	}
	dtime := (time - this.BASE_TIME) % MILLIS_PER_DAY
	return int(dtime)
}
func (this *DateTimeHelper) getDateStartTime(time int64) int64 {
	if time < this.BASE_TIME {
		return 0
	}
	dtime := (time - this.BASE_TIME) % MILLIS_PER_DAY
	return time - dtime
}
func (this *DateTimeHelper) getHour(time int64) int {
	return this.getDateMillis(time) / MILLIS_PER_HOUR
}
func (this *DateTimeHelper) getMM(time int64) int {
	dtime := this.getDateMillis(time) % MILLIS_PER_HOUR
	return (dtime / MILLIS_PER_MINUTE)
}
func (this *DateTimeHelper) getTimeUnit(time int64) int64 {
	return (time - this.BASE_TIME)
}
func (this *DateTimeHelper) getDateUnit(time int64) int64 {
	return (time - this.BASE_TIME) / MILLIS_PER_DAY
}
func (this *DateTimeHelper) getTenMinUnit(time int64) int64 {
	return (time - this.BASE_TIME) / MILLIS_PER_TEN_MINUTE
}
func (this *DateTimeHelper) getFiveMinUnit(time int64) int64 {
	return (time - this.BASE_TIME) / MILLIS_PER_FIVE_MINUTE
}
func (this *DateTimeHelper) getMinUnit(time int64) int64 {
	return (time - this.BASE_TIME) / MILLIS_PER_MINUTE
}
func (this *DateTimeHelper) getHourUnit(time int64) int64 {
	return (time - this.BASE_TIME) / MILLIS_PER_HOUR
}
func (this *DateTimeHelper) reverseHourUnit(unit int64) int64 {
	return (unit * MILLIS_PER_HOUR) + this.BASE_TIME
}
func (this *DateTimeHelper) reverseUnit(unit int64, millis int64) int64 {
	return (unit * millis) + this.BASE_TIME
}

func mk2(n int) string {
	switch n {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
func mk3(n int) string {
	switch n {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		return "00" + strconv.Itoa(n)
	}
	if n < 100 {
		return "0" + strconv.Itoa(n)
	} else {
		return strconv.Itoa(n)
	}
}
func (this *DateTimeHelper) ymdhms(time int64) string {
	if time < this.BASE_TIME {
		return "20000101000000"
	}
	idx := (int)((time - this.BASE_TIME) / MILLIS_PER_DAY)
	dtime := (int)((time - this.BASE_TIME) % MILLIS_PER_DAY)
	hh := (int)(dtime / MILLIS_PER_HOUR)
	dtime = (int)(dtime % MILLIS_PER_HOUR)
	mm := (int)(dtime / MILLIS_PER_MINUTE)
	dtime = (int)(dtime % MILLIS_PER_MINUTE)
	ss := (int)(dtime / MILLIS_PER_SECOND)
	//sss := (int)(dtime % 1000)

	var buffer bytes.Buffer
	buffer.WriteString(this.dateTable[idx].date)
	buffer.WriteString(mk2(hh))
	buffer.WriteString(mk2(mm))
	buffer.WriteString(mk2(ss))
	return buffer.String()
}
func (this *DateTimeHelper) getYear(time int64, delta int) int64 {
	if time < this.BASE_TIME {
		return this.BASE_TIME
	}
	idx := (int)((time - this.BASE_TIME) / MILLIS_PER_DAY)
	year := this.dateTable[idx].yyyy + delta
	mm := this.dateTable[idx].mm
	dd := this.dateTable[idx].dd
	if year < 2000 {
		return this.BASE_TIME
	}
	if year > 2100 {
		return this.dateTable[this.LAST_DATE].time
	}

	dtime := (time - this.BASE_TIME) % MILLIS_PER_DAY

	return this.table[year-2000][mm-1][dd-1].time + dtime
}
func (this *DateTimeHelper) getMonth(time int64, delta int) int64 {
	if time < this.BASE_TIME {
		return this.BASE_TIME
	}
	idx := (int)((time - this.BASE_TIME) / MILLIS_PER_DAY)
	year := this.dateTable[idx].yyyy
	mm := this.dateTable[idx].mm
	dd := this.dateTable[idx].dd
	delta = delta + mm - 1
	deltaYear := delta / 12
	deltaMM := delta % 12
	if deltaMM < 0 {
		deltaYear--
		deltaMM = 12 + deltaMM
	}
	year += deltaYear
	mm = deltaMM + 1
	if year < 2000 {
		return this.BASE_TIME
	}
	if year > 2100 {
		return this.dateTable[this.LAST_DATE].time
	}

	dtime := (time - this.BASE_TIME) % MILLIS_PER_DAY
	return this.table[year-2000][mm-1][dd-1].time + dtime
}
func (this *DateTimeHelper) getDate(time int64, delta int) int64 {
	if time < this.BASE_TIME {
		return this.BASE_TIME
	}
	idx := (int)((time - this.BASE_TIME) / MILLIS_PER_DAY)
	dtime := (time - this.BASE_TIME) % MILLIS_PER_DAY
	idx += delta
	if idx < 0 {
		return this.BASE_TIME
	}
	if idx >= this.LAST_DATE {
		return this.dateTable[this.LAST_DATE].time
	}
	return this.dateTable[idx].time + dtime
}
func (this *DateTimeHelper) logtime(time int64) string {
	if time < this.BASE_TIME {
		return "00:00:00.000"
	}
	//idx := (int)((time - this.BASE_TIME) / MILLIS_PER_DAY)
	dtime := (int)((time - this.BASE_TIME) % MILLIS_PER_DAY)
	hh := (int)(dtime / MILLIS_PER_HOUR)
	dtime = (int)(dtime % MILLIS_PER_HOUR)
	mm := (int)(dtime / MILLIS_PER_MINUTE)
	dtime = (int)(dtime % MILLIS_PER_MINUTE)
	ss := (int)(dtime / MILLIS_PER_SECOND)
	sss := (int)(dtime % 1000)

	var buffer bytes.Buffer
	buffer.WriteString(mk2(hh))
	buffer.WriteString(":")
	buffer.WriteString(mk2(mm))
	buffer.WriteString(":")
	buffer.WriteString(mk2(ss))
	buffer.WriteString(".")
	buffer.WriteString(mk2(sss))
	return buffer.String()
}
func (this *DateTimeHelper) hhmmss(time int64) string {
	if time < this.BASE_TIME {
		return "000000"
	}
	dtime := (int)((time - this.BASE_TIME) % MILLIS_PER_DAY)
	hh := (int)(dtime / MILLIS_PER_HOUR)
	dtime = (int)(dtime % MILLIS_PER_HOUR)
	mm := (int)(dtime / MILLIS_PER_MINUTE)
	dtime = (int)(dtime % MILLIS_PER_MINUTE)
	ss := (int)(dtime / MILLIS_PER_SECOND)
	return fmt.Sprintf("%02d%02d%02d", hh, mm, ss)
}
func (this *DateTimeHelper) hhmm(time int64) string {
	if time < this.BASE_TIME {
		return "0000"
	}
	dtime := (int)((time - this.BASE_TIME) % MILLIS_PER_DAY)
	hh := (int)(dtime / MILLIS_PER_HOUR)
	dtime = (int)(dtime % MILLIS_PER_HOUR)
	mm := (int)(dtime / MILLIS_PER_MINUTE)
	return fmt.Sprintf("%02d%02d", hh, mm)
}
