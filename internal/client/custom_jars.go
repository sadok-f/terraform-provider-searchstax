package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CustomJarsList struct {
	Results []CustomJar `json:"results"`
}

type CustomJar struct {
	Name string `json:"name,omitempty"`
}

func (c *Client) GetCustomJars(accountName, deploymentID string) (*CustomJarsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/solr/custom-jars/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := CustomJarsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadCustomJar uploads metadata for a custom jar.
// Note: mock API expects JSON and returns {"uploaded": true}.
func (c *Client) UploadCustomJar(accountName, deploymentID string, jar CustomJar) error {
	rb, err := json.Marshal(jar)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/solr/custom-jars/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Uploaded bool `json:"uploaded"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Uploaded {
		return fmt.Errorf("custom jar not uploaded")
	}
	return nil
}

func (c *Client) DeleteCustomJar(accountName, deploymentID, jarName string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/solr/custom-jars/%s/", c.HostURL, accountName, deploymentID, jarName), nil)
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
		return fmt.Errorf("custom jar not deleted")
	}
	return nil
}
