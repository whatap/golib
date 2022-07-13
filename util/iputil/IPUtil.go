package iputil

import (
	"bytes"
	"encoding/hex"
	"net"
	"strconv"
	"strings"

	"github.com/whatap/golib/io"
)

func ToStringFrInt(ip int32) string {
	return ToString(io.ToBytesInt(ip))
}

func ToStringInt(ip int32) string {
	return ToString(io.ToBytesInt(ip))
}

func ToString(ip []byte) string {
	if ip == nil || len(ip) == 0 {
		return "0.0.0.0"
	}
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(uint(ip[0]))))
	buffer.WriteString(".")
	buffer.WriteString(strconv.Itoa(int(uint(ip[1]))))
	buffer.WriteString(".")
	buffer.WriteString(strconv.Itoa(int(uint(ip[2]))))
	buffer.WriteString(".")
	buffer.WriteString(strconv.Itoa(int(uint(ip[3]))))
	return buffer.String()
}

func ToBytes(ip string) []byte {
	if ip == "" {
		return []byte{0, 0, 0, 0}
	}
	result := []byte{0, 0, 0, 0}
	s := strings.Split(ip, ".")
	if len(s) != 4 {
		return []byte{0, 0, 0, 0}
	}
	for i := 0; i < 4; i++ {
		if val, err := strconv.Atoi(s[i]); err == nil {
			result[i] = (byte)(val & 0xff)
		}
	}
	return result
}

func ToBytesFrInt(ip int32) []byte {
	return io.ToBytesInt(ip)
}
func ToInt(ip []byte) int32 {
	return io.ToInt(ip, 0)
}

func IsOK(ip []byte) bool {
	return ip != nil && len(ip) == 4
}

func IsNotLocal(ip []byte) bool {
	return IsOK(ip) && uint(ip[0]) != 127
}

func ParseHexString(ipport string) ([]byte, error) {
	words := strings.Split(ipport, ":")
	parsedbytes, err := hex.DecodeString(words[0])
	if err != nil {
		return nil, err
	}
	parsedLength := len(parsedbytes)
	ipbytes := make([]byte, 6)
	ipbytes[3] = parsedbytes[parsedLength-4]
	ipbytes[2] = parsedbytes[parsedLength-3]
	ipbytes[1] = parsedbytes[parsedLength-2]
	ipbytes[0] = parsedbytes[parsedLength-1]

	portbytes, err := hex.DecodeString(words[1])
	ipbytes[4] = portbytes[0]
	ipbytes[5] = portbytes[1]

	return ipbytes, nil
}

func GetIPsToString() string {
	rt := make([]string, 0)

	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, i := range ifaces {
		if i.Flags&net.FlagUp != 1 {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				if !v.IP.IsLoopback() {
					rt = append(rt, v.IP.String())
				}
			case *net.IPNet:
				if !v.IP.IsLoopback() {
					rt = append(rt, v.IP.String())
				}
			}
		}
	}

	return strings.Join(rt, ",")
}

func LocalAddresses() (rt []net.IP) {
	rt = make([]net.IP, 0)

	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range ifaces {
		if i.Flags&net.FlagUp != 1 {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				//logutil.Printf("WAIPUtil005", "----- %v : %s (%s) %t %d [%v,%v]\n", i.Name, v, v.IP.DefaultMask(), v.IP.IsLoopback(), len(v.IP), v.IP.To4(), v.IP.To16)
				if !v.IP.IsLoopback() && v.IP.To4() != nil {
					//logutil.Println("WAIPUtil006", v.IP)
					rt = append(rt, v.IP)
				}

			case *net.IPNet:
				//logutil.Printf("WAIPUtil007", "------ %v : %s [%v/%v] %t %d [%v,%v]\n", i.Name, v, v.IP, v.Mask, v.IP.IsLoopback(), len(v.IP), v.IP.To4(), v.IP.To16())
				if !v.IP.IsLoopback() && v.IP.To4() != nil {
					//logutil.Println("WAIPUtil008", v.IP, v.IP[12], v.IP[13], v.IP[14], v.IP[15])
					rt = append(rt, v.IP)
				}
			}
		}
	}

	return rt
}

func lookup(ip net.IP) bool {
	_, err := net.LookupHost(ip.String())
	//arr, err := net.LookupAddr(it.String())
	if err != nil {
		return false
	}
	return true
}
