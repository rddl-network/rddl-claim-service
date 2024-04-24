package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rddl-network/rddl-claim-service/client"
	"github.com/rddl-network/rddl-claim-service/service"
	"github.com/stretchr/testify/assert"
)

func TestGetClaim(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/claim/1", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		res := service.GetClaimResponse{
			ID:           1,
			Beneficiary:  "beneficiary",
			Amount:       "1.00000000",
			LiquidTXHash: "liquidTXID",
			ClaimID:      1,
		}
		resBytes, err := json.Marshal(res)
		assert.NoError(t, err)
		_, err = w.Write(resBytes)
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	c := client.NewRCClient(mockServer.URL, mockServer.Client())
	res, err := c.GetClaim(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.ID)
	assert.Equal(t, "beneficiary", res.Beneficiary)
	assert.Equal(t, "1.00000000", res.Amount)
	assert.Equal(t, "liquidTXID", res.LiquidTXHash)
	assert.Equal(t, 1, res.ClaimID)
}

func TestPostClaim(t *testing.T) {
	t.Parallel()

	req := service.PostClaimRequest{
		ClaimID:     1,
		Beneficiary: "beneficiary",
		Amount:      "1.00000000",
	}

	expectedRes := service.PostClaimResponse{
		ID:   1,
		TxID: "liquidTXID",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/claim", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		bytes, err := json.Marshal(expectedRes)
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(bytes)
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	c := client.NewRCClient(mockServer.URL, mockServer.Client())
	res, err := c.PostClaim(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedRes.ID, res.ID)
	assert.Equal(t, expectedRes.TxID, res.TxID)
}
