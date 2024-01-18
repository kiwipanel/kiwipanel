package config

import (
	"time"
)

type Config struct {
	Server *serverConfig
}

type serverConfig struct {
	Addr         string
	VERSION      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}
