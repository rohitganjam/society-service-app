# Technology Stack & Architecture

**Version:** 2.0
**Date:** November 17, 2025
**Architecture:** Clean Service-Oriented with Flutter Mobile + Go Backend

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Architecture Overview](#2-architecture-overview)
3. [Technology Stack](#3-technology-stack)
4. [Backend Services (Go on Railway)](#4-backend-services-go-on-railway)
5. [Mobile Applications (Flutter)](#5-mobile-applications-flutter)
6. [Web Applications (Next.js)](#6-web-applications-nextjs)
7. [Database & Storage (Supabase)](#7-database--storage-supabase)
8. [Edge Functions & Webhooks](#8-edge-functions--webhooks)
9. [Third-Party Integrations](#9-third-party-integrations)
10. [Deployment Strategy](#10-deployment-strategy)
11. [Development Workflow](#11-development-workflow)
12. [Infrastructure Costs](#12-infrastructure-costs)
13. [Security & Performance](#13-security--performance)

---

## 1. Executive Summary

### 1.1 Architecture Philosophy

**Clean Service-Oriented Architecture** with strict separation of concerns:

- **Mobile Apps (Flutter)**: Pure UI rendering + API calls, zero business logic
- **Backend API (Go)**: All business logic, validation, orchestration
- **Database (Supabase)**: PostgreSQL with Row Level Security
- **Edge Functions (Supabase)**: Webhooks, cron jobs, background tasks
- **Web Admin (Next.js)**: Dashboard UI + shared backend API

### 1.2 Key Principles

✅ **Separation of Concerns**: Business logic lives ONLY in backend
✅ **Single Source of Truth**: Backend API is the only entry point for data operations
✅ **Thin Clients**: Mobile/Web apps are presentation layers
✅ **API-First Design**: All features exposed via REST/GraphQL APIs
✅ **Serverless**: Auto-scaling, zero server management
✅ **Type Safety**: TypeScript across backend and frontend

### 1.3 Platform Summary

| Component | Technology | Platform | Purpose |
|-----------|-----------|----------|---------|
| **Mobile Apps** | Flutter | iOS/Android | UI rendering only |
| **Backend API** | Go 1.21+ (Gin/Chi) | Railway/VPS | All business logic |
| **Web Admin** | Next.js 14 | Vercel | Admin dashboard |
| **Database** | PostgreSQL | Supabase | Data storage |
| **Edge Functions** | Deno | Supabase | Webhooks, background tasks |
| **File Storage** | S3-compatible | Supabase | Images, documents |
| **Authentication** | Supabase Auth | Supabase | User management |
| **Payments** | Manual (V2: Razorpay) | - | Manual confirmation |
| **Notifications** | Firebase Cloud Messaging | - | Push notifications |
| **Email** | SendGrid/Resend | - | Transactional emails |
| **Monitoring** | Sentry | - | Error tracking |
| **Cron Jobs** | Go cron (robfig/cron) | Railway/VPS | Auto-close orders |

**Total Platforms to Manage:** 5 (Railway/VPS, Vercel, Supabase, Email Provider, Sentry)
**V2 Platforms:** Razorpay (in-app payments)

---

## 2. Architecture Overview

### 2.1 System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENTS                                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │  Resident   │  │   Vendor    │  │   Admin Web Dashboard   │ │
│  │ Flutter App │  │ Flutter App │  │      (Next.js 14)       │ │
│  │             │  │             │  │                         │ │
│  │  • UI Only  │  │  • UI Only  │  │    • UI Only            │ │
│  │  • API Calls│  │  • API Calls│  │    • API Calls          │ │
│  └──────┬──────┘  └──────┬──────┘  └────────────┬────────────┘ │
│         │                │                       │               │
└─────────┼────────────────┼───────────────────────┼───────────────┘
          │                │                       │
          │                │                       │
          └────────────────┴───────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                   RAILWAY/VPS (Go Backend)                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │          Go HTTP API Server (Gin/Chi Framework)           │  │
│  ├───────────────────────────────────────────────────────────┤  │
│  │                                                           │  │
│  │  /api/v1/                                                │  │
│  │  ├─ auth/*           - Authentication                    │  │
│  │  ├─ residents/*      - Resident operations              │  │
│  │  ├─ vendors/*        - Vendor operations                │  │
│  │  ├─ orders/*         - Order management                 │  │
│  │  ├─ payments/*       - Payment processing               │  │
│  │  ├─ categories/*     - Service categories               │  │
│  │  ├─ societies/*      - Society management               │  │
│  │  └─ admin/*          - Admin operations                 │  │
│  │                                                           │  │
│  │  Business Logic Layer:                                   │  │
│  │  ├─ Order Validation                                     │  │
│  │  ├─ Pricing Calculation                                  │  │
│  │  ├─ Workflow Orchestration                              │  │
│  │  ├─ Payment Processing                                   │  │
│  │  ├─ Notification Triggers                               │  │
│  │  └─ Analytics & Reporting                               │  │
│  │                                                           │  │
│  └───────────────────┬───────────────────────────────────────┘  │
│                      │                                           │
└──────────────────────┼───────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    SUPABASE (Backend Services)                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────┐  │
│  │   PostgreSQL     │  │  Edge Functions  │  │   Storage    │  │
│  │   Database       │  │     (Deno)       │  │  (S3-like)   │  │
│  ├──────────────────┤  ├──────────────────┤  ├──────────────┤  │
│  │                  │  │                  │  │              │  │
│  │ • All tables     │  │ • Webhooks       │  │ • Images     │  │
│  │ • RLS policies   │  │ • Cron jobs      │  │ • Documents  │  │
│  │ • Triggers       │  │ • Async tasks    │  │ • Photos     │  │
│  │ • Functions      │  │ • Integrations   │  │              │  │
│  │                  │  │                  │  │              │  │
│  └──────────────────┘  └──────────────────┘  └──────────────┘  │
│                                                                  │
│  ┌──────────────────┐  ┌──────────────────┐                    │
│  │  Supabase Auth   │  │    Realtime      │                    │
│  ├──────────────────┤  ├──────────────────┤                    │
│  │                  │  │                  │                    │
│  │ • JWT tokens     │  │ • Live updates   │                    │
│  │ • User sessions  │  │ • Subscriptions  │                    │
│  │ • OTP/Phone auth │  │ • Pub/Sub        │                    │
│  │                  │  │                  │                    │
│  └──────────────────┘  └──────────────────┘                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                    EXTERNAL SERVICES                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │  Razorpay    │  │     FCM      │  │    Twilio/MSG91      │  │
│  │  Payments    │  │Push Notifs   │  │        SMS           │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐                            │
│  │  SendGrid/   │  │    Sentry    │                            │
│  │  Resend      │  │  Monitoring  │                            │
│  │  Email       │  │              │                            │
│  └──────────────┘  └──────────────┘                            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Request Flow Example

**Resident creates an order:**

```
1. Resident App (Flutter)
   └─> Calls: POST /api/v1/orders
       Body: { laundry_id, items[], pickup_time }

2. Backend API (Node.js on Vercel)
   ├─> Validates request (auth, data)
   ├─> Calculates pricing (business logic)
   ├─> Checks vendor availability
   ├─> Calculates delivery estimates
   ├─> Creates order in database
   ├─> Triggers notification (via edge function)
   └─> Returns: { order_id, total, estimated_delivery }

3. Supabase PostgreSQL
   ├─> Inserts into orders table
   ├─> Inserts into order_items table
   ├─> Inserts into order_service_status table
   └─> Returns success

4. Supabase Edge Function (async)
   ├─> Triggered by new order
   ├─> Sends push notification to vendor
   ├─> Sends SMS to vendor (if enabled)
   └─> Logs notification

5. Resident App
   └─> Displays order confirmation
```

**Key Point:** Flutter app NEVER talks directly to database. All operations go through Backend API.

---

## 3. Technology Stack

### 3.1 Mobile Applications

**Technology:** Flutter 3.x

**Why Flutter:**
- ✅ Single codebase for iOS + Android
- ✅ Native performance (compiled to native ARM code)
- ✅ Excellent UI framework with Material Design & Cupertino
- ✅ Hot reload for fast development
- ✅ Growing ecosystem with strong Indian market support
- ✅ Better performance than React Native
- ✅ Dart language - type-safe, modern, easy to learn

**Flutter Packages:**
```yaml
dependencies:
  flutter_riverpod: ^2.4.0        # State management
  go_router: ^12.0.0              # Navigation
  dio: ^5.4.0                     # HTTP client for API calls
  freezed: ^2.4.0                 # Immutable models
  json_annotation: ^4.8.0         # JSON serialization
  flutter_secure_storage: ^9.0.0  # Secure token storage
  razorpay_flutter: ^1.3.0        # Razorpay integration
  firebase_messaging: ^14.7.0     # Push notifications
  image_picker: ^1.0.0            # Camera/gallery
  cached_network_image: ^3.3.0    # Image caching
  intl: ^0.18.0                   # Internationalization
```

**Architecture Pattern:** Clean Architecture
```
lib/
├── core/
│   ├── api/
│   │   └── api_client.dart      # Dio HTTP client
│   ├── models/
│   │   ├── order.dart
│   │   ├── vendor.dart
│   │   └── ...
│   └── providers/
│       └── auth_provider.dart
│
├── features/
│   ├── orders/
│   │   ├── data/
│   │   │   └── orders_repository.dart  # API calls only
│   │   ├── domain/
│   │   │   └── order_model.dart
│   │   └── presentation/
│   │       ├── screens/
│   │       └── widgets/
│   │
│   └── vendors/
│       └── ...
│
└── main.dart
```

**No Business Logic in Flutter:**
- ❌ No price calculations
- ❌ No validation beyond basic input
- ❌ No workflow logic
- ✅ Only UI rendering
- ✅ Only API calls to backend
- ✅ Only local state management

---

### 3.2 Backend API Server

**Technology:** Go 1.21+

**Why Go:**
- ✅ High performance and low memory footprint
- ✅ Built-in concurrency (goroutines) perfect for high-traffic APIs
- ✅ Static typing and compile-time error checking
- ✅ Fast compilation and deployment
- ✅ Excellent standard library (net/http, crypto, json)
- ✅ Easy deployment as single binary
- ✅ Strong ecosystem for web services

**Framework:** Gin or Chi (recommended: Gin for performance)

**Core Dependencies:**
```go
// go.mod
module society-service-api

go 1.21

require (
    // Web Framework
    github.com/gin-gonic/gin v1.9.1

    // Database & Supabase
    github.com/jackc/pgx/v5 v5.5.0
    github.com/supabase-community/gotrue-go v1.0.0
    github.com/supabase-community/storage-go v0.7.0

    // Validation
    github.com/go-playground/validator/v10 v10.16.0

    // Environment & Config
    github.com/joho/godotenv v1.5.1
    github.com/spf13/viper v1.18.1

    // Cron Jobs
    github.com/robfig/cron/v3 v3.0.1

    // HTTP Client
    github.com/go-resty/resty/v2 v2.11.0

    // Middleware
    github.com/gin-contrib/cors v1.5.0
    github.com/gin-contrib/gzip v0.0.6

    // Logging
    github.com/sirupsen/logrus v1.9.3
    go.uber.org/zap v1.26.0

    // Error Tracking
    github.com/getsentry/sentry-go v0.25.0

    // Email
    github.com/resend/resend-go/v2 v2.0.0

    // Payments (V2)
    github.com/razorpay/razorpay-go v1.2.0

    // Testing
    github.com/stretchr/testify v1.8.4
)
```

**Project Structure:**
```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── handlers/                   # HTTP handlers (controllers)
│   │   ├── auth/
│   │   │   ├── login.go
│   │   │   ├── verify_otp.go
│   │   │   └── refresh_token.go
│   │   │
│   │   ├── orders/
│   │   │   ├── create.go           # POST /api/v1/orders
│   │   │   ├── get.go              # GET /api/v1/orders/:id
│   │   │   ├── list.go             # GET /api/v1/orders
│   │   │   ├── update_status.go    # PATCH /api/v1/orders/:id/status
│   │   │   └── cancel.go           # POST /api/v1/orders/:id/cancel
│   │   │
│   │   ├── vendors/
│   │   │   ├── register.go
│   │   │   ├── rate_cards.go
│   │   │   ├── services.go
│   │   │   └── analytics.go
│   │   │
│   │   ├── categories/
│   │   │   ├── list.go
│   │   │   └── services.go
│   │   │
│   │   └── admin/
│   │       ├── societies.go
│   │       ├── subscriptions.go
│   │       └── reports.go
│   │
│   ├── middleware/
│   │   ├── auth.go                 # JWT verification
│   │   ├── validate.go             # Request validation
│   │   ├── error_handler.go        # Global error handling
│   │   ├── rate_limit.go           # Rate limiting
│   │   └── cors.go                 # CORS configuration
│   │
│   ├── services/
│   │   ├── order_service.go        # Order business logic
│   │   ├── pricing_service.go      # Price calculation
│   │   ├── workflow_service.go     # Service workflows
│   │   ├── payment_service.go      # Payment handling
│   │   └── notification_service.go # Send notifications
│   │
│   ├── repositories/
│   │   ├── order_repository.go     # Database operations
│   │   ├── vendor_repository.go
│   │   └── user_repository.go
│   │
│   ├── models/
│   │   ├── order.go
│   │   ├── vendor.go
│   │   ├── user.go
│   │   └── response.go             # Standard API responses
│   │
│   ├── database/
│   │   ├── postgres.go             # PostgreSQL connection
│   │   └── migrations/             # SQL migration files
│   │
│   ├── config/
│   │   └── config.go               # Configuration management
│   │
│   └── utils/
│       ├── logger.go               # Structured logging
│       ├── validator.go            # Input validation
│       └── helpers.go
│
├── pkg/                            # Public packages
│   ├── supabase/
│   │   ├── client.go               # Supabase client wrapper
│   │   └── auth.go                 # Supabase auth helpers
│   │
│   └── errors/
│       └── errors.go               # Custom error types
│
├── scripts/
│   └── migrate.sh                  # Database migration script
│
├── .env.example
├── Dockerfile
├── docker-compose.yml              # For local development
├── go.mod
├── go.sum
└── Makefile                        # Build and deployment commands
```

**Example API Endpoint:**
```go
// internal/handlers/orders/create.go
package orders

import (
    "net/http"
    "society-service-api/internal/models"
    "society-service-api/internal/services"
    "society-service-api/internal/middleware"

    "github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
    VendorID           string      `json:"vendor_id" binding:"required,uuid"`
    SocietyID          int         `json:"society_id" binding:"required,min=1"`
    Items              []OrderItem `json:"items" binding:"required,min=1,dive"`
    PickupDatetime     string      `json:"pickup_datetime" binding:"required"`
    PickupAddress      string      `json:"pickup_address" binding:"required"`
    DeliveryPreference string      `json:"delivery_preference" binding:"required,oneof=SINGLE PARTIAL"`
}

type OrderItem struct {
    ServiceID  int     `json:"service_id" binding:"required,min=1"`
    ItemName   string  `json:"item_name" binding:"required"`
    Quantity   int     `json:"quantity" binding:"required,min=1"`
    UnitPrice  float64 `json:"unit_price" binding:"required,min=0"`
}

func CreateOrder(orderService *services.OrderService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user ID from auth middleware
        userID := middleware.GetUserID(c)

        var req CreateOrderRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, models.ErrorResponse{
                Success: false,
                Error: models.Error{
                    Code:    "INVALID_REQUEST",
                    Message: "Invalid request data",
                    Details: map[string]interface{}{"error": err.Error()},
                },
            })
            return
        }

        // Business logic in service layer
        order, err := orderService.CreateOrder(c.Request.Context(), userID, req)
        if err != nil {
            c.JSON(http.StatusBadRequest, models.ErrorResponse{
                Success: false,
                Error: models.Error{
                    Code:    "ORDER_CREATION_FAILED",
                    Message: err.Error(),
                },
            })
            return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
            Success: true,
            Data:    order,
        })
    }
}
```

**Service Layer Example:**
```go
// internal/services/order_service.go
package services

import (
    "context"
    "errors"
    "society-service-api/internal/models"
    "society-service-api/internal/repositories"
)

type OrderService struct {
    orderRepo        *repositories.OrderRepository
    pricingService   *PricingService
    workflowService  *WorkflowService
    notificationSvc  *NotificationService
}

func NewOrderService(
    orderRepo *repositories.OrderRepository,
    pricingService *PricingService,
    workflowService *WorkflowService,
    notificationSvc *NotificationService,
) *OrderService {
    return &OrderService{
        orderRepo:       orderRepo,
        pricingService:  pricingService,
        workflowService: workflowService,
        notificationSvc: notificationSvc,
    }
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, req CreateOrderRequest) (*models.Order, error) {
    // 1. Validate vendor availability
    vendor, err := s.orderRepo.GetVendorAvailability(ctx, req.VendorID, req.PickupDatetime)
    if err != nil {
        return nil, err
    }

    if !vendor.IsAvailable {
        return nil, errors.New("vendor not available at requested time")
    }

    // 2. Calculate pricing (business logic)
    const pricing = await PricingService.calculateOrderTotal(orderData.items);

    // 3. Calculate delivery estimate
    const deliveryEstimate = await WorkflowService.calculateDeliveryDate(
      orderData.items
    );

    // 4. Create order in database
    const order = await OrderRepository.createOrder({
      resident_id: userId,
      laundry_id: orderData.laundry_id,
      society_id: orderData.society_id,
      estimated_price: pricing.total,
      expected_delivery_date: deliveryEstimate,
      pickup_datetime: orderData.pickup_datetime,
      pickup_address: orderData.pickup_address
    });

    // 5. Create order items
    await OrderRepository.createOrderItems(order.order_id, orderData.items);

    // 6. Create service status tracking
    await WorkflowService.initializeServiceTracking(order.order_id, orderData.items);

    // 7. Trigger notification (async)
    NotificationService.sendOrderNotification(order.order_id, 'vendor');

    return {
      order_id: order.order_id,
      order_number: order.order_number,
      total: pricing.total,
      estimated_delivery: deliveryEstimate,
      service_breakdown: pricing.breakdown
    };
  }
}
```

**All Business Logic Lives Here:**
- ✅ Order validation and creation
- ✅ Pricing calculation
- ✅ Delivery time estimation
- ✅ Workflow state management
- ✅ Payment processing
- ✅ Notification triggers
- ✅ Analytics calculation
- ✅ Report generation

---

### 3.3 Web Admin Dashboard

**Technology:** Next.js 14 (App Router)

**Why Next.js:**
- ✅ React-based, easy to develop
- ✅ Server-side rendering for better performance
- ✅ Perfect Vercel integration (zero-config deployment)
- ✅ Built-in API routes (for admin-specific operations)
- ✅ File-based routing
- ✅ Great developer experience

**Key Dependencies:**
```json
{
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "typescript": "^5.3.0",

    "axios": "^1.6.0",
    "@tanstack/react-query": "^5.8.0",

    "shadcn/ui": "latest",
    "tailwindcss": "^3.3.0",
    "recharts": "^2.10.0",

    "react-hook-form": "^7.48.0",
    "zod": "^3.22.4",
    "date-fns": "^2.30.0"
  }
}
```

**Project Structure:**
```
admin-web/
├── app/
│   ├── (auth)/
│   │   └── login/
│   │       └── page.tsx
│   │
│   ├── (dashboard)/
│   │   ├── page.tsx                    # Dashboard home
│   │   ├── societies/
│   │   │   ├── page.tsx                # List societies
│   │   │   ├── [id]/page.tsx           # Society details
│   │   │   └── new/page.tsx            # Add society
│   │   │
│   │   ├── vendors/
│   │   │   ├── page.tsx                # Pending approvals
│   │   │   └── [id]/page.tsx
│   │   │
│   │   ├── orders/
│   │   │   └── page.tsx                # Order monitoring
│   │   │
│   │   ├── subscriptions/
│   │   │   ├── page.tsx
│   │   │   └── invoices/page.tsx
│   │   │
│   │   ├── categories/
│   │   │   └── page.tsx                # Manage categories
│   │   │
│   │   └── analytics/
│   │       └── page.tsx
│   │
│   └── layout.tsx
│
├── components/
│   ├── ui/                             # shadcn components
│   ├── societies/
│   ├── vendors/
│   └── orders/
│
├── lib/
│   ├── api-client.ts                   # Axios wrapper (calls backend API)
│   └── utils.ts
│
└── hooks/
    ├── use-societies.ts
    ├── use-vendors.ts
    └── use-orders.ts
```

**Admin Dashboard calls Backend API:**
```typescript
// lib/api-client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL, // Points to Vercel backend
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add auth token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default apiClient;

// Usage in components:
export const getSocieties = async () => {
  const { data } = await apiClient.get('/api/v1/admin/societies');
  return data;
};
```

**No Direct Database Access:** Admin dashboard also goes through Backend API for all operations.

---

## 4. Backend Services (Go on Railway)

### 4.1 Railway/Server Configuration

**Dockerfile:**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/api ./

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api"]
```

**Environment Variables:**
```bash
# .env.example
PORT=8080
ENV=production

# Database
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key
DATABASE_URL=postgresql://postgres:[password]@db.your-project.supabase.co:5432/postgres

# Authentication
JWT_SECRET=your-jwt-secret
JWT_EXPIRY=24h

# External Services
RAZORPAY_KEY_ID=your-razorpay-key-id
RAZORPAY_KEY_SECRET=your-razorpay-key-secret
RESEND_API_KEY=your-resend-api-key
SENTRY_DSN=your-sentry-dsn

# Logging
LOG_LEVEL=info
```

### 4.2 API Design Patterns

**RESTful API Structure:**

```
Base URL: https://api.yourapp.com/api/v1

Authentication:
POST   /auth/login                 # Send OTP
POST   /auth/verify-otp            # Verify OTP, get JWT
POST   /auth/refresh               # Refresh JWT token
POST   /auth/logout                # Invalidate session

Categories:
GET    /categories                 # List all categories
GET    /categories/:id/services    # List services in category

Residents:
GET    /residents/me               # Get current user
PATCH  /residents/me               # Update profile
GET    /residents/me/orders        # Order history

Vendors:
POST   /vendors/register           # Vendor registration
GET    /vendors/:id                # Get vendor details
GET    /vendors/:id/services       # Services offered
GET    /vendors/:id/rate-card      # Get rate card
POST   /vendors/:id/rate-card      # Create/update rate card
GET    /vendors/search             # Search vendors by society/category
GET    /vendors/me/dashboard       # Vendor dashboard data
GET    /vendors/me/analytics       # Vendor analytics

Orders:
POST   /orders                     # Create order
GET    /orders/:id                 # Get order details
GET    /orders                     # List orders (filtered)
PATCH  /orders/:id/status          # Update order status
PATCH  /orders/:id/service-status  # Update service type status
POST   /orders/:id/cancel          # Cancel order
POST   /orders/:id/approve-count   # Approve count change

Payments:
POST   /payments                   # Record payment
GET    /payments/:id               # Payment details
GET    /payments                   # Payment history

Societies:
GET    /societies                  # List societies
GET    /societies/:id              # Society details
POST   /societies                  # Create society (admin)
PATCH  /societies/:id              # Update society (admin)

Subscriptions:
GET    /subscriptions              # List subscriptions (admin)
GET    /subscriptions/:id          # Subscription details
POST   /subscriptions/:id/invoice  # Generate invoice
PATCH  /subscriptions/:id/status   # Update status (admin)

Admin:
GET    /admin/dashboard            # Platform-wide stats
GET    /admin/vendors/pending      # Pending approvals
POST   /admin/vendors/:id/approve  # Approve vendor
POST   /admin/vendors/:id/reject   # Reject vendor
GET    /admin/orders               # All orders monitoring
GET    /admin/disputes             # All disputes
```

**Standard Response Format:**
```typescript
// Success
{
  "success": true,
  "data": { ... },
  "meta": {
    "timestamp": "2025-11-17T10:30:00Z",
    "request_id": "uuid"
  }
}

// Error
{
  "success": false,
  "error": {
    "code": "INVALID_INPUT",
    "message": "Pickup time must be in the future",
    "details": { ... }
  },
  "meta": {
    "timestamp": "2025-11-17T10:30:00Z",
    "request_id": "uuid"
  }
}

// Paginated
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### 4.3 Authentication Flow

**JWT-based with Supabase Auth:**

```typescript
// Backend: api/v1/auth/login.ts
export default async (req: Request, res: Response) => {
  const { phone } = req.body;

  // Send OTP via Supabase Auth
  const { data, error } = await supabase.auth.signInWithOtp({
    phone: phone
  });

  if (error) {
    return res.status(400).json({
      success: false,
      error: { message: 'Failed to send OTP' }
    });
  }

  res.json({
    success: true,
    data: { message: 'OTP sent successfully' }
  });
};

// Backend: api/v1/auth/verify-otp.ts
export default async (req: Request, res: Response) => {
  const { phone, otp } = req.body;

  // Verify OTP with Supabase
  const { data, error } = await supabase.auth.verifyOtp({
    phone: phone,
    token: otp,
    type: 'sms'
  });

  if (error) {
    return res.status(401).json({
      success: false,
      error: { message: 'Invalid OTP' }
    });
  }

  // Return JWT token
  res.json({
    success: true,
    data: {
      access_token: data.session.access_token,
      refresh_token: data.session.refresh_token,
      user: data.user
    }
  });
};

// Mobile App: Uses token for all API calls
dio.options.headers['Authorization'] = 'Bearer $token';
```

### 4.4 Business Logic Services

**Order Service - Complete Example:**

```typescript
// src/services/order-service.ts
import { supabase } from '@/utils/supabase';
import { PricingService } from './pricing-service';
import { WorkflowService } from './workflow-service';

export class OrderService {
  /**
   * Create new order with full validation and processing
   */
  static async createOrder(residentId: string, orderData: CreateOrderDTO) {
    // 1. Validate vendor exists and is active
    const vendor = await this.validateVendor(orderData.laundry_id);

    // 2. Validate services are offered by vendor
    await this.validateVendorServices(
      orderData.laundry_id,
      orderData.items.map(i => i.service_id)
    );

    // 3. Validate society access
    await this.validateSocietyAccess(residentId, orderData.society_id);

    // 4. Calculate pricing
    const pricing = await PricingService.calculateOrderTotal(orderData.items);

    // 5. Calculate delivery estimate
    const serviceIds = [...new Set(orderData.items.map(i => i.service_id))];
    const deliveryEstimate = await WorkflowService.calculateDeliveryDate(serviceIds);

    // 6. Generate order number
    const orderNumber = await this.generateOrderNumber();

    // 7. Create order (transaction)
    const { data: order, error } = await supabase
      .from('orders')
      .insert({
        order_number: orderNumber,
        resident_id: residentId,
        laundry_id: orderData.laundry_id,
        society_id: orderData.society_id,
        status: 'PICKUP_SCHEDULED',
        estimated_price: pricing.total,
        expected_delivery_date: deliveryEstimate,
        pickup_datetime: orderData.pickup_datetime,
        pickup_address: orderData.pickup_address,
        has_multiple_services: serviceIds.length > 1
      })
      .select()
      .single();

    if (error) throw new Error('Failed to create order');

    // 8. Create order items
    const itemsToInsert = orderData.items.map(item => ({
      order_id: order.order_id,
      service_id: item.service_id,
      item_name: item.item_name,
      quantity: item.quantity,
      unit_price: item.unit_price,
      total_price: item.quantity * item.unit_price
    }));

    await supabase.from('order_items').insert(itemsToInsert);

    // 9. Initialize service tracking
    const serviceGroups = this.groupItemsByService(orderData.items);
    const statusRecords = serviceGroups.map(group => ({
      order_id: order.order_id,
      service_id: group.service_id,
      item_count: group.total_items,
      total_amount: group.total_amount,
      status: 'PICKUP_SCHEDULED',
      expected_delivery_date: group.expected_delivery
    }));

    await supabase.from('order_service_status').insert(statusRecords);

    // 10. Return formatted response
    return {
      order_id: order.order_id,
      order_number: order.order_number,
      status: order.status,
      total: pricing.total,
      expected_delivery: deliveryEstimate,
      service_breakdown: pricing.breakdown
    };
  }

  /**
   * Update order service status
   */
  static async updateServiceStatus(
    orderId: string,
    serviceId: number,
    status: string,
    updatedBy: string
  ) {
    // Validate status transition
    await WorkflowService.validateStatusTransition(orderId, serviceId, status);

    // Update service status
    const { error } = await supabase
      .from('order_service_status')
      .update({
        status: status,
        [`${this.getStatusField(status)}_at`]: new Date().toISOString()
      })
      .eq('order_id', orderId)
      .eq('service_id', serviceId);

    if (error) throw new Error('Failed to update status');

    // Check if all services in order are complete
    const allComplete = await this.checkOrderCompletion(orderId);

    if (allComplete) {
      await this.completeOrder(orderId);
    }

    return { success: true };
  }

  // ... more methods
}
```

**Pricing Service:**

```typescript
// src/services/pricing-service.ts
export class PricingService {
  static async calculateOrderTotal(items: OrderItemDTO[]) {
    let total = 0;
    const breakdown: any[] = [];

    // Group by service type
    const grouped = items.reduce((acc, item) => {
      if (!acc[item.service_id]) acc[item.service_id] = [];
      acc[item.service_id].push(item);
      return acc;
    }, {});

    for (const [serviceId, serviceItems] of Object.entries(grouped)) {
      const serviceTotal = serviceItems.reduce(
        (sum, item) => sum + (item.quantity * item.unit_price),
        0
      );

      const service = await this.getServiceDetails(parseInt(serviceId));

      breakdown.push({
        service_id: parseInt(serviceId),
        service_name: service.service_name,
        item_count: serviceItems.reduce((sum, i) => sum + i.quantity, 0),
        subtotal: serviceTotal
      });

      total += serviceTotal;
    }

    return {
      total,
      breakdown
    };
  }

  // Apply discounts, coupons, etc. (future)
  static async applyDiscounts(total: number, discountCode?: string) {
    // Business logic for discounts
    return total;
  }
}
```

---

## 5. Mobile Applications (Flutter)

### 5.1 Flutter Architecture

**Clean Architecture Pattern:**

```dart
// lib/core/api/api_client.dart
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiClient {
  final Dio _dio;
  final FlutterSecureStorage _storage;

  static const String baseUrl = 'https://api.yourapp.com/api/v1';

  ApiClient()
      : _dio = Dio(BaseOptions(
          baseUrl: baseUrl,
          connectTimeout: Duration(seconds: 30),
          receiveTimeout: Duration(seconds: 30),
        )),
        _storage = FlutterSecureStorage() {
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        // Add auth token
        final token = await _storage.read(key: 'auth_token');
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        return handler.next(options);
      },
      onError: (error, handler) async {
        // Handle 401 unauthorized
        if (error.response?.statusCode == 401) {
          // Refresh token or logout
        }
        return handler.next(error);
      },
    ));
  }

  // Generic HTTP methods
  Future<T> get<T>(String path, {Map<String, dynamic>? params}) async {
    final response = await _dio.get(path, queryParameters: params);
    return response.data['data'] as T;
  }

  Future<T> post<T>(String path, {Map<String, dynamic>? data}) async {
    final response = await _dio.post(path, data: data);
    return response.data['data'] as T;
  }

  // ... patch, delete methods
}
```

**Repository Pattern (API Calls Only):**

```dart
// lib/features/orders/data/orders_repository.dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../core/api/api_client.dart';
import '../domain/order_model.dart';

part 'orders_repository.g.dart';

@riverpod
class OrdersRepository extends _$OrdersRepository {
  @override
  FutureOr<void> build() {}

  // Create order - calls backend API
  Future<Order> createOrder(CreateOrderRequest request) async {
    final apiClient = ref.read(apiClientProvider);

    final response = await apiClient.post<Map<String, dynamic>>(
      '/orders',
      data: request.toJson(),
    );

    return Order.fromJson(response);
  }

  // Get order details
  Future<Order> getOrder(String orderId) async {
    final apiClient = ref.read(apiClientProvider);

    final response = await apiClient.get<Map<String, dynamic>>(
      '/orders/$orderId',
    );

    return Order.fromJson(response);
  }

  // List orders
  Future<List<Order>> listOrders({
    int page = 1,
    int limit = 20,
    String? status,
  }) async {
    final apiClient = ref.read(apiClientProvider);

    final response = await apiClient.get<List<dynamic>>(
      '/orders',
      params: {
        'page': page,
        'limit': limit,
        if (status != null) 'status': status,
      },
    );

    return response.map((json) => Order.fromJson(json)).toList();
  }

  // Update order status
  Future<void> updateOrderStatus(String orderId, String status) async {
    final apiClient = ref.read(apiClientProvider);

    await apiClient.patch(
      '/orders/$orderId/status',
      data: {'status': status},
    );
  }
}
```

**UI Layer (Pure Presentation):**

```dart
// lib/features/orders/presentation/screens/create_order_screen.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../data/orders_repository.dart';
import '../../domain/order_model.dart';

class CreateOrderScreen extends ConsumerStatefulWidget {
  final String vendorId;

  const CreateOrderScreen({required this.vendorId});

  @override
  ConsumerState<CreateOrderScreen> createState() => _CreateOrderScreenState();
}

class _CreateOrderScreenState extends ConsumerState<CreateOrderScreen> {
  final List<OrderItem> _cart = [];
  DateTime? _pickupTime;

  Future<void> _createOrder() async {
    if (_cart.isEmpty) {
      // Show error
      return;
    }

    final request = CreateOrderRequest(
      laundryId: widget.vendorId,
      items: _cart,
      pickupDatetime: _pickupTime!.toIso8601String(),
      pickupAddress: 'A-404', // From user profile
    );

    try {
      // Call backend API via repository
      final order = await ref.read(ordersRepositoryProvider.notifier)
          .createOrder(request);

      // Navigate to order details
      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (_) => OrderDetailsScreen(orderId: order.orderId),
        ),
      );
    } catch (e) {
      // Show error
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to create order: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Create Order')),
      body: Column(
        children: [
          // Cart items list
          Expanded(
            child: ListView.builder(
              itemCount: _cart.length,
              itemBuilder: (context, index) {
                final item = _cart[index];
                return ListTile(
                  title: Text(item.itemName),
                  subtitle: Text('${item.quantity} x ₹${item.unitPrice}'),
                  trailing: Text('₹${item.quantity * item.unitPrice}'),
                );
              },
            ),
          ),

          // Total
          Padding(
            padding: EdgeInsets.all(16),
            child: Text(
              'Total: ₹${_calculateTotal()}',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
            ),
          ),

          // Create order button
          ElevatedButton(
            onPressed: _createOrder,
            child: Text('Create Order'),
          ),
        ],
      ),
    );
  }

  double _calculateTotal() {
    // Simple UI calculation (backend will validate)
    return _cart.fold(0, (sum, item) => sum + (item.quantity * item.unitPrice));
  }
}
```

**Key Points:**
- ✅ Flutter app ONLY renders UI and calls APIs
- ✅ All validation happens in backend
- ✅ All calculations happen in backend
- ✅ Flutter just displays data from API responses
- ✅ No business logic in Dart code

---

## 6. Web Applications (Next.js)

### 6.1 Admin Dashboard Architecture

**API Calls via React Query:**

```typescript
// lib/api/societies.ts
import apiClient from './api-client';

export const getSocieties = async (params?: {
  page?: number;
  limit?: number;
  search?: string;
}) => {
  const { data } = await apiClient.get('/api/v1/admin/societies', { params });
  return data;
};

export const approveSociety = async (societyId: number) => {
  const { data } = await apiClient.post(
    `/api/v1/admin/societies/${societyId}/approve`
  );
  return data;
};
```

**React Component:**

```typescript
// app/(dashboard)/societies/page.tsx
'use client';

import { useQuery, useMutation } from '@tanstack/react-query';
import { getSocieties, approveSociety } from '@/lib/api/societies';
import { Button } from '@/components/ui/button';

export default function SocietiesPage() {
  const { data: societies, isLoading } = useQuery({
    queryKey: ['societies'],
    queryFn: () => getSocieties()
  });

  const approveMutation = useMutation({
    mutationFn: approveSociety,
    onSuccess: () => {
      // Refetch societies
      queryClient.invalidateQueries(['societies']);
    }
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Societies</h1>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {societies?.data.map((society) => (
            <tr key={society.id}>
              <td>{society.name}</td>
              <td>{society.status}</td>
              <td>
                {society.status === 'PENDING' && (
                  <Button
                    onClick={() => approveMutation.mutate(society.id)}
                  >
                    Approve
                  </Button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

**No Direct Database Access:** Next.js app calls Backend API for everything.

---

## 7. Database & Storage (Supabase)

### 7.1 Supabase Services Used

**PostgreSQL Database:**
- Primary data store (PostgreSQL 15+)
- All tables, relationships, constraints
- Row Level Security (RLS) for data isolation
- Database functions and triggers
- Full-text search capabilities
- **ltree extension** for hierarchical tree structures (society hierarchy model)

**Supabase Auth:**
- Phone-based OTP authentication
- JWT token generation
- Session management
- User metadata storage

**Supabase Storage:**
- S3-compatible file storage
- Image uploads (vendor photos, pickup photos)
- Document storage (invoices, reports)
- CDN for fast delivery

**Supabase Realtime:**
- Live order updates
- Real-time notifications
- Pub/Sub for events

**Access Pattern:**
- Backend API uses Supabase service key (full access)
- Mobile apps DO NOT access Supabase directly
- All database operations through Backend API

### 7.2 Database Connection

**Backend uses Supabase JS Client:**

```typescript
// src/utils/supabase.ts
import { createClient } from '@supabase/supabase-js';

// Service key for backend (full access)
export const supabase = createClient(
  process.env.SUPABASE_URL!,
  process.env.SUPABASE_SERVICE_KEY!, // NOT anon key
  {
    auth: {
      autoRefreshToken: false,
      persistSession: false
    }
  }
);

// Usage in repositories
export class OrderRepository {
  static async createOrder(orderData: any) {
    const { data, error } = await supabase
      .from('orders')
      .insert(orderData)
      .select()
      .single();

    if (error) throw error;
    return data;
  }
}
```

**PostgreSQL Extensions Used:**

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";      -- UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";       -- Encryption
CREATE EXTENSION IF NOT EXISTS "pg_trgm";        -- Full-text search
CREATE EXTENSION IF NOT EXISTS "btree_gin";      -- GIN indexes
CREATE EXTENSION IF NOT EXISTS "ltree";          -- Hierarchical trees
```

**ltree Extension for Generic Hierarchy:**

The `ltree` extension is critical for the flexible society hierarchy model:

- **Purpose:** Efficiently store and query tree structures (society → building → floor → unit)
- **Path Format:** Materialized paths like `1.2.4.6` (society.building.floor.unit)
- **Operators:**
  - `<@` : Path is ancestor (e.g., `'1.2.4.6' <@ '1.2'` checks if unit is in Building A)
  - `@>` : Path is descendant
  - `~` : Pattern matching
- **Indexes:** GIST indexes on ltree columns enable fast tree traversal
- **Benefits:**
  - O(log n) ancestor/descendant queries
  - No recursive CTEs needed
  - Automatic path validation

**Example Query:**
```sql
-- Find all vendors assigned to resident's hierarchy path
SELECT v.*
FROM vendors v
JOIN vendor_service_areas vsa ON v.vendor_id = vsa.vendor_id
JOIN hierarchy_nodes hn ON vsa.node_id = hn.node_id
WHERE resident_unit_path <@ hn.path;  -- Fast ltree comparison
```

**RLS Policies (Defense in Depth):**

Even though backend uses service key, RLS provides additional security:

```sql
-- Residents can only see their own orders
CREATE POLICY "Residents view own orders"
ON orders FOR SELECT
USING (auth.uid() = resident_id);

-- Vendors can only see orders assigned to them
CREATE POLICY "Vendors view assigned orders"
ON orders FOR SELECT
USING (auth.uid() = laundry_id);

-- Admins can see all orders (via service key bypass)
```

### 7.3 File Upload Flow

```
1. Mobile App requests upload URL
   └─> POST /api/v1/uploads/request
       Backend generates signed URL from Supabase Storage

2. Mobile App uploads file directly to Supabase Storage
   └─> PUT https://supabase.co/storage/v1/object/...
       Uses signed URL (no auth needed)

3. Mobile App sends file URL to backend
   └─> POST /api/v1/orders/:id/photos
       Backend saves URL reference in database
```

---

## 8. Edge Functions & Webhooks

### 8.1 Supabase Edge Functions

**Purpose:**
- Async background tasks
- Webhook handling
- Scheduled jobs (cron)
- Third-party integrations

**Technology:** Deno (TypeScript)

**Example: Send Push Notification**

```typescript
// supabase/functions/send-notification/index.ts
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';

serve(async (req) => {
  try {
    const { user_id, title, body, data } = await req.json();

    // Get user's FCM token from database
    const supabase = createClient(
      Deno.env.get('SUPABASE_URL')!,
      Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
    );

    const { data: user } = await supabase
      .from('users')
      .select('fcm_token')
      .eq('user_id', user_id)
      .single();

    if (!user?.fcm_token) {
      return new Response('No FCM token', { status: 404 });
    }

    // Send via Firebase Cloud Messaging
    const fcmResponse = await fetch(
      'https://fcm.googleapis.com/fcm/send',
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `key=${Deno.env.get('FCM_SERVER_KEY')}`
        },
        body: JSON.stringify({
          to: user.fcm_token,
          notification: { title, body },
          data: data
        })
      }
    );

    return new Response(JSON.stringify({ success: true }), {
      headers: { 'Content-Type': 'application/json' }
    });
  } catch (error) {
    return new Response(JSON.stringify({ error: error.message }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});
```

**Triggered from Backend:**

```typescript
// Backend: src/services/notification-service.ts
export class NotificationService {
  static async sendOrderNotification(orderId: string, recipientType: string) {
    // Call Supabase Edge Function
    const { data, error } = await supabase.functions.invoke(
      'send-notification',
      {
        body: {
          user_id: userId,
          title: 'New Order',
          body: `You have a new order #${orderNumber}`,
          data: { order_id: orderId }
        }
      }
    );

    if (error) console.error('Notification failed:', error);
  }
}
```

### 8.2 Cron Jobs (Scheduled Tasks)

**Example: Monthly Invoice Generation**

```typescript
// supabase/functions/generate-monthly-invoices/index.ts
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';

serve(async (req) => {
  const supabase = createClient(
    Deno.env.get('SUPABASE_URL')!,
    Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
  );

  // Get all active subscriptions
  const { data: subscriptions } = await supabase
    .from('society_subscriptions')
    .select('*')
    .eq('status', 'ACTIVE')
    .lte('next_billing_date', new Date().toISOString());

  for (const subscription of subscriptions) {
    // Generate invoice
    await supabase.from('subscription_invoices').insert({
      subscription_id: subscription.subscription_id,
      society_id: subscription.society_id,
      amount: subscription.monthly_fee,
      due_date: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000), // 7 days
      status: 'PENDING'
    });

    // Update next billing date
    await supabase
      .from('society_subscriptions')
      .update({
        next_billing_date: new Date(
          subscription.next_billing_date.getTime() + 30 * 24 * 60 * 60 * 1000
        )
      })
      .eq('subscription_id', subscription.subscription_id);
  }

  return new Response(
    JSON.stringify({ invoices_generated: subscriptions.length }),
    { headers: { 'Content-Type': 'application/json' } }
  );
});
```

**Configure Cron Schedule in Supabase Dashboard:**
```
Function: generate-monthly-invoices
Schedule: 0 2 1 * * (2 AM on 1st of every month)
```

### 8.3 Webhook Handlers

**Example: Razorpay Payment Webhook**

```typescript
// supabase/functions/razorpay-webhook/index.ts
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';
import { createHmac } from 'https://deno.land/std@0.168.0/node/crypto.ts';

serve(async (req) => {
  const supabase = createClient(
    Deno.env.get('SUPABASE_URL')!,
    Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
  );

  // Verify webhook signature
  const signature = req.headers.get('x-razorpay-signature');
  const body = await req.text();

  const expectedSignature = createHmac('sha256', Deno.env.get('RAZORPAY_WEBHOOK_SECRET')!)
    .update(body)
    .digest('hex');

  if (signature !== expectedSignature) {
    return new Response('Invalid signature', { status: 401 });
  }

  const payload = JSON.parse(body);
  const event = payload.event;

  // Handle payment success
  if (event === 'payment.captured') {
    const paymentId = payload.payload.payment.entity.id;
    const amount = payload.payload.payment.entity.amount / 100; // Paise to rupees

    // Update payment record
    await supabase
      .from('payments')
      .update({
        status: 'COMPLETED',
        razorpay_payment_id: paymentId,
        paid_at: new Date().toISOString()
      })
      .eq('razorpay_order_id', payload.payload.payment.entity.order_id);

    // Update order status
    // ... additional logic
  }

  return new Response(JSON.stringify({ received: true }), {
    headers: { 'Content-Type': 'application/json' }
  });
});
```

**Register webhook URL with Razorpay:**
```
https://your-project.supabase.co/functions/v1/razorpay-webhook
```

---

## 9. Third-Party Integrations

### 9.1 Razorpay (Payments)

**Backend Integration:**

```typescript
// src/services/payment-service.ts
import Razorpay from 'razorpay';

const razorpay = new Razorpay({
  key_id: process.env.RAZORPAY_KEY_ID!,
  key_secret: process.env.RAZORPAY_KEY_SECRET!
});

export class PaymentService {
  static async createPaymentLink(orderId: string, amount: number) {
    const order = await razorpay.orders.create({
      amount: amount * 100, // Convert to paise
      currency: 'INR',
      receipt: orderId,
      notes: {
        order_id: orderId
      }
    });

    return {
      razorpay_order_id: order.id,
      amount: order.amount,
      currency: order.currency
    };
  }

  static async verifyPayment(
    razorpayOrderId: string,
    razorpayPaymentId: string,
    razorpaySignature: string
  ) {
    const text = `${razorpayOrderId}|${razorpayPaymentId}`;
    const expectedSignature = crypto
      .createHmac('sha256', process.env.RAZORPAY_KEY_SECRET!)
      .update(text)
      .digest('hex');

    return expectedSignature === razorpaySignature;
  }
}
```

**Mobile App Integration:**

```dart
// Flutter: Razorpay checkout
import 'package:razorpay_flutter/razorpay_flutter.dart';

Future<void> initiatePayment(String orderId, double amount) async {
  // 1. Get Razorpay order ID from backend
  final response = await apiClient.post('/payments/create', data: {
    'order_id': orderId,
    'amount': amount
  });

  final razorpayOrderId = response['razorpay_order_id'];

  // 2. Open Razorpay checkout
  final razorpay = Razorpay();

  razorpay.on(Razorpay.EVENT_PAYMENT_SUCCESS, (PaymentSuccessResponse response) {
    // 3. Verify payment with backend
    apiClient.post('/payments/verify', data: {
      'razorpay_order_id': response.orderId,
      'razorpay_payment_id': response.paymentId,
      'razorpay_signature': response.signature
    });
  });

  razorpay.open({
    'key': 'rzp_test_xxxxx', // From env
    'amount': (amount * 100).toInt(),
    'name': 'Laundry App',
    'order_id': razorpayOrderId,
    'prefill': {
      'contact': userPhone,
      'email': userEmail
    }
  });
}
```

### 9.2 Firebase Cloud Messaging (Push Notifications)

**Backend: Trigger Notification**

```typescript
// Backend calls Supabase Edge Function
await supabase.functions.invoke('send-notification', {
  body: {
    user_id: vendorId,
    title: 'New Order',
    body: 'You have a new order',
    data: { order_id: orderId, type: 'new_order' }
  }
});
```

**Mobile App: Receive Notification**

```dart
// Flutter: Initialize FCM
import 'package:firebase_messaging/firebase_messaging.dart';

Future<void> initializeNotifications() async {
  final messaging = FirebaseMessaging.instance;

  // Request permission
  await messaging.requestPermission();

  // Get FCM token
  final token = await messaging.getToken();

  // Send token to backend
  await apiClient.post('/users/me/fcm-token', data: {'token': token});

  // Listen for messages
  FirebaseMessaging.onMessage.listen((RemoteMessage message) {
    // Show in-app notification
    showNotification(
      title: message.notification?.title ?? '',
      body: message.notification?.body ?? ''
    );
  });

  // Handle notification tap
  FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
    // Navigate to relevant screen
    final orderId = message.data['order_id'];
    navigateToOrder(orderId);
  });
}
```

### 9.3 SMS (Twilio / MSG91)

**Edge Function: Send SMS**

```typescript
// supabase/functions/send-sms/index.ts
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';

serve(async (req) => {
  const { phone, message } = await req.json();

  // Using MSG91 (Indian SMS provider)
  const response = await fetch('https://api.msg91.com/api/v5/flow/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'authkey': Deno.env.get('MSG91_AUTH_KEY')!
    },
    body: JSON.stringify({
      flow_id: Deno.env.get('MSG91_FLOW_ID'),
      sender: Deno.env.get('MSG91_SENDER_ID'),
      mobiles: phone,
      message: message
    })
  });

  return new Response(JSON.stringify({ success: true }), {
    headers: { 'Content-Type': 'application/json' }
  });
});
```

### 9.4 Email Service (SendGrid / Resend)

**Why Email is Needed:**
- Vendor registration approval/rejection notifications
- Society subscription invoices (monthly)
- Order confirmation receipts (optional, backup to SMS/push)
- Payment overdue reminders for society admins
- Dispute notifications to all parties
- Admin reports and analytics summaries

**Recommended Provider: Resend**
- Developer-friendly API
- Excellent deliverability
- React Email template support
- Free tier: 100 emails/day (3000/month)
- Paid: $20/month for 50k emails

**Alternative: SendGrid**
- More established, larger scale
- Free tier: 100 emails/day
- More complex setup but very reliable

**Backend Integration (Resend):**

```typescript
// src/services/email-service.ts
import { Resend } from 'resend';

const resend = new Resend(process.env.RESEND_API_KEY);

export class EmailService {
  /**
   * Send vendor approval notification
   */
  static async sendVendorApproval(
    vendorEmail: string,
    vendorName: string,
    societyName: string
  ) {
    try {
      await resend.emails.send({
        from: 'noreply@yourapp.com',
        to: vendorEmail,
        subject: `Vendor Registration Approved - ${societyName}`,
        html: `
          <h2>Congratulations ${vendorName}!</h2>
          <p>Your vendor registration for <strong>${societyName}</strong> has been approved.</p>
          <p>You can now start receiving orders from residents.</p>
          <p>Login to your vendor app to get started.</p>
        `
      });

      return { success: true };
    } catch (error) {
      console.error('Email send failed:', error);
      // Don't throw - email is non-critical
      return { success: false, error };
    }
  }

  /**
   * Send vendor rejection notification
   */
  static async sendVendorRejection(
    vendorEmail: string,
    vendorName: string,
    societyName: string,
    reason?: string
  ) {
    await resend.emails.send({
      from: 'noreply@yourapp.com',
      to: vendorEmail,
      subject: `Vendor Registration Update - ${societyName}`,
      html: `
        <h2>Hello ${vendorName},</h2>
        <p>Thank you for your interest in serving <strong>${societyName}</strong>.</p>
        <p>Unfortunately, your registration was not approved at this time.</p>
        ${reason ? `<p><strong>Reason:</strong> ${reason}</p>` : ''}
        <p>You can contact the society admin for more information.</p>
      `
    });
  }

  /**
   * Send monthly subscription invoice to society admin
   */
  static async sendSubscriptionInvoice(
    adminEmail: string,
    societyName: string,
    invoiceData: {
      invoice_number: string;
      amount: number;
      due_date: string;
      billing_period: string;
    }
  ) {
    await resend.emails.send({
      from: 'billing@yourapp.com',
      to: adminEmail,
      subject: `Invoice #${invoiceData.invoice_number} - ${societyName}`,
      html: `
        <h2>Monthly Subscription Invoice</h2>
        <p><strong>Society:</strong> ${societyName}</p>
        <p><strong>Invoice Number:</strong> ${invoiceData.invoice_number}</p>
        <p><strong>Billing Period:</strong> ${invoiceData.billing_period}</p>
        <p><strong>Amount:</strong> ₹${invoiceData.amount}</p>
        <p><strong>Due Date:</strong> ${invoiceData.due_date}</p>
        <br>
        <p>Please make payment to avoid service interruption.</p>
        <p>Payment can be made via bank transfer or UPI.</p>
      `
    });
  }

  /**
   * Send payment overdue reminder
   */
  static async sendPaymentOverdueReminder(
    adminEmail: string,
    societyName: string,
    invoiceNumber: string,
    overdueAmount: number,
    daysPastDue: number
  ) {
    await resend.emails.send({
      from: 'billing@yourapp.com',
      to: adminEmail,
      subject: `Payment Overdue - Invoice #${invoiceNumber}`,
      html: `
        <h2>Payment Reminder</h2>
        <p>Dear ${societyName} Admin,</p>
        <p>Your invoice #${invoiceNumber} is <strong>${daysPastDue} days overdue</strong>.</p>
        <p><strong>Amount Due:</strong> ₹${overdueAmount}</p>
        <p>Please make payment immediately to avoid service suspension.</p>
        <p>If payment has already been made, please ignore this reminder.</p>
      `
    });
  }

  /**
   * Send dispute notification
   */
  static async sendDisputeNotification(
    recipientEmail: string,
    recipientName: string,
    orderNumber: string,
    disputeDetails: string
  ) {
    await resend.emails.send({
      from: 'support@yourapp.com',
      to: recipientEmail,
      subject: `Dispute Raised - Order #${orderNumber}`,
      html: `
        <h2>Dispute Notification</h2>
        <p>Hello ${recipientName},</p>
        <p>A dispute has been raised for order <strong>#${orderNumber}</strong>.</p>
        <p><strong>Details:</strong> ${disputeDetails}</p>
        <p>Please login to your app to respond to this dispute.</p>
      `
    });
  }

  /**
   * Send order confirmation (optional, as backup to push/SMS)
   */
  static async sendOrderConfirmation(
    residentEmail: string,
    orderNumber: string,
    orderDetails: {
      vendor_name: string;
      total: number;
      pickup_time: string;
      expected_delivery: string;
    }
  ) {
    await resend.emails.send({
      from: 'orders@yourapp.com',
      to: residentEmail,
      subject: `Order Confirmation - #${orderNumber}`,
      html: `
        <h2>Order Confirmed!</h2>
        <p><strong>Order Number:</strong> ${orderNumber}</p>
        <p><strong>Service Provider:</strong> ${orderDetails.vendor_name}</p>
        <p><strong>Total:</strong> ₹${orderDetails.total}</p>
        <p><strong>Pickup Time:</strong> ${orderDetails.pickup_time}</p>
        <p><strong>Expected Delivery:</strong> ${orderDetails.expected_delivery}</p>
        <br>
        <p>Track your order in the app for real-time updates.</p>
      `
    });
  }
}
```

**Usage in Backend API:**

```typescript
// api/v1/admin/vendors/approve.ts
import { EmailService } from '@/services/email-service';

export default async (req: Request, res: Response) => {
  const { vendor_id } = req.params;

  // Approve vendor in database
  const vendor = await VendorRepository.approveVendor(vendor_id);

  // Send email notification (non-blocking)
  EmailService.sendVendorApproval(
    vendor.email,
    vendor.business_name,
    vendor.society_name
  ).catch(err => console.error('Email failed:', err));

  res.json({ success: true });
};
```

**Edge Function: Monthly Invoice Email (Cron Job)**

```typescript
// supabase/functions/send-monthly-invoices/index.ts
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';

serve(async (req) => {
  const supabase = createClient(
    Deno.env.get('SUPABASE_URL')!,
    Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
  );

  // Get all pending invoices
  const { data: invoices } = await supabase
    .from('subscription_invoices')
    .select(`
      *,
      society_subscriptions(
        society_id,
        societies(name, admin_email)
      )
    `)
    .eq('status', 'PENDING')
    .gte('due_date', new Date().toISOString());

  for (const invoice of invoices) {
    // Send email via Resend API
    await fetch('https://api.resend.com/emails', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${Deno.env.get('RESEND_API_KEY')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        from: 'billing@yourapp.com',
        to: invoice.society_subscriptions.societies.admin_email,
        subject: `Invoice #${invoice.invoice_number}`,
        html: `Invoice details...`
      })
    });
  }

  return new Response(JSON.stringify({ sent: invoices.length }));
});
```

**Cost Estimate:**
```
Resend Free Tier: 3,000 emails/month (sufficient for start)

Expected usage (100 societies):
- Vendor approvals: ~50/month
- Monthly invoices: 100/month
- Payment reminders: ~20/month
- Disputes: ~30/month
- Order confirmations (optional): ~2,000/month

Total: ~2,200 emails/month (within free tier)

At scale (500 societies):
- Total: ~10,000 emails/month
- Cost: $20/month (Resend paid plan)
```

### 9.5 Error Tracking & Monitoring (Sentry)

**Why Sentry:**
- **Real-time error tracking** across all platforms
- **Stack traces** with source maps
- **User context** (which user hit the error)
- **Release tracking** (which deployment caused issues)
- **Performance monitoring** (slow API endpoints)
- **Alerts** via email/Slack when critical errors occur

**Backend Integration:**

```typescript
// src/utils/sentry.ts
import * as Sentry from '@sentry/node';
import { ProfilingIntegration } from '@sentry/profiling-node';

Sentry.init({
  dsn: process.env.SENTRY_DSN,
  environment: process.env.NODE_ENV,
  integrations: [
    new ProfilingIntegration(),
  ],
  tracesSampleRate: 0.1, // 10% of requests tracked
  profilesSampleRate: 0.1,
});

// Express middleware
export const sentryRequestHandler = Sentry.Handlers.requestHandler();
export const sentryErrorHandler = Sentry.Handlers.errorHandler();
```

```typescript
// api/index.ts (Express app)
import express from 'express';
import { sentryRequestHandler, sentryErrorHandler } from '@/utils/sentry';

const app = express();

// Sentry must be first middleware
app.use(sentryRequestHandler);

// ... your routes

// Sentry error handler must be before other error middleware
app.use(sentryErrorHandler);

// Global error handler
app.use((err, req, res, next) => {
  // Sentry already captured the error
  console.error(err);
  res.status(500).json({
    success: false,
    error: { message: 'Internal server error' }
  });
});
```

**Capture Custom Errors:**

```typescript
// src/services/order-service.ts
import * as Sentry from '@sentry/node';

export class OrderService {
  static async createOrder(userId: string, orderData: any) {
    try {
      // Order creation logic
      const order = await OrderRepository.createOrder(orderData);
      return order;
    } catch (error) {
      // Capture error with context
      Sentry.captureException(error, {
        tags: {
          service: 'order-creation',
          user_id: userId
        },
        extra: {
          order_data: orderData
        },
        level: 'error'
      });

      throw error;
    }
  }
}
```

**Flutter Mobile App Integration:**

```yaml
# pubspec.yaml
dependencies:
  sentry_flutter: ^7.14.0
```

```dart
// lib/main.dart
import 'package:sentry_flutter/sentry_flutter.dart';

Future<void> main() async {
  await SentryFlutter.init(
    (options) {
      options.dsn = 'YOUR_SENTRY_DSN';
      options.environment = 'production';
      options.tracesSampleRate = 0.1;
    },
    appRunner: () => runApp(MyApp()),
  );
}

// Capture errors manually
try {
  await createOrder(orderData);
} catch (error, stackTrace) {
  await Sentry.captureException(
    error,
    stackTrace: stackTrace,
    hint: Hint.withMap({
      'order_data': orderData,
      'user_id': userId,
    }),
  );

  // Show error to user
  showErrorDialog(error.toString());
}
```

**Next.js Admin Web Integration:**

```bash
npm install @sentry/nextjs
npx @sentry/wizard@latest -i nextjs
```

```typescript
// sentry.client.config.ts
import * as Sentry from '@sentry/nextjs';

Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  environment: process.env.NODE_ENV,
  tracesSampleRate: 0.1,

  // Replay sessions for debugging
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1.0,
});
```

**Supabase Edge Functions:**

```typescript
// supabase/functions/_shared/sentry.ts
import * as Sentry from 'https://deno.land/x/sentry/index.mjs';

Sentry.init({
  dsn: Deno.env.get('SENTRY_DSN'),
  environment: 'production',
});

export const captureEdgeFunctionError = (error: Error, context?: any) => {
  Sentry.captureException(error, {
    tags: { runtime: 'deno-edge' },
    extra: context,
  });
};
```

**Cost Estimate:**
```
Sentry Pricing:
- Developer (Free): 5,000 errors/month, 1 user
- Team ($26/month): 50,000 errors/month, unlimited users
- Business ($80/month): 500,000 errors/month

Expected usage:
- Early stage: <5,000 errors/month (FREE)
- Growth: 10,000-50,000 errors/month ($26/month)

Recommended: Start with free tier, upgrade to Team when needed
```

**Sentry Alerts Configuration:**

```yaml
# .sentryclirc
[alerts]
- name: Critical API Errors
  conditions:
    - event.level: error
    - event.tags.severity: critical
  actions:
    - email: team@yourapp.com
    - slack: #alerts-channel

- name: High Error Rate
  conditions:
    - event.frequency: >100/hour
  actions:
    - email: oncall@yourapp.com
```

---

## 10. Deployment Strategy

### 10.1 Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      PRODUCTION                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │   Backend    │  │   Admin Web  │  │  Mobile Apps    │   │
│  │   Node.js    │  │   Next.js    │  │    Flutter      │   │
│  │              │  │              │  │                 │   │
│  │   Vercel     │  │   Vercel     │  │  App Store      │   │
│  │  Production  │  │  Production  │  │  Play Store     │   │
│  └──────────────┘  └──────────────┘  └─────────────────┘   │
│         │                 │                                  │
│         └─────────────────┴──────────────────┐              │
│                                               │              │
│                                               ▼              │
│                                    ┌─────────────────────┐  │
│                                    │    Supabase         │  │
│                                    │    Production       │  │
│                                    │                     │  │
│                                    │  • PostgreSQL       │  │
│                                    │  • Edge Functions   │  │
│                                    │  • Storage          │  │
│                                    │  • Auth             │  │
│                                    └─────────────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                      STAGING                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐  ┌──────────────┐                         │
│  │   Backend    │  │   Admin Web  │                         │
│  │   Vercel     │  │   Vercel     │                         │
│  │   Preview    │  │   Preview    │                         │
│  └──────────────┘  └──────────────┘                         │
│         │                 │                                  │
│         └─────────────────┴──────────────────┐              │
│                                               ▼              │
│                                    ┌─────────────────────┐  │
│                                    │    Supabase         │  │
│                                    │    Staging          │  │
│                                    └─────────────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 10.2 Backend API Deployment (Railway/VPS)

**Option 1: Railway.app (Recommended for MVP)**

1. **Connect GitHub repository to Railway**
2. **Configure environment variables** in Railway dashboard:
   ```
   PORT=8080
   ENV=production
   SUPABASE_URL=https://xxx.supabase.co
   SUPABASE_SERVICE_KEY=xxx
   DATABASE_URL=postgresql://...
   JWT_SECRET=xxx
   RAZORPAY_KEY_ID=xxx
   RAZORPAY_KEY_SECRET=xxx
   FCM_SERVER_KEY=xxx
   ```

3. **Railway will auto-detect Dockerfile and deploy**

4. **Deploy:**
   ```bash
   # Automatic on git push to main
   git push origin main

   # Or use Railway CLI
   railway up
   ```

**Option 2: VPS (DigitalOcean/AWS EC2)**

1. **Build Docker image:**
   ```bash
   docker build -t society-api:latest .
   ```

2. **Run container:**
   ```bash
   docker run -d \
     -p 8080:8080 \
     --env-file .env \
     --name society-api \
     society-api:latest
   ```

3. **Setup reverse proxy (Nginx):**
   ```nginx
   server {
       listen 80;
       server_name api.yourapp.com;

       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

4. **Setup SSL with Let's Encrypt:**
   ```bash
   certbot --nginx -d api.yourapp.com
   ```

**URL Structure:**
```
Production: https://api.yourapp.com/api/v1/*
Staging: https://api-staging.yourapp.com/api/v1/*
```

### 10.3 Admin Web Deployment (Vercel)

**Setup:**

1. **Connect GitHub repository**
2. **Configure build settings:**
   ```
   Framework: Next.js
   Build Command: npm run build
   Output Directory: .next
   Install Command: npm install
   ```

3. **Environment variables:**
   ```
   NEXT_PUBLIC_API_URL=https://api.yourapp.com
   ```

4. **Deploy:** Automatic on push to main

**URL:**
```
Production: https://admin.yourapp.com
```

### 10.4 Mobile App Deployment (Flutter)

**iOS (App Store):**

```bash
# 1. Build release version
flutter build ios --release

# 2. Open Xcode
open ios/Runner.xcworkspace

# 3. Archive and upload to App Store Connect

# 4. Submit for review
```

**Android (Play Store):**

```bash
# 1. Build release APK/AAB
flutter build appbundle --release

# 2. Sign with release keystore
# (configured in android/app/build.gradle)

# 3. Upload to Play Console

# 4. Submit for review
```

**Environment Configuration:**

```dart
// lib/config/env.dart
class Environment {
  static const String apiUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'https://api.yourapp.com/api/v1'
  );
}

// Build with environment
flutter build apk --dart-define=API_URL=https://api.yourapp.com/api/v1
```

### 10.5 Supabase Setup

**Production Project:**

1. **Create project** on supabase.com
2. **Run migrations:**
   ```bash
   supabase db push --linked
   ```

3. **Deploy edge functions:**
   ```bash
   supabase functions deploy send-notification
   supabase functions deploy razorpay-webhook
   supabase functions deploy generate-monthly-invoices
   ```

4. **Configure cron jobs** in Supabase dashboard

5. **Set up storage buckets:**
   ```
   - order-photos (public)
   - vendor-documents (private)
   - invoices (private)
   ```

### 10.6 CI/CD Pipeline

**GitHub Actions Workflow:**

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: npm install
      - run: npm test
      - uses: amondnet/vercel-action@v25
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}

  deploy-admin:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: cd admin-web && npm install
      - run: cd admin-web && npm run build
      - uses: amondnet/vercel-action@v25
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_ADMIN_PROJECT_ID }}

  deploy-functions:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: supabase/setup-cli@v1
      - run: supabase functions deploy --project-ref ${{ secrets.SUPABASE_PROJECT_REF }}
```

---

## 11. Development Workflow

### 11.1 Local Development Setup

**Prerequisites:**
```bash
# Install Node.js 20 LTS
nvm install 20
nvm use 20

# Install Flutter
# Download from flutter.dev

# Install Supabase CLI
brew install supabase/tap/supabase
```

**Setup Steps:**

```bash
# 1. Clone repository
git clone https://github.com/yourorg/laundry-app.git
cd laundry-app

# 2. Install backend dependencies
npm install

# 3. Setup environment variables
cp .env.example .env
# Edit .env with your values

# 4. Start Supabase locally
supabase start

# 5. Run database migrations
supabase db reset

# 6. Start backend dev server
npm run dev
# Backend runs on http://localhost:3000

# 7. Start admin web
cd admin-web
npm install
npm run dev
# Admin runs on http://localhost:3001

# 8. Run Flutter app
cd ../resident-app
flutter pub get
flutter run
```

### 11.2 Development URLs

```
Backend API: http://localhost:3000/api/v1
Admin Web: http://localhost:3001
Supabase Studio: http://localhost:54323
PostgreSQL: postgresql://postgres:postgres@localhost:54322/postgres
```

### 11.3 Testing Strategy

**Backend API Tests:**

```typescript
// tests/orders/create-order.test.ts
import { describe, it, expect } from 'vitest';
import request from 'supertest';
import app from '@/api/index';

describe('POST /api/v1/orders', () => {
  it('creates order successfully', async () => {
    const response = await request(app)
      .post('/api/v1/orders')
      .set('Authorization', `Bearer ${testToken}`)
      .send({
        laundry_id: 'test-vendor-id',
        items: [
          {
            service_id: 1,
            item_name: 'Shirt',
            quantity: 5,
            unit_price: 10
          }
        ],
        pickup_datetime: '2025-11-18T10:00:00Z',
        pickup_address: 'A-404'
      });

    expect(response.status).toBe(201);
    expect(response.body.success).toBe(true);
    expect(response.body.data).toHaveProperty('order_id');
  });

  it('validates required fields', async () => {
    const response = await request(app)
      .post('/api/v1/orders')
      .set('Authorization', `Bearer ${testToken}`)
      .send({});

    expect(response.status).toBe(400);
    expect(response.body.success).toBe(false);
  });
});
```

**Flutter Widget Tests:**

```dart
// test/widgets/order_card_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:resident_app/features/orders/presentation/widgets/order_card.dart';

void main() {
  testWidgets('OrderCard displays order details', (tester) async {
    final order = Order(
      orderId: '123',
      orderNumber: 'ORD001',
      status: 'PICKUP_SCHEDULED',
      total: 350,
    );

    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: OrderCard(order: order),
        ),
      ),
    );

    expect(find.text('ORD001'), findsOneWidget);
    expect(find.text('₹350'), findsOneWidget);
  });
}
```

**Run Tests:**

```bash
# Backend tests
npm test

# Flutter tests
flutter test
```

---

## 12. Infrastructure Costs

### 12.1 Monthly Cost Breakdown

**Vercel (Backend + Admin Web):**
```
Plan: Pro ($20/month)
Includes:
  - Unlimited deployments
  - 100GB bandwidth
  - Serverless function execution
  - Custom domains
  - Team collaboration

Expected usage:
  - Backend API: ~50GB bandwidth
  - Admin Web: ~20GB bandwidth
  - Well within limits

Cost: $20/month
```

**Supabase:**
```
Plan: Pro ($25/month)
Includes:
  - 8GB database
  - 100GB bandwidth
  - 100GB file storage
  - 500K edge function invocations
  - 50GB egress
  - Daily backups

Expected usage (100 societies):
  - Database: ~2GB
  - Storage: ~20GB (photos)
  - Bandwidth: ~30GB
  - Edge functions: ~100K/month

Cost: $25/month
```

**Razorpay:**
```
Transaction fee: 2% per transaction
Payment gateway: Free to integrate

Expected usage:
  - 0 transactions (residents pay vendors directly)
  - Only used for society subscription payments

Cost: Minimal (~₹500/month or $6)
```

**Firebase Cloud Messaging:**
```
Free tier: Unlimited notifications

Cost: $0
```

**SMS (MSG91):**
```
Cost per SMS: ₹0.20 ($0.0024)
OTPs per month: ~2000 (100 societies × 20 users × 1 OTP/month)

Cost: ₹400/month ($5)
```

**Domain & SSL:**
```
Domain: $12/year
SSL: Free (Vercel provides)

Cost: $1/month
```

**Total Monthly Cost:**
```
Vercel: $20
Supabase: $25
Razorpay: $6
FCM: $0
SMS: $5
Email (Resend): $0 (free tier, 3k/month)
Sentry: $0 (free tier, 5k errors/month)
Domain: $1
-----------
Total: $57/month (~₹4,750/month)
```

**Cost at Scale (500 societies):**
```
Vercel: $20 (same)
Supabase: $50 (Pro+ for more database/bandwidth)
Razorpay: $30 (more subscription payments)
SMS: $20 (more OTPs)
Email (Resend): $20 (50k emails)
Sentry: $26 (error tracking)
Domain: $1
-----------
Total: $167/month (~₹14,000/month)
```

**Revenue vs Cost:**
```
50 societies @ ₹10k avg = ₹5L/month revenue
Infrastructure cost: ₹5k/month
Profit margin: 99%

500 societies @ ₹10k avg = ₹50L/month revenue
Infrastructure cost: ₹14k/month
Profit margin: 97%
```

### 12.2 Cost Optimization

**Strategies:**
1. **Vercel:** Use edge caching for static content
2. **Supabase:** Optimize queries, use connection pooling
3. **Storage:** Compress images before upload, use CDN
4. **SMS:** Implement rate limiting on OTP requests
5. **Edge Functions:** Use cron jobs instead of realtime triggers where possible
6. **Email:** Batch non-urgent emails, use free tier efficiently
7. **Sentry:** Filter out non-critical errors, use sampling for performance tracking

---

## 13. Security & Performance

### 13.1 Security Measures

**API Security:**

```typescript
// Rate limiting
import rateLimit from 'express-rate-limit';

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // 100 requests per window
  message: 'Too many requests, please try again later'
});

app.use('/api/', limiter);

// Helmet for security headers
import helmet from 'helmet';
app.use(helmet());

// CORS configuration
import cors from 'cors';
app.use(cors({
  origin: [
    'https://admin.yourapp.com',
    'https://yourapp.com'
  ],
  credentials: true
}));

// Input validation with Zod
import { z } from 'zod';

const createOrderSchema = z.object({
  laundry_id: z.string().uuid(),
  items: z.array(z.object({
    service_id: z.number().int().positive(),
    quantity: z.number().int().positive().max(100)
  })).min(1).max(50)
});
```

**Database Security:**

```sql
-- Row Level Security policies
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Residents view own orders"
ON orders FOR SELECT
USING (auth.uid() = resident_id);

-- Encrypted columns for sensitive data
CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE users
ADD COLUMN phone_encrypted BYTEA;
```

**Mobile App Security:**

```dart
// Secure storage for tokens
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

final storage = FlutterSecureStorage();

// Store token
await storage.write(key: 'auth_token', value: token);

// Retrieve token
final token = await storage.read(key: 'auth_token');

// SSL pinning (production)
import 'package:dio/dio.dart';

final dio = Dio();
dio.httpClientAdapter = IOHttpClientAdapter(
  onHttpClientCreate: (client) {
    client.badCertificateCallback = (cert, host, port) {
      // Validate certificate
      return validateCertificate(cert);
    };
    return client;
  }
);
```

### 13.2 Performance Optimization

**Backend API:**

```typescript
// Response compression
import compression from 'compression';
app.use(compression());

// Database connection pooling
import { createClient } from '@supabase/supabase-js';

const supabase = createClient(url, key, {
  db: {
    schema: 'public'
  },
  global: {
    headers: {
      'x-connection-pool': 'true'
    }
  }
});

// Caching with Redis (future)
import Redis from 'ioredis';
const redis = new Redis(process.env.REDIS_URL);

// Cache vendor details
const getCachedVendor = async (vendorId: string) => {
  const cached = await redis.get(`vendor:${vendorId}`);
  if (cached) return JSON.parse(cached);

  const vendor = await fetchVendorFromDB(vendorId);
  await redis.setex(`vendor:${vendorId}`, 3600, JSON.stringify(vendor));
  return vendor;
};
```

**Database Optimization:**

```sql
-- Indexes for common queries
CREATE INDEX idx_orders_resident ON orders(resident_id, created_at DESC);
CREATE INDEX idx_orders_vendor ON orders(laundry_id, status, created_at DESC);
CREATE INDEX idx_orders_society ON orders(society_id, created_at DESC);

-- Materialized view for analytics
CREATE MATERIALIZED VIEW vendor_analytics AS
SELECT
  laundry_id,
  COUNT(*) as total_orders,
  SUM(final_price) as total_revenue,
  AVG(rating) as avg_rating
FROM orders
GROUP BY laundry_id;

-- Refresh periodically
REFRESH MATERIALIZED VIEW vendor_analytics;
```

**Flutter App Performance:**

```dart
// Image caching
import 'package:cached_network_image/cached_network_image.dart';

CachedNetworkImage(
  imageUrl: vendor.photoUrl,
  placeholder: (context, url) => CircularProgressIndicator(),
  errorWidget: (context, url, error) => Icon(Icons.error),
  cacheManager: CacheManager(
    Config(
      'vendor_images',
      stalePeriod: Duration(days: 7),
      maxNrOfCacheObjects: 100,
    ),
  ),
);

// Lazy loading lists
ListView.builder(
  itemCount: orders.length,
  itemBuilder: (context, index) {
    return OrderCard(order: orders[index]);
  },
);
```

---

## Summary

### Technology Stack Overview

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Mobile** | Flutter | iOS + Android apps (UI only) |
| **Backend** | Node.js + Express | All business logic |
| **Admin Web** | Next.js 14 | Admin dashboard |
| **Database** | PostgreSQL (Supabase) | Data storage |
| **Edge Functions** | Deno (Supabase) | Webhooks, crons |
| **Storage** | Supabase Storage | File uploads |
| **Auth** | Supabase Auth | User authentication |
| **Hosting** | Vercel | Backend + Web |
| **Payments** | Razorpay | Payment gateway |
| **Notifications** | Firebase FCM | Push notifications |
| **Email** | SendGrid/Resend | Transactional emails |
| **Monitoring** | Sentry | Error tracking |

### Key Architecture Decisions

✅ **Clean separation:** Mobile apps are thin clients, all logic in backend
✅ **API-first:** Single backend API serves mobile and web
✅ **Serverless:** Zero server management, auto-scaling
✅ **Type-safe:** TypeScript everywhere
✅ **Cost-effective:** ~$60/month for 100 societies
✅ **Scalable:** Can handle 1000+ societies without architecture changes

### Deployment Flow

```
Code Push → GitHub
    ↓
GitHub Actions CI/CD
    ↓
├─→ Vercel (Backend API)
├─→ Vercel (Admin Web)
├─→ Supabase (Edge Functions)
└─→ App Store / Play Store (Mobile Apps)
```

**Zero downtime deployments. Automatic rollbacks on failure.**

---

**End of Document**
