package service

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"

	elements "github.com/rddl-network/elements-rpc"
)

type RDDLClaimService struct {
	db          *leveldb.DB
	router      *gin.Engine
	queue       map[string]RedeemClaim
	dataChannel chan RedeemClaim
}

func NewRDDLClaimService(db *leveldb.DB) *RDDLClaimService {
	service := &RDDLClaimService{
		db:          db,
		router:      gin.Default(),
		queue:       make(map[string]RedeemClaim),
		dataChannel: make(chan RedeemClaim),
	}
	service.registerRoutes()
	return service
}

func (rcs *RDDLClaimService) Load() (err error) {
	claims, err := rcs.GetAllUnconfirmedClaims()
	for _, claim := range claims {
		rcs.queue[claim.LiquidTXHash] = claim
		rcs.dataChannel <- claim
	}
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

	go pollConfirmations(rcs.dataChannel)
	for {
		claim := <-rcs.dataChannel
		err := rcs.ConfirmClaim(claim.Id)
		if err != nil {
			log.Println("error while persisting claim confirmation: ", err)
			continue
		}
		delete(rcs.queue, claim.LiquidTXHash)
	}
}

func pollConfirmations(c chan RedeemClaim) {
	cfg := config.GetConfig()
	claims := make([]RedeemClaim, 0)
	for {
		claim := <-c
		claims = append(claims, claim)

		for _, rc := range claims {
			confirmations, err := getTxConfirmations(rc.LiquidTXHash)
			if err != nil {
				log.Println("error while fetching tx confirmations: ", err)
			}
			if confirmations >= cfg.Confirmations {
				c <- rc
			}
		}

		time.Sleep(time.Duration(cfg.WaitPeriod) * time.Second)
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
