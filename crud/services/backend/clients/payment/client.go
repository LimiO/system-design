package payment

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

func (c *Client) AddBalance(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("POST", data, fmt.Sprintf("http://%s/balance/add", c.GetAddr()), header)
}

func (c *Client) SubBalance(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("POST", data, fmt.Sprintf("http://%s/balance/sub", c.GetAddr()), header)
}

func (c *Client) GetBalance(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("GET", data, fmt.Sprintf("http://%s/balance", c.GetAddr()), header)
}
