package main

import (
	"fmt"

	_ "github.com/whatap/golib/config"
	_ "github.com/whatap/golib/config/conffile"
	_ "github.com/whatap/golib/io"
	_ "github.com/whatap/golib/lang/pack"
	_ "github.com/whatap/golib/lang/pack/udp"
	_ "github.com/whatap/golib/lang/ref"
	_ "github.com/whatap/golib/lang/service"
	_ "github.com/whatap/golib/lang/step"
	_ "github.com/whatap/golib/lang/value"
	_ "github.com/whatap/golib/lang/variable"
	_ "github.com/whatap/golib/logger"
	_ "github.com/whatap/golib/logger/logfile"
	_ "github.com/whatap/golib/logsink/zip"

	_ "github.com/whatap/golib/net"
	_ "github.com/whatap/golib/net/oneway"
	_ "github.com/whatap/golib/net/udp"
	_ "github.com/whatap/golib/util/bitutil"
	_ "github.com/whatap/golib/util/cmdutil"
	_ "github.com/whatap/golib/util/compare"
	_ "github.com/whatap/golib/util/dateutil"
	_ "github.com/whatap/golib/util/hash"
	_ "github.com/whatap/golib/util/hexa32"
	_ "github.com/whatap/golib/util/hll"
	_ "github.com/whatap/golib/util/hmap"
	_ "github.com/whatap/golib/util/iputil"
	_ "github.com/whatap/golib/util/keygen"
	_ "github.com/whatap/golib/util/list"
	_ "github.com/whatap/golib/util/mathutil"
	_ "github.com/whatap/golib/util/openstack"
	_ "github.com/whatap/golib/util/paramtext"
	_ "github.com/whatap/golib/util/queue"
	_ "github.com/whatap/golib/util/stringutil"
	_ "github.com/whatap/golib/util/urlutil"
	_ "github.com/whatap/golib/util/uuidutil"
)

func main() {
	fmt.Println("Whatap Golang common library")
}
