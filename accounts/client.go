package accounts

import (
	"context"
	"net/http"
)

func NewClient(c *http.Client) Service {
	return &httpClient{c: c}
}

type httpClient struct {
	c *http.Client
}

func (c *httpClient) Create(ctx context.Context, account *Data) (*CreateSuccess, error) {
	return nil, nil
}

func (c *httpClient) Fetch(ctx context.Context, id string) (*FetchSuccess, error) {
	return nil, nil
}

func (c *httpClient) Delete(ctx context.Context, id string, version string) error {
	return nil
}

type createRequest struct {
	Data *Data `json:"data"`
}
