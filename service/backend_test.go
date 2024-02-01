package service_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	elements "github.com/rddl-network/elements-rpc"
	elementsmocks "github.com/rddl-network/elements-rpc/utils/mocks"
	"github.com/rddl-network/rddl-claim-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func createNRedeemClaim(app *service.RDDLClaimService, n int) []service.RedeemClaim {
	items := make([]service.RedeemClaim, n)
	for i := range items {
		items[i].Amount = "10000.00000000"
		items[i].Beneficiary = fmt.Sprintf("liquidAddress%d", i)
		items[i].LiquidTXHash = fmt.Sprintf("liquidTxHash%d", i)
		items[i].ClaimID = i
		id, _ := app.CreateUnconfirmedClaim(items[i])
		items[i].ID = id
	}
	return items
}

func setupService(t *testing.T) (app *service.RDDLClaimService, db *leveldb.DB, router *gin.Engine) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal("Error opening in-memory LevelDB: ", err)
	}
	router = gin.Default()

	elements.Client = &elementsmocks.MockClient{}

	app = service.NewRDDLClaimService(db, router)
	return
}

func TestGetUnconfirmedClaim(t *testing.T) {
	t.Parallel()
	app, db, _ := setupService(t)
	defer db.Close()

	items := createNRedeemClaim(app, 10)
	for _, item := range items {
		rc, err := app.GetUnconfirmedClaim(item.ID)
		assert.NoError(t, err)
		assert.Equal(t, item, rc)
	}
}

func TestGetAllUnconfirmedClaims(t *testing.T) {
	t.Parallel()
	app, db, _ := setupService(t)
	defer db.Close()

	items := createNRedeemClaim(app, 20)
	claims, err := app.GetAllUnconfirmedClaims()
	assert.NoError(t, err)
	assert.Equal(t, items, claims)
}

func TestDeleteUnconfirmedClaim(t *testing.T) {
	t.Parallel()
	app, db, _ := setupService(t)
	defer db.Close()

	items := createNRedeemClaim(app, 1)
	err := app.DeleteUnconfirmedClaim(items[0].ID)
	assert.NoError(t, err)

	_, err = app.GetUnconfirmedClaim(items[0].ID)
	assert.Error(t, err)
	assert.Equal(t, leveldb.ErrNotFound, err)
}

func TestConfirmClaim(t *testing.T) {
	t.Parallel()
	app, db, _ := setupService(t)
	defer db.Close()

	items := createNRedeemClaim(app, 1)
	err := app.ConfirmClaim(items[0].ID)
	assert.NoError(t, err)

	_, err = app.GetUnconfirmedClaim(items[0].ID)
	assert.Error(t, err)
	assert.Equal(t, leveldb.ErrNotFound, err)

	rc, err := app.GetConfirmedClaim(items[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, items[0], rc)
}

func TestGetAllConfirmedClaims(t *testing.T) {
	t.Parallel()
	app, db, _ := setupService(t)
	defer db.Close()

	items := createNRedeemClaim(app, 10)
	for _, item := range items {
		err := app.ConfirmClaim(item.ID)
		assert.NoError(t, err)
	}

	claims, err := app.GetAllConfirmedClaims()
	assert.NoError(t, err)
	assert.Equal(t, items, claims)
}
