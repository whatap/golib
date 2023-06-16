package shellarg

import (
	"strings"

	"github.com/whatap/golib/util/castutil"
	"github.com/whatap/golib/util/hmap"
)

type ShellArg struct {
	parameter  *hmap.StringKeyLinkedMap
	parameter2 *hmap.StringKeyLinkedMap
	Tags       *hmap.StringKeyLinkedMap
}

func HasNextValue(i int, args []string) bool {
	return i+1 < len(args) && strings.HasPrefix(args[i+1], "-") == false
}

func NewShellArg(args []string) *ShellArg {
	this := new(ShellArg)
	this.Tags = hmap.NewStringKeyLinkedMap()
	this.parameter = hmap.NewStringKeyLinkedMap()
	this.parameter2 = hmap.NewStringKeyLinkedMap()

	max := len(args)
	i := 0
	for i < max {
		if strings.HasPrefix(args[i], "-tag.") {
			this.Tags.Put(args[i], args[i][5:])
		}
		if HasNextValue(i, args) {
			key := args[i]
			this.parameter.Put(key, args[i+1])
			i++
			if HasNextValue(i, args) {
				this.parameter2.Put(key, args[i+1])
				i++
			}
		} else {
			this.parameter.Put(args[i], "")
		}
		i++
	}
	return this
}

func (this *ShellArg) HasKey(key string) bool {
	return this.parameter.ContainsKey(key)
}

func (this *ShellArg) Keys() hmap.StringEnumer {
	return this.parameter.Keys()
}

func (this *ShellArg) Get(key string, defaultValue string) string {
	s := this.parameter.Get(key)
	if s == nil {
		return defaultValue
	} else {
		return s.(string)
	}
}

func (this *ShellArg) GetInt(key string, defaultValue int32) int32 {
	s := this.parameter.Get(key)
	if s == nil {
		return defaultValue
	} else {
		return castutil.CInt(s)
	}
}
func (this *ShellArg) GetLong(key string, defaultValue int64) int64 {
	s := this.parameter.Get(key)
	if s == nil {
		return defaultValue
	} else {
		return castutil.CLong(s)
	}
}
func (this *ShellArg) GetBoolean(key string, defaultValue bool) bool {
	s := this.parameter.Get(key)
	if s == nil {
		return defaultValue
	} else {
		return castutil.CBool(s)
	}
}

func (this *ShellArg) Get2(key string) string {
	return this.parameter2.Get(key).(string)
}

func (this *ShellArg) Put(key string, value string) {
	this.parameter.Put(key, value)
}
