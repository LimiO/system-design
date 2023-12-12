package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"onlinestore/pkg/clients"
	"onlinestore/services/paymentservice/types"
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

func (c *Client) AddBalance(username string, amount int, header http.Header) error {
	url := fmt.Sprintf("http://%s/balance/add", c.GetAddr())
	response := &types.AddBalanceResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.AddBalanceRequest{
			Username: username,
			Amount:   amount,
		}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}

func (c *Client) SubBalance(username string, amount int, header http.Header) error {
	url := fmt.Sprintf("http://%s/balance/sub", c.GetAddr())
	response := &types.SubBalanceResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.AddBalanceRequest{
			Username: username,
			Amount:   amount}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}

func (c *Client) Commit(reserveID int, status int, header http.Header) error {
	url := fmt.Sprintf("http://%s/balance/commit", c.GetAddr())
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

func (c *Client) ReserveBalance(username string, amount int, header http.Header) (int, error) {
	url := fmt.Sprintf("http://%s/balance/reserve", c.GetAddr())
	response := &types.ReserveBalanceResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.ReserveBalanceRequest{
			Username: username,
			Amount:   amount,
		}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.ReserveID, nil
}

func (c *Client) GetBalance(username string, header http.Header) (int, error) {
	url := fmt.Sprintf("http://%s/balance", c.GetAddr())
	response := &types.GetBalanceResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.GetBalanceRequest{
			Username: username,
		}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Get(url)

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.Balance, nil
}
