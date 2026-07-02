package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type BackupsList struct {
	Results []Backup `json:"results"`
}

type Backup struct {
	ID string `json:"id,omitempty"`
}

func (c *Client) GetAccountBackups(accountName string) (*BackupsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/backup/", c.HostURL, accountName), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := BackupsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteAccountBackup(accountName, backupUID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/backup/%s/", c.HostURL, accountName, backupUID), nil)
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
		return fmt.Errorf("account backup not deleted")
	}
	return nil
}

type RestoreRequest struct {
	BackupID string `json:"backup_id,omitempty"`
}

type RestoreResponse struct {
	RestoreID string `json:"restore_id"`
	Queued    bool   `json:"queued"`
	Status    string `json:"status,omitempty"`
}

func (c *Client) CreateAccountRestore(accountName string, reqBody RestoreRequest) (*RestoreResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/restore/", c.HostURL, accountName), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := RestoreResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if !out.Queued {
		return nil, fmt.Errorf("restore not queued")
	}
	return &out, nil
}

func (c *Client) GetDeploymentBackups(accountName, deploymentID string) (*BackupsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/backup/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := BackupsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type CreateBackupResponse struct {
	BackupID string `json:"backup_id"`
	Queued   bool   `json:"queued"`
}

func (c *Client) CreateDeploymentBackup(accountName, deploymentID string, reqBody map[string]any) (*CreateBackupResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/backup/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := CreateBackupResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if !out.Queued {
		return nil, fmt.Errorf("deployment backup not queued")
	}
	return &out, nil
}

func (c *Client) DeleteDeploymentBackup(accountName, deploymentID, backupUID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/backup/%s/", c.HostURL, accountName, deploymentID, backupUID), nil)
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
		return fmt.Errorf("deployment backup not deleted")
	}
	return nil
}

type BackupSchedulesList struct {
	Results []BackupSchedule `json:"results"`
}

type BackupSchedule struct {
	ID          string   `json:"id,omitempty"`
	Days        []string `json:"days,omitempty"`
	Time        string   `json:"time,omitempty"`
	Retention   int      `json:"retention,omitempty"`
	Frequency   int      `json:"frequency,omitempty"`
	RegionID    string   `json:"region_id,omitempty"`
	Collections []string `json:"collections,omitempty"`
}

func (c *Client) GetBackupSchedules(accountName, deploymentID string) (*BackupSchedulesList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/backup/schedule/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := BackupSchedulesList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateBackupSchedule(accountName, deploymentID string, reqBody map[string]any) (*BackupSchedule, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/backup/schedule/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	// A non-2xx response is already surfaced as an error by doRequest, so a
	// successful call here means the schedule was created. The API returns the
	// created schedule object rather than a {"created": true} flag.
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := BackupSchedule{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteBackupSchedule(accountName, deploymentID, scheduleUID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/backup/schedule/%s/", c.HostURL, accountName, deploymentID, scheduleUID), nil)
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
		return fmt.Errorf("backup schedule not deleted")
	}
	return nil
}

func (c *Client) CreateDeploymentRestore(accountName, deploymentID string, reqBody RestoreRequest) (*RestoreResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/restore/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := RestoreResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if !out.Queued {
		return nil, fmt.Errorf("deployment restore not queued")
	}
	return &out, nil
}

func (c *Client) GetDeploymentRestoreStatus(accountName, deploymentID string, reqBody RestoreRequest) (*RestoreResponse, error) {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/restore/status/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := RestoreResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
