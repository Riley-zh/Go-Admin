package repository

import (
	"encoding/json"

	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// resourceRepository implements ResourceRepository
type resourceRepository struct {
	db *gorm.DB
}

// NewResourceRepository creates a new resource repository
func NewResourceRepository() ResourceRepository {
	return &resourceRepository{
		db: database.GetDB(),
	}
}

// Create creates a new resource
func (r *resourceRepository) Create(resource *model.Resource) error {
	return r.db.Create(resource).Error
}

// Update updates a resource
func (r *resourceRepository) Update(resource *model.Resource) error {
	return r.db.Save(resource).Error
}

// Delete deletes a resource
func (r *resourceRepository) Delete(id uint) error {
	return r.db.Delete(&model.Resource{}, id).Error
}

// GetByID gets a resource by ID
func (r *resourceRepository) GetByID(id uint) (*model.Resource, error) {
	var resource model.Resource
	err := r.db.First(&resource, id).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetByName gets a resource by name
func (r *resourceRepository) GetByName(name string) (*model.Resource, error) {
	var resource model.Resource
	err := r.db.Where("name = ?", name).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// List lists resources with pagination
func (r *resourceRepository) List(query *ResourceQuery) ([]*model.Resource, int64, error) {
	var resources []*model.Resource
	var total int64
	
	db := r.db.Model(&model.Resource{})
	
	// Apply filters
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.ParentID != nil {
		db = db.Where("parent_id = ?", *query.ParentID)
	}
	if query.Status >= 0 {
		db = db.Where("status = ?", query.Status)
	}
	
	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}
	
	// Get results
	if err := db.Find(&resources).Error; err != nil {
		return nil, 0, err
	}
	
	return resources, total, nil
}

// GetChildren gets child resources
func (r *resourceRepository) GetChildren(parentID uint) ([]*model.Resource, error) {
	var resources []*model.Resource
	err := r.db.Where("parent_id = ?", parentID).Find(&resources).Error
	return resources, err
}

// actionRepository implements ActionRepository
type actionRepository struct {
	db *gorm.DB
}

// NewActionRepository creates a new action repository
func NewActionRepository() ActionRepository {
	return &actionRepository{
		db: database.GetDB(),
	}
}

// Create creates a new action
func (r *actionRepository) Create(action *model.Action) error {
	return r.db.Create(action).Error
}

// Update updates an action
func (r *actionRepository) Update(action *model.Action) error {
	return r.db.Save(action).Error
}

// Delete deletes an action
func (r *actionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Action{}, id).Error
}

// GetByID gets an action by ID
func (r *actionRepository) GetByID(id uint) (*model.Action, error) {
	var action model.Action
	err := r.db.First(&action, id).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// GetByName gets an action by name
func (r *actionRepository) GetByName(name string) (*model.Action, error) {
	var action model.Action
	err := r.db.Where("name = ?", name).First(&action).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// List lists actions with pagination
func (r *actionRepository) List(query *ActionQuery) ([]*model.Action, int64, error) {
	var actions []*model.Action
	var total int64
	
	db := r.db.Model(&model.Action{})
	
	// Apply filters
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	
	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}
	
	// Get results
	if err := db.Find(&actions).Error; err != nil {
		return nil, 0, err
	}
	
	return actions, total, nil
}

// permissionRepository implements PermissionRepository
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository() PermissionRepository {
	return &permissionRepository{
		db: database.GetDB(),
	}
}

// Create creates a new permission
func (r *permissionRepository) Create(permission *model.PermissionExtended) error {
	return r.db.Create(permission).Error
}

// Update updates a permission
func (r *permissionRepository) Update(permission *model.PermissionExtended) error {
	return r.db.Save(permission).Error
}

// Delete deletes a permission
func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.PermissionExtended{}, id).Error
}

// GetByID gets a permission by ID
func (r *permissionRepository) GetByID(id uint) (*model.PermissionExtended, error) {
	var permission model.PermissionExtended
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByRoleResourceAction gets a permission by role, resource, and action
func (r *permissionRepository) GetByRoleResourceAction(roleID, resourceID, actionID uint) (*model.PermissionExtended, error) {
	var permission model.PermissionExtended
	err := r.db.Where("role_id = ? AND resource_id = ? AND action_id = ?", roleID, resourceID, actionID).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByRoleID gets permissions by role ID
func (r *permissionRepository) GetByRoleID(roleID uint) ([]*model.PermissionExtended, error) {
	var permissions []*model.PermissionExtended
	err := r.db.Where("role_id = ?", roleID).Find(&permissions).Error
	return permissions, err
}

// GetByResourceID gets permissions by resource ID
func (r *permissionRepository) GetByResourceID(resourceID uint) ([]*model.PermissionExtended, error) {
	var permissions []*model.PermissionExtended
	err := r.db.Where("resource_id = ?", resourceID).Find(&permissions).Error
	return permissions, err
}

// GetByActionID gets permissions by action ID
func (r *permissionRepository) GetByActionID(actionID uint) ([]*model.PermissionExtended, error) {
	var permissions []*model.PermissionExtended
	err := r.db.Where("action_id = ?", actionID).Find(&permissions).Error
	return permissions, err
}

// CreateUserAttribute creates a user attribute
func (r *permissionRepository) CreateUserAttribute(attribute *model.UserAttribute) error {
	return r.db.Create(attribute).Error
}

// UpdateUserAttribute updates a user attribute
func (r *permissionRepository) UpdateUserAttribute(attribute *model.UserAttribute) error {
	return r.db.Save(attribute).Error
}

// GetUserAttribute gets a user attribute
func (r *permissionRepository) GetUserAttribute(userID uint, key string) (*model.UserAttribute, error) {
	var attribute model.UserAttribute
	err := r.db.Where("user_id = ? AND key = ?", userID, key).First(&attribute).Error
	if err != nil {
		return nil, err
	}
	return &attribute, nil
}

// GetUserAttributes gets all attributes for a user
func (r *permissionRepository) GetUserAttributes(userID uint) (map[string]interface{}, error) {
	var attributes []*model.UserAttribute
	err := r.db.Where("user_id = ?", userID).Find(&attributes).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, attr := range attributes {
		// Convert value based on type
		switch attr.Type {
		case "number":
			var numValue float64
			if err := json.Unmarshal([]byte(attr.Value), &numValue); err == nil {
				result[attr.Key] = numValue
			} else {
				result[attr.Key] = attr.Value
			}
		case "boolean":
			var boolValue bool
			if err := json.Unmarshal([]byte(attr.Value), &boolValue); err == nil {
				result[attr.Key] = boolValue
			} else {
				result[attr.Key] = attr.Value
			}
		default:
			result[attr.Key] = attr.Value
		}
	}

	return result, nil
}

// CreateResourceAttribute creates a resource attribute
func (r *permissionRepository) CreateResourceAttribute(attribute *model.ResourceAttribute) error {
	return r.db.Create(attribute).Error
}

// UpdateResourceAttribute updates a resource attribute
func (r *permissionRepository) UpdateResourceAttribute(attribute *model.ResourceAttribute) error {
	return r.db.Save(attribute).Error
}

// GetResourceAttribute gets a resource attribute
func (r *permissionRepository) GetResourceAttribute(resourceID uint, key string) (*model.ResourceAttribute, error) {
	var attribute model.ResourceAttribute
	err := r.db.Where("resource_id = ? AND key = ?", resourceID, key).First(&attribute).Error
	if err != nil {
		return nil, err
	}
	return &attribute, nil
}

// GetResourceAttributes gets all attributes for a resource
func (r *permissionRepository) GetResourceAttributes(resourceID uint) (map[string]interface{}, error) {
	var attributes []*model.ResourceAttribute
	err := r.db.Where("resource_id = ?", resourceID).Find(&attributes).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, attr := range attributes {
		// Convert value based on type
		switch attr.Type {
		case "number":
			var numValue float64
			if err := json.Unmarshal([]byte(attr.Value), &numValue); err == nil {
				result[attr.Key] = numValue
			} else {
				result[attr.Key] = attr.Value
			}
		case "boolean":
			var boolValue bool
			if err := json.Unmarshal([]byte(attr.Value), &boolValue); err == nil {
				result[attr.Key] = boolValue
			} else {
				result[attr.Key] = attr.Value
			}
		default:
			result[attr.Key] = attr.Value
		}
	}

	return result, nil
}

// CreateRoleHierarchy creates a role hierarchy relationship
func (r *permissionRepository) CreateRoleHierarchy(hierarchy *model.RoleHierarchy) error {
	return r.db.Create(hierarchy).Error
}

// DeleteRoleHierarchy deletes a role hierarchy relationship
func (r *permissionRepository) DeleteRoleHierarchy(id uint) error {
	return r.db.Delete(&model.RoleHierarchy{}, id).Error
}

// GetRoleHierarchy gets a role hierarchy relationship
func (r *permissionRepository) GetRoleHierarchy(parentID, childID uint) (*model.RoleHierarchy, error) {
	var hierarchy model.RoleHierarchy
	err := r.db.Where("parent_id = ? AND child_id = ?", parentID, childID).First(&hierarchy).Error
	if err != nil {
		return nil, err
	}
	return &hierarchy, nil
}

// GetRoleHierarchyByParent gets role hierarchy by parent ID
func (r *permissionRepository) GetRoleHierarchyByParent(parentID uint) ([]*model.RoleHierarchy, error) {
	var hierarchies []*model.RoleHierarchy
	err := r.db.Where("parent_id = ?", parentID).Find(&hierarchies).Error
	return hierarchies, err
}

// GetRoleHierarchyByChild gets role hierarchy by child ID
func (r *permissionRepository) GetRoleHierarchyByChild(childID uint) ([]*model.RoleHierarchy, error) {
	var hierarchies []*model.RoleHierarchy
	err := r.db.Where("child_id = ?", childID).Find(&hierarchies).Error
	return hierarchies, err
}

// CreateAuditLog creates an audit log entry
func (r *permissionRepository) CreateAuditLog(log *model.PermissionAuditLog) error {
	return r.db.Create(log).Error
}

// GetAuditLogs gets audit logs for a user
func (r *permissionRepository) GetAuditLogs(userID uint, limit int) ([]*model.PermissionAuditLog, error) {
	var logs []*model.PermissionAuditLog
	db := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		db = db.Limit(limit)
	}
	err := db.Find(&logs).Error
	return logs, err
}

// CheckUserPermission checks if a user has permission for a resource and action
func (r *permissionRepository) CheckUserPermission(userID uint, resource, action string) (bool, error) {
	var count int64
	
	// Query to check if user has direct permission
	query := `
		SELECT COUNT(*) FROM permissions p
		JOIN user_roles ur ON p.role_id = ur.role_id
		JOIN resources res ON p.resource_id = res.id
		JOIN actions act ON p.action_id = act.id
		WHERE ur.user_id = ? AND res.name = ? AND act.name = ? AND p.status = 1
	`
	
	err := r.db.Raw(query, userID, resource, action).Scan(&count).Error
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// CreatePermission creates a new permission (for Permission model)
func (r *permissionRepository) CreatePermission(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

// UpdatePermission updates a permission (for Permission model)
func (r *permissionRepository) UpdatePermission(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

// DeletePermission deletes a permission (for Permission model)
func (r *permissionRepository) DeletePermission(permission *model.Permission) error {
	return r.db.Delete(permission).Error
}

// GetPermissionByID gets a permission by ID (for Permission model)
func (r *permissionRepository) GetPermissionByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetPermissionByName gets a permission by name (for Permission model)
func (r *permissionRepository) GetPermissionByName(name string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// ListPermissions lists permissions with pagination (for Permission model)
func (r *permissionRepository) ListPermissions(page, pageSize int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	// Count total records
	r.db.Model(&model.Permission{}).Count(&total)

	// Get paginated results
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&permissions).Error
	
	return permissions, total, err
}

// AssignPermissionToRole assigns a permission to a role
func (r *permissionRepository) AssignPermissionToRole(roleID, permissionID uint) error {
	rolePermission := &model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.Create(rolePermission).Error
}

// RemovePermissionFromRole removes a permission from a role
func (r *permissionRepository) RemovePermissionFromRole(roleID, permissionID uint) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&model.RolePermission{}).Error
}

// GetPermissionsByRoleID gets permissions assigned to a role
func (r *permissionRepository) GetPermissionsByRoleID(roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	
	err := r.db.Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Find(&permissions).Error
	
	return permissions, err
}

// GetPermissionsByUserID gets permissions for a user through their roles
func (r *permissionRepository) GetPermissionsByUserID(userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	
	err := r.db.Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Distinct("p.*").
		Find(&permissions).Error
	
	return permissions, err
}