package config

import "bitbucket.org/creativeadvtech/project-template/pkg/common"

// Config is responsible for application startup configuration.
// `envconfig` is a specific tag for https://github.com/kelseyhightower/envconfig package.
type Config struct {
	common.Config
	common.DbConfig
	common.SentryConfig
}
