package lang

import (
	"fmt"

	"github.com/whatap/golib/util/hmap"
)

type POID struct {
	PCode int64
	Oid   int32
}

func NewPOID(pcode int64, oid int32) *POID {
	p := new(POID)
	p.PCode = pcode
	p.Oid = oid
	return p
}

func (this *POID) GetPCode() int64 {
	return this.PCode
}

func (this *POID) GetOid() int32 {
	return this.Oid
}

func (this *POID) Hash() uint {
	prime := 31
	result := 1
	result = prime*result + int(this.Oid)
	result = prime*result + int(this.PCode^int64(uint64(this.PCode)>>32))
	return uint(result)
}

func (this *POID) Equals(obj hmap.LinkedKey) bool {
	if obj == nil {
		return false
	}
	if this == obj.(*POID) {
		return true
	}

	//		if (getClass() != obj.getClass())
	//			return false;
	other := obj.(*POID)
	if this.Oid != other.Oid {
		return false
	}
	if this.PCode != other.PCode {
		return false
	}
	return true
}

func (this *POID) ToString() string {
	return fmt.Sprintln("[", this.PCode, ",", this.Oid, "]")
}

func (this *POID) CompareTo(n *POID) int {
	o := n
	v1 := this.PCode - o.PCode
	if v1 != 0 {
		if v1 > 0 {
			return 1
		} else {
			return -1
		}
	}
	v1 = int64(this.Oid - o.Oid)
	if v1 != 0 {
		if v1 > 0 {
			return 1
		} else {
			return -1
		}
	}
	return 0
}
