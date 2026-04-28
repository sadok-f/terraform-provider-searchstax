package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type UsersList struct {
	Results []User `json:"results"`
}

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (c *Client) GetUsers() (*UsersList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/", c.RestHostURL()), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	out := UsersList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type InviteUserRequest struct {
	Email     string `json:"email"`
	Role      string `json:"role"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func (c *Client) InviteUser(reqBody InviteUserRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/add-user", c.RestHostURL()), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Invited bool `json:"invited"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Invited {
		return fmt.Errorf("user invite failed")
	}
	return nil
}

type ChangeUserPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

func (c *Client) ChangeUserPassword(reqBody ChangeUserPasswordRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/change-password", c.RestHostURL()), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Changed bool `json:"changed"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Changed {
		return fmt.Errorf("password change failed")
	}
	return nil
}

type SetUserRoleRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (c *Client) SetUserRole(reqBody SetUserRoleRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/set-role/", c.RestHostURL()), strings.NewReader(string(rb)))
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
		return fmt.Errorf("set role failed")
	}
	return nil
}

type DeleteUserRequest struct {
	Email string `json:"email"`
}

func (c *Client) DeleteUser(reqBody DeleteUserRequest) error {
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/", c.RestHostURL()), strings.NewReader(string(rb)))
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
		return fmt.Errorf("delete user failed")
	}
	return nil
}
