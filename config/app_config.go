package config

// GlobalConf is the interface that application configs must implement
type GlobalConf interface {
	GetStaticConfig() interface{}
	GetDynamicConfig() interface{}
}

// AppConfig is the top-level application configuration
type AppConfig struct {
	StaticConfig  StaticConfig
	DynamicConfig DynamicConfig
}

func (c *AppConfig) GetStaticConfig() interface{} {
	return &c.StaticConfig
}

func (c *AppConfig) GetDynamicConfig() interface{} {
	return &c.DynamicConfig
}

type StaticConfig struct {
	App struct {
		Name               string  `mapstructure:"name"`
		Port               string  `mapstructure:"port"`
		Environment        string  `mapstructure:"env"`
		LogLevel           string  `mapstructure:"log_level"`
		MetricSamplingRate float64 `mapstructure:"metric_sampling_rate"`
	} `mapstructure:"app"`
	Auth struct {
		JWTSecret      string `mapstructure:"jwt_secret"`
		GoogleClientID string `mapstructure:"google_client_id"`
	} `mapstructure:"auth"`
}

type DynamicConfig struct{}
