package service_test

import (
	"fmt"
	"testing"

	"github.com/rddl-network/rddl-claim-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func TestBackend(t *testing.T) {
	memDB, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal("Error opening in-memory LevelDB: ", err)
	}
	defer memDB.Close()

	app := service.NewRDDLClaimService(memDB)

	id, err := app.PutUnconfirmedClaim()
	assert.NoError(t, err)

	valBytes, err := memDB.Get([]byte(fmt.Sprintf("claim:%d", id)), nil)
	assert.NoError(t, err)
	assert.Equal(t, "value:1", string(valBytes))

	id, err = app.PutUnconfirmedClaim()
	assert.NoError(t, err)

	err = app.ConfirmClaim(id)
	assert.NoError(t, err)

	valBytes, err = memDB.Get([]byte(fmt.Sprintf("claim:%d", id)), nil)
	assert.Equal(t, leveldb.ErrNotFound, err)
	assert.Equal(t, []uint8{}, valBytes)

	valBytes, err = memDB.Get([]byte(fmt.Sprintf("confirmedClaim:%d", id)), nil)
	assert.NoError(t, err)
	assert.Equal(t, "value:2", string(valBytes))
}
