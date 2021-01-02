package config

import "time"

type MySqlConfig struct {
	Dsn         string
	MaxOpenNum  int
	MaxIdleNum  int
	MaxLifetime time.Duration
}
