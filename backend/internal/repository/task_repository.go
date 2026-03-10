package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, taskID uuid.UUID) error
	FindByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error)
	FindAll(ctx context.Context, filter TaskListFilter) (tasks []models.Task, total int64, err error)

	// Notes
	CreateNote(ctx context.Context, note *models.TaskNote) error
	UpdateNote(ctx context.Context, note *models.TaskNote) error
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
	IsArchived *bool
	Pagination PaginationParams
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	if err := r.db.WithContext(ctx).Save(task).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
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
		Preload("Notes").
		Preload("Attachments").
		First(&task, "id = ?", taskID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &task, nil
}

func (r *taskRepository) FindAll(ctx context.Context, filter TaskListFilter) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Task{})

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
	if filter.Search != "" {
		search := "%" + sanitizeLike(filter.Search) + "%"
		query = query.Where("(title ILIKE ? OR description ILIKE ?)", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	err := query.
		Order(filter.Pagination.GetOrderBy()).
		Offset(filter.Pagination.GetOffset()).
		Limit(filter.Pagination.Limit).
		Find(&tasks).Error

	if err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return tasks, total, nil
}

func (r *taskRepository) CreateNote(ctx context.Context, note *models.TaskNote) error {
	if err := r.db.WithContext(ctx).Create(note).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
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

func (r *taskRepository) UpdateNote(ctx context.Context, note *models.TaskNote) error {
	if err := r.db.WithContext(ctx).Save(note).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
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
