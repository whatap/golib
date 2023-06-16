package topology

import (
	"bytes"
	"fmt"

	//"strconv"
	"strings"

	"github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/castutil"
	"github.com/whatap/golib/util/hmap"
	//"gitlab.whatap.io/go/agent/util/logutil"
)

/**
 <pre>
attr
  {type=java,was=JAVA}
listen
	192.168.219.162:3306
	192.168.219.162:7283
	192.168.219.162:8888
	192.168.219.162:21300
outter
	59.30.186.24:27985
	91.108.56.199:80
	91.108.56.199:443
	54.243.172.216:443
	104.85.170.42:443
	107.23.165.206:443
	113.29.141.66:5223
	17.188.166.16:5223
	64.233.188.188:443
	108.177.97.125:5222
</pre>
*/

type NODE struct {
	Attr   *value.MapValue
	listen *hmap.LinkedSet
	outter *hmap.LinkedSet
}

func NewNODE() *NODE {
	p := new(NODE)
	p.Attr = value.NewMapValue()
	p.listen = hmap.NewLinkedSet()
	p.outter = hmap.NewLinkedSet()
	return p

}

func (this *NODE) ToString() string {
	return fmt.Sprintf("%s \nlisten\n %s \noutter\n %s", this.Attr, this.toString(this.listen, "\t"), this.toString(this.outter, "\t"))
}

func (this *NODE) toString(t *hmap.LinkedSet, space string) string {
	var buf bytes.Buffer

	it := t.Keys()
	for it.HasMoreElements() {
		buf.WriteString(space)

		tmp := it.NextElement()
		if tmp != nil {
			key := tmp.(*LINK)

			buf.WriteString(space)
			buf.WriteString(key.ToString())
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func (this *NODE) IsAttachable(k *LINK) bool {
	en := this.listen.Keys()
	for en.HasMoreElements() {
		tmp := en.NextElement()
		if tmp != nil {
			n := tmp.(*LINK)
			if n != nil && n.Include(k) {
				return true
			}
		}
	}
	return false
}

func (this *NODE) ToBytes() []byte {
	out := io.NewDataOutputX()
	out.WriteByte(0) // ver 0
	// out.writeValue(this.Attr)
	out.WriteByte(this.Attr.GetValueType())
	this.Attr.Write(out)
	this.toLinkBytes(this.listen, out)
	this.toLinkBytes(this.outter, out)

	return out.ToByteArray()

}

func (this *NODE) ToObject(b []byte) *NODE {
	in := io.NewDataInputX(b)
	//ver := in.ReadByte() // ver는 일단 더미
	in.ReadByte() // ver는 일단 더미

	mv := value.NewMapValue()
	mv.Read(in)
	this.Attr = mv
	this.listen = this.toLinkObject(in)
	this.outter = this.toLinkObject(in)

	return this

}

func (this *NODE) toLinkObject(in *io.DataInputX) *hmap.LinkedSet {
	data := hmap.NewLinkedSet()
	sz := int(in.ReadDecimal())
	for i := 0; i < sz; i++ {
		data.Put(NewLINK().ToObject(in))
	}
	return data
}

func (this *NODE) toLinkBytes(data *hmap.LinkedSet, out *io.DataOutputX) {
	out.WriteDecimal(int64(data.Size()))
	if data.Size() == 0 {
		return
	}
	en := data.Keys()
	for en.HasMoreElements() {
		tmp := en.NextElement()
		if tmp != nil {
			k := tmp.(*LINK)
			k.ToBytes(out)
		}
	}
}

func (this *NODE) AddListen(localIpSet *hmap.StringSet, listenAddr string) {
	ipo := this.getIPPORT(listenAddr)
	if ipo == nil || ipo.IsLocal127() {
		return
	}

	if ipo.IP == "*" || ipo.IP == "0.0.0.0" || ipo.IP == "::" {
		en := localIpSet.Keys()
		for en.HasMoreElements() {
			localIp := en.NextString()

			k := CreateLINK(localIp, int(castutil.CInt(ipo.Port)))
			if k != nil {
				this.listen.Put(k)
			}
		}
	} else {
		k := CreateLINK(ipo.IP, int(castutil.CInt(ipo.Port)))
		if k != nil {
			this.listen.Put(k)
		}
	}
}

func (this *NODE) getIPPORT(listenAddr string) (rt *IPO) {
	defer func() {
		if r := recover(); r != nil {
			// logutil.Println("WALT001", "Recover ", r)
			rt = nil
		}
	}()

	ipo := NewIPO()
	x := strings.LastIndex(listenAddr, ":")
	if x < 0 {
		x = strings.LastIndex(listenAddr, ".")
	}
	ipo.IP = listenAddr[0:x]
	ipo.Port = listenAddr[x+1:]

	return ipo
}

func (this *NODE) AddOutter(local, remote string) {
	localIPO := this.getIPPORT(local)
	if localIPO == nil || localIPO.IsIPv6() {
		return
	}

	if this.hasListen(localIPO.IP, localIPO.Port) {
		return
	}

	remoteIPO := this.getIPPORT(remote)
	if remoteIPO == nil || remoteIPO.IsIPv6() || remoteIPO.IsLocal127() {
		return
	}

	k := CreateLINK(remoteIPO.IP, int(castutil.CInt(remoteIPO.Port)))
	if k != nil {
		this.outter.Put(k)
	}
}

func (this *NODE) hasListen(ip, port string) bool {
	k := CreateLINK(ip, int(castutil.CInt(port)))
	if k != nil {
		return false
	} else {
		return this.listen.Contains(k)
	}
}

/*


	public JSONObject toJSON() {
		JSONObject o = new JSONObject();
		o.put("attr", new JSONObject(this.Attr.toJSONString()));
		o.put("listen", toJSON(this.listen));
		o.put("outter", toJSON(this.outter));
		return o;
	}

	private JSONArray toJSON(LinkedSet<LINK> data) {
		JSONArray out = new JSONArray();
		Enumeration<LINK> en = data.elements();
		while (en.hasMoreElements()) {
			LINK k = en.nextElement();
			out.put(k.toString());
		}
		return out;
	}

}
**/

type IPO struct {
	IP   string
	Port string
}

func NewIPO() *IPO {
	p := new(IPO)
	return p
}

func (this *IPO) IsIPv6() bool {
	return this.IP != "" && strings.Index(this.IP, ":") >= 0
}

func (this *IPO) IsLocal127() bool {
	return "127.0.0.1" == this.IP
}
