package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rddl-network/rddl-claim-service/types"
	"github.com/stretchr/testify/assert"
)

func TestGetClaimRoute(t *testing.T) {
	app, _, router := setupService(t)

	items := createNRedeemClaim(app, 1)
	itemBytes, err := json.Marshal(items[0])
	assert.NoError(t, err)

	tests := []struct {
		name   string
		id     string
		code   int
		err    bool
		errMsg string
	}{
		{
			name: "valid request",
			id:   "1",
			code: 200,
			err:  false,
		},
		{
			name:   "not found",
			id:     "2",
			code:   404,
			err:    true,
			errMsg: "{\"error\":\"no claim found for id 2\"}",
		},
		{
			name:   "invalid request",
			id:     "foo",
			code:   400,
			err:    true,
			errMsg: "{\"error\":\"id must be a number\"}",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/claim/"+tc.id, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.code, w.Code)
			if tc.err {
				assert.Equal(t, tc.errMsg, w.Body.String())
			} else {
				assert.Equal(t, string(itemBytes), w.Body.String())
			}
		})
	}
}

func TestPostClaimRoute(t *testing.T) {
	_, _, router := setupService(t)

	tests := []struct {
		name    string
		reqBody types.PostClaimRequest
		resBody string
		code    int
	}{
		{
			name: "valid request",
			reqBody: types.PostClaimRequest{
				Beneficiary: "liquid-address",
				Amount:      1000000000000,
				ClaimID:     1,
			},
			resBody: "{\"id\":1,\"tx-id\":\"0000000000000000000000000000000000000000000000000000000000000000\"}",
			code:    200,
		},
		{
			name:    "invalid request",
			reqBody: types.PostClaimRequest{},
			resBody: "{\"error\":\"Key: 'PostClaimRequest.Beneficiary' Error:Field validation for 'Beneficiary' failed on the 'required' tag\\nKey: 'PostClaimRequest.Amount' Error:Field validation for 'Amount' failed on the 'required' tag\\nKey: 'PostClaimRequest.ClaimID' Error:Field validation for 'ClaimID' failed on the 'required' tag\"}",
			code:    400,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			bodyBytes, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/claim", bytes.NewBuffer(bodyBytes))
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.code, w.Code)
			assert.Equal(t, tc.resBody, w.Body.String())
		})
	}
}
