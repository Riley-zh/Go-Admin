package repository

import (
	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// TaskRepository handles task data operations
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		db: database.GetDB(),
	}
}

// Create saves a new task
func (r *TaskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

// GetByID retrieves a task by its ID
func (r *TaskRepository) GetByID(id uint) (*model.Task, error) {
	var task model.Task
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List retrieves tasks with pagination
func (r *TaskRepository) List(page, pageSize int, status string) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := r.db.Model(&model.Task{})

	// Apply filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tasks).Error
	return tasks, total, err
}

// Update updates a task
func (r *TaskRepository) Update(task *model.Task) error {
	return r.db.Save(task).Error
}

// Delete removes a task by its ID
func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&model.Task{}, id).Error
}

// GetAllActiveTasks retrieves all active tasks
func (r *TaskRepository) GetAllActiveTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("status = ?", "active").Find(&tasks).Error
	return tasks, err
}
