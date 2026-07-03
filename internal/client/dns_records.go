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

// UnmarshalJSON tolerates a "ttl" returned as either a JSON number (real API)
// or a string, normalizing it to a string.
func (d *DNSRecord) UnmarshalJSON(data []byte) error {
	type alias DNSRecord
	aux := &struct {
		TTL json.RawMessage `json:"ttl,omitempty"`
		*alias
	}{alias: (*alias)(d)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	d.TTL = coerceID(aux.TTL)
	return nil
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
	// The real API returns the updated DNS record object (no "updated" flag); a
	// non-2xx status is already an error, so reaching here means it was updated.
	out := DNSRecord{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
