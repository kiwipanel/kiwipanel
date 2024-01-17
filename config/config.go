package config

import (
	"html/template"
	"log"
	"sync"
	"time"
)

type Config struct {
	Server *serverConfig
}

type serverConfig struct {
	Addr         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Test          string
	InProduction  bool
	AUTH_USER     string
	AUTH_PASSWORD string
}

var (
	once     sync.Once // create sync.Once primitive
	instance *Config   // create nil Config struct

)
