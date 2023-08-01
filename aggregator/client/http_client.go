package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boyanivskyy/toll-calculator/types"
	"github.com/sirupsen/logrus"
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
	req, err := http.NewRequest(http.MethodPost, c.Endpoint+"/aggregate", bytes.NewReader(b))
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
	invoiceRequest := types.GetInvoiceRequest{
		OBUID: int32(id),
	}
	b, err := json.Marshal(&invoiceRequest)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s?obuId=%d", c.Endpoint, "invoice", id)
	logrus.Infof("HERE endpoint %s", endpoint)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}

	invoice := &types.Invoice{}
	if err := json.NewDecoder(resp.Body).Decode(invoice); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return invoice, nil
}
