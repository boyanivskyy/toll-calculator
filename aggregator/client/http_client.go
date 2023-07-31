package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boyanivskyy/toll-calculator/types"
)

type HttpClient struct {
	Endpoint string
}

func NewHttpClient(endpoint string) *HttpClient {
	return &HttpClient{
		Endpoint: endpoint,
	}
}

func (c *HttpClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}

	return nil
}

func (c *HttpClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	return &types.Invoice{
		OBUID:         1,
		TotalDistance: 123.123,
		TotalAmount:   1230,
	}, nil
}
