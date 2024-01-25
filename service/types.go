package service

type RedeemClaim struct {
	Id           int    `json:"id"`
	Beneficiary  string `json:"beneficiary"`
	Amount       string `json:"amount"`
	LiquidTXHash string `json:"liquidTxHash"`
}
