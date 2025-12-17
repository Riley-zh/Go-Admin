package repository

import (
	"errors"

	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// MenuRepository defines the menu repository interface
type MenuRepository interface {
	Create(menu *model.Menu) error
	GetByID(id uint) (*model.Menu, error)
	GetByName(name string) (*model.Menu, error)
	Update(menu *model.Menu) error
	Delete(id uint) error
	List(page, pageSize int) ([]*model.Menu, int64, error)
	ListByParentID(parentID uint) ([]*model.Menu, error)
	ListAll() ([]*model.Menu, error)
}

// menuRepository implements MenuRepository interface
type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new menu repository
func NewMenuRepository() MenuRepository {
	return &menuRepository{
		db: database.GetDB(),
	}
}

// Create creates a new menu
func (r *menuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

// GetByID gets a menu by ID
func (r *menuRepository) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Where("id = ? AND status = ?", id, 1).First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

// GetByName gets a menu by name
func (r *menuRepository) GetByName(name string) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Where("name = ? AND status = ?", name, 1).First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

// Update updates a menu
func (r *menuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

// Delete deletes a menu (soft delete)
func (r *menuRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Menu{}).Error
}

// List lists menus with pagination
func (r *menuRepository) List(page, pageSize int) ([]*model.Menu, int64, error) {
	var menus []*model.Menu
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&model.Menu{}).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&menus).Error
	if err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

// ListByParentID lists menus by parent ID
func (r *menuRepository) ListByParentID(parentID uint) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.db.Where("parent_id = ? AND status = ?", parentID, 1).Order("sort ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// ListAll lists all menus
func (r *menuRepository) ListAll() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.db.Where("status = ?", 1).Order("sort ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
