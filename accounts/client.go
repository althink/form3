package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func NewClient(c *http.Client, baseURL url.URL) Service {
	return &httpClient{httpClient: c, baseURL: baseURL}
}

type httpClient struct {
	httpClient *http.Client
	baseURL    url.URL
	UserAgent  string
}

func (c *httpClient) Create(ctx context.Context, account *Data) (*CreateSuccess, error) {
	req, err := c.newRequest("POST", "organisation/accounts", createRequest{Data: account})
	if err != nil {
		return nil, err
	}
	var res CreateSuccess
	resp, err := c.do(req, &res)
	if err != nil {
		return nil, err
	}
	err = checkStatusCode(resp)
	return &res, err
}

func (c *httpClient) Fetch(ctx context.Context, id string) (*FetchSuccess, error) {
	url := fmt.Sprintf("organisation/accounts/%s", id)
	req, err := c.newRequest("GET", url, nil)
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
	url := fmt.Sprintf("organisation/accounts/%s?version=%d", id, ver)
	req, err := c.newRequest("DELETE", url, nil)
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
		return errors.New(fmt.Sprintf("Invalid status code: %v", resp.StatusCode))
	}
	return nil
}

func (c *httpClient) newRequest(method, path string, body interface{}) (*http.Request, error) {
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
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
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

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

type createRequest struct {
	Data *Data `json:"data"`
}
