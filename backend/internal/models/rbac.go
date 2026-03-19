package models

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user role in the system.
type RoleModel struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string       `gorm:"type:varchar(50);unique;not null" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time    `gorm:"type:timestamptz;not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time    `gorm:"type:timestamptz;not null;default:now()" json:"updatedAt"`
}

// TableName overrides the table name for RoleModel
func (RoleModel) TableName() string {
	return "roles"
}

// Permission represents a specific action on a resource.
type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Resource    string    `gorm:"type:varchar(100);not null" json:"resource"`
	Action      string    `gorm:"type:varchar(50);not null" json:"action"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updatedAt"`
}

// CasbinRule represents a rule in Casbin's storage.
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"type:varchar(100);index"`
	V0    string `gorm:"type:varchar(100);index"`
	V1    string `gorm:"type:varchar(100);index"`
	V2    string `gorm:"type:varchar(100);index"`
	V3    string `gorm:"type:varchar(100);index"`
	V4    string `gorm:"type:varchar(100);index"`
	V5    string `gorm:"type:varchar(100);index"`
}

// TableName overrides the table name for CasbinRule
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
