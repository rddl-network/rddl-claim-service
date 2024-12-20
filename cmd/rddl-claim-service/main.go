package main

import (
	"bytes"
	"html/template"
	stdlog "log"
	"os"

	log "github.com/rddl-network/go-utils/logger"
	"github.com/rddl-network/shamir-coordinator-service/client"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/go-utils/tls"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/rddl-network/rddl-claim-service/service"
	"github.com/spf13/viper"

	"github.com/planetmint/planetmint-go/app"
	"github.com/planetmint/planetmint-go/lib"
)

var libConfig *lib.Config

func loadConfig(path string) (cfg *config.Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("toml")

	err = v.ReadInConfig()
	if err == nil {
		cfg = config.GetConfig()
		cfg.ServicePort = v.GetInt("service-port")
		cfg.ServiceHost = v.GetString("service-host")
		cfg.CertsPath = v.GetString("certs-path")
		cfg.DBPath = v.GetString("db-path")
		cfg.RPCHost = v.GetString("rpc-host")
		cfg.RPCUser = v.GetString("rpc-user")
		cfg.RPCPass = v.GetString("rpc-pass")
		cfg.Asset = v.GetString("asset")
		cfg.Wallet = v.GetString("wallet")
		cfg.Confirmations = v.GetInt64("confirmations")
		cfg.WaitPeriod = v.GetInt("wait-period")
		cfg.PlanetmintAddress = v.GetString("planetmint-address")
		cfg.PlanetmintChainID = v.GetString("planetmint-chain-id")
		cfg.ShamirHost = v.GetString("shamir-host")
		cfg.LogLevel = v.GetString("log-level")
		return
	}
	stdlog.Println("no config file found")

	tmpl := template.New("appConfigFileTemplate")
	configTemplate, err := tmpl.Parse(config.DefaultConfigTemplate)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	if err = configTemplate.Execute(&buffer, config.GetConfig()); err != nil {
		return
	}

	if err = v.ReadConfig(&buffer); err != nil {
		return
	}
	if err = v.SafeWriteConfig(); err != nil {
		return
	}

	stdlog.Println("default config file created. please adapt it and restart the application. exiting...")
	os.Exit(0)
	return
}

func main() {
	config, err := loadConfig("./")
	if err != nil {
		stdlog.Fatalf("fatal error loading config file: %s", err)
	}

	encodingConfig := app.MakeEncodingConfig()

	libConfig = lib.GetConfig()
	libConfig.SetEncodingConfig(encodingConfig)

	planetmintChainID := config.PlanetmintChainID
	if planetmintChainID == "" {
		stdlog.Fatalf("chain id must not be empty")
	}
	libConfig.SetChainID(planetmintChainID)

	db, err := service.InitDB(config)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	logger := log.GetLogger(config.LogLevel)
	mTLSClient, err := tls.Get2WayTLSClient(config.CertsPath)
	if err != nil {
		stdlog.Fatalf("fatal error setting up mutual TLS shareholder client")
	}
	shamir := client.NewSCClient(config.ShamirHost, mTLSClient)
	pmClient := service.NewPlanetmintClient()
	service := service.NewRDDLClaimService(db, router, shamir, logger, pmClient)

	service.Run(config)
}
