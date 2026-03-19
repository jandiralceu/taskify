package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepositoryUpdateSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	repo := NewTaskRepository(gormDB)
	taskID := uuid.New()
	title := "Updated Title"
	desc := "Updated Description"
	status := models.TaskStatusInProgress
	priority := models.TaskPriorityHigh

	params := UpdateTaskParams{
		Title:       &title,
		Description: &desc,
		Status:      &status,
		Priority:    &priority,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "description"=$1,"priority"=$2,"status"=$3,"title"=$4,"updated_at"=$5 WHERE id = $6`)).
		WithArgs(desc, priority, status, title, sqlmock.AnyArg(), taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE id = $1 ORDER BY "tasks"."id" LIMIT $2`)).
		WithArgs(taskID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "priority", "updated_at"}).
			AddRow(taskID, title, desc, status, priority, time.Now()))

	task, err := repo.Update(context.Background(), taskID, params)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, title, task.Title)
	assert.Equal(t, status, task.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaskRepositoryUpdateNotFound(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	repo := NewTaskRepository(gormDB)
	taskID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks"`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	task, err := repo.Update(context.Background(), taskID, UpdateTaskParams{})

	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, apperrors.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
