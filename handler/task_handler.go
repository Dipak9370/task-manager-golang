package handler

import (
	"net/http"
	"strconv"
	"task-manager/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{Service: service}
}

// CreateTask godoc
//	@Summary	Create task
//	@Tags		Tasks
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		map[string]string	true	"Create Task"
//	@Success	200		{object}	map[string]string
//	@Router		/tasks [post]

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	c.ShouldBindJSON(&req)

	userID := c.GetString("user_id")

	err := h.Service.CreateTask(c, req.Title, req.Description, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task created"})
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	task, err := h.Service.GetTask(c, id, userID, role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// ListTasks godoc
//	@Summary	List tasks
//	@Tags		Tasks
//	@Security	BearerAuth
//	@Produce	json
//	@Router		/tasks [get]

func (h *TaskHandler) ListTasks(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tasks, err := h.Service.ListTasks(c, userID, role, status, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// DeleteTask godoc
//	@Summary	Delete task
//	@Tags		Tasks
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id	path	string	true	"Task ID"
//	@Router		/tasks/{id} [delete]

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	err := h.Service.DeleteTask(c, id, userID, role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
