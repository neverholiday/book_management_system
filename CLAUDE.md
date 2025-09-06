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

### API Versioning and Route Organization
- **Versioning Structure**: Use `/api/v{version}/` for all API routes (e.g., `/api/v1/`, `/api/v2/`)
- **Version Groups**: Create separate Echo groups for each API version to enable parallel versions
- **Backward Compatibility**: Maintain older API versions during transition periods
- **Route Groups**: Group related routes together within each version
- **Health Check**: Always include `/healthz` endpoint for service health monitoring (no versioning)

#### API Group Structure Pattern
```go
// Main API group
apiGroup := e.Group("/api")

// Version-specific groups
v1Group := apiGroup.Group("/v1")
v2Group := apiGroup.Group("/v2")  // Future version

// Feature groups within versions
authV1 := v1Group.Group("/auth")
usersV1 := v1Group.Group("/users")
booksV1 := v1Group.Group("/books")

authV2 := v2Group.Group("/auth")  // Different implementation
```

#### Versioning Best Practices
- **Start with v1**: Always begin with `/api/v1/` even for initial release
- **Breaking Changes**: Increment major version for breaking changes (v1 → v2)
- **Non-breaking Changes**: Keep same version for backward-compatible additions
- **Migration Path**: Provide clear migration documentation between versions
- **Deprecation Notice**: Add deprecation headers for older versions before removal

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

## Git and Version Control Rules

### .gitignore Requirements
Always create a comprehensive `.gitignore` file that excludes:

```gitignore
# Go-specific files
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/
*.test
*.out
go.work

# Environment and configuration files (CRITICAL)
.env
.env.*
*.env
config.local.*
settings.local.*

# IDE and editor files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs and temporary files
*.log
logs/
tmp/
temp/
.tmp/

# Database files
*.db
*.sqlite
*.sqlite3

# Build artifacts
build/
*.tar.gz
*.zip

# Certificate and security files
*.pem
*.key
*.crt
*.cert
```

### Git Security Rules
- **Never Commit Secrets**: Environment files, API keys, passwords, certificates
- **Early Setup**: Create `.gitignore` before first commit to prevent accidents
- **Comprehensive Coverage**: Include IDE, OS, build artifacts, and language-specific files
- **Documentation**: Comment critical exclusions in `.gitignore` for team awareness

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

## Go Coding Conventions

### Variable Naming
- **Short Names**: Use abbreviated variable names for commonly used variables
  - `config` → `cfg` 
  - `database` → `db`
  - `context` → `ctx`
  - `request` → `req`
  - `response` → `resp`
- **Modern Types**: Use modern Go type aliases (Go 1.18+)
  - `interface{}` → `any`
- **Clean Structure**: Prioritize readability and concise code

### Struct Organization
- **No Line Breaks**: Never add line breaks within ANY struct definition
- **Compact Format**: All struct fields should be consecutive without empty lines
- **Universal Rule**: Applies to all structs (Config, models, DTOs, etc.)

```go
// Good - compact and clean (applies to ALL structs)
type User struct {
    ID          string     `gorm:"column:id"`
    Email       string     `gorm:"column:email"`
    PasswordHash string    `gorm:"column:password_hash"`
    FirstName   string     `gorm:"column:first_name"`
    LastName    string     `gorm:"column:last_name"`
    Role        string     `gorm:"column:role"`
    Status      string     `gorm:"column:status"`
    CreatedDate time.Time  `gorm:"column:created_date"`
    UpdatedDate time.Time  `gorm:"column:updated_date"`
    DeletedDate *time.Time `gorm:"column:deleted_date"`
}

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

// Bad - unnecessary line breaks (NEVER do this for any struct)
type User struct {
    ID          string     `gorm:"column:id"`
    Email       string     `gorm:"column:email"`
    
    FirstName   string     `gorm:"column:first_name"`
    LastName    string     `gorm:"column:last_name"`
    
    CreatedDate time.Time  `gorm:"column:created_date"`
    UpdatedDate time.Time  `gorm:"column:updated_date"`
}
```

### Error Handling
- **Main Function**: Always use `panic(err)` for error handling in main.go
- **Consistent Pattern**: Use the standard `if err != nil { panic(err) }` format
- **No log.Fatal**: Never use `log.Fatal()` - use panic instead

```go
// Good - panic pattern
if err := envconfig.Process("BOOKMS", &cfg); err != nil {
    panic(err)
}

if err := e.Start(cfg.ServerAddress()); err != nil {
    panic(err)
}

// Bad - log.Fatal pattern  
if err := envconfig.Process("BOOKMS", &cfg); err != nil {
    log.Fatal("Failed to load configuration:", err)
}
```

### Package Location Best Practices

#### Shared Packages (`pkg/` Directory)
- **Multi-Service Usage**: Use `pkg/` for code that can be imported by multiple services within the monorepo
- **Shared Libraries**: Place reusable packages that multiple `cmd/` services need in `pkg/`
- **Interface-based Design**: Use interfaces to decouple pkg packages from specific implementations
- **No Service Dependencies**: `pkg/` packages should never import from `cmd/` directories
- **Examples**: `pkg/auth`, `pkg/logger`, `pkg/metrics`, `pkg/database` for cross-service functionality

#### Service-Specific Packages (`cmd/<service>/` Directory)
- **Single Service Usage**: Keep packages that are only used by one service within that service's directory
- **Service Boundaries**: Each service maintains its own `models/`, `apis/`, `repositories/` packages
- **No Cross-Service Imports**: Services should not import from other `cmd/<service>/` directories
- **Clear Ownership**: Service-specific packages belong to that service and follow its lifecycle
- **Examples**: `cmd/server_api/models/`, `cmd/worker/handlers/`, `cmd/cli/commands/`

#### Decision Guidelines
- **Multiple Services Need It**: Use `pkg/` (e.g., JWT auth, logging, metrics)
- **Single Service Only**: Keep in `cmd/<service>/` (e.g., specific API models, handlers)
- **Future Sharing Possible**: Start in `cmd/<service>/`, move to `pkg/` when needed by second service
- **External Dependencies**: Business logic specific to one service stays in `cmd/<service>/`

### Package Naming Best Practices
- **Descriptive Names**: Use specific, descriptive package names that indicate purpose
- **Avoid Generic Names**: Never use `utils`, `helpers`, `common`, or `shared` packages
- **Purpose-driven**: Name packages by what they provide, not how they're used
- **Examples**: `auth` (not `utils`), `http` (not `helpers`), `storage` (not `common`)

### Code Style Priorities
- **Conciseness**: Favor shorter, cleaner code over verbose implementations
- **Readability**: Code should be self-documenting through good naming
- **Consistency**: Apply naming conventions uniformly across the codebase

### Code Cleanliness Rules
- **No Comments by Default**: Write clean, self-documenting code without comments
- **Descriptive Naming**: Use clear function and variable names that explain purpose
- **Comment Only When Complex**: Add comments only for complicated logic that needs explanation
- **Self-Documenting Code**: Code should be readable without requiring comments to understand

```go
// Good - clean, self-documenting code
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ? AND deleted_date IS NULL", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Bad - unnecessary comments
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    // Create a user variable to store the result
    var user models.User
    // Query the database for user by email, excluding soft deleted records
    err := r.db.Where("email = ? AND deleted_date IS NULL", email).First(&user).Error
    // Check if there was an error
    if err != nil {
        return nil, err
    }
    // Return the user
    return &user, nil
}
```

### When to Add Comments
- **Complex Business Logic**: Algorithms or calculations that aren't immediately obvious
- **External Dependencies**: Integration points with third-party services
- **Performance Considerations**: Code optimized in non-obvious ways
- **Regulatory Requirements**: Code that implements specific compliance rules

## Logging System Rules

### Standard Logging with slog
- **Primary Logger**: Always use Go's standard `log/slog` package for structured logging
- **JSON Format**: Use `slog.NewJSONHandler` for production-ready structured logs
- **Default Setup**: Set slog as the default logger with `slog.SetDefault()`
- **Context Logging**: Use `slog.InfoContext()` and `slog.ErrorContext()` for request-scoped logging

### GORM Logging Integration
- **Package**: Use `github.com/orandin/slog-gorm` for GORM-slog integration
- **Simple Setup**: `gormLogger := slogGorm.New()` provides automatic slog integration
- **No Custom Loggers**: Avoid writing custom GORM logger adapters - use established packages

### Echo Logging Integration
- **Request Logging**: Use `middleware.RequestLoggerWithConfig` with custom `LogValuesFunc`
- **Structured Data**: Log method, URI, status, latency, remote IP as structured fields
- **Error Distinction**: Separate logging for successful requests vs errors
- **Context Awareness**: Use request context for all HTTP-related logging

### slog Structured Logging Style
- **Multi-line**: Break slog calls across lines when multiple fields are present
- **Field Alignment**: Each key-value pair on its own line
- **Descriptive Keys**: Use clear, consistent field names

```go
// Good - structured multi-line slog
slog.Info("Database connection established",
    "max_open_conns", cfg.DBMaxOpenConns,
    "max_idle_conns", cfg.DBMaxIdleConns,
    "conn_max_lifetime", cfg.DBConnMaxLifetime,
)
```

## Docker Development Rules

### Development Database Setup
- **Docker Compose**: Always provide `docker-compose.yml` for development databases
- **Generic Credentials**: Use simple, generic credentials for development (postgres/password/myapp)
- **Health Checks**: Include health check configuration for database readiness
- **Persistent Volumes**: Use named volumes for data persistence
- **Standard Ports**: Expose standard database ports (5432 for PostgreSQL)

### Container Standards
- **Alpine Images**: Prefer Alpine variants for smaller image sizes (postgres:15-alpine)
- **Container Naming**: Use descriptive but simple container names
- **Environment Variables**: Match container environment with application .env configuration

### Database Schema Management
- **Schema Location**: Store complete database schema in `init/init.sql`
- **Docker Integration**: Mount init.sql to `/docker-entrypoint-initdb.d/init.sql` in PostgreSQL container
- **Automatic Initialization**: PostgreSQL will execute init.sql on first container startup
- **Schema Control**: Maintain full control over table creation, indexes, and constraints
- **No Auto-Migration**: Avoid GORM AutoMigrate in favor of explicit schema management

### Database Initialization Pattern
```yaml
# docker-compose.yml volume configuration
volumes:
  - postgres_data:/var/lib/postgresql/data
  - ./init/init.sql:/docker-entrypoint-initdb.d/init.sql
```

### Schema File Structure
```sql
-- init/init.sql structure
-- Create tables with explicit column definitions
CREATE TABLE table_name (
    id VARCHAR(100) PRIMARY KEY,
    -- ... other columns
    created_date timestamptz NOT NULL,
    updated_date timestamptz NOT NULL,
    deleted_date timestamptz
);

-- Create indexes for performance
CREATE INDEX idx_table_field ON table_name(field);
CREATE UNIQUE INDEX idx_table_unique ON table_name(unique_field);
```

### Database Reset Process
- **Complete Reset**: `docker-compose down -v && docker-compose up -d`
- **Volume Removal**: `-v` flag removes named volumes and triggers re-initialization
- **Fresh Start**: New containers will execute init.sql automatically

```yaml
# Example docker-compose.yml
services:
  postgres:
    image: postgres:15-alpine
    container_name: postgres_db
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d myapp"]
      interval: 5s
      timeout: 3s
      retries: 5

volumes:
  postgres_data:
```

## API Implementation Patterns

### API Structure and Organization
- **Package Location**: All API handlers go in `cmd/<service>/apis/` directory
- **File Naming**: Use `<entity>.go` naming convention (e.g., `healthz.go`, `user.go`, `book.go`)
- **Constructor Pattern**: Use `New<Entity>API(dependencies)` constructor functions
- **Setup Method**: Each API struct should have a `Setup(group *echo.Group)` method for route registration

### API Constructor Pattern
```go
// Example API structure
type HealthzAPI struct {
    db *gorm.DB
}

func NewHealthzAPI(db *gorm.DB) *HealthzAPI {
    return &HealthzAPI{
        db: db,
    }
}

func (api *HealthzAPI) Setup(group *echo.Group) {
    group.GET("/healthz", api.healthCheck)
}
```

### Route Registration Pattern
- **Group-based Registration**: Use Echo groups for organizing routes
- **Dependency Injection**: Pass required dependencies (database, services) to constructors
- **Method Binding**: Bind HTTP methods to struct methods in Setup function
- **Chainable Setup**: API setup should be chainable with constructor

```go
// Main.go route registration pattern
rootg := e.Group("")
apis.NewHealthzAPI(db).Setup(rootg)

apiV1 := e.Group("/api/v1")
apis.NewUserAPI(db).Setup(apiV1)
apis.NewBookAPI(db).Setup(apiV1)
```

### API Method Conventions
- **Handler Signature**: All handlers should follow `func (api *API) method(c echo.Context) error`
- **Error Handling**: Use consistent error response format
- **Response Format**: Follow established success/error response patterns from previous sections
- **Context Usage**: Use `c.Request().Context()` for request-scoped operations

## Multi-line Function Call Style Rules

### Consistent Multi-line Formatting
- **Universal Application**: Apply to all function calls, method calls, and struct initialization
- **Parameter Alignment**: Each parameter gets its own line when breaking across lines
- **Closing Parenthesis**: Place closing parenthesis on separate line aligned with function call

```go
// Good - consistent multi-line style
err := envconfig.Process(
    "BOOKMS",
    &cfg,
)

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
```

## CLAUDE.md Maintenance Rules

### Conflict Prevention
- **Always Check**: Before adding new rules or sections, search the entire document for conflicting guidance
- **Remove Conflicts**: When updating practices, remove or update conflicting sections
- **Consistency Review**: New additions should align with existing patterns and conventions
- **Single Source of Truth**: Each topic should have only one authoritative section

### Update Process
1. Read the entire CLAUDE.md before making changes
2. Search for existing guidance on the same topic
3. Identify and resolve conflicts
4. Update or remove outdated sections
5. Verify consistency across the document
