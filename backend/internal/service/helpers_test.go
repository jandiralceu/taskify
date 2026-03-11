package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/jandiralceu/taskify/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

// MockCacheManager is a shared mock for pkg.CacheManager used across service tests.
type MockCacheManager struct {
	mock.Mock
}

func (m *MockCacheManager) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheManager) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCacheManager) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheManager) DeletePrefix(ctx context.Context, prefix string) error {
	args := m.Called(ctx, prefix)
	return args.Error(0)
}

func (m *MockCacheManager) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCacheManager) GetClient() *redis.Client {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.Client)
}

// MockUserRepository is a mock implementation of repository.UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, userID uuid.UUID, params repository.UpdateUserParams) (*models.User, error) {
	args := m.Called(ctx, userID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL *string) (*models.User, error) {
	args := m.Called(ctx, userID, avatarURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, filter repository.UserListFilter) ([]models.User, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) ChangePassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error {
	args := m.Called(ctx, userID, newHashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockPasswordHasher is a mock implementation of pkg.PasswordHasher.
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Verify(password, hash string) (bool, error) {
	args := m.Called(password, hash)
	return args.Bool(0), args.Error(1)
}

// MockTaskRepository is a mock implementation of repository.TaskRepository.
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(ctx context.Context, taskID uuid.UUID, params repository.UpdateTaskParams) (*models.Task, error) {
	args := m.Called(ctx, taskID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(ctx context.Context, taskID uuid.UUID) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, taskID uuid.UUID) (*models.Task, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAll(ctx context.Context, filter repository.TaskListFilter) ([]models.Task, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]models.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskRepository) CreateNote(ctx context.Context, params repository.CreateNoteParams) (*models.TaskNote, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskNote), args.Error(1)
}

func (m *MockTaskRepository) UpdateNote(ctx context.Context, noteID uuid.UUID, params repository.UpdateNoteParams) (*models.TaskNote, error) {
	args := m.Called(ctx, noteID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskNote), args.Error(1)
}

func (m *MockTaskRepository) DeleteNote(ctx context.Context, noteID uuid.UUID) error {
	args := m.Called(ctx, noteID)
	return args.Error(0)
}

func (m *MockTaskRepository) FindNoteByID(ctx context.Context, noteID uuid.UUID) (*models.TaskNote, error) {
	args := m.Called(ctx, noteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskNote), args.Error(1)
}

func (m *MockTaskRepository) GetNotesByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TaskNote, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]models.TaskNote), args.Error(1)
}

func (m *MockTaskRepository) CreateAttachment(ctx context.Context, attachment *models.TaskAttachment) error {
	args := m.Called(ctx, attachment)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteAttachment(ctx context.Context, attachmentID uuid.UUID) error {
	args := m.Called(ctx, attachmentID)
	return args.Error(0)
}

func (m *MockTaskRepository) FindAttachmentByID(ctx context.Context, attachmentID uuid.UUID) (*models.TaskAttachment, error) {
	args := m.Called(ctx, attachmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskAttachment), args.Error(1)
}

func (m *MockTaskRepository) GetAttachmentsByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TaskAttachment, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]models.TaskAttachment), args.Error(1)
}
