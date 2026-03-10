package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
	TaskStatusBlocked    TaskStatus = "blocked"
)

type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

type Task struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title          string       `gorm:"type:varchar(255);not null" json:"title"`
	Description    string       `gorm:"type:text" json:"description"`
	Status         TaskStatus   `gorm:"type:task_status;not null;default:'pending'" json:"status"`
	Priority       TaskPriority `gorm:"type:task_priority;not null;default:'medium'" json:"priority"`
	CreatedBy      uuid.UUID    `gorm:"type:uuid;not null" json:"created_by"`
	AssignedTo     *uuid.UUID   `gorm:"type:uuid" json:"assigned_to"`
	DueDate        *time.Time   `gorm:"type:timestamptz" json:"due_date"`
	CompletedAt    *time.Time   `gorm:"type:timestamptz" json:"completed_at"`
	EstimatedHours *float64     `gorm:"type:decimal(5,2)" json:"estimated_hours"`
	ActualHours    *float64     `gorm:"type:decimal(5,2)" json:"actual_hours"`
	IsArchived     bool         `gorm:"not null;default:false" json:"is_archived"`
	CreatedAt      time.Time    `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Associations
	Creator     User             `gorm:"foreignKey:CreatedBy" json:"-"`
	Assignee    *User            `gorm:"foreignKey:AssignedTo" json:"-"`
	Notes       []TaskNote       `gorm:"foreignKey:TaskID" json:"notes,omitempty"`
	Attachments []TaskAttachment `gorm:"foreignKey:TaskID" json:"attachments,omitempty"`
}

type TaskNote struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null" json:"task_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

type TaskAttachment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null" json:"task_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileSize  int64     `gorm:"type:bigint;not null" json:"file_size"`
	MimeType  string    `gorm:"type:varchar(100);not null" json:"mime_type"`
	FilePath  string    `gorm:"type:varchar(500);not null" json:"file_path"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
}
