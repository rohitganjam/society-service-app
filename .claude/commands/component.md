# Scaffold Component

Create a new component: $ARGUMENTS

## Argument Format

Parse the argument to determine platform and details:

- `flutter:resident:{feature}:{name}` - Flutter widget for resident app
- `flutter:vendor:{feature}:{name}` - Flutter widget for vendor app
- `web:society:{feature}:{name}` - React component for society admin
- `web:platform:{feature}:{name}` - React component for platform admin

Examples:
- `flutter:resident:orders:OrderCard`
- `web:society:vendors:VendorTable`

---

## Flutter Component

### Widget File

**Location**: `apps/{app}_app/lib/features/{feature}/presentation/widgets/{snake_case_name}.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class {PascalCaseName} extends ConsumerWidget {
  const {PascalCaseName}({
    super.key,
    // Add required parameters
  });

  // Add final fields for parameters

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Container(
      // TODO: Implement widget
      child: const Placeholder(),
    );
  }
}
```

### Screen File (if creating a screen)

**Location**: `apps/{app}_app/lib/features/{feature}/presentation/screens/{snake_case_name}_screen.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

class {PascalCaseName}Screen extends ConsumerWidget {
  const {PascalCaseName}Screen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Watch relevant providers
    // final dataAsync = ref.watch(dataProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('{Title}'),
      ),
      body: const Center(
        child: Text('TODO: Implement screen'),
      ),
    );
  }
}
```

### Widget Test File

**Location**: `apps/{app}_app/test/features/{feature}/widgets/{snake_case_name}_test.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:{app}_app/features/{feature}/presentation/widgets/{snake_case_name}.dart';

void main() {
  group('{PascalCaseName}', () {
    testWidgets('renders correctly', (tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: Scaffold(
              body: {PascalCaseName}(
                // Add required parameters
              ),
            ),
          ),
        ),
      );

      // Add assertions
      expect(find.byType({PascalCaseName}), findsOneWidget);
    });

    testWidgets('handles tap interaction', (tester) async {
      var tapped = false;

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Scaffold(
              body: {PascalCaseName}(
                onTap: () => tapped = true,
              ),
            ),
          ),
        ),
      );

      await tester.tap(find.byType({PascalCaseName}));
      await tester.pumpAndSettle();

      expect(tapped, isTrue);
    });

    testWidgets('displays data correctly', (tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: Scaffold(
              body: {PascalCaseName}(
                // Add test data
              ),
            ),
          ),
        ),
      );

      // Verify displayed data
      expect(find.text('Expected Text'), findsOneWidget);
    });
  });
}
```

### Provider (if needed)

**Location**: `apps/{app}_app/lib/features/{feature}/presentation/providers/{snake_case_name}_provider.dart`

```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:{app}_app/features/{feature}/data/{feature}_repository.dart';
import 'package:{app}_app/features/{feature}/domain/{feature}_model.dart';

part '{snake_case_name}_provider.g.dart';

@riverpod
class {PascalCaseName}Notifier extends _${PascalCaseName}Notifier {
  @override
  {State}State build() => const {State}State.initial();

  Future<void> load() async {
    state = const {State}State.loading();
    try {
      final repository = ref.read({feature}RepositoryProvider);
      final data = await repository.getData();
      state = {State}State.success(data);
    } catch (e) {
      state = {State}State.error(e.toString());
    }
  }
}
```

---

## React Component (Next.js)

### Component File

**Location**: `apps/{app}-admin-web/src/components/{feature}/{kebab-case-name}.tsx`

```typescript
'use client';

import { FC } from 'react';
import { cn } from '@/lib/utils';

interface {PascalCaseName}Props {
  // Define props
  className?: string;
}

export const {PascalCaseName}: FC<{PascalCaseName}Props> = ({
  className,
  ...props
}) => {
  return (
    <div className={cn('', className)} {...props}>
      {/* TODO: Implement component */}
    </div>
  );
};
```

### Page Component (if creating a page)

**Location**: `apps/{app}-admin-web/src/app/(dashboard)/{feature}/page.tsx`

```typescript
import { Suspense } from 'react';
import { {PascalCaseName} } from '@/components/{feature}/{kebab-case-name}';
import { {PascalCaseName}Skeleton } from '@/components/{feature}/{kebab-case-name}-skeleton';

export default function {PascalCaseName}Page() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{Title}</h1>
          <p className="text-muted-foreground">
            {Description}
          </p>
        </div>
      </div>

      <Suspense fallback={<{PascalCaseName}Skeleton />}>
        <{PascalCaseName} />
      </Suspense>
    </div>
  );
}
```

### Component Test File

**Location**: `apps/{app}-admin-web/src/components/{feature}/__tests__/{kebab-case-name}.test.tsx`

```typescript
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { {PascalCaseName} } from '../{kebab-case-name}';

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

describe('{PascalCaseName}', () => {
  it('renders correctly', () => {
    render(<{PascalCaseName} />, { wrapper: createWrapper() });

    // Add assertions
    expect(screen.getByRole('...')).toBeInTheDocument();
  });

  it('handles user interaction', async () => {
    const user = userEvent.setup();
    const onAction = jest.fn();

    render(<{PascalCaseName} onAction={onAction} />, { wrapper: createWrapper() });

    await user.click(screen.getByRole('button'));

    expect(onAction).toHaveBeenCalled();
  });

  it('displays data correctly', () => {
    const testData = {
      // Test data
    };

    render(<{PascalCaseName} data={testData} />, { wrapper: createWrapper() });

    expect(screen.getByText('Expected Text')).toBeInTheDocument();
  });

  it('shows loading state', () => {
    render(<{PascalCaseName} isLoading />, { wrapper: createWrapper() });

    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });

  it('shows error state', () => {
    render(<{PascalCaseName} error="Error message" />, { wrapper: createWrapper() });

    expect(screen.getByText('Error message')).toBeInTheDocument();
  });
});
```

### Hook (if data fetching needed)

**Location**: `apps/{app}-admin-web/src/hooks/use-{kebab-case-name}.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { {Type} } from '@/types';

const {UPPER_CASE}_KEY = ['{feature}'] as const;

export function use{PascalCaseName}() {
  return useQuery({
    queryKey: {UPPER_CASE}_KEY,
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<{Type}[]>>(
        '/api/v1/{feature}'
      );
      return data.data;
    },
  });
}

export function useCreate{PascalCaseName}() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (request: Create{Type}Request) => {
      const { data } = await apiClient.post<ApiResponse<{Type}>>(
        '/api/v1/{feature}',
        request
      );
      return data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: {UPPER_CASE}_KEY });
    },
  });
}
```

---

## After Creating

### Flutter

1. Run code generation:
   ```bash
   cd apps/{app}_app
   dart run build_runner build --delete-conflicting-outputs
   ```

2. Run tests:
   ```bash
   flutter test test/features/{feature}/
   ```

3. Check analysis:
   ```bash
   flutter analyze
   ```

### React

1. Run type check:
   ```bash
   cd apps/{app}-admin-web
   npm run type-check
   ```

2. Run tests:
   ```bash
   npm test -- --testPathPattern="{feature}"
   ```

3. Run lint:
   ```bash
   npm run lint
   ```
