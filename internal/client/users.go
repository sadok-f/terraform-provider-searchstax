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
	// The real API wraps the list as {"success": true, "users": [...]}; the mock
	// historically used {"results": [...]}. Accept either.
	var wrapper struct {
		Users   []User `json:"users"`
		Results []User `json:"results"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, err
	}
	if len(wrapper.Users) > 0 {
		out.Results = wrapper.Users
	} else {
		out.Results = wrapper.Results
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
	// The real API returns {"success": true, "message": ...}; a non-2xx status is
	// already an error, so reaching here means the invite was sent.
	if _, err := c.doRequest(req); err != nil {
		return err
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
	// The real API returns {"success": true, "message": ...}; a non-2xx status is
	// already an error, so reaching here means the password was changed.
	if _, err := c.doRequest(req); err != nil {
		return err
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
	// The real API returns an array of per-user result objects; a non-2xx status
	// is already an error, so reaching here means the role was updated.
	if _, err := c.doRequest(req); err != nil {
		return err
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
	// The real API returns {"id": ..., "success": true, "message": ...}; a
	// non-2xx status is already an error, so reaching here means the user was
	// deleted.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}
