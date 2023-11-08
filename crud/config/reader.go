package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func Read(path string, config interface{}) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %v", path, err)
	}
	if err = yaml.Unmarshal(content, config); err != nil {
		return fmt.Errorf("failed to unmarshal data from path %q: %v", path, err)
	}
	return nil
}
