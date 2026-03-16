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
	IsBlocked      bool         `gorm:"not null;default:false" json:"isBlocked"`
	CreatedBy      uuid.UUID    `gorm:"type:uuid;not null" json:"createdBy"`
	AssignedTo     *uuid.UUID   `gorm:"type:uuid" json:"assignedTo"`
	DueDate        *time.Time   `gorm:"type:timestamptz" json:"dueDate"`
	CompletedAt    *time.Time   `gorm:"type:timestamptz" json:"completedAt"`
	EstimatedHours *float64     `gorm:"type:decimal(5,2)" json:"estimatedHours"`
	ActualHours    *float64     `gorm:"type:decimal(5,2)" json:"actualHours"`
	IsArchived     bool         `gorm:"not null;default:false" json:"isArchived"`
	CreatedAt      time.Time    `gorm:"type:timestamptz;not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time    `gorm:"type:timestamptz;not null;default:now()" json:"updatedAt"`

	// UI Fields (populated via subqueries)
	NotesCount       int `gorm:"-" json:"notesCount"`
	AttachmentsCount int `gorm:"-" json:"attachmentsCount"`

	// Associations
	Creator     User             `gorm:"foreignKey:CreatedBy" json:"-"`
	Assignee    *User            `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	Notes       []TaskNote       `gorm:"foreignKey:TaskID" json:"notes,omitempty"`
	Attachments []TaskAttachment `gorm:"foreignKey:TaskID" json:"attachments,omitempty"`
}

type TaskNote struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null" json:"taskId"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updatedAt"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

type TaskAttachment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null" json:"taskId"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"fileName"`
	FileSize  int64     `gorm:"type:bigint;not null" json:"fileSize"`
	MimeType  string    `gorm:"type:varchar(100);not null" json:"mimeType"`
	FilePath  string    `gorm:"type:text;not null" json:"filePath"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updatedAt"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
}
