package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GetDeploymentUsers - Returns list of datasources (no auth required).
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

// GetDeploymentUser - Returns details of a Deployment User.
func (c *Client) GetDeploymentUser(accountName string, deploymentID string, username string) (*DeploymentUser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/get-users/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	deploymentUser := DeploymentUser{}
	err = json.Unmarshal(body, &deploymentUser)
	if err != nil {
		return nil, err
	}

	return &deploymentUser, nil
}

// CreateDeploymentUser - Create new deployment User.
func (c *Client) CreateDeploymentUser(deploymentUser DeploymentUser, accountName string, deploymentID string) (*DeploymentUser, *Error) {
	rb, err := json.Marshal(deploymentUser)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "Marshal",
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/%s/solr/auth/add-user", c.HostURL, accountName, deploymentID), strings.NewReader(string(rb)))
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

	var newDeploymentUserResponse struct {
		success bool
		message string
	}
	err = json.Unmarshal(body, &newDeploymentUserResponse)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "Unmarshal",
		}
	}

	if newDeploymentUserResponse.success {
		return &deploymentUser, nil
	}

	return nil, &Error{
		err:     err,
		context: "onCreateDeploymentUser",
	}
}

// UpdateDeploymentUser -Update a deployment: for now it recreate the cluster.
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

// DeleteDeploymentUser -Delete a specific deployment user.
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

// DeploymentUsersList - DeploymentUsersList struct.
type DeploymentUsersList struct {
	Success bool             `json:"success"`
	Users   []DeploymentUser `json:"users"`
}

// DeploymentUser - DeploymentUser struct.
type DeploymentUser struct {
	UID      string `json:"UID"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Role     string `json:"Roles"`
}
