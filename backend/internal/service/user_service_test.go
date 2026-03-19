package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/dto"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByName(ctx context.Context, name string) (*models.RoleModel, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) GetPermissionsByRole(ctx context.Context, roleName string) ([]models.Permission, error) {
	args := m.Called(ctx, roleName)
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.RoleModel, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.RoleModel), args.Error(1)
}

func (m *MockRoleRepository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) ClearUserRoles(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	t.Run("DefaultRole", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRoleRepo := new(MockRoleRepository)
		mockHasher := new(MockPasswordHasher)
		svc := NewUserService(mockRepo, mockRoleRepo, mockHasher, "/tmp", nil)

		ctx := context.Background()
		user := &models.User{
			FirstName:    "Test",
			LastName:     "User",
			Email:        "test@example.com",
			PasswordHash: "plain-password",
		}

		mockHasher.On("Hash", "plain-password").Return("hashed-password", nil)
		mockRoleRepo.On("FindByName", ctx, "employee").Return(&models.RoleModel{Name: "employee"}, nil)
		mockRepo.On("Create", ctx, user).Return(nil)

		err := svc.Create(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, "hashed-password", user.PasswordHash)
		assert.Len(t, user.Roles, 1)
		assert.Equal(t, "employee", user.Roles[0].Name)
		mockHasher.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("SpecificRole", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRoleRepo := new(MockRoleRepository)
		mockHasher := new(MockPasswordHasher)
		svc := NewUserService(mockRepo, mockRoleRepo, mockHasher, "/tmp", nil)

		ctx := context.Background()
		user := &models.User{
			FirstName:    "Admin",
			LastName:     "User",
			Email:        "admin@example.com",
			PasswordHash: "secret",
			Role:         "admin",
		}

		mockHasher.On("Hash", "secret").Return("hashed-secret", nil)
		mockRoleRepo.On("FindByName", ctx, "admin").Return(&models.RoleModel{Name: "admin"}, nil)
		mockRepo.On("Create", ctx, user).Return(nil)

		err := svc.Create(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, "admin", user.Roles[0].Name)
	})
}

func TestFindAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := NewUserService(mockRepo, new(MockRoleRepository), mockHasher, "/tmp", nil)

	ctx := context.Background()
	req := dto.GetUserListRequest{
		FirstName: "test",
	}

	users := []models.User{{FirstName: "test", LastName: "user"}}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("repository.UserListFilter")).
		Return(users, int64(1), nil)

	res, err := svc.FindAll(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Len(t, res.Data, 1)
	mockRepo.AssertExpectations(t)
}

func TestFindUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := NewUserService(mockRepo, new(MockRoleRepository), mockHasher, "/tmp", nil)

	ctx := context.Background()
	userID := uuid.New()
	user := &models.User{ID: userID, FirstName: "test", LastName: "user"}

	mockRepo.On("FindByID", ctx, userID).Return(user, nil)

	res, err := svc.FindByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	mockCache := new(MockCacheManager)
	svc := NewUserService(mockRepo, new(MockRoleRepository), mockHasher, "/tmp", mockCache)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("Delete", ctx, userID).Return(nil)
	mockCache.On("DeletePrefix", ctx, mock.AnythingOfType("string")).Return(nil)

	err := svc.Delete(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestChangePassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockHasher := new(MockPasswordHasher)
		svc := NewUserService(mockRepo, new(MockRoleRepository), mockHasher, "/tmp", nil)

		ctx := context.Background()
		userID := uuid.New()
		req := dto.ChangePasswordRequest{
			OldPassword: "old",
			NewPassword: "new",
		}

		user := &models.User{ID: userID, PasswordHash: "old-hashed"}

		mockRepo.On("FindByID", ctx, userID).Return(user, nil)
		mockHasher.On("Verify", "old", "old-hashed").Return(true, nil)
		mockHasher.On("Hash", "new").Return("new-hashed", nil)
		mockRepo.On("ChangePassword", ctx, userID, "new-hashed").Return(nil)

		err := svc.ChangePassword(ctx, userID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("IncorrectOldPassword", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockHasher := new(MockPasswordHasher)
		svc := NewUserService(mockRepo, new(MockRoleRepository), mockHasher, "/tmp", nil)

		ctx := context.Background()
		userID := uuid.New()
		req := dto.ChangePasswordRequest{
			OldPassword: "wrong",
			NewPassword: "new",
		}

		user := &models.User{ID: userID, PasswordHash: "old-hashed"}

		mockRepo.On("FindByID", ctx, userID).Return(user, nil)
		mockHasher.On("Verify", "wrong", "old-hashed").Return(false, nil)

		err := svc.ChangePassword(ctx, userID, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperrors.ErrUnauthorized))
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("BasicProfile", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRoleRepo := new(MockRoleRepository)
		mockHasher := new(MockPasswordHasher)
		svc := NewUserService(mockRepo, mockRoleRepo, mockHasher, "/tmp", nil)

		ctx := context.Background()
		userID := uuid.New()
		newName := "Updated"
		req := dto.UpdateUserRequest{
			FirstName: &newName,
		}

		expectedUser := &models.User{ID: userID, FirstName: newName}
		mockRepo.On("Update", ctx, userID, mock.AnythingOfType("repository.UpdateUserParams")).
			Return(expectedUser, nil)

		res, err := svc.Update(ctx, userID, req)

		assert.NoError(t, err)
		assert.Equal(t, newName, res.FirstName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("WithRoleChange", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRoleRepo := new(MockRoleRepository)
		svc := NewUserService(mockRepo, mockRoleRepo, nil, "/tmp", nil)

		ctx := context.Background()
		userID := uuid.New()
		newRole := "admin"
		req := dto.UpdateUserRequest{
			Role: &newRole,
		}

		roleID := uuid.New()
		mockRoleRepo.On("FindByName", ctx, "admin").Return(&models.RoleModel{ID: roleID, Name: "admin"}, nil)
		mockRoleRepo.On("ClearUserRoles", ctx, userID).Return(nil)
		mockRoleRepo.On("AssignRoleToUser", ctx, userID, roleID).Return(nil)

		expectedUser := &models.User{ID: userID, Roles: []models.RoleModel{}}
		mockRepo.On("Update", ctx, userID, mock.AnythingOfType("repository.UpdateUserParams")).
			Return(expectedUser, nil)

		res, err := svc.Update(ctx, userID, req)

		assert.NoError(t, err)
		assert.Equal(t, "admin", res.Role)
		assert.Len(t, res.Roles, 1)
		mockRoleRepo.AssertExpectations(t)
	})
}
