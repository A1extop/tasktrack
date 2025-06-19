package repository

import (
	"context"
	"sync"
	"taskTrack/internal/domain"
	"taskTrack/internal/services/task/models"
	"time"
)

type ITaskTrackRepository interface {
	Get(ctx context.Context, id int) (*models.Task, error)
	GetAll(ctx context.Context) (*[]models.Task, error)
	Create(ctx context.Context, task *models.Task) (int, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, name string) (bool, error)
	UpdateStatus(ctx context.Context, id int, status models.TaskStatus) error
}
type TaskTrackRepository struct {
	tasks     map[int]*models.Task
	nameIndex map[string]struct{}
	mu        sync.RWMutex
	idGen     int
}

func (r *TaskTrackRepository) Exists(ctx context.Context, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.nameIndex[name]
	return exists, nil
}
func NewTaskRepository() ITaskTrackRepository {
	return &TaskTrackRepository{
		tasks:     make(map[int]*models.Task),
		nameIndex: make(map[string]struct{}),
		idGen:     1,
	}
}
func (r *TaskTrackRepository) Get(ctx context.Context, id int) (*models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, domain.ErrTaskNotFound
	}
	return task, nil
}

func (r *TaskTrackRepository) GetAll(ctx context.Context) (*[]models.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			result = append(result, *task)
		}
	}
	return &result, nil
}

func (r *TaskTrackRepository) Create(ctx context.Context, task *models.Task) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	task.Id = r.idGen
	r.tasks[task.Id] = task
	r.idGen++
	r.nameIndex[task.Name] = struct{}{}
	return task.Id, nil
}
func (r *TaskTrackRepository) UpdateStatus(ctx context.Context, id int, status models.TaskStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return domain.ErrTaskNotFound
	}
	r.tasks[id].Status = status
	if status == models.TaskFailed || status == models.TaskDone {
		r.tasks[id].CompletedAt = time.Now()
	}
	return nil
}
func (r *TaskTrackRepository) Update(ctx context.Context, task *models.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.Id]; !exists {
		return domain.ErrTaskNotFound
	}
	if task.Description != "" {
		r.tasks[task.Id].Description = task.Description
	}
	if task.Name != "" {
		r.tasks[task.Id].Name = task.Name
	}
	if task.Status != "" {
		r.tasks[task.Id].Status = task.Status
	}
	return nil
}

func (r *TaskTrackRepository) Delete(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return domain.ErrTaskNotFound
	}

	delete(r.nameIndex, task.Name)
	delete(r.tasks, id)

	return nil
}
