package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type DeploymentHealth struct {
	Status         string `json:"status"`
	Level          string `json:"level"`
	DeploymentUID  string `json:"deployment_uid"`
}

func (c *Client) GetDeploymentHealth(accountName, deploymentID string) (*DeploymentHealth, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/deployment-health/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := DeploymentHealth{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type CollectionsHealth struct {
	Success       bool     `json:"success"`
	Healthy       bool     `json:"healthy"`
	Error         string   `json:"error"`
	DeploymentUID string   `json:"deployment_uid"`
	Collections   []string `json:"collections"`
}

func (c *Client) GetCollectionsHealth(accountName, deploymentID string) (*CollectionsHealth, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/collection-health/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := CollectionsHealth{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type DeploymentServersList struct {
	Results []DeploymentServer `json:"results"`
}

type DeploymentServer struct {
	SN             int64  `json:"sn"`
	Node           string `json:"node"`
	PrivateAddress string `json:"private_address"`
	DNSAddress     string `json:"dns_address"`
	Status         string `json:"status"`
	StatusDetails  string `json:"status_details"`
	Solr           bool   `json:"solr"`
	Zookeeper      bool   `json:"zookeeper"`
	Role           string `json:"role"`
}

func (c *Client) GetDeploymentServers(accountName, deploymentID string) (*DeploymentServersList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/server/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := DeploymentServersList{}
	if err := json.Unmarshal(body, &out); err != nil {
		// Real API may return a bare array.
		var servers []DeploymentServer
		if err := json.Unmarshal(body, &servers); err != nil {
			return nil, err
		}
		out.Results = servers
	}
	return &out, nil
}

type ServerHostStatus struct {
	Level  string `json:"level"`
	Status string `json:"status"`
	Node   string `json:"node"`
}

func (c *Client) GetServerHostStatus(accountName, deploymentID, node string) (*ServerHostStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/server/%s/host-status/", c.HostURL, accountName, deploymentID, node), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := ServerHostStatus{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type RollingRestartRequest struct {
	Solr      bool `json:"solr"`
	Zookeeper bool `json:"zookeeper"`
}

type RollingRestartResponse struct {
	Detail  string `json:"detail"`
	Message string `json:"message"`
	Queued  bool   `json:"queued"`
}

func (c *Client) RollingRestart(accountName, deploymentID string, reqBody RollingRestartRequest) (*RollingRestartResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/account/%s/deployment/%s/rolling-restart/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := RollingRestartResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) StartSolr(accountName, deploymentID, node string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/server/%s/start-solr/", c.HostURL, accountName, deploymentID, node), nil)
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Success bool `json:"success"`
		Started bool `json:"started"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Success && !resp.Started {
		return fmt.Errorf("solr node not started")
	}
	return nil
}

func (c *Client) StopSolr(accountName, deploymentID, node string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/server/%s/stop-solr/", c.HostURL, accountName, deploymentID, node), nil)
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Success bool `json:"success"`
		Stopped bool `json:"stopped"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Success && !resp.Stopped {
		return fmt.Errorf("solr node not stopped")
	}
	return nil
}

type PlansList struct {
	Count    int32  `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Plan `json:"results"`
}

type Plan struct {
	Name                 string       `json:"name"`
	Plan                 string       `json:"plan"`
	Description          string       `json:"description"`
	PlanType             string       `json:"plan_type"`
	Application          string       `json:"-"`
	ApplicationVersions  []string     `json:"application_versions"`
	PlanRegions          []PlanRegion `json:"plan_regions"`
	TrialAvailable       bool         `json:"trial_available"`
}

type PlanRegion struct {
	Price                          float64 `json:"price"`
	AdditionalApplicationNodePrice float64 `json:"additional_application_node_price"`
	AdditionalZookeeperNodePrice   float64 `json:"additional_zookeeper_node_price"`
	RegionID                       string  `json:"region_id"`
	CloudProvider                  string  `json:"cloud_provider"`
	CloudProviderID                string  `json:"cloud_provider_id"`
}

func (p *Plan) UnmarshalJSON(data []byte) error {
	type planAlias Plan
	aux := struct {
		Application json.RawMessage `json:"application"`
		*planAlias
	}{
		planAlias: (*planAlias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.Application) > 0 {
		switch aux.Application[0] {
		case '"':
			_ = json.Unmarshal(aux.Application, &p.Application)
		case '{':
			var app struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(aux.Application, &app); err == nil {
				p.Application = app.Name
			}
		}
	}
	if p.Name == "" {
		p.Name = p.Plan
	}
	return nil
}

func (c *Client) GetPlans(accountName, application, planType string, page int) (*PlansList, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	if application != "" {
		q.Set("application", application)
	}
	if planType != "" {
		q.Set("plan_type", planType)
	}
	reqURL := fmt.Sprintf("%s/account/%s/plan/", c.HostURL, accountName)
	if encoded := q.Encode(); encoded != "" {
		reqURL += "?" + encoded
	}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := PlansList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	for i := range out.Results {
		if out.Results[i].Name == "" {
			out.Results[i].Name = out.Results[i].Plan
		}
	}
	return &out, nil
}

type UsageList struct {
	Year    int         `json:"year"`
	Month   int         `json:"month"`
	Results []UsageItem `json:"results"`
}

type UsageItem struct {
	StartDate     string   `json:"startDate"`
	EndDate       string   `json:"endDate"`
	ObjectID      string   `json:"objectID"`
	SKU           string   `json:"SKU"`
	Currency      string   `json:"currency"`
	Amount        string   `json:"amount"`
	Usage         int      `json:"usage"`
	TagCollection []string `json:"tagCollection"`
}

func (c *Client) GetUsage(accountName string, year, month int) (*UsageList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/usage/%d/%d/", c.HostURL, accountName, year, month), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := UsageList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if out.Year == 0 {
		out.Year = year
	}
	if out.Month == 0 {
		out.Month = month
	}
	return &out, nil
}

type UsageExtendedItem struct {
	StartDate          string   `json:"startDate"`
	EndDate            string   `json:"endDate"`
	ObjectID           string   `json:"objectID"`
	SKU                string   `json:"SKU"`
	Usage              int      `json:"usage"`
	TagCollection      []string `json:"tagCollection"`
	Currency           string   `json:"currency"`
	Amount             any      `json:"amount"`
	ExtendedAttributes []string `json:"extendedAttributes"`
}

type UsageExtendedList struct {
	Year    int                 `json:"year"`
	Month   int                 `json:"month"`
	Results []UsageExtendedItem `json:"results"`
}

func (c *Client) GetUsageExtended(accountName string, year, month int) (*UsageExtendedList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/usage-extended/%d/%d/", c.HostURL, accountName, year, month), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := UsageExtendedList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if out.Year == 0 {
		out.Year = year
	}
	if out.Month == 0 {
		out.Month = month
	}
	return &out, nil
}
