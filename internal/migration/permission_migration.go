package migration

import (
	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// MigratePermissionTables creates the permission-related tables
func MigratePermissionTables() error {
	db := database.GetDB()
	
	// Auto migrate the new permission tables
	err := db.AutoMigrate(
		&model.Resource{},
		&model.Action{},
		&model.RoleHierarchy{},
		&model.PermissionExtended{},
		&model.UserAttribute{},
		&model.ResourceAttribute{},
		&model.PermissionAuditLog{},
	)
	
	if err != nil {
		return err
	}

	// Create indexes for better performance
	if err := createPermissionIndexes(db); err != nil {
		return err
	}

	// Insert default data
	if err := insertDefaultPermissionData(db); err != nil {
		return err
	}

	return nil
}

// createPermissionIndexes creates indexes for permission tables
func createPermissionIndexes(db *gorm.DB) error {
	// Resource indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_type ON resources(type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_parent_id ON resources(parent_id)")
	
	// Action indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_actions_category ON actions(category)")
	
	// Role hierarchy indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_role_hierarchies_parent_child ON role_hierarchies(parent_id, child_id)")
	
	// Extended permission indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_permissions_extended_role_resource_action ON permissions_extended(role_id, resource_id, action_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_permissions_extended_status ON permissions_extended(status)")
	
	// User attribute indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_user_attributes_user_key ON user_attributes(user_id, key)")
	
	// Resource attribute indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resource_attributes_resource_key ON resource_attributes(resource_id, key)")
	
	// Audit log indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_permission_audit_logs_user_id ON permission_audit_logs(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_permission_audit_logs_operation ON permission_audit_logs(operation)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_permission_audit_logs_created_at ON permission_audit_logs(created_at)")
	
	return nil
}

// insertDefaultPermissionData inserts default permission data
func insertDefaultPermissionData(db *gorm.DB) error {
	// Insert default actions
	defaultActions := []model.Action{
		{Name: "create", Description: "Create resource", Category: "crud"},
		{Name: "read", Description: "Read resource", Category: "crud"},
		{Name: "update", Description: "Update resource", Category: "crud"},
		{Name: "delete", Description: "Delete resource", Category: "crud"},
		{Name: "manage", Description: "Manage resource", Category: "system"},
		{Name: "approve", Description: "Approve resource", Category: "business"},
		{Name: "reject", Description: "Reject resource", Category: "business"},
	}
	
	for _, action := range defaultActions {
		var existing model.Action
		if err := db.Where("name = ?", action.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&action).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	
	// Insert default resources
	defaultResources := []model.Resource{
		{Name: "user", Description: "User management", Type: "system", Path: "/users"},
		{Name: "role", Description: "Role management", Type: "system", Path: "/roles"},
		{Name: "permission", Description: "Permission management", Type: "system", Path: "/permissions"},
		{Name: "resource", Description: "Resource management", Type: "system", Path: "/resources"},
		{Name: "audit", Description: "Audit logs", Type: "system", Path: "/audit"},
	}
	
	for _, resource := range defaultResources {
		var existing model.Resource
		if err := db.Where("name = ?", resource.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&resource).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	
	return nil
}