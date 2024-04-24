package service

type RedeemClaim struct {
	ID           int    `json:"id"`
	Beneficiary  string `json:"beneficiary"`
	Amount       uint64 `json:"amount"`
	LiquidTXHash string `json:"liquid-tx-hash"`
	ClaimID      int    `json:"claim-id"`
}

type PostClaimRequest struct {
	Beneficiary string `binding:"required" json:"beneficiary"`
	Amount      uint64 `binding:"required" json:"amount"`
	ClaimID     int    `binding:"required" json:"claim-id"`
}

type PostClaimResponse struct {
	ID   int    `binding:"required" json:"id"`
	TxID string `binding:"required" json:"tx-id"`
}

type GetClaimResponse struct {
	ID           int    `binding:"required" json:"id"`
	Beneficiary  string `binding:"required" json:"beneficiary"`
	Amount       uint64 `binding:"required" json:"amount"`
	LiquidTXHash string `binding:"required" json:"liquid-tx-hash"`
	ClaimID      int    `binding:"required" json:"claim-id"`
}
