# Generate Tests

Generate tests for: $ARGUMENTS

## Process

1. **Analyze the target file** to understand:
   - What functions/methods need testing
   - Dependencies that need mocking
   - Edge cases to cover
   - Error conditions

2. **Generate comprehensive tests** including:
   - Happy path tests
   - Error/failure cases
   - Edge cases
   - Boundary conditions

---

## Go Test Generation

### Service Test Template

**Location**: `backend/internal/services/{name}_test.go`

```go
package services

import (
    "context"
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "society-service-app/backend/internal/models"
)

// Mock repository
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

// Add other mock methods...

func Test{Entity}Service_Create(t *testing.T) {
    tests := []struct {
        name      string
        input     *models.Create{Entity}Request
        mockSetup func(*Mock{Entity}Repository)
        want      *models.{Entity}
        wantErr   error
    }{
        {
            name: "creates entity successfully",
            input: &models.Create{Entity}Request{
                Name: "Test Entity",
            },
            mockSetup: func(m *Mock{Entity}Repository) {
                m.On("Create", mock.Anything, mock.AnythingOfType("*models.{Entity}")).
                    Return(&models.{Entity}{
                        ID:   "123",
                        Name: "Test Entity",
                    }, nil)
            },
            want: &models.{Entity}{
                ID:   "123",
                Name: "Test Entity",
            },
            wantErr: nil,
        },
        {
            name: "returns error when name is empty",
            input: &models.Create{Entity}Request{
                Name: "",
            },
            mockSetup: func(m *Mock{Entity}Repository) {},
            want:      nil,
            wantErr:   ErrInvalidName,
        },
        {
            name: "returns error when repository fails",
            input: &models.Create{Entity}Request{
                Name: "Test Entity",
            },
            mockSetup: func(m *Mock{Entity}Repository) {
                m.On("Create", mock.Anything, mock.Anything).
                    Return(nil, errors.New("database error"))
            },
            want:    nil,
            wantErr: errors.New("database error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(Mock{Entity}Repository)
            tt.mockSetup(mockRepo)

            svc := New{Entity}Service(mockRepo)
            got, err := svc.Create(context.Background(), tt.input)

            if tt.wantErr != nil {
                require.Error(t, err)
                if tt.wantErr.Error() != "" {
                    assert.Contains(t, err.Error(), tt.wantErr.Error())
                }
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.want.ID, got.ID)
            assert.Equal(t, tt.want.Name, got.Name)
            mockRepo.AssertExpectations(t)
        })
    }
}

func Test{Entity}Service_GetByID(t *testing.T) {
    tests := []struct {
        name      string
        id        string
        mockSetup func(*Mock{Entity}Repository)
        want      *models.{Entity}
        wantErr   error
    }{
        {
            name: "returns entity when found",
            id:   "123",
            mockSetup: func(m *Mock{Entity}Repository) {
                m.On("GetByID", mock.Anything, "123").
                    Return(&models.{Entity}{
                        ID:   "123",
                        Name: "Test",
                    }, nil)
            },
            want: &models.{Entity}{
                ID:   "123",
                Name: "Test",
            },
            wantErr: nil,
        },
        {
            name: "returns error when not found",
            id:   "nonexistent",
            mockSetup: func(m *Mock{Entity}Repository) {
                m.On("GetByID", mock.Anything, "nonexistent").
                    Return(nil, errors.New("not found"))
            },
            want:    nil,
            wantErr: Err{Entity}NotFound,
        },
        {
            name: "returns error for empty ID",
            id:   "",
            mockSetup: func(m *Mock{Entity}Repository) {},
            want:    nil,
            wantErr: ErrInvalidID,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(Mock{Entity}Repository)
            tt.mockSetup(mockRepo)

            svc := New{Entity}Service(mockRepo)
            got, err := svc.GetByID(context.Background(), tt.id)

            if tt.wantErr != nil {
                require.Error(t, err)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.want.ID, got.ID)
        })
    }
}
```

### Handler Test Template

**Location**: `backend/internal/handlers/{name}_test.go`

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

func setupRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    return gin.New()
}

func Test{Entity}Handler_Create(t *testing.T) {
    tests := []struct {
        name       string
        body       interface{}
        mockSetup  func(*MockService)
        wantStatus int
        wantBody   map[string]interface{}
    }{
        {
            name: "returns 201 on success",
            body: map[string]interface{}{
                "name": "Test Entity",
            },
            mockSetup: func(m *MockService) {
                m.On("Create", mock.Anything, mock.Anything).
                    Return(&models.{Entity}{ID: "123", Name: "Test Entity"}, nil)
            },
            wantStatus: http.StatusCreated,
            wantBody: map[string]interface{}{
                "success": true,
            },
        },
        {
            name:       "returns 400 for invalid JSON",
            body:       "invalid json",
            mockSetup:  func(m *MockService) {},
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "returns 400 for missing required fields",
            body: map[string]interface{}{},
            mockSetup:  func(m *MockService) {},
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "returns 409 on conflict",
            body: map[string]interface{}{
                "name": "Existing Entity",
            },
            mockSetup: func(m *MockService) {
                m.On("Create", mock.Anything, mock.Anything).
                    Return(nil, services.Err{Entity}Exists)
            },
            wantStatus: http.StatusConflict,
        },
        {
            name: "returns 500 on internal error",
            body: map[string]interface{}{
                "name": "Test Entity",
            },
            mockSetup: func(m *MockService) {
                m.On("Create", mock.Anything, mock.Anything).
                    Return(nil, errors.New("database error"))
            },
            wantStatus: http.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := new(MockService)
            tt.mockSetup(mockSvc)

            handler := New{Entity}Handler(mockSvc)
            router := setupRouter()
            router.POST("/{entities}", handler.Create)

            bodyBytes, _ := json.Marshal(tt.body)
            req := httptest.NewRequest("POST", "/{entities}", bytes.NewReader(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.wantStatus, w.Code)

            if tt.wantBody != nil {
                var response map[string]interface{}
                json.Unmarshal(w.Body.Bytes(), &response)
                for key, value := range tt.wantBody {
                    assert.Equal(t, value, response[key])
                }
            }
        })
    }
}

func Test{Entity}Handler_GetByID(t *testing.T) {
    tests := []struct {
        name       string
        id         string
        mockSetup  func(*MockService)
        wantStatus int
    }{
        {
            name: "returns 200 when found",
            id:   "123",
            mockSetup: func(m *MockService) {
                m.On("GetByID", mock.Anything, "123").
                    Return(&models.{Entity}{ID: "123"}, nil)
            },
            wantStatus: http.StatusOK,
        },
        {
            name: "returns 404 when not found",
            id:   "nonexistent",
            mockSetup: func(m *MockService) {
                m.On("GetByID", mock.Anything, "nonexistent").
                    Return(nil, services.Err{Entity}NotFound)
            },
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := new(MockService)
            tt.mockSetup(mockSvc)

            handler := New{Entity}Handler(mockSvc)
            router := setupRouter()
            router.GET("/{entities}/:id", handler.GetByID)

            req := httptest.NewRequest("GET", "/{entities}/"+tt.id, nil)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.wantStatus, w.Code)
        })
    }
}
```

---

## Flutter Test Generation

### Repository Test Template

**Location**: `apps/{app}_app/test/features/{feature}/data/{feature}_repository_test.dart`

```dart
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:{app}_app/core/api/api_client.dart';
import 'package:{app}_app/features/{feature}/data/{feature}_repository.dart';
import 'package:{app}_app/features/{feature}/domain/{feature}_model.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  late MockApiClient mockClient;
  late {Feature}Repository repository;

  setUp(() {
    mockClient = MockApiClient();
    repository = {Feature}Repository(mockClient);
  });

  group('{Feature}Repository', () {
    group('getAll', () {
      test('returns list on success', () async {
        // Arrange
        when(() => mockClient.get(any())).thenAnswer(
          (_) async => Response(
            requestOptions: RequestOptions(path: ''),
            statusCode: 200,
            data: {
              'success': true,
              'data': [
                {'id': '1', 'name': 'Test 1'},
                {'id': '2', 'name': 'Test 2'},
              ],
            },
          ),
        );

        // Act
        final result = await repository.getAll();

        // Assert
        expect(result, isA<List<{Model}>>());
        expect(result.length, 2);
        verify(() => mockClient.get('/api/v1/{feature}')).called(1);
      });

      test('throws exception on API error', () async {
        when(() => mockClient.get(any())).thenThrow(
          DioException(
            requestOptions: RequestOptions(path: ''),
            response: Response(
              requestOptions: RequestOptions(path: ''),
              statusCode: 500,
            ),
          ),
        );

        expect(() => repository.getAll(), throwsA(isA<DioException>()));
      });
    });

    group('getById', () {
      test('returns item on success', () async {
        when(() => mockClient.get(any())).thenAnswer(
          (_) async => Response(
            requestOptions: RequestOptions(path: ''),
            statusCode: 200,
            data: {
              'success': true,
              'data': {'id': '1', 'name': 'Test'},
            },
          ),
        );

        final result = await repository.getById('1');

        expect(result.id, '1');
        expect(result.name, 'Test');
      });

      test('throws on 404', () async {
        when(() => mockClient.get(any())).thenThrow(
          DioException(
            requestOptions: RequestOptions(path: ''),
            response: Response(
              requestOptions: RequestOptions(path: ''),
              statusCode: 404,
            ),
          ),
        );

        expect(() => repository.getById('999'), throwsA(isA<DioException>()));
      });
    });

    group('create', () {
      test('returns created item on success', () async {
        final request = Create{Model}Request(name: 'New Item');

        when(() => mockClient.post(any(), data: any(named: 'data'))).thenAnswer(
          (_) async => Response(
            requestOptions: RequestOptions(path: ''),
            statusCode: 201,
            data: {
              'success': true,
              'data': {'id': '123', 'name': 'New Item'},
            },
          ),
        );

        final result = await repository.create(request);

        expect(result.id, '123');
        expect(result.name, 'New Item');
      });
    });
  });
}
```

### Widget Test Template

**Location**: `apps/{app}_app/test/features/{feature}/presentation/widgets/{widget}_test.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:{app}_app/features/{feature}/presentation/widgets/{widget}.dart';

void main() {
  group('{Widget}', () {
    testWidgets('renders correctly with data', (tester) async {
      final testData = {Model}(
        id: '1',
        name: 'Test Item',
      );

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Scaffold(
              body: {Widget}(data: testData),
            ),
          ),
        ),
      );

      expect(find.text('Test Item'), findsOneWidget);
    });

    testWidgets('calls onTap when tapped', (tester) async {
      var tapped = false;
      final testData = {Model}(id: '1', name: 'Test');

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Scaffold(
              body: {Widget}(
                data: testData,
                onTap: () => tapped = true,
              ),
            ),
          ),
        ),
      );

      await tester.tap(find.byType({Widget}));
      await tester.pumpAndSettle();

      expect(tapped, isTrue);
    });

    testWidgets('shows loading state', (tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: Scaffold(
              body: {Widget}.loading(),
            ),
          ),
        ),
      );

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('shows error state', (tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: Scaffold(
              body: {Widget}.error(message: 'Error occurred'),
            ),
          ),
        ),
      );

      expect(find.text('Error occurred'), findsOneWidget);
    });
  });
}
```

---

## React Test Generation

### Component Test Template

**Location**: `apps/{app}-admin-web/src/components/{feature}/__tests__/{component}.test.tsx`

```typescript
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { {Component} } from '../{component}';

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('{Component}', () => {
  describe('rendering', () => {
    it('renders correctly with data', () => {
      const testData = {
        id: '1',
        name: 'Test Item',
      };

      render(<{Component} data={testData} />, { wrapper: createWrapper() });

      expect(screen.getByText('Test Item')).toBeInTheDocument();
    });

    it('renders empty state when no data', () => {
      render(<{Component} data={[]} />, { wrapper: createWrapper() });

      expect(screen.getByText(/no items/i)).toBeInTheDocument();
    });

    it('renders loading state', () => {
      render(<{Component} isLoading />, { wrapper: createWrapper() });

      expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('renders error state', () => {
      render(<{Component} error="Something went wrong" />, { wrapper: createWrapper() });

      expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    });
  });

  describe('interactions', () => {
    it('calls onClick when clicked', async () => {
      const user = userEvent.setup();
      const onClick = jest.fn();
      const testData = { id: '1', name: 'Test' };

      render(<{Component} data={testData} onClick={onClick} />, { wrapper: createWrapper() });

      await user.click(screen.getByRole('button'));

      expect(onClick).toHaveBeenCalledWith('1');
    });

    it('calls onSubmit with form data', async () => {
      const user = userEvent.setup();
      const onSubmit = jest.fn();

      render(<{Component} onSubmit={onSubmit} />, { wrapper: createWrapper() });

      await user.type(screen.getByLabelText(/name/i), 'Test Name');
      await user.click(screen.getByRole('button', { name: /submit/i }));

      await waitFor(() => {
        expect(onSubmit).toHaveBeenCalledWith(
          expect.objectContaining({ name: 'Test Name' })
        );
      });
    });
  });

  describe('validation', () => {
    it('shows error for required field', async () => {
      const user = userEvent.setup();

      render(<{Component} />, { wrapper: createWrapper() });

      await user.click(screen.getByRole('button', { name: /submit/i }));

      await waitFor(() => {
        expect(screen.getByText(/required/i)).toBeInTheDocument();
      });
    });
  });
});
```

### Hook Test Template

**Location**: `apps/{app}-admin-web/src/hooks/__tests__/use-{feature}.test.ts`

```typescript
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { use{Feature}, useCreate{Feature} } from '../use-{feature}';
import { apiClient } from '@/lib/api-client';

jest.mock('@/lib/api-client');

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('use{Feature}', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('fetches data successfully', async () => {
    const mockData = [{ id: '1', name: 'Test' }];

    (apiClient.get as jest.Mock).mockResolvedValueOnce({
      data: { success: true, data: mockData },
    });

    const { result } = renderHook(() => use{Feature}(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(mockData);
  });

  it('handles error state', async () => {
    (apiClient.get as jest.Mock).mockRejectedValueOnce(new Error('Network error'));

    const { result } = renderHook(() => use{Feature}(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));

    expect(result.current.error).toBeDefined();
  });
});

describe('useCreate{Feature}', () => {
  it('creates item successfully', async () => {
    const newItem = { id: '123', name: 'New Item' };

    (apiClient.post as jest.Mock).mockResolvedValueOnce({
      data: { success: true, data: newItem },
    });

    const { result } = renderHook(() => useCreate{Feature}(), {
      wrapper: createWrapper(),
    });

    await result.current.mutateAsync({ name: 'New Item' });

    expect(apiClient.post).toHaveBeenCalledWith(
      '/api/v1/{feature}',
      { name: 'New Item' }
    );
  });
});
```

---

## Run Tests After Generation

### Go
```bash
cd backend && make test
```

### Flutter
```bash
cd apps/{app}_app && flutter test
```

### React
```bash
cd apps/{app}-admin-web && npm test
```
