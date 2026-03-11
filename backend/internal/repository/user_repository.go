package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, userID uuid.UUID, params UpdateUserParams) (*models.User, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL *string) (*models.User, error)
	FindAll(ctx context.Context, filter UserListFilter) (users []models.User, total int64, err error)
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

var _ UserRepository = (*userRepository)(nil)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type UserListFilter struct {
	FirstName  string
	LastName   string
	Email      string
	Role       models.Role
	Pagination PaginationParams
}

type UpdateUserParams struct {
	FirstName *string
	LastName  *string
	IsActive  *bool
}

const (
	userIDQuery = "id = ?"
)

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, userID uuid.UUID, params UpdateUserParams) (*models.User, error) {
	updates := make(map[string]interface{})

	if params.FirstName != nil {
		updates["first_name"] = *params.FirstName
	}
	if params.LastName != nil {
		updates["last_name"] = *params.LastName
	}
	if params.IsActive != nil {
		updates["is_active"] = *params.IsActive
	}

	var user models.User
	result := r.db.WithContext(ctx).Model(&user).Where(userIDQuery, userID).Updates(updates)
	if result.Error != nil {
		return nil, mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	if err := r.db.WithContext(ctx).First(&user, userIDQuery, userID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}

	return &user, nil
}

func (r *userRepository) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL *string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Model(&user).Where(userIDQuery, userID).Update("avatar_url", avatarURL)
	if result.Error != nil {
		return nil, mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	if err := r.db.WithContext(ctx).First(&user, userIDQuery, userID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}

	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, userIDQuery, userID)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userIDQuery, userID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, filter UserListFilter) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	if filter.FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+sanitizeLike(filter.FirstName)+"%")
	}
	if filter.LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+sanitizeLike(filter.LastName)+"%")
	}
	if filter.Email != "" {
		query = query.Where("email ILIKE ?", "%"+sanitizeLike(filter.Email)+"%")
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	err := query.
		Order(filter.Pagination.GetOrderBy()).
		Offset(filter.Pagination.GetOffset()).
		Limit(filter.Pagination.Limit).
		Find(&users).Error

	if err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return users, total, nil
}

func (r *userRepository) ChangePassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where(userIDQuery, userID).
		Update("password_hash", newHashedPassword)

	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
