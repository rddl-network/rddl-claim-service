package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/planetmint/planetmint-go/util"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/rddl-network/rddl-claim-service/types"
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

	var resBody types.GetClaimResponse
	resBody.ID = rc.ID
	resBody.Beneficiary = rc.Beneficiary
	resBody.LiquidTXHash = rc.LiquidTXHash
	resBody.Amount = rc.Amount
	resBody.ClaimID = rc.ClaimID

	c.JSON(http.StatusOK, resBody)
}

func (rcs *RDDLClaimService) postClaim(c *gin.Context) {
	cfg := config.GetConfig()

	// Read and buffer the body for multiple uses
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var requestBody types.PostClaimRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// manual validation if claim-id is set on requestBody
	if err := checkClaimIDSet(bodyBytes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rcs.logger.Info("msg", "received claim request", "beneficiary", requestBody.Beneficiary, "amount", requestBody.Amount)

	res, err := rcs.shamir.SendTokens(context.Background(), requestBody.Beneficiary, util.UintValueToRDDLTokenString(requestBody.Amount), cfg.Asset)
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

	var resBody types.PostClaimResponse
	resBody.ID = id
	resBody.TxID = res.TxID

	c.JSON(http.StatusOK, resBody)
}

func checkClaimIDSet(bodyBytes []byte) (err error) {
	if !strings.Contains(string(bodyBytes), "claim-id") {
		return errors.New("missing claim-id")
	}
	return
}
