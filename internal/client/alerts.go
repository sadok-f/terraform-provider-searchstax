package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AlertMetricsList struct {
	Results []AlertMetric `json:"results"`
}

type AlertMetric struct {
	Metric string `json:"metric"`
	Unit   string `json:"unit"`
}

func (c *Client) GetAlertMetrics(accountName string) (*AlertMetricsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/alerts/metrics/", c.HostURL, accountName), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := AlertMetricsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type IncidentsList struct {
	Results []Incident `json:"results"`
}

type Incident struct {
	ID int64 `json:"id,omitempty"`
}

func (c *Client) GetIncidents(accountName, deploymentID string) (*IncidentsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/incidents/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := IncidentsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type AlertsList struct {
	Results []Alert `json:"results"`
}

type Alert struct {
	ID int64 `json:"id,omitempty"`
}

func (c *Client) GetAlerts(accountName, deploymentID string) (*AlertsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := AlertsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateAlert(accountName, deploymentID string, alert Alert) (int64, error) {
	rb, err := json.Marshal(alert)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return 0, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return 0, err
	}
	var resp struct {
		Created bool  `json:"created"`
		ID      int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}
	if !resp.Created {
		return 0, fmt.Errorf("alert not created")
	}
	return resp.ID, nil
}

func (c *Client) UpdateAlert(accountName, deploymentID string, id int64, alert Alert) error {
	rb, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/%d/", c.HostURL, accountName, deploymentID, id), strings.NewReader(string(rb)))
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
		return fmt.Errorf("alert not updated")
	}
	return nil
}

func (c *Client) DeleteAlert(accountName, deploymentID string, id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/%d/", c.HostURL, accountName, deploymentID, id), nil)
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
		return fmt.Errorf("alert not deleted")
	}
	return nil
}

type HeartbeatsList struct {
	Results []Heartbeat `json:"results"`
}

type Heartbeat struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Host string `json:"host,omitempty"`
}

func (c *Client) GetHeartbeats(accountName, deploymentID string) (*HeartbeatsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/heartbeat/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := HeartbeatsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateHeartbeat(accountName, deploymentID string, hb Heartbeat) (int64, error) {
	rb, err := json.Marshal(hb)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/heartbeat/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return 0, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return 0, err
	}
	var resp struct {
		Created bool  `json:"created"`
		ID      int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}
	if !resp.Created {
		return 0, fmt.Errorf("heartbeat not created")
	}
	return resp.ID, nil
}

func (c *Client) GetHeartbeat(accountName, deploymentID string, id int64) (*Heartbeat, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/heartbeat/%d/", c.HostURL, accountName, deploymentID, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := Heartbeat{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateHeartbeat(accountName, deploymentID string, id int64, hb Heartbeat) error {
	rb, err := json.Marshal(hb)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/heartbeat/%d/", c.HostURL, accountName, deploymentID, id), strings.NewReader(string(rb)))
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
		return fmt.Errorf("heartbeat not updated")
	}
	return nil
}

func (c *Client) DeleteHeartbeat(accountName, deploymentID string, id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/heartbeat/%d/", c.HostURL, accountName, deploymentID, id), nil)
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
		return fmt.Errorf("heartbeat not deleted")
	}
	return nil
}
