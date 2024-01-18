package config

import (
	"log"
	"text/template"

	"gorm.io/gorm"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Test          string
	InProduction  bool
	AUTH_USER     string
	AUTH_PASSWORD string
	DB            *gorm.DB
}
