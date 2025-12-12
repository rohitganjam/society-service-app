# Full-Stack Feature Implementation

Implement feature: $ARGUMENTS

## Overview

This command orchestrates full-stack feature implementation across all layers:
1. Database schema changes
2. Backend API
3. Web admin dashboards
4. Mobile apps
5. Comprehensive testing

---

## Phase 1: Analysis

### 1.1 Understand the Feature

Answer these questions:
- What user problem does this solve?
- Which user types are affected?
  - [ ] Resident
  - [ ] Vendor
  - [ ] Society Admin
  - [ ] Platform Admin
- What data needs to be stored/retrieved?
- What are the main user flows?

### 1.2 Identify Affected Layers

```markdown
## Feature: {Feature Name}

### Affected Layers
- [ ] Database: {Yes/No - describe schema changes}
- [ ] Backend API: {Yes/No - list endpoints}
- [ ] Mobile - Resident App: {Yes/No - list screens}
- [ ] Mobile - Vendor App: {Yes/No - list screens}
- [ ] Web - Society Admin: {Yes/No - list pages}
- [ ] Web - Platform Admin: {Yes/No - list pages}

### Dependencies
- {List dependencies between layers}

### API Contract
{Define request/response formats}
```

---

## Phase 2: Database (if needed)

### 2.1 Create Migration

**Location**: `backend/migrations/{timestamp}_{description}.sql`

```sql
-- +migrate Up
-- Description: {Description}

CREATE TABLE {table_name} (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Add columns
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add indexes
CREATE INDEX idx_{table_name}_{column} ON {table_name}({column});

-- Add foreign key constraints
ALTER TABLE {table_name}
ADD CONSTRAINT fk_{table_name}_{ref} FOREIGN KEY ({column})
REFERENCES {ref_table}(id);

-- +migrate Down
DROP TABLE IF EXISTS {table_name};
```

### 2.2 Update Shared Types

**TypeScript** (`packages/shared-types/src/models/{entity}.ts`):
```typescript
export interface {Entity} {
  id: string;
  // fields
  created_at: string;
  updated_at: string;
}

export interface Create{Entity}Request {
  // request fields
}
```

---

## Phase 3: Backend API

### 3.1 Create Model

**Location**: `backend/internal/models/{entity}.go`

### 3.2 Create Repository

**Location**: `backend/internal/repositories/{entity}_repository.go`

### 3.3 Create Service

**Location**: `backend/internal/services/{entity}_service.go`

### 3.4 Create Handler

**Location**: `backend/internal/handlers/{entity}_handler.go`

### 3.5 Register Routes

**Update**: `backend/cmd/api/main.go`

### 3.6 Create Tests

- `backend/internal/services/{entity}_service_test.go`
- `backend/internal/handlers/{entity}_handler_test.go`

### 3.7 Verify Backend

```bash
cd backend
make fmt
make lint
make test
```

---

## Phase 4: Web Admin (if applicable)

### 4.1 Create API Hooks

**Location**: `apps/{app}-admin-web/src/hooks/use-{feature}.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';

export function use{Feature}() {
  return useQuery({
    queryKey: ['{feature}'],
    queryFn: async () => {
      const { data } = await apiClient.get('/api/v1/{feature}');
      return data.data;
    },
  });
}

export function useCreate{Feature}() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (request) => apiClient.post('/api/v1/{feature}', request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['{feature}'] });
    },
  });
}
```

### 4.2 Create Components

**Location**: `apps/{app}-admin-web/src/components/{feature}/`

- `{feature}-table.tsx` - List/table view
- `{feature}-form.tsx` - Create/edit form
- `{feature}-card.tsx` - Card component

### 4.3 Create Pages

**Location**: `apps/{app}-admin-web/src/app/(dashboard)/{feature}/`

- `page.tsx` - List page
- `[id]/page.tsx` - Detail page
- `new/page.tsx` - Create page

### 4.4 Create Tests

**Location**: `apps/{app}-admin-web/src/components/{feature}/__tests__/`

### 4.5 Verify Web

```bash
cd apps/{app}-admin-web
npm run type-check
npm run lint
npm test
```

---

## Phase 5: Mobile Apps (if applicable)

### 5.1 Create Model

**Location**: `apps/{app}_app/lib/features/{feature}/domain/{feature}_model.dart`

```dart
import 'package:freezed_annotation/freezed_annotation.dart';

part '{feature}_model.freezed.dart';
part '{feature}_model.g.dart';

@freezed
class {Feature} with _${Feature} {
  const factory {Feature}({
    required String id,
    // fields
  }) = _{Feature};

  factory {Feature}.fromJson(Map<String, dynamic> json) =>
      _${Feature}FromJson(json);
}
```

### 5.2 Create Repository

**Location**: `apps/{app}_app/lib/features/{feature}/data/{feature}_repository.dart`

### 5.3 Create Provider

**Location**: `apps/{app}_app/lib/features/{feature}/presentation/providers/{feature}_provider.dart`

### 5.4 Create Screens

**Location**: `apps/{app}_app/lib/features/{feature}/presentation/screens/`

### 5.5 Create Widgets

**Location**: `apps/{app}_app/lib/features/{feature}/presentation/widgets/`

### 5.6 Update Router

**Update**: `apps/{app}_app/lib/core/router/app_router.dart`

### 5.7 Create Tests

**Location**: `apps/{app}_app/test/features/{feature}/`

### 5.8 Verify Mobile

```bash
cd apps/{app}_app
dart run build_runner build --delete-conflicting-outputs
flutter analyze
flutter test
```

---

## Phase 6: Integration Testing

### 6.1 API Contract Verification

Ensure consistency across all platforms:

| Field | Backend (Go) | Web (TypeScript) | Mobile (Dart) |
|-------|--------------|------------------|---------------|
| ID | `string` | `string` | `String` |
| Date | `time.Time` → ISO8601 | `string` → `Date` | `DateTime` |
| Enum | `string` const | `string` union | `enum` |

### 6.2 Run Full Test Suite

```bash
# Backend
cd backend && make test

# Web
cd apps/society-admin-web && npm test
cd apps/platform-admin-web && npm test

# Mobile
cd apps/resident_app && flutter test
cd apps/vendor_app && flutter test
```

### 6.3 E2E Testing

Create E2E tests for critical user flows:

**Web (Playwright)**:
```typescript
test('can create {feature}', async ({ page }) => {
  await page.goto('/{feature}/new');
  await page.fill('[name="field"]', 'value');
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('/{feature}');
});
```

**Mobile (Integration)**:
```dart
testWidgets('can create {feature}', (tester) async {
  await tester.pumpWidget(const MyApp());
  await tester.tap(find.text('Create'));
  await tester.pumpAndSettle();
  // assertions
});
```

---

## Completion Checklist

### Backend
- [ ] Model created
- [ ] Repository created (interface + implementation)
- [ ] Service created with business logic
- [ ] Handler created with proper error handling
- [ ] Routes registered
- [ ] Unit tests written and passing
- [ ] Handler tests written and passing

### Web Admin
- [ ] API hooks created
- [ ] Components created
- [ ] Pages created
- [ ] Component tests written and passing
- [ ] Type checking passing
- [ ] Linting passing

### Mobile Apps
- [ ] Models created (Freezed)
- [ ] Repository created
- [ ] Providers created
- [ ] Screens created
- [ ] Widgets created
- [ ] Router updated
- [ ] Widget tests written and passing
- [ ] Code generated
- [ ] Analysis passing

### Cross-Platform
- [ ] API contract consistent
- [ ] Error codes consistent
- [ ] Date formats consistent (ISO 8601)
- [ ] ID formats consistent (UUID)
- [ ] All tests passing

---

## Post-Implementation

1. **Update Documentation**
   - API documentation
   - Feature documentation

2. **Code Review**
   - Review for security issues
   - Review for performance
   - Review for consistency

3. **Create PR**
   ```bash
   git add .
   git commit -m "feat: Add {feature}

   - Add database migration for {table}
   - Add API endpoints: GET/POST/PUT/DELETE /{feature}
   - Add web admin pages for {feature} management
   - Add mobile screens for {feature}
   - Add comprehensive tests for all layers"
   ```
