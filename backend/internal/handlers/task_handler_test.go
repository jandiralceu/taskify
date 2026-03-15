package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/dto"
	"github.com/jandiralceu/taskify/internal/middleware"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock TaskService ---

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTaskRequest) (*models.Task, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) Update(ctx context.Context, taskID uuid.UUID, req dto.UpdateTaskRequest) (*models.Task, error) {
	args := m.Called(ctx, taskID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) Delete(ctx context.Context, taskID uuid.UUID) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *MockTaskService) GetByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) GetAll(ctx context.Context, req dto.GetTaskListRequest) ([]models.Task, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskService) AddNote(ctx context.Context, taskID, userID uuid.UUID, req dto.CreateTaskNoteRequest) (*models.TaskNote, error) {
	args := m.Called(ctx, taskID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskNote), args.Error(1)
}

func (m *MockTaskService) UpdateNote(ctx context.Context, noteID, userID uuid.UUID, req dto.UpdateTaskNoteRequest) (*models.TaskNote, error) {
	args := m.Called(ctx, noteID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskNote), args.Error(1)
}

func (m *MockTaskService) DeleteNote(ctx context.Context, noteID, userID uuid.UUID) error {
	args := m.Called(ctx, noteID, userID)
	return args.Error(0)
}

func (m *MockTaskService) GetNotes(ctx context.Context, taskID uuid.UUID) ([]models.TaskNote, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]models.TaskNote), args.Error(1)
}

func (m *MockTaskService) AddAttachment(ctx context.Context, taskID, userID uuid.UUID, file io.Reader, filename string, size int64, mimeType string) (*models.TaskAttachment, error) {
	args := m.Called(ctx, taskID, userID, file, filename, size, mimeType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskAttachment), args.Error(1)
}

func (m *MockTaskService) DeleteAttachment(ctx context.Context, attachmentID, userID uuid.UUID) error {
	args := m.Called(ctx, attachmentID, userID)
	return args.Error(0)
}

func (m *MockTaskService) GetAttachments(ctx context.Context, taskID uuid.UUID) ([]models.TaskAttachment, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]models.TaskAttachment), args.Error(1)
}

// =====================
// CreateTask Tests
// =====================

func TestCreateTask_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	userID := uuid.New()
	req := dto.CreateTaskRequest{
		Title: "New Task",
	}
	resTask := &models.Task{ID: uuid.New(), Title: req.Title, CreatedBy: userID}

	mockService.On("Create", mock.Anything, userID, req).Return(resTask, nil)

	router := setupRouter()
	router.POST("/tasks", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.CreateTask(c)
	})

	w := performRequest(router, "POST", "/tasks", req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var actualTask models.Task
	json.Unmarshal(w.Body.Bytes(), &actualTask)
	assert.Equal(t, resTask.ID, actualTask.ID)
	mockService.AssertExpectations(t)
}

func TestCreateTask_BadRequest(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	router := setupRouter()
	router.POST("/tasks", handler.CreateTask)

	// Missing required title
	req := map[string]string{}
	w := performRequest(router, "POST", "/tasks", req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =====================
// UpdateTask Tests
// =====================

func TestUpdateTask_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	taskID := uuid.New()
	newTitle := "Updated Title"
	req := dto.UpdateTaskRequest{
		Title: &newTitle,
	}
	resTask := &models.Task{ID: taskID, Title: newTitle}

	mockService.On("Update", mock.Anything, taskID, req).Return(resTask, nil)

	router := setupRouter()
	router.PATCH("/tasks/:id", handler.UpdateTask)

	w := performRequest(router, "PATCH", fmt.Sprintf("/tasks/%s", taskID), req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateTask_InvalidID(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	router := setupRouter()
	router.PATCH("/tasks/:id", handler.UpdateTask)

	w := performRequest(router, "PATCH", "/tasks/invalid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =====================
// GetTask Tests
// =====================

func TestGetTask_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	taskID := uuid.New()
	resTask := &models.Task{ID: taskID, Title: "Test Task"}

	mockService.On("GetByID", mock.Anything, taskID).Return(resTask, nil)

	router := setupRouter()
	router.GET("/tasks/:id", handler.GetTask)

	w := performRequest(router, "GET", fmt.Sprintf("/tasks/%s", taskID), nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetTask_NotFound(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	taskID := uuid.New()
	mockService.On("GetByID", mock.Anything, taskID).Return(nil, apperrors.ErrNotFound)

	router := setupRouter()
	router.GET("/tasks/:id", handler.GetTask)

	w := performRequest(router, "GET", fmt.Sprintf("/tasks/%s", taskID), nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// =====================
// DeleteTask Tests
// =====================

func TestDeleteTask_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	taskID := uuid.New()
	mockService.On("Delete", mock.Anything, taskID).Return(nil)

	router := setupRouter()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	w := performRequest(router, "DELETE", fmt.Sprintf("/tasks/%s", taskID), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

// =====================
// ListTasks Tests
// =====================

func TestListTasks_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	res := []models.Task{{ID: uuid.New(), Title: "Task 1"}}

	mockService.On("GetAll", mock.Anything, mock.AnythingOfType("dto.GetTaskListRequest")).Return(res, nil)

	router := setupRouter()
	router.GET("/tasks", handler.ListTasks)

	w := performRequest(router, "GET", "/tasks", nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

// =====================
// AddNote Tests
// =====================

func TestAddNote_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	taskID := uuid.New()
	userID := uuid.New()
	req := dto.CreateTaskNoteRequest{Content: "Test Note"}
	resNote := &models.TaskNote{ID: uuid.New(), Content: req.Content}

	mockService.On("AddNote", mock.Anything, taskID, userID, req).Return(resNote, nil)

	router := setupRouter()
	router.POST("/tasks/:id/notes", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.AddNote(c)
	})

	w := performRequest(router, "POST", fmt.Sprintf("/tasks/%s/notes", taskID), req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

// =====================
// GetNotes Tests
// =====================

func TestGetNotes_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	taskID := uuid.New()
	resNotes := []models.TaskNote{{ID: uuid.New(), Content: "Note 1"}}

	mockService.On("GetNotes", mock.Anything, taskID).Return(resNotes, nil)

	router := setupRouter()
	router.GET("/tasks/:id/notes", handler.GetNotes)

	w := performRequest(router, "GET", fmt.Sprintf("/tasks/%s/notes", taskID), nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

// =====================
// UpdateNote Tests
// =====================

func TestUpdateNote_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	noteID := uuid.New()
	userID := uuid.New()
	req := dto.UpdateTaskNoteRequest{Content: "Updated Note"}
	resNote := &models.TaskNote{ID: noteID, Content: req.Content}

	mockService.On("UpdateNote", mock.Anything, noteID, userID, req).Return(resNote, nil)

	router := setupRouter()
	router.PATCH("/tasks/notes/:noteId", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.UpdateNote(c)
	})

	w := performRequest(router, "PATCH", fmt.Sprintf("/tasks/notes/%s", noteID), req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// =====================
// DeleteNote Tests
// =====================

func TestDeleteNote_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	noteID := uuid.New()
	userID := uuid.New()

	mockService.On("DeleteNote", mock.Anything, noteID, userID).Return(nil)

	router := setupRouter()
	router.DELETE("/tasks/notes/:noteId", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.DeleteNote(c)
	})

	w := performRequest(router, "DELETE", fmt.Sprintf("/tasks/notes/%s", noteID), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// =====================
// AddAttachment Tests
// =====================

func TestAddAttachment_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	taskID := uuid.New()
	userID := uuid.New()
	resAttachment := &models.TaskAttachment{ID: uuid.New(), FileName: "test.txt"}

	mockService.On("AddAttachment", mock.Anything, taskID, userID, mock.Anything, "test.txt", mock.Anything, mock.Anything).
		Return(resAttachment, nil)

	router := setupRouter()
	router.POST("/tasks/:id/attachments", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.AddAttachment(c)
	})

	body, contentType := createMultipartForm("file", "test.txt", "content")
	w := performRequestWithContentType(router, "POST", fmt.Sprintf("/tasks/%s/attachments", taskID), body, contentType)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// =====================
// GetAttachments Tests
// =====================

func TestGetAttachments_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	taskID := uuid.New()
	resAttachments := []models.TaskAttachment{{ID: uuid.New(), FileName: "file.png"}}

	mockService.On("GetAttachments", mock.Anything, taskID).Return(resAttachments, nil)

	router := setupRouter()
	router.GET("/tasks/:id/attachments", handler.GetAttachments)

	w := performRequest(router, "GET", fmt.Sprintf("/tasks/%s/attachments", taskID), nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

// =====================
// DeleteAttachment Tests
// =====================

func TestDeleteAttachment_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	attachmentID := uuid.New()
	userID := uuid.New()

	mockService.On("DeleteAttachment", mock.Anything, attachmentID, userID).Return(nil)

	router := setupRouter()
	router.DELETE("/tasks/attachments/:attachmentId", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.DeleteAttachment(c)
	})

	w := performRequest(router, "DELETE", fmt.Sprintf("/tasks/attachments/%s", attachmentID), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
