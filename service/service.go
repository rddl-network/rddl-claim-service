package service

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"

	elements "github.com/rddl-network/elements-rpc"
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

func (rcs *RDDLClaimService) Load() (claims []RedeemClaim, err error) {
	claims, err = rcs.GetAllUnconfirmedClaims()

	return
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

func getTxConfirmations(txID string) (confirmations int64, err error) {
	cfg := config.GetConfig()

	url := fmt.Sprintf("http://%s:%s@%s/wallet/%s", cfg.RPCUser, cfg.RPCPass, cfg.RPCHost, cfg.Wallet)
	tx, err := elements.GetTransaction(url, []string{txID})
	if err != nil {
		return
	}

	return tx.Confirmations, err
}
