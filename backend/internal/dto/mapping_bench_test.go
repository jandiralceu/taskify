package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/models"
)

// UserResponse is what we would typically return to the client
type UserResponse struct {
	ID        uuid.UUID   `json:"id"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     string      `json:"email"`
	Role      string `json:"role"`
	CreatedAt time.Time   `json:"createdAt"`
}

// MapUserToResponse simulates the conversion logic
func MapUserToResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      "employee",
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
