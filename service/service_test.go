package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	elements "github.com/rddl-network/elements-rpc"
	elementsmocks "github.com/rddl-network/elements-rpc/utils/mocks"
	log "github.com/rddl-network/go-logger"
	"github.com/rddl-network/rddl-claim-service/config"
	"github.com/rddl-network/rddl-claim-service/service"
	"github.com/rddl-network/rddl-claim-service/testutil"
	"github.com/rddl-network/rddl-claim-service/types"
	shamir "github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var (
	pmMock     *testutil.MockIPlanetmintClient
	shamirMock *testutil.MockIShamirCoordinatorClient
)

func setupServiceWithMocks(t *testing.T) (app *service.RDDLClaimService, db *leveldb.DB, router *gin.Engine) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal("Error opening in-memory LevelDB: ", err)
	}
	router = gin.Default()
	ctrl := gomock.NewController(t)
	shamirMock = testutil.NewMockIShamirCoordinatorClient(ctrl)
	pmMock = testutil.NewMockIPlanetmintClient(ctrl)
	elements.Client = &elementsmocks.MockClient{}
	logger := log.GetLogger(log.DEBUG)
	app = service.NewRDDLClaimService(db, router, shamirMock, logger, pmMock)
	return
}

func TestIntegration(t *testing.T) {
	cfg := config.GetConfig()
	cfg.WaitPeriod = 1
	cfg.Confirmations = 0

	app, _, router := setupServiceWithMocks(t)

	// define Request
	reqBody := types.PostClaimRequest{
		Beneficiary: "address",
		Amount:      1000000000000,
		ClaimID:     1,
	}

	// ShamirCoordinator.SendTokens should be called exactly once with reqBody.Amount{1000000000000} converted to "10000.00000000" liquid float
	mockRes := shamir.SendTokensResponse{TxID: "0000000000000000000000000000000000000000000000000000000000000000"}
	shamirMock.EXPECT().SendTokens(gomock.Any(), gomock.Any(), "10000.00000000").Times(1).Return(mockRes, nil)

	// Planetmint Confirmations shall be sent exactly once with reqBody.ClaimID{1} and reqBody.Beneficiary{"address"}
	pmMock.EXPECT().SendConfirmation(reqBody.ClaimID, reqBody.Beneficiary).Times(1).Return(sdk.TxResponse{
		TxHash: "0000000000000000000000000000000000000000000000000000000000000000",
	}, nil)

	// Start service seperate thread
	go app.Run(cfg)

	// Send PostClaimRequest to service
	bodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/claim", bytes.NewBuffer(bodyBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"id\":1,\"tx-id\":\"0000000000000000000000000000000000000000000000000000000000000000\"}", w.Body.String())

	var res types.PostClaimResponse
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)

	// Make sure unconfirmed claim is written to disk
	rc, err := app.GetUnconfirmedClaim(res.ID)
	assert.NoError(t, err)

	assert.Equal(t, reqBody.Amount, rc.Amount)
	assert.Equal(t, reqBody.Beneficiary, rc.Beneficiary)
	assert.Equal(t, res.TxID, rc.LiquidTXHash)

	// wait for polling
	time.Sleep(2 * time.Second)

	// Make sure unconfirmed claim is deleted
	rc, err = app.GetUnconfirmedClaim(res.ID)
	assert.Error(t, err)
	assert.Equal(t, leveldb.ErrNotFound, err)

	// Make sure deleted unconfirmed claim is stored as confirmed claim
	rc, err = app.GetConfirmedClaim(res.ID)
	assert.NoError(t, err)
}
