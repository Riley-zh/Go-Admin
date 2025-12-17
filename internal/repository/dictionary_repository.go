package repository

import (
	"errors"

	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// DictionaryRepository defines the dictionary repository interface
type DictionaryRepository interface {
	// Dictionary operations
	CreateDictionary(dictionary *model.Dictionary) error
	GetDictionaryByID(id uint) (*model.Dictionary, error)
	GetDictionaryByName(name string) (*model.Dictionary, error)
	UpdateDictionary(dictionary *model.Dictionary) error
	DeleteDictionary(id uint) error
	ListDictionaries(page, pageSize int) ([]*model.Dictionary, int64, error)

	// DictionaryItem operations
	CreateDictionaryItem(item *model.DictionaryItem) error
	GetDictionaryItemByID(id uint) (*model.DictionaryItem, error)
	UpdateDictionaryItem(item *model.DictionaryItem) error
	DeleteDictionaryItem(id uint) error
	ListDictionaryItems(dictionaryID, page, pageSize int) ([]*model.DictionaryItem, int64, error)
	ListAllDictionaryItems(dictionaryID int) ([]*model.DictionaryItem, error)
	GetDictionaryItemByValue(dictionaryID uint, value string) (*model.DictionaryItem, error)
}

// dictionaryRepository implements DictionaryRepository interface
type dictionaryRepository struct {
	db *gorm.DB
}

// NewDictionaryRepository creates a new dictionary repository
func NewDictionaryRepository() DictionaryRepository {
	return &dictionaryRepository{
		db: database.GetDB(),
	}
}

// CreateDictionary creates a new dictionary
func (r *dictionaryRepository) CreateDictionary(dictionary *model.Dictionary) error {
	return r.db.Create(dictionary).Error
}

// GetDictionaryByID gets a dictionary by ID
func (r *dictionaryRepository) GetDictionaryByID(id uint) (*model.Dictionary, error) {
	var dictionary model.Dictionary
	err := r.db.Where("id = ? AND status = ?", id, 1).First(&dictionary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dictionary, nil
}

// GetDictionaryByName gets a dictionary by name
func (r *dictionaryRepository) GetDictionaryByName(name string) (*model.Dictionary, error) {
	var dictionary model.Dictionary
	err := r.db.Where("name = ? AND status = ?", name, 1).First(&dictionary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dictionary, nil
}

// UpdateDictionary updates a dictionary
func (r *dictionaryRepository) UpdateDictionary(dictionary *model.Dictionary) error {
	return r.db.Save(dictionary).Error
}

// DeleteDictionary deletes a dictionary (soft delete)
func (r *dictionaryRepository) DeleteDictionary(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Dictionary{}).Error
}

// ListDictionaries lists dictionaries with pagination
func (r *dictionaryRepository) ListDictionaries(page, pageSize int) ([]*model.Dictionary, int64, error) {
	var dictionaries []*model.Dictionary
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&model.Dictionary{}).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&dictionaries).Error
	if err != nil {
		return nil, 0, err
	}

	return dictionaries, total, nil
}

// CreateDictionaryItem creates a new dictionary item
func (r *dictionaryRepository) CreateDictionaryItem(item *model.DictionaryItem) error {
	return r.db.Create(item).Error
}

// GetDictionaryItemByID gets a dictionary item by ID
func (r *dictionaryRepository) GetDictionaryItemByID(id uint) (*model.DictionaryItem, error) {
	var item model.DictionaryItem
	err := r.db.Where("id = ? AND status = ?", id, 1).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// UpdateDictionaryItem updates a dictionary item
func (r *dictionaryRepository) UpdateDictionaryItem(item *model.DictionaryItem) error {
	return r.db.Save(item).Error
}

// DeleteDictionaryItem deletes a dictionary item (soft delete)
func (r *dictionaryRepository) DeleteDictionaryItem(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.DictionaryItem{}).Error
}

// ListDictionaryItems lists dictionary items with pagination
func (r *dictionaryRepository) ListDictionaryItems(dictionaryID, page, pageSize int) ([]*model.DictionaryItem, int64, error) {
	var items []*model.DictionaryItem
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&model.DictionaryItem{}).Where("dictionary_id = ? AND status = ?", dictionaryID, 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("dictionary_id = ? AND status = ?", dictionaryID, 1).Order("sort ASC").Offset(offset).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ListAllDictionaryItems lists all dictionary items
func (r *dictionaryRepository) ListAllDictionaryItems(dictionaryID int) ([]*model.DictionaryItem, error) {
	var items []*model.DictionaryItem
	err := r.db.Where("dictionary_id = ? AND status = ?", dictionaryID, 1).Order("sort ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetDictionaryItemByValue gets a dictionary item by value
func (r *dictionaryRepository) GetDictionaryItemByValue(dictionaryID uint, value string) (*model.DictionaryItem, error) {
	var item model.DictionaryItem
	err := r.db.Where("dictionary_id = ? AND value = ? AND status = ?", dictionaryID, value, 1).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
