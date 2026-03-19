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
	// GetUserPermissions retrieves all permissions for a user across all their roles.
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
}

// refreshTokenCacheKeyPrefix is the Redis key prefix used for refresh tokens.
// Format: refresh_token:{userID}:{token}
const refreshTokenCacheKeyPrefix = "refresh_token:"

type userService struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	hasher     pkg.PasswordHasher
	uploadPath string
	cache      pkg.CacheManager
}

var _ UserService = (*userService)(nil)

// NewUserService initializes a UserService with its required repository and hasher dependencies.
func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, hasher pkg.PasswordHasher, uploadPath string, cache pkg.CacheManager) UserService {
	return &userService{userRepo: userRepo, roleRepo: roleRepo, hasher: hasher, uploadPath: uploadPath, cache: cache}
}

// Create performs password hashing using the injected hasher before persisting the user through the repository.
func (s *userService) Create(ctx context.Context, user *models.User) error {
	hashedPassword, err := s.hasher.Hash(user.PasswordHash)
	if err != nil {
		return apperrors.ErrInternal
	}
	user.PasswordHash = hashedPassword

	// Assign specified or default role 'employee' if no roles assigned
	if len(user.Roles) == 0 {
		roleName := "employee"
		if user.Role != "" {
			roleName = user.Role
		}
		role, err := s.roleRepo.FindByName(ctx, roleName)
		if err == nil && role != nil {
			user.Roles = append(user.Roles, *role)
		}
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	if len(user.Roles) > 0 {
		user.Role = user.Roles[0].Name
	}
	return nil
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

	for i := range users {
		if len(users[i].Roles) > 0 {
			users[i].Role = users[i].Roles[0].Name
		}
	}

	return dto.NewPaginatedResponse(users, total, filter.Pagination.Page, filter.Pagination.Limit), nil
}

// FindByID retrieves a user by their primary key from the repository.
func (s *userService) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(user.Roles) > 0 {
		user.Role = user.Roles[0].Name
	}

	return user, nil
}

// FindByEmail retrieves a user by their email address for authentication or identification purposes.
func (s *userService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// Delete removes the user record and invalidates all their active sessions in the cache.
func (s *userService) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return err
	}

	prefix := fmt.Sprintf("%s%s:", refreshTokenCacheKeyPrefix, userID.String())
	_ = s.cache.DeletePrefix(ctx, prefix)

	return nil
}

// Update handles partial updates to the user profile.
func (s *userService) Update(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*models.User, error) {
	params := repository.UpdateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  req.IsActive,
	}

	user, err := s.userRepo.Update(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	if len(user.Roles) > 0 {
		user.Role = user.Roles[0].Name
	}

	return user, nil
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
	diskPath := filepath.Join(avatarDir, uniqueName)

	// publicPath is what gets stored in the DB and returned to clients.
	// The router serves the uploads directory at /uploads, so a path like
	// /uploads/avatars/avatar_<uuid>.jpg is directly accessible via the API.
	publicPath := "/uploads/avatars/" + uniqueName

	// 3. Ensure directory exists
	if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
		return "", apperrors.ErrStorage
	}

	// 4. Create file on disk
	dst, err := os.Create(diskPath)
	if err != nil {
		return "", apperrors.ErrStorage
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", apperrors.ErrStorage
	}

	// 5. Persist public path in DB
	oldAvatar := user.AvatarURL
	if _, err := s.userRepo.UpdateAvatar(ctx, userID, &publicPath); err != nil {
		os.Remove(diskPath) // Cleanup on DB failure
		return "", err
	}

	// 6. Delete old avatar file from disk if one existed
	if oldAvatar != nil {
		// oldAvatar is a public path (/uploads/avatars/...), reconstruct disk path
		oldDiskPath := filepath.Join(s.uploadPath, "avatars", filepath.Base(*oldAvatar))
		os.Remove(oldDiskPath)
	}

	return publicPath, nil
}

func (s *userService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	var permissions []string
	permMap := make(map[string]bool)

	for _, role := range roles {
		perms, err := s.roleRepo.GetPermissionsByRole(ctx, role.Name)
		if err != nil {
			continue
		}
		for _, p := range perms {
			permStr := fmt.Sprintf("%s:%s", p.Resource, p.Action)
			if !permMap[permStr] {
				permMap[permStr] = true
				permissions = append(permissions, permStr)
			}
		}
	}

	if permissions == nil {
		return []string{}, nil
	}

	return permissions, nil
}
