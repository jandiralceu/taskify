//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jandiralceu/taskify/internal/dto"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserManagementIntegration(t *testing.T) {
	ts, db, cleanup := setupApp(t)
	defer cleanup()

	baseURL := ts.URL
	adminEmail := "superadmin@example.com"
	password := "SecurePass123!"

	// Create and login as Admin
	signUpUser(t, baseURL, "Super", "Admin", adminEmail, password, "admin")
	adminToken, _ := signInUser(t, baseURL, adminEmail, password)

	t.Run("Admin can list and search users", func(t *testing.T) {
		// Create a second user to search for
		signUpUser(t, baseURL, "Searchable", "User", "search@example.com", "Pass123!", "employee")

		resp := authedRequest(t, "GET", baseURL+"/api/v1/users?first_name=Searchable", adminToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var listResp dto.UserListResponse
		_ = json.NewDecoder(resp.Body).Decode(&listResp)
		assert.NotEmpty(t, listResp.Data)
		assert.Equal(t, "Searchable", listResp.Data[0].FirstName)
	})

	t.Run("User can change their own password", func(t *testing.T) {
		userEmail := "changepass@example.com"
		signUpUser(t, baseURL, "Pass", "Changer", userEmail, "old-pass-123", "employee")
		token, _ := signInUser(t, baseURL, userEmail, "old-pass-123")

		req := dto.ChangePasswordRequest{
			OldPassword: "old-pass-123",
			NewPassword: "new-pass-456",
		}

		resp := authedRequest(t, "PATCH", baseURL+"/api/v1/users/change-password", token, req)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify login with new password
		newToken, _ := signInUser(t, baseURL, userEmail, "new-pass-456")
		assert.NotEmpty(t, newToken)
	})

	t.Run("Admin can delete users", func(t *testing.T) {
		userEmail := "tobedeleted@example.com"
		signUpUser(t, baseURL, "To be", "Deleted", userEmail, "Pass123!", "employee")

		var user models.User
		db.Where("email = ?", userEmail).First(&user)

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/users/%s", baseURL, user.ID), adminToken, nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify user is gone
		var checkUser models.User
		err := db.Where("id = ?", user.ID).First(&checkUser).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("RBAC: Employee cannot delete users", func(t *testing.T) {
		opEmail := "op_delete_bad@example.com"
		signUpUser(t, baseURL, "Employee", "User", opEmail, "Pass123!", "employee")
		opToken, _ := signInUser(t, baseURL, opEmail, "Pass123!")

		var user models.User
		db.Where("email = ?", adminEmail).First(&user)

		resp := authedRequest(t, "DELETE", fmt.Sprintf("%s/api/v1/users/%s", baseURL, user.ID), opToken, nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
