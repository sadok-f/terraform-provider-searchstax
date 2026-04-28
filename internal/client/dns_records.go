package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DNSRecordsList struct {
	Results []DNSRecord `json:"results"`
}

type DNSRecord struct {
	Name       string `json:"name"`
	Deployment string `json:"deployment"`
	TTL        string `json:"ttl"`
}

func (c *Client) GetDNSRecords(accountName string) (*DNSRecordsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/dns-record/", c.HostURL, accountName), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := DNSRecordsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetDNSRecord(accountName string, name string) (*DNSRecord, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/dns-record/%s/", c.HostURL, accountName, name), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := DNSRecord{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type AssociateDNSRecordRequest struct {
	Deployment string `json:"deployment"`
	TTL        string `json:"ttl,omitempty"`
}

func (c *Client) AssociateDNSRecord(accountName string, name string, reqBody AssociateDNSRecordRequest) (*DNSRecord, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/account/%s/dns-record/%s/", c.HostURL, accountName, name), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Updated bool `json:"updated"`
		DNSRecord
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if !resp.Updated {
		return nil, fmt.Errorf("dns record not updated")
	}
	out := resp.DNSRecord
	return &out, nil
}
