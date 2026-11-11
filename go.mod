module github.com/whatap/golib

go 1.25.0

// 서버의 Go 가 이보다 낮을 때 자동으로 내려받는 툴체인. 지정하지 않으면 go1.25.0 을 받는데,
// 그 버전은 `go test -cover` 가 `no such tool "covdata"` 로 실패한다.
toolchain go1.26.10

require (
	github.com/google/uuid v1.3.0
	github.com/magiconair/properties v1.8.7
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/stretchr/testify v1.9.0
	golang.org/x/text v0.41.0
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sys v0.44.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
