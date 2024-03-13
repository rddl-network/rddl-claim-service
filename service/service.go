package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
	planetmint "github.com/planetmint/planetmint-go/lib"
	daotypes "github.com/planetmint/planetmint-go/x/dao/types"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/syndtr/goleveldb/leveldb"

	elements "github.com/rddl-network/elements-rpc"
)

type RDDLClaimService struct {
	db     *leveldb.DB
	router *gin.Engine
	claims SafeClaims
	shamir IShamirClient
}

type SafeClaims struct {
	mut  sync.Mutex
	list []RedeemClaim
}

func NewRDDLClaimService(db *leveldb.DB, router *gin.Engine, shamir IShamirClient) *RDDLClaimService {
	service := &RDDLClaimService{
		db:     db,
		router: router,
		claims: SafeClaims{
			list: make([]RedeemClaim, 0),
		},
		shamir: shamir,
	}
	service.registerRoutes()
	return service
}

func (rcs *RDDLClaimService) Load() (err error) {
	claims, err := rcs.GetAllUnconfirmedClaims()
	rcs.claims.mut.Lock()
	rcs.claims.list = claims
	rcs.claims.mut.Unlock()
	return
}

func (rcs *RDDLClaimService) Run(cfg *config.Config) error {
	go rcs.pollConfirmations(cfg.WaitPeriod, cfg.Confirmations)
	return rcs.router.Run(fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort))
}

func sendConfirmation(claimID int, beneficiary string) (err error) {
	cfg := config.GetConfig()
	addressString := cfg.PlanetmintAddress
	addr := sdk.MustAccAddressFromBech32(addressString)
	msg := daotypes.NewMsgConfirmRedeemClaim(addressString, uint64(claimID), beneficiary)

	_, err = planetmint.BroadcastTxWithFileLock(addr, msg)
	return
}

func (rcs *RDDLClaimService) pollConfirmations(waitPeriod int, confirmations int64) {
	ticker := time.NewTicker(time.Duration(waitPeriod) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rcs.claims.mut.Lock()
		for i, claim := range rcs.claims.list {
			txConfirmations, err := getTxConfirmations(claim.LiquidTXHash)
			if err != nil {
				log.Println("error while fetching tx confirmations: ", err)
			}
			if txConfirmations >= confirmations {
				rcs.claims.list = append(rcs.claims.list[:i], rcs.claims.list[i+1:]...)
				err := sendConfirmation(claim.ClaimID, claim.Beneficiary)
				if err != nil {
					log.Println("error while sending claim confirmation: ", err)
				}
				err = rcs.ConfirmClaim(claim.ID)
				if err != nil {
					log.Println("error while persisting claim confirmation: ", err)
				}
			}
		}
		rcs.claims.mut.Unlock()
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
