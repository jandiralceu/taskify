package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/dto"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/jandiralceu/taskify/internal/pkg"
	"github.com/jandiralceu/taskify/internal/repository"
)

// UserService defines the business logic contract for user management.
type UserService interface {
	// Create registers a new user, ensuring the password is securely hashed.
	Create(ctx context.Context, user *models.User) error
	// FindAll returns a paginated list of users based on search criteria.
	FindAll(ctx context.Context, req dto.GetUserListRequest) (*dto.PaginatedResponse[models.User], error)
	// FindByID retrieves a single user by their unique UUID.
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	// FindByEmail locates a user using their unique email address.
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	// ChangePassword updates a user's password after verifying the old one.
	ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error
	// Delete removes a user from the system.
	Delete(ctx context.Context, userID uuid.UUID) error
	// Update modifies general user profile information.
	Update(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*models.User, error)
	// UpdateAvatar uploads a profile picture and returns the URL/path.
	UpdateAvatar(ctx context.Context, userID uuid.UUID, file io.Reader, filename string) (string, error)
}

type userService struct {
	userRepo   repository.UserRepository
	hasher     pkg.PasswordHasher
	uploadPath string
}

var _ UserService = (*userService)(nil)

// NewUserService initializes a UserService with its required repository and hasher dependencies.
func NewUserService(userRepo repository.UserRepository, hasher pkg.PasswordHasher, uploadPath string) UserService {
	return &userService{userRepo: userRepo, hasher: hasher, uploadPath: uploadPath}
}

// Create performs password hashing using the injected hasher before persisting the user through the repository.
func (s *userService) Create(ctx context.Context, user *models.User) error {
	hashedPassword, err := s.hasher.Hash(user.PasswordHash)
	if err != nil {
		return apperrors.ErrInternal
	}
	user.PasswordHash = hashedPassword

	return s.userRepo.Create(ctx, user)
}

// FindAll delegates the retrieval of the user list to the repository.
func (s *userService) FindAll(ctx context.Context, req dto.GetUserListRequest) (*dto.PaginatedResponse[models.User], error) {
	filter := repository.UserListFilter{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
		Pagination: repository.PaginationParams{
			Page:  req.GetPage(),
			Limit: req.GetLimit(),
			Sort:  req.GetSort("created_at", "first_name", "last_name", "email"),
			Order: req.GetOrder(),
		},
	}

	users, total, err := s.userRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return dto.NewPaginatedResponse(users, total, filter.Pagination.Page, filter.Pagination.Limit), nil
}

// FindByID retrieves a user by their primary key from the repository.
func (s *userService) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// FindByEmail retrieves a user by their email address for authentication or identification purposes.
func (s *userService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// Delete removes the user record identified by the unique ID.
func (s *userService) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

// Update handles partial updates to the user profile.
func (s *userService) Update(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*models.User, error) {
	params := repository.UpdateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  req.IsActive,
	}

	return s.userRepo.Update(ctx, userID, params)
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	match, err := s.hasher.Verify(req.OldPassword, user.PasswordHash)
	if err != nil {
		return err
	}

	if !match {
		return apperrors.ErrUnauthorized
	}

	newHashedPassword, err := s.hasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.ChangePassword(ctx, userID, newHashedPassword)
}

func (s *userService) UpdateAvatar(ctx context.Context, userID uuid.UUID, file io.Reader, filename string) (string, error) {
	// 1. Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// 2. Generate unique filename
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("avatar_%s%s", uuid.New().String(), ext)
	avatarDir := filepath.Join(s.uploadPath, "avatars")
	filePath := filepath.Join(avatarDir, uniqueName)

	// 3. Ensure directory exists
	if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
		return "", apperrors.ErrStorage
	}

	// 4. Create file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", apperrors.ErrStorage
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", apperrors.ErrStorage
	}

	// 5. Update user in DB
	oldAvatar := user.AvatarURL
	if _, err := s.userRepo.UpdateAvatar(ctx, userID, &filePath); err != nil {
		os.Remove(filePath) // Cleanup
		return "", err
	}

	// 6. Delete old avatar if exists
	if oldAvatar != nil {
		os.Remove(*oldAvatar)
	}

	return filePath, nil
}
