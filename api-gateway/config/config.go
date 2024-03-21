package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string //develop, staging, production

	RedisHost string
	RedisPort int

	UserServiceHost string
	UserServicePort int

	ProductServiceHost string
	ProductServicePort int

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	//context timeout in seconds
	CtxTimeOut int

	LogLevel string
	HTTPPort string

	AccessTokenTimeout  int //minutes
	RefreshTokenTimeout int //hours
	AuthConfigPath      string
	AuthCSVPath         string

	SignInKey string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":4040"))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "db"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "auth"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "mubina2007"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "userservice"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 8081))

	c.ProductServiceHost = cast.ToString(getOrReturnDefault("PRODUCT_SERVICE_HOST", "productservice"))
	c.ProductServicePort = cast.ToInt(getOrReturnDefault("PRODUCT_SERVICE_PORT", 5050))

	c.AccessTokenTimeout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 1000))
	c.RefreshTokenTimeout = cast.ToInt(getOrReturnDefault("REFRESH_TOKEN_TIMEOUT", 1000))

	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_CONFIG_PATH", "./config/auth.conf"))
	c.AuthCSVPath = cast.ToString(getOrReturnDefault("AUTH_CSV_PATH", "./config/auth.csv"))

	c.CtxTimeOut = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	c.SignInKey = cast.ToString(getOrReturnDefault("SIGN_IN_KEY", "golang-exam-3d-month"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return defaultValue
	}

	return defaultValue
}
