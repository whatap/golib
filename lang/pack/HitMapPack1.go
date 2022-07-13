package pack

import (
	"github.com/whatap/golib/io"
)

const (
	HITMAP_LENGTH = 120
	// 2.5초
	HITMAP_VERTICAL_INDEX = 20
	// 5초
	HITMAP_HORIZONTAL_INDEX = 60
)

type HitMapPack1 struct {
	AbstractPack
	Hit   []int32
	Error []int32
}

func NewHitMapPack1() *HitMapPack1 {
	p := new(HitMapPack1)
	p.Hit = make([]int32, HITMAP_LENGTH)
	p.Error = make([]int32, HITMAP_LENGTH)
	return p
}

func (this *HitMapPack1) GetPackType() int16 {
	return PACK_HITMAP_1
}

func (this *HitMapPack1) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteByte(1)
	for i := 0; i < HITMAP_LENGTH; i++ {
		dout.WriteShort(int16(this.Hit[i]))
		dout.WriteShort(int16(this.Error[i]))
	}
}
func (this *HitMapPack1) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	ver := din.ReadByte()
	if ver == 1 {
		this.Hit = make([]int32, HITMAP_LENGTH)
		this.Error = make([]int32, HITMAP_LENGTH)
		for i := 0; i < HITMAP_LENGTH; i++ {
			this.Hit[i] = int32(din.ReadShort()) & 0xffff
			this.Error[i] = int32(din.ReadShort()) & 0xffff
		}
	}
}

func (this *HitMapPack1) Add(time int, isError bool) {
	idx := index(time)

	this.Hit[idx]++
	if isError {
		this.Error[idx]++
	}
}

// java에는 없지만 index, time 을 외부에서 쓰기 위해 추가
func (this *HitMapPack1) HitMapIndex(time int) int {
	return index(time)
}

func (this *HitMapPack1) HitMapTime(index int) int {
	return time(index)
}

func index(time int) int {
	x := time / 10000
	switch x {
	case 0:
		if time < 5000 {
			return time / 125
		} else {
			return 40 + (time-5000)/250
		}
	case 1:
		return 60 + (time-10000)/500
	case 2, 3:
		return 80 + (time-20000)/1000
	case 4, 5, 6, 7:
		return 100 + (time-40000)/2000
	default:
		return 119
	}
}

func time(index int) int {
	if index < 40 {
		return index * 125
	}
	_time_ := 5000
	if index < 60 {
		return _time_ + (index-40)*250
	}
	_time_ = 10000
	if index < 80 {
		return _time_ + (index-60)*500
	}
	_time_ = 20000
	if index < 100 {
		return _time_ + (index-80)*1000
	}
	_time_ = 40000
	if index < 120 {
		return _time_ + (index-100)*2000
	} else {
		return 80000
	}
}
