package service

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

type RDDLClaimService struct {
	db     *leveldb.DB
	router *gin.Engine
}

func NewRDDLClaimService(db *leveldb.DB) *RDDLClaimService {
	service := &RDDLClaimService{
		db:     db,
		router: gin.Default(),
	}
	service.registerRoutes()
	return service
}

func (rcs *RDDLClaimService) Run(config *viper.Viper) {
	bindAddress := config.GetString("service-bind")
	servicePort := config.GetString("service-port")
	err := rcs.router.Run(fmt.Sprintf("%s:%s", bindAddress, servicePort))
	if err != nil {
		log.Fatalf("fatal error starting router: %s", err)
		panic(err)
	}
}

func (rcs *RDDLClaimService) registerRoutes() {
	rcs.router.POST("/claim", rcs.postClaim)
}
