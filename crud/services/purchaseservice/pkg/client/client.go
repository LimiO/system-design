package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"onlinestore/pkg/clients"
	"onlinestore/pkg/models"
	"onlinestore/services/purchaseservice/types"
)

type Client struct {
	clients.BaseClient
}

func NewClient(addr string) *Client {
	return &Client{
		BaseClient: clients.BaseClient{
			Addr: addr,
			C:    resty.New().SetAllowGetMethodPayload(true),
		},
	}
}

func (c *Client) GetOrder(username string, orderID int, header http.Header) (*models.Order, error) {
	url := fmt.Sprintf("http://%s/order", c.GetAddr())
	response := &types.GetOrderResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.GetOrderRequest{
			OrderID:  orderID,
			Username: username,
		}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Get(url)

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.Order, nil
}

func (c *Client) GetOrders(count int, header http.Header) ([]*models.Order, error) {
	url := fmt.Sprintf("http://%s/orders", c.GetAddr())
	response := &types.GetOrdersResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.GetOrdersRequest{Count: count}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Get(url)

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.Orders, nil
}

func (c *Client) Buy(count int, price int, productID int, header http.Header) (int, int, error) {
	url := fmt.Sprintf("http://%s/buy", c.GetAddr())
	response := &types.BuyResponse{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.BuyRequest{
			Count:     count,
			Price:     price,
			ProductID: productID,
		}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return 0, 0, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response.OrderID, response.Total, nil
}

func (c *Client) Commit(orderID int, status int, header http.Header) error {
	url := fmt.Sprintf("http://%s/commit", c.GetAddr())
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.CommitOrderRequest{
			OrderID: orderID,
			Status:  status,
		}).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}
