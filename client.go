package opennms

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"
)

// defaultAccept keeps text/plain acceptable at a lower q — some
// OpenNMS versions return 406 on /count endpoints if the Accept
// header does not include text/plain. JSON stays preferred via the
// implicit q=1.0.
const defaultAccept = "application/json, text/plain;q=0.9"

// Client is a thin synchronous wrapper for the OpenNMS REST API.
//
// All responses are decoded generic values: map[string]any for JSON
// objects, int for plain-text counts, string for text bodies. Failed
// requests return an *APIError.
type Client struct {
	v1URL      string
	v2URL      string
	username   string
	password   string
	httpClient *http.Client
	retries    int
}

type clientConfig struct {
	timeout    time.Duration
	noTimeout  bool
	retries    int
	skipVerify bool
	httpClient *http.Client
}

// Option configures a Client created by NewClient.
type Option func(*clientConfig)

// WithTimeout sets the read timeout for all HTTP requests. The
// connect timeout is capped at min(d, 10s) so unreachable hosts fail
// fast. Defaults to 30s. Pass 0 to disable timeouts entirely.
func WithTimeout(d time.Duration) Option {
	return func(c *clientConfig) {
		c.timeout = d
		c.noTimeout = d == 0
	}
}

// WithRetries sets the number of retries on connection errors and
// HTTP 500/502/503/504, with exponential backoff (0.5s factor).
// Defaults to 3. Pass 0 to disable retries.
func WithRetries(n int) Option {
	return func(c *clientConfig) { c.retries = n }
}

// WithInsecureSkipVerify disables SSL certificate verification. Only
// for self-signed certs in dev/test environments — not recommended in
// production.
func WithInsecureSkipVerify() Option {
	return func(c *clientConfig) { c.skipVerify = true }
}

// WithHTTPClient replaces the underlying *http.Client entirely. The
// timeout, retry-transport, and TLS options are then ignored (retries
// still apply — they are implemented above the transport).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *clientConfig) { c.httpClient = hc }
}

// NewClient returns a client for the OpenNMS instance at baseURL
// (e.g. "https://onms.example.com" — without the /opennms context
// path), authenticating with HTTP Basic auth. The username needs at
// minimum the "rest" role.
func NewClient(baseURL, username, password string, opts ...Option) *Client {
	cfg := clientConfig{timeout: 30 * time.Second, retries: 3}
	for _, o := range opts {
		o(&cfg)
	}
	hc := cfg.httpClient
	if hc == nil {
		connectTimeout := cfg.timeout
		if cfg.noTimeout || connectTimeout > 10*time.Second {
			connectTimeout = 10 * time.Second
		}
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: connectTimeout,
			}).DialContext,
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 20,
		}
		if cfg.skipVerify {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		timeout := cfg.timeout
		if cfg.noTimeout {
			timeout = 0
		}
		hc = &http.Client{Transport: transport, Timeout: timeout}
	}
	base := strings.TrimRight(baseURL, "/")
	return &Client{
		v1URL:      base + "/opennms/rest",
		v2URL:      base + "/opennms/api/v2",
		username:   username,
		password:   password,
		httpClient: hc,
		retries:    cfg.retries,
	}
}
