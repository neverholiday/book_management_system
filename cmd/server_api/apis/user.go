package apis

import (
	"book-management-system/cmd/server_api/models"
	"book-management-system/cmd/server_api/repositories"
	"book-management-system/pkg/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserAPI struct {
	userRepo *repositories.UserRepository
	authMw   *auth.Middleware
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Role      string `json:"role" validate:"required,oneof=admin member"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Role      *string `json:"role,omitempty" validate:"omitempty,oneof=admin member"`
	Status    *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

type UserListResponse struct {
	Users  []UserDetail `json:"users"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

type UserDetail struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

func NewUserAPI(userRepo *repositories.UserRepository, authMw *auth.Middleware) *UserAPI {
	return &UserAPI{
		userRepo: userRepo,
		authMw:   authMw,
	}
}

func (api *UserAPI) Setup(group *echo.Group) {
	group.POST("", api.createUser, api.authMw.RequireAdmin())
	group.GET("", api.getUsers, api.authMw.RequireAdmin())
	group.GET("/:id", api.getUserByID, api.authMw.RequireAdmin())
	group.PUT("/:id", api.updateUser, api.authMw.RequireAdmin())
	group.DELETE("/:id", api.deleteUser, api.authMw.RequireAdmin())
}

func (api *UserAPI) createUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid request format",
		})
	}
	exists, err := api.userRepo.EmailExists(req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error checking email availability",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, models.Response{
			Message: "Email already exists",
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error processing password",
		})
	}
	user := &models.User{
		ID:           generateID(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		Status:       "active",
	}
	err = api.userRepo.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error creating user",
		})
	}
	response := models.Response{
		Data: UserDetail{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Role:        user.Role,
			Status:      user.Status,
			CreatedDate: user.CreatedDate,
			UpdatedDate: user.UpdatedDate,
		},
		Message: "User created successfully",
	}
	return c.JSON(http.StatusCreated, response)
}

func (api *UserAPI) getUsers(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}
	role := c.QueryParam("role")
	status := c.QueryParam("status")
	var users []models.User
	var err error
	if role != "" {
		users, err = api.userRepo.GetByRole(role, limit, offset)
	} else if status != "" {
		users, err = api.userRepo.GetByStatus(status, limit, offset)
	} else {
		users, err = api.userRepo.GetAll(limit, offset)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error retrieving users",
		})
	}
	total, err := api.userRepo.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error counting users",
		})
	}
	userDetails := make([]UserDetail, len(users))
	for i, user := range users {
		userDetails[i] = UserDetail{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Role:        user.Role,
			Status:      user.Status,
			CreatedDate: user.CreatedDate,
			UpdatedDate: user.UpdatedDate,
		}
	}
	response := models.Response{
		Data: UserListResponse{
			Users:  userDetails,
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
		Message: "Users retrieved successfully",
	}
	return c.JSON(http.StatusOK, response)
}

func (api *UserAPI) getUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := api.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, models.Response{
				Message: "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error retrieving user",
		})
	}
	response := models.Response{
		Data: UserDetail{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Role:        user.Role,
			Status:      user.Status,
			CreatedDate: user.CreatedDate,
			UpdatedDate: user.UpdatedDate,
		},
		Message: "User retrieved successfully",
	}
	return c.JSON(http.StatusOK, response)
}

func (api *UserAPI) updateUser(c echo.Context) error {
	id := c.Param("id")
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid request format",
		})
	}
	user, err := api.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, models.Response{
				Message: "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error retrieving user",
		})
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	err = api.userRepo.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error updating user",
		})
	}
	response := models.Response{
		Data: UserDetail{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Role:        user.Role,
			Status:      user.Status,
			CreatedDate: user.CreatedDate,
			UpdatedDate: user.UpdatedDate,
		},
		Message: "User updated successfully",
	}
	return c.JSON(http.StatusOK, response)
}

func (api *UserAPI) deleteUser(c echo.Context) error {
	id := c.Param("id")
	_, err := api.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, models.Response{
				Message: "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error retrieving user",
		})
	}
	err = api.userRepo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error deleting user",
		})
	}
	response := models.Response{
		Message: "User deleted successfully",
	}
	return c.JSON(http.StatusOK, response)
}