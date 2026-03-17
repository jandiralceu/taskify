package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, taskID uuid.UUID, params UpdateTaskParams) (*models.Task, error)
	Delete(ctx context.Context, taskID uuid.UUID) error
	FindByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error)
	FindAll(ctx context.Context, filter TaskListFilter) ([]models.Task, error)

	// Notes
	CreateNote(ctx context.Context, params CreateNoteParams) (*models.TaskNote, error)
	UpdateNote(ctx context.Context, noteID uuid.UUID, params UpdateNoteParams) (*models.TaskNote, error)
	DeleteNote(ctx context.Context, noteID uuid.UUID) error
	FindNoteByID(ctx context.Context, noteID uuid.UUID) (*models.TaskNote, error)
	GetNotesByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TaskNote, error)

	// Attachments
	CreateAttachment(ctx context.Context, attachment *models.TaskAttachment) error
	DeleteAttachment(ctx context.Context, attachmentID uuid.UUID) error
	FindAttachmentByID(ctx context.Context, attachmentID uuid.UUID) (*models.TaskAttachment, error)
	GetAttachmentsByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TaskAttachment, error)
}

type taskRepository struct {
	db *gorm.DB
}

var _ TaskRepository = (*taskRepository)(nil)

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

type TaskListFilter struct {
	Status     models.TaskStatus
	Priority   models.TaskPriority
	CreatedBy  *uuid.UUID
	AssignedTo *uuid.UUID
	Search     string
	IsBlocked  *bool
	IsArchived *bool
	Sort       string
	Order      string
}

type UpdateTaskParams struct {
	Title          *string
	Description    *string
	Status         *models.TaskStatus
	Priority       *models.TaskPriority
	IsBlocked      *bool
	AssignedTo     *uuid.UUID
	DueDate        *time.Time
	EstimatedHours *float64
	ActualHours    *float64
	CompletedAt    *time.Time
	IsArchived     *bool
}

type CreateNoteParams struct {
	TaskID  uuid.UUID
	UserID  uuid.UUID
	Content string
}

type UpdateNoteParams struct {
	Content string
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

func (r *taskRepository) Update(ctx context.Context, taskID uuid.UUID, params UpdateTaskParams) (*models.Task, error) {
	updates := make(map[string]interface{})

	if params.Title != nil {
		updates["title"] = *params.Title
	}
	if params.Description != nil {
		updates["description"] = *params.Description
	}
	if params.Status != nil {
		updates["status"] = *params.Status
	}
	if params.Priority != nil {
		updates["priority"] = *params.Priority
	}
	if params.IsBlocked != nil {
		updates["is_blocked"] = *params.IsBlocked
	}
	if params.AssignedTo != nil {
		updates["assigned_to"] = params.AssignedTo
	}
	if params.DueDate != nil {
		updates["due_date"] = params.DueDate
	}
	if params.EstimatedHours != nil {
		updates["estimated_hours"] = params.EstimatedHours
	}
	if params.ActualHours != nil {
		updates["actual_hours"] = params.ActualHours
	}
	if params.CompletedAt != nil {
		updates["completed_at"] = params.CompletedAt
	}
	if params.IsArchived != nil {
		updates["is_archived"] = *params.IsArchived
	}

	var task models.Task
	result := r.db.WithContext(ctx).Model(&task).Where("id = ?", taskID).Updates(updates)
	if result.Error != nil {
		return nil, mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	if err := r.db.WithContext(ctx).Preload("Assignee").First(&task, "id = ?", taskID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}

	// Set counts
	var notesCount, attachmentsCount int64
	r.db.WithContext(ctx).Model(&models.TaskNote{}).Where("task_id = ?", taskID).Count(&notesCount)
	r.db.WithContext(ctx).Model(&models.TaskAttachment{}).Where("task_id = ?", taskID).Count(&attachmentsCount)
	task.NotesCount = int(notesCount)
	task.AttachmentsCount = int(attachmentsCount)

	return &task, nil
}

func (r *taskRepository) Delete(ctx context.Context, taskID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Task{}, "id = ?", taskID)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *taskRepository) FindByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error) {
	var task models.Task
	if err := r.db.WithContext(ctx).
		Preload("Notes.User").
		Preload("Attachments.User").
		Preload("Assignee").
		First(&task, "id = ?", taskID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}

	task.NotesCount = len(task.Notes)
	task.AttachmentsCount = len(task.Attachments)

	return &task, nil
}

func (r *taskRepository) FindAll(ctx context.Context, filter TaskListFilter) ([]models.Task, error) {
	var tasks []models.Task

	// Select all fields from tasks plus counts as virtual columns
	query := r.db.WithContext(ctx).Model(&models.Task{}).
		Select("tasks.*, " +
			"(SELECT COUNT(*) FROM task_notes WHERE task_notes.task_id = tasks.id) AS notes_count, " +
			"(SELECT COUNT(*) FROM task_attachments WHERE task_attachments.task_id = tasks.id) AS attachments_count")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.CreatedBy != nil {
		query = query.Where("created_by = ?", *filter.CreatedBy)
	}
	if filter.AssignedTo != nil {
		query = query.Where("assigned_to = ?", *filter.AssignedTo)
	}
	if filter.IsArchived != nil {
		query = query.Where("is_archived = ?", *filter.IsArchived)
	}
	if filter.IsBlocked != nil {
		query = query.Where("is_blocked = ?", *filter.IsBlocked)
	}
	if filter.Search != "" {
		search := "%" + sanitizeLike(filter.Search) + "%"
		query = query.Where("(title ILIKE ? OR description ILIKE ?)", search, search)
	}

	orderBy := "created_at DESC"
	if filter.Sort != "" {
		// Map frontend field names to database columns
		columnMap := map[string]string{
			"createdAt": "created_at",
			"title":     "title",
			"priority":  "priority",
			"dueDate":   "due_date",
		}

		dbColumn, ok := columnMap[filter.Sort]
		if !ok {
			// Fallback to the provided string if not in map, 
			// though it's better to be strict for security.
			dbColumn = filter.Sort 
		}

		order := filter.Order
		if order == "" {
			order = "desc"
		}
		orderBy = dbColumn + " " + order
	}

	if err := query.Order(orderBy).Preload("Assignee").Find(&tasks).Error; err != nil {
		return nil, mapDatabaseError(err)
	}

	return tasks, nil
}

func (r *taskRepository) CreateNote(ctx context.Context, params CreateNoteParams) (*models.TaskNote, error) {
	note := &models.TaskNote{
		TaskID:  params.TaskID,
		UserID:  params.UserID,
		Content: params.Content,
	}
	if err := r.db.WithContext(ctx).Create(note).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return note, nil
}

func (r *taskRepository) GetNotesByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TaskNote, error) {
	var notes []models.TaskNote
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at desc").Find(&notes).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return notes, nil
}

func (r *taskRepository) FindNoteByID(ctx context.Context, noteID uuid.UUID) (*models.TaskNote, error) {
	var note models.TaskNote
	if err := r.db.WithContext(ctx).First(&note, "id = ?", noteID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &note, nil
}

func (r *taskRepository) UpdateNote(ctx context.Context, noteID uuid.UUID, params UpdateNoteParams) (*models.TaskNote, error) {
	var note models.TaskNote
	result := r.db.WithContext(ctx).Model(&note).Where("id = ?", noteID).Update("content", params.Content)
	if result.Error != nil {
		return nil, mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	if err := r.db.WithContext(ctx).First(&note, "id = ?", noteID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}

	return &note, nil
}

func (r *taskRepository) DeleteNote(ctx context.Context, noteID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.TaskNote{}, "id = ?", noteID)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *taskRepository) CreateAttachment(ctx context.Context, attachment *models.TaskAttachment) error {
	if err := r.db.WithContext(ctx).Create(attachment).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

func (r *taskRepository) DeleteAttachment(ctx context.Context, attachmentID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.TaskAttachment{}, "id = ?", attachmentID)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *taskRepository) FindAttachmentByID(ctx context.Context, attachmentID uuid.UUID) (*models.TaskAttachment, error) {
	var attachment models.TaskAttachment
	if err := r.db.WithContext(ctx).First(&attachment, "id = ?", attachmentID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &attachment, nil
}

func (r *taskRepository) GetAttachmentsByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TaskAttachment, error) {
	var attachments []models.TaskAttachment
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at desc").Find(&attachments).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return attachments, nil
}
