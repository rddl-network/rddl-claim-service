package config

import "sync"

const DefaultConfigTemplate = `
service-bind="{{ .ServiceBind }}"
service-port={{ .ServicePort }}
db-path="{{ .DBPath }}"
rpc-host="{{ .RPCHost }}"
rpc-user="{{ .RPCUser }}"
rpc-pass="{{ .RPCPass }}"
asset="{{ .Asset }}"
wallet="{{ .Wallet }}"
`

type Config struct {
	ServicePort int    `mapstructure:"service-port"`
	ServiceBind string `mapstructure:"service-bind"`
	DBPath      string `mapstructure:"db-path"`
	RPCHost     string `mapstructure:"rpc-host"`
	RPCUser     string `mapstructure:"rpc-user"`
	RPCPass     string `mapstructure:"rpc-pass"`
	Asset       string `mapstructure:"asset"`
	Wallet      string `mapstructure:"wallet"`
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
		RPCHost:     "planetmint-go-testnet-3.rddl.io:18884",
		RPCUser:     "user",
		RPCPass:     "password",
		Asset:       "7add40beb27df701e02ee85089c5bc0021bc813823fedb5f1dcb5debda7f3da9",
		Wallet:      "rddl2plmnt",
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}
