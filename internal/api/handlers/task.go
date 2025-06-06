package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/siluk00/task_scheduler/internal/domain"
	"github.com/siluk00/task_scheduler/internal/repository"
)

type taskHandler struct {
	repo repository.TaskRepository
}

func NewTaskHandler(repo repository.TaskRepository) *taskHandler {
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

//TODO: Other handlers like UpdateTask, ListTasks, DeleteTask, etc.
