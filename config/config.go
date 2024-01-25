package config

import "sync"

const DefaultConfigTemplate = `
service-bind="{{ .ServiceBind }}"
service-port={{ .ServicePort }}
db-path="{{ .DBPath }}"
`

type Config struct {
	ServicePort int    `mapstructure:"service-port"`
	ServiceBind string `mapstructure:"service-bind"`
	DBPath      string `mapstructure:"db-path"`
}

// global singleton
var (
	config     *Config
	initConfig sync.Once
)

// DefaultConfig returns RDDL-2-PLMNT default config
func DefaultConfig() *Config {
	return &Config{
		ServicePort: 8080,
		ServiceBind: "localhost",
		DBPath:      "./data",
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}
