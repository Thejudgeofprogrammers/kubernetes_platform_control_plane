package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	VersionAPI        string
	AllowUnauthorized string
	AllowForbidden    string
}

func LoadEnv() *Config {
	rootDir, _ := os.Getwd()
	nameEnv := os.Getenv("CONFIG_FILE")
	log.Println("env:", nameEnv)
	if nameEnv == "" {
		nameEnv = ".env.dev"
	}
	path := filepath.Join(rootDir, nameEnv)

	if _, err := os.Stat(path); err != nil {
		log.Fatal("Error file env not exists")
	}

	err := godotenv.Load(path)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{
		Port:              getenv("PORT", "8000"),
		VersionAPI:        getenv("VERSION_API", "v1"),
		AllowUnauthorized: getenv("ALLOW_CHECK_401", "True"),
		AllowForbidden:    getenv("ALLOW_CHECK_403", "True"),
	}

	return config
}

func getenv(k string, v string) string {
	e := os.Getenv(k)
	if e == "" {
		return v
	}
	return e
}
