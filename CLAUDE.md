# Society Service Platform - Claude Code Guidelines

## Project Overview

Multi-service platform for residential society management. Residents order services (laundry, vehicle care, home services) from vendors assigned to their society.

## Architecture

```
society-service-app/
├── backend/                    # Go 1.23 + Gin + pgx
├── apps/
│   ├── resident_app/           # Flutter - Resident mobile app
│   ├── vendor_app/             # Flutter - Vendor mobile app
│   ├── society-admin-web/      # Next.js 14 - Society admin dashboard
│   └── platform-admin-web/     # Next.js 14 - Platform admin dashboard
├── packages/
│   └── shared-types/           # Shared TypeScript types
├── supabase/                   # Database migrations & edge functions
└── docs/service_app_docs/      # Comprehensive documentation
```

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.23, Gin, pgx/v5, SQLC |
| Mobile | Flutter 3.10+, Riverpod, GoRouter, Dio, Freezed |
| Web | Next.js 14, React Query, React Hook Form, Zod, shadcn/ui |
| Database | PostgreSQL 15 (Supabase) with ltree extension |
| Payments | Razorpay |
| Notifications | Firebase Cloud Messaging, MSG91 (SMS) |

## Testing Requirements

**ALL code changes MUST include:**

1. **Unit Tests** - Business logic, services, providers, repositories
2. **Component/Widget Tests** - UI components and screens
3. **E2E Tests** - Critical user flows

### Test Commands

```bash
# Backend (Go)
cd backend && make test
cd backend && make test-coverage

# Mobile (Flutter)
cd apps/resident_app && flutter test
cd apps/vendor_app && flutter test

# Web (Next.js)
cd apps/society-admin-web && npm test
cd apps/platform-admin-web && npm test
```

### Coverage Requirements
- Backend: Minimum 80% coverage
- Mobile: Widget tests for all screens
- Web: Component tests for all components

---

## Coding Conventions

### Go Backend

**File Naming:** `snake_case.go`
**Types:** `PascalCase`
**Functions:** `PascalCase` (exported), `camelCase` (unexported)

**Clean Architecture Layers:**
```
cmd/api/           → Entry point, routes
internal/handlers/ → HTTP request/response handling
internal/services/ → Business logic
internal/repositories/ → Database operations (future)
internal/models/   → Data structures
internal/middleware/ → HTTP middleware
internal/utils/    → Helpers
```

**Handler Pattern:**
```go
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

    result, err := h.service.CreateOrder(ctx, &req)
    if err != nil {
        utils.RespondError(c, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), nil)
        return
    }

    utils.RespondSuccess(c, http.StatusCreated, result, "Order created")
}
```

**Test Pattern (Table-Driven):**
```go
func TestOrderService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   *CreateOrderRequest
        want    *Order
        wantErr bool
    }{
        {"valid order", &CreateOrderRequest{...}, &Order{...}, false},
        {"invalid input", &CreateOrderRequest{}, nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup, Execute, Assert
        })
    }
}
```

---

### Flutter Mobile

**Architecture:** Feature-first with clean architecture

```
lib/
├── core/
│   ├── api/              # Dio client, interceptors
│   ├── config/           # Environment, constants
│   ├── providers/        # Global Riverpod providers
│   └── router/           # GoRouter configuration
├── features/
│   └── {feature}/
│       ├── data/         # Repositories (API calls ONLY)
│       ├── domain/       # Models (Freezed)
│       └── presentation/
│           ├── screens/
│           └── widgets/
└── shared/widgets/       # Reusable components
```

**State Management (Riverpod):**
```dart
@riverpod
class OrderNotifier extends _$OrderNotifier {
  @override
  OrderState build() => const OrderState.initial();

  Future<void> createOrder(CreateOrderRequest request) async {
    state = const OrderState.loading();
    try {
      final order = await ref.read(orderRepositoryProvider).create(request);
      state = OrderState.success(order);
    } catch (e) {
      state = OrderState.error(e.toString());
    }
  }
}
```

**Models (Freezed):**
```dart
@freezed
class Order with _$Order {
  const factory Order({
    required String orderId,
    required String orderNumber,
    required OrderStatus status,
    @Default([]) List<OrderItem> items,
  }) = _Order;

  factory Order.fromJson(Map<String, dynamic> json) => _$OrderFromJson(json);
}
```

**Repository Pattern (API calls only, NO business logic):**
```dart
class OrderRepository {
  final ApiClient _client;
  OrderRepository(this._client);

  Future<Order> create(CreateOrderRequest request) async {
    final response = await _client.post('/orders', data: request.toJson());
    return Order.fromJson(response.data['data']);
  }
}
```

**Widget Test Pattern:**
```dart
void main() {
  testWidgets('OrderCard displays order number', (tester) async {
    await tester.pumpWidget(
      ProviderScope(
        child: MaterialApp(home: OrderCard(order: mockOrder)),
      ),
    );
    expect(find.text('Order #123'), findsOneWidget);
  });
}
```

---

### Next.js Web

**Architecture:** App Router with feature-based organization

```
src/
├── app/
│   ├── (auth)/           # Auth routes (login)
│   │   └── login/
│   ├── (dashboard)/      # Protected dashboard routes
│   │   ├── layout.tsx
│   │   ├── vendors/
│   │   ├── orders/
│   │   └── analytics/
│   └── layout.tsx
├── components/
│   ├── ui/               # shadcn/ui components
│   └── {feature}/        # Feature-specific components
├── hooks/                # React Query hooks
├── lib/
│   ├── api-client.ts     # Axios configuration
│   ├── utils.ts
│   └── validations.ts    # Zod schemas
└── types/
```

**Data Fetching (React Query):**
```typescript
export function useVendors(societyId: string) {
  return useQuery({
    queryKey: ['vendors', societyId],
    queryFn: async () => {
      const { data } = await apiClient.get(`/societies/${societyId}/vendors`);
      return data.data as Vendor[];
    },
  });
}

export function useApproveVendor() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (vendorId: string) =>
      apiClient.post(`/vendors/${vendorId}/approve`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vendors'] });
    },
  });
}
```

**Forms (React Hook Form + Zod):**
```typescript
const vendorSchema = z.object({
  businessName: z.string().min(2, 'Required'),
  phone: z.string().regex(/^[6-9]\d{9}$/, 'Invalid phone'),
});

function VendorForm() {
  const form = useForm<z.infer<typeof vendorSchema>>({
    resolver: zodResolver(vendorSchema),
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        {/* Form fields */}
      </form>
    </Form>
  );
}
```

**Component Test Pattern:**
```typescript
describe('VendorCard', () => {
  it('renders vendor name', () => {
    render(<VendorCard vendor={mockVendor} />);
    expect(screen.getByText('Test Vendor')).toBeInTheDocument();
  });
});
```

---

## API Standards

### Response Format (Success)
```json
{
  "success": true,
  "data": { },
  "message": "Optional message",
  "meta": {
    "timestamp": "2025-01-01T00:00:00Z",
    "request_id": "uuid-v4"
  }
}
```

### Response Format (Error)
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": { },
    "metadata": {
      "timestamp": "2025-01-01T00:00:00Z",
      "request_id": "uuid-v4"
    }
  }
}
```

### Common Error Codes
- `INVALID_REQUEST` - Malformed request body
- `UNAUTHORIZED` - Missing or invalid auth token
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource not found
- `CONFLICT` - Resource conflict (duplicate, etc.)
- `INTERNAL_ERROR` - Server error

---

## Database Conventions

- **Primary Keys:** UUID (preferred) or SERIAL with `GENERATED ALWAYS AS IDENTITY`
- **Timestamps:** `TIMESTAMP WITH TIME ZONE`, columns: `created_at`, `updated_at`
- **Soft Deletes:** `deleted_at` column (nullable)
- **Hierarchy:** Use ltree extension for tree structures
- **Naming:** `snake_case` for tables and columns

---

## Development Commands

```bash
# Backend
cd backend
make dev              # Hot reload development
make build            # Build binary
make test             # Run tests
make test-coverage    # Tests with coverage report
make lint             # Run linter
make fmt              # Format code

# Flutter
cd apps/resident_app  # or vendor_app
flutter run           # Run app
flutter test          # Run tests
flutter analyze       # Lint code
dart run build_runner build --delete-conflicting-outputs  # Generate code

# Web
cd apps/society-admin-web  # or platform-admin-web
npm run dev           # Development server
npm run build         # Production build
npm test              # Run tests
npm run lint          # Lint code
npm run format        # Format code

# Database
supabase start        # Start local Supabase
supabase db reset     # Reset database
./scripts/generate-types.sh  # Generate TypeScript types
```

---

## Environment Variables

### Backend (.env)
```
PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgresql://...
SUPABASE_URL=https://...
SUPABASE_ANON_KEY=...
JWT_SECRET=...
JWT_EXPIRY_HOURS=24
RAZORPAY_KEY_ID=...
RAZORPAY_KEY_SECRET=...
FCM_SERVER_KEY=...
```

### Web Apps (.env.local)
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_SUPABASE_URL=...
NEXT_PUBLIC_SUPABASE_ANON_KEY=...
```

---

## Key Domain Concepts

### User Types
- **Resident** - Orders services, lives in society
- **Vendor** - Provides services, assigned to societies/buildings
- **Society Admin** - Manages single society
- **Platform Admin** - Manages entire platform

### Order Flow
1. Resident creates order (one category per order)
2. Vendor accepts and schedules pickup
3. Workflow progresses per service type
4. Delivery and payment confirmation
5. Rating and completion

### Hierarchy Model (ltree)
Societies use generic tree structure:
- Society → Building → Floor → Unit
- Society → Phase → House
- Supports any depth, queried with path operations

---

## Agent System

This project uses specialized Claude Code agents for AI-assisted development. Each agent has deep knowledge of its platform's patterns, conventions, and testing requirements.

### Available Agents

| Agent | File | Scope | Purpose |
|-------|------|-------|---------|
| **Backend API** | `.claude/commands/backend-api.md` | `backend/**` | Go API development with clean architecture |
| **Mobile** | `.claude/commands/mobile.md` | `apps/*_app/**` | Flutter development with Riverpod, Freezed |
| **Frontend Web** | `.claude/commands/frontend-web.md` | `apps/*-web/**` | Next.js dashboards with React Query |
| **Orchestrator** | `.claude/commands/orchestrator.md` | All layers | Full-stack features, API contracts |

### How to Use Agents

**Option 1: Use Slash Commands (Recommended)**

Slash commands automatically apply the right agent's guidelines:

```
/api POST /orders/{id}/rating - Submit rating for completed order
/component flutter:resident:orders:RatingCard
/component web:society:vendors:VendorApprovalTable
/test backend/internal/services/order_service.go
/feature Add vendor rating system
/migrate create_ratings_table
```

**Option 2: Reference Agent in Conversation**

For more control, reference an agent directly:

```
Following the backend-api agent guidelines, create the order
cancellation endpoint with proper validation and error handling
```

```
Using the mobile agent patterns, implement the order tracking
screen with real-time status updates
```

```
Using the orchestrator approach, implement the payment confirmation
feature across backend, mobile, and web
```

### Slash Commands Reference

| Command | Purpose | Example |
|---------|---------|---------|
| `/api` | Create new API endpoint with handler, service, repository, and tests | `/api GET /vendors - List vendors with filtering` |
| `/component` | Scaffold component with test file | `/component flutter:vendor:orders:OrderCard` |
| `/test` | Generate comprehensive tests for any file | `/test backend/internal/handlers/order_handler.go` |
| `/feature` | Full-stack implementation across all layers | `/feature Add order cancellation with refunds` |
| `/migrate` | Create database migration with up/down | `/migrate add_rating_to_orders` |

### Agent Capabilities

**Backend API Agent** (`.claude/commands/backend-api.md`)
- Creates handlers, services, repositories following clean architecture
- Generates table-driven unit tests
- Generates handler tests with httptest
- Uses proper error codes and response formats
- Follows Go naming conventions

**Mobile Agent** (`.claude/commands/mobile.md`)
- Creates Freezed models with JSON serialization
- Implements Riverpod providers (StateNotifier, AsyncNotifier)
- Sets up GoRouter navigation
- Creates widget tests with ProviderScope
- Follows feature-first architecture

**Frontend Web Agent** (`.claude/commands/frontend-web.md`)
- Creates React Query hooks for data fetching
- Implements forms with React Hook Form + Zod
- Uses shadcn/ui components
- Creates component tests with Testing Library
- Follows Next.js App Router patterns

**Orchestrator Agent** (`.claude/commands/orchestrator.md`)
- Breaks down features into platform-specific tasks
- Defines API contracts before implementation
- Ensures type consistency across platforms
- Coordinates testing at all layers
- Verifies cross-platform integration

### Example Workflows

**Creating a New Feature (Full-Stack)**
```
Using the orchestrator, implement vendor ratings:
- Residents can rate vendors 1-5 stars after order completion
- Vendors see their average rating on profile
- Society admins see rating analytics
```

The orchestrator will:
1. Define API contract for `/orders/{id}/rating`
2. Create database migration for ratings table
3. Implement backend endpoint with tests
4. Implement mobile rating UI with tests
5. Implement web analytics with tests
6. Verify cross-platform consistency

**Creating a Backend Endpoint**
```
/api POST /orders/{id}/cancel - Cancel order with reason

Request: { "reason": "string" }
Response: { "order": Order, "refund_status": "pending" | "processed" }
Errors: ORDER_NOT_FOUND, ORDER_ALREADY_COMPLETED, CANCELLATION_NOT_ALLOWED
```

**Creating a Mobile Screen**
```
Using the mobile agent, create OrderTrackingScreen that:
- Shows real-time order status with timeline
- Displays vendor contact info
- Has pull-to-refresh
- Shows loading and error states
Include widget tests for all states
```

**Creating a Web Dashboard Component**
```
Using the frontend-web agent, create VendorAnalyticsCard that:
- Shows vendor's rating trend (chart)
- Displays order count and revenue
- Has date range filter
- Uses React Query for data fetching
Include component tests
```

### Testing Requirements (Enforced by All Agents)

Every code change includes tests:

| Layer | Required Tests |
|-------|---------------|
| Backend | Unit tests for services, handler tests for endpoints |
| Mobile | Widget tests for screens, provider tests for state |
| Web | Component tests, hook tests for React Query |
| Full-stack | E2E tests for critical user flows |

### Files Reference

```
.claude/
├── settings.json              # Hooks and permissions config
└── commands/
    ├── backend-api.md         # Backend agent prompt
    ├── mobile.md              # Mobile agent prompt
    ├── frontend-web.md        # Web agent prompt
    ├── orchestrator.md        # Orchestrator agent prompt
    ├── api.md                 # /api slash command
    ├── component.md           # /component slash command
    ├── test.md                # /test slash command
    ├── feature.md             # /feature slash command
    └── migrate.md             # /migrate slash command
```
