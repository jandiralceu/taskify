package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/dto"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// =====================
// Create Tests
// =====================

func TestUserRepositoryCreateSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	userID := uuid.New()
	user := &models.User{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: "hashed_password",
		Role:         models.RoleAdmin,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("first_name","last_name","email","password_hash","role","is_active","avatar_url") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id","created_at","updated_at"`)).
		WithArgs(user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(userID, time.Now(), time.Now()))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryCreateError(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	user := &models.User{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		PasswordHash: "hashed_password",
		Role:         models.RoleAdmin,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	err := repo.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrConflict, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// FindByID Tests
// =====================

func TestUserRepositoryFindByIDSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	userID := uuid.New()
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password_hash", "role", "is_active", "created_at", "updated_at"}).
		AddRow(userID, "John", "Doe", "john@example.com", "hash", "admin", true, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	user, err := repo.FindByID(context.Background(), userID)

	assert.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindByIDNotFound(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	userID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(userID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByID(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, apperrors.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// FindByEmail Tests
// =====================

func TestUserRepositoryFindByEmailSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	userID := uuid.New()
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password_hash", "role", "is_active", "created_at", "updated_at"}).
		AddRow(userID, "John", "Doe", "john@example.com", "hash", "admin", true, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("john@example.com", 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail(context.Background(), "john@example.com")

	assert.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindByEmailNotFound(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("unknown@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail(context.Background(), "unknown@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, apperrors.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// Delete Tests
// =====================

func TestUserRepositoryDeleteSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryDeleteNotFound(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// FindAll Tests
// =====================

func TestUserRepositoryFindAllSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	userID1 := uuid.New()
	userID2 := uuid.New()
	now := time.Now()

	req := dto.GetUserListRequest{
		PaginationRequest: dto.PaginationRequest{
			Page:  1,
			Limit: 10,
		},
	}

	// Count query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Find query
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password_hash", "role", "is_active", "created_at", "updated_at"}).
		AddRow(userID1, "User", "One", "user1@example.com", "hash1", "employee", true, now, now).
		AddRow(userID2, "User", "Two", "user2@example.com", "hash2", "employee", true, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(rows)

	filter := UserListFilter{
		Pagination: PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at",
			Order: "desc",
		},
	}

	users, total, err := repo.FindAll(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
	assert.Equal(t, "User", users[0].FirstName)
	assert.Equal(t, "User", users[1].FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindAllEmpty(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)

	req := dto.GetUserListRequest{
		PaginationRequest: dto.PaginationRequest{
			Page:  1,
			Limit: 10,
		},
	}

	// Count query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Find query
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password_hash", "role", "is_active", "created_at", "updated_at"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(rows)

	filter := UserListFilter{
		Pagination: PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at",
			Order: "desc",
		},
	}

	users, total, err := repo.FindAll(context.Background(), filter)

	assert.NoError(t, err)
	assert.Len(t, users, 0)
	assert.Equal(t, int64(0), total)
	mock.ExpectationsWereMet()
}

// =====================
// Update Tests
// =====================

func TestUserRepositoryUpdateSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)
	ctx := context.Background()
	userID := uuid.New()
	newName := "Updated Name"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "first_name"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(newName, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name"}).AddRow(userID, newName))

	params := UpdateUserParams{
		FirstName: &newName,
	}

	user, err := repo.Update(ctx, userID, params)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newName, user.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryUpdateNotFound(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)
	ctx := context.Background()
	userID := uuid.New()
	newName := "Updated Name"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "first_name"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(newName, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	params := UpdateUserParams{
		FirstName: &newName,
	}

	user, err := repo.Update(ctx, userID, params)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryUpdateAvatarSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)
	ctx := context.Background()
	userID := uuid.New()
	avatarURL := "http://example.com/avatar.png"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "avatar_url"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(&avatarURL, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "avatar_url"}).AddRow(userID, avatarURL))

	user, err := repo.UpdateAvatar(ctx, userID, &avatarURL)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, &avatarURL, user.AvatarURL)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryUpdateAvatarNotFound(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(gormDB)
	ctx := context.Background()
	userID := uuid.New()
	avatarURL := "http://example.com/avatar.png"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "avatar_url"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(&avatarURL, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	user, err := repo.UpdateAvatar(ctx, userID, &avatarURL)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// ChangePassword Tests
// =====================

func TestUserRepositoryChangePasswordSuccess(t *testing.T) {
	db, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	newHash := "new-hashed-password"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "password_hash"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(newHash, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.ChangePassword(ctx, userID, newHash)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func TestUserRepositoryChangePasswordNotFound(t *testing.T) {
	db, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()
	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "password_hash"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs("hash", sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.ChangePassword(ctx, userID, "hash")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
	mock.ExpectationsWereMet()
}
