package usecase

import (
	"context"
	"taskTrack/internal/domain"
	"taskTrack/internal/services/task/models"
	"taskTrack/internal/services/task/repository"
	"time"
)

type ITaskTrackUsecase interface {
	GetTaskById(ctx context.Context, id int) (*models.Task, error)
	GetAllTasks(ctx context.Context) (*[]models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) (int, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTaskById(ctx context.Context, id int) error
}

type TaskTrackUsecase struct {
	repo repository.ITaskTrackRepository
}

func NewTaskTrackUsecase(repo repository.ITaskTrackRepository) ITaskTrackUsecase {
	return &TaskTrackUsecase{repo: repo}

}
func (t *TaskTrackUsecase) DeleteTaskById(ctx context.Context, id int) error {
	task, err := t.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if task.Status == models.TaskRunning {
		return domain.ErrCannotDeleteTask
	}
	return t.repo.Delete(ctx, id)
}

func (t *TaskTrackUsecase) CreateTask(ctx context.Context, task *models.Task) (int, error) {
	if task.Name == "" {
		return -1, domain.ErrNamesTaskIsEmpty
	}
	if len(task.Name) > 100 {
		return -1, domain.ErrNameTooLong
	}

	if exists, err := t.repo.Exists(ctx, task.Name); err != nil || exists {
		return -1, domain.ErrTaskAlreadyExists
	}

	if task.Status == "" {
		task.Status = models.TaskPending
	} else if !models.IsValidStatus(task.Status) {
		return -1, domain.ErrInvalidStatus
	}
	task.CreatedAt = time.Now()
	id, err := t.repo.Create(ctx, task)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (t *TaskTrackUsecase) GetAllTasks(ctx context.Context) (*[]models.Task, error) {
	return t.repo.GetAll(ctx)
}
func (t *TaskTrackUsecase) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}
	return t.repo.Get(ctx, id)
}
func (t *TaskTrackUsecase) UpdateTask(ctx context.Context, task *models.Task) error {
	if task.Name == "" && task.Description == "" && task.Status == "" {
		return domain.ErrDataIsEmpty
	}
	_, err := t.GetTaskById(ctx, task.Id)
	if err != nil {
		return domain.ErrTaskNotFound
	}

	if task.Status != "" {
		if err := task.SetStatus(task.Status); err != nil {
			return domain.ErrInvalidStatus
		}
	}

	return t.repo.Update(ctx, task)
}
