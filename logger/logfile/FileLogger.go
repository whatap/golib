package logfile

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"context"
	"io/ioutil"
	"path/filepath"
	"runtime/debug"

	"github.com/whatap/golib/config"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/logger"
	"github.com/whatap/golib/util/ansi"
	"github.com/whatap/golib/util/dateutil"
	"github.com/whatap/golib/util/hmap"
	"github.com/whatap/golib/util/stringutil"
)

const (
	HOME_ENV_KEY             = "WHATAP_HOME"
	defaultLevel             = logger.LOG_LEVEL_WARN
	defaultRotationEnabled   = true
	defaultCacheInterval     = 10
	defaultKeepDays          = 7
	defaultOname             = "boot"
	defaultLogID             = "whatap"
	defaultLogIDPrefixLength = 10
)

type LogData struct {
	Before int64
	Next   int64
	Text   string
}

func NewLogData(pre, next int64, text string) *LogData {
	p := new(LogData)
	p.Before = pre
	p.Next = next
	p.Text = text

	return p
}

type FileLogger struct {
	ctx            context.Context
	cancel         context.CancelFunc
	configObserver *config.ConfigObserver
	conf           fileLoggerConfig
	myLog          *log.Logger

	lastLog          *hmap.StringLongLinkedMap
	lock             sync.Mutex
	logfile          *os.File
	last             int64
	lastDataUnit     int64
	lastFileRotation bool
}

var defaultFileLogger *FileLogger
var fileLoggerMutex sync.Mutex

// 패키지 로드
//func init() {
//    logger = NewLogger()
//}

func GetFileLogger(opts ...FileLoggerOption) *FileLogger {
	fileLoggerMutex.Lock()
	defer fileLoggerMutex.Unlock()

	if defaultFileLogger != nil {
		return defaultFileLogger
	} else {
		defaultFileLogger = NewFileLogger(opts...)
		return defaultFileLogger
	}
}

func NewFileLogger(opts ...FileLoggerOption) *FileLogger {
	p := new(FileLogger)
	p.myLog = log.New(os.Stdout, "", log.LstdFlags)
	p.conf = defaultFileLoggerConfig()
	for _, opt := range opts {
		opt.apply(&p.conf)
	}
	p.lastLog = hmap.NewStringLongLinkedMap().SetMax(1000)
	if p.conf.homePath == "" {
		p.conf.homePath = os.Getenv("WHATAP_HOME")
		if p.conf.homePath == "" {
			p.conf.homePath = "./"
		}
	}

	if p.configObserver != nil {
		p.configObserver.Add(fmt.Sprintf("FileLogger-%s-%s-%s", p.conf.homePath, p.conf.logID, p.conf.oname), p)
	}
	// open File
	p.openFile()
	// 파일 로그 설
	go p.run()

	return p
}

func (this *FileLogger) GetLogFile() *os.File {
	return this.logfile
}
func (this *FileLogger) GetLogFilePath() string {
	return filepath.Join(this.conf.homePath, this.logfile.Name())
}

func (this *FileLogger) run() {
	this.last = dateutil.Now()
	this.lastDataUnit = dateutil.GetDateUnitNow()
	this.lastFileRotation = this.conf.rotationEnabled
	for {
		this.process()
		time.Sleep(10000 * time.Millisecond)
	}
}

func (this *FileLogger) process() {
	this.lock.Lock()
	defer func() {
		this.lock.Unlock()
		if r := recover(); r != nil {
			this.Error("WA10005", " Recover", r) //, string(debug.Stack()))
		}
	}()

	now := dateutil.Now()
	//fmt.Printf("FileLogger process oname=%s, now=%d \r\n", this.oname, now)

	//if now > this.last+dateutil.MILLIS_PER_HOUR {
	if now > this.last+dateutil.MILLIS_PER_MINUTE {
		this.last = now
		this.clearOldLog()
	}

	if (this.lastFileRotation != this.conf.rotationEnabled) || (this.lastDataUnit != dateutil.GetDateUnitNow()) || (this.logfile == nil) {
		this.logfile.Close()
		this.logfile = nil
		this.lastFileRotation = this.conf.rotationEnabled
		this.lastDataUnit = dateutil.GetDateUnitNow()
	}
	this.openFile()
}

func (this *FileLogger) openFile() {
	defer func() {
		if r := recover(); r != nil {
			this.Error("WA10004", "openFile Recover", r)
		}
	}()

	if this.logfile == nil {
		//fmt.Println("Logger open file", "oname=", this.oname, "filname=", fmt.Sprintf("whatap-%s-%s.log", this.oname, dateutil.YYYYMMDD(dateutil.Now())))
		// 로그파일 오픈
		home := this.conf.homePath

		if _, err := os.Stat(filepath.Join(home, "logs")); err != nil {
			if os.IsNotExist(err) {
				// file does not exist
				os.Mkdir(filepath.Join(home, "logs"), os.ModePerm)
			} else {
				// other error
			}
		}
		var file *os.File
		var err error
		if this.conf.rotationEnabled {
			file, err = os.OpenFile(filepath.Join(home, "logs", fmt.Sprintf("%s-%s-%s.log", this.conf.logID, this.conf.oname, dateutil.YYYYMMDD(dateutil.Now()))), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
		} else {
			file, err = os.OpenFile(filepath.Join(home, "logs", fmt.Sprintf("%s-%s.log", this.conf.logID, this.conf.oname)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
		}
		this.logfile = file
		//fmt.Println("Logger open file", this.logfile)

		// 로거를 파일로그로 변경
		if this.conf.IsStdout {
			multi := io.MultiWriter(this.logfile, os.Stdout)
			this.myLog.SetOutput(multi)
		} else {
			this.myLog.SetOutput(this.logfile)
		}
		this.myLog.SetFlags(log.Ldate | log.Ltime)

		this.myLog.Println("")
		this.myLog.Println("## OPEN LOG FILE ", this.conf.oname, "", dateutil.TimeStampNow()+" ##")
		this.myLog.Println("")
	}

	//defer logfile.Close()
}

func (this *FileLogger) Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	id := stringutil.Truncate(s, defaultLogIDPrefixLength)
	if this.checkOk(id, this.conf.cacheInterval) == false {
		return
	}
	this.myLog.Println(ansi.Red(fmt.Sprintf("%s %s", "[Error]", s)))
}
func (this *FileLogger) Error(args ...interface{}) {
	s := fmt.Sprintln(args...)
	id := stringutil.Truncate(s, defaultLogIDPrefixLength)
	if this.checkOk(id, this.conf.cacheInterval) == false {
		return
	}
	this.myLog.Println(ansi.Red(fmt.Sprintf("%s %s", "[Error]", s)))
}
func (this *FileLogger) Warnf(format string, args ...interface{}) {
	if this.conf.level > logger.LOG_LEVEL_WARN {
		return
	}
	s := fmt.Sprintf(format, args...)
	id := stringutil.Truncate(s, defaultLogIDPrefixLength)
	if this.checkOk(id, this.conf.cacheInterval) == false {
		return
	}
	this.myLog.Println("[Warn] ", s)
}
func (this *FileLogger) Warn(args ...interface{}) {
	if this.conf.level > logger.LOG_LEVEL_WARN {
		return
	}
	s := fmt.Sprintln(args...)
	id := stringutil.Truncate(s, defaultLogIDPrefixLength)
	if this.checkOk(id, this.conf.cacheInterval) == false {
		return
	}
	this.myLog.Println("[Warn] ", s)
}
func (this *FileLogger) Infof(format string, args ...interface{}) {
	if this.conf.level > logger.LOG_LEVEL_INFO {
		return
	}
	s := fmt.Sprintf(format, args...)
	id := stringutil.Truncate(s, defaultLogIDPrefixLength)
	if this.checkOk(id, this.conf.cacheInterval) == false {
		return
	}
	this.myLog.Println("[Info] ", s)
}
func (this *FileLogger) Info(args ...interface{}) {
	if this.conf.level > logger.LOG_LEVEL_INFO {
		return
	}
	s := fmt.Sprintln(args...)
	id := stringutil.Truncate(s, defaultLogIDPrefixLength)
	if this.checkOk(id, this.conf.cacheInterval) == false {
		return
	}
	this.myLog.Println("[Info] ", s)
}
func (this *FileLogger) Infoln(args ...interface{}) {
	if this.conf.level > logger.LOG_LEVEL_INFO {
		return
	}
	s := fmt.Sprintln(args...)
	id := stringutil.Truncate(s, defaultLogIDPrefixLength)
	if this.checkOk(id, this.conf.cacheInterval) == false {
		return
	}
	this.myLog.Println("[Info] ", s)
}

//debug level, nocache
func (this *FileLogger) Debugf(format string, args ...interface{}) {
	if this.conf.level > logger.LOG_LEVEL_DEBUG {
		return
	}
	s := fmt.Sprintf(format, args...)
	this.myLog.Println("[Debug] ", s)
}

//debug level, nocache
func (this *FileLogger) Debug(args ...interface{}) {
	if this.conf.level > logger.LOG_LEVEL_DEBUG {
		return
	}
	s := fmt.Sprintln(args...)
	this.myLog.Println("[Debug] ", s)
}

// 첫번째 인수는 무조건 String으로 ID 값을 넣어야 함( WA111 형식)
// 해당 ID로 중복 확인.

func (this *FileLogger) Printf(id string, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	this.println(id, this.build(id, s))
}
func (this *FileLogger) Println(id string, args ...interface{}) {
	s := fmt.Sprintln(args...)
	this.println(id, this.build(id, s))
}

func (this *FileLogger) println(id, message string) {
	if this.checkOk(id, this.conf.cacheInterval) == false {
		return
	}
	this.myLog.Println(this.build(id, message))
}

func (this *FileLogger) build(id, message string) string {
	return fmt.Sprint("[", id, "] ", message)
}

func (this *FileLogger) getCallStack() string {
	defer func() {
		if r := recover(); r != nil {
			this.Error("WA10001", "getCallStack Recover ", r)
		}
	}()
	return string(debug.Stack())
}

func (this *FileLogger) checkOk(id string, sec int) bool {
	if sec > 0 {
		last := this.lastLog.Get(id)
		now := dateutil.Now()
		if now < (last + int64(sec)*1000) {
			return false
		}
		this.lastLog.Put(id, now)
	}
	return true
}

func (this *FileLogger) sysout(message string) {
	fmt.Println(message)
}

func (this *FileLogger) PrintlnStd(msg string, sysout bool) {
	defer func() {
		if r := recover(); r != nil {
			this.myLog.Println("WA10002", "println Recover", r)
		}
	}()
	if sysout {
		fmt.Println(msg)
	} else {
		this.myLog.Println(msg)
	}
}

// func Update(oname string) {
// 	logger.update(oname)
// }
func (this *FileLogger) update(oname string) {
	defer func() {
		if r := recover(); r != nil {
			this.myLog.Println("WA10003", "Update Recover", r)
		}
	}()

	oname = strings.TrimSpace(oname)
	if oname == this.conf.oname {
		return
	}

	this.conf.oname = oname
	this.openFile()
}

func (this *FileLogger) clearOldLog() {
	if this.conf.rotationEnabled == false {
		return
	}
	if this.conf.keepDays <= 0 {
		return
	}
	whatapPrefix := this.conf.logID
	nowUnit := dateutil.GetDateUnitNow()

	home := this.conf.homePath
	searchDir := filepath.Join(home, "logs")

	// Get filelist
	files, _ := ioutil.ReadDir(searchDir)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		// prefix 구분
		//fmt.Printf("file=%s", f.Name())
		if !strings.HasPrefix(name, whatapPrefix+"-") {
			continue
		}
		// oname 을 구분하지 않고 날짜만 확인 해서 모두 정리
		x := strings.LastIndex(name, ".")
		if x < 0 {
			continue
		}

		s := strings.LastIndex(name, "-")
		//s >= x-1  적어도 한 문자는 slice 되게
		if s < 0 || s >= x-1 {
			continue
		}
		date := name[s+1 : x]

		//fmt.Printf("file=%s, date=%s", f.Name(), date)

		if len(date) != 8 {
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					this.myLog.Println("WA10006", " File Delete Error", r)
				}
			}()

			d := dateutil.GetYmdTime(date)
			fileUnit := dateutil.GetDateUnit(d)
			if nowUnit-fileUnit > int64(this.conf.keepDays) {
				//fmt.Println("File Remove", filepath.Join(searchDir,f.Name()))
				err := os.Remove(filepath.Join(searchDir, f.Name()))
				if err != nil {
					this.Error("WA10007", " File Remove Error", err)
				}
			}
		}()
	}
}

func (this *FileLogger) GetLogFiles() *value.MapValue {
	out := value.NewMapValue()

	whatapPrefix := this.conf.logID + "-" + this.conf.oname
	home := this.conf.homePath
	searchDir := filepath.Join(home, "logs")

	// Get filelist
	files, _ := ioutil.ReadDir(searchDir)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()

		x := strings.Index(name, ".")
		if x < 0 {
			continue
		}
		if name != "whatap-hook.log" {
			if !strings.HasPrefix(name, whatapPrefix+"-") {
				continue
			}
			date := name[len(whatapPrefix)+1 : x]

			if len(date) != 8 {
				continue
			}
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					this.Println("WA10008", " File Delete Error", r)
				}
			}()
			out.Put(f.Name(), value.NewDecimalValue(f.Size()))
		}()

		if out.Size() >= 100 {
			break
		}
	}

	searchDotnetPath := filepath.Join(os.Getenv("ProgramData"), "WhaTap", "dotnet-profiler.log")
	fi, err := os.Stat(searchDotnetPath)
	if err == nil {
		out.Put(fi.Name(), value.NewDecimalValue(fi.Size()))
	}

	return out
}

func (this *FileLogger) Read(file string, endpos int64, length int64) *LogData {
	var ret string

	if file == "" || length == 0 {
		return nil
	}

	// 폴더 없을 때 발생하는 오류를 임시로 조정.
	// logs 폴더 없는 경우 생성
	if _, err := os.Stat(filepath.Join(this.conf.homePath, "logs")); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			os.Mkdir(filepath.Join(this.conf.homePath, "logs"), os.ModePerm)
		} else {
			// other error
		}
	}

	searchFilePath := filepath.Join(this.conf.homePath, "logs", file)
	f, err := os.Open(searchFilePath)
	if err != nil {
		this.Error("Read log file ", err)
		return nil
	}

	fInfo, err := f.Stat()
	if fInfo.Size() < endpos {
		return nil
	}

	if endpos < 0 {
		endpos = fInfo.Size()
	}
	start := int64(math.Max(0, float64(endpos-length)))

	available := fInfo.Size() - start
	readable := int(math.Min(float64(available), float64(length)))
	//readable = int(math.Min(math.MinInt16, float64(readable)))

	buff := make([]byte, readable)

	n, err := f.ReadAt(buff, start)

	//this.Println("FileLogger Read ", "file=", file, ",size=", fInfo.Size(), "readable=", readable, ",endpos=", endpos, ",start=", start, ",length=", length, "read=", n) //, ",result=" + string(buff));

	if err != nil {
		this.Error("WA1000901", " Read Error ", err)
		return nil
	}
	ret = string(buff)

	next := start + int64(n)

	if (next + length) > fInfo.Size() {
		next = -1
	} else {
		next += length
	}

	defer func() {
		f.Close()
		if r := recover(); r != nil {
			this.Error("WA10009", " Read Recover", r)
			ret = ""
		}
	}()

	//return ret
	return NewLogData(start, next, ret)
}

func (this *FileLogger) ApplyConfig(conf config.Config) {
	//this.Pcode = conf.GetLong("pcode")
	this.conf.rotationEnabled = conf.GetBoolean("log_rotation_enabled", true)
	this.conf.keepDays = int(conf.GetInt("log_keep_days", 7))
	this.conf.cacheInterval = int(conf.GetInt("_log_interval", 10))
}
