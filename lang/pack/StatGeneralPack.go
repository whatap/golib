package pack

import (
	"fmt"
	"sync"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/hmap"
	"github.com/whatap/golib/util/list"
)

type StatGeneralPack struct {
	AbstractPack
	Id            string
	dataBytes     []byte
	dataBytesSize int
	data          *hmap.StringKeyLinkedMap
	lock          sync.Mutex
}

func NewStatGeneralPack() *StatGeneralPack {
	p := new(StatGeneralPack)
	p.data = hmap.NewStringKeyLinkedMap()
	return p
}

func (this *StatGeneralPack) GetPackType() int16 {
	return PACK_STAT_GENERAL
}

func (this *StatGeneralPack) ToString() string {
	return fmt.Sprintln("StatGeneralPack ", this.AbstractPack.ToString(), ",data=", this.data.Size(), ",length=", this.dataBytesSize, ",bytes=", len(this.dataBytes))
}

func (this *StatGeneralPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)

	dout.WriteText(this.Id)
	//this.data.size() > 0 && ArrayUtil.isEmpty(this.dataBytes
	if this.data.Size() > 0 && (this.dataBytes == nil || len(this.dataBytes) == 0) {
		this.dataBytes = this.writeTable(this.data)
		this.dataBytesSize = len(this.dataBytes)
	}

	dout.WriteInt3(int32(this.dataBytesSize))
	dout.WriteBytes(this.dataBytes)
}

func (this *StatGeneralPack) writeTable(data *hmap.StringKeyLinkedMap) []byte {
	dd := io.NewDataOutputX()

	dd.WriteShort(int16(data.Size()))

	en := data.Entries()
	for en.HasMoreElements() {
		ent := en.NextElement().(*hmap.StringKeyLinkedEntry)
		dd.WriteText(ent.GetKey())
		dd.WriteByte(ent.GetValue().(list.AnyList).GetType())
		ent.GetValue().(list.AnyList).Write(dd)

	}
	return dd.ToByteArray()
}

func (this *StatGeneralPack) readTable(bytes []byte, data *hmap.StringKeyLinkedMap) int {
	in := io.NewDataInputX(bytes)
	cnt := int(in.ReadShort())
	for i := 0; i < cnt; i++ {
		key := in.ReadText()
		a := this.create(in.ReadByte())
		a.Read(in)
		data.Put(key, a)
	}
	return cnt
}

func (this *StatGeneralPack) create(t byte) list.AnyList {
	switch t {
	case list.ANYLIST_INT:
		return list.NewIntListDefault()
	case list.ANYLIST_LONG:
		return list.NewLongListDefault()
	case list.ANYLIST_FLOAT:
		return list.NewFloatListDefault()
	case list.ANYLIST_DOUBLE:
		return list.NewDoubleListDefault()
	default:
		return list.NewStringListDefault()
	}
}

func Sort(data *hmap.StringKeyLinkedMap, sortKey string, asc bool) *hmap.StringKeyLinkedMap {
	any := data.Get(sortKey).(list.AnyList)
	if any == nil {
		return data
	}
	ord := any.Sorting(asc)

	data2 := hmap.NewStringKeyLinkedMap()
	en := data.Keys()
	for en.HasMoreElements() {
		key := en.NextString()
		oldList := data.Get(key).(list.AnyList)
		newList := oldList.Filtering(ord)
		data2.Put(key, newList)
	}
	return data2

}

func SortAnyList(data *hmap.StringKeyLinkedMap, sortKey string, asc bool, sortKey2 string, asc2 bool) *hmap.StringKeyLinkedMap {
	f1 := data.Get(sortKey).(list.AnyList)
	if f1 == nil {
		return data
	}

	f2 := data.Get(sortKey2).(list.AnyList)
	if f2 == nil {
		return data
	}

	ord := f1.SortingAnyList(asc, f2, asc2)
	data2 := hmap.NewStringKeyLinkedMap()
	en := data.Keys()
	for en.HasMoreElements() {
		key := en.NextString()
		oldList := data.Get(key).(list.AnyList)
		newList := oldList.Filtering(ord)
		data2.Put(key, newList)
	}
	return data2
}

func (this *StatGeneralPack) Sort(sortKey string, asc bool) {
	this.data = Sort(this.data, sortKey, asc)
}

func (this *StatGeneralPack) SortAnyList(sortKey string, asc bool, sortKey2 string, asc2 bool) {
	this.data = SortAnyList(this.data, sortKey, asc, sortKey2, asc2)
}

func (this *StatGeneralPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)
	this.Id = din.ReadText()
	this.dataBytesSize = int(din.ReadInt3())
	this.dataBytes = din.ReadBytes(int32(this.dataBytesSize))
	//return this
}

func (this *StatGeneralPack) GetDataTable() *hmap.StringKeyLinkedMap {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.unpack()
	return this.data
}

func (this *StatGeneralPack) unpack() {
	if this.dataBytes != nil {
		this.readTable(this.dataBytes, this.data)
		this.dataBytes = nil
		this.dataBytesSize = 0
	}
}

func (this *StatGeneralPack) Put(key string, data list.AnyList) {
	this.data.Put(key, data)
}

func (this *StatGeneralPack) Get(key string) list.AnyList {
	this.unpack()
	return this.data.Get(key).(list.AnyList)
}

func (this *StatGeneralPack) IsEmpty() bool {
	return this.dataBytesSize == 0 && this.data.IsEmpty()
}

func (this *StatGeneralPack) Iterate(h func(a []string, b []list.AnyList, c int)) {
	this.unpack()
	if this.data.Size() == 0 {
		return
	}
	title := make([]string, this.data.Size())
	values := make([]list.AnyList, this.data.Size())
	en := this.data.Entries()
	for i := 0; en.HasMoreElements(); i++ {
		ent := en.NextElement().(*hmap.StringKeyLinkedEntry)
		title[i] = ent.GetKey()
		values[i] = ent.GetValue().(list.AnyList)
	}

	n := values[0].Size()
	for i := 0; i < n; i++ {
		h(title, values, i)
	}
}
