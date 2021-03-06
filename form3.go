package form3

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/althink/form3/accounts"
)

var defaultUrl string = "http://localhost:8080/v1/"

type Form3 struct {
	Accounts accounts.Service

	baseURL    url.URL
	httpClient *http.Client
}

type Option func(*Form3)

func WithHTTPClient(c *http.Client) Option {
	return func(f3 *Form3) {
		f3.httpClient = c
	}
}

func WithBaseURL(u url.URL) Option {
	return func(f3 *Form3) {
		f3.baseURL = u
	}
}

// NewClient creates new Form3 client.
func NewClient(opts ...Option) (*Form3, error) {
	url, err := url.Parse(defaultUrl)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	f3 := &Form3{
		httpClient: &http.Client{},
		baseURL:    *url,
	}

	for _, o := range opts {
		o(f3)
	}

	if !strings.HasSuffix(f3.baseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash: %q", f3.baseURL.String())
	}

	f3.Accounts = accounts.NewClient(f3.httpClient, f3.baseURL)

	return f3, nil
}
