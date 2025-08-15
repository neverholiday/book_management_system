package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	UserContextKey = "user"
)

type Middleware struct {
	jwt *JWT
}

func NewMiddleware(jwt *JWT) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}

func (m *Middleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := m.extractToken(c)
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Authorization header is required",
				})
			}
			claims, err := m.jwt.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid or expired token",
				})
			}
			c.Set(UserContextKey, claims)
			return next(c)
		}
	}
}

func (m *Middleware) RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := m.GetUserFromContext(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Authentication required",
				})
			}
			if user.Role != role {
				return c.JSON(http.StatusForbidden, map[string]string{
					"message": "Insufficient permissions",
				})
			}
			return next(c)
		}
	}
}

func (m *Middleware) RequireAdmin() echo.MiddlewareFunc {
	return m.RequireRole("admin")
}

func (m *Middleware) extractToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func (m *Middleware) GetUserFromContext(c echo.Context) *Claims {
	user, ok := c.Get(UserContextKey).(*Claims)
	if !ok {
		return nil
	}
	return user
}