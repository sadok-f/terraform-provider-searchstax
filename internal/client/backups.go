package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// coerceID converts a JSON id that may be a string or a number into a string.
// The SearchStax API returns numeric ids for backups and schedules, while the
// acceptance-test mock returns strings.
func coerceID(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return ""
	}
	return strings.Trim(s, `"`)
}

// decodeResults unmarshals a list response that is either a bare JSON array of
// items (real API) or a {"results": [...]} wrapper (mock API) into results,
// which must be a pointer to a slice.
func decodeResults(body []byte, results any) error {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		if err := json.Unmarshal(body, results); err != nil {
			return fmt.Errorf("%w (body: %s)", err, truncate(body, 512))
		}
		return nil
	}
	wrapper := struct {
		Results json.RawMessage `json:"results"`
	}{}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return fmt.Errorf("%w (body: %s)", err, truncate(body, 512))
	}
	if len(bytes.TrimSpace(wrapper.Results)) == 0 {
		return nil
	}
	if err := json.Unmarshal(wrapper.Results, results); err != nil {
		return fmt.Errorf("%w (body: %s)", err, truncate(body, 512))
	}
	return nil
}

type BackupsList struct {
	Results []Backup `json:"results"`
}

type Backup struct {
	ID string `json:"id,omitempty"`
}

// UnmarshalJSON tolerates an "id" returned as either a JSON string or a number.
func (b *Backup) UnmarshalJSON(data []byte) error {
	type alias Backup
	aux := &struct {
		ID json.RawMessage `json:"id,omitempty"`
		*alias
	}{alias: (*alias)(b)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	b.ID = coerceID(aux.ID)
	return nil
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
	if err := decodeResults(body, &out.Results); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteAccountBackup(accountName, backupUID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/backup/%s/", c.HostURL, accountName, backupUID), nil)
	if err != nil {
		return err
	}
	// The real API returns 204 with an empty body; the mock returns
	// {"deleted": true}. doRequest already errors on any non-2xx status, so
	// reaching here means the backup was deleted.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

type RestoreRequest struct {
	BackupID string `json:"backup_id,omitempty"`
}

// RestoreResponse mirrors the real SearchStax API, which confirms restore
// create and reports restore status with a single "message" string.
type RestoreResponse struct {
	Message string `json:"message,omitempty"`
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
	// The API confirms with {"message": "restore begun"}. A non-2xx status is
	// already surfaced as an error by doRequest, so reaching here means the
	// restore was accepted.
	out := RestoreResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
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
	if err := decodeResults(body, &out.Results); err != nil {
		return nil, err
	}
	return &out, nil
}

type CreateBackupResponse struct {
	BackupID string `json:"backup_id"`
}

// UnmarshalJSON reads the new backup id from either "backup_id" (mock) or "id"
// (real API), tolerating a string or numeric value.
func (r *CreateBackupResponse) UnmarshalJSON(data []byte) error {
	var aux struct {
		BackupID json.RawMessage `json:"backup_id,omitempty"`
		ID       json.RawMessage `json:"id,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if id := coerceID(aux.BackupID); id != "" {
		r.BackupID = id
	} else {
		r.BackupID = coerceID(aux.ID)
	}
	return nil
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
	return &out, nil
}

func (c *Client) DeleteDeploymentBackup(accountName, deploymentID, backupUID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/backup/%s/", c.HostURL, accountName, deploymentID, backupUID), nil)
	if err != nil {
		return err
	}
	// The real API confirms with {"message": ..., "success": "true"} while the
	// mock returns {"deleted": true}. A non-2xx status is already an error, so
	// reaching here means the backup was deleted.
	if _, err := c.doRequest(req); err != nil {
		return err
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

// UnmarshalJSON tolerates an "id" returned as either a JSON string or a number,
// which the SearchStax API uses interchangeably for schedule identifiers.
func (b *BackupSchedule) UnmarshalJSON(data []byte) error {
	type alias BackupSchedule
	aux := &struct {
		ID json.RawMessage `json:"id,omitempty"`
		*alias
	}{alias: (*alias)(b)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	b.ID = coerceID(aux.ID)
	return nil
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
	if err := decodeResults(body, &out.Results); err != nil {
		return nil, err
	}
	return &out, nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}

func (c *Client) CreateBackupSchedule(accountName, deploymentID string, reqBody map[string]any) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/backup/schedule/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	// The API responds with {"message": "Backup Scheduled Successfully"} and does
	// not echo the schedule id. A non-2xx status is already surfaced as an error
	// by doRequest, so reaching here means the schedule was created.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteBackupSchedule(accountName, deploymentID, scheduleUID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/backup/schedule/%s/", c.HostURL, accountName, deploymentID, scheduleUID), nil)
	if err != nil {
		return err
	}
	// A successful delete returns {"success": "The backup schedule has been
	// deleted!"}. doRequest already turns any non-2xx status into an error, so
	// reaching here without error means the schedule was deleted.
	if _, err := c.doRequest(req); err != nil {
		return err
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
	// The API confirms with a {"message": ...} body. A non-2xx status is already
	// an error, so reaching here means the restore was accepted.
	out := RestoreResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
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
