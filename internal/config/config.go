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
	Minio   MinioConfig `json:"minio"`
}

type DBConfig struct {
	DSN string `json:"dsn"`
}

type HttpService struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type MinioConfig struct {
	Endpoint        string `json:"endpoint"`
	PublicEndpoint  string `json:"public_endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

func Must(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read config file: %v", err))
	}

	c := Config{}

	if err := json.Unmarshal(b, &c); err != nil {
		panic(fmt.Errorf("unmarshall config: %v", err))
	}

	return c
}
