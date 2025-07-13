package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
