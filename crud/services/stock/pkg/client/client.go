package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"onlinestore/pkg/clients"
	"onlinestore/services/stock/types"
)

type Client struct {
	clients.BaseClient
}

func NewClient(addr string) *Client {
	client := &Client{
		BaseClient: clients.BaseClient{
			Addr: addr,
			C:    resty.New(),
		},
	}
	client.C.SetAllowGetMethodPayload(true)
	return client
}

func (c *Client) ReserveCount(productID int, count int, header http.Header) (int, error) {
	url := fmt.Sprintf("http://%s/products/reserve", c.GetAddr())
	response := &types.ReserveResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.ReserveRequest{
			ProductID: productID,
			Count:     count,
		}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return response.ReserveID, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.ReserveID, nil
}

func (c *Client) Commit(reserveID int, status int, header http.Header) error {
	url := fmt.Sprintf("http://%s/products/commit", c.GetAddr())
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.CommitRequest{
			ReserveID: reserveID,
			Status:    status,
		}).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}

func (c *Client) AddCount(productID int, count int, header http.Header) error {
	url := fmt.Sprintf("http://%s/products/add", c.GetAddr())
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.AddCountRequest{
			ProductID: productID,
			Count:     count,
		}).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}

func (c *Client) GetCount(productID int, header http.Header) (int, error) {
	url := fmt.Sprintf("http://%s/product", c.GetAddr())
	response := &types.GetCountResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.GetCountRequest{
			ProductID: productID,
		}).
		SetResult(response).
		SetHeader("Authorization", header.Get("Authorization")).
		Get(url)

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.Count, nil
}
