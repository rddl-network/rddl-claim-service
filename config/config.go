package config

import "sync"

const DefaultConfigTemplate = `
service-bind="{{ .ServiceBind }}"
service-port={{ .ServicePort }}
`

type Config struct {
	ServicePort int    `mapstructure:"service-port"`
	ServiceBind string `mapstructure:"service-bind"`
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
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}
