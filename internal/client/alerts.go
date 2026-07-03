package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// decodeNamedList unmarshals a list response that may be wrapped under a
// specific key (real API, e.g. {"alerts": [...]}), under "results" (older mock),
// or returned as a bare array, into results (a pointer to a slice).
func decodeNamedList(body []byte, key string, results any) error {
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(body, &wrapper); err == nil {
		if raw, ok := wrapper[key]; ok && len(bytes.TrimSpace(raw)) > 0 {
			return json.Unmarshal(raw, results)
		}
		if raw, ok := wrapper["results"]; ok && len(bytes.TrimSpace(raw)) > 0 {
			return json.Unmarshal(raw, results)
		}
		return nil
	}
	return json.Unmarshal(body, results)
}

type AlertMetricsList struct {
	Results []AlertMetric `json:"results"`
}

type AlertMetric struct {
	Metric string
	Unit   string
}

// UnmarshalJSON reads the real API shape {"name", "unit": [...], "description"}
// (and tolerates the older {"metric", "unit"} mock shape), reducing the unit
// list to its first value.
func (m *AlertMetric) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name   string          `json:"name"`
		Metric string          `json:"metric"`
		Unit   json.RawMessage `json:"unit"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw.Name != "" {
		m.Metric = raw.Name
	} else {
		m.Metric = raw.Metric
	}
	unit := bytes.TrimSpace(raw.Unit)
	if len(unit) > 0 && unit[0] == '[' {
		var units []string
		if err := json.Unmarshal(unit, &units); err == nil && len(units) > 0 {
			m.Unit = units[0]
		}
	} else if len(unit) > 0 {
		_ = json.Unmarshal(unit, &m.Unit)
	}
	return nil
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
	if err := decodeResults(body, &out.Results); err != nil {
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
	if err := decodeNamedList(body, "incidents", &out.Results); err != nil {
		return nil, err
	}
	return &out, nil
}

type AlertsList struct {
	Results []Alert `json:"results"`
}

type Alert struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
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
	if err := decodeNamedList(body, "alerts", &out.Results); err != nil {
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
	// The real API confirms with {"message": "Success!"} and does not return the
	// new id, so look it up from the alerts list (by name when available).
	if _, err := c.doRequest(req); err != nil {
		return 0, err
	}
	list, err := c.GetAlerts(accountName, deploymentID)
	if err != nil {
		return 0, err
	}
	for _, a := range list.Results {
		if alert.Name != "" && a.Name == alert.Name {
			return a.ID, nil
		}
	}
	if len(list.Results) > 0 {
		return list.Results[len(list.Results)-1].ID, nil
	}
	return 0, nil
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
	// The real API returns {"message": ..., "success": "true"}; a non-2xx status
	// is already an error, so reaching here means the alert was updated.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteAlert(accountName, deploymentID string, id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/%d/", c.HostURL, accountName, deploymentID, id), nil)
	if err != nil {
		return err
	}
	// The real API returns {"message": ..., "success": "true"}; a non-2xx status
	// is already an error, so reaching here means the alert was deleted.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

type HeartbeatsList struct {
	Results []Heartbeat `json:"results"`
}

type Heartbeat struct {
	ID             int64    `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Host           string   `json:"host,omitempty"`
	Interval       string   `json:"interval,omitempty"`
	MaxAlerts      string   `json:"max_alerts,omitempty"`
	Email          []string `json:"email,omitempty"`
	WebhookTrigger int64    `json:"webhook_trigger,omitempty"`
	WebhookResolve int64    `json:"webhook_resolve,omitempty"`
	Status         string   `json:"status,omitempty"`
}

// UnmarshalJSON tolerates "interval" and "max_alerts" returned as JSON numbers
// (real API) or strings, normalizing them to strings.
func (h *Heartbeat) UnmarshalJSON(data []byte) error {
	type alias Heartbeat
	aux := &struct {
		Interval  json.RawMessage `json:"interval,omitempty"`
		MaxAlerts json.RawMessage `json:"max_alerts,omitempty"`
		*alias
	}{alias: (*alias)(h)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	h.Interval = coerceID(aux.Interval)
	h.MaxAlerts = coerceID(aux.MaxAlerts)
	return nil
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
	if err := decodeNamedList(body, "alerts", &out.Results); err != nil {
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
	// The real API confirms with {"message": "Success!"} and does not return the
	// new id, so look it up from the heartbeat list (by name when available).
	if _, err := c.doRequest(req); err != nil {
		return 0, err
	}
	list, err := c.GetHeartbeats(accountName, deploymentID)
	if err != nil {
		return 0, err
	}
	for _, h := range list.Results {
		if hb.Name != "" && h.Name == hb.Name {
			return h.ID, nil
		}
	}
	if len(list.Results) > 0 {
		return list.Results[len(list.Results)-1].ID, nil
	}
	return 0, nil
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
	// The real API returns {"message": ..., "success": "true"}; a non-2xx status
	// is already an error, so reaching here means the heartbeat was updated.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteHeartbeat(accountName, deploymentID string, id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/alerts/heartbeat/%d/", c.HostURL, accountName, deploymentID, id), nil)
	if err != nil {
		return err
	}
	// The real API returns {"message": ..., "success": "true"}; a non-2xx status
	// is already an error, so reaching here means the heartbeat was deleted.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}
