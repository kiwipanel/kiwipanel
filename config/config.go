package config

import "sync"

var (
	once     sync.Once // create sync.Once primitive
	instance *Config   // create nil Config struct
	env      *ENV
	mode     bool
)
