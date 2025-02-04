package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	URL            string
	MigrationsPath string
}

type RedisConfig struct {
	Address  string
	Password string
}

func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("DB_URL", "postgres://postgres:postgres@postgres:5432/shortener")
	viper.SetDefault("DB_MIGRATIONS_PATH", "migrations")
	viper.SetDefault("REDIS_ADDRESS", "redis:6379")
	viper.SetDefault("REDIS_PASSWORD", "")

	return &Config{
		Server: ServerConfig{
			Port: viper.GetInt("SERVER_PORT"),
		},
		DB: DatabaseConfig{
			URL:            viper.GetString("DB_URL"),
			MigrationsPath: viper.GetString("DB_MIGRATIONS_PATH"),
		},
		Redis: RedisConfig{
			Address:  viper.GetString("REDIS_ADDRESS"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
	}
}
