package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/dto"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/jandiralceu/taskify/internal/repository"
)

type TaskService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateTaskRequest) (*models.Task, error)
	Update(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskRequest) (*models.Task, error)
	Delete(ctx context.Context, taskID uuid.UUID) error
	GetByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error)
	GetAll(ctx context.Context, req dto.GetTaskListRequest) ([]models.Task, error)

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
		IsBlocked:      req.IsBlocked,
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

	params := repository.UpdateTaskParams{
		Title:          req.Title,
		Description:    req.Description,
		Status:         req.Status,
		Priority:       req.Priority,
		IsBlocked:      req.IsBlocked,
		AssignedTo:     req.AssignedTo,
		DueDate:        req.DueDate,
		EstimatedHours: req.EstimatedHours,
		ActualHours:    req.ActualHours,
		IsArchived:     req.IsArchived,
	}

	// Business logic: automatically set completed_at when status changes to completed
	if req.Status != nil && *req.Status == models.TaskStatusCompleted && task.CompletedAt == nil {
		now := time.Now()
		params.CompletedAt = &now
	}

	return s.taskRepo.Update(ctx, taskID, params)
}

func (s *taskService) Delete(ctx context.Context, taskID uuid.UUID) error {
	return s.taskRepo.Delete(ctx, taskID)
}

func (s *taskService) GetByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error) {
	return s.taskRepo.FindByID(ctx, taskID)
}

func (s *taskService) GetAll(ctx context.Context, req dto.GetTaskListRequest) ([]models.Task, error) {
	filter := repository.TaskListFilter{
		Status:     req.Status,
		Priority:   req.Priority,
		CreatedBy:  req.CreatedBy,
		AssignedTo: req.AssignedTo,
		Search:     req.Search,
		IsBlocked:  req.IsBlocked,
		IsArchived: req.IsArchived,
		Sort:       req.Sort,
		Order:      req.Order,
	}

	return s.taskRepo.FindAll(ctx, filter)
}

func (s *taskService) AddNote(ctx context.Context, taskID, userID uuid.UUID, req dto.CreateTaskNoteRequest) (*models.TaskNote, error) {
	// Verify task exists
	_, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	params := repository.CreateNoteParams{
		TaskID:  taskID,
		UserID:  userID,
		Content: req.Content,
	}

	return s.taskRepo.CreateNote(ctx, params)
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

	params := repository.UpdateNoteParams{
		Content: req.Content,
	}

	return s.taskRepo.UpdateNote(ctx, noteID, params)
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
		return nil, apperrors.ErrStorage
	}

	// 4. Create local file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, apperrors.ErrStorage
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, apperrors.ErrStorage
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
