package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

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
	var requestBody PostClaimRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rcs.logger.Info("msg", "received claim request", "beneficiary", requestBody.Beneficiary, "amount", requestBody.Amount)

	res, err := rcs.shamir.SendTokens(context.Background(), requestBody.Beneficiary, requestBody.Amount)
	if err != nil {
		rcs.logger.Error("msg", "failed to send tx", "beneficiary", requestBody.Beneficiary, "amount", requestBody.Amount)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send tx"})
		return
	}
	rcs.logger.Info("msg", "tokens sent", "TxID", res.TxID)

	rc := RedeemClaim{
		Beneficiary:  requestBody.Beneficiary,
		Amount:       requestBody.Amount,
		LiquidTXHash: res.TxID,
		ClaimID:      requestBody.ClaimID,
	}

	// store claim
	id, err := rcs.CreateUnconfirmedClaim(rc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store claim"})
		return
	}

	rcs.claims.mut.Lock()
	rcs.claims.list = append(rcs.claims.list, rc)
	rcs.claims.mut.Unlock()

	// TODO: return PostClaimResponse
	c.JSON(http.StatusOK, gin.H{"message": "claim enqueued", "id": id, "hash": res.TxID})
}
