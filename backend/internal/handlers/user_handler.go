package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/taskify/internal/apperrors"
	"github.com/jandiralceu/taskify/internal/dto"
	"github.com/jandiralceu/taskify/internal/middleware"
	"github.com/jandiralceu/taskify/internal/models"
	"github.com/jandiralceu/taskify/internal/service"
)

// UserHandler provides endpoints for general user management tasks.
type UserHandler struct {
	userService service.UserService
	_           *models.User
}

// NewUserHandler creates a new instance of UserHandler with the specified service.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// FindAllUsers godoc
// @Summary      List all users
// @Description  Get a paginated list of users. Supports filtering by name, email, and role.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page    query  int     false  "Page number"
// @Param        limit   query  int     false  "Number of items per page"
// @Param        first_name  query  string  false  "Filter by first name"
// @Param        last_name   query  string  false  "Filter by last name"
// @Param        email       query  string  false  "Filter by email"
// @Param        role        query  string  false  "Filter by role (admin, employee)"
// @Success      200     {object}  dto.UserListResponse
// @Failure      400     {object}  ProblemDetails
// @Failure      401     {object}  ProblemDetails
// @Failure      403     {object}  ProblemDetails
// @Failure      429     {object}  ProblemDetails
// @Failure      500     {object}  ProblemDetails
// @Security     Bearer
// @Router       /users [get]
func (h *UserHandler) FindAllUsers(c *gin.Context) {
	var req dto.GetUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	users, err := h.userService.FindAll(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// FindUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieve details for a single user using their unique ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ProblemDetails
// @Failure      401  {object}  ProblemDetails
// @Failure      403  {object}  ProblemDetails
// @Failure      404  {object}  ProblemDetails
// @Failure      429  {object}  ProblemDetails
// @Security     Bearer
// @Router       /users/{id} [get]
func (h *UserHandler) FindUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	user, err := h.userService.FindByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Permanently remove a user from the system by their unique ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  ProblemDetails
// @Failure      401  {object}  ProblemDetails
// @Failure      403  {object}  ProblemDetails
// @Failure      404  {object}  ProblemDetails
// @Failure      500  {object}  ProblemDetails
// @Failure      429  {object}  ProblemDetails
// @Security     Bearer
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ChangePassword godoc
// @Summary      Change own password
// @Description  Updates the authenticated user's password. Requires the old password for verification.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.ChangePasswordRequest true "Password change data"
// @Success      204 "No Content"
// @Failure      400  {object}  ProblemDetails
// @Failure      401  {object}  ProblemDetails
// @Failure      403  {object}  ProblemDetails
// @Failure      500  {object}  ProblemDetails
// @Failure      429  {object}  ProblemDetails
// @Security     Bearer
// @Router       /users/change-password [patch]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		RespondWithError(c, apperrors.ErrUnauthorized)
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, req); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
