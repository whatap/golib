package lang

import (
	"fmt"

	"github.com/whatap/golib/util/hmap"
)

type PKOID struct {
	PCode int64
	OKind int32
	Oid   int32
}

func NewPKOID(pcode int64, okind int32, oid int32) *PKOID {
	p := new(PKOID)
	p.PCode = pcode
	p.OKind = okind
	p.Oid = oid
	return p
}

func (this *PKOID) GetPCode() int64 {
	return this.PCode
}

func (this *PKOID) GetOid() int32 {
	return this.Oid
}

func (this *PKOID) Hash() uint {
	prime := 31
	result := 1
	result = prime*result + int(this.Oid)
	result = prime*result + int(this.OKind)
	result = prime*result + int(this.PCode^int64(uint64(this.PCode)>>32))
	return uint(result)
}

func (this *PKOID) Equals(obj hmap.LinkedKey) bool {
	if obj == nil {
		return false
	}
	if this == obj.(*PKOID) {
		return true
	}

	//		if (getClass() != obj.getClass())
	//			return false;
	other := obj.(*PKOID)
	if this.Oid != other.Oid {
		return false
	}
	if this.OKind != other.OKind {
		return false
	}
	if this.PCode != other.PCode {
		return false
	}
	return true
}

func (this *PKOID) ToString() string {
	return fmt.Sprintln("[", this.PCode, ",", this.OKind, ",", this.Oid, "]")
}

func (this *PKOID) CompareTo(n *PKOID) int {
	o := n
	v1 := this.PCode - o.PCode
	if v1 != 0 {
		if v1 > 0 {
			return 1
		} else {
			return -1
		}
	}
	v1 = int64(this.OKind - o.OKind)
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
