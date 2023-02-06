package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func (c Config) DatabaseConnString() string {
	return fmt.Sprintf(`postgres://%s:%s@%s:%d`,
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}

func New(filename string) (*Config, error) {
	f, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = json.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
