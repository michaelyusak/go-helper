package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func InitFromJson[T any](log *logrus.Logger, configFilePath string) T {
	logHeading := "[go-helper][config][InitFromJson]"

	var config T

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%s[os.ReadFile] error: %s", logHeading, err.Error()),
		}).Fatal("error initiating config file")

		return config
	}

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%s[json.Unmarshal] error: %s", logHeading, err.Error()),
		}).Fatal("error initiating config file")

		return config
	}

	return config
}
