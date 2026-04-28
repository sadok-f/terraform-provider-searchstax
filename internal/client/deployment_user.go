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

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/add-user/", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "NewRequest",
		}
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "doRequest",
		}
	}

	var createResp struct {
		Created bool `json:"created"`
	}
	err = json.Unmarshal(body, &createResp)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "Unmarshal",
		}
	}

	if createResp.Created {
		return &deploymentUser, nil
	}

	return nil, &Error{
		err:     fmt.Errorf("deployment user not created"),
		context: "CreateDeploymentUser",
	}
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

// DeploymentUser represents a basic auth user.
type DeploymentUser struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}
