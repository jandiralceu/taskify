package dto

import (
	"github.com/jandiralceu/taskify/internal/models"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type CreateUserRequest struct {
	FirstName string      `json:"firstName" binding:"required,min=2,max=100"`
	LastName  string      `json:"lastName" binding:"required,min=2,max=100"`
	Email     string      `json:"email" binding:"required,email,max=255"`
	Password  string      `json:"password" binding:"required,min=8"`
	Role      string `json:"role" binding:"required,oneof=admin employee"`
}

type GetUserListRequest struct {
	PaginationRequest
	FirstName string `form:"firstName" binding:"omitempty"`
	LastName  string `form:"lastName" binding:"omitempty"`
	Email     string `form:"email" binding:"omitempty,email"`
	Role      string `form:"role" binding:"omitempty,oneof=admin employee"`
}


type UserListResponse PaginatedResponse[models.User]

type UpdateUserRequest struct {
	FirstName *string `json:"firstName" binding:"omitempty,min=2,max=100"`
	LastName  *string `json:"lastName" binding:"omitempty,min=2,max=100"`
	IsActive  *bool   `json:"isActive" binding:"omitempty"`
}
