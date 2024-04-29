package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/rddl-network/shamir-coordinator-service/client"
	"github.com/syndtr/goleveldb/leveldb"

	stdlog "log"

	elements "github.com/rddl-network/elements-rpc"
	log "github.com/rddl-network/go-logger"
)

type RDDLClaimService struct {
	db       *leveldb.DB
	router   *gin.Engine
	shamir   client.IShamirCoordinatorClient
	pmClient IPlanetmintClient
	logger   log.AppLogger
}

func NewRDDLClaimService(db *leveldb.DB, router *gin.Engine, shamir client.IShamirCoordinatorClient, logger log.AppLogger, pmClient IPlanetmintClient) *RDDLClaimService {
	service := &RDDLClaimService{
		db:       db,
		router:   router,
		shamir:   shamir,
		logger:   logger,
		pmClient: pmClient,
	}
	service.registerRoutes()
	return service
}

func (rcs *RDDLClaimService) Run(cfg *config.Config) {
	go rcs.pollConfirmations(cfg.WaitPeriod, cfg.Confirmations)
	err := rcs.router.Run(fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort))
	if err != nil {
		stdlog.Panicf("error starting router: %s", err)
	}
}

func (rcs *RDDLClaimService) pollConfirmations(waitPeriod int, confirmations int64) {
	ticker := time.NewTicker(time.Duration(waitPeriod) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		claims, err := rcs.GetAllUnconfirmedClaims()
		if err != nil {
			rcs.logger.Error("msg", "error while reading unconfirmed claims: "+err.Error())
		}
		for _, claim := range claims {
			txConfirmations, err := getTxConfirmations(claim.LiquidTXHash)
			if err != nil {
				rcs.logger.Error("msg", "error while fetching tx confirmations: "+err.Error())
			}
			rcs.logger.Debug("msg", "fetchted liquid confirmations", "TxID", claim.LiquidTXHash, "confirmations", txConfirmations)
			if txConfirmations >= confirmations {
				rcs.logger.Info("msg", "sufficient confirmations reached, removing from polling list", "TxID", claim.LiquidTXHash)
				err = rcs.ConfirmClaim(claim.ID)
				if err != nil {
					rcs.logger.Error("msg", "error while persisting claim confirmation: "+err.Error())
				}
				txResponse, err := rcs.pmClient.SendConfirmation(claim.ClaimID, claim.Beneficiary)
				if err != nil {
					rcs.logger.Error("msg", "error while sending claim confirmation: "+err.Error())
				}
				rcs.logger.Info("msg", "claim confirmation sent to Planetmint", "txResponse", txResponse.String())
			}
		}
	}
}

func getTxConfirmations(txID string) (confirmations int64, err error) {
	cfg := config.GetConfig()

	url := fmt.Sprintf("http://%s:%s@%s/wallet/%s", cfg.RPCUser, cfg.RPCPass, cfg.RPCHost, cfg.Wallet)
	tx, err := elements.GetTransaction(url, []string{`"` + txID + `"`})
	if err != nil {
		return
	}

	return tx.Confirmations, err
}
