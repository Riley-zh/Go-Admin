package service

import (
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/errors"
)

// DictionaryService defines the dictionary service interface
type DictionaryService interface {
	// Dictionary operations
	CreateDictionary(name, title, description string) (*model.Dictionary, error)
	GetDictionaryByID(id uint) (*model.Dictionary, error)
	GetDictionaryByName(name string) (*model.Dictionary, error)
	UpdateDictionary(dictionary *model.Dictionary) error
	DeleteDictionary(id uint) error
	ListDictionaries(page, pageSize int) ([]*model.Dictionary, int64, error)

	// DictionaryItem operations
	CreateDictionaryItem(dictionaryID uint, label, value string, sort, status int) (*model.DictionaryItem, error)
	GetDictionaryItemByID(id uint) (*model.DictionaryItem, error)
	UpdateDictionaryItem(item *model.DictionaryItem) error
	DeleteDictionaryItem(id uint) error
	ListDictionaryItems(dictionaryID, page, pageSize int) ([]*model.DictionaryItem, int64, error)
	ListAllDictionaryItems(dictionaryID int) ([]*model.DictionaryItem, error)
	GetDictionaryItemByValue(dictionaryID uint, value string) (*model.DictionaryItem, error)
}

// dictionaryService implements DictionaryService interface
type dictionaryService struct {
	dictRepo repository.DictionaryRepository
}

// NewDictionaryService creates a new dictionary service
func NewDictionaryService() DictionaryService {
	return &dictionaryService{
		dictRepo: repository.NewDictionaryRepository(),
	}
}

// CreateDictionary creates a new dictionary
func (s *dictionaryService) CreateDictionary(name, title, description string) (*model.Dictionary, error) {
	// Check if dictionary name already exists
	existingDict, err := s.dictRepo.GetDictionaryByName(name)
	if err != nil {
		return nil, err
	}
	if existingDict != nil {
		return nil, errors.Conflict("Dictionary name already exists", "字典名称已存在")
	}

	// Create dictionary
	dictionary := &model.Dictionary{
		Name:        name,
		Title:       title,
		Description: description,
		Status:      1,
	}

	err = s.dictRepo.CreateDictionary(dictionary)
	if err != nil {
		return nil, err
	}

	return dictionary, nil
}

// GetDictionaryByID gets a dictionary by ID
func (s *dictionaryService) GetDictionaryByID(id uint) (*model.Dictionary, error) {
	dictionary, err := s.dictRepo.GetDictionaryByID(id)
	if err != nil {
		return nil, err
	}
	if dictionary == nil {
		return nil, errors.NotFound("Dictionary not found", "字典不存在")
	}

	return dictionary, nil
}

// GetDictionaryByName gets a dictionary by name
func (s *dictionaryService) GetDictionaryByName(name string) (*model.Dictionary, error) {
	dictionary, err := s.dictRepo.GetDictionaryByName(name)
	if err != nil {
		return nil, err
	}
	if dictionary == nil {
		return nil, errors.NotFound("Dictionary not found", "字典不存在")
	}

	return dictionary, nil
}

// UpdateDictionary updates a dictionary
func (s *dictionaryService) UpdateDictionary(dictionary *model.Dictionary) error {
	// Check if dictionary exists
	existingDict, err := s.dictRepo.GetDictionaryByID(dictionary.ID)
	if err != nil {
		return err
	}
	if existingDict == nil {
		return errors.NotFound("Dictionary not found", "字典不存在")
	}

	// Check if dictionary name already exists (excluding current dictionary)
	if dictionary.Name != existingDict.Name {
		otherDict, err := s.dictRepo.GetDictionaryByName(dictionary.Name)
		if err != nil {
			return err
		}
		if otherDict != nil {
			return errors.Conflict("Dictionary name already exists", "字典名称已存在")
		}
	}

	// Update dictionary
	return s.dictRepo.UpdateDictionary(dictionary)
}

// DeleteDictionary deletes a dictionary
func (s *dictionaryService) DeleteDictionary(id uint) error {
	// Check if dictionary exists
	existingDict, err := s.dictRepo.GetDictionaryByID(id)
	if err != nil {
		return err
	}
	if existingDict == nil {
		return errors.NotFound("Dictionary not found", "字典不存在")
	}

	// Delete dictionary
	return s.dictRepo.DeleteDictionary(id)
}

// ListDictionaries lists dictionaries with pagination
func (s *dictionaryService) ListDictionaries(page, pageSize int) ([]*model.Dictionary, int64, error) {
	return s.dictRepo.ListDictionaries(page, pageSize)
}

// CreateDictionaryItem creates a new dictionary item
func (s *dictionaryService) CreateDictionaryItem(dictionaryID uint, label, value string, sort, status int) (*model.DictionaryItem, error) {
	// Check if dictionary exists
	dictionary, err := s.dictRepo.GetDictionaryByID(dictionaryID)
	if err != nil {
		return nil, err
	}
	if dictionary == nil {
		return nil, errors.NotFound("Dictionary not found", "字典不存在")
	}

	// Check if item value already exists in this dictionary
	existingItem, err := s.dictRepo.GetDictionaryItemByValue(dictionaryID, value)
	if err != nil {
		return nil, err
	}
	if existingItem != nil {
		return nil, errors.Conflict("Dictionary item value already exists", "字典项值已存在")
	}

	// Create dictionary item
	item := &model.DictionaryItem{
		DictionaryID: dictionaryID,
		Label:        label,
		Value:        value,
		Sort:         sort,
		Status:       status,
	}

	err = s.dictRepo.CreateDictionaryItem(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetDictionaryItemByID gets a dictionary item by ID
func (s *dictionaryService) GetDictionaryItemByID(id uint) (*model.DictionaryItem, error) {
	item, err := s.dictRepo.GetDictionaryItemByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.NotFound("Dictionary item not found", "字典项不存在")
	}

	return item, nil
}

// UpdateDictionaryItem updates a dictionary item
func (s *dictionaryService) UpdateDictionaryItem(item *model.DictionaryItem) error {
	// Check if item exists
	existingItem, err := s.dictRepo.GetDictionaryItemByID(item.ID)
	if err != nil {
		return err
	}
	if existingItem == nil {
		return errors.NotFound("Dictionary item not found", "字典项不存在")
	}

	// Check if item value already exists in this dictionary (excluding current item)
	if item.Value != existingItem.Value {
		otherItem, err := s.dictRepo.GetDictionaryItemByValue(item.DictionaryID, item.Value)
		if err != nil {
			return err
		}
		if otherItem != nil {
			return errors.Conflict("Dictionary item value already exists", "字典项值已存在")
		}
	}

	// Update dictionary item
	return s.dictRepo.UpdateDictionaryItem(item)
}

// DeleteDictionaryItem deletes a dictionary item
func (s *dictionaryService) DeleteDictionaryItem(id uint) error {
	// Check if item exists
	existingItem, err := s.dictRepo.GetDictionaryItemByID(id)
	if err != nil {
		return err
	}
	if existingItem == nil {
		return errors.NotFound("Dictionary item not found", "字典项不存在")
	}

	// Delete dictionary item
	return s.dictRepo.DeleteDictionaryItem(id)
}

// ListDictionaryItems lists dictionary items with pagination
func (s *dictionaryService) ListDictionaryItems(dictionaryID, page, pageSize int) ([]*model.DictionaryItem, int64, error) {
	// Check if dictionary exists
	dictionary, err := s.dictRepo.GetDictionaryByID(uint(dictionaryID))
	if err != nil {
		return nil, 0, err
	}
	if dictionary == nil {
		return nil, 0, errors.NotFound("Dictionary not found", "字典不存在")
	}

	return s.dictRepo.ListDictionaryItems(dictionaryID, page, pageSize)
}

// ListAllDictionaryItems lists all dictionary items
func (s *dictionaryService) ListAllDictionaryItems(dictionaryID int) ([]*model.DictionaryItem, error) {
	// Check if dictionary exists
	dictionary, err := s.dictRepo.GetDictionaryByID(uint(dictionaryID))
	if err != nil {
		return nil, err
	}
	if dictionary == nil {
		return nil, errors.NotFound("Dictionary not found", "字典不存在")
	}

	return s.dictRepo.ListAllDictionaryItems(dictionaryID)
}

// GetDictionaryItemByValue gets a dictionary item by value
func (s *dictionaryService) GetDictionaryItemByValue(dictionaryID uint, value string) (*model.DictionaryItem, error) {
	// Check if dictionary exists
	dictionary, err := s.dictRepo.GetDictionaryByID(dictionaryID)
	if err != nil {
		return nil, err
	}
	if dictionary == nil {
		return nil, errors.NotFound("Dictionary not found", "字典不存在")
	}

	item, err := s.dictRepo.GetDictionaryItemByValue(dictionaryID, value)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.NotFound("Dictionary item not found", "字典项不存在")
	}

	return item, nil
}
