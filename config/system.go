package config

type System struct {
	Name           string `mapstructure:"name"`
	DSN            string `mapstructure:"dsn"`
	PprofHost      string `mapstructure:"pprof-host"`
	JaegerHost     string `mapstructure:"jaeger-host"`
	ServerPort     string `mapstructure:"server-port"`
	HttpPort       string `mapstructure:"http-port"`
	PrometheusHost string `mapstructure:"prometheus-host"`
}

type Conf interface {
	GetSystemConfig() *System
}

type Helper struct {
	System
}

func (ch Helper) GetSystemConfig() *System {
	return &ch.System
}
