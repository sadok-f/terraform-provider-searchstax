package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CreateAPIKeyRequest struct {
	Scope []string `json:"scope"`
}

type CreateAPIKeyResponse struct {
	APIKey string `json:"apikey"`
}

func (c *Client) CreateAPIKey(accountName string, reqBody CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/apikey/", c.HostURL, accountName), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := CreateAPIKeyResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if out.APIKey == "" {
		return nil, fmt.Errorf("api key not created")
	}
	return &out, nil
}

type AssociateAPIKeyRequest struct {
	APIKey     string `json:"apikey"`
	Deployment string `json:"deployment"`
}

type APIKeyDeploymentsResponse struct {
	Deployments []string `json:"deployments"`
}

func (c *Client) AssociateAPIKey(accountName string, reqBody AssociateAPIKeyRequest) (*APIKeyDeploymentsResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/apikey/associate/", c.HostURL, accountName), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := APIKeyDeploymentsResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DisassociateAPIKey(accountName string, reqBody AssociateAPIKeyRequest) (*APIKeyDeploymentsResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/apikey/disassociate/", c.HostURL, accountName), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := APIKeyDeploymentsResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type APIKeyDeploymentsRequest struct {
	APIKey string `json:"apikey"`
}

func (c *Client) GetAPIKeyDeployments(accountName string, reqBody APIKeyDeploymentsRequest) (*APIKeyDeploymentsResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/apikey/deployments/", c.HostURL, accountName), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := APIKeyDeploymentsResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if len(out.Deployments) == 0 {
		var mock struct {
			Results []struct {
				UID string `json:"uid"`
			} `json:"results"`
		}
		if err := json.Unmarshal(body, &mock); err == nil {
			for _, r := range mock.Results {
				out.Deployments = append(out.Deployments, r.UID)
			}
		}
	}
	return &out, nil
}

type DeploymentAPIKeysRequest struct {
	Deployment string `json:"deployment"`
}

type DeploymentAPIKeysResponse struct {
	APIKey []string `json:"apikey"`
}

func (c *Client) GetDeploymentAPIKeys(accountName string, reqBody DeploymentAPIKeysRequest) (*DeploymentAPIKeysResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/apikey/list/", c.HostURL, accountName), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := DeploymentAPIKeysResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if len(out.APIKey) == 0 {
		var mock struct {
			Results []struct {
				APIKey string `json:"apikey"`
			} `json:"results"`
		}
		if err := json.Unmarshal(body, &mock); err == nil {
			for _, r := range mock.Results {
				out.APIKey = append(out.APIKey, r.APIKey)
			}
		}
	}
	return &out, nil
}

type RevokeAPIKeyRequest struct {
	APIKey string `json:"apikey"`
}

func (c *Client) RevokeAPIKey(accountName string, reqBody RevokeAPIKeyRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/apikey/revoke/", c.HostURL, accountName), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Success string `json:"success"`
		Revoked bool   `json:"revoked"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if resp.Success == "" && !resp.Revoked {
		return fmt.Errorf("api key not revoked")
	}
	return nil
}
