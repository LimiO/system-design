package user

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

func (c *Client) GetUser(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("GET", data, fmt.Sprintf("http://%s/user", c.GetAddr()), header)
}

func (c *Client) PutUser(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("PUT", data, fmt.Sprintf("http://%s/user", c.GetAddr()), header)
}

func (c *Client) DeleteUser(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("DELETE", data, fmt.Sprintf("http://%s/user", c.GetAddr()), header)
}

func (c *Client) PostUser(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("POST", data, fmt.Sprintf("http://%s/user", c.GetAddr()), header)
}
