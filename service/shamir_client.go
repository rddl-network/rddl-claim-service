package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IShamirClient interface {
	IssueTransaction(amount string, address string) (hash string, err error)
}

type ShamirClient struct {
	host string
}

func NewShamirClient(host string) *ShamirClient {
	return &ShamirClient{host: host}
}

type SendTokensRequest struct {
	Recipient string `json:"recipient"`
	Amount    string `json:"amount"`
}

type SendTokensResponse struct {
	TxID string `binding:"required" json:"tx-id"`
}

func (sc *ShamirClient) IssueTransaction(amount string, address string) (hash string, err error) {
	url := fmt.Sprintf("https://%s/send", sc.host)

	body := &SendTokensRequest{
		Recipient: address,
		Amount:    amount,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var sendTokenResp SendTokensResponse
	err = json.Unmarshal(bodyBytes, &sendTokenResp)
	if err != nil {
		return "", err
	}

	return sendTokenResp.TxID, nil
}
