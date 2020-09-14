package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/quentin-fox/fox"
)

func GetConfig() (fox.Config, error) {
	var config fox.Config
	env := os.Getenv("GO_ENV")

	if env == "" {
		env = "default"
	}

	filename := getFileName(env)
	exists := doesFileExist(filename)

	if !exists {
		return config, errors.New("Config file does not exist: " + filename)
	}

	err := readIntoConfig(filename, &config)

	return config, err
}

func getFileName(env string) string {
	return "config/" + env + ".json"
}

func doesFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func readIntoConfig(filename string, config *fox.Config) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)

	return err
}




