package service

import (
	"github.com/gin-gonic/gin"
)

type ClaimRequestBody struct {
	Beneficiary string `json:"beneficiary"`
	Amount      string `json:"amount"`
}

func (rcs *RDDLClaimService) postClaim(c *gin.Context) {
	var requestBody ClaimRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		return
	}

	// send tx to liquid

	// store claim
}
