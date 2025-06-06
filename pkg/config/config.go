package config

import "os"

// struct AppConfig
// Armazena as configurações da aplicação
type AppConfig struct {
	RedisAddress string `json:"redis_address"`
	RedisMQURL   string `json:"redis_mq_url"`
	ServerPort   string `json:"server_port"`
}

// Carrega as variáveis de ambiente e retorna a configuração da aplicação
func LoadConfig() *AppConfig {
	return &AppConfig{
		RedisAddress: getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisMQURL:   getEnv("REDIS_MQ_URL", "amqp://guest:guest@localhost:5672/"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
	}
}

// PARA IMPLEMENTAR:
// Usar godotenv ou similar para ler variáveis de ambiente a partir de um arquivo .env
func getEnv(key, defaultValue string) string {
	// A função os.LookupEnv verifica se a variável de ambiente existe
	// e retorna o valor se existir, ou o valor padrão se não existir.
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
