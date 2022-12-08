package conffile

import (
	"context"
	_ "log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/whatap/golib/config"
	_ "github.com/whatap/golib/logger"
)

func TestDefaultFileParser(t *testing.T) {
	var notnil bool
	DeleteConfigFile("whatap.conf")

	err := OpenConfigFile("whatap.conf", "debug=false\n", 0)
	assert.NoError(t, err)

	conf := GetConfig()
	defer conf.Destroy()

	notnil = assert.NotNil(t, conf)
	if !notnil {
		return
	}
	assert.Equal(t, "whatap.conf", conf.GetConfFile())
	assert.Equal(t, false, conf.GetBoolean("debug", true))
}

func TestDefaultFileParserOberver(t *testing.T) {
	var notnil bool
	DeleteConfigFile("whatap.conf")
	err := OpenConfigFile("whatap.conf", "debug=false\n", 0)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.TODO())
	observer := config.GetConfigObserver()

	notnil = assert.NotNil(t, observer)
	if !notnil {
		return
	}

	// add target
	noti := func(b bool) bool {
		if b {
			cancel()
		}
		return b
	}

	target := &CallbackFromObserver{
		NotiFunc: noti,
	}
	observer.Add("TestTarget", target)

	//conf := GetConfig(WithConfigObserver(observer), WithLogger(logger.NewDefaultLogger()))
	conf := GetConfig(WithConfigObserver(observer))
	defer conf.Destroy()

	notnil = assert.NotNil(t, conf)
	if !notnil {
		return
	}

	assert.Equal(t, "whatap.conf", conf.GetConfFile())
	assert.Equal(t, false, conf.GetBoolean("debug", true))

	go OpenConfigFile("whatap.conf", "noti=true\n", 5)

	for {
		select {
		case <-ctx.Done():
			time.Sleep(2 * time.Second)
			return
		case <-time.After(30 * time.Second):
			// 30초 경과면 무조건 실패
			assert.Equal(t, true, false)
			return
		}
	}
}

func TestDeleteConfigFile(t *testing.T) {
	assert.NoError(t, DeleteConfigFile("whatap.conf"))
}

type CallbackFromObserver struct {
	NotiFunc func(bool) bool
}

func (this *CallbackFromObserver) ApplyConfig(conf config.Config) {
	this.NotiFunc(conf.GetBoolean("noti", false))
}

func DeleteConfigFile(path string) error {
	return os.Remove(path)
}
func OpenConfigFile(path string, str string, delay int) error {
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}

	if f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		defer f.Close()
		f.WriteString(str)
		return nil
	} else {
		return err
	}
}
