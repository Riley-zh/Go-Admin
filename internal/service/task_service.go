package service

import (
	"fmt"
	"sync"
	"time"

	"go-admin/internal/model"
	"go-admin/internal/repository"

	"github.com/robfig/cron/v3"
)

// TaskService handles task scheduling business logic
type TaskService struct {
	taskRepo *repository.TaskRepository
	cron     *cron.Cron
	entries  map[uint]cron.EntryID
	mu       sync.RWMutex
}

// NewTaskService creates a new task service
func NewTaskService() *TaskService {
	service := &TaskService{
		taskRepo: repository.NewTaskRepository(),
		cron:     cron.New(),
		entries:  make(map[uint]cron.EntryID),
	}

	// Start cron scheduler
	service.cron.Start()

	// Load existing active tasks
	service.loadActiveTasks()

	return service
}

// loadActiveTasks loads all active tasks from database and schedules them
func (s *TaskService) loadActiveTasks() {
	tasks, err := s.taskRepo.GetAllActiveTasks()
	if err != nil {
		// Log error but continue
		fmt.Printf("Failed to load active tasks: %v\n", err)
		return
	}

	for _, task := range tasks {
		s.scheduleTask(&task)
	}
}

// scheduleTask schedules a task using cron
func (s *TaskService) scheduleTask(task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing entry if any
	if entryID, exists := s.entries[task.ID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, task.ID)
	}

	// Schedule new entry
	entryID, err := s.cron.AddFunc(task.CronExpr, func() {
		s.executeTask(task)
	})
	if err != nil {
		return fmt.Errorf("failed to schedule task: %w", err)
	}

	// Store entry ID
	s.entries[task.ID] = entryID

	// Update next run time
	nextRun := s.cron.Entry(entryID).Next
	task.NextRun = &nextRun
	if err := s.taskRepo.Update(task); err != nil {
		fmt.Printf("Failed to update task next run time: %v\n", err)
	}

	return nil
}

// executeTask executes a scheduled task
func (s *TaskService) executeTask(task *model.Task) {
	now := time.Now()

	// Update last run time
	task.LastRun = &now

	// Execute task based on handler
	switch task.Handler {
	case "sample_task":
		s.executeSampleTask(task)
	default:
		// Log unknown handler
		fmt.Printf("Unknown task handler: %s\n", task.Handler)
		task.Status = "error"
	}

	// Update task
	if err := s.taskRepo.Update(task); err != nil {
		fmt.Printf("Failed to update task: %v\n", err)
	}
}

// executeSampleTask executes a sample task for demonstration
func (s *TaskService) executeSampleTask(task *model.Task) {
	fmt.Printf("Executing sample task: %s (ID: %d)\n", task.Name, task.ID)

	// Simulate some work
	time.Sleep(1 * time.Second)

	fmt.Printf("Sample task completed: %s (ID: %d)\n", task.Name, task.ID)
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(name, description, cronExpr, handler string, userID uint) (*model.Task, error) {
	task := &model.Task{
		Name:        name,
		Description: description,
		CronExpr:    cronExpr,
		Handler:     handler,
		Status:      "active",
		CreatedBy:   userID,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Schedule task if it's active
	if task.Status == "active" {
		s.scheduleTask(task)
	}

	return task, nil
}

// GetTaskByID retrieves a task by its ID
func (s *TaskService) GetTaskByID(id uint) (*model.Task, error) {
	return s.taskRepo.GetByID(id)
}

// ListTasks retrieves tasks with pagination
func (s *TaskService) ListTasks(page, pageSize int, status string) ([]model.Task, int64, error) {
	return s.taskRepo.List(page, pageSize, status)
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(id uint, name, description, cronExpr, handler, status string) (*model.Task, error) {
	// First get the existing task
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Update fields
	task.Name = name
	task.Description = description
	task.CronExpr = cronExpr
	task.Handler = handler
	task.Status = status

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Reschedule or remove task based on status
	s.mu.Lock()
	if entryID, exists := s.entries[id]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, id)
	}
	s.mu.Unlock()

	if task.Status == "active" {
		s.scheduleTask(task)
	}

	return task, nil
}

// DeleteTask removes a task by its ID
func (s *TaskService) DeleteTask(id uint) error {
	// Remove from scheduler
	s.mu.Lock()
	if entryID, exists := s.entries[id]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, id)
	}
	s.mu.Unlock()

	// Remove from database
	if err := s.taskRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// RunTaskImmediately runs a task immediately
func (s *TaskService) RunTaskImmediately(id uint) error {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Execute task in a goroutine to avoid blocking
	go s.executeTask(task)

	return nil
}
