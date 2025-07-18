package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github-runner-manager/internal/types"
)

// Load loads configuration from file or returns default
func Load(filename string) (types.Config, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return getDefaultConfig(), nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return types.Config{}, fmt.Errorf("reading config: %w", err)
	}

	var config types.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return types.Config{}, fmt.Errorf("parsing config: %w", err)
	}

	return config, nil
}

// Save saves configuration to file
func Save(filename string, config types.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}
