package service

import (
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/errors"
)

// MenuService defines the menu service interface
type MenuService interface {
	CreateMenu(name, title, icon, path, component, redirect, permission string, parentID, sort, status, hidden int) (*model.Menu, error)
	GetMenuByID(id uint) (*model.Menu, error)
	GetMenuByName(name string) (*model.Menu, error)
	UpdateMenu(menu *model.Menu) error
	DeleteMenu(id uint) error
	ListMenus(page, pageSize int) ([]*model.Menu, int64, error)
	GetMenuTree() ([]*MenuTreeNode, error)
}

// MenuTreeNode represents a menu tree node
type MenuTreeNode struct {
	*model.Menu
	Children []*MenuTreeNode `json:"children"`
}

// menuService implements MenuService interface
type menuService struct {
	menuRepo repository.MenuRepository
}

// NewMenuService creates a new menu service
func NewMenuService() MenuService {
	return &menuService{
		menuRepo: repository.NewMenuRepository(),
	}
}

// CreateMenu creates a new menu
func (s *menuService) CreateMenu(name, title, icon, path, component, redirect, permission string, parentID, sort, status, hidden int) (*model.Menu, error) {
	// Check if menu name already exists
	existingMenu, err := s.menuRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if existingMenu != nil {
		return nil, errors.Conflict("Menu name already exists", "菜单名称已存在")
	}

	// Create menu
	menu := &model.Menu{
		Name:       name,
		Title:      title,
		Icon:       icon,
		Path:       path,
		Component:  component,
		Redirect:   redirect,
		Permission: permission,
		ParentID:   uint(parentID),
		Sort:       sort,
		Status:     status,
		Hidden:     hidden,
	}

	err = s.menuRepo.Create(menu)
	if err != nil {
		return nil, err
	}

	return menu, nil
}

// GetMenuByID gets a menu by ID
func (s *menuService) GetMenuByID(id uint) (*model.Menu, error) {
	menu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, errors.NotFound("Menu not found", "菜单不存在")
	}

	return menu, nil
}

// GetMenuByName gets a menu by name
func (s *menuService) GetMenuByName(name string) (*model.Menu, error) {
	menu, err := s.menuRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, errors.NotFound("Menu not found", "菜单不存在")
	}

	return menu, nil
}

// UpdateMenu updates a menu
func (s *menuService) UpdateMenu(menu *model.Menu) error {
	// Check if menu exists
	existingMenu, err := s.menuRepo.GetByID(menu.ID)
	if err != nil {
		return err
	}
	if existingMenu == nil {
		return errors.NotFound("Menu not found", "菜单不存在")
	}

	// Check if menu name already exists (excluding current menu)
	if menu.Name != existingMenu.Name {
		otherMenu, err := s.menuRepo.GetByName(menu.Name)
		if err != nil {
			return err
		}
		if otherMenu != nil {
			return errors.Conflict("Menu name already exists", "菜单名称已存在")
		}
	}

	// Update menu
	return s.menuRepo.Update(menu)
}

// DeleteMenu deletes a menu
func (s *menuService) DeleteMenu(id uint) error {
	// Check if menu exists
	existingMenu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existingMenu == nil {
		return errors.NotFound("Menu not found", "菜单不存在")
	}

	// Delete menu
	return s.menuRepo.Delete(id)
}

// ListMenus lists menus with pagination
func (s *menuService) ListMenus(page, pageSize int) ([]*model.Menu, int64, error) {
	return s.menuRepo.List(page, pageSize)
}

// GetMenuTree gets menu tree
func (s *menuService) GetMenuTree() ([]*MenuTreeNode, error) {
	// Get all menus
	menus, err := s.menuRepo.ListAll()
	if err != nil {
		return nil, err
	}

	// Build menu tree
	return buildMenuTree(menus, 0), nil
}

// buildMenuTree builds menu tree from flat menu list
func buildMenuTree(menus []*model.Menu, parentID uint) []*MenuTreeNode {
	var tree []*MenuTreeNode

	// Find children of the parent
	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := &MenuTreeNode{
				Menu:     menu,
				Children: buildMenuTree(menus, menu.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}
