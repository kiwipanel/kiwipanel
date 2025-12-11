package config

import (
	"log"
	"sync"
	"text/template"

	"gorm.io/gorm"
)

var (
	once sync.Once // create sync.Once primitive
	mode bool
)

type AppConfig struct {
	KIWIPANEL_MODE string
	UseCache       bool
	TemplateCache  map[string]*template.Template
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Test           string
	InProduction   bool
	AUTH_USER      string
	AUTH_PASSWORD  string
	DB             *gorm.DB
}
