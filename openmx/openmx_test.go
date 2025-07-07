package openmx

import (
	// "fmt"
	// "os"
	// "strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/whatap/golib/lang/pack"
	"github.com/whatap/golib/lang/pack/open"
	"github.com/whatap/golib/logger"
	// "github.com/whatap/golib/net"
	// "github.com/whatap/golib/net/oneway"
)

var (
	// prod
	accesskey string = "x4g2e21sj8a9m-x2iglc38nmd1fa-x2vunh5jbms4u1"
	servers          = []string{"13.124.11.223:6600", "13.209.172.35:6600"}
	pcode            = int64(16462)

	// dev
	// accesskey string = "x41nf22n9j9kf-xhldoeeo9f018-z53eh3v2kmi5s2"
	// servers          = []string{"15.165.146.117:6600", "15.165.146.117:6600"}
	// pcode            = int64(1775)

	Log = logger.NewDefaultLogger()
)

func SendOneway(p pack.Pack) {
	// Log.SetLevel(logger.LOG_LEVEL_DEBUG)
	// client := oneway.GetOneWayTcpClient(
	// 	oneway.WithLicense(accesskey),
	// 	oneway.WithPcode(pcode),
	// 	oneway.WithServers(servers),
	// 	oneway.WithOid(2345667),
	// 	oneway.WithLogger(Log))

	// fmt.Println("client.Oid=", client.Oid)
	// p.SetPCODE(client.Pcode)
	// p.SetOID(client.Oid)
	// p.SetTime(time.Now().UnixMilli())

	// if err := client.SendFlush(p, true); err != nil {
	// 	fmt.Println("error = ", err)
	// }
}

func TestOpenMetricSend(t *testing.T) {
	om := New()
	// hostname, _ := os.Hostname()
	label := []*open.Label{}
	idx := 0
	tm := time.Now().UnixMilli() / 1000 * 1000

	// for {
	name := "test_test_gauge"
	om.AddHelp(name, "help asdfggggg", "gauge")
	m := open.NewOpenMxWithLabel(name, tm, 3.33+float64(idx), label...)
	om.AddMetric(m)
	idx++
	// if idx > 5 {
	// 	break
	// }
	// time.Sleep(5 * time.Second)
	// }

	assert.Equal(t, 1, len(om.helps))
	assert.Equal(t, 1, len(om.metrics))

	om.Send(tm, SendOneway)

	// time.Sleep(5 * time.Second)

}

func TestOpenMetricSummerySend(t *testing.T) {
	om := New()
	// hostname, _ := os.Hostname()
	tm := time.Now().UnixMilli() / 1000 * 1000

	label := []*open.Label{}
	name := "request_latency_summary_sum"
	om.AddHelp(name, "help request_latency_summary_sum", "summary")
	label = append(label, open.NewLabel("endpoint", "/api"))
	m := open.NewOpenMxWithLabel(name, tm, 450.000000, label...)
	om.AddMetric(m)

	label = []*open.Label{}
	name = "request_latency_summary_count"
	om.AddHelp(name, "help request_latency_summary_count", "summary")
	label = append(label, open.NewLabel("endpoint", "/api"))
	m = open.NewOpenMxWithLabel(name, tm, 200.000000, label...)
	om.AddMetric(m)

	label = []*open.Label{}
	name = "request_latency_summary"
	om.AddHelp(name, "help request_latency_summary_count", "summary")
	label = append(label, open.NewLabel("endpoint", "/api"))
	label = append(label, open.NewLabel("quantile", "0.50"))
	m = open.NewOpenMxWithLabel(name, tm, 200.000000, label...)
	om.AddMetric(m)

	label = []*open.Label{}
	name = "request_latency_summary"
	om.AddHelp(name, "help request_latency_summary_count", "summary")
	label = append(label, open.NewLabel("endpoint", "/api"))
	label = append(label, open.NewLabel("quantile", "0.90"))
	m = open.NewOpenMxWithLabel(name, tm, 2.300000, label...)
	om.AddMetric(m)

	label = []*open.Label{}
	name = "request_latency_summary"
	om.AddHelp(name, "help request_latency_summary_count", "summary")
	label = append(label, open.NewLabel("endpoint", "/api"))
	label = append(label, open.NewLabel("quantile", "0.99"))
	m = open.NewOpenMxWithLabel(name, tm, 3.000000, label...)
	om.AddMetric(m)

	assert.Equal(t, 5, len(om.helps))
	assert.Equal(t, 5, len(om.metrics))

	om.Send(tm, SendOneway)

	// time.Sleep(5 * time.Second)

}
