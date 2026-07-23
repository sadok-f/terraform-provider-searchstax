package client

import (
	"encoding/json"
	"errors"
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
	// The real API returns a bare JSON array; decodeResults also tolerates a
	// {"results": [...]} wrapper.
	out := IPFiltersList{}
	if err := decodeResults(body, &out.Results); err != nil {
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
	return c.ipFilterAction("add-cidr-ip", accountName, deploymentID, reqBody)
}

func (c *Client) UpdateIPFilter(accountName, deploymentID string, reqBody IPFilterUpsertRequest) error {
	return c.ipFilterAction("update-cidr-ip", accountName, deploymentID, reqBody)
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
	// The real API confirms with {"detail": null, "success": true}; a non-2xx
	// status is already an error, so reaching here means the filter was deleted.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *Client) ipFilterAction(action, accountName, deploymentID string, reqBody IPFilterUpsertRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/ip-filter/%s/", c.HostURL, accountName, deploymentID, action), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	// The real API confirms with {"detail": null, "success": true}; a non-2xx
	// status is already an error, so reaching here means the action succeeded.
	if _, err := c.doRequest(req); err != nil {
		// The API can return HTTP 400 with a no-op message when the same CIDR
		// configuration already exists. Treat that specific response as success
		// so update/create operations remain idempotent.
		var httpErr *HTTPStatusError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusBadRequest {
			body := strings.ToLower(httpErr.Body)
			if strings.Contains(body, "already exist") && strings.Contains(body, "no change performed") {
				return nil
			}
		}
		return err
	}
	return nil
}
