package config

import (
	"encoding/json"
	"os"
)

func InitFromJson[T any](configFilePath string) (T, error) {
	var config T

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configData, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
