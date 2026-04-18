package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	EnvFile         string
	Port            string
	VersionAPI      string
	AllowForbidden  string
	secret          string
	Exp             int
	Ref_time        int
	RedisAddr       string
	redisPassword   string
	RedisDB         int
	ExpireEmailCode int
	Namespace       string

	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	EmailAdmin    string
	FullNameAdmin string

	ImageAPIClient string

	ProxyConnectTimeout string
	ProxyReadTimeout    string
	ProxySendTimeout    string

	MaxMetricsPerClient int

	BaseURLIngress string
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
		EnvFile:         getenv("ENV_FILE", "dev"),
		Port:            getenv("PORT", "8000"),
		VersionAPI:      getenv("VERSION_API", "v1"),
		AllowForbidden:  getenv("ALLOW_CHECK_403", "True"),
		secret:          getenv("SECRET", "1984"),
		Exp:             getenvInt("EXPIRE_JWT", "600"),
		Ref_time:        getenvInt("REFRESH_TIME_JWT", "604800"),
		RedisAddr:       getenv("REDIS_ADDR", "localhost:6379"),
		redisPassword:   getenv("REDIS_PASSWORD", "1984"),
		RedisDB:         getenvInt("REDIS_DB", "0"),
		ExpireEmailCode: getenvInt("EXPIRE_EMAIL_CODE", "300"),
		Namespace:       getenv("NAMESPACE", "default"),

		SMTPHost: getenv("SMTP_HOST", "smtp.mail.ru"),
		SMTPPort: getenv("SMTP_PORT", "587"),
		SMTPUser: getenv("SMTP_USER", "your@mail.ru"),
		SMTPPass: getenv("SMTP_PASS", "app_password"),
		SMTPFrom: getenv("SMTP_FROM", "your@mail.ru"),

		EmailAdmin:    getenv("EMAIL_ADMIN", "temp@mail.ru"),
		FullNameAdmin: getenv("FULL_NAME_ADMIN", "Ivan Ivanov Ivanovich"),

		ImageAPIClient: getenv("IMAGE_API_CLIENT", "api-client-runtime"),

		ProxyConnectTimeout: getenv("GLOBAL_PROXY_CONNECT_TIMEOUT", "30"),
		ProxyReadTimeout:    getenv("GLOBAL_PROXY_READ_TIMEOUT", "120"),
		ProxySendTimeout:    getenv("GLOBAL_PROXY_SEND_TIMEOUT", "120"),
	
		MaxMetricsPerClient: getenvInt("MAX_METRICS_PER_CLIENT", "1000"),

		BaseURLIngress: getenv("BASE_URL_INGRESS", "http://localhost:8080"),
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

func (c *Config) GetSecret() string {
	return c.secret
}

func (c *Config) GetRedisPassword() string {
	return c.redisPassword
}
