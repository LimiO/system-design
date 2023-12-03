package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Read[T interface{}](path string) (*T, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %v", path, err)
	}
	config := new(T)
	if err = yaml.Unmarshal(content, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data from path %q: %v", path, err)
	}
	return config, nil
}
