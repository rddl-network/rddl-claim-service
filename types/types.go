package types

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
