package web

import (
	"encoding/json"
	"fmt"
	"io"
)

func DecodeHttpBody[T interface{}](body io.ReadCloser) (*T, error) {
	t := new(T)
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(t); err != nil {
		return nil, fmt.Errorf("failed to read data: %v", err)
	}
	return t, nil
}
