package conffile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/util/dateutil"
	"github.com/whatap/golib/util/hash"
	"github.com/whatap/golib/util/stringutil"
)

const ()

type FileParser interface {
	Read(filePath string) (map[string]string, error)
	Write(fielPath string, m *map[string]string) error
}

type FileConfig struct {
	m              map[string]string
	conf           *fileConfigConfig
	last_file_time int64
	last_check     int64
	Debug          bool
}

var confInstance *FileConfig = nil
var mutex = sync.Mutex{}

func GetConfig(opts ...FileConfigOption) config.Config {
	mutex.Lock()
	defer mutex.Unlock()
	if confInstance != nil {
		return confInstance
	}
	confInstance = newFileConfig(opts...)
	return confInstance
}
func newFileConfig(opts ...FileConfigOption) *FileConfig {
	c := new(FileConfig)
	c.m = make(map[string]string)
	c.last_file_time = -1
	c.last_check = 0

	c.conf = defaultFileConfigConfig()

	for _, opt := range opts {
		opt.apply(c.conf)
	}

	if c.conf.ctx == nil {
		c.conf.ctx, c.conf.cancel = context.WithCancel(context.Background())
	} else if c.conf.cancel == nil {
		c.conf.ctx, c.conf.cancel = context.WithCancel(c.conf.ctx)
	}
	c.reload()
	go c.run()

	return c
}
func (this *FileConfig) Destroy() {
	this.conf.cancel()
	confInstance = nil
}
func (this *FileConfig) GetWhatapHome() string {
	if this.conf.homePath != "" {
		return this.conf.homePath
	}
	home := os.Getenv("WHATAP_HOME")
	if home == "" {
		home = "."
	}
	return home
}
func (this *FileConfig) GetConfFile() string {
	home := this.GetWhatapHome()
	// config 파일이 WHATAP_HOME 과 다른 경로에 있을 경우 설정.
	confHome := os.Getenv("WHATAP_CONFIG_HOME")
	if confHome != "" {
		home = confHome
	}

	confName := os.Getenv("WHATAP_CONFIG")
	if confName == "" {
		confName = "whatap.conf"
	}

	return filepath.Join(home, confName)
}

func (this *FileConfig) run() {
	for {
		select {
		case <-this.conf.ctx.Done():
			this.conf.Log.Warn("FileConfig Destroy")
			return
		default:
			time.Sleep(3000 * time.Millisecond)
			this.reload()
		}
	}
}

func (this *FileConfig) ApplyDefault() {

	this.m["enabled"] = "true"
	this.m["net_udp_port"] = "6600"
	this.m["transaction_enabled"] = "true"
	this.m["profile_http_header_enabled"] = "false"
	this.m["profile_http_header_url_prefix"] = "/"
	this.m["profile_http_parameter_enabled"] = "false"
	this.m["profile_http_parameter_url_prefix"] = "/"

	this.m["profile_sql_param_enabled"] = "false"

	this.m["trace_user_enabled"] = "true"
	this.m["trace_user_using_ip"] = "false"
	this.m["trace_user_header_ticket"] = ""
	this.m["trace_user_set_cookie"] = "false"
	this.m["trace_user_cookie_limit"] = "2048"
	this.m["trace_user_cookie_keys"] = ""
	this.m["trace_http_client_ip_header_key_enabled"] = "true"
	this.m["trace_http_client_ip_header_key"] = "x-forwarded-for"

	this.m["mtrace_enabled"] = "true"
	this.m["mtrace_caller_key"] = "x-wtap-mst"
	this.m["mtrace_callee_key"] = "x-wtap-tx"
	this.m["mtrace_info_key"] = "x-wtap-inf"
	this.m["mtrace_poid_key"] = "x-wtap-po"
	this.m["mtrace_spec_key"] = "x-wtap-sp"
	this.m["mtrace_spec_key1"] = "x-wtap-sp1"
	this.m["mtrace_send_url_length"] = "80"
	this.m["mtrace_spec"] = "ver1.0"
	this.m["mtrace_rate"] = "10"

	this.m["tx_max_count"] = "8000"

	this.m["debug"] = "false"

	// this.ConfGo.ApplyDefault(m)
	// this.ConfGoGrpc.ApplyDefault(m)
	// this.ApplyConfig(m)
}

func (this *FileConfig) apply(newM map[string]string) {
	for k, v := range newM {
		if v != "" {
			this.m[k] = v
		}
	}
}

func (this *FileConfig) reload() {
	this.conf.Log.Println("reload")
	// 종료 되지 않도록  Recover
	defer func() {
		if r := recover(); r != nil {
			this.conf.Log.Error("WA211 Recover", r) //, string(debug.Stack()))
		}
	}()

	now := dateutil.Now()
	if now < this.last_check+3000 {
		return
	}
	this.last_check = now
	path := this.GetConfFile()

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		if this.last_file_time == -1 {
			this.conf.Log.Error("WA212", "fail to load config file")
			return
		} else if this.last_file_time == 0 {
			return
		}
		// 파일이 없는 경우 기본 설정으로 초기.
		this.last_file_time = 0
		this.m = make(map[string]string)
		this.ApplyDefault()
		this.conf.Log.Warn("WA213", " Default Config: ", this.GetConfFile())
		return
	}
	// is change?
	new_time := stat.ModTime().Unix()
	if this.last_file_time == new_time {
		return
	}

	this.last_file_time = new_time
	if m, err := this.conf.Parser.Read(path); err == nil {
		this.apply(m)
	} else {
		return
	}

	// Observer run
	if this.conf.configObserver != nil {
		this.conf.configObserver.Run(this)
	}

	this.conf.Log.Warn("WA214", "Reload Config: ", this.GetConfFile())
}

// interface ConfigObaserver
func (this *FileConfig) ApplyConfig(m map[string]string) {
	if m != nil {
		for k, v := range m {
			this.m[k] = v
		}
	}
}

func (this *FileConfig) GetKeys() []string {
	keys := make([]string, 0)
	for k, _ := range this.m {
		keys = append(keys, k)
	}
	return keys
}

func (this *FileConfig) GetValue(key string) string { return this.getValue(key) }
func (this *FileConfig) getValue(key string) string {
	if v, ok := this.m[key]; ok {
		return strings.TrimSpace(v)
	}
	return os.Getenv(key)
}
func (this *FileConfig) GetValueDef(key, def string) string { return this.getValueDef(key, def) }
func (this *FileConfig) getValueDef(key string, def string) string {
	v := this.getValue(key)

	if v == "" {
		return def
	}

	return v
}
func (this *FileConfig) GetBoolean(key string, def bool) bool {
	return this.getBoolean(key, def)
}
func (this *FileConfig) getBoolean(key string, def bool) bool {
	v := this.getValue(key)
	if v == "" {
		return def
	}
	value, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return value
}
func (this *FileConfig) GetInt(key string, def int) int32 {
	return this.getInt(key, def)
}
func (this *FileConfig) getInt(key string, def int) int32 {
	v := this.getValue(key)
	if v == "" {
		return int32(def)
	}
	value, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return int32(def)
	}
	return int32(value)
}

func (this *FileConfig) GetIntSet(key, defaultValue, deli string) []int32 {
	set := make([]int32, 0)
	vv := stringutil.Tokenizer(this.GetValueDef(key, defaultValue), deli)
	if vv != nil {
		for _, x := range vv {
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Continue
					}
				}()
				if xx, err := strconv.Atoi(strings.TrimSpace(x)); err != nil {
					set = append(set, int32(xx))
				}
			}()
		}
	}
	return set
}

func (this *FileConfig) GetStringHashSet(key, defaultValue, deli string) []int32 {
	set := make([]int32, 0)
	vv := stringutil.Tokenizer(this.GetValueDef(key, defaultValue), deli)
	if vv != nil {
		for _, x := range vv {
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Continue
					}
				}()
				xx := hash.HashStr(strings.TrimSpace(x))
				set = append(set, xx)
			}()
		}
	}
	return set
}

func (this *FileConfig) GetStringHashCodeSet(key, defaultValue, deli string) []int32 {
	set := make([]int32, 0)
	vv := stringutil.Tokenizer(this.GetValueDef(key, defaultValue), deli)
	if vv != nil {
		for _, x := range vv {
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Continue
					}
				}()
				xx := stringutil.HashCode(strings.TrimSpace(x))
				set = append(set, int32(xx))
			}()
		}
	}
	return set
}
func (this *FileConfig) GetLong(key string, def int64) int64 {
	return this.getLong(key, def)
}
func (this *FileConfig) getLong(key string, def int64) int64 {
	v := this.getValue(key)
	if v == "" {
		return def
	}
	value, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return value
}
func (this *FileConfig) GetStringArray(key, def, deli string) []string {
	return this.getStringArray(key, def, deli)
}
func (this *FileConfig) getStringArray(key, def, deli string) []string {
	v := this.getValueDef(key, def)
	if v == "" {
		return []string{}
	}
	tokens := stringutil.Tokenizer(v, deli)
	// trim Space
	trimTokens := make([]string, 0)
	for _, v := range tokens {
		trimTokens = append(trimTokens, strings.TrimSpace(v))
	}
	return trimTokens
}
func (this *FileConfig) GetFloat(key string, def float32) float32 {
	return this.getFloat(key, def)
}
func (this *FileConfig) getFloat(key string, def float32) float32 {
	v := this.getValue(key)
	if v == "" {
		return float32(def)
	}
	value, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return float32(def)
	}
	return float32(value)
}

func (this *FileConfig) InArray(str string, list []string) bool {
	for _, it := range list {
		if strings.TrimSpace(str) == strings.TrimSpace(it) {
			return true
		}
	}
	return false
}

func (this *FileConfig) SetValues(keyValues *map[string]string) {
	path := this.GetConfFile()

	// 쓰기 전에 최종 conf merge
	tmp, err := this.conf.Parser.Read(path)
	if err != nil {
		this.conf.Log.Error("Paser read error ", err)
		return
	}

	// append prefix and suffix
	for k, v := range *keyValues {
		if this.isExcludeKeys(k) {
			continue
		}
		if this.conf.prefix != "" {
			if !strings.HasPrefix(k, this.conf.prefix) {
				k = fmt.Sprintf("%s%s", this.conf.prefix, k)
			}
		}
		if this.conf.suffix != "" {
			if !strings.HasSuffix(k, this.conf.suffix) {
				k = fmt.Sprintf("%s%s", k, this.conf.suffix)
			}
		}
		tmp[k] = v
	}

	//this.apply(tmp)

	err = this.conf.Parser.Write(path, &tmp)
	if err != nil {
		this.conf.Log.Error("Paser write error ", err)
		return
	}

}

func (this *FileConfig) isExcludeKeys(key string) bool {
	for _, k := range this.conf.excludeKeys {
		if k == key {
			return true
		}
	}
	return false
}

// func SearchKey(keyPrefix string) *map[string]string {
// 	keyValues := map[string]string{}
// 	for _, key := range prop.Keys() {
// 		if strings.HasPrefix(key, keyPrefix) {
// 			if v, ok := prop.Get(key); ok {
// 				keyValues[key] = v
// 			}
// 		}
// 	}

// 	return &keyValues
// }

// func FilterPrefix(keyPrefix string) map[string]string {
// 	keyValues := make(map[string]string)
// 	//php prefix whatap.
// 	if this.AppType == lang.APP_TYPE_PHP {
// 		if !strings.HasPrefix(keyPrefix, "whatap.") {
// 			keyPrefix = "whatap." + keyPrefix
// 		}
// 	} else if this.AppType == lang.APP_TYPE_BSM_PHP {
// 		if !strings.HasPrefix(keyPrefix, "opsnowbsm.") {
// 			keyPrefix = "opsnowbsm." + keyPrefix
// 		}
// 	}
// 	pp := prop.FilterPrefix(keyPrefix)
// 	for _, key := range pp.Keys() {
// 		keyValues[key] = pp.GetString(key, "")
// 	}
// 	return keyValues
// }

// func cutOut(val, delim string) string {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			this.this.Log.Println("WA217", " Recover ", r)
// 		}
// 	}()
// 	if val == "" {
// 		return val
// 	}
// 	x := strings.LastIndex(val, delim)
// 	if x <= 0 {
// 		return ""
// 	}
// 	//return val.substring(0, x);
// 	return val[0:x]

// }

// func toHashSet(key, def string) *hmap.IntSet {
// 	set := hmap.NewIntSet()
// 	vv := strings.Split(getValueDef(key, def), ",")
// 	if vv != nil {
// 		for _, x := range vv {
// 			func() {
// 				defer func() {
// 					if r := recover(); r != nil {
// 						this.this.Log.Infoln("WA218", " Recover ", r)
// 					}
// 				}()

// 				x = strings.TrimSpace(x)
// 				if len(x) > 0 {
// 					xx := hash.HashStr(x)
// 					set.Put(xx)
// 				}
// 			}()
// 		}
// 	}
// 	return set
// }

// func toStringSet(key, def string) *hmap.StringSet {
// 	set := hmap.NewStringSet()
// 	vv := strings.Split(getValueDef(key, def), ",")
// 	if vv != nil {
// 		for _, x := range vv {
// 			func() {
// 				defer func() {
// 					if r := recover(); r != nil {
// 						this.this.Log.Infoln("WA219", " Recover ", r)
// 					}
// 				}()
// 				x = strings.TrimSpace(x)
// 				if len(x) > 0 {
// 					set.Put(x)
// 				}
// 			}()
// 		}
// 	}
// 	return set
// }

// func IsIgnoreTrace(hash int32, service string) bool {
// 	if this.TraceIgnoreUrlSet.Contains(hash) {
// 		return true
// 	}
// 	if this.IsTraceIgnoreUrlPrefix {
// 		if strings.HasPrefix(service, this.TraceIgnoreUrlPrefix) {
// 			return true
// 		}
// 	}
// 	return false
// }

func (this *FileConfig) ToString() string {
	return this.String()
}
func (this *FileConfig) String() string {
	sb := stringutil.NewStringBuffer()
	for k, v := range this.m {
		sb.Append(k).Append("=").AppendLine(v)
	}
	return sb.ToString()
}
