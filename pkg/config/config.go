package config

import "time"

type Config struct {
	BaseURL  string
	Email    string
	Password string
	Proxy    string
	Timeout  time.Duration
	Debug    bool
}
