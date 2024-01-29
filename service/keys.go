package service

import (
	"encoding/binary"
)

const (
	ClaimKeyPrefix          = "Claim/"
	ConfirmedClaimKeyPrefix = "ConfirmedClaim/"
	CountKey                = "Count"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func ClaimKey(id int) []byte {
	var key []byte

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))

	prefixBytes := []byte(ClaimKeyPrefix)
	key = append(key, prefixBytes...)
	key = append(key, buf...)

	return key
}

func ConfirmedClaimKey(id int) []byte {
	var key []byte

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))

	prefixBytes := []byte(ConfirmedClaimKeyPrefix)
	key = append(key, prefixBytes...)
	key = append(key, buf...)
	return key
}
