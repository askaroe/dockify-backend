package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Host string `json:"host" envconfig:"host"`
	Port string `json:"port" envconfig:"port"`
	PostgresConfig
	MindsporeModelURL string `json:"mindspore_model_url" envconfig:"mindspore_model_url"`
}

type PostgresConfig struct {
	DbHost     string `json:"db_host" envconfig:"db_host"`
	DbName     string `json:"db_name" envconfig:"db_name"`
	DbPort     string `json:"db_port" envconfig:"db_port"`
	DbUsername string `json:"db_username" envconfig:"db_username"`
	DbPassword string `json:"db_password" envconfig:"db_password"`
	DbSslmode  string `json:"db_sslmode" envconfig:"db_sslmode"`
}

func getConfigsFromJSON() (*Config, error) {
	var filePath string
	if os.Getenv("config") == "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		filePath = pwd + "/config/config.json"
	} else {
		filePath = os.Getenv("config")
	}

	file, err := os.Open(filePath)

	if err != nil {
		return &Config{}, err
	}

	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)

	if err != nil {
		return &Config{}, err
	}

	return &config, err
}

func GetConfig() (*Config, error) {
	return getConfigsFromJSON()
}
