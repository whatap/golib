package config

type Configure interface {
	ApplyConfig(conf Config)
}

type ConfigObserver struct {
	observer map[string]Configure
}

func GetConfigObserver() *ConfigObserver {
	return defaultConfigObserver
}

var (
	defaultConfigObserver = NewConfigObserver()
)

func NewConfigObserver() *ConfigObserver {
	p := new(ConfigObserver)
	p.observer = make(map[string]Configure)
	return p
}

func (this *ConfigObserver) Add(cls string, conf Configure) {
	this.observer[cls] = conf
}

func (this *ConfigObserver) Run(conf Config) {
	for _, v := range this.observer {
		v.ApplyConfig(conf)
	}
}

// type Config interface {
// 	GetLong(k string, def int64) int64
// }

// type Runnable interface {
// 	Run()
// }

// var observer map[string]Runnable = make(map[string]Runnable)

// func AddConfObserver(cls string, run Runnable) {
// 	//fmt.Println("Add=", cls)
// 	observer[cls] = run
// }
// func RunConfObserver() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			logutil.Println("WA10500"," Recover", r)
// 		}
// 	}()

// 	//fmt.Println("Run=")

// 	for _, v := range observer {
// 		//fmt.Println("Run=", k)
// 		v.Run()
// 	}
// }
