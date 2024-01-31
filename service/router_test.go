package service_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClaimRoute(t *testing.T) {
	t.Parallel()
	app, _, router := setupService(t)

	items := createNRedeemClaim(app, 1)
	itemBytes, _ := json.Marshal(items[0])

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
			req, _ := http.NewRequest("GET", fmt.Sprintf("/claim/%s", tc.id), nil)
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
