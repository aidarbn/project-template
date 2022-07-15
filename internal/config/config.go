package config

import "time"

// Config is responsible for application startup configuration.
// `envconfig` is a specific tag for https://github.com/kelseyhightower/envconfig package.
type Config struct {
	LogLevel           string        `envconfig:"LOG_LEVEL" default:"debug"`
	ServerPort         int           `envconfig:"SERVER_PORT" default:"8080"`
	ServerReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"15s"`
	ServerWriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"15s"`
	DBUser             string        `envconfig:"DB_USER" default:"postgres"`
	DBPass             string        `envconfig:"DB_PASS" default:"postgres"`
	DBHost             string        `envconfig:"DB_HOST" default:"db"`
	DBPort             int           `envconfig:"DB_PORT" default:"5432"`
	DBName             string        `envconfig:"DB_NAME" default:"postgres"`
	SentryDSN          string        `envconfig:"SENTRY_DSN"`
	SentryENV          string        `envconfig:"SENTRY_ENV" default:"staging"`
}

type TestConfig struct {
	DBUser string `envconfig:"TEST_DB_USER" default:"test_postgres"`
	DBPass string `envconfig:"TEST_DB_PASS" default:"test_postgres"`
	DBHost string `envconfig:"TEST_DB_HOST" default:"db"`
	DBPort int    `envconfig:"TEST_DB_PORT" default:"5432"`
	DBName string `envconfig:"TEST_DB_NAME" default:"test_postgres"`
}
