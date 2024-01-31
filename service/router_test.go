package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rddl-network/rddl-claim-service/service"
	"github.com/stretchr/testify/assert"
)

func TestGetClaimRoute(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("/claim/%s", tc.id), nil)
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
	t.Parallel()
	_, _, router := setupService(t)

	tests := []struct {
		name    string
		reqBody service.ClaimRequestBody
		resBody string
		code    int
	}{
		{
			name: "valid request",
			reqBody: service.ClaimRequestBody{
				Beneficiary: "liquid-address",
				Amount:      "10000.00000",
			},
			resBody: "{\"hash\":\"0000000000000000000000000000000000000000000000000000000000000000\",\"id\":1,\"message\":\"claim enqueued\"}",
			code:    200,
		},
		{
			name:    "invalid request",
			reqBody: service.ClaimRequestBody{},
			resBody: "{\"error\":\"Key: 'ClaimRequestBody.Beneficiary' Error:Field validation for 'Beneficiary' failed on the 'required' tag\\nKey: 'ClaimRequestBody.Amount' Error:Field validation for 'Amount' failed on the 'required' tag\"}",
			code:    400,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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
