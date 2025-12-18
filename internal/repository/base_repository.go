package repository

import (
	"errors"

	"go-admin/internal/database"

	"gorm.io/gorm"
)

// BaseModel represents the base model interface
type BaseModel interface {
	GetID() uint
	GetStatus() int
	SetName(string)
	GetName() string
}

// BaseRepository defines the base repository interface
type BaseRepository[T BaseModel] interface {
	Create(entity T) error
	GetByID(id uint) (T, error)
	GetByName(name string) (T, error)
	Update(entity T) error
	Delete(id uint) error
	List(page, pageSize int) ([]T, int64, error)
	Count() (int64, error)
	Exists(id uint) (bool, error)
}

// baseRepository implements BaseRepository interface
type baseRepository[T BaseModel] struct {
	db     *gorm.DB
	entity T // Used for type inference
}

// NewBaseRepository creates a new base repository
func NewBaseRepository[T BaseModel](entity T) BaseRepository[T] {
	return &baseRepository[T]{
		db:     database.GetDB(),
		entity: entity,
	}
}

// Create creates a new entity
func (r *baseRepository[T]) Create(entity T) error {
	return r.db.Create(entity).Error
}

// GetByID gets an entity by ID
func (r *baseRepository[T]) GetByID(id uint) (T, error) {
	var entity T
	err := r.db.Where("id = ? AND status = ?", id, 1).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, nil // Return zero value with nil error for not found
		}
		return entity, err
	}
	return entity, nil
}

// GetByName gets an entity by name
func (r *baseRepository[T]) GetByName(name string) (T, error) {
	var entity T
	err := r.db.Where("name = ? AND status = ?", name, 1).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, nil // Return zero value with nil error for not found
		}
		return entity, err
	}
	return entity, nil
}

// Count counts the number of entities
func (r *baseRepository[T]) Count() (int64, error) {
	var count int64
	err := r.db.Model(r.entity).Where("status = ?", 1).Count(&count).Error
	return count, err
}

// Exists checks if an entity exists by ID
func (r *baseRepository[T]) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(r.entity).Where("id = ? AND status = ?", id, 1).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update updates an entity
func (r *baseRepository[T]) Update(entity T) error {
	return r.db.Save(entity).Error
}

// Delete deletes an entity (soft delete)
func (r *baseRepository[T]) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(r.entity).Error
}

// List lists entities with pagination
func (r *baseRepository[T]) List(page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(r.entity).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
