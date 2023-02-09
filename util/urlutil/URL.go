package urlutil

import (
	//	"log"
	_ "fmt"
	"net/url"
	"strconv"
	"strings"
)

// java URL 처럼
type URL struct {
	Url      string
	Protocol string
	Host     string
	RawPath  string
	Path     string
	RawPort  string
	Port     int
	RawQuery string
	Query    string
	File     string
}

func NewURL(url string) *URL {
	p := new(URL)
	p.Url = url

	p.process()

	return p
}
func (this *URL) process1() {
	if u, err := url.Parse(this.Url); err == nil {
		this.Protocol = u.Scheme
		this.Host = u.Host
		this.RawPort = u.Port()
		if this.RawPort == "" {
			if strings.ToLower(this.Protocol) == "https" {
				this.Port = 443
			} else {
				this.Port = 80
			}
		} else {
			this.Port, _ = strconv.Atoi(this.RawPort)
		}
		pos := strings.LastIndex(u.Path, "/")
		if pos > -1 {
			this.Path = string(u.Path[pos:])
		}

		this.Path = u.Path
		if this.Path == "" {
			this.Path = "/"
		}
		this.Query = u.RawQuery

		// Path, Query URLDecode추가 (net/url)
		this.RawPath = this.Path
		this.Path, err = url.PathUnescape(this.RawPath)
		if err != nil {
			this.Path = this.RawPath
		}
		this.RawQuery = this.Query
		this.Query, err = url.QueryUnescape(this.RawQuery)
		if err != nil {
			this.Query = this.RawQuery
		}

	} else {

	}
}
func (this *URL) process() {
	var tmp string
	var pos int
	var err error

	// query 먼저 분리
	var urlStr string
	if pos := strings.Index(this.Url, "?"); pos > -1 {
		this.Query = string(this.Url[pos+1:])
		urlStr = string(this.Url[0:pos])
	} else {
		urlStr = this.Url
	}
	if pos := strings.LastIndex(urlStr, "/"); pos > -1 {
		this.File = string(urlStr[pos:])
	}

	// Protocol
	tmp = urlStr
	pos = strings.Index(tmp, "://")
	if pos > -1 {
		this.Protocol = tmp[0:pos]
		tmp = tmp[pos+3:]
	}

	// Host, Port
	pos = strings.Index(tmp, "/")
	if pos > -1 {
		this.Host = tmp[0:pos]
		tmp = tmp[pos:]

		// Port 분리
		pos = strings.Index(this.Host, ":")
		if pos > -1 {
			this.RawPort = this.Host[pos+1:]
			this.Port, _ = strconv.Atoi(this.Host[pos+1:])
			this.Host = this.Host[0:pos]
		} else {
			if this.Protocol == "https" {
				this.Port = 443
			} else {
				this.Port = 80
			}
		}

	} else {
		this.Host = tmp
		if this.Protocol == "https" {
			this.Port = 443
		} else {
			this.Port = 80
		}
	}

	// Path, File, Query
	pos = strings.Index(tmp, "?")
	if pos > -1 {
		this.Path = tmp[0:pos]
		//this.File = tmp
		//this.Query = tmp[pos+1:]
	} else {
		this.Path = tmp
		//this.File = tmp
		//this.Query = ""
	}

	// Path, Query URLDecode추가 (net/url)
	this.RawPath = this.Path
	this.Path, err = url.PathUnescape(this.RawPath)
	if err != nil {
		this.Path = this.RawPath
	}
	this.RawQuery = this.Query
	this.Query, err = url.QueryUnescape(this.RawQuery)
	if err != nil {
		this.Query = this.RawQuery
	}

}

// URL Decode 된  url 정보를 다시 조합해서 출력
func (this *URL) String() string {
	rt := ""
	if this.Protocol != "" {
		rt = rt + this.Protocol + "://"
	}
	rt = rt + this.Host

	if this.RawPort != "" {
		rt = rt + ":" + this.RawPort
	}

	rt = rt + this.Path

	if this.Query != "" {
		rt = rt + "?" + this.Query
	}
	return rt
}

func (this *URL) Domain() string {
	rt := ""
	if this.Protocol != "" {
		rt = rt + this.Protocol + "://"
	}
	rt = rt + this.Host
	return rt
}

func (this *URL) DomainPath() string {
	rt := ""
	if this.Protocol != "" {
		rt = rt + this.Protocol + "://"
	}
	rt = rt + this.Host

	if this.RawPort != "" {
		rt = rt + ":" + this.RawPort
	}

	rt = rt + this.Path
	return rt
}
