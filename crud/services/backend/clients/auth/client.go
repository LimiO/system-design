package auth

import (
	"fmt"
	"net/http"

	"onlinestore/services/backend/clients"
)

type Client struct {
	clients.BaseClient
}

func NewClient(addr string) *Client {
	return &Client{
		BaseClient: clients.BaseClient{
			Addr: addr,
		},
	}
}

func (c *Client) GetToken(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("GET", data, fmt.Sprintf("http://%s/token", c.GetAddr()), header)
}

func (c *Client) Register(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("POST", data, fmt.Sprintf("http://%s/register", c.GetAddr()), header)
}

func (c *Client) Unregister(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("POST", data, fmt.Sprintf("http://%s/unregister", c.GetAddr()), header)
}
