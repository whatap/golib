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
