package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const accountsBasePath = "organisation/accounts"

func NewClient(c *http.Client, baseURL url.URL) Service {
	return &httpClient{httpClient: c, baseURL: baseURL}
}

type httpClient struct {
	httpClient *http.Client
	baseURL    url.URL
}

func (c *httpClient) Create(ctx context.Context, account *Data) (*CreateSuccess, error) {
	req, err := c.newRequest(ctx, "POST", accountsBasePath, createRequest{Data: account})
	if err != nil {
		return nil, err
	}
	var res CreateSuccess
	resp, err := c.do(req, &res)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 409 {
		return nil, &AccountAlreadyExistsError{ID: account.ID}
	}

	err = checkStatusCode(resp)
	return &res, err
}

func (c *httpClient) Fetch(ctx context.Context, id string) (*FetchSuccess, error) {
	url := fmt.Sprintf("%s/%s", accountsBasePath, id)
	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var res FetchSuccess
	resp, err := c.do(req, &res)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, &AccountNotFoundError{ID: id}
	}

	err = checkStatusCode(resp)
	return &res, err
}

func (c *httpClient) Delete(ctx context.Context, id string, ver int64) error {
	url := fmt.Sprintf("%s/%s?version=%d", accountsBasePath, id, ver)
	req, err := c.newRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return &AccountNotFoundError{ID: id}
	} else if resp.StatusCode == 409 {
		return &InvalidVersionError{Ver: ver}
	}
	return checkStatusCode(resp)
}

func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HttpStatusError{StatusCode: resp.StatusCode}
	}
	return nil
}

func (c *httpClient) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Date", time.Now().Format(time.RFC1123))
	return req.WithContext(ctx), nil
}

func (c *httpClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		r := InvalidDataError{}
		err := json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			return resp, err
		}
		return resp, &r
	}

	if v != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

type createRequest struct {
	Data *Data `json:"data"`
}
