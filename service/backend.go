package service

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

func InitDB(config *viper.Viper) (db *leveldb.DB, err error) {
	dbPath := config.GetString("db-path")
	return leveldb.OpenFile(dbPath, nil)
}

func (rcs *RDDLClaimService) incrementCount() (count int, err error) {
	countBytes, err := rcs.db.Get([]byte("count"), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return 0, err
	}

	if countBytes == nil {
		count = 1
	} else {
		count, err = strconv.Atoi(string(countBytes))
		if err != nil {
			return 0, err
		}
		count = count + 1
	}

	rcs.db.Put([]byte("count"), []byte(strconv.Itoa(count)), nil)

	return count, nil
}

func (rcs *RDDLClaimService) GetUnconfirmedClaim(id int) (claim []byte, err error) {
	key := []byte(fmt.Sprintf("claim:%d", id))
	return rcs.db.Get(key, nil)
}

func (rcs *RDDLClaimService) PutUnconfirmedClaim() (id int, err error) {
	id, err = rcs.incrementCount()
	if err != nil {
		return
	}

	key := []byte(fmt.Sprintf("claim:%d", id))

	if err := rcs.db.Put(key, []byte(fmt.Sprintf("value:%d", id)), nil); err != nil {
		return 0, err
	}

	return id, nil
}

func (rcs *RDDLClaimService) DeleteUnconfirmClaim(id int) (err error) {
	key := []byte(fmt.Sprintf("claim:%d", id))
	return rcs.db.Delete(key, nil)
}

func (rcs *RDDLClaimService) ConfirmClaim(id int) (err error) {
	val, err := rcs.GetUnconfirmedClaim(id)
	if err != nil {
		return
	}

	err = rcs.DeleteUnconfirmClaim(id)
	if err != nil {
		return
	}

	key := []byte(fmt.Sprintf("confirmedClaim:%d", id))

	return rcs.db.Put(key, val, nil)
}
