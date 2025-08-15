package apis

import (
	"book-management-system/cmd/server_api/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HealthzAPI struct {
	db *gorm.DB
}

func NewHealthzAPI(db *gorm.DB) *HealthzAPI {
	return &HealthzAPI{
		db: db,
	}
}

func (a *HealthzAPI) Setup(g *echo.Group) {
	g.GET("/healthz", a.checkHealth)
}

func (a *HealthzAPI) checkHealth(c echo.Context) error {

	sqlDB, err := a.db.DB()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			models.Response{
				Message: err.Error(),
			},
		)
	}

	err = sqlDB.Ping()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			models.Response{
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.Response{
			Message: "healthy",
		},
	)
}
