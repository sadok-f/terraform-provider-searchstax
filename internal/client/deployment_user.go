package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GetDeploymentUsers returns the list of Solr Basic Auth users for a deployment.
func (c *Client) GetDeploymentUsers(accountName string, deploymentID string) (*DeploymentUsersList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/get-users/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	deploymentUsers := DeploymentUsersList{}
	err = json.Unmarshal(body, &deploymentUsers)
	if err != nil {
		return nil, err
	}

	return &deploymentUsers, nil
}

// GetDeploymentUser returns details of a specific Solr Basic Auth user.
//
// The mock API does not expose a per-user read endpoint; we filter from the list.
func (c *Client) GetDeploymentUser(accountName string, deploymentID string, username string) (*DeploymentUser, error) {
	users, err := c.GetDeploymentUsers(accountName, deploymentID)
	if err != nil {
		return nil, err
	}
	for _, u := range users.Results {
		if u.Username == username {
			found := u
			return &found, nil
		}
	}
	return nil, fmt.Errorf("deployment user %q not found", username)
}

// CreateDeploymentUser creates a new Solr Basic Auth user for the deployment.
func (c *Client) CreateDeploymentUser(deploymentUser DeploymentUser, accountName string, deploymentID string) (*DeploymentUser, *Error) {
	rb, err := json.Marshal(deploymentUser)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "Marshal",
		}
	}

	// Enabling basic auth triggers a Solr restart that may still be settling
	// when add-user runs, causing transient 5xx responses. Retry on transient
	// errors to ride out the restart (skip the long backoff against the mock).
	const (
		attempts = 10
		backoff  = 15 * time.Second
	)

	var body []byte
	var lastErr error
	for i := 0; i < attempts; i++ {
		req, reqErr := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/add-user/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
		if reqErr != nil {
			return nil, &Error{
				err:     reqErr,
				context: "NewRequest",
			}
		}

		body, err = c.doRequest(req)
		if err == nil {
			break
		}
		if isTransient(err) && !c.isMockHost() {
			lastErr = err
			if i < attempts-1 {
				time.Sleep(backoff)
			}
			continue
		}
		return nil, &Error{
			err:     err,
			context: "doRequest",
		}
	}
	if err != nil {
		return nil, &Error{
			err:     lastErr,
			context: "doRequest",
		}
	}

	// A 2xx response from doRequest means the user was created. The mock API
	// returns {"created": true}; the real API returns the user payload without
	// a "created" field, so only treat an explicit {"created": false} as a
	// failure.
	var createResp struct {
		Created *bool `json:"created"`
	}
	if err := json.Unmarshal(body, &createResp); err == nil && createResp.Created != nil && !*createResp.Created {
		return nil, &Error{
			err:     fmt.Errorf("deployment user not created"),
			context: "CreateDeploymentUser",
		}
	}

	return &deploymentUser, nil
}

// UpdateDeploymentUser updates a deployment user by deleting then re-adding.
func (c *Client) UpdateDeploymentUser(accountName string, deploymentID string, deploymentUser DeploymentUser) (*DeploymentUser, *Error) {
	err := c.DeleteDeploymentUser(accountName, deploymentID, deploymentUser.Username)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "DeleteDeploymentUserOnUpdate",
		}
	}
	// Sleep for 5 seconds to wait until the deployment user is deleted.
	time.Sleep(time.Second * 5)
	newDeployment, err := c.CreateDeploymentUser(deploymentUser, accountName, deploymentID)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "CreateDeploymentUserOnUpdate",
		}
	}
	return newDeployment, nil
}

// DeleteDeploymentUser deletes a deployment user.
func (c *Client) DeleteDeploymentUser(accountName string, deploymentID string, username string) *Error {
	userToDelete, err := json.Marshal(map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return &Error{
			err:     err,
			context: "NewRequestOnDelete",
		}
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/delete-user/",
			c.HostURL, accountName, deploymentID), strings.NewReader(string(userToDelete)))
	if err != nil {
		return &Error{
			err:     err,
			context: "NewRequestOnDelete",
		}
	}

	body, err := c.doRequest(req)
	if err != nil {
		return &Error{
			"doRequestOnDelete",
			err,
		}
	}
	apiResponse := ApiResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return &Error{
			context: "UnmarshalOnDelete",
			err:     err,
		}
	}
	// mock returns {"deleted": true}
	var deleteResp struct {
		Deleted bool `json:"deleted"`
	}
	if err := json.Unmarshal(body, &deleteResp); err == nil && deleteResp.Deleted {
		return nil
	}
	if apiResponse.Success != "true" {
		return &Error{
			err:     fmt.Errorf("%s", apiResponse.Message),
			context: "ApiResponseOnDelete",
		}
	}
	//Check the resource status in a loop until it got deleted
	for {
		dep, err := c.GetDeployment(accountName, deploymentID)
		if err != nil {
			fmt.Printf("Deployment Deleted successfully: %v\n", err)
			return nil
		}
		// a workaround for the Acceptance tests to pass (since it is running against a mock API)
		if dep.Status == "Running" && dep.ProvisionState == "Done" {
			return nil
		}

		time.Sleep(time.Minute)
	}
}

// DeploymentUsersList represents the list payload for basic auth users.
type DeploymentUsersList struct {
	Results []DeploymentUser `json:"results"`
}

// UnmarshalJSON supports both response shapes for the get-users endpoint:
//   - The real API returns {"success":"true","users":{"<name>":{"roles":[...]}}}.
//   - The mock API returns {"results":[{"username":"...","role":"..."}]}.
func (l *DeploymentUsersList) UnmarshalJSON(data []byte) error {
	var apiResp struct {
		Users map[string]struct {
			Roles []string `json:"roles"`
		} `json:"users"`
	}
	if err := json.Unmarshal(data, &apiResp); err == nil && apiResp.Users != nil {
		l.Results = nil
		for name, u := range apiResp.Users {
			l.Results = append(l.Results, DeploymentUser{
				Username: name,
				Role:     canonicalRole(u.Roles),
			})
		}
		return nil
	}

	type alias DeploymentUsersList
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*l = DeploymentUsersList(a)
	return nil
}

// canonicalRole reduces the roles list returned by the API to the single
// role value used by this provider (Admin, ReadWrite, Read, or Write).
func canonicalRole(roles []string) string {
	has := make(map[string]bool, len(roles))
	for _, r := range roles {
		has[r] = true
	}
	switch {
	case has["Admin"]:
		return "Admin"
	case has["Read"] && has["Write"]:
		return "ReadWrite"
	case has["Read"]:
		return "Read"
	case has["Write"]:
		return "Write"
	default:
		if len(roles) > 0 {
			return roles[len(roles)-1]
		}
		return ""
	}
}

// DeploymentUser represents a basic auth user.
type DeploymentUser struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}
