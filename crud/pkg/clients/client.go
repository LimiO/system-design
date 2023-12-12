package clients

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type BaseClient struct {
	C    *resty.Client
	Addr string
}

func (b *BaseClient) GetAddr() string {
	return b.Addr
}

func (b *BaseClient) DoRequest(method string, data []byte, addr string, headers http.Header) (*http.Response, error) {
	dst := make([]byte, len(data))
	copy(dst, data)
	responseBody := bytes.NewBuffer(dst)
	req, err := http.NewRequest(method, addr, responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key := range headers {
		req.Header.Set(key, headers.Get(key))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	return resp, nil
}
