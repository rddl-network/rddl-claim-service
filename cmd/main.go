package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func postClaim(c *gin.Context) {}

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

func startWebService(config *viper.Viper) {
	router := gin.Default()
	router.POST("/claim", postClaim)

	bindAddress := config.GetString("service-bind")
	servicePort := config.GetString("service-port")
	_ = router.Run(fmt.Sprintf("%s:%s", bindAddress, servicePort))
}

func main() {
	config, err := loadConfig("./")
	if err != nil {
		log.Fatalf("fatal error loading config file: %s", err)
	}

	startWebService(config)
}
