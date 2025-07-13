package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siluk00/task_scheduler/internal/domain"
)

// ListTasks lists all tasks, maybe filtered by status
// @Summary Lists tasks
// @Description Lists tasks, it can be filtered by status
// @Tags tasks
// @Produce json
// @Param status query string false "Status for filering (pending, running, completed, failed)"
// @Success 200 {array} domain.Task
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
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
