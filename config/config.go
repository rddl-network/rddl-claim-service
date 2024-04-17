package config

import (
	"sync"

	log "github.com/rddl-network/go-logger"
)

const DefaultConfigTemplate = `
service-host="{{ .ServiceHost }}"
service-port={{ .ServicePort }}
db-path="{{ .DBPath }}"
rpc-host="{{ .RPCHost }}"
rpc-user="{{ .RPCUser }}"
rpc-pass="{{ .RPCPass }}"
asset="{{ .Asset }}"
wallet="{{ .Wallet }}"
confirmations={{ .Confirmations }}
wait-period={{ .WaitPeriod }}
planetmint-address="{{ .PlanetmintAddress }}"
shamir-host="{{ .ShamirHost }}"
log-level="{{ .LogLevel }}"
`

type Config struct {
	ServicePort       int    `mapstructure:"service-port"`
	ServiceHost       string `mapstructure:"service-host"`
	DBPath            string `mapstructure:"db-path"`
	RPCHost           string `mapstructure:"rpc-host"`
	RPCUser           string `mapstructure:"rpc-user"`
	RPCPass           string `mapstructure:"rpc-pass"`
	Asset             string `mapstructure:"asset"`
	Wallet            string `mapstructure:"wallet"`
	Confirmations     int64  `mapstructure:"confirmations"`
	WaitPeriod        int    `mapstructure:"wait-period"`
	PlanetmintAddress string `mapstructure:"planetmint-address"`
	ShamirHost        string `mapstructure:"shamir-host"`
	LogLevel          string `mapstructure:"log-level"`
}

// global singleton
var (
	config     *Config
	initConfig sync.Once
)

// DefaultConfig returns RDDL-2-PLMNT default config
func DefaultConfig() *Config {
	return &Config{
		ServicePort:       8080,
		ServiceHost:       "localhost",
		DBPath:            "./data",
		RPCHost:           "planetmint-go-testnet-3.rddl.io:18884",
		RPCUser:           "user",
		RPCPass:           "password",
		Asset:             "7add40beb27df701e02ee85089c5bc0021bc813823fedb5f1dcb5debda7f3da9",
		Wallet:            "pop",
		Confirmations:     10,
		WaitPeriod:        10,
		PlanetmintAddress: "plmnt15xuq0yfxtd70l7jzr5hg722sxzcqqdcr8ptpl5",
		ShamirHost:        "localhost:9091",
		LogLevel:          log.DEBUG,
	}
}

func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}
