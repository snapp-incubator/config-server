package config

import "time"

func Default() Config {
	return Config{
		API: API{Port: 8080, GracefulTimeout: time.Second * 5},
	}
}
