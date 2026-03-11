package dto

import (
	"github.com/jandiralceu/taskify/internal/models"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type CreateUserRequest struct {
	FirstName string      `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string      `json:"last_name" binding:"required,min=2,max=100"`
	Email     string      `json:"email" binding:"required,email,max=255"`
	Password  string      `json:"password" binding:"required,min=8"`
	Role      models.Role `json:"role" binding:"required,oneof=admin employee"`
}

type GetUserListRequest struct {
	PaginationRequest
	FirstName string      `form:"first_name" binding:"omitempty"`
	LastName  string      `form:"last_name" binding:"omitempty"`
	Email     string      `form:"email" binding:"omitempty,email"`
	Role      models.Role `form:"role" binding:"omitempty,oneof=admin employee"`
}

type UserListResponse PaginatedResponse[models.User]

type UpdateUserRequest struct {
	FirstName *string `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" binding:"omitempty,min=2,max=100"`
	IsActive  *bool   `json:"is_active" binding:"omitempty"`
}
