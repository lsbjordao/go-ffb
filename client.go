package ffb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://servicos.jbrj.gov.br/v2/flora/"

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	userAgent  string
}

type Option func(*Client) error

func WithBaseURL(raw string) Option {
	return func(c *Client) error {
		u, err := url.Parse(raw)
		if err != nil {
			return fmt.Errorf("ffb: invalid base URL: %w", err)
		}
		if u.Scheme == "" || u.Host == "" {
			return errors.New("ffb: base URL must include scheme and host")
		}
		if !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		c.baseURL = u
		return nil
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("ffb: HTTP client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		c.userAgent = strings.TrimSpace(ua)
		return nil
	}
}

func NewClient(opts ...Option) (*Client, error) {
	base, err := url.Parse(DefaultBaseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		baseURL:    base,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userAgent:  "go-ffb/0.1.0",
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type APIError struct {
	StatusCode int
	URL        string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("ffb: API returned HTTP %d for %s", e.StatusCode, e.URL)
	}
	return fmt.Sprintf("ffb: API returned HTTP %d for %s: %s", e.StatusCode, e.URL, e.Body)
}

func (c *Client) get(ctx context.Context, parts ...string) ([]byte, error) {
	if c == nil || c.baseURL == nil || c.httpClient == nil {
		return nil, errors.New("ffb: uninitialized client")
	}

	u := *c.baseURL
	escapedPath := strings.TrimSuffix(u.EscapedPath(), "/")
	plainPath := strings.TrimSuffix(u.Path, "/")
	for _, p := range parts {
		escapedPath += "/" + url.PathEscape(p)
		plainPath += "/" + p
	}
	u.Path = plainPath
	u.RawPath = escapedPath

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ffb: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ffb: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ffb: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		const max = 1024
		msg := strings.TrimSpace(string(body))
		if len(msg) > max {
			msg = msg[:max] + "…"
		}
		return nil, &APIError{StatusCode: resp.StatusCode, URL: u.String(), Body: msg}
	}
	if !json.Valid(body) {
		return nil, errors.New("ffb: API returned invalid JSON")
	}
	return body, nil
}

func decodeStrings(body []byte) ([]string, error) {
	var values []string
	if err := json.Unmarshal(body, &values); err != nil {
		return nil, fmt.Errorf("ffb: expected JSON array of strings: %w", err)
	}
	return values, nil
}
