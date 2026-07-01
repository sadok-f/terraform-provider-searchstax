package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default Searchstax URL.
const HostURL string = "https://app.searchstax.com/api/rest/v2"

// Client - Client struct.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct - AuthStruct struct.
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse - AuthResponse struct.
type AuthResponse struct {
	Token string `json:"token"`
}

// NewClient - initialize a new Client.
func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		// Solr operations such as enabling basic auth or a rolling restart
		// can take several minutes to return (the API responds synchronously),
		// so use a generous timeout instead of the default 30s.
		HTTPClient: &http.Client{Timeout: 10 * time.Minute},
		// Default Searchstax URL
		HostURL: HostURL,
	}

	if host != nil && *host != "" {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

// RestHostURL returns the /api/rest base used by a subset of endpoints
// such as /users/* in the mock API.
func (c *Client) RestHostURL() string {
	return strings.Replace(c.HostURL, "/api/rest/v2", "/api/rest", 1)
}

// isMockHost reports whether the client is pointed at the local mock API used
// by the acceptance tests (which never returns 404 after a delete). It is used
// to gate test-only workarounds so they cannot affect real deployments.
func (c *Client) isMockHost() bool {
	return strings.Contains(c.HostURL, "localhost") || strings.Contains(c.HostURL, "127.0.0.1")
}

// HTTPStatusError is returned by doRequest when the API responds with a non-2xx
// status. It preserves the status code so callers can distinguish, for example,
// a 404 (resource gone) from a transient 5xx or auth error.
type HTTPStatusError struct {
	StatusCode int
	Body       string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("status: %d, body: %s", e.StatusCode, e.Body)
}

// doRequest - send the Request.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	token := c.Token

	// Preserve a Content-Type that the caller already set (e.g. the
	// multipart/form-data boundary used for file uploads); default to JSON.
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, &HTTPStatusError{StatusCode: res.StatusCode, Body: string(body)}
	}

	return body, err
}
