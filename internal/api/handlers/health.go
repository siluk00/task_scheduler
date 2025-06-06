package handlers

import "github.com/gin-gonic/gin"

// O método healthCheck é um manipulador de rota que responde a requisições GET na rota /health.
// Ele retorna um JSON simples indicando que o servidor está funcionando corretamente.
// Isso é útil para verificar se o servidor está ativo e respondendo.
// Ele pode ser usado por ferramentas de monitoramento ou health checks de serviços externos.
// O manipulador recebe um contexto (c *gin.Context) que contém informações sobre a requisição
// e a resposta.
func (h *taskHandler) HealthCheck(c *gin.Context) {
	//gin.H é um atalho para criar um mapa de strings para JSON.
	// Ele é usado para construir respostas JSON de forma simples e rápida.
	// Aqui, estamos retornando um JSON com o status "ok".
	c.JSON(200, gin.H{"status": "ok"})
}
