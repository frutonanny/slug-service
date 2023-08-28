package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	Arg = "config"
)

type Config struct {
	DB      DBConfig    `json:"db"`
	Service HttpService `json:"service"`
}

type DBConfig struct {
	DSN string `json:"dsn"`
}

type HttpService struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

func Must(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read config file: %w", err))
	}

	c := Config{}

	if err := json.Unmarshal(b, &c); err != nil {
		panic(fmt.Errorf("unmarshall config: %w", err))
	}

	return c
}
