package purchases

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

func (c *Client) GetOrder(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("GET", data, fmt.Sprintf("http://%s/order", c.GetAddr()), header)
}

func (c *Client) GetOrders(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("GET", data, fmt.Sprintf("http://%s/orders", c.GetAddr()), header)
}

func (c *Client) Buy(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("POST", data, fmt.Sprintf("http://%s/buy", c.GetAddr()), header)
}

func (c *Client) Commit(data []byte, header http.Header) (*http.Response, error) {
	return c.DoRequest("POST", data, fmt.Sprintf("http://%s/commit", c.GetAddr()), header)
}
