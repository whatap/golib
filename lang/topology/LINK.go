package topology

import (
	"fmt"
	"net"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/util/compare"
	"github.com/whatap/golib/util/hash"
	"github.com/whatap/golib/util/hmap"
)

//import java.net.InetAddress;
//import whatap.util.CompareUtil;

type LINK struct {
	IP   net.IP
	Port int
}

func NewLINK() *LINK {
	p := new(LINK)
	return p
}

func (this *LINK) ToString() string {
	if this.IP == nil {
		return fmt.Sprintf("0.0.0.0:%d", this.Port)
	}

	arr, err := net.LookupIP(this.IP.String())
	if err != nil {
		// logutil.Println("WALT101", "LINK ip=", this.IP.String(), ",lookup Error=", err)
		return fmt.Sprintf("0.0.0.0:%d", this.Port)
	}
	for _, val := range arr {
		ip := val.To4()
		return fmt.Sprintf("%s:%d", ip.String(), this.Port)
	}

	return fmt.Sprintf("0.0.0.0:%d", this.Port)
}

func CreateLINK(ipStr string, port int) *LINK {
	k := NewLINK()
	k.IP = make([]byte, 4)

	arr, err := net.LookupIP(ipStr)
	if err != nil {
		// logutil.Println("WALT102", "LINK ip=", ipStr, ",lookup Error=", err)
		return k
	}
	// 처음 한개만 설정. 추후 필요하면 여러개 설정.
	for _, val := range arr {
		ip := val.To4()
		copy(k.IP, ip)
		break
	}

	k.Port = port
	return k
}

func (this *LINK) Include(k *LINK) bool {
	if compare.EqualBytes(this.IP, k.IP) == false {
		return false
	}

	if this.Port == 0 {
		return true
	}
	return this.Port == k.Port
}

// LinkedKey Interface for put LinkedSet
func (this *LINK) Hash() uint {
	return uint(this.HashCode())
}

// LinkedKey Interface for put LinkedSet
//
//	func (this *LINK) Equals(h LinkedKey) bool {
//		return this.
//	}
func (this *LINK) Equals(h hmap.LinkedKey) bool {
	var k *LINK
	if h != nil {
		k = h.(*LINK)
	}
	if compare.EqualBytes(this.IP, k.IP) == false {
		return false
	}
	return this.Port == k.Port
}

func (this *LINK) HashCode() int32 {
	return (hash.Hash(this.IP) | int32(this.Port))
}

func (this *LINK) ToBytes(out *io.DataOutputX) *LINK {
	out.WriteBlob(this.IP)
	out.WriteInt(int32(this.Port))
	return this
}

func (this *LINK) ToObject(in *io.DataInputX) *LINK {
	this.IP = in.ReadBlob()
	this.Port = int(in.ReadInt())
	return this
}
