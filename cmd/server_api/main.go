package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBHost                string `envconfig:"DB_HOST" required:"true"`
	DBPort                int    `envconfig:"DB_PORT" required:"true"`
	DBUser                string `envconfig:"DB_USER" required:"true"`
	DBPassword            string `envconfig:"DB_PASSWORD" required:"true"`
	DBName                string `envconfig:"DB_NAME" required:"true"`
	DBMaxOpenConns        int    `envconfig:"DB_MAX_OPEN_CONNS" required:"true"`
	DBMaxIdleConns        int    `envconfig:"DB_MAX_IDLE_CONNS" required:"true"`
	DBConnMaxLifetime     int    `envconfig:"DB_CONN_MAX_LIFETIME" required:"true"`
	ServerHost            string `envconfig:"SERVER_HOST" required:"true"`
	ServerPort            string `envconfig:"SERVER_PORT" required:"true"`
	JWTSecret             string `envconfig:"JWT_SECRET" required:"true"`
	JWTExpiryHours        int    `envconfig:"JWT_EXPIRY_HOURS" required:"true"`
	JWTRefreshExpiryHours int    `envconfig:"JWT_REFRESH_EXPIRY_HOURS" required:"true"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

func init() {
	os.Setenv("TZ", "UTC")
}

func main() {

	var cfg Config
	if err := envconfig.Process("BOOKMS", &cfg); err != nil {
		panic(err)
	}

	db, err := setupDatabase(&cfg)
	if err != nil {
		panic(err)
	}

	e := setupServer()

	setupRoutes(e, db)

	log.Printf("Server starting on %s", cfg.ServerAddress())
	if err := e.Start(cfg.ServerAddress()); err != nil {
		panic(err)
	}
}

func setupDatabase(config *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(config.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.DBConnMaxLifetime) * time.Second)

	log.Printf("Database connection established - MaxOpen: %d, MaxIdle: %d, MaxLifetime: %ds",
		config.DBMaxOpenConns, config.DBMaxIdleConns, config.DBConnMaxLifetime)

	return db, nil
}

func setupServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}

func setupRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]any{
			"data": map[string]any{
				"status":    "healthy",
				"timestamp": time.Now().UTC(),
				"version":   "1.0.0",
			},
			"message": "Service is healthy",
		})
	})

	api := e.Group("/api/v1")

	_ = api
}

