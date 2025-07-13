// @title Task Scheduler API
// @version 1.0
// @deescription API for task scheduling
// @host localhost:8080
// @BasePath /

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siluk00/task_scheduler/internal/domain"
	"github.com/siluk00/task_scheduler/internal/repository"
)

type taskHandler struct {
	repo repository.TaskHandler
}

func NewTaskHandler(repo repository.TaskHandler) *taskHandler {
	return &taskHandler{
		repo: repo,
	}
}

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

// DeleteTask removes a task
// @Summary Delete a task
// @Description Delete a task by its id
// @Tags tasks
// @Produce json
// @Param id path string true "task id"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *taskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	taskToDelete, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	if taskToDelete == nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	//c.Request.Context() is used to get the context of the current request.
	// The context is used to pass request-scoped values, deadlines, and cancellation signals.
	// The Request context is sent instead of the context of the handler function because
	// it allows the repository to access the request context, which may contain information
	// about the request, such as the request ID, user, etc.
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// C.Status is a method to set the HTTP status code of the response.
	// http.StatusNoContent is the status code for No Content (204).
	// It indicates that the request was successful, but there is no content to return.
	// This is typically used for DELETE requests where no response body is needed.
	// 204 is the status code for No Content
	c.Status(http.StatusNoContent)
}

func (h *taskHandler) ListTasks(c *gin.Context) {
	status := domain.TaskStatus(c.Query("status"))

	tasks, err := h.repo.List(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	if tasks == nil {
		tasks = []*domain.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}

// GetScheduledTasks lists tasks in a time span
// @Summary List Scheduled Tasks
// @Description List scheduled tasks in a determined timespan
// @Tags taks
// @Produce json
// @Param from query string true "begin date (RFC3339)"
// @Param to query string true "end date (RFC3339)"
// @Success 200 {array} domain.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/scheduled [get]
func (h *taskHandler) GetScheduledTasks(c *gin.Context) {
	// this is the time.RFC3339 format: "2006-01-02T15:04:05Z07:00"
	from, err := time.Parse(time.RFC3339, c.Query("from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
		return
	}

	to, err := time.Parse(time.RFC3339, c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
		return
	}

	if to.Before(from) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to date before from date"})
		return
	}

	tasks, err := h.repo.FindScheduled(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if tasks == nil {
		tasks = []*domain.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}
