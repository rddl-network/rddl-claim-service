package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	elements "github.com/rddl-network/elements-rpc"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/syndtr/goleveldb/leveldb"
)

type ClaimRequestBody struct {
	Beneficiary string `binding:"required" json:"beneficiary"`
	Amount      string `binding:"required" json:"amount"`
}

func (rcs *RDDLClaimService) registerRoutes() {
	rcs.router.POST("/claim", rcs.postClaim)
	rcs.router.GET("/claim/:id", rcs.getClaim)
}

func (rcs *RDDLClaimService) getClaim(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	rc, err := rcs.GetUnconfirmedClaim(id)
	if errors.Is(err, leveldb.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no claim found for id %d", id)})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rc)
}

func (rcs *RDDLClaimService) postClaim(c *gin.Context) {
	cfg := config.GetConfig()

	var requestBody ClaimRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// send tx to liquid
	url := fmt.Sprintf("http://%s:%s@%s/wallet/%s", cfg.RPCUser, cfg.RPCPass, cfg.RPCHost, cfg.Wallet)
	hex, err := elements.SendToAddress(url, []string{
		requestBody.Beneficiary,
		requestBody.Amount,
		`""`,
		`""`,
		"false",
		"true",
		"null",
		`"unset"`,
		"false",
		`"` + cfg.Asset + `"`,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send tx"})
		return
	}

	rc := RedeemClaim{
		Beneficiary:  requestBody.Beneficiary,
		Amount:       requestBody.Amount,
		LiquidTXHash: hex,
	}

	// store claim
	id, err := rcs.CreateUnconfirmedClaim(rc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store claim"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "claim enqueued", "id": id, "hash": hex})
}
