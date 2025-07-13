package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/siluk00/task_scheduler/internal/repository"
	"github.com/siluk00/task_scheduler/internal/repository/redis"
	"github.com/siluk00/task_scheduler/pkg/config"
)

type Server struct {
	config   *config.AppConfig
	router   *gin.Engine
	taskRepo repository.TaskHandler
	//Adicionar serviços/repositorios aqui
}

// Creates the server client and tests it, creates the redis repository
// Creates the server and setup the routes
func NewServer(cfg *config.AppConfig) (*Server, error) {
	//New Client with Options like Address, Password, DB
	// O redis.NewClient cria um novo cliente Redis com as opções fornecidas.
	rdb := goRedis.NewClient(&goRedis.Options{
		Addr:     cfg.RedisAddress,
		Password: "", // Senha se necessário
		DB:       0,  // Usar o banco de dados padrão
	})

	// Rdb ping is used to check if the Redis server is reachable.
	//The context.Background() is used to create a context for the ping operation.
	// It is a simple way to ensure that the Redis server is running and accessible.
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	taskRepo := redis.NewTaskRepository(rdb)

	server := &Server{
		config:   cfg,
		router:   gin.Default(),
		taskRepo: taskRepo,
	}

	server.setupRoutes()
	return server, nil
}

// Runs the server
func (s *Server) Start() error {
	return s.router.Run(":" + s.config.ServerPort)
}
