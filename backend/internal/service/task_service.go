package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
)

type TaskService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateTaskRequest) (*models.Task, error)
	Update(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskRequest) (*models.Task, error)
	Delete(ctx context.Context, taskID uuid.UUID) error
	GetByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error)
	GetAll(ctx context.Context, req dto.GetTaskListRequest) (dto.PaginatedResponse[models.Task], error)

	// Notes
	AddNote(ctx context.Context, taskID, userID uuid.UUID, req dto.CreateTaskNoteRequest) (*models.TaskNote, error)
	UpdateNote(ctx context.Context, noteID, userID uuid.UUID, req dto.UpdateTaskNoteRequest) (*models.TaskNote, error)
	DeleteNote(ctx context.Context, noteID, userID uuid.UUID) error
	GetNotes(ctx context.Context, taskID uuid.UUID) ([]models.TaskNote, error)

	// Attachments
	AddAttachment(ctx context.Context, taskID, userID uuid.UUID, file io.Reader, filename string, size int64, mimeType string) (*models.TaskAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID, userID uuid.UUID) error
	GetAttachments(ctx context.Context, taskID uuid.UUID) ([]models.TaskAttachment, error)
}

type taskService struct {
	taskRepo   repository.TaskRepository
	uploadPath string
}

var _ TaskService = (*taskService)(nil)

func NewTaskService(taskRepo repository.TaskRepository, uploadPath string) TaskService {
	return &taskService{
		taskRepo:   taskRepo,
		uploadPath: uploadPath,
	}
}

func (s *taskService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTaskRequest) (*models.Task, error) {
	status := req.Status
	if status == "" {
		status = models.TaskStatusPending
	}

	priority := req.Priority
	if priority == "" {
		priority = models.TaskPriorityMedium
	}

	task := &models.Task{
		Title:          req.Title,
		Description:    req.Description,
		Status:         status,
		Priority:       priority,
		CreatedBy:      userID,
		AssignedTo:     req.AssignedTo,
		DueDate:        req.DueDate,
		EstimatedHours: req.EstimatedHours,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) Update(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskRequest) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
		if *req.Status == models.TaskStatusCompleted && task.CompletedAt == nil {
			now := time.Now()
			task.CompletedAt = &now
		}
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.AssignedTo != nil {
		task.AssignedTo = req.AssignedTo
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.EstimatedHours != nil {
		task.EstimatedHours = req.EstimatedHours
	}
	if req.ActualHours != nil {
		task.ActualHours = req.ActualHours
	}
	if req.IsArchived != nil {
		task.IsArchived = *req.IsArchived
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) Delete(ctx context.Context, taskID uuid.UUID) error {
	return s.taskRepo.Delete(ctx, taskID)
}

func (s *taskService) GetByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error) {
	return s.taskRepo.FindByID(ctx, taskID)
}

func (s *taskService) GetAll(ctx context.Context, req dto.GetTaskListRequest) (dto.PaginatedResponse[models.Task], error) {
	filter := repository.TaskListFilter{
		Status:     req.Status,
		Priority:   req.Priority,
		CreatedBy:  req.CreatedBy,
		AssignedTo: req.AssignedTo,
		Search:     req.Search,
		Pagination: repository.PaginationParams{
			Page:  req.GetPage(),
			Limit: req.GetLimit(),
			Sort:  req.GetSort("created_at", "title", "due_date", "priority", "status"),
			Order: req.GetOrder(),
		},
	}

	tasks, total, err := s.taskRepo.FindAll(ctx, filter)
	if err != nil {
		return dto.PaginatedResponse[models.Task]{}, err
	}

	return dto.NewPaginatedResponse(tasks, total, filter.Pagination.Page, filter.Pagination.Limit), nil
}

func (s *taskService) AddNote(ctx context.Context, taskID, userID uuid.UUID, req dto.CreateTaskNoteRequest) (*models.TaskNote, error) {
	// Verify task exists
	_, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	note := &models.TaskNote{
		TaskID:  taskID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.taskRepo.CreateNote(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *taskService) GetNotes(ctx context.Context, taskID uuid.UUID) ([]models.TaskNote, error) {
	return s.taskRepo.GetNotesByTaskID(ctx, taskID)
}

func (s *taskService) UpdateNote(ctx context.Context, noteID, userID uuid.UUID, req dto.UpdateTaskNoteRequest) (*models.TaskNote, error) {
	note, err := s.taskRepo.FindNoteByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	// Ownership check
	if note.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	note.Content = req.Content

	if err := s.taskRepo.UpdateNote(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *taskService) DeleteNote(ctx context.Context, noteID, userID uuid.UUID) error {
	note, err := s.taskRepo.FindNoteByID(ctx, noteID)
	if err != nil {
		return err
	}

	// Ownership check
	if note.UserID != userID {
		return apperrors.ErrForbidden
	}

	return s.taskRepo.DeleteNote(ctx, noteID)
}

func (s *taskService) AddAttachment(ctx context.Context, taskID, userID uuid.UUID, file io.Reader, filename string, size int64, mimeType string) (*models.TaskAttachment, error) {
	// 1. Verify task exists
	_, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// 2. Generate unique filename to avoid collisions
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(s.uploadPath, uniqueName)

	// 3. Ensure upload directory exists
	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 4. Create local file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to stream file to disk: %w", err)
	}

	// 5. Build attachment record (save relative or absolute path, usually relative is better for portability)
	// But for now we use what's configured
	attachment := &models.TaskAttachment{
		TaskID:   taskID,
		UserID:   userID,
		FileName: filename,
		FileSize: size,
		MimeType: mimeType,
		FilePath: filePath,
	}

	if err := s.taskRepo.CreateAttachment(ctx, attachment); err != nil {
		// Rollback: delete file if DB fails
		os.Remove(filePath)
		return nil, err
	}

	return attachment, nil
}

func (s *taskService) DeleteAttachment(ctx context.Context, attachmentID, userID uuid.UUID) error {
	attachment, err := s.taskRepo.FindAttachmentByID(ctx, attachmentID)
	if err != nil {
		return err
	}

	// Permission check
	if attachment.UserID != userID {
		return apperrors.ErrForbidden
	}

	// 1. Delete from DB
	if err := s.taskRepo.DeleteAttachment(ctx, attachmentID); err != nil {
		return err
	}

	// 2. Delete from disk (best effort, don't fail if file is already gone)
	_ = os.Remove(attachment.FilePath)

	return nil
}

func (s *taskService) GetAttachments(ctx context.Context, taskID uuid.UUID) ([]models.TaskAttachment, error) {
	return s.taskRepo.GetAttachmentsByTaskID(ctx, taskID)
}
