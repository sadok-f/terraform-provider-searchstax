package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type IPFiltersList struct {
	Count    int64      `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []IPFilter `json:"results"`
}

type IPFilter struct {
	Services    []string `json:"services"`
	CIDRIP      string   `json:"cidr_ip"`
	Description string   `json:"description"`
}

func (c *Client) GetIPFilters(accountName, deploymentID string) (*IPFiltersList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/ip-filter/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := IPFiltersList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type IPFilterUpsertRequest struct {
	Services    []string `json:"services"`
	CIDRIP      string   `json:"cidr_ip"`
	Description string   `json:"description,omitempty"`
}

func (c *Client) AddIPFilter(accountName, deploymentID string, reqBody IPFilterUpsertRequest) error {
	return c.ipFilterAction("add-cidr-ip", accountName, deploymentID, reqBody, "added")
}

func (c *Client) UpdateIPFilter(accountName, deploymentID string, reqBody IPFilterUpsertRequest) error {
	return c.ipFilterAction("update-cidr-ip", accountName, deploymentID, reqBody, "updated")
}

type IPFilterDeleteRequest struct {
	CIDRIP string `json:"cidr_ip"`
}

func (c *Client) DeleteIPFilter(accountName, deploymentID string, reqBody IPFilterDeleteRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/ip-filter/delete-cidr-ip/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Deleted bool `json:"deleted"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Deleted {
		return fmt.Errorf("ip filter not deleted")
	}
	return nil
}

func (c *Client) ipFilterAction(action, accountName, deploymentID string, reqBody IPFilterUpsertRequest, expectKey string) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/ip-filter/%s/", c.HostURL, accountName, deploymentID, action), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	// mock uses one of: {"added":true}, {"updated":true}
	var m map[string]bool
	if err := json.Unmarshal(body, &m); err != nil {
		return err
	}
	if !m[expectKey] {
		return fmt.Errorf("ip filter action %s failed", action)
	}
	return nil
}
