package service

type RedeemClaim struct {
	ID           int    `json:"id"`
	Beneficiary  string `json:"beneficiary"`
	Amount       string `json:"amount"`
	LiquidTXHash string `json:"liquid-tx-hash"`
	ClaimID      int    `json:"claim-id"`
}
