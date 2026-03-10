//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthFlowAndRBAC(t *testing.T) {
	ts, _, cleanup := setupApp(t)
	defer cleanup()

	baseURL := ts.URL

	t.Run("Full Auth Flow: Register -> SignIn -> Refresh -> SignOut", func(t *testing.T) {
		email := "tester@example.com"
		password := "Password123!"

		// Register
		signUpUser(t, baseURL, "First", "Last", email, password, "admin")

		// SignIn
		accessToken, refreshToken := signInUser(t, baseURL, email, password)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)

		// Refresh Token
		refreshBody := map[string]string{"refreshToken": refreshToken}
		resp := authedRequest(t, "POST", baseURL+"/api/v1/auth/refresh", "", refreshBody)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var refreshResult map[string]string
		json.NewDecoder(resp.Body).Decode(&refreshResult)
		newAccessToken := refreshResult["accessToken"]
		newRefreshToken := refreshResult["refreshToken"]
		assert.NotEmpty(t, newAccessToken)
		assert.NotEqual(t, refreshToken, newRefreshToken)

		// SignOut
		signOutBody := map[string]string{"refreshToken": newRefreshToken}
		resp = authedRequest(t, "POST", baseURL+"/api/v1/auth/signout", "", signOutBody)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("RBAC: Admin can access restricted routes, Guest cannot", func(t *testing.T) {
		adminEmail := "admin_test@example.com"
		password := "Pass123!"

		// Create Admin
		signUpUser(t, baseURL, "Admin", "User", adminEmail, password, "admin")
		adminToken, _ := signInUser(t, baseURL, adminEmail, password)

		// Admin attempts to list users
		resp := authedRequest(t, "GET", baseURL+"/api/v1/users", adminToken, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Admin should be able to list users")

		// Guest attempts to list users (no token)
		resp = authedRequest(t, "GET", baseURL+"/api/v1/users", "", nil)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Guest should be unauthorized")
	})

	t.Run("RBAC: Accessing non-existent route returns 404 (with auth)", func(t *testing.T) {
		adminEmail := "admin_404@example.com"
		password := "Pass123!"
		signUpUser(t, baseURL, "Admin", "404", adminEmail, password, "admin")
		adminToken, _ := signInUser(t, baseURL, adminEmail, password)

		resp := authedRequest(t, "GET", baseURL+"/api/v1/non-existent", adminToken, nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
