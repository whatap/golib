package stringutil

import (
	//	"log"
	"bytes"
	"fmt"
	"strings"
)

type StringBuffer struct {
	buffer bytes.Buffer
	indent int
}

func NewStringBuffer() *StringBuffer {
	ret := new(StringBuffer)
	ret.indent = 0
	return ret
}

/// append a string to the tail of this buffer
func (sb *StringBuffer) Append(str string) *StringBuffer {
	sb.buffer.WriteString(strings.Repeat("\t", sb.indent) + str)
	return sb
}

func (sb *StringBuffer) AppendFormat(str string, args ...interface{}) *StringBuffer {
	return sb.Append(fmt.Sprintf(str, args...))
}

func (sb *StringBuffer) AppendLine(str string) *StringBuffer {
	return sb.Append(str + "\n")
}

func (sb *StringBuffer) AppendLineIndent(str string) *StringBuffer {
	sb.AppendLine(str)
	sb.indent++
	return sb
}

func (sb *StringBuffer) AppendLineClose(str string) *StringBuffer {
	sb.indent--
	return sb.AppendLine(str)
}

func (sb *StringBuffer) AppendClose(str string) *StringBuffer {
	sb.indent--
	return sb.Append(str)
}

/// append a line as comment.
func (sb *StringBuffer) AppendComment(str string) *StringBuffer {
	return sb.AppendLine("/// " + str)
}

func (sb *StringBuffer) ToString() string {
	return sb.buffer.String()
}

/// clear .
func (sb *StringBuffer) Clear() {
	sb.buffer.Reset()
	sb.indent = 0
}
