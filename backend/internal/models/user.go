package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system.
type User struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName    string      `gorm:"type:varchar(100);not null" json:"firstName"`
	LastName     string      `gorm:"type:varchar(100);not null" json:"lastName"`
	Email        string      `gorm:"type:varchar(255);not null;unique" json:"email"`
	PasswordHash string      `gorm:"type:text;not null" json:"-"`
	Roles        []RoleModel `gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;References:ID;JoinReferences:RoleID" json:"roles,omitempty"`
	Role         string      `gorm:"-" json:"role,omitempty"` // Virtual field for role specification on creation
	IsActive     bool        `gorm:"not null;default:true" json:"isActive"`

	AvatarURL    *string     `gorm:"type:text" json:"avatarUrl"`
	CreatedAt    time.Time   `gorm:"type:timestamptz;not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"type:timestamptz;not null;default:now()" json:"updatedAt"`
}
