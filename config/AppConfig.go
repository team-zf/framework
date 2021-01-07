package config

type AppConfig struct {
	Settings map[string]interface{}
	Logger   *LoggerConfig
	Table    *TableConfig
	Redis    *RedisConfig
}
