# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Documentation

Comprehensive documentation for this repository is available in the `docs/` directory

Always refer to these specification documents when implementing features to ensure consistency with the defined architecture and data models.

## Project Status

This is a Go backend monorepo using Echo framework, PostgreSQL database with GORM, and envconfig for configuration management. The project follows a service-oriented architecture within a monorepo structure.

## Technology Stack

- **Language**: Go 1.19+
- **Framework**: Echo (Go web framework)
- **Database**: PostgreSQL 12+
- **ORM**: GORM
- **Configuration**: envconfig with required environment variables

## Project Structure Rules

### Monorepo Structure
```
your-app/
├── cmd/
│   └── <service_name>/
│       ├── apis/
│       │   └── <entity>.go
│       ├── models/
│       │   └── <entity>.go
│       ├── repositories/
│       │   └── <entity>.go
│       └── main.go
├── go.mod
├── go.sum
├── .env
└── README.md
```

### Service Organization Rules

1. **Service Location**: All services live under `cmd/<service_name>/`
2. **Self-contained Services**: Each service contains its own `apis/`, `models/`, `repositories/`
3. **Entry Point**: Each service has its own `main.go` with embedded configuration
4. **Naming Convention**: Use descriptive service names like `server_api`, `worker`, `cli`, `server_admin`

### Directory Structure Within Services

- **`apis/`** - HTTP handlers and route definitions
- **`models/`** - Database models and business entities  
- **`repositories/`** - Data access layer with CRUD operations
- **`main.go`** - Application entry point with configuration and setup

## Configuration Rules

### Environment Variables Pattern
- **Prefix**: Use `PROGPREFIX_` for all environment variables
- **Required**: All environment variables must be marked as `required:"true"`
- **Database Variables**: Always include these core database settings:
  ```bash
  PROGPREFIX_DB_HOST=localhost
  PROGPREFIX_DB_PORT=5432
  PROGPREFIX_DB_USER=username
  PROGPREFIX_DB_PASSWORD=password
  PROGPREFIX_DB_NAME=dbname
  PROGPREFIX_DB_MAX_OPEN_CONNS=25
  PROGPREFIX_DB_MAX_IDLE_CONNS=5
  PROGPREFIX_DB_CONN_MAX_LIFETIME=300
  ```

### Configuration Structure
- **Location**: Embed configuration struct in each service's `main.go`
- **Type Safety**: Use envconfig with proper Go types
- **Helper Methods**: Include `DSN()` and `ServerAddress()` methods
- **Connection Pool**: Always configure database connection pooling

```go
// Example configuration structure
type Config struct {
    DBHost            string `envconfig:"DB_HOST" required:"true"`
    DBPort            int    `envconfig:"DB_PORT" required:"true"`
    DBUser            string `envconfig:"DB_USER" required:"true"`
    DBPassword        string `envconfig:"DB_PASSWORD" required:"true"`
    DBName            string `envconfig:"DB_NAME" required:"true"`
    DBMaxOpenConns    int    `envconfig:"DB_MAX_OPEN_CONNS" required:"true"`
    DBMaxIdleConns    int    `envconfig:"DB_MAX_IDLE_CONNS" required:"true"`
    DBConnMaxLifetime int    `envconfig:"DB_CONN_MAX_LIFETIME" required:"true"`
}
```

## Database Model Rules

### Primary Key and ID Fields
- **ID Type**: All primary keys must be `VARCHAR(100)` - never use auto-increment integers
- **ID Generation**: Let the application generate ID values, not the database
- **GORM ID Field**: Use `string` type in Go structs for ID fields

### Column Design Rules
- **No Default Values**: Never specify DEFAULT values in database schema - let the program choose defaults
- **NOT NULL Required**: All entity properties must be NOT NULL (except soft delete fields)
- **Explicit Nullability**: Only allow NULL for optional fields like soft delete timestamps

### Timestamp Fields
- **Time Type**: Use `timestamptz` (PostgreSQL) for all timestamp fields
- **Standard Fields**: Include `created_date`, `updated_date`, and `deleted_date` for all entities
- **Soft Delete**: `deleted_date` should be NULL by default - when NOT NULL, record is considered deleted
- **UTC Timezone**: All time fields should use UTC

### GORM Tags
- **Column Mapping Only**: Use only `gorm:"column:<column_name>"` tags
- **No Constraints**: Do not include database constraints in struct tags
- **No JSON Tags**: Handle JSON serialization separately if needed
- **Explicit Naming**: Always specify column names explicitly

```go
// Example model structure following the rules
type Entity struct {
    ID          string     `gorm:"column:id"`
    Name        string     `gorm:"column:name"`
    CreatedDate time.Time  `gorm:"column:created_date"`
    UpdatedDate time.Time  `gorm:"column:updated_date"`
    DeletedDate *time.Time `gorm:"column:deleted_date"`
}
```

### Model Conventions
- **Standard Fields**: Include `ID` (string), `CreatedDate`, `UpdatedDate`, `DeletedDate` for all entities
- **Soft Delete Logic**: Use `deleted_date IS NULL` for active records, `IS NOT NULL` for deleted
- **Package Location**: Models go in `cmd/<service>/models/<entity>.go`

## Repository Pattern Rules

### Repository Structure
- **Interface Definition**: Define repository interface in the same file
- **CRUD Operations**: Include Create, Read, Update, Delete methods
- **Query Methods**: Add specific query methods as needed
- **Error Handling**: Return errors from repository methods

```go
// Example repository pattern
type EntityRepository struct {
    db *gorm.DB
}

func NewEntityRepository(db *gorm.DB) *EntityRepository {
    return &EntityRepository{db: db}
}

// CRUD methods
func (r *EntityRepository) Create(entity *models.Entity) error
func (r *EntityRepository) GetByID(id uint) (*models.Entity, error)
func (r *EntityRepository) GetAll(limit, offset int) ([]models.Entity, error)
func (r *EntityRepository) Update(entity *models.Entity) error
func (r *EntityRepository) Delete(id uint) error
```

## API Handler Rules

### Handler Structure
- **Constructor Pattern**: Use `NewEntityHandler(repo)` constructor
- **HTTP Methods**: Map to CRUD operations (POST=Create, GET=Read, PUT=Update, DELETE=Delete)
- **Response Format**: Follow consistent JSON response format
- **Status Codes**: Use appropriate HTTP status codes

### Response Format Rules
- **Success Response (200 OK)**: Always return `{"data": <response_data>, "message": "<message from API>"}`
- **Error Response (4xx, 5xx)**: Always return `{"message": "<response or some error message>"}`
- **No Success Field**: Never include `success` boolean field in responses
- **Consistent Structure**: Success responses always have both `data` and `message` fields
- **Error Simplicity**: Error responses only contain `message` field

### Route Organization
- **RESTful Patterns**: Follow REST conventions for URLs
- **API Versioning**: Use `/api/v1/` prefix for main API routes
- **Route Groups**: Group related routes together
- **Health Check**: Always include `/healthz` endpoint for service health monitoring

```go
// Example route structure
// Health check (no versioning)
e.GET("/healthz", handler.HealthCheck)

// Main API routes
api := e.Group("/api/v1")
api.POST("/entities", handler.CreateEntity)           // CREATE
api.GET("/entities", handler.GetEntities)             // READ (all)
api.GET("/entities/:id", handler.GetEntity)           // READ (by ID)
api.PUT("/entities/:id", handler.UpdateEntity)        // UPDATE
api.DELETE("/entities/:id", handler.DeleteEntity)     // DELETE
```

### Health Check Requirements
- **Endpoint Path**: `/healthz` (not versioned, directly on root)
- **Method**: GET only
- **Authentication**: None required (public endpoint)
- **Response**: JSON with service status, timestamp, and version
- **Purpose**: Used by load balancers, monitoring systems, and deployment tools

## Development Commands Pattern

```bash
# Service-specific commands
go run cmd/<service_name>/main.go
go build -o bin/<service_name> cmd/<service_name>/main.go
go test ./cmd/<service_name>/...

# Repository-wide commands
go mod tidy
go test ./...
go fmt ./...
golangci-lint run
```

## Database Setup Rules

### UTC Enforcement
- **Application Level**: Set `time.Local = time.UTC` in main.go
- **Database Connection**: Include `TimeZone=UTC` in DSN
- **GORM Configuration**: Use custom `NowFunc` that returns UTC

### Connection Pool Configuration
- **Required Settings**: Always configure MaxOpenConns, MaxIdleConns, ConnMaxLifetime
- **Environment Driven**: Make pool settings configurable via environment variables
- **Logging**: Log connection pool configuration at startup

## Architecture Guidelines

### Monorepo Benefits
- **Shared Dependencies**: Single go.mod for all services
- **Code Reuse**: Common utilities can be shared across services
- **Atomic Changes**: Deploy related changes across services together
- **Consistent Tooling**: Same linting, testing, and build processes

### Service Independence
- **Self-contained**: Each service should be independently runnable
- **Database Per Service**: Each service can have its own database/schema
- **Configuration Isolation**: Each service manages its own configuration

### Scaling Patterns
- **Add New Services**: Create new directories under `cmd/`
- **Shared Libraries**: Use `pkg/` for shared utilities
- **Internal Packages**: Use `internal/` for non-exportable code

## Code Quality Rules

### Error Handling
- **Repository Level**: Handle database errors in repositories
- **Handler Level**: Convert to appropriate HTTP responses
- **Consistent Format**: Use consistent error response format

### Testing Strategy
- **Unit Tests**: Test repository and handler logic separately
- **Integration Tests**: Test database interactions
- **Service Tests**: Test complete HTTP endpoints

### Security Considerations
- **Input Validation**: Validate all input at handler level
- **SQL Injection**: GORM provides protection, but validate inputs
- **Environment Variables**: Never commit .env files
- **Rate Limiting**: Implement per-endpoint rate limiting as needed

This structure provides a scalable, maintainable foundation for Go backend services in a monorepo architecture.

## Task Management and Documentation

### Development Journal - Sprint System
- **Location**: Always create a `journal/` directory in the project root
- **Sprint Files**: Use `SPRINT_YYYY-MM-DD.md` format (e.g., `SPRINT_2025-08-10.md`, `SPRINT_2025-08-17.md`)
- **Task Organization**: Group related tasks in sprints for flexible management
- **Simple Format**: Keep task descriptions concise with short action descriptions

### Sprint File Format
```markdown
# SPRINT YYYY-MM-DD - Sprint Name

**Started:** YYYY-MM-DD
**Status:** 🚧 IN PROGRESS / ✅ COMPLETED

## Tasks

- [x] **Task XX**: Short description
  - Brief action taken
  - Key files modified

- [ ] **Task XX**: Short description
  - What needs to be done

## Progress: X/Y completed
```

### Sprint Management
- **Flexibility**: Tasks can be moved between sprints as needed
- **Breakdown**: Large tasks can be split across multiple sprints
- **Progress Tracking**: Simple checkbox format with progress counter
- **README Index**: Maintain journal/README.md with current sprint status

This approach provides lightweight task tracking while maintaining development audit trail.
