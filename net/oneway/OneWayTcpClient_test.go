package oneway

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whatap/golib/config"
	wio "github.com/whatap/golib/io"
	"github.com/whatap/golib/lang/pack"
	wnet "github.com/whatap/golib/net"
	whash "github.com/whatap/golib/util/hash"
)

var (
	testServer  = ""
	testPcode   = make([]int64, 0)
	testLicense = make([]string, 0)
	testOid     = int32(0)
	setOption   = false
)

// TODO UT Case 검토 필요

// Connection 제외 로직 검증 케이스
// 실제 외부에서 호출하는 코드 X, 검증용

func TestGetOneWayClientWithServers(t *testing.T) {
	assert := assert.New(t)
	servers := make([]string, 0)
	servers = append(servers, "127.0.0.1:6600")
	pcode := int64(00000)
	oid := int32(123)
	license := "abcdefg"

	p := pack.NewTagCountPack()
	p.SetPCODE(10)

	oneway := newOneWayTcpClient(WithServers(servers), WithPcode(12345), WithPcode(pcode), WithOid(oid), WithLicense(license), WithUseQueue())
	defer oneway.Close()
	defer oneway.Destroy()

	assert.Equal(servers, oneway.Servers)
	assert.Equal(pcode, oneway.Pcode)
	assert.Equal(oid, oneway.Oid)
	assert.Equal(license, oneway.License)
	assert.True(oneway.UseQueue)

	oneway.Send(p, wnet.WithLicense("hijklmn"))

	if tmp := oneway.Queue.GetTimeout(int(queueMaxWaitTime)); tmp != nil {
		if tcpSend, ok := tmp.(*wnet.TcpSend); ok {
			dout := oneway.makeData(&wnet.TcpSend{Pack: p, Opts: tcpSend.Opts})

			//Send 옵션으로 전달된 License로 Data 생성 확인
			v := whash.Hash64Str("hijklmn")
			vArr := wio.ToBytesLong(v)
			v2 := dout.ToByteArray()[10:18]
			assert.Equal(vArr, v2)
			assert.True(dout.ToByteArray()[9] == 10)
		}

	}

	// new 과정에서 만들어진 옵션값 변동 없음
	assert.Equal(oneway.License, license)

}

// Test 환경 필요
// 실제 사용 사례

func TestSingleConnect(t *testing.T) {

	if setOption == false {
		t.Skip("Neet To Test Option..")
	}

	assert := assert.New(t)
	servers := make([]string, 0)
	servers = append(servers, testServer)

	p := pack.NewCounterPack1()
	p.SetPCODE(int64(testPcode[0]))

	oneway := GetOneWayTcpClient(WithServers(servers), WithLicense(testLicense[0]), WithOid(testOid))
	defer func() {
		oneway.Close()
		oneway.Destroy()
	}()

	assert.Equal(oneway.Servers[0], servers[0])
	assert.Equal(oneway.License, testLicense[0])
	assert.Equal(oneway.Oid, testOid)

	err := oneway.Send(p)

	assert.Nil(err)
}

func TestMultiConnect(t *testing.T) {

	if setOption == false {
		t.Skip("Neet To Test Option..")
	}

	assert := assert.New(t)
	servers := make([]string, 0)
	servers = append(servers, testServer)

	var wg sync.WaitGroup

	oneway := GetOneWayTcpClient(WithServers(servers), WithLicense(testLicense[0]), WithOid(testOid))
	defer func() {
		oneway.Close()
		oneway.Destroy()
	}()

	for i := 0; i < len(testPcode); i++ {
		wg.Add(1)
		go func(idx int) {
			p := pack.NewCounterPack1()
			p.SetPCODE(testPcode[idx])

			err := oneway.Send(p, wnet.WithLicense(testLicense[idx]))

			assert.Nil(err)
			wg.Done()
		}(i)
	}
	wg.Wait()

	//assert 초기값 변화없음 확인
	assert.Equal(oneway.Servers[0], servers[0])
	assert.Equal(oneway.License, testLicense[0])
	assert.Equal(oneway.Oid, testOid)
}

func TestConfigApply(t *testing.T) {
	assert := assert.New(t)
	conf := &config.MockConfig{}
	servers := make([]string, 0)
	servers = append(servers, testServer)

	oneway := GetOneWayTcpClient(WithServers(servers), WithLicense(testLicense[0]), WithOid(testOid))
	defer func() {
		oneway.Close()
		oneway.Destroy()
	}()

	conf.Mock.Test(t)
	conf.On("GetValue", mock.AnythingOfType("string")).
		Return(func(s string) string {
			if s == "license" {
				return testLicense[1]
			} else if s == "whatap.server.host" {
				return strings.Split(testServer, ":")[0]

			}
			return s
		})

	conf.On("GetValueDef", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(func(s string, v string) string {
			return v
		})

	conf.On("GetLong", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
		Return(func(s string, v int64) int64 {
			return v
		})

	conf.On("GetInt", mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Return(func(s string, v int) int32 {
			return int32(v)
		})

	oneway.ApplyConfig(conf)

	assert.Equal(oneway.License, testLicense[1])
}

func TestMain(m *testing.M) {
	// Option 처리
	testServer = "15.165.146.117:6600"
	testPcode = append(testPcode, 78)
	testPcode = append(testPcode, 79)
	testLicense = append(testLicense, "x2jgg66m4jlck-z6l4o2nb3cckq0-x5jfk4ktaqmfth")
	testLicense = append(testLicense, "x2jogsfu22s10-x52c23sgmuab82-z38tcf60s1i202")
	testOid = 123

	if testServer != "" && len(testPcode) != 0 && len(testLicense) != 0 && testOid != 0 {
		setOption = true
	}
	m.Run()

}
