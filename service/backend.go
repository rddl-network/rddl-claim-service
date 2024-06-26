package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var dbMutex sync.Mutex

type RedeemClaim struct {
	ID           int    `json:"id"`
	Beneficiary  string `json:"beneficiary"`
	Amount       uint64 `json:"amount"`
	LiquidTXHash string `json:"liquid-tx-hash"`
	ClaimID      int    `json:"claim-id"`
}

func InitDB(cfg *config.Config) (db *leveldb.DB, err error) {
	return leveldb.OpenFile(cfg.DBPath, nil)
}

func (rcs *RDDLClaimService) incrementCount() (count int, err error) {
	countBytes, err := rcs.db.Get(KeyPrefix(CountKey), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return 0, err
	}

	if countBytes == nil {
		count = 1
	} else {
		count, err = strconv.Atoi(string(countBytes))
		if err != nil {
			return 0, err
		}
		count++
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()
	err = rcs.db.Put(KeyPrefix(CountKey), []byte(strconv.Itoa(count)), nil)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (rcs *RDDLClaimService) GetUnconfirmedClaim(id int) (claim RedeemClaim, err error) {
	key := ClaimKey(id)
	valBytes, err := rcs.db.Get(key, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(valBytes, &claim)
	return
}

func (rcs *RDDLClaimService) GetAllUnconfirmedClaims() (claims []RedeemClaim, err error) {
	iter := rcs.db.NewIterator(util.BytesPrefix([]byte(ClaimKeyPrefix)), nil)
	defer iter.Release()
	for iter.Next() {
		var claim RedeemClaim
		claimBytes := iter.Value()
		err = json.Unmarshal(claimBytes, &claim)
		if err != nil {
			return
		}
		claims = append(claims, claim)
	}
	return
}

func (rcs *RDDLClaimService) CreateUnconfirmedClaim(rc RedeemClaim) (id int, err error) {
	id, err = rcs.incrementCount()
	if err != nil {
		return
	}

	rc.ID = id

	key := ClaimKey(id)
	val, err := json.Marshal(rc)
	if err != nil {
		return 0, err
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()
	if err := rcs.db.Put(key, val, nil); err != nil {
		return 0, err
	}

	return id, nil
}

func (rcs *RDDLClaimService) DeleteUnconfirmedClaim(id int) (err error) {
	key := ClaimKey(id)
	dbMutex.Lock()
	defer dbMutex.Unlock()
	return rcs.db.Delete(key, nil)
}

func (rcs *RDDLClaimService) ConfirmClaim(id int) (err error) {
	claim, err := rcs.GetUnconfirmedClaim(id)
	if err != nil {
		return
	}

	err = rcs.DeleteUnconfirmedClaim(id)
	if err != nil {
		return
	}

	key := ConfirmedClaimKey(id)
	val, err := json.Marshal(claim)
	if err != nil {
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()
	return rcs.db.Put(key, val, nil)
}

func (rcs *RDDLClaimService) GetConfirmedClaim(id int) (claim RedeemClaim, err error) {
	key := ConfirmedClaimKey(id)
	valBytes, err := rcs.db.Get(key, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(valBytes, &claim)
	return
}

func (rcs *RDDLClaimService) GetAllConfirmedClaims() (claims []RedeemClaim, err error) {
	iter := rcs.db.NewIterator(util.BytesPrefix([]byte(ConfirmedClaimKeyPrefix)), nil)
	defer iter.Release()
	for iter.Next() {
		var claim RedeemClaim
		claimBytes := iter.Value()
		err = json.Unmarshal(claimBytes, &claim)
		if err != nil {
			return
		}
		claims = append(claims, claim)
	}
	return
}
