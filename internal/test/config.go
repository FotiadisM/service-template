package test

import "github.com/FotiadisM/service-template/internal/config"

func NewConfig() *config.Config {
	return &config.Config{
		Inst: config.Instrumentation{},
		Server: config.Server{
			Reflection:             false,
			DisableRESTTranscoding: false,
		},
		DB:      config.DB{},
		Logging: config.Logging{},
		Cors:    config.Cors{},
		Redis:   config.Redis{},
	}
}
