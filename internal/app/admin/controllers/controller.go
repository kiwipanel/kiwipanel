package controllers

import "github.com/kiwipanel/kiwpanel/config"

type Controller struct {
	config *config.AppConfig
}

func New(cf *config.AppConfig) *Controller {
	return &Controller{
		config: cf,
	}
}
