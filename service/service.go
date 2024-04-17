package service

import (
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
	planetmint "github.com/planetmint/planetmint-go/lib"
	daotypes "github.com/planetmint/planetmint-go/x/dao/types"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/rddl-network/shamir-coordinator-service/client"
	"github.com/syndtr/goleveldb/leveldb"

	elements "github.com/rddl-network/elements-rpc"
	log "github.com/rddl-network/go-logger"
)

type RDDLClaimService struct {
	db     *leveldb.DB
	router *gin.Engine
	claims SafeClaims
	shamir client.IShamirCoordinatorClient
	logger log.AppLogger
}

type SafeClaims struct {
	mut  sync.Mutex
	list []RedeemClaim
}

func NewRDDLClaimService(db *leveldb.DB, router *gin.Engine, shamir client.IShamirCoordinatorClient, logger log.AppLogger) *RDDLClaimService {
	service := &RDDLClaimService{
		db:     db,
		router: router,
		claims: SafeClaims{
			list: make([]RedeemClaim, 0),
		},
		shamir: shamir,
		logger: logger,
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

func sendConfirmation(claimID int, beneficiary string) (txResponse sdk.TxResponse, err error) {
	cfg := config.GetConfig()
	addressString := cfg.PlanetmintAddress
	addr := sdk.MustAccAddressFromBech32(addressString)
	msg := daotypes.NewMsgConfirmRedeemClaim(addressString, uint64(claimID), beneficiary)

	out, err := planetmint.BroadcastTxWithFileLock(addr, msg)
	if err != nil {
		return
	}

	return planetmint.GetTxResponseFromOut(out)
}

func (rcs *RDDLClaimService) pollConfirmations(waitPeriod int, confirmations int64) {
	ticker := time.NewTicker(time.Duration(waitPeriod) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rcs.claims.mut.Lock()
		for i, claim := range rcs.claims.list {
			txConfirmations, err := getTxConfirmations(claim.LiquidTXHash)
			if err != nil {
				rcs.logger.Error("msg", "error while fetching tx confirmations: "+err.Error())
			}
			rcs.logger.Debug("msg", "fetchted liquid confirmations", "TxID", claim.LiquidTXHash, "confirmations", txConfirmations)
			if txConfirmations >= confirmations {
				rcs.logger.Info("msg", "sufficient confirmations reached, removing from polling list", "TxID", claim.LiquidTXHash)
				rcs.claims.list = append(rcs.claims.list[:i], rcs.claims.list[i+1:]...)
				txResponse, err := sendConfirmation(claim.ClaimID, claim.Beneficiary)
				if err != nil {
					rcs.logger.Error("msg", "error while sending claim confirmation: "+err.Error())
				}
				err = rcs.ConfirmClaim(claim.ID)
				if err != nil {
					rcs.logger.Error("msg", "error while persisting claim confirmation: "+err.Error())
				}
				rcs.logger.Info("msg", "claim confirmation sent to Planetmint", "txResponse", txResponse.String())
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
