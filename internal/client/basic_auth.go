package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) EnableBasicAuth(accountName, deploymentID string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/enable/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return false, err
	}
	body, err := c.doRequest(req)
	if err != nil {
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

func (c *Client) DisableBasicAuth(accountName, deploymentID string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/disable/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return false, err
	}
	body, err := c.doRequest(req)
	if err != nil {
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
