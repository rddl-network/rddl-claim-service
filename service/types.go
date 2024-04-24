package service

type RedeemClaim struct {
	ID           int    `json:"id"`
	Beneficiary  string `json:"beneficiary"`
	Amount       string `json:"amount"`
	LiquidTXHash string `json:"liquid-tx-hash"`
	ClaimID      int    `json:"claim-id"`
}

type PostClaimRequest struct {
	Beneficiary string `binding:"required" json:"beneficiary"`
	Amount      string `binding:"required" json:"amount"`
	ClaimID     int    `binding:"required" json:"claim-id"`
}

type PostClaimResponse struct {
	ID   string `binding:"required" json:"id"`
	TxID string `binding:"required" json:"tx-id"`
}

type GetClaimResponse struct {
	ID           int    `json:"id"`
	Beneficiary  string `json:"beneficiary"`
	Amount       string `json:"amount"`
	LiquidTXHash string `json:"liquid-tx-hash"`
	ClaimID      int    `json:"claim-id"`
}
