package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-admin/internal/logger"
	"go-admin/internal/model"
	"go-admin/internal/repository"

	"github.com/Knetic/govaluate"
	"go.uber.org/zap"
)

// PermissionService defines the permission service interface
type PermissionService interface {
	// Resource management
	CreateResource(ctx context.Context, resource *model.Resource) error
	UpdateResource(ctx context.Context, resource *model.Resource) error
	DeleteResource(ctx context.Context, id uint) error
	GetResource(ctx context.Context, id uint) (*model.Resource, error)
	GetResourceByName(ctx context.Context, name string) (*model.Resource, error)
	ListResources(ctx context.Context, query *repository.ResourceQuery) ([]*model.Resource, int64, error)

	// Action management
	CreateAction(ctx context.Context, action *model.Action) error
	UpdateAction(ctx context.Context, action *model.Action) error
	DeleteAction(ctx context.Context, id uint) error
	GetAction(ctx context.Context, id uint) (*model.Action, error)
	GetActionByName(ctx context.Context, name string) (*model.Action, error)
	ListActions(ctx context.Context, query *repository.ActionQuery) ([]*model.Action, int64, error)

	// Permission management (for Permission model - used by handler)
	CreatePermission(name, description, resource, action string) (*model.Permission, error)
	GetPermissionByID(id uint) (*model.Permission, error)
	UpdatePermission(permission *model.Permission) error
	DeletePermission(id uint) error
	ListPermissions(page, pageSize int) ([]*model.Permission, int64, error)
	AssignPermissionToRole(roleID, permissionID uint) error
	RemovePermissionFromRole(roleID, permissionID uint) error
	GetPermissionsByRoleID(roleID uint) ([]*model.Permission, error)
	GetPermissionsByUserID(userID uint) ([]*model.Permission, error)

	// Permission management (for PermissionExtended model - role-based permissions)
	GrantPermission(ctx context.Context, roleID, resourceID, actionID uint, conditions *model.PermissionCondition) error
	RevokePermission(ctx context.Context, roleID, resourceID, actionID uint) error
	CheckPermission(ctx context.Context, userID uint, resource, action string, context map[string]interface{}) (bool, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]*PermissionInfo, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]*PermissionInfo, error)

	// Role hierarchy
	AddRoleInheritance(ctx context.Context, parentID, childID uint) error
	RemoveRoleInheritance(ctx context.Context, parentID, childID uint) error
	GetRoleHierarchy(ctx context.Context, roleID uint) ([]*model.Role, error)
	GetRoleChildren(ctx context.Context, roleID uint) ([]*model.Role, error)

	// Attributes management
	SetUserAttribute(ctx context.Context, userID uint, key, value, attrType string) error
	GetUserAttributes(ctx context.Context, userID uint) (map[string]interface{}, error)
	SetResourceAttribute(ctx context.Context, resourceID uint, key, value, attrType string) error
	GetResourceAttributes(ctx context.Context, resourceID uint) (map[string]interface{}, error)

	// Audit logging
	LogPermissionCheck(ctx context.Context, userID uint, resource, action string, result bool, reason string, context map[string]interface{}) error
}

// PermissionInfo represents permission information
type PermissionInfo struct {
	Resource   *model.Resource            `json:"resource"`
	Action     *model.Action              `json:"action"`
	Conditions *model.PermissionCondition `json:"conditions,omitempty"`
	Priority   int                        `json:"priority"`
}

// PermissionCheckContext represents the context for permission checking
type PermissionCheckContext struct {
	User        *model.User
	Resource    *model.Resource
	Action      *model.Action
	UserAttrs   map[string]interface{}
	ResAttrs    map[string]interface{}
	Environment map[string]interface{}
}

// permissionService implements PermissionService
type permissionService struct {
	resourceRepo   repository.ResourceRepository
	actionRepo     repository.ActionRepository
	permissionRepo repository.PermissionRepository
	roleRepo       repository.RoleRepository
	userRepo       repository.UserRepository
}

// NewPermissionService creates a new permission service
func NewPermissionService() PermissionService {
	return &permissionService{
		resourceRepo:   repository.NewResourceRepository(),
		actionRepo:     repository.NewActionRepository(),
		permissionRepo: repository.NewPermissionRepository(),
		roleRepo:       repository.NewRoleRepository(),
		userRepo:       repository.NewUserRepository(),
	}
}

// CreateResource creates a new resource
func (s *permissionService) CreateResource(ctx context.Context, resource *model.Resource) error {
	if resource.Name == "" {
		return fmt.Errorf("resource name cannot be empty")
	}

	// Check if resource already exists
	existing, err := s.resourceRepo.GetByName(resource.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("resource with name %s already exists", resource.Name)
	}

	return s.resourceRepo.Create(resource)
}

// UpdateResource updates a resource
func (s *permissionService) UpdateResource(ctx context.Context, resource *model.Resource) error {
	if resource.ID == 0 {
		return fmt.Errorf("resource ID cannot be empty")
	}

	return s.resourceRepo.Update(resource)
}

// DeleteResource deletes a resource
func (s *permissionService) DeleteResource(ctx context.Context, id uint) error {
	// Check if resource has child resources
	children, err := s.resourceRepo.GetChildren(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return fmt.Errorf("cannot delete resource with child resources")
	}

	// Check if resource has permissions
	permissions, err := s.permissionRepo.GetByResourceID(id)
	if err != nil {
		return err
	}
	if len(permissions) > 0 {
		return fmt.Errorf("cannot delete resource with associated permissions")
	}

	return s.resourceRepo.Delete(id)
}

// GetResource gets a resource by ID
func (s *permissionService) GetResource(ctx context.Context, id uint) (*model.Resource, error) {
	return s.resourceRepo.GetByID(id)
}

// GetResourceByName gets a resource by name
func (s *permissionService) GetResourceByName(ctx context.Context, name string) (*model.Resource, error) {
	return s.resourceRepo.GetByName(name)
}

// ListResources lists resources with pagination
func (s *permissionService) ListResources(ctx context.Context, query *repository.ResourceQuery) ([]*model.Resource, int64, error) {
	return s.resourceRepo.List(query)
}

// CreateAction creates a new action
func (s *permissionService) CreateAction(ctx context.Context, action *model.Action) error {
	if action.Name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	// Check if action already exists
	existing, err := s.actionRepo.GetByName(action.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("action with name %s already exists", action.Name)
	}

	return s.actionRepo.Create(action)
}

// UpdateAction updates an action
func (s *permissionService) UpdateAction(ctx context.Context, action *model.Action) error {
	if action.ID == 0 {
		return fmt.Errorf("action ID cannot be empty")
	}

	return s.actionRepo.Update(action)
}

// DeleteAction deletes an action
func (s *permissionService) DeleteAction(ctx context.Context, id uint) error {
	// Check if action has permissions
	permissions, err := s.permissionRepo.GetByActionID(id)
	if err != nil {
		return err
	}
	if len(permissions) > 0 {
		return fmt.Errorf("cannot delete action with associated permissions")
	}

	return s.actionRepo.Delete(id)
}

// GetAction gets an action by ID
func (s *permissionService) GetAction(ctx context.Context, id uint) (*model.Action, error) {
	return s.actionRepo.GetByID(id)
}

// GetActionByName gets an action by name
func (s *permissionService) GetActionByName(ctx context.Context, name string) (*model.Action, error) {
	return s.actionRepo.GetByName(name)
}

// ListActions lists actions with pagination
func (s *permissionService) ListActions(ctx context.Context, query *repository.ActionQuery) ([]*model.Action, int64, error) {
	return s.actionRepo.List(query)
}

// GrantPermission grants a permission to a role
func (s *permissionService) GrantPermission(ctx context.Context, roleID, resourceID, actionID uint, conditions *model.PermissionCondition) error {
	// Check if role exists
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("role not found: %v", err)
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	// Check if resource exists
	resource, err := s.resourceRepo.GetByID(resourceID)
	if err != nil {
		return fmt.Errorf("resource not found: %v", err)
	}
	if resource == nil {
		return fmt.Errorf("resource not found")
	}

	// Check if action exists
	action, err := s.actionRepo.GetByID(actionID)
	if err != nil {
		return fmt.Errorf("action not found: %v", err)
	}
	if action == nil {
		return fmt.Errorf("action not found")
	}

	// Check if permission already exists
	existing, err := s.permissionRepo.GetByRoleResourceAction(roleID, resourceID, actionID)
	if err == nil && existing != nil {
		return fmt.Errorf("permission already exists for this role, resource, and action combination")
	}

	// Create permission
	permission := &model.PermissionExtended{
		RoleID:     roleID,
		ResourceID: resourceID,
		ActionID:   actionID,
		Status:     1,
	}

	if conditions != nil {
		conditionsStr, err := conditions.MarshalConditions()
		if err != nil {
			return fmt.Errorf("failed to marshal conditions: %v", err)
		}
		permission.Conditions = conditionsStr
	}

	return s.permissionRepo.Create(permission)
}

// RevokePermission revokes a permission from a role
func (s *permissionService) RevokePermission(ctx context.Context, roleID, resourceID, actionID uint) error {
	permission, err := s.permissionRepo.GetByRoleResourceAction(roleID, resourceID, actionID)
	if err != nil {
		return fmt.Errorf("permission not found: %v", err)
	}
	if permission == nil {
		return fmt.Errorf("permission not found")
	}

	return s.permissionRepo.Delete(permission.ID)
}

// CheckPermission checks if a user has permission to perform an action on a resource
func (s *permissionService) CheckPermission(ctx context.Context, userID uint, resource, action string, context map[string]interface{}) (bool, error) {
	// Get user information
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		return false, fmt.Errorf("user not found")
	}

	// Get resource information
	resourceObj, err := s.resourceRepo.GetByName(resource)
	if err != nil {
		return false, fmt.Errorf("failed to get resource: %v", err)
	}
	if resourceObj == nil {
		return false, fmt.Errorf("resource not found")
	}

	// Get action information
	actionObj, err := s.actionRepo.GetByName(action)
	if err != nil {
		return false, fmt.Errorf("failed to get action: %v", err)
	}
	if actionObj == nil {
		return false, fmt.Errorf("action not found")
	}

	// Get user roles
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %v", err)
	}

	if len(roles) == 0 {
		// Log permission denial
		s.LogPermissionCheck(ctx, userID, resource, action, false, "User has no roles", context)
		return false, nil
	}

	// Get user attributes
	userAttrs, err := s.permissionRepo.GetUserAttributes(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user attributes: %v", err)
	}

	// Get resource attributes
	resourceAttrs, err := s.permissionRepo.GetResourceAttributes(resourceObj.ID)
	if err != nil {
		return false, fmt.Errorf("failed to get resource attributes: %v", err)
	}

	// Check permissions for each role
	for _, role := range roles {
		// Get role permissions
		permissions, err := s.permissionRepo.GetByRoleID(role.ID)
		if err != nil {
			continue
		}

		// Check each permission
		for _, permission := range permissions {
			if permission.ResourceID == resourceObj.ID && permission.ActionID == actionObj.ID && permission.Status == 1 {
				// Check conditions if present
				if permission.Conditions != "" {
					var conditions model.PermissionCondition
					if err := conditions.UnmarshalConditions(permission.Conditions); err != nil {
						logger.Error("Failed to unmarshal permission conditions", zap.Error(err))
						continue
					}

					// Evaluate conditions
					if !s.evaluateConditions(&conditions, userAttrs, resourceAttrs, context) {
						continue
					}
				}

				// Permission granted
				s.LogPermissionCheck(ctx, userID, resource, action, true, "Permission granted", context)
				return true, nil
			}
		}
	}

	// Permission denied
	s.LogPermissionCheck(ctx, userID, resource, action, false, "No matching permission found", context)
	return false, nil
}

// evaluateConditions evaluates permission conditions
func (s *permissionService) evaluateConditions(conditions *model.PermissionCondition, userAttrs, resourceAttrs, env map[string]interface{}) bool {
	// Check resource attributes
	if conditions.ResourceAttributes != nil {
		for key, expectedValue := range conditions.ResourceAttributes {
			if actualValue, exists := resourceAttrs[key]; !exists || actualValue != expectedValue {
				return false
			}
		}
	}

	// Check user attributes
	if conditions.UserAttributes != nil {
		for key, expectedValue := range conditions.UserAttributes {
			if actualValue, exists := userAttrs[key]; !exists || actualValue != expectedValue {
				return false
			}
		}
	}

	// Check environment conditions
	if conditions.Environment != nil {
		for key, expectedValue := range conditions.Environment {
			if actualValue, exists := env[key]; !exists || actualValue != expectedValue {
				return false
			}
		}
	}

	// Check expression if present
	if conditions.Expression != "" {
		// Simple expression evaluation (can be extended with a proper expression engine)
		return s.evaluateExpression(conditions.Expression, userAttrs, resourceAttrs, env)
	}

	return true
}

// evaluateExpression evaluates a simple expression using govaluate
func (s *permissionService) evaluateExpression(expression string, userAttrs, resourceAttrs, env map[string]interface{}) bool {
	// This is a simple implementation - in production, use a proper expression engine
	// For now, just return true for simple expressions
	if strings.TrimSpace(expression) == "" {
		return true
	}

	// Create parameters for the expression
	parameters := make(map[string]interface{})

	// Add user attributes with user. prefix
	for key, value := range userAttrs {
		parameters["user."+key] = value
	}

	// Add resource attributes with resource. prefix
	for key, value := range resourceAttrs {
		parameters["resource."+key] = value
	}

	// Add environment variables with env. prefix
	for key, value := range env {
		parameters["env."+key] = value
	}

	// Create the expression
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		logger.Error("Failed to parse expression",
			zap.String("expression", expression),
			zap.Error(err))
		// Default to deny if expression is invalid
		return false
	}

	// Evaluate the expression
	result, err := expr.Evaluate(parameters)
	if err != nil {
		logger.Error("Failed to evaluate expression",
			zap.String("expression", expression),
			zap.Error(err))
		// Default to deny if evaluation fails
		return false
	}

	// Convert result to boolean
	boolResult, ok := result.(bool)
	if !ok {
		logger.Warn("Expression did not evaluate to boolean",
			zap.String("expression", expression),
			zap.Any("result", result))
		// Default to deny if result is not boolean
		return false
	}

	logger.Debug("Expression evaluation result",
		zap.String("expression", expression),
		zap.Bool("result", boolResult))

	return boolResult
}

// GetUserPermissions gets all permissions for a user
func (s *permissionService) GetUserPermissions(ctx context.Context, userID uint) ([]*PermissionInfo, error) {
	// Get user roles
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %v", err)
	}

	var permissions []*PermissionInfo

	for _, role := range roles {
		rolePerms, err := s.GetRolePermissions(ctx, role.ID)
		if err != nil {
			continue
		}
		permissions = append(permissions, rolePerms...)
	}

	return permissions, nil
}

// GetRolePermissions gets all permissions for a role
func (s *permissionService) GetRolePermissions(ctx context.Context, roleID uint) ([]*PermissionInfo, error) {
	// Get permissions
	permissions, err := s.permissionRepo.GetByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %v", err)
	}

	var result []*PermissionInfo

	for _, permission := range permissions {
		if permission.Status != 1 {
			continue
		}

		// Get resource
		resource, err := s.resourceRepo.GetByID(permission.ResourceID)
		if err != nil {
			continue
		}

		// Get action
		action, err := s.actionRepo.GetByID(permission.ActionID)
		if err != nil {
			continue
		}

		info := &PermissionInfo{
			Resource: resource,
			Action:   action,
			Priority: permission.Priority,
		}

		// Parse conditions if present
		if permission.Conditions != "" {
			var conditions model.PermissionCondition
			if err := conditions.UnmarshalConditions(permission.Conditions); err == nil {
				info.Conditions = &conditions
			}
		}

		result = append(result, info)
	}

	return result, nil
}

// AddRoleInheritance adds a role inheritance relationship
func (s *permissionService) AddRoleInheritance(ctx context.Context, parentID, childID uint) error {
	// Check if both roles exist
	parentRole, err := s.roleRepo.GetByID(parentID)
	if err != nil || parentRole == nil {
		return fmt.Errorf("parent role not found")
	}

	childRole, err := s.roleRepo.GetByID(childID)
	if err != nil || childRole == nil {
		return fmt.Errorf("child role not found")
	}

	// Check if relationship already exists
	existing, err := s.permissionRepo.GetRoleHierarchy(parentID, childID)
	if err == nil && existing != nil {
		return fmt.Errorf("role inheritance relationship already exists")
	}

	// Create relationship
	hierarchy := &model.RoleHierarchy{
		ParentID: parentID,
		ChildID:  childID,
	}

	return s.permissionRepo.CreateRoleHierarchy(hierarchy)
}

// RemoveRoleInheritance removes a role inheritance relationship
func (s *permissionService) RemoveRoleInheritance(ctx context.Context, parentID, childID uint) error {
	hierarchy, err := s.permissionRepo.GetRoleHierarchy(parentID, childID)
	if err != nil {
		return fmt.Errorf("role inheritance relationship not found: %v", err)
	}
	if hierarchy == nil {
		return fmt.Errorf("role inheritance relationship not found")
	}

	return s.permissionRepo.DeleteRoleHierarchy(hierarchy.ID)
}

// GetRoleHierarchy gets the role hierarchy for a role
func (s *permissionService) GetRoleHierarchy(ctx context.Context, roleID uint) ([]*model.Role, error) {
	return s.roleRepo.GetRoleHierarchy(roleID)
}

// GetRoleChildren gets the child roles for a role
func (s *permissionService) GetRoleChildren(ctx context.Context, roleID uint) ([]*model.Role, error) {
	return s.roleRepo.GetRoleChildren(roleID)
}

// SetUserAttribute sets a user attribute
func (s *permissionService) SetUserAttribute(ctx context.Context, userID uint, key, value, attrType string) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	// Check if attribute already exists
	existing, err := s.permissionRepo.GetUserAttribute(userID, key)
	if err == nil && existing != nil {
		// Update existing attribute
		existing.Value = value
		existing.Type = attrType
		return s.permissionRepo.UpdateUserAttribute(existing)
	}

	// Create new attribute
	attribute := &model.UserAttribute{
		UserID: userID,
		Key:    key,
		Value:  value,
		Type:   attrType,
	}

	return s.permissionRepo.CreateUserAttribute(attribute)
}

// GetUserAttributes gets all attributes for a user
func (s *permissionService) GetUserAttributes(ctx context.Context, userID uint) (map[string]interface{}, error) {
	return s.permissionRepo.GetUserAttributes(userID)
}

// SetResourceAttribute sets a resource attribute
func (s *permissionService) SetResourceAttribute(ctx context.Context, resourceID uint, key, value, attrType string) error {
	// Check if resource exists
	resource, err := s.resourceRepo.GetByID(resourceID)
	if err != nil || resource == nil {
		return fmt.Errorf("resource not found")
	}

	// Check if attribute already exists
	existing, err := s.permissionRepo.GetResourceAttribute(resourceID, key)
	if err == nil && existing != nil {
		// Update existing attribute
		existing.Value = value
		existing.Type = attrType
		return s.permissionRepo.UpdateResourceAttribute(existing)
	}

	// Create new attribute
	attribute := &model.ResourceAttribute{
		ResourceID: resourceID,
		Key:        key,
		Value:      value,
		Type:       attrType,
	}

	return s.permissionRepo.CreateResourceAttribute(attribute)
}

// GetResourceAttributes gets all attributes for a resource
func (s *permissionService) GetResourceAttributes(ctx context.Context, resourceID uint) (map[string]interface{}, error) {
	return s.permissionRepo.GetResourceAttributes(resourceID)
}

// LogPermissionCheck logs a permission check operation
func (s *permissionService) LogPermissionCheck(ctx context.Context, userID uint, resource, action string, result bool, reason string, context map[string]interface{}) error {
	// Get resource ID
	resourceObj, err := s.resourceRepo.GetByName(resource)
	if err != nil {
		resourceObj = &model.Resource{ID: 0}
	}

	// Get action ID
	actionObj, err := s.actionRepo.GetByName(action)
	if err != nil {
		actionObj = &model.Action{ID: 0}
	}

	// Create audit log
	auditLog := &model.PermissionAuditLog{
		UserID:     userID,
		ResourceID: resourceObj.ID,
		ActionID:   actionObj.ID,
		Operation:  "check",
		Result:     result,
		Reason:     reason,
	}

	// Add context if available
	if context != nil {
		contextStr, err := json.Marshal(context)
		if err == nil {
			auditLog.Context = string(contextStr)
		}
	}

	return s.permissionRepo.CreateAuditLog(auditLog)
}

// CreatePermission creates a new permission (for Permission model)
func (s *permissionService) CreatePermission(name, description, resource, action string) (*model.Permission, error) {
	if name == "" {
		return nil, fmt.Errorf("permission name cannot be empty")
	}
	if resource == "" {
		return nil, fmt.Errorf("resource cannot be empty")
	}
	if action == "" {
		return nil, fmt.Errorf("action cannot be empty")
	}

	// Check if permission already exists
	existing, err := s.permissionRepo.GetPermissionByName(name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("permission with name %s already exists", name)
	}

	// Create permission
	permission := &model.Permission{
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
	}

	err = s.permissionRepo.CreatePermission(permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// GetPermissionByID gets a permission by ID (for Permission model)
func (s *permissionService) GetPermissionByID(id uint) (*model.Permission, error) {
	return s.permissionRepo.GetPermissionByID(id)
}

// UpdatePermission updates a permission (for Permission model)
func (s *permissionService) UpdatePermission(permission *model.Permission) error {
	if permission.ID == 0 {
		return fmt.Errorf("permission ID cannot be empty")
	}
	if permission.Name == "" {
		return fmt.Errorf("permission name cannot be empty")
	}
	if permission.Resource == "" {
		return fmt.Errorf("resource cannot be empty")
	}
	if permission.Action == "" {
		return fmt.Errorf("action cannot be empty")
	}

	return s.permissionRepo.UpdatePermission(permission)
}

// DeletePermission deletes a permission (for Permission model)
func (s *permissionService) DeletePermission(id uint) error {
	// Get the permission first
	permission, err := s.permissionRepo.GetPermissionByID(id)
	if err != nil {
		return fmt.Errorf("permission not found: %v", err)
	}
	if permission == nil {
		return fmt.Errorf("permission not found")
	}

	// Check if permission is assigned to any roles
	rolePermissions, err := s.permissionRepo.GetPermissionsByRoleID(0) // This will get all role-permission relationships
	if err != nil {
		return err
	}

	// Check if this permission is assigned to any role
	for _, rp := range rolePermissions {
		if rp.ID == id {
			return fmt.Errorf("cannot delete permission that is assigned to roles")
		}
	}

	return s.permissionRepo.DeletePermission(permission)
}

// ListPermissions lists permissions with pagination (for Permission model)
func (s *permissionService) ListPermissions(page, pageSize int) ([]*model.Permission, int64, error) {
	return s.permissionRepo.ListPermissions(page, pageSize)
}

// AssignPermissionToRole assigns a permission to a role
func (s *permissionService) AssignPermissionToRole(roleID, permissionID uint) error {
	// Check if role exists
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil || role == nil {
		return fmt.Errorf("role not found")
	}

	// Check if permission exists
	permission, err := s.permissionRepo.GetPermissionByID(permissionID)
	if err != nil || permission == nil {
		return fmt.Errorf("permission not found")
	}

	// Check if assignment already exists
	permissions, err := s.permissionRepo.GetPermissionsByRoleID(roleID)
	if err != nil {
		return err
	}

	for _, p := range permissions {
		if p.ID == permissionID {
			return fmt.Errorf("permission is already assigned to this role")
		}
	}

	return s.permissionRepo.AssignPermissionToRole(roleID, permissionID)
}

// RemovePermissionFromRole removes a permission from a role
func (s *permissionService) RemovePermissionFromRole(roleID, permissionID uint) error {
	return s.permissionRepo.RemovePermissionFromRole(roleID, permissionID)
}

// GetPermissionsByRoleID gets permissions assigned to a role
func (s *permissionService) GetPermissionsByRoleID(roleID uint) ([]*model.Permission, error) {
	return s.permissionRepo.GetPermissionsByRoleID(roleID)
}

// GetPermissionsByUserID gets permissions for a user (through their roles)
func (s *permissionService) GetPermissionsByUserID(userID uint) ([]*model.Permission, error) {
	// Get user roles
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %v", err)
	}

	var allPermissions []*model.Permission
	permissionMap := make(map[uint]*model.Permission) // To avoid duplicates

	for _, role := range roles {
		permissions, err := s.permissionRepo.GetPermissionsByRoleID(role.ID)
		if err != nil {
			continue
		}

		for _, permission := range permissions {
			if _, exists := permissionMap[permission.ID]; !exists {
				permissionMap[permission.ID] = permission
				allPermissions = append(allPermissions, permission)
			}
		}
	}

	return allPermissions, nil
}
