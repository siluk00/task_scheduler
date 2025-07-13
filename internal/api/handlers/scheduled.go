package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siluk00/task_scheduler/internal/domain"
)

// GetScheduledTasks lists tasks in a time span
// @Summary List Scheduled Tasks
// @Description List scheduled tasks in a determined timespan
// @Tags tasks
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
