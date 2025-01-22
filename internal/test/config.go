package test

import "github.com/FotiadisM/mock-microservice/internal/config"

func NewConfig() *config.Config {
	return &config.Config{
		Server:  config.Server{},
		DB:      config.DB{},
		Logging: config.Logging{},
	}
}
