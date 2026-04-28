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
		return fmt.Errorf("tags not updated")
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
		return fmt.Errorf("tags not deleted")
	}
	return nil
}

type DeploymentsByTagList struct {
	Results []DeploymentTags `json:"results"`
}

type DeploymentTags struct {
	UID  string   `json:"uid"`
	Tags []string `json:"tags"`
}

type GetDeploymentsByTagRequest struct {
	Tags []string `json:"tags"`
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
	out := DeploymentsByTagList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
