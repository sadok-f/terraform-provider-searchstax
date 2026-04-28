package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhooksList struct {
	Results []Webhook `json:"results"`
}

type Webhook struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Paused bool   `json:"paused"`
}

func (c *Client) GetWebhooks(accountName string) (*WebhooksList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/webhook/", c.HostURL, accountName), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := WebhooksList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
