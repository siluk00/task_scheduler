package config

import "os"

// struct AppConfig
// Stores the configurations of the application
type AppConfig struct {
	RedisAddress string `json:"redis_address"`
	RedisMQURL   string `json:"redis_mq_url"`
	ServerPort   string `json:"server_port"`
}

// Load ambient variables onto the AppConfig struct
func LoadConfig() *AppConfig {
	return &AppConfig{
		RedisAddress: getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisMQURL:   getEnv("REDIS_MQ_URL", "amqp://user:password@localhost:5672/"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
	}
}

// TODO: Use viper here instead
func getEnv(key, defaultValue string) string {
	// A função os.LookupEnv verifica se a variável de ambiente existe
	// e retorna o valor se existir, ou o valor padrão se não existir.
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
