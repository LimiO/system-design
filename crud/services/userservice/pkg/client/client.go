package client

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"

	"onlinestore/pkg/clients"
	"onlinestore/pkg/models"
	"onlinestore/services/userservice/types"
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

func convertUser(user *models.User) map[string]string {
	m := map[string]string{
		"username":   user.Username,
		"email":      user.Email,
		"phone":      strconv.Itoa(user.Phone),
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}
	if user.Password != "" {
		m["password"] = user.Password
	}
	return m
}

func (c *Client) GetUser(username string, header http.Header) (*models.User, error) {
	url := fmt.Sprintf("http://%s/user", c.GetAddr())
	response := &models.User{}
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(map[string]string{
			"username": username,
		}).
		SetResult(response).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Get(url)

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return response, nil
}

func (c *Client) PutUser(user *models.User, header http.Header) error {
	url := fmt.Sprintf("http://%s/user", c.GetAddr())
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(user).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Put(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}

func (c *Client) PostUser(user *models.User, header http.Header) error {
	url := fmt.Sprintf("http://%s/user", c.GetAddr())
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(user).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Post(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}

func (c *Client) DeleteUser(username string, header http.Header) error {
	url := fmt.Sprintf("http://%s/user", c.GetAddr())
	resp, err := c.C.R().ForceContentType("application/json").
		SetBody(types.DeleteUser{
			Username: username,
		}).
		SetHeaders(map[string]string{"Authorization": header.Get("Authorization")}).
		Delete(url)

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request %d: %v", resp.StatusCode(), err)
	}

	return nil
}
