package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json: "db_url"`
	CurrentUserName string `json: "current_user_name"`
}

const configFileName = "/.gatorconfig.json"

func getConfigPath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path += configFileName
	return path, nil
}

func Read() (Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer data.Close()

	decoder := json.NewDecoder(data)
	Cfg := Config{}

	err = decoder.Decode(&Cfg)
	if err != nil {
		return Config{}, err
	}

	return Cfg, nil

}

func write(cfg Config) error {

	fullPath, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetUser(usr string) error {
	c.CurrentUserName = usr
	return write(*c)
}
