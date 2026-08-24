package config

import (
	"encoding/json"
	"os"
)

type Client struct {
	ID   string `json:"id"`
	Flow string `json:"flow"`
}

type Config struct {
	Listen      string   `json:"listen"`
	RealitySock string   `json:"reality_sock"`
	Clients     []Client `json:"clients"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
