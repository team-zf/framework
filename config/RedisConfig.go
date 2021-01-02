package config

type RedisConfig struct {
	Addr        string
	Password    string
	MaxActive   int
	MaxIdle     int
	IdleTimeout int
}
