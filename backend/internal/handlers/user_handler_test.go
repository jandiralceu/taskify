package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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

// =====================
// FindAllUsers Tests
// =====================

func TestFindAllUsersSuccess(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	users := []models.User{
		{ID: uuid.New(), FirstName: "User", LastName: "1", Email: "user1@example.com"},
		{ID: uuid.New(), FirstName: "User", LastName: "2", Email: "user2@example.com"},
	}

	response := dto.PaginatedResponse[models.User]{
		Data:       users,
		Total:      2,
		TotalPages: 1,
		Page:       1,
		Limit:      10,
	}

	mockService.On("FindAll", mock.Anything, mock.AnythingOfType("dto.GetUserListRequest")).
		Return(&response, nil)

	router := setupRouter()
	router.GET("/users", handler.FindAllUsers)

	w := performRequest(router, "GET", "/users?page=1&limit=10", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse dto.PaginatedResponse[models.User]
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), actualResponse.Total)
	assert.Len(t, actualResponse.Data, 2)
	assert.Equal(t, "User", actualResponse.Data[0].FirstName)
	mockService.AssertExpectations(t)
}

func TestFindAllUsersServiceError(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	mockService.On("FindAll", mock.Anything, mock.AnythingOfType("dto.GetUserListRequest")).
		Return(nil, errors.New("database error"))

	router := setupRouter()
	router.GET("/users", handler.FindAllUsers)

	w := performRequest(router, "GET", "/users", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// =====================
// FindUserByID Tests
// =====================

func TestFindUserByIDSuccess(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	userID := uuid.New()
	user := &models.User{ID: userID, FirstName: "John", LastName: "Doe", Email: "john@example.com"}

	mockService.On("FindByID", mock.Anything, userID).Return(user, nil)

	router := setupRouter()
	router.GET("/users/:id", handler.FindUserByID)

	w := performRequest(router, "GET", fmt.Sprintf("/users/%s", userID), nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.NoError(t, err)
	assert.Equal(t, userID, actualUser.ID)
	assert.Equal(t, "John", actualUser.FirstName)
	mockService.AssertExpectations(t)
}

func TestFindUserByIDInvalidID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	router := setupRouter()
	router.GET("/users/:id", handler.FindUserByID)

	w := performRequest(router, "GET", "/users/invalid-uuid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "FindByID")
}

func TestFindUserByIDNotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	userID := uuid.New()
	mockService.On("FindByID", mock.Anything, userID).Return(nil, apperrors.ErrNotFound)

	router := setupRouter()
	router.GET("/users/:id", handler.FindUserByID)

	w := performRequest(router, "GET", fmt.Sprintf("/users/%s", userID), nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// =====================
// DeleteUser Tests
// =====================

func TestDeleteUserSuccess(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	userID := uuid.New()
	mockService.On("Delete", mock.Anything, userID).Return(nil)

	router := setupRouter()
	router.DELETE("/users/:id", handler.DeleteUser)

	w := performRequest(router, "DELETE", fmt.Sprintf("/users/%s", userID), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteUserNotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	userID := uuid.New()
	mockService.On("Delete", mock.Anything, userID).Return(apperrors.ErrNotFound)

	router := setupRouter()
	router.DELETE("/users/:id", handler.DeleteUser)

	w := performRequest(router, "DELETE", fmt.Sprintf("/users/%s", userID), nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// =====================
// ChangePassword Tests
// =====================

func TestChangePasswordSuccess(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	userID := uuid.New()
	req := dto.ChangePasswordRequest{
		OldPassword: "old-password",
		NewPassword: "new-password",
	}

	mockService.On("ChangePassword", mock.Anything, userID, req).Return(nil)

	router := setupRouter()
	router.PATCH("/users/change-password", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.ChangePassword(c)
	})

	w := performRequest(router, "PATCH", "/users/change-password", req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestChangePasswordUnauthorized(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	router := setupRouter()
	router.PATCH("/users/change-password", handler.ChangePassword)

	req := dto.ChangePasswordRequest{
		OldPassword: "old-password",
		NewPassword: "new-password",
	}

	w := performRequest(router, "PATCH", "/users/change-password", req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChangePasswordBadRequest(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	router := setupRouter()
	router.PATCH("/users/change-password", handler.ChangePassword)

	// Missing required fields
	req := map[string]string{
		"oldPassword": "old",
	}

	w := performRequest(router, "PATCH", "/users/change-password", req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =====================
// UpdateUser Tests
// =====================

func TestUpdateUserSuccess(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	userID := uuid.New()
	newName := "John Updated"
	req := dto.UpdateUserRequest{
		FirstName: &newName,
	}

	resUser := &models.User{ID: userID, FirstName: newName}

	mockService.On("Update", mock.Anything, userID, req).Return(resUser, nil)

	router := setupRouter()
	router.PATCH("/users/profile", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.UpdateUser(c)
	})

	w := performRequest(router, "PATCH", "/users/profile", req)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.NoError(t, err)
	assert.Equal(t, newName, actualUser.FirstName)
	mockService.AssertExpectations(t)
}

// =====================
// UpdateAvatar Tests
// =====================

func TestUpdateAvatarSuccess(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	userID := uuid.New()
	avatarPath := "/uploads/avatars/avatar.png"

	mockService.On("UpdateAvatar", mock.Anything, userID, mock.Anything, "avatar.png").
		Return(avatarPath, nil)

	router := setupRouter()
	router.POST("/users/avatar", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.UpdateAvatar(c)
	})

	// Create multipart form data
	body, contentType := createMultipartForm("avatar", "avatar.png", "fake-image-binary")
	w := performRequestWithContentType(router, "POST", "/users/avatar", body, contentType)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, avatarPath, resp["avatarUrl"])
	mockService.AssertExpectations(t)
}

func TestUpdateAvatarNoFile(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	router := setupRouter()
	router.POST("/users/avatar", handler.UpdateAvatar)

	w := performRequestWithContentType(router, "POST", "/users/avatar", nil, "multipart/form-data")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAvatarServiceError(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService, nil)

	userID := uuid.New()
	mockService.On("UpdateAvatar", mock.Anything, userID, mock.Anything, "avatar.png").
		Return("", errors.New("storage error"))

	router := setupRouter()
	router.POST("/users/avatar", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		handler.UpdateAvatar(c)
	})

	body, contentType := createMultipartForm("avatar", "avatar.png", "fake-image-binary")
	w := performRequestWithContentType(router, "POST", "/users/avatar", body, contentType)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
