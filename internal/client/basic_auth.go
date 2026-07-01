package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (c *Client) EnableBasicAuth(accountName, deploymentID string) (bool, error) {
	const (
		attempts = 10
		backoff  = 15 * time.Second
	)

	var lastErr error
	for i := 0; i < attempts; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/enable/", c.HostURL, accountName, deploymentID), nil)
		if err != nil {
			return false, err
		}
		body, err := c.doRequest(req)
		if err != nil {
			// Retry on transient (5xx / network) errors — the cluster may be
			// briefly unavailable while a restart is in progress.
			if isTransient(err) {
				lastErr = err
				if i < attempts-1 {
					time.Sleep(backoff)
				}
				continue
			}
			return false, err
		}
		var resp struct {
			Enabled bool `json:"enabled"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			return false, err
		}
		return resp.Enabled, nil
	}
	return false, fmt.Errorf("basic auth not enabled after %d attempts: %w", attempts, lastErr)
}

func (c *Client) DisableBasicAuth(accountName, deploymentID string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/disable/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return false, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		// The real API returns 400 "No basic authentication enabled for this
		// deployment." when auth is already disabled. Treat that as success so
		// destroy is idempotent.
		var httpErr *HTTPStatusError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusBadRequest &&
			strings.Contains(httpErr.Body, "No basic authentication enabled") {
			return false, nil
		}
		return false, err
	}
	var resp struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, err
	}
	return resp.Enabled, nil
}

type SetBasicAuthPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) SetBasicAuthPassword(accountName, deploymentID string, reqBody SetBasicAuthPasswordRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/set-password/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Updated bool `json:"updated"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Updated {
		return fmt.Errorf("basic auth password not updated")
	}
	return nil
}

type SetBasicAuthRoleRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// IsBasicAuthEnabled reports whether Solr basic auth appears enabled for a deployment.
func (c *Client) IsBasicAuthEnabled(accountName, deploymentID string) (bool, error) {
	_, err := c.GetDeploymentUsers(accountName, deploymentID)
	return err == nil, nil
}

func (c *Client) SetBasicAuthRole(accountName, deploymentID string, reqBody SetBasicAuthRoleRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/set-role/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Updated bool `json:"updated"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Updated {
		return fmt.Errorf("basic auth role not updated")
	}
	return nil
}
