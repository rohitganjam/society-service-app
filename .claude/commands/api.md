# Create API Endpoint

Create a new API endpoint: $ARGUMENTS

## Process

### Step 1: Analyze Requirements

Determine from the arguments:
- **HTTP Method**: GET, POST, PUT, PATCH, DELETE
- **Resource**: What entity (orders, vendors, residents, etc.)
- **Action**: CRUD operation or custom action
- **Route**: RESTful path

### Step 2: Create Files in Order

#### A. Model (if new entity)

**Location**: `backend/internal/models/{entity}.go`

```go
package models

import "time"

type {Entity} struct {
    ID          string    `json:"id"`
    // Add fields based on requirements
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Request/Response types
type Create{Entity}Request struct {
    // Request fields with validation tags
    Name string `json:"name" binding:"required,min=2"`
}

type Update{Entity}Request struct {
    Name *string `json:"name,omitempty"`
}

type {Entity}Filter struct {
    Status string `form:"status"`
    Limit  int    `form:"limit"`
    Offset int    `form:"offset"`
}
```

#### B. Repository Interface

**Location**: `backend/internal/repositories/{entity}_repository.go`

```go
package repositories

import (
    "context"
    "society-service-app/backend/internal/models"
)

type {Entity}Repository interface {
    Create(ctx context.Context, entity *models.{Entity}) (*models.{Entity}, error)
    GetByID(ctx context.Context, id string) (*models.{Entity}, error)
    Update(ctx context.Context, entity *models.{Entity}) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter *models.{Entity}Filter) ([]*models.{Entity}, error)
}
```

#### C. Repository Implementation

**Location**: `backend/internal/repositories/{entity}_repository_pg.go`

```go
package repositories

import (
    "context"
    "society-service-app/backend/internal/database"
    "society-service-app/backend/internal/models"
)

type pg{Entity}Repository struct {
    db *database.DB
}

func New{Entity}Repository(db *database.DB) {Entity}Repository {
    return &pg{Entity}Repository{db: db}
}

func (r *pg{Entity}Repository) Create(ctx context.Context, entity *models.{Entity}) (*models.{Entity}, error) {
    query := `
        INSERT INTO {entities} (field1, field2, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    err := r.db.Pool.QueryRow(ctx, query,
        entity.Field1,
        entity.Field2,
    ).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)

    if err != nil {
        return nil, err
    }
    return entity, nil
}

func (r *pg{Entity}Repository) GetByID(ctx context.Context, id string) (*models.{Entity}, error) {
    query := `
        SELECT id, field1, field2, created_at, updated_at
        FROM {entities}
        WHERE id = $1 AND deleted_at IS NULL
    `
    var entity models.{Entity}
    err := r.db.Pool.QueryRow(ctx, query, id).Scan(
        &entity.ID,
        &entity.Field1,
        &entity.Field2,
        &entity.CreatedAt,
        &entity.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &entity, nil
}
```

#### D. Service

**Location**: `backend/internal/services/{entity}_service.go`

```go
package services

import (
    "context"
    "errors"
    "society-service-app/backend/internal/models"
    "society-service-app/backend/internal/repositories"
)

var (
    Err{Entity}NotFound = errors.New("{entity} not found")
    Err{Entity}Exists   = errors.New("{entity} already exists")
)

type {Entity}Service struct {
    repo repositories.{Entity}Repository
}

func New{Entity}Service(repo repositories.{Entity}Repository) *{Entity}Service {
    return &{Entity}Service{repo: repo}
}

func (s *{Entity}Service) Create(ctx context.Context, req *models.Create{Entity}Request) (*models.{Entity}, error) {
    // Business logic and validation
    entity := &models.{Entity}{
        Field1: req.Field1,
        Field2: req.Field2,
    }

    return s.repo.Create(ctx, entity)
}

func (s *{Entity}Service) GetByID(ctx context.Context, id string) (*models.{Entity}, error) {
    entity, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, Err{Entity}NotFound
    }
    return entity, nil
}
```

#### E. Handler

**Location**: `backend/internal/handlers/{entity}_handler.go`

```go
package handlers

import (
    "context"
    "errors"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "society-service-app/backend/internal/models"
    "society-service-app/backend/internal/services"
    "society-service-app/backend/internal/utils"
)

type {Entity}Handler struct {
    service *services.{Entity}Service
}

func New{Entity}Handler(service *services.{Entity}Service) *{Entity}Handler {
    return &{Entity}Handler{service: service}
}

func (h *{Entity}Handler) Create(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    var req models.Create{Entity}Request
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), nil)
        return
    }

    result, err := h.service.Create(ctx, &req)
    if err != nil {
        if errors.Is(err, services.Err{Entity}Exists) {
            utils.RespondError(c, http.StatusConflict, "ALREADY_EXISTS", err.Error(), nil)
            return
        }
        utils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create", nil)
        return
    }

    utils.RespondSuccess(c, http.StatusCreated, result, "{Entity} created successfully")
}

func (h *{Entity}Handler) GetByID(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    id := c.Param("id")
    if id == "" {
        utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", "ID is required", nil)
        return
    }

    result, err := h.service.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, services.Err{Entity}NotFound) {
            utils.RespondError(c, http.StatusNotFound, "NOT_FOUND", "{Entity} not found", nil)
            return
        }
        utils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get", nil)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, result, "")
}

func (h *{Entity}Handler) List(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    var filter models.{Entity}Filter
    if err := c.ShouldBindQuery(&filter); err != nil {
        utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), nil)
        return
    }

    results, err := h.service.List(ctx, &filter)
    if err != nil {
        utils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list", nil)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, results, "")
}
```

#### F. Register Routes

**Update**: `backend/cmd/api/main.go`

```go
// In the routes setup section:
{entity}Repo := repositories.New{Entity}Repository(db)
{entity}Service := services.New{Entity}Service({entity}Repo)
{entity}Handler := handlers.New{Entity}Handler({entity}Service)

v1.GET("/{entities}", {entity}Handler.List)
v1.GET("/{entities}/:id", {entity}Handler.GetByID)
v1.POST("/{entities}", {entity}Handler.Create)
v1.PUT("/{entities}/:id", {entity}Handler.Update)
v1.DELETE("/{entities}/:id", {entity}Handler.Delete)
```

### Step 3: Create Tests

#### Service Tests

**Location**: `backend/internal/services/{entity}_service_test.go`

```go
package services

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "society-service-app/backend/internal/models"
)

type Mock{Entity}Repository struct {
    mock.Mock
}

func (m *Mock{Entity}Repository) Create(ctx context.Context, entity *models.{Entity}) (*models.{Entity}, error) {
    args := m.Called(ctx, entity)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.{Entity}), args.Error(1)
}

func (m *Mock{Entity}Repository) GetByID(ctx context.Context, id string) (*models.{Entity}, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.{Entity}), args.Error(1)
}

func TestCreate{Entity}(t *testing.T) {
    tests := []struct {
        name      string
        req       *models.Create{Entity}Request
        mockSetup func(*Mock{Entity}Repository)
        wantErr   bool
    }{
        {
            name: "success",
            req:  &models.Create{Entity}Request{Field1: "test"},
            mockSetup: func(m *Mock{Entity}Repository) {
                m.On("Create", mock.Anything, mock.Anything).Return(&models.{Entity}{
                    ID:     "123",
                    Field1: "test",
                }, nil)
            },
            wantErr: false,
        },
        {
            name: "validation error",
            req:  &models.Create{Entity}Request{Field1: ""},
            mockSetup: func(m *Mock{Entity}Repository) {},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(Mock{Entity}Repository)
            tt.mockSetup(mockRepo)

            svc := New{Entity}Service(mockRepo)
            result, err := svc.Create(context.Background(), tt.req)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
            }
        })
    }
}
```

#### Handler Tests

**Location**: `backend/internal/handlers/{entity}_handler_test.go`

```go
package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCreate{Entity}Handler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    tests := []struct {
        name       string
        body       interface{}
        mockSetup  func(*MockService)
        wantStatus int
    }{
        {
            name: "success",
            body: map[string]string{"field1": "test"},
            mockSetup: func(m *MockService) {
                m.On("Create", mock.Anything, mock.Anything).Return(&models.{Entity}{ID: "123"}, nil)
            },
            wantStatus: http.StatusCreated,
        },
        {
            name:       "invalid json",
            body:       "invalid",
            mockSetup:  func(m *MockService) {},
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := new(MockService)
            tt.mockSetup(mockSvc)

            handler := New{Entity}Handler(mockSvc)

            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

            bodyBytes, _ := json.Marshal(tt.body)
            c.Request = httptest.NewRequest("POST", "/{entities}", bytes.NewReader(bodyBytes))
            c.Request.Header.Set("Content-Type", "application/json")

            handler.Create(c)

            assert.Equal(t, tt.wantStatus, w.Code)
        })
    }
}
```

### Step 4: Verify

Run all checks:

```bash
cd backend
make fmt          # Format code
make lint         # Check linting
make test         # Run tests
```
