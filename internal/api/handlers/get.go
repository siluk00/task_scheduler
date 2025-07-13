package handlers

import "github.com/gin-gonic/gin"

// GetTask Gets a task by its ID
// @Summary Gets a Task
// @Description gets a task by its ID
// @Tags tasks
// @Produce json
// @Param id path string true "taskId"
// @Success 200 {object} domain.Task
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *taskHandler) GetTask(c *gin.Context) {
	//Param is the way to get a path parameter from the URL.
	id := c.Param("id")
	// the context passed to FindById is the context of the current request because
	// it contains information about the request, such as the request ID, user, etc.
	// the context of c doesn't have a deadline or cancellation signal, so it will
	// not timeout or cancel the operation. And it doesnt have information about the request,
	// so it will not be able to access the request ID, user, etc.
	task, err := h.repo.FindById(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(404, "Task not found")
		return
	}

	c.JSON(200, task)
}
