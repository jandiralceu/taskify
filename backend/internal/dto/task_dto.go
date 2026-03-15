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
	IsBlocked      bool                `json:"isBlocked" binding:"omitempty"`
	AssignedTo     *uuid.UUID          `json:"assignedTo" binding:"omitempty"`
	DueDate        *time.Time          `json:"dueDate" binding:"omitempty"`
	EstimatedHours *float64            `json:"estimatedHours" binding:"omitempty,min=0"`
}

type UpdateTaskRequest struct {
	Title          *string              `json:"title" binding:"omitempty,min=3,max=255"`
	Description    *string              `json:"description" binding:"omitempty"`
	Status         *models.TaskStatus   `json:"status" binding:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority       *models.TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high critical"`
	IsBlocked      *bool                `json:"isBlocked" binding:"omitempty"`
	AssignedTo     *uuid.UUID           `json:"assignedTo" binding:"omitempty"`
	DueDate        *time.Time           `json:"dueDate" binding:"omitempty"`
	EstimatedHours *float64             `json:"estimatedHours" binding:"omitempty,min=0"`
	ActualHours    *float64             `json:"actualHours" binding:"omitempty,min=0"`
	IsArchived     *bool                `json:"isArchived" binding:"omitempty"`
}

type TaskResponse struct {
	ID             uuid.UUID           `json:"id"`
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	Status         models.TaskStatus   `json:"status"`
	Priority       models.TaskPriority `json:"priority"`
	IsBlocked      bool                `json:"isBlocked"`
	CreatedBy      uuid.UUID           `json:"createdBy"`
	AssignedTo     *uuid.UUID          `json:"assignedTo"`
	DueDate        *time.Time          `json:"dueDate"`
	CompletedAt    *time.Time          `json:"completedAt"`
	EstimatedHours *float64            `json:"estimatedHours"`
	ActualHours    *float64            `json:"actualHours"`
	IsArchived     bool                `json:"isArchived"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
}

type CreateTaskNoteRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateTaskNoteRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type TaskNoteResponse struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"taskId"`
	UserID    uuid.UUID `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetTaskListRequest struct {
	Status     models.TaskStatus   `form:"status" binding:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority   models.TaskPriority `form:"priority" binding:"omitempty,oneof=low medium high critical"`
	IsBlocked  *bool               `form:"isBlocked" binding:"omitempty"`
	IsArchived *bool               `form:"isArchived" binding:"omitempty"`
	CreatedBy  *uuid.UUID          `form:"createdBy" binding:"omitempty"`
	AssignedTo *uuid.UUID          `form:"assignedTo" binding:"omitempty"`
	Search     string              `form:"search" binding:"omitempty"`
	Sort       string              `form:"sort" binding:"omitempty"`
	Order      string              `form:"order" binding:"omitempty,oneof=asc desc"`
}
