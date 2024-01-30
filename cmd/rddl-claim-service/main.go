package main

import (
	"log"

	"github.com/rddl-network/rddl-claim-service/service"
	"github.com/spf13/viper"
)

func loadConfig(path string) (v *viper.Viper, err error) {
	v = viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("toml")

	err = v.ReadInConfig()
	if err == nil {
		return
	}

	return
}

func main() {
	config, err := loadConfig("./")
	if err != nil {
		log.Fatalf("fatal error loading config file: %s", err)
	}

	db, err := service.InitDB(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service := service.NewRDDLClaimService(db)

	service.Run(config)
	err = service.Load()
	if err != nil {
		log.Fatalf("fatal error loading claims: %s", err)
	}
}
