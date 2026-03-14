package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

// User represents a user account in the system.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName    string    `gorm:"type:varchar(100);not null" json:"firstName"`
	LastName     string    `gorm:"type:varchar(100);not null" json:"lastName"`
	Email        string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`
	Role         Role      `gorm:"type:user_role;not null;default:'employee'" json:"role"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	AvatarURL    *string   `gorm:"type:text" json:"avatarUrl"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updatedAt"`
}
