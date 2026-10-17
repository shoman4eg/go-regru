// Package regru is a client for REG.API 2 (https://www.reg.ru/reseller/api2doc).
//
// Usage:
//
//	client := regru.NewClient(nil).WithAuthPassword("test", "test")
//	answer, _, err := client.Domain.Check(ctx, &regru.DomainCheckOptions{...})
//
// Client TLS certificate ("дополнительная проверка по ssl-сертификату") is
// configured on the *http.Client passed to NewClient via its Transport.
package regru

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "https://api.reg.ru/api/regru2/"
	userAgent      = "go-regru"
)

// Client manages communication with REG.API 2.
type Client struct {
	client *http.Client

	// BaseURL must have a trailing slash.
	BaseURL   *url.URL
	UserAgent string

	// Lang selects the language of error_text: "en" (default), "ru" or "th".
	Lang string

	username string
	password string
	signer   *rsa.PrivateKey

	common service

	User    *UserService
	Bill    *BillService
	Service *ServiceService
	Domain  *DomainService
	Zone    *ZoneService
	DNSSEC  *DNSSECService
	Hosting *HostingService
	Folder  *FolderService
	Shop    *ShopService
}

type service struct {
	client *Client
}

// NewClient returns a new REG.API client. If httpClient is nil, http.DefaultClient is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.User = (*UserService)(&c.common)
	c.Bill = (*BillService)(&c.common)
	c.Service = (*ServiceService)(&c.common)
	c.Domain = (*DomainService)(&c.common)
	c.Zone = (*ZoneService)(&c.common)
	c.DNSSEC = (*DNSSECService)(&c.common)
	c.Hosting = (*HostingService)(&c.common)
	c.Folder = (*FolderService)(&c.common)
	c.Shop = (*ShopService)(&c.common)
	return c
}

// WithAuthPassword sets username/password authentication.
// password may be the account password or the alternative API password.
func (c *Client) WithAuthPassword(username, password string) *Client {
	c.username, c.password, c.signer = username, password, nil
	return c
}

// WithSignature sets signature authentication. The matching certificate must be
// uploaded in the API settings and presented via the http.Client TLS config.
func (c *Client) WithSignature(username string, key *rsa.PrivateKey) *Client {
	c.username, c.password, c.signer = username, "", key
	return c
}

// NewRequest builds a POST request for the API function at path (e.g. "domain/check").
// params are serialized to JSON and passed as input_data; nil means no parameters.
func (c *Client) NewRequest(ctx context.Context, path string, params any) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("output_format", "json")
	if c.username != "" {
		form.Set("username", c.username)
	}
	if c.password != "" {
		form.Set("password", c.password)
	}
	if c.Lang != "" {
		form.Set("lang", c.Lang)
	}

	var input []byte
	if params != nil {
		if input, err = json.Marshal(params); err != nil {
			return nil, err
		}
		form.Set("input_format", "json")
		form.Set("input_data", string(input))
	}

	if c.signer != nil {
		sig, err := sign(c.signer, form, input)
		if err != nil {
			return nil, err
		}
		form.Set("sig", sig)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Response wraps the HTTP response of an API call.
type Response struct {
	*http.Response
}

type envelope struct {
	Result      string          `json:"result"`
	Answer      json.RawMessage `json:"answer"`
	ErrorCode   string          `json:"error_code"`
	ErrorText   string          `json:"error_text"`
	ErrorParams map[string]any  `json:"error_params"`
}

// Do sends req and decodes the "answer" field into v (if v is non-nil).
// An API-level error is returned as *ErrorResponse.
func (c *Client) Do(req *http.Request, v any) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response := &Response{Response: resp}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	var env envelope
	if err := json.Unmarshal(body, &env); err != nil || env.Result == "" {
		return response, fmt.Errorf("regru: %s %s: unexpected response (HTTP %d): %.200s",
			req.Method, req.URL, resp.StatusCode, bytes.TrimSpace(body))
	}
	if env.Result != "success" {
		return response, &ErrorResponse{Response: resp, Code: env.ErrorCode, Text: env.ErrorText, Params: env.ErrorParams}
	}
	if v != nil && len(env.Answer) > 0 && string(env.Answer) != "null" {
		if err := json.Unmarshal(env.Answer, v); err != nil {
			return response, fmt.Errorf("regru: decode answer: %w", err)
		}
	}
	return response, nil
}

// call is the shared body of every API method: build, send, decode into a new T.
func call[T any](ctx context.Context, c *Client, path string, params any) (*T, *Response, error) {
	req, err := c.NewRequest(ctx, path, params)
	if err != nil {
		return nil, nil, err
	}
	out := new(T)
	resp, err := c.Do(req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}
