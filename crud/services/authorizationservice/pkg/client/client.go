package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"onlinestore/pkg/clients"
	"onlinestore/services/authorizationservice/types"
)

type Client struct {
	clients.BaseClient
}

func NewClient(addr string) *Client {
	return &Client{
		BaseClient: clients.BaseClient{
			Addr: addr,
			C:    resty.New(),
		},
	}
}

func (c *Client) GetToken(username string, password string, header http.Header) (string, error) {
	url := fmt.Sprintf("http://%s/token", c.GetAddr())
	response := &types.TokenResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.TokenRequest{
			Username: username,
			Password: password}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Get(url)

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.Token, nil
}

func (c *Client) Register(username string, password string, header http.Header) (string, error) {
	header.Set("Content-Type", "application/json")

	url := fmt.Sprintf("http://%s/register", c.GetAddr())
	response := &types.TokenResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.TokenRequest{
			Username: username,
			Password: password}).
		SetResult(response).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.Token, nil
}

func (c *Client) Unregister(username string, password string, header http.Header) error {
	url := fmt.Sprintf("http://%s/unregister", c.GetAddr())
	response := &types.DeleteTokenResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.TokenRequest{
			Username: username,
			Password: password,
		}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}
