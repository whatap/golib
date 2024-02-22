package dateutil

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
)

const (
	DATEFORMAT_YEAR        = 'y'
	DATEFORMAT_MONTH       = 'm'
	DATEFORMAT_DAY         = 'd'
	DATEFORMAT_HOUR        = 'H'
	DATEFORMAT_MINUTE      = 'M'
	DATEFORMAT_SECOND      = 'S'
	DATEFORMAT_MILLISECOND = 's'
)

type DateFormat struct {
	formatStr []rune
	dateStr   string
	date      map[rune]int
}

func NewDateFormat(formatStr string) *DateFormat {
	p := new(DateFormat)
	p.formatStr = []rune(formatStr)
	p.date = make(map[rune]int)
	return p
}

func (this *DateFormat) Format() string {
	return this.format(time.Now())
}
func (this *DateFormat) FormatTime(t time.Time) string {
	return this.format(t)
}

func (this *DateFormat) format(t time.Time) string {
	ret := this.formatStr
	var buf bytes.Buffer
	for _, ch := range ret {
		switch ch {
		case DATEFORMAT_YEAR:
			buf.WriteString(LPadInt(t.Year(), 4))
		case DATEFORMAT_MONTH:
			buf.WriteString(LPadInt(int(t.Month()), 2))
		case DATEFORMAT_DAY:
			buf.WriteString(LPadInt(t.Day(), 2))
		case DATEFORMAT_HOUR:
			buf.WriteString(LPadInt(t.Hour(), 2))
		case DATEFORMAT_MINUTE:
			buf.WriteString(LPadInt(t.Minute(), 2))
		case DATEFORMAT_SECOND:
			buf.WriteString(LPadInt(t.Second(), 2))
		case DATEFORMAT_MILLISECOND:
			buf.WriteString(LPadInt(int((t.UnixNano()/1000000)%1000), 3))
		default:
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func (this *DateFormat) Parse(dateStr string) (int64, error) {
	r := bytes.NewReader([]byte(dateStr))
	sz := len(dateStr)
	now := time.Now()
	for i, ch := range this.formatStr {
		if i >= sz {
			break
		}
		switch ch {
		case DATEFORMAT_YEAR:
			if v, err := this.ToInt(r, 4); err == nil {
				this.date[ch] = v
			} else {
				return 0, fmt.Errorf("Parse error ch=%s, err=%s ", string(ch), err)
			}
		case DATEFORMAT_MONTH:
			if v, err := this.ToInt(r, 2); err == nil {
				this.date[ch] = v
			} else {
				return 0, fmt.Errorf("Parse error ch=%s, err=%s ", string(ch), err)
			}
		case DATEFORMAT_DAY:
			if v, err := this.ToInt(r, 2); err == nil {
				this.date[ch] = v
			} else {
				return 0, fmt.Errorf("Parse error ch=%s, err=%s ", string(ch), err)
			}
		case DATEFORMAT_HOUR:
			if v, err := this.ToInt(r, 2); err == nil {
				this.date[ch] = v
			} else {
				return 0, fmt.Errorf("Parse error ch=%s, err=%s ", string(ch), err)
			}
		case DATEFORMAT_MINUTE:
			if v, err := this.ToInt(r, 2); err == nil {
				this.date[ch] = v
			} else {
				return 0, fmt.Errorf("Parse error ch=%s, err=%s ", string(ch), err)
			}
		case DATEFORMAT_SECOND:
			if v, err := this.ToInt(r, 2); err == nil {
				this.date[ch] = v
			} else {
				return 0, fmt.Errorf("Parse error ch=%s, err=%s ", string(ch), err)
			}
		case DATEFORMAT_MILLISECOND:
			if v, err := this.ToInt(r, 3); err == nil {
				this.date[ch] = v
			} else {
				return 0, fmt.Errorf("Parse error ch=%s, err=%s ", string(ch), err)
			}
		default:
			r.ReadRune()
		}
	}

	// format에 없는 경우, 현재 시간 기준
	if _, ok := this.date[DATEFORMAT_YEAR]; !ok {
		this.date[DATEFORMAT_YEAR] = now.Year()
	}
	if _, ok := this.date[DATEFORMAT_MONTH]; !ok {
		this.date[DATEFORMAT_MONTH] = int(now.Month())
	}
	if _, ok := this.date[DATEFORMAT_DAY]; !ok {
		this.date[DATEFORMAT_DAY] = now.Day()
	}
	if _, ok := this.date[DATEFORMAT_HOUR]; !ok {
		this.date[DATEFORMAT_HOUR] = now.Hour()
	}
	if _, ok := this.date[DATEFORMAT_MINUTE]; !ok {
		this.date[DATEFORMAT_MINUTE] = now.Minute()
	}
	if _, ok := this.date[DATEFORMAT_SECOND]; !ok {
		this.date[DATEFORMAT_SECOND] = now.Second()
	}
	if _, ok := this.date[DATEFORMAT_MILLISECOND]; !ok {
		this.date[DATEFORMAT_MILLISECOND] = int(now.UnixNano() / 1000000 % 1000)
	}

	d := time.Date(this.date[DATEFORMAT_YEAR], time.Month(this.date[DATEFORMAT_MONTH]), this.date[DATEFORMAT_DAY], this.date[DATEFORMAT_HOUR],
		this.date[DATEFORMAT_MINUTE], this.date[DATEFORMAT_SECOND], this.date[DATEFORMAT_MILLISECOND]*1000000, time.Now().Location())
	tm := d.UnixNano() / 1000000
	return tm, nil
}

func (this *DateFormat) ToInt(rd *bytes.Reader, size int) (int, error) {
	buf := make([]byte, size)
	if n, err := rd.Read(buf); err != nil {
		if err == io.EOF {
			return 0, nil
		}
		return 0, err
	} else {
		if n != size {
			return 0, fmt.Errorf("Error read size, want=%d, have=%d", size, n)
		}
	}

	if v, err := strconv.Atoi(string(buf)); err != nil {
		return 0, err
	} else {
		return v, nil
	}
}

func padding(n int, ch string) string {
	buf := bytes.Buffer{}
	for i := 0; i < n; i++ {
		buf.WriteString(ch)
	}
	return buf.String()
}

func LPadInt(v, size int) string {
	var ret string
	ret = fmt.Sprintf("%d", v)
	if len(ret) > size {
		return ret
	}
	return padding(size-len(ret), "0") + ret
}
