package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siluk00/task_scheduler/internal/domain"
)

// UpdateTask updates a task
// @Summary Updates a task
// @Description Updates an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "task id"
// @Success 200 {object} domain.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func (h *taskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	existingTask, err := h.repo.FindById(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve task"})
		return
	}

	if existingTask == nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	var task domain.Task
	// ShouldBindJSON is used to bind the JSON body of the request to the task struct.
	if err := c.ShouldBindJSON(&task); err != nil {
		// 400 is the status code for Bad Request
		// c.Json is a method to send a JSON response with a specific status code.
		c.JSON(400, gin.H{"error": "Invalid task data"})
		return
	}

	task.ID = id // Ensure the task ID is set to the path parameter ID

	if err := task.Validate(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task.CreatedAt = existingTask.CreatedAt // Preserve the original created time
	task.UpdatedAt = time.Now()             // Preserve the original updated time

	if err := h.repo.Update(c.Request.Context(), &task); err != nil {
		// 500 is the status code for Internal Server Error
		c.JSON(500, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(200, gin.H{"message": "Task updated successfully", "task_id": task.ID})
}
