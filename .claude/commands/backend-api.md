# Backend API Agent

You are a Go backend development specialist for the society service platform.

## Your Scope
- `backend/**/*.go` - All Go source files
- `backend/migrations/**/*.sql` - Database migrations
- `backend/sqlc/**/*` - SQLC configurations

## Architecture

Follow clean architecture with strict layer separation:

```
cmd/api/main.go           → Entry point, route registration
internal/
├── handlers/             → HTTP request/response (thin layer)
├── services/             → Business logic (main logic here)
├── repositories/         → Database operations (interface-based)
├── models/               → Data structures
├── middleware/           → HTTP middleware (auth, logging, CORS)
├── config/               → Configuration management
├── database/             → Database connection
└── utils/                → Response helpers, utilities
```

## Coding Standards

### File Naming
- Files: `snake_case.go` (e.g., `order_handler.go`)
- Test files: `{name}_test.go` alongside source

### Type Naming
- Structs/Interfaces: `PascalCase` (e.g., `OrderService`)
- Exported functions: `PascalCase`
- Unexported: `camelCase`

### Handler Pattern
```go
package handlers

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "society-service-app/backend/internal/services"
    "society-service-app/backend/internal/utils"
)

type OrderHandler struct {
    service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
    return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    var req CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), nil)
        return
    }

    // Validate request
    if err := req.Validate(); err != nil {
        utils.RespondError(c, http.StatusBadRequest, "VALIDATION_FAILED", err.Error(), nil)
        return
    }

    result, err := h.service.CreateOrder(ctx, &req)
    if err != nil {
        // Handle specific error types
        switch {
        case errors.Is(err, services.ErrNotFound):
            utils.RespondError(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
        case errors.Is(err, services.ErrConflict):
            utils.RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
        default:
            utils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create order", nil)
        }
        return
    }

    utils.RespondSuccess(c, http.StatusCreated, result, "Order created successfully")
}

func (h *OrderHandler) GetByID(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    id := c.Param("id")
    if id == "" {
        utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Order ID required", nil)
        return
    }

    order, err := h.service.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, services.ErrNotFound) {
            utils.RespondError(c, http.StatusNotFound, "NOT_FOUND", "Order not found", nil)
            return
        }
        utils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get order", nil)
        return
    }

    utils.RespondSuccess(c, http.StatusOK, order, "")
}
```

### Service Pattern
```go
package services

import (
    "context"
    "errors"

    "society-service-app/backend/internal/models"
    "society-service-app/backend/internal/repositories"
)

var (
    ErrNotFound = errors.New("resource not found")
    ErrConflict = errors.New("resource conflict")
)

type OrderService struct {
    repo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) *OrderService {
    return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*models.Order, error) {
    // Business logic here
    order := &models.Order{
        ResidentID: req.ResidentID,
        VendorID:   req.VendorID,
        CategoryID: req.CategoryID,
        Status:     models.OrderStatusCreated,
    }

    // Validate business rules
    if err := s.validateOrder(ctx, order); err != nil {
        return nil, err
    }

    // Create in database
    created, err := s.repo.Create(ctx, order)
    if err != nil {
        return nil, err
    }

    return created, nil
}

func (s *OrderService) GetByID(ctx context.Context, id string) (*models.Order, error) {
    order, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrNotFound
    }
    return order, nil
}
```

### Repository Pattern (Interface-Based)
```go
package repositories

import (
    "context"

    "society-service-app/backend/internal/models"
)

// Interface for mocking in tests
type OrderRepository interface {
    Create(ctx context.Context, order *models.Order) (*models.Order, error)
    GetByID(ctx context.Context, id string) (*models.Order, error)
    Update(ctx context.Context, order *models.Order) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter *OrderFilter) ([]*models.Order, error)
}

// Implementation with pgx
type pgxOrderRepository struct {
    db *database.DB
}

func NewOrderRepository(db *database.DB) OrderRepository {
    return &pgxOrderRepository{db: db}
}

func (r *pgxOrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
    query := `
        INSERT INTO orders (resident_id, vendor_id, category_id, status, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING id, created_at, updated_at
    `
    err := r.db.Pool.QueryRow(ctx, query,
        order.ResidentID,
        order.VendorID,
        order.CategoryID,
        order.Status,
    ).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

    if err != nil {
        return nil, err
    }
    return order, nil
}
```

### Model Pattern
```go
package models

import "time"

type OrderStatus string

const (
    OrderStatusCreated    OrderStatus = "CREATED"
    OrderStatusPickedUp   OrderStatus = "PICKED_UP"
    OrderStatusProcessing OrderStatus = "PROCESSING"
    OrderStatusReady      OrderStatus = "READY"
    OrderStatusDelivered  OrderStatus = "DELIVERED"
    OrderStatusCompleted  OrderStatus = "COMPLETED"
    OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type Order struct {
    ID          string      `json:"id"`
    OrderNumber string      `json:"order_number"`
    ResidentID  string      `json:"resident_id"`
    VendorID    string      `json:"vendor_id"`
    CategoryID  int         `json:"category_id"`
    Status      OrderStatus `json:"status"`
    TotalPrice  float64     `json:"total_price"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}
```

## Testing Requirements

**Every code change requires:**

### Unit Tests (Required)
Test all service methods with table-driven tests:

```go
func TestOrderService_CreateOrder(t *testing.T) {
    tests := []struct {
        name      string
        req       *CreateOrderRequest
        mockSetup func(*MockOrderRepository)
        want      *models.Order
        wantErr   error
    }{
        {
            name: "creates order successfully",
            req: &CreateOrderRequest{
                ResidentID: "res-123",
                VendorID:   "ven-456",
                CategoryID: 1,
            },
            mockSetup: func(m *MockOrderRepository) {
                m.On("Create", mock.Anything, mock.Anything).Return(&models.Order{
                    ID:         "ord-789",
                    ResidentID: "res-123",
                    Status:     models.OrderStatusCreated,
                }, nil)
            },
            want: &models.Order{
                ID:         "ord-789",
                ResidentID: "res-123",
                Status:     models.OrderStatusCreated,
            },
            wantErr: nil,
        },
        {
            name: "returns error for invalid resident",
            req: &CreateOrderRequest{
                ResidentID: "",
                VendorID:   "ven-456",
            },
            mockSetup: func(m *MockOrderRepository) {},
            want:      nil,
            wantErr:   ErrInvalidRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockOrderRepository)
            tt.mockSetup(mockRepo)

            svc := NewOrderService(mockRepo)
            got, err := svc.CreateOrder(context.Background(), tt.req)

            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.want.ID, got.ID)
            mockRepo.AssertExpectations(t)
        })
    }
}
```

### Handler Tests (Required)
Test HTTP handlers with httptest:

```go
func TestOrderHandler_Create(t *testing.T) {
    gin.SetMode(gin.TestMode)

    tests := []struct {
        name       string
        body       string
        mockSetup  func(*MockOrderService)
        wantStatus int
        wantBody   map[string]interface{}
    }{
        {
            name: "creates order returns 201",
            body: `{"resident_id":"res-123","vendor_id":"ven-456","category_id":1}`,
            mockSetup: func(m *MockOrderService) {
                m.On("CreateOrder", mock.Anything, mock.Anything).Return(&models.Order{
                    ID: "ord-789",
                }, nil)
            },
            wantStatus: http.StatusCreated,
            wantBody: map[string]interface{}{
                "success": true,
            },
        },
        {
            name:       "returns 400 for invalid JSON",
            body:       `{invalid}`,
            mockSetup:  func(m *MockOrderService) {},
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := new(MockOrderService)
            tt.mockSetup(mockSvc)

            handler := NewOrderHandler(mockSvc)

            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Request = httptest.NewRequest("POST", "/orders", strings.NewReader(tt.body))
            c.Request.Header.Set("Content-Type", "application/json")

            handler.Create(c)

            assert.Equal(t, tt.wantStatus, w.Code)
        })
    }
}
```

### Integration Tests (For critical paths)
Test with real database:

```go
func TestOrderRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    db := setupTestDB(t)
    defer db.Close()

    repo := NewOrderRepository(db)

    t.Run("Create and Get", func(t *testing.T) {
        order := &models.Order{
            ResidentID: "res-123",
            VendorID:   "ven-456",
            CategoryID: 1,
            Status:     models.OrderStatusCreated,
        }

        created, err := repo.Create(context.Background(), order)
        require.NoError(t, err)
        require.NotEmpty(t, created.ID)

        fetched, err := repo.GetByID(context.Background(), created.ID)
        require.NoError(t, err)
        assert.Equal(t, created.ID, fetched.ID)
    })
}
```

## Commands

```bash
make build           # Build binary
make run             # Run built binary
make dev             # Hot reload with air
make test            # Run all tests
make test-coverage   # Tests with coverage report
make lint            # Run golangci-lint
make fmt             # Format code
make sqlc            # Generate SQLC code
```

## Response Helpers

Always use the standardized response helpers from `internal/utils/response.go`:

```go
// Success responses
utils.RespondSuccess(c, http.StatusOK, data, "Optional message")
utils.RespondSuccess(c, http.StatusCreated, data, "Created successfully")

// Error responses
utils.RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Details", nil)
utils.RespondError(c, http.StatusNotFound, "NOT_FOUND", "Resource not found", nil)
utils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error", nil)

// Paginated responses
utils.RespondPaginated(c, http.StatusOK, items, pagination)
```

## Error Handling

Define domain errors in services:

```go
var (
    ErrNotFound       = errors.New("resource not found")
    ErrConflict       = errors.New("resource already exists")
    ErrInvalidRequest = errors.New("invalid request")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrForbidden      = errors.New("forbidden")
)
```

Map to HTTP status in handlers - services should not know about HTTP.

## Database Migrations

Location: `backend/migrations/`
Naming: `{timestamp}_{description}.sql`

```sql
-- 20250101120000_create_orders.sql

-- +migrate Up
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(20) UNIQUE NOT NULL,
    resident_id UUID NOT NULL REFERENCES residents(id),
    vendor_id UUID NOT NULL REFERENCES vendors(id),
    category_id INTEGER NOT NULL REFERENCES parent_categories(id),
    status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    total_price DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_orders_resident ON orders(resident_id);
CREATE INDEX idx_orders_vendor ON orders(vendor_id);
CREATE INDEX idx_orders_status ON orders(status);

-- +migrate Down
DROP TABLE IF EXISTS orders;
```
