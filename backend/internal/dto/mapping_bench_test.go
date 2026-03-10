package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
)

// UserResponse is what we would typically return to the client
type UserResponse struct {
	ID        uuid.UUID   `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Role      models.Role `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
}

// MapUserToResponse simulates the conversion logic
func MapUserToResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

func BenchmarkMappingLargeList(b *testing.B) {
	// Setup: create a list of 1000 models
	size := 1000
	users := make([]models.User, size)
	for i := 0; i < size; i++ {
		users[i] = models.User{
			ID:        uuid.New(),
			FirstName: "Jermaine",
			LastName:  "Cole",
			Email:     "jcole@fakeemail.com",
			Role:      models.RoleEmployee,
			CreatedAt: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		responses := make([]UserResponse, size)
		for j := 0; j < size; j++ {
			responses[j] = MapUserToResponse(users[j])
		}
	}
}
