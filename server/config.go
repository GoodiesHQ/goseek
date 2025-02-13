package server

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type GoSeekConfig struct {
	Root    string   `yaml:"root"`
	Port    uint16   `yaml:"port"`
	ApiKeys []string `yaml:"apikeys"`
}

func (config *GoSeekConfig) LoadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", filePath, err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return nil
}
