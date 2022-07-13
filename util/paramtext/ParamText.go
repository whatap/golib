package paramtext

import (
	"bytes"
	//"fmt"
	"strings"
)

type REF struct {
	key string
}
type ParamText struct {
	startBrace string
	endBrace   string
	tokenList  []interface{}
	keyList    []string
	origin     string
}

func (this *ParamText) ToStringMap(p map[string]string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(this.tokenList); i++ {
		t := this.tokenList[i]
		switch t.(type) {
		case string:
			buffer.WriteString(t.(string))
		default:
			r := t.(*REF)
			if p == nil {
				buffer.WriteString(this.startBrace)
				buffer.WriteString(r.key)
				buffer.WriteString(this.endBrace)
			} else {
				v, ok := p[r.key]
				if ok {
					buffer.WriteString(v)
				} else {
					buffer.WriteString(this.startBrace)
					buffer.WriteString(r.key)
					buffer.WriteString(this.endBrace)
				}
			}
		}
	}
	return buffer.String()
}

func (this *ParamText) GetKeys() []string {
	return this.keyList
}

func (this *ParamText) GetOriginal() string {
	return this.origin
}

func (this *ParamText) ToStringStr(p string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(this.tokenList); i++ {
		t := this.tokenList[i]
		switch t.(type) {
		case string:
			buffer.WriteString(t.(string))
		default:
			buffer.WriteString(p)
		}
	}
	return buffer.String()
}

func NewParamText(plainText string) *ParamText {
	return NewParamTextBrace(plainText, "${", "}")
}
func NewParamTextBrace(plainText string, startBrace string, endBrace string) *ParamText {
	p := new(ParamText)
	p.startBrace = startBrace
	p.endBrace = endBrace
	p.keyList = []string{}
	p.origin = plainText

	for len(plainText) > 0 {
		pos := strings.Index(plainText, startBrace)
		if pos < 0 {
			p.tokenList = append(p.tokenList, plainText)
			return p
		} else if pos > 0 {
			p.tokenList = append(p.tokenList, plainText[0:pos])
			plainText = plainText[pos:]
		} else {
			pos += len(startBrace)
			org := plainText
			plainText = plainText[pos:]
			nextPos := strings.Index(plainText, endBrace)
			if nextPos < 0 {
				p.tokenList = append(p.tokenList, org)
				return p
			}
			argName := strings.TrimSpace(plainText[0:nextPos])
			p.keyList = append(p.keyList, argName)
			p.tokenList = append(p.tokenList, &REF{argName})
			plainText = plainText[nextPos+len(endBrace):]
		}
	}

	return p
}
