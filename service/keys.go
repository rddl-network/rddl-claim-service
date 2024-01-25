package service

import "fmt"

const (
	ClaimKeyPrefix          = "Claim/"
	ConfirmedClaimKeyPrefix = "ConfirmedClaim/"
	CountKey                = "Count"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func ClaimKey(id int) []byte {
	return []byte(fmt.Sprintf("%s%d", ClaimKeyPrefix, id))
}

func ConfirmedClaimKey(id int) []byte {
	return []byte(fmt.Sprintf("%s%d", ConfirmedClaimKeyPrefix, id))
}
