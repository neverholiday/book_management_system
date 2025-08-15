package main

import (
	"book-management-system/cmd/server_api/apis"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf(
		"%s:%s",
		c.ServerHost,
		c.ServerPort,
	)
}

func init() {
	os.Setenv("TZ", "UTC")
}

func main() {

	var cfg Config
	err := envconfig.Process(
		"BOOKMS",
		&cfg,
	)
	if err != nil {
		panic(err)
	}

	gormLogger := slogGorm.New()

	db, err := gorm.Open(
		postgres.Open(
			cfg.DSN(),
		),
		&gorm.Config{
			Logger: gormLogger,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		},
	)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	defer sqlDB.Close()

	sqlDB.SetMaxOpenConns(
		cfg.DBMaxOpenConns,
	)
	sqlDB.SetMaxIdleConns(
		cfg.DBMaxIdleConns,
	)
	sqlDB.SetConnMaxLifetime(
		time.Duration(
			cfg.DBConnMaxLifetime,
		) * time.Second,
	)

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	slog.Info(
		"Database connection established",
		"max_open_conns", cfg.DBMaxOpenConns,
		"max_idle_conns", cfg.DBMaxIdleConns,
		"conn_max_lifetime", cfg.DBConnMaxLifetime,
	)

	e := echo.New()
	e.Use(
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogStatus:   true,
			LogURI:      true,
			LogError:    true,
			LogLatency:  true,
			LogMethod:   true,
			LogRemoteIP: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if v.Error == nil {
					slog.InfoContext(c.Request().Context(), "request",
						"method", v.Method,
						"uri", v.URI,
						"status", v.Status,
						"latency", v.Latency,
						"remote_ip", v.RemoteIP,
					)
				} else {
					slog.ErrorContext(c.Request().Context(), "request_error",
						"method", v.Method,
						"uri", v.URI,
						"status", v.Status,
						"latency", v.Latency,
						"remote_ip", v.RemoteIP,
						"error", v.Error,
					)
				}
				return nil
			},
		}),
	)
	e.Use(
		middleware.Recover(),
	)

	rootg := e.Group("")
	apis.NewHealthzAPI(
		db,
	).Setup(
		rootg,
	)

	slog.Info("Server starting", "address", cfg.ServerAddress())
	err = e.Start(
		cfg.ServerAddress(),
	)
	if err != nil {
		panic(err)
	}

}
