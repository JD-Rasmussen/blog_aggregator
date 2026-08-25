package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `"postgres://example"`
	Current_user_name string `json:"current_user_name"`
}

//func (c *Config) GetConfig() (string, string) {
//	return c.Db_url, c.Current_user_name
//}

func (c *Config) SetUser(name string) {
	c.Current_user_name = name
	fileDir, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Error getting config file path:", err)
		return
	}
	// Write the updated config to the file
	data, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Error marshaling config:", err)
		return
	}
	err = os.WriteFile(fileDir, data, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting user home directory: %v", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {

	fileDir, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Error getting config file path: %v", err)
	}

	data, err := os.ReadFile(fileDir)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading config file: %v", err)
	}
	cfg := Config{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing config file: %v", err)
	}
	return cfg, nil
}
