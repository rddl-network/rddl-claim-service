package service

import (
	"fmt"
	"log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
	planetmint "github.com/planetmint/planetmint-go/lib"
	daotypes "github.com/planetmint/planetmint-go/x/dao/types"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"

	elements "github.com/rddl-network/elements-rpc"
)

type RDDLClaimService struct {
	db          *leveldb.DB
	router      *gin.Engine
	queue       map[string]RedeemClaim
	receiveChan chan RedeemClaim
	confirmChan chan RedeemClaim
}

func NewRDDLClaimService(db *leveldb.DB, router *gin.Engine) *RDDLClaimService {
	service := &RDDLClaimService{
		db:          db,
		router:      router,
		queue:       make(map[string]RedeemClaim),
		receiveChan: make(chan RedeemClaim),
		confirmChan: make(chan RedeemClaim),
	}
	service.registerRoutes()
	return service
}

func (rcs *RDDLClaimService) Load() (err error) {
	claims, err := rcs.GetAllUnconfirmedClaims()
	for _, claim := range claims {
		rcs.queue[claim.LiquidTXHash] = claim
		rcs.receiveChan <- claim
	}
	return
}

func (rcs *RDDLClaimService) Run(config *viper.Viper) {
	bindAddress := config.GetString("service-bind")
	servicePort := config.GetString("service-port")
	err := rcs.router.Run(fmt.Sprintf("%s:%s", bindAddress, servicePort))
	if err != nil {
		log.Fatalf("fatal error starting router: %s", err)
	}

	go pollConfirmations(rcs.receiveChan, rcs.confirmChan)
	for {
		claim := <-rcs.confirmChan
		err := rcs.ConfirmClaim(claim.ID)
		if err != nil {
			log.Println("error while persisting claim confirmation: ", err)
			continue
		}
		delete(rcs.queue, claim.LiquidTXHash)
	}
}

func sendConfirmation(cfg *viper.Viper, beneficiary string) (err error) {
	addressString := cfg.GetString("planetmint-address")
	addr := sdk.MustAccAddressFromBech32(addressString)
	msg := daotypes.NewMsgConfirmRedeemClaim(addressString, 0, beneficiary)

	_, err = planetmint.BroadcastTxWithFileLock(addr, msg)
	return
}

func pollConfirmations(receive chan RedeemClaim, confirm chan RedeemClaim) {
	cfg := config.GetConfig()
	claims := make(map[string]RedeemClaim)
	for {
		claim := <-receive
		claims[claim.LiquidTXHash] = claim

		for _, rc := range claims {
			confirmations, err := getTxConfirmations(rc.LiquidTXHash)
			if err != nil {
				log.Println("error while fetching tx confirmations: ", err)
			}
			if confirmations >= cfg.Confirmations {
				confirm <- rc
				delete(claims, rc.LiquidTXHash)
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
