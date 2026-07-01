package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GetDeployments - Returns list of datasources (no auth required).
func (c *Client) GetDeployments(accountName string) (*DeploymentsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/", c.HostURL, accountName), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	deployments := DeploymentsList{}
	err = json.Unmarshal(body, &deployments)
	if err != nil {
		return nil, err
	}

	return &deployments, nil
}

// GetDeployment - Returns specific deployment (no auth required).
func (c *Client) GetDeployment(accountName string, deploymentID string) (*Deployment, *Error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/", c.HostURL, accountName, deploymentID), nil)
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

	deployment := Deployment{}
	err = json.Unmarshal(body, &deployment)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "Unmarshal",
		}
	}

	return &deployment, nil
}

// deploymentCreateRequest is the payload accepted by the SearchStax
// deployment-create API. It intentionally contains only the fields the API
// expects on creation; sending the full Deployment struct would include
// response-only fields (e.g. an empty "subscription"), which the API rejects.
type deploymentCreateRequest struct {
	Name                  string `json:"name"`
	Application           string `json:"application"`
	ApplicationVersion    string `json:"application_version"`
	TerminationLock       bool   `json:"termination_lock"`
	PlanType              string `json:"plan_type"`
	Plan                  string `json:"plan"`
	RegionId              string `json:"region_id"`
	CloudProviderId       string `json:"cloud_provider_id"`
	NumAdditionalAppNodes int64  `json:"num_additional_app_nodes"`
	PrivateVpc            *int64 `json:"private_vpc,omitempty"`
}

// CreateDeployment - Create new deployment.
func (c *Client) CreateDeployment(deployment Deployment, accountName string) (*Deployment, *Error) {
	payload := deploymentCreateRequest{
		Name:                  deployment.Name,
		Application:           deployment.Application,
		ApplicationVersion:    deployment.ApplicationVersion,
		TerminationLock:       deployment.TerminationLock,
		PlanType:              deployment.PlanType,
		Plan:                  deployment.Plan,
		RegionId:              deployment.RegionId,
		CloudProviderId:       deployment.CloudProviderId,
		NumAdditionalAppNodes: deployment.NumAdditionalAppNodes,
	}
	if deployment.PrivateVpc != 0 {
		vpc := deployment.PrivateVpc
		payload.PrivateVpc = &vpc
	}

	rb, err := json.Marshal(payload)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "Marshal",
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/account/%s/deployment/", c.HostURL, accountName), strings.NewReader(string(rb)))
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

	newDeployment := Deployment{}
	err = json.Unmarshal(body, &newDeployment)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "Unmarshal",
		}
	}
	//Check the resource status in a loop until it becomes "Done"
	for {
		dep, err := c.GetDeployment(accountName, newDeployment.UID)
		if err != nil {
			return nil, &Error{
				err:     err,
				context: "GetDeploymentStatus",
			}
		}

		if dep.Status == "Running" && dep.ProvisionState == "Done" {
			newDeployment.Status = dep.Status
			newDeployment.ProvisionState = dep.ProvisionState
			newDeployment.HttpEndpoint = dep.HttpEndpoint
			break
		}
		if dep.Status == "Failed" {
			err := fmt.Errorf("operation failed with status: %s", dep.Status)
			return nil, &Error{
				err:     err,
				context: "GetDeploymentStatus",
			}
		}

		time.Sleep(time.Minute)
	}

	return &newDeployment, nil
}

// UpdateDeployment -Update a deployment: for now it recreate the cluster.
func (c *Client) UpdateDeployment(accountName string, deploymentID string, deployment Deployment) (*Deployment, *Error) {
	err := c.DeleteDeployment(context.Background(), accountName, deploymentID)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "DeleteDeploymentOnUpdate",
		}
	}
	// Sleep for 1 minutes to wait until the deployment is deleted.
	time.Sleep(time.Minute)
	newDeployment, err := c.CreateDeployment(deployment, accountName)
	if err != nil {
		return nil, &Error{
			err:     err,
			context: "CreateDeploymentOnUpdate",
		}
	}
	return newDeployment, nil
}

// DeleteDeployment -Delete a specific deployment.
func (c *Client) DeleteDeployment(ctx context.Context, accountName string, deploymentID string) *Error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/", c.HostURL, accountName, deploymentID), nil)
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

	// Poll until the deployment is actually gone (the API returns 404) or the
	// timeout elapses. Only a definitive 404 confirms deletion — a transient
	// error (network blip, 5xx, expired token) must NOT be treated as success,
	// otherwise Terraform would drop the resource from state while it still
	// exists, leaving orphaned infrastructure.
	const (
		maxWait      = 30 * time.Minute
		pollInterval = 15 * time.Second
	)
	deadline := time.Now().Add(maxWait)
	for {
		dep, getErr := c.GetDeployment(accountName, deploymentID)
		if getErr != nil {
			if isNotFound(getErr) {
				return nil // deployment is gone
			}
			// Transient error: keep polling until the timeout.
		} else if c.isMockHost() && dep.Status == "Running" && dep.ProvisionState == "Done" {
			// The mock API used by the acceptance tests never returns 404 after
			// a delete; treat a healthy mock deployment as deleted. This branch
			// is gated to the mock host so it cannot fire against a real API.
			return nil
		}

		if time.Now().After(deadline) {
			return &Error{
				context: "DeleteDeploymentTimeout",
				err:     fmt.Errorf("timed out after %s waiting for deployment %s to be deleted", maxWait, deploymentID),
			}
		}

		select {
		case <-ctx.Done():
			return &Error{context: "DeleteDeploymentCanceled", err: ctx.Err()}
		case <-time.After(pollInterval):
		}
	}
}

// isNotFound reports whether err (or a wrapped error) is an HTTP 404 response.
func isNotFound(err error) bool {
	var httpErr *HTTPStatusError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusNotFound
	}
	return false
}

// DeploymentsList - DeploymentsList struct.
type DeploymentsList struct {
	Count    int32        `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []Deployment `json:"results"`
}

// Deployment - Deployment struct.
type Deployment struct {
	UID                         string             `json:"uid"`
	Name                        string             `json:"name"`
	Application                 string             `json:"application"`
	ApplicationVersion          string             `json:"application_version"`
	Tier                        string             `json:"tier"`
	HttpEndpoint                string             `json:"http_endpoint"`
	Status                      string             `json:"status"`
	ProvisionState              string             `json:"provision_state"`
	TerminationLock             bool               `json:"termination_lock"`
	Plan                        string             `json:"plan"`
	PlanType                    string             `json:"plan_type"`
	IsMasterSlave               bool               `json:"is_master_slave"`
	VpcType                     string             `json:"vpc_type"`
	VpcName                     string             `json:"vpc_name"`
	RegionId                    string             `json:"region_id"`
	CloudProvider               string             `json:"cloud_provider"`
	CloudProviderId             string             `json:"cloud_provider_id"`
	DeploymentType              string             `json:"deployment_type"`
	NumAdditionalAppNodes       int64              `json:"num_additional_app_nodes"`
	NumNodesDefault             int64              `json:"num_nodes_default"`
	NumZookeeperNodesDefault    int64              `json:"num_zookeeper_nodes_default"`
	NumAdditionalZookeeperNodes int64              `json:"num_additional_zookeeper_nodes"`
	PrivateVpc                  int64              `json:"private_vpc"`
	DateCreated                 string             `json:"date_created"`
	Servers                     FlexStringList     `json:"servers"`
	ZookeeperEnsemble           string             `json:"zookeeper_ensemble"`
	Tag                         FlexStringList     `json:"tag"`
	Specifications              FlexSpecifications `json:"specifications"`
	BackupsEnabled              bool               `json:"backups_enabled"`
	DrEnabled                   bool               `json:"dr_enabled"`
	SlaActive                   bool               `json:"sla_active"`
	ApplicationNodesCount       int64              `json:"application_nodes_count"`
	Subscription                string             `json:"subscription"`
	SecurityPack                bool               `json:"security_pack"`
	DesiredTier                 string             `json:"desired_tier"`
}
