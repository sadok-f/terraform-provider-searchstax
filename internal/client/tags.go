package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type TagsList struct {
	Tags []string `json:"tags"`
}

func (c *Client) GetTags(accountName, deploymentID string) (*TagsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/tags/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := TagsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type UpdateTagsRequest struct {
	Tags []string `json:"tags"`
}

func (c *Client) AddOrUpdateTags(accountName, deploymentID string, reqBody UpdateTagsRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/tags/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	// The real API returns the updated {"deployment", "tags"} object; a non-2xx
	// status is already an error, so reaching here means the tags were saved.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteTags(accountName, deploymentID string, reqBody UpdateTagsRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/tags/delete/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	// The real API returns the updated {"deployment", "tags"} object; a non-2xx
	// status is already an error, so reaching here means the tags were removed.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

type DeploymentsByTagList struct {
	Results []DeploymentTags `json:"results"`
}

type DeploymentTags struct {
	Deployment string   `json:"deployment"`
	UID        string   `json:"uid"`
	Tags       []string `json:"tags"`
}

func (d DeploymentTags) DeploymentUID() string {
	if d.Deployment != "" {
		return d.Deployment
	}
	return d.UID
}

type GetDeploymentsByTagRequest struct {
	Tags     []string `json:"tags"`
	Operator string   `json:"operator,omitempty"`
}

func (c *Client) GetDeploymentsByTag(accountName string, reqBody GetDeploymentsByTagRequest) (*DeploymentsByTagList, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/tags/get-deployments/", c.HostURL, accountName), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	// The real API returns a bare JSON array; decodeResults also tolerates a
	// {"results": [...]} wrapper.
	out := DeploymentsByTagList{}
	if err := decodeResults(body, &out.Results); err != nil {
		return nil, err
	}
	return &out, nil
}
