package taskprocessor

//в этом пакете происходит имитация выполнения задачи
import (
	"context"
	"log"
	"math/rand"
	"sync"
	"taskTrack/internal/services/task/models"
	"taskTrack/internal/services/task/repository"
	"time"
)

// в целом, можно вынести в конфиг
const (
	defaultTaskChanSize   = 100
	defaultDispatcherTick = 1 * time.Second
	minProcessingTime     = 3 * time.Minute
	maxProcessingTime     = 5 * time.Minute
)

type Processor struct {
	repo       repository.ITaskTrackRepository
	workerPool int
	stop       chan struct{}
	wg         sync.WaitGroup
	tasks      chan models.Task
}

func NewProcessor(repo repository.ITaskTrackRepository, workerPool int) *Processor {
	return &Processor{
		repo:       repo,
		workerPool: workerPool,
		stop:       make(chan struct{}),
		tasks:      make(chan models.Task, defaultTaskChanSize),
	}
}

func (p *Processor) Start(ctx context.Context) {
	p.wg.Add(p.workerPool + 1)

	go p.dispatcher(ctx)

	for i := 0; i < p.workerPool; i++ {
		go p.worker(ctx, i)
	}
}

func (p *Processor) Stop() {
	close(p.stop)
	p.wg.Wait()
}

func (p *Processor) dispatcher(ctx context.Context) {
	defer p.wg.Done()
	ticker := time.NewTicker(defaultDispatcherTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.fetchAndDispatchTasks(ctx)
		case <-p.stop:
			log.Println("Dispatcher stop")
			close(p.tasks)
			return
		case <-ctx.Done():
			log.Println("Dispatcher canceled-  context")
			close(p.tasks)
			return
		}
	}
}

func (p *Processor) fetchAndDispatchTasks(ctx context.Context) {
	tasks, err := p.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		return
	}

	for _, task := range *tasks {
		if task.Status == models.TaskPending {
			if err := p.repo.UpdateStatus(ctx, task.Id, models.TaskRunning); err != nil {
				log.Printf("Failed to update task %d status: %v", task.Id, err)
				continue
			}
			select {
			case p.tasks <- task:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (p *Processor) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	defer log.Printf("Worker %d exiting", id)

	for task := range p.tasks {
		p.processTask(ctx, task)
	}
}

func (p *Processor) processTask(ctx context.Context, task models.Task) {
	processingTime := randomDuration(minProcessingTime, maxProcessingTime)

	select {
	case <-time.After(processingTime):
		p.finalizeTask(ctx, task)
	case <-ctx.Done():
		_ = p.repo.UpdateStatus(ctx, task.Id, models.TaskFailed)
	}
}

func (p *Processor) finalizeTask(ctx context.Context, task models.Task) {
	status := models.TaskDone
	if rand.Intn(2) == 1 {
		status = models.TaskFailed
	}

	if err := p.repo.UpdateStatus(ctx, task.Id, status); err != nil {
		log.Printf("Failed to update task %d status: %v", task.Id, err)
	}
}

func randomDuration(minimum, maximum time.Duration) time.Duration {
	return minimum + time.Duration(rand.Int63n(int64(maximum-minimum)))
}
