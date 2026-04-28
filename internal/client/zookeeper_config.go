package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ZookeeperConfigsList struct {
	Count    int64             `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []ZookeeperConfig `json:"results"`
}

type ZookeeperConfig struct {
	Name    string   `json:"name"`
	Created string   `json:"created,omitempty"`
	Files   []string `json:"files,omitempty"`
}

func (c *Client) GetZookeeperConfigs(accountName, deploymentID string) (*ZookeeperConfigsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/zookeeper-config/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := ZookeeperConfigsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadZookeeperConfig creates/uploads a new config.
// Mock expects JSON and returns {"uploaded": true, "name": "..."}.
func (c *Client) UploadZookeeperConfig(accountName, deploymentID string, cfg ZookeeperConfig) (*ZookeeperConfig, error) {
	rb, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/zookeeper-config/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Uploaded bool   `json:"uploaded"`
		Name     string `json:"name"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if !resp.Uploaded {
		return nil, fmt.Errorf("zookeeper config not uploaded")
	}
	out := ZookeeperConfig{Name: resp.Name}
	return &out, nil
}

func (c *Client) GetZookeeperConfig(accountName, deploymentID, name string) (*ZookeeperConfig, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/zookeeper-config/%s/", c.HostURL, accountName, deploymentID, name), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := ZookeeperConfig{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteZookeeperConfig(accountName, deploymentID, name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/zookeeper-config/%s/", c.HostURL, accountName, deploymentID, name), nil)
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
		return fmt.Errorf("zookeeper config not deleted")
	}
	return nil
}

type ZookeeperConfigDownload struct {
	Download string `json:"download"`
	Note     string `json:"note,omitempty"`
}

func (c *Client) DownloadZookeeperConfig(accountName, deploymentID, name string) (*ZookeeperConfigDownload, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/zookeeper-config/%s/download/", c.HostURL, accountName, deploymentID, name), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := ZookeeperConfigDownload{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
