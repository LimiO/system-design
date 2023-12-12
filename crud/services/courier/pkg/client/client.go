package client

import (
	"fmt"
	"net/http"
	"onlinestore/pkg/models"

	"github.com/go-resty/resty/v2"

	"onlinestore/pkg/clients"
	"onlinestore/services/courier/types"
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

func (c *Client) ReserveCourier(header http.Header) error {
	url := fmt.Sprintf("http://%s/courier/reserve", c.GetAddr())
	response := &types.ReserveCourierResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		SetResult(response).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to send request %d: %v", resp.StatusCode(), err)
	}

	return nil
}

func (c *Client) UnreserveCourier(courier string, header http.Header) error {
	url := fmt.Sprintf("http://%s/courier/unreserve", c.GetAddr())
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.UnreserveCourierRequest{
			Username: courier,
		}).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to send request %d: %v", resp.StatusCode(), err)
	}

	return nil
}

func (c *Client) GetCourier(username string, header http.Header) (*models.Courier, error) {
	url := fmt.Sprintf("http://%s/courier", c.GetAddr())
	response := &models.Courier{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.GetCourierRequest{
			Username: username,
		}).
		SetResult(response).
		SetHeader("Authorization", header.Get("Authorization")).
		Get(url)

	if resp.StatusCode() != http.StatusOK {
		return response, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response, nil
}
