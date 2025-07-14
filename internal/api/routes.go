package api

import (
	_ "github.com/siluk00/task_scheduler/docs"
	"github.com/siluk00/task_scheduler/internal/api/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Aqui o gin configura as rotas do servidor.
// As rotas são definidas usando o método GET, POST, PUT, DELETE, etc. do gin.Engine.
// Cada rota é associada a uma função (handler) que será executada quando a rota for acessada.

func (s *Server) setupRoutes() {

	//s.router.GET("/metrics", s.metricsHandler)

	taskHandler := handlers.NewTaskHandler(s.taskRepo)

	s.router.GET("/health", taskHandler.HealthCheck)

	taskGroup := s.router.Group("/tasks")
	{
		taskGroup.POST("/", taskHandler.CreateTask)
		taskGroup.GET("/:id", taskHandler.GetTask)
		taskGroup.GET("/health", taskHandler.HealthCheck)
		taskGroup.PUT("/:id", taskHandler.UpdateTask)
		taskGroup.DELETE("/:id", taskHandler.DeleteTask)
		taskGroup.GET("/", taskHandler.ListTasks)
		taskGroup.GET("/scheduled", taskHandler.GetScheduledTasks)
		taskGroup.PUT("/tasks/:id/execute", taskHandler.ExecuteTask)
	}

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
