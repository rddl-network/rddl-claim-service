package service_test

import (
	"encoding/json"
	"testing"

	"github.com/rddl-network/rddl-claim-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// TODO: implement test init and test all functions seperately
func TestBackend(t *testing.T) {
	memDB, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal("Error opening in-memory LevelDB: ", err)
	}
	defer memDB.Close()

	app := service.NewRDDLClaimService(memDB)

	rc := service.RedeemClaim{
		Beneficiary:  "beneficiary",
		Amount:       "1000.00000000",
		LiquidTXHash: "liquidTxHash",
	}

	id, err := app.CreateUnconfirmedClaim(rc)
	assert.NoError(t, err)

	rc.Id = id
	rcBytes, err := json.Marshal(rc)
	assert.NoError(t, err)

	valBytes, err := memDB.Get(service.ClaimKey(id), nil)
	assert.NoError(t, err)
	assert.Equal(t, string(rcBytes), string(valBytes))

	id, err = app.CreateUnconfirmedClaim(rc)
	assert.NoError(t, err)

	err = app.ConfirmClaim(id)
	assert.NoError(t, err)

	valBytes, err = memDB.Get(service.ClaimKey(id), nil)
	assert.Equal(t, leveldb.ErrNotFound, err)
	assert.Equal(t, []uint8{}, valBytes)

	rc.Id = id
	rcBytes, err = json.Marshal(rc)
	assert.NoError(t, err)
	valBytes, err = memDB.Get([]byte(service.ConfirmedClaimKey(id)), nil)
	assert.NoError(t, err)
	assert.Equal(t, string(rcBytes), string(valBytes))

	app.CreateUnconfirmedClaim(rc)
	app.CreateUnconfirmedClaim(rc)
	app.CreateUnconfirmedClaim(rc)
	claims := make([]string, 0)

	iter := memDB.NewIterator(util.BytesPrefix([]byte(service.ClaimKeyPrefix)), nil)
	for iter.Next() {
		valBytes := iter.Value()
		val := string(valBytes)
		claims = append(claims, val)
	}
	assert.Len(t, claims, 4)
	iter.Release()

	confirmedClaims := make([]string, 0)
	iter = memDB.NewIterator(util.BytesPrefix([]byte(service.ConfirmedClaimKeyPrefix)), nil)
	for iter.Next() {
		valBytes := iter.Value()
		val := string(valBytes)
		confirmedClaims = append(confirmedClaims, val)
	}
	assert.Len(t, confirmedClaims, 1)
	iter.Release()

	iter = memDB.NewIterator(util.BytesPrefix([]byte(service.ClaimKeyPrefix)), nil)
	iter.Seek([]byte{2})
	for iter.Next() {
		valBytes := iter.Value()
		val := string(valBytes)
		claims = append(claims, val)
	}
	assert.Len(t, claims, 7)
	iter.Release()
}
