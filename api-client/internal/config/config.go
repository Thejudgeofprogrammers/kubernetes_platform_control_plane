package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	BaseURL      string
	TimeoutMs    int
	RetryCount   int
	RetryBackoff int
	AuthType     string
	AuthRef      string

	ClientID string
}

func LoadEnv() *Config {
	config := &Config{
		BaseURL:      getenv("BASE_URL", ""),
		TimeoutMs:    getenvInt("TIMEOUT_MS", "1000"),
		RetryCount:   getenvInt("RETRY_COUNT", "3"),
		RetryBackoff: getenvInt("RETRY_BACKOFF", "100"),
		AuthType:     getenv("AUTH_TYPE", "none"),
		AuthRef:      getenv("AUTH_REF", ""),

		ClientID: getenv("CLIENT_ID", ""),
	}

	if config.ClientID == "" {
		log.Fatal("CLIENT_ID is required")
	}

	return config
}

func getenvInt(k string, v string) int {
	e := os.Getenv(k)
	if e == "" {
		num, _ := strconv.Atoi(v)
		return num
	}
	num, err := strconv.Atoi(e)
	if err != nil {
		num, _ = strconv.Atoi(v)
		fmt.Printf("key: %s=%s not a number, default=%s", k, e, v)
		return num
	}
	return num
}

func getenv(k string, v string) string {
	e := os.Getenv(k)
	if e == "" {
		return v
	}
	return e
}
