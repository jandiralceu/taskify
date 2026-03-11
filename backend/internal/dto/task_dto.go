package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/models"
)

type CreateTaskRequest struct {
	Title          string              `json:"title" binding:"required,min=3,max=255"`
	Description    string              `json:"description" binding:"omitempty"`
	Status         models.TaskStatus   `json:"status" binding:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority       models.TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high critical"`
	IsBlocked      bool                `json:"is_blocked" binding:"omitempty"`
	AssignedTo     *uuid.UUID          `json:"assigned_to" binding:"omitempty"`
	DueDate        *time.Time          `json:"due_date" binding:"omitempty"`
	EstimatedHours *float64            `json:"estimated_hours" binding:"omitempty,min=0"`
}

type UpdateTaskRequest struct {
	Title          *string              `json:"title" binding:"omitempty,min=3,max=255"`
	Description    *string              `json:"description" binding:"omitempty"`
	Status         *models.TaskStatus   `json:"status" binding:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority       *models.TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high critical"`
	IsBlocked      *bool                `json:"is_blocked" binding:"omitempty"`
	AssignedTo     *uuid.UUID           `json:"assigned_to" binding:"omitempty"`
	DueDate        *time.Time           `json:"due_date" binding:"omitempty"`
	EstimatedHours *float64             `json:"estimated_hours" binding:"omitempty,min=0"`
	ActualHours    *float64             `json:"actual_hours" binding:"omitempty,min=0"`
	IsArchived     *bool                `json:"is_archived" binding:"omitempty"`
}

type TaskResponse struct {
	ID             uuid.UUID           `json:"id"`
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	Status         models.TaskStatus   `json:"status"`
	Priority       models.TaskPriority `json:"priority"`
	IsBlocked      bool                `json:"is_blocked"`
	CreatedBy      uuid.UUID           `json:"created_by"`
	AssignedTo     *uuid.UUID          `json:"assigned_to"`
	DueDate        *time.Time          `json:"due_date"`
	CompletedAt    *time.Time          `json:"completed_at"`
	EstimatedHours *float64            `json:"estimated_hours"`
	ActualHours    *float64            `json:"actual_hours"`
	IsArchived     bool                `json:"is_archived"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type CreateTaskNoteRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateTaskNoteRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type TaskNoteResponse struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetTaskListRequest struct {
	PaginationRequest
	Status     models.TaskStatus   `form:"status" binding:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority   models.TaskPriority `form:"priority" binding:"omitempty,oneof=low medium high critical"`
	IsBlocked  *bool               `form:"is_blocked" binding:"omitempty"`
	IsArchived *bool               `form:"is_archived" binding:"omitempty"`
	CreatedBy  *uuid.UUID          `form:"created_by" binding:"omitempty"`
	AssignedTo *uuid.UUID          `form:"assigned_to" binding:"omitempty"`
	Search     string              `form:"search" binding:"omitempty"`
}
