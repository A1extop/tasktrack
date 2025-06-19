package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskTrack/internal/config"
	"taskTrack/internal/services/task/models"
	"taskTrack/internal/services/task/usecase"
	"time"
)

type TaskHandler struct {
	config  *config.Config
	service usecase.ITaskTrackUsecase
}

func NewTaskHandler(config *config.Config, engine *gin.RouterGroup, taskService usecase.ITaskTrackUsecase) {
	handler := TaskHandler{
		config:  config,
		service: taskService,
	}
	router := engine.Group("/task")
	{
		router.GET(
			"",
			handler.getAllTasks)
		router.GET(
			"/:id",
			handler.getTaskById)
		router.POST(
			"",
			handler.createTask)
		router.PATCH(
			"/:id",
			handler.updateTask)
		router.DELETE(
			"/:id",
			handler.deleteTask)
	}
}

func (th *TaskHandler) getAllTasks(ctx *gin.Context) {
	tasks, err := th.service.GetAllTasks(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(*tasks) == 0 {
		ctx.JSON(http.StatusNotFound, tasks)
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

func (th *TaskHandler) getTaskById(ctx *gin.Context) {
	taskId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := th.service.GetTaskById(ctx, taskId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	var duration time.Duration
	if task.Status == models.TaskDone && !task.CompletedAt.IsZero() {
		duration = task.CompletedAt.Sub(task.CreatedAt)
	}
	response := gin.H{
		"id":          task.Id,
		"name":        task.Name,
		"description": task.Description,
		"status":      task.Status,
		"created_at":  task.CreatedAt,
		"duration":    duration.String(),
	}
	ctx.JSON(http.StatusOK, response)
}

func (th *TaskHandler) createTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId, err := th.service.CreateTask(ctx, &task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task_id": taskId})
}

func (th *TaskHandler) updateTask(ctx *gin.Context) {
	taskId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.Id = taskId
	err = th.service.UpdateTask(ctx, &task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task_id": taskId})
}

func (th *TaskHandler) deleteTask(ctx *gin.Context) {
	taskId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = th.service.DeleteTaskById(ctx, taskId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
