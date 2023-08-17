package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetPrivateVpc - Returns list of private_vpc_list (no auth required).
func (c *Client) GetPrivateVpc(accountName string) (*PrivateVpcList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/privatevpc/", c.HostURL, accountName), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	privateVpcs := PrivateVpcList{}
	err = json.Unmarshal(body, &privateVpcs)
	if err != nil {
		return nil, err
	}

	return &privateVpcs, nil
}

// PrivateVpcList - PrivateVpcList struct.
type PrivateVpcList struct {
	Count    int32        `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []PrivateVpc `json:"results"`
}

// PrivateVpc - PrivateVpc struct.
type PrivateVpc struct {
	ID           int64  `json:"id"`
	Account      string `json:"account"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Region       string `json:"region"`
	AddressSpace string `json:"address_space"`
}
