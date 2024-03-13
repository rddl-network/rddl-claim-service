package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

type ClaimRequestBody struct {
	Beneficiary string `binding:"required" json:"beneficiary"`
	Amount      string `binding:"required" json:"amount"`
	ClaimID     int    `binding:"required" json:"claim-id"`
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
	var requestBody ClaimRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := rcs.shamir.IssueTransaction(requestBody.Amount, requestBody.Beneficiary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send tx"})
		return
	}

	rc := RedeemClaim{
		Beneficiary:  requestBody.Beneficiary,
		Amount:       requestBody.Amount,
		LiquidTXHash: hash,
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

	c.JSON(http.StatusOK, gin.H{"message": "claim enqueued", "id": id, "hash": hash})
}
