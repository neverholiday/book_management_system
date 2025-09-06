package apis

import (
	"book-management-system/cmd/server_api/models"
	"book-management-system/cmd/server_api/repositories"
	"book-management-system/pkg/auth"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthAPI struct {
	userRepo *repositories.UserRepository
	jwt      *auth.JWT
	authMw   *auth.Middleware
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	User         *UserProfile     `json:"user"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresAt    time.Time        `json:"expires_at"`
}

type UserProfile struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

func NewAuthAPI(userRepo *repositories.UserRepository, jwt *auth.JWT) *AuthAPI {
	return &AuthAPI{
		userRepo: userRepo,
		jwt:      jwt,
		authMw:   auth.NewMiddleware(jwt),
	}
}

func (api *AuthAPI) Setup(group *echo.Group) {
	group.POST("/register", api.register)
	group.POST("/login", api.login)
	group.POST("/refresh", api.refresh)
	group.GET("/profile", api.profile, api.authMw.RequireAuth())
}

func (api *AuthAPI) register(c echo.Context) error {
	var req RegisterRequest
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
			Message: "Email already registered",
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
		Role:         "member",
		Status:       "active",
	}
	err = api.userRepo.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error creating user account",
		})
	}
	tokens, err := api.jwt.GenerateTokenPair(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error generating authentication tokens",
		})
	}
	response := models.Response{
		Data: AuthResponse{
			User: &UserProfile{
				ID:        user.ID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
				Status:    user.Status,
			},
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresAt:    time.Now().Add(time.Hour * 24),
		},
		Message: "Account created successfully",
	}
	return c.JSON(http.StatusCreated, response)
}

func (api *AuthAPI) login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid request format",
		})
	}
	user, err := api.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusUnauthorized, models.Response{
				Message: "Invalid email or password",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error during authentication",
		})
	}
	if user.Status != "active" {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Message: "Account is not active",
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid email or password",
		})
	}
	tokens, err := api.jwt.GenerateTokenPair(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error generating authentication tokens",
		})
	}
	response := models.Response{
		Data: AuthResponse{
			User: &UserProfile{
				ID:        user.ID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
				Status:    user.Status,
			},
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresAt:    time.Now().Add(time.Hour * 24),
		},
		Message: "Login successful",
	}
	return c.JSON(http.StatusOK, response)
}

func (api *AuthAPI) refresh(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Invalid request format",
		})
	}
	userID, err := api.jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid refresh token",
		})
	}
	user, err := api.userRepo.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Message: "User not found",
		})
	}
	if user.Status != "active" {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Message: "Account is not active",
		})
	}
	tokens, err := api.jwt.GenerateTokenPair(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error generating authentication tokens",
		})
	}
	response := models.Response{
		Data: AuthResponse{
			User: &UserProfile{
				ID:        user.ID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
				Status:    user.Status,
			},
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresAt:    time.Now().Add(time.Hour * 24),
		},
		Message: "Tokens refreshed successfully",
	}
	return c.JSON(http.StatusOK, response)
}

func (api *AuthAPI) profile(c echo.Context) error {
	claims := api.authMw.GetUserFromContext(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Message: "Authentication required",
		})
	}
	user, err := api.userRepo.GetByID(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Message: "User not found",
		})
	}
	response := models.Response{
		Data: UserProfile{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Status:    user.Status,
		},
		Message: "User profile retrieved successfully",
	}
	return c.JSON(http.StatusOK, response)
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + time.Now().Format("000000")
}