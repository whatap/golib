package paramtext

import (
	"bytes"
	//"fmt"
	"strings"
)

/**
문자열 내부 key=value 조합의 value 값을 치환

separate : key=value(구분자:seperate)key=value
kvPipe : key(pipe)value
ToStringMap : k,v Map을 받아서 해당 key의 value를 모두 변경
ToStringStr : k,v string으로 하나의 키만 변환.
*/

type ParamKV struct {
	separate  string
	kvPipe    string
	tokenList []string
	kvMap     map[string]string
	origin    string
}

func NewParamKV(plainText string) *ParamKV {
	return NewParamKVSeperate(plainText, " ", "=")
}

func NewParamKVSeperate(plainText string, separate string, kvPipe string) *ParamKV {
	p := new(ParamKV)
	p.separate = separate
	p.kvPipe = kvPipe
	p.tokenList = strings.Split(plainText, separate)
	p.kvMap = make(map[string]string)
	p.origin = plainText

	if p.tokenList != nil {
		for _, kv := range p.tokenList {
			k, v := ToPair(kv, p.kvPipe)
			if k != "" {
				p.kvMap[k] = v
			}
		}
	}
	return p
}

func (this *ParamKV) ToStringMap(p map[string]string) string {
	// key 의 value 치환
	for key, val := range p {
		if this.ExistsKey(key) {
			this.kvMap[key] = val
		}
	}

	return this.ToString()
}

func (this *ParamKV) ToStringStr(k, v string) string {

	if this.ExistsKey(k) {
		this.kvMap[k] = v
	}

	return this.ToString()
}

func (this *ParamKV) ToString() string {
	var buffer bytes.Buffer

	if this.tokenList != nil {
		kvLen := len(this.tokenList)
		for i, kv := range this.tokenList {
			k, _ := ToPair(kv, this.kvPipe)
			if k != "" {
				buffer.WriteString(k)
				buffer.WriteString(this.kvPipe)
				buffer.WriteString(this.kvMap[k])

			} else {
				buffer.WriteString(kv)
			}

			if i < kvLen-1 {
				buffer.WriteString(this.separate)
			}
		}
	}

	return buffer.String()
}

func (this *ParamKV) ExistsKey(key string) bool {
	_, exists := this.kvMap[key]
	return exists
}

func (this *ParamKV) GetKeys() []string {
	keys := make([]string, len(this.kvMap))
	i := 0
	for k := range this.kvMap {
		keys[i] = k
		i++
	}
	return keys
}

func (this *ParamKV) GetValue(key string) string {
	if this.ExistsKey(key) {
		return this.kvMap[key]
	}
	return ""
}

func (this *ParamKV) GetOriginal() string {
	return this.origin
}

// get k, v
func ToPair(s string, sep string) (k, v string) {
	pos := 0
	pos = strings.Index(strings.ToLower(s), strings.ToLower(sep))
	if pos != -1 {
		k = s[0:pos]
		v = s[pos+len(sep):]
	} else {
		k = ""
		v = ""
	}

	return strings.TrimSpace(k), strings.TrimSpace(v)
}
