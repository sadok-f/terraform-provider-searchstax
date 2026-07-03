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
	// The real API returns {"configs": ["name1", "name2"], "success": "true"};
	// older mocks used a {"results": [{...}]} wrapper.
	var named struct {
		Configs []string `json:"configs"`
	}
	out := ZookeeperConfigsList{}
	if err := json.Unmarshal(body, &named); err == nil && named.Configs != nil {
		for _, name := range named.Configs {
			out.Results = append(out.Results, ZookeeperConfig{Name: name})
		}
		return &out, nil
	}
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
	// The real API returns {"configs": [...], "success": "true"}; a non-2xx status
	// is already an error, so reaching here means the config was uploaded.
	if _, err := c.doRequest(req); err != nil {
		return nil, err
	}
	out := ZookeeperConfig{Name: cfg.Name}
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
	// The real API returns {"name": ..., "configs": [file paths], "success": ...};
	// the "configs" field is the list of files in the config. Fall back to a
	// "files" field for older mocks.
	var raw struct {
		Name    string   `json:"name"`
		Configs []string `json:"configs"`
		Files   []string `json:"files"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	out := ZookeeperConfig{Name: raw.Name, Files: raw.Files}
	if len(raw.Configs) > 0 {
		out.Files = raw.Configs
	}
	return &out, nil
}

func (c *Client) DeleteZookeeperConfig(accountName, deploymentID, name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/zookeeper-config/%s/", c.HostURL, accountName, deploymentID, name), nil)
	if err != nil {
		return err
	}
	// The real API returns {"message": ..., "success": "true"}; a non-2xx status
	// is already an error, so reaching here means the config was deleted.
	if _, err := c.doRequest(req); err != nil {
		return err
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
