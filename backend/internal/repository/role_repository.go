package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/models"
	"gorm.io/gorm"
)

// RoleRepository defines the interface for RBAC database operations.
type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*models.RoleModel, error)
	GetPermissionsByRole(ctx context.Context, roleName string) ([]models.Permission, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.RoleModel, error)
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
}

type roleRepository struct {
	db *gorm.DB
}

var _ RoleRepository = (*roleRepository)(nil)

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*models.RoleModel, error) {
	var role models.RoleModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &role, nil
}

func (r *roleRepository) GetPermissionsByRole(ctx context.Context, roleName string) ([]models.Permission, error) {
	var role models.RoleModel
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", roleName).First(&role).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return role.Permissions, nil
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.RoleModel, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return user.Roles, nil
}

func (r *roleRepository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	// Using a raw query or GORM Association to insert into user_roles
	sql := "INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING"
	if err := r.db.WithContext(ctx).Exec(sql, userID, roleID).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}
