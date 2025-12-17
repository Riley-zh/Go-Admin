package service

import (
	"go-admin/internal/repository"
	"go-admin/pkg/errors"
)

// BaseService defines the base service interface
type BaseService[T repository.BaseModel] interface {
	GetByID(id uint) (T, error)
	Delete(id uint) error
}

// baseService implements BaseService interface
type baseService[T repository.BaseModel] struct {
	repo repository.BaseRepository[T]
}

// NewBaseService creates a new base service
func NewBaseService[T repository.BaseModel](entity T) BaseService[T] {
	return &baseService[T]{
		repo: repository.NewBaseRepository(entity),
	}
}

// GetByID gets an entity by ID
func (s *baseService[T]) GetByID(id uint) (T, error) {
	entity, err := s.repo.GetByID(id)
	if err != nil {
		return entity, err
	}

	// Check if entity is nil (zero value)
	if entity.GetID() == 0 {
		return entity, errors.NotFound("Entity not found", "实体不存在")
	}

	return entity, nil
}

// Delete deletes an entity
func (s *baseService[T]) Delete(id uint) error {
	// Check if entity exists
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// CheckExists checks if an entity exists by ID
func CheckExists[T repository.BaseModel](repo repository.BaseRepository[T], id uint) error {
	entity, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	if entity.GetID() == 0 {
		return errors.NotFound("Entity not found", "实体不存在")
	}

	return nil
}
