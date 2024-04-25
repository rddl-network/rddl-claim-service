package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/rddl-network/rddl-claim-service/types"
)

type IRCClient interface {
	GetClaim(ctx context.Context, id int) (res types.GetClaimResponse, err error)
	PostClaim(ctx context.Context, req types.PostClaimRequest) (res types.PostClaimResponse, err error)
}

type RCClient struct {
	baseURL string
	client  *http.Client
}

func NewRCClient(baseURL string, client *http.Client) *RCClient {
	if client == nil {
		client = &http.Client{}
	}
	return &RCClient{
		baseURL: baseURL,
		client:  client,
	}
}

func (rcc *RCClient) GetClaim(ctx context.Context, id int) (res types.GetClaimResponse, err error) {
	err = rcc.doRequest(ctx, http.MethodGet, rcc.baseURL+"/claim/"+strconv.Itoa(id), nil, &res)
	return
}

func (rcc *RCClient) PostClaim(ctx context.Context, req types.PostClaimRequest) (res types.PostClaimResponse, err error) {
	err = rcc.doRequest(ctx, http.MethodPost, rcc.baseURL+"/claim", nil, &res)
	return
}

func (rcc *RCClient) doRequest(ctx context.Context, method, url string, body interface{}, response interface{}) (err error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := rcc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &httpError{StatusCode: resp.StatusCode}
	}

	if response != nil {
		return json.NewDecoder(resp.Body).Decode(response)
	}

	return
}

type httpError struct {
	StatusCode int
}

func (e *httpError) Error() string {
	return http.StatusText(e.StatusCode)
}
