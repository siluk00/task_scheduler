package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/siluk00/task_scheduler/internal/domain"
)

// CreateTask createss a new task
// @Summary Creates a new Task
// @Description Creates a new task in the system
// @Tags tasks
// @Accpet json
// @Produce json
// @Param task body domain.Task true "taskData"
// @Success 201 {object} domain.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (h *taskHandler) CreateTask(c *gin.Context) {
	var task domain.Task

	// ShouldBindJSON is a method to get the body of the HTTP request and bind it to the task struct.
	// It automatically handles JSON decoding and validation.
	// Returns an error if the JSON is invalid or does not match the struct fields.
	// BindJson differs from shouldBindJSON in that it does not return an error if the binding fails,
	if err := c.ShouldBindJSON(&task); err != nil {
		//c.Json is a method to send a JSON response with a specific status code.
		// 400 is the status code for Bad Request
		c.JSON(400, gin.H{"error": "Invalid task data"})
		return
	}

	if err := task.Validate(); err != nil {
		// 400 is the status code for Bad Request
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if task.Status == "" {
		task.Status = domain.TaskStatusPending // Default status if not provided
	}

	if err := h.repo.Create(c.Request.Context(), &task); err != nil {
		// 500 is the status code for Internal Server Error
		c.JSON(500, gin.H{"error": "Failed to create task"})
		return
	}
	//201 is the status code for Created
	c.JSON(201, gin.H{"message": "Task created successfully", "task_id": task.ID})
}
