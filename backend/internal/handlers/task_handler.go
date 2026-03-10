package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/middleware"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask godoc
// @Summary      Create a new task
// @Description  Creates a new task associated with the authenticated user.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTaskRequest true "Task details"
// @Success      201 {object} models.Task
// @Security     Bearer
// @Router       /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	userID := middleware.GetUserID(c)
	task, err := h.taskService.Create(c.Request.Context(), userID, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Updates an existing task's details.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path   string true "Task UUID"
// @Param        request body dto.UpdateTaskRequest true "Updated task details"
// @Success      200 {object} models.Task
// @Security     Bearer
// @Router       /tasks/{id} [patch]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	task, err := h.taskService.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetTask godoc
// @Summary      Get task by ID
// @Description  Retrieves a single task by its unique ID, including notes and attachments.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path   string true "Task UUID"
// @Success      200 {object} models.Task
// @Security     Bearer
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	task, err := h.taskService.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Permanently removes a task.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path   string true "Task UUID"
// @Success      204 "No Content"
// @Security     Bearer
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.taskService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTasks godoc
// @Summary      List tasks
// @Description  Get a paginated list of tasks with optional filtering.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        status      query string false "Filter by status"
// @Param        priority    query string false "Filter by priority"
// @Param        search      query string false "Search in title/description"
// @Param        assigned_to query string false "Filter by assigned user ID"
// @Param        page        query int    false "Page number"
// @Param        limit       query int    false "Items per page"
// @Success      200 {object} dto.PaginatedResponse[models.Task]
// @Security     Bearer
// @Router       /tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	var req dto.GetTaskListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	tasks, err := h.taskService.GetAll(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// AddNote godoc
// @Summary      Add a note to a task
// @Description  Creates a new note associated with the specified task.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task UUID"
// @Param        request body dto.CreateTaskNoteRequest true "Note content"
// @Success      201 {object} models.TaskNote
// @Security     Bearer
// @Router       /tasks/{id}/notes [post]
func (h *TaskHandler) AddNote(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.CreateTaskNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	userID := middleware.GetUserID(c)
	note, err := h.taskService.AddNote(c.Request.Context(), taskID, userID, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetNotes godoc
// @Summary      Get all notes for a task
// @Description  Retrieves all notes associated with a specific task.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task UUID"
// @Success      200 {array} models.TaskNote
// @Security     Bearer
// @Router       /tasks/{id}/notes [get]
func (h *TaskHandler) GetNotes(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	notes, err := h.taskService.GetNotes(c.Request.Context(), taskID)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, notes)
}

// UpdateNote godoc
// @Summary      Update a task note
// @Description  Updates the content of an existing note. Only the author can update it.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        noteId path string true "Note UUID"
// @Param        request body dto.UpdateTaskNoteRequest true "New note content"
// @Success      200 {object} models.TaskNote
// @Security     Bearer
// @Router       /tasks/notes/{noteId} [patch]
func (h *TaskHandler) UpdateNote(c *gin.Context) {
	idStr := c.Param("noteId")
	noteID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.UpdateTaskNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	userID := middleware.GetUserID(c)
	note, err := h.taskService.UpdateNote(c.Request.Context(), noteID, userID, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNote godoc
// @Summary      Delete a task note
// @Description  Permanently removes a note. Only the author can delete it.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        noteId path string true "Note UUID"
// @Success      204 "No Content"
// @Security     Bearer
// @Router       /tasks/notes/{noteId} [delete]
func (h *TaskHandler) DeleteNote(c *gin.Context) {
	idStr := c.Param("noteId")
	noteID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.taskService.DeleteNote(c.Request.Context(), noteID, userID); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// AddAttachment godoc
// @Summary      Add an attachment to a task
// @Description  Uploads a file and attaches it to the specified task.
// @Tags         tasks
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path   string  true  "Task UUID"
// @Param        file formData file    true  "File to upload"
// @Success      201  {object}  models.TaskAttachment
// @Security     Bearer
// @Router       /tasks/{id}/attachments [post]
func (h *TaskHandler) AddAttachment(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidInput)
		return
	}
	defer file.Close()

	userID := middleware.GetUserID(c)
	attachment, err := h.taskService.AddAttachment(c.Request.Context(), taskID, userID, file, header.Filename, header.Size, header.Header.Get("Content-Type"))
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

// GetAttachments godoc
// @Summary      Get all attachments for a task
// @Description  Retrieves all attachments associated with a specific task.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task UUID"
// @Success      200 {array} models.TaskAttachment
// @Security     Bearer
// @Router       /tasks/{id}/attachments [get]
func (h *TaskHandler) GetAttachments(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	attachments, err := h.taskService.GetAttachments(c.Request.Context(), taskID)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, attachments)
}

// DeleteAttachment godoc
// @Summary      Delete an attachment
// @Description  Permanently removes an attachment and its file from disk. Only the author can delete it.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        attachmentId path string true "Attachment UUID"
// @Success      204 "No Content"
// @Security     Bearer
// @Router       /tasks/attachments/{attachmentId} [delete]
func (h *TaskHandler) DeleteAttachment(c *gin.Context) {
	idStr := c.Param("attachmentId")
	attachmentID, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.taskService.DeleteAttachment(c.Request.Context(), attachmentID, userID); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
