package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/dto"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/jandiralceu/taskify/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskServiceCreate(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	userID := uuid.New()

	t.Run("SuccessWithDefaults", func(t *testing.T) {
		req := dto.CreateTaskRequest{
			Title: "Test Task",
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(task *models.Task) bool {
			return task.Title == req.Title &&
				task.Status == models.TaskStatusPending &&
				task.Priority == models.TaskPriorityMedium &&
				task.CreatedBy == userID
		})).Return(nil).Once()

		task, err := svc.Create(ctx, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, models.TaskStatusPending, task.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepoError", func(t *testing.T) {
		req := dto.CreateTaskRequest{Title: "Error Task"}
		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error")).Once()

		task, err := svc.Create(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskServiceUpdate(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	taskID := uuid.New()

	t.Run("SuccessAndSetCompletedAt", func(t *testing.T) {
		existingTask := &models.Task{ID: taskID, Status: models.TaskStatusPending}
		mockRepo.On("FindByID", ctx, taskID).Return(existingTask, nil).Once()

		completedStatus := models.TaskStatusCompleted
		req := dto.UpdateTaskRequest{Status: &completedStatus}

		mockRepo.On("Update", ctx, taskID, mock.MatchedBy(func(params repository.UpdateTaskParams) bool {
			return params.Status != nil && *params.Status == models.TaskStatusCompleted && params.CompletedAt != nil
		})).Return(&models.Task{Status: models.TaskStatusCompleted}, nil).Once()

		task, err := svc.Update(ctx, taskID, req)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, taskID).Return(nil, apperrors.ErrNotFound).Once()
		req := dto.UpdateTaskRequest{}

		task, err := svc.Update(ctx, taskID, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperrors.ErrNotFound))
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskServiceGetByID(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	taskID := uuid.New()

	mockRepo.On("FindByID", ctx, taskID).Return(&models.Task{ID: taskID}, nil).Once()
	task, err := svc.GetByID(ctx, taskID)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, taskID, task.ID)
	mockRepo.AssertExpectations(t)
}

func TestTaskServiceDelete(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	taskID := uuid.New()

	mockRepo.On("Delete", ctx, taskID).Return(nil).Once()
	err := svc.Delete(ctx, taskID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTaskServiceGetAll(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	req := dto.GetTaskListRequest{}

	tasks := []models.Task{{Title: "Task 1"}}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("repository.TaskListFilter")).
		Return(tasks, nil).Once()

	res, err := svc.GetAll(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Task 1", res[0].Title)
	mockRepo.AssertExpectations(t)
}

func TestTaskServiceAddNote(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	taskID := uuid.New()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, taskID).Return(&models.Task{ID: taskID}, nil).Once()
		req := dto.CreateTaskNoteRequest{Content: "Test Note"}
		
		mockRepo.On("CreateNote", ctx, repository.CreateNoteParams{
			TaskID:  taskID,
			UserID:  userID,
			Content: req.Content,
		}).Return(&models.TaskNote{Content: req.Content}, nil).Once()

		note, err := svc.AddNote(ctx, taskID, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, note)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskServiceGetNotes(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	taskID := uuid.New()

	mockRepo.On("GetNotesByTaskID", ctx, taskID).Return([]models.TaskNote{{Content: "note"}}, nil).Once()
	notes, err := svc.GetNotes(ctx, taskID)

	assert.NoError(t, err)
	assert.Len(t, notes, 1)
	mockRepo.AssertExpectations(t)
}

func TestTaskServiceUpdateNote(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	noteID := uuid.New()
	userID := uuid.New()

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.On("FindNoteByID", ctx, noteID).Return(&models.TaskNote{UserID: uuid.New()}, nil).Once()
		req := dto.UpdateTaskNoteRequest{Content: "New Content"}

		note, err := svc.UpdateNote(ctx, noteID, userID, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperrors.ErrForbidden))
		assert.Nil(t, note)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskServiceDeleteNote(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	noteID := uuid.New()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindNoteByID", ctx, noteID).Return(&models.TaskNote{UserID: userID}, nil).Once()
		mockRepo.On("DeleteNote", ctx, noteID).Return(nil).Once()

		err := svc.DeleteNote(ctx, noteID, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.On("FindNoteByID", ctx, noteID).Return(&models.TaskNote{UserID: uuid.New()}, nil).Once()

		err := svc.DeleteNote(ctx, noteID, userID)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskServiceAddAttachment(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	tempDir, _ := os.MkdirTemp("", "taskify-uploads")
	defer func() { _ = os.RemoveAll(tempDir) }()
	
	svc := NewTaskService(mockRepo, tempDir)
	ctx := context.Background()
	taskID := uuid.New()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, taskID).Return(&models.Task{ID: taskID}, nil).Once()
		
		content := "file content"
		reader := strings.NewReader(content)
		filename := "test.txt"

		mockRepo.On("CreateAttachment", ctx, mock.AnythingOfType("*models.TaskAttachment")).Return(nil).Once()

		attachment, err := svc.AddAttachment(ctx, taskID, userID, reader, filename, int64(len(content)), "text/plain")

		assert.NoError(t, err)
		assert.NotNil(t, attachment)
		assert.True(t, strings.HasSuffix(attachment.FileName, filename))
		
		// Check file exists on disk
		diskPath := filepath.Join(tempDir, "attachments", filepath.Base(attachment.FilePath))
		_, err = os.Stat(diskPath)
		assert.NoError(t, err)
		
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskServiceDeleteAttachment(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	tempDir, _ := os.MkdirTemp("", "taskify-uploads")
	defer func() { _ = os.RemoveAll(tempDir) }()

	svc := NewTaskService(mockRepo, tempDir)
	ctx := context.Background()
	attachmentID := uuid.New()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Create a dummy file in the 'attachments' subfolder
		attachmentDir := filepath.Join(tempDir, "attachments")
		_ = os.MkdirAll(attachmentDir, 0755)
		
		fileName := "to_delete.txt"
		diskPath := filepath.Join(attachmentDir, fileName)
		_ = os.WriteFile(diskPath, []byte("test"), 0644)

		attachment := &models.TaskAttachment{
			ID:       attachmentID,
			UserID:   userID,
			FilePath: "/uploads/attachments/" + fileName,
		}

		mockRepo.On("FindAttachmentByID", ctx, attachmentID).Return(attachment, nil).Once()
		mockRepo.On("DeleteAttachment", ctx, attachmentID).Return(nil).Once()

		err := svc.DeleteAttachment(ctx, attachmentID, userID)

		assert.NoError(t, err)
		// Check file is gone
		_, err = os.Stat(diskPath)
		assert.True(t, os.IsNotExist(err))
		
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskServiceGetAttachments(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	svc := NewTaskService(mockRepo, "/tmp")
	ctx := context.Background()
	taskID := uuid.New()

	mockRepo.On("GetAttachmentsByTaskID", ctx, taskID).Return([]models.TaskAttachment{{FileName: "file"}}, nil).Once()
	attachments, err := svc.GetAttachments(ctx, taskID)

	assert.NoError(t, err)
	assert.Len(t, attachments, 1)
	mockRepo.AssertExpectations(t)
}
