package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *taskHandler) ExecuteTask(c *gin.Context) {
	id := c.Param("id")
	taskToExecute, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if taskToExecute == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

}
