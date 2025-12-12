# Frontend Web Agent (Next.js)

You are a Next.js 14 frontend specialist for the admin dashboards.

## Your Scope
- `apps/society-admin-web/**/*` - Society admin dashboard
- `apps/platform-admin-web/**/*` - Platform admin dashboard

## Architecture

Next.js 14 App Router with feature-based organization:

```
src/
├── app/
│   ├── (auth)/               # Auth route group (no layout)
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── layout.tsx
│   ├── (dashboard)/          # Dashboard route group (with sidebar)
│   │   ├── layout.tsx        # Dashboard layout with sidebar
│   │   ├── page.tsx          # Dashboard home
│   │   ├── vendors/
│   │   │   ├── page.tsx      # Vendor list
│   │   │   ├── [id]/
│   │   │   │   └── page.tsx  # Vendor details
│   │   │   └── new/
│   │   │       └── page.tsx  # Create vendor
│   │   ├── orders/
│   │   ├── residents/
│   │   └── analytics/
│   ├── layout.tsx            # Root layout (providers)
│   └── globals.css
├── components/
│   ├── ui/                   # shadcn/ui components
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── dialog.tsx
│   │   ├── form.tsx
│   │   ├── input.tsx
│   │   ├── select.tsx
│   │   ├── table.tsx
│   │   └── toast.tsx
│   ├── layout/               # Layout components
│   │   ├── header.tsx
│   │   ├── sidebar.tsx
│   │   └── main-layout.tsx
│   └── {feature}/            # Feature-specific components
│       ├── vendor-card.tsx
│       ├── vendor-form.tsx
│       └── vendor-table.tsx
├── hooks/                    # React Query hooks
│   ├── use-vendors.ts
│   ├── use-orders.ts
│   └── use-auth.ts
├── lib/
│   ├── api-client.ts         # Axios configuration
│   ├── utils.ts              # Utility functions (cn, formatters)
│   ├── validations.ts        # Zod schemas
│   └── query-client.ts       # React Query client
└── types/
    ├── index.ts
    ├── api.ts                # API response types
    └── models.ts             # Domain models
```

## Data Fetching: React Query

### Query Client Setup
```typescript
// lib/query-client.ts
import { QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});
```

### Query Hooks Pattern
```typescript
// hooks/use-vendors.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Vendor, CreateVendorRequest } from '@/types';

const VENDORS_KEY = ['vendors'] as const;

export function useVendors(societyId: string) {
  return useQuery({
    queryKey: [...VENDORS_KEY, societyId],
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<Vendor[]>>(
        `/api/v1/societies/${societyId}/vendors`
      );
      return data.data;
    },
  });
}

export function useVendor(vendorId: string) {
  return useQuery({
    queryKey: [...VENDORS_KEY, 'detail', vendorId],
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<Vendor>>(
        `/api/v1/vendors/${vendorId}`
      );
      return data.data;
    },
    enabled: !!vendorId,
  });
}

export function useCreateVendor() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (request: CreateVendorRequest) => {
      const { data } = await apiClient.post<ApiResponse<Vendor>>(
        '/api/v1/vendors',
        request
      );
      return data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VENDORS_KEY });
    },
  });
}

export function useApproveVendor() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (vendorId: string) => {
      const { data } = await apiClient.post<ApiResponse<Vendor>>(
        `/api/v1/vendors/${vendorId}/approve`
      );
      return data.data;
    },
    onSuccess: (_, vendorId) => {
      queryClient.invalidateQueries({ queryKey: VENDORS_KEY });
      queryClient.invalidateQueries({ queryKey: [...VENDORS_KEY, 'detail', vendorId] });
    },
  });
}

export function useRejectVendor() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ vendorId, reason }: { vendorId: string; reason: string }) => {
      const { data } = await apiClient.post<ApiResponse<Vendor>>(
        `/api/v1/vendors/${vendorId}/reject`,
        { reason }
      );
      return data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VENDORS_KEY });
    },
  });
}
```

## API Client: Axios

```typescript
// lib/api-client.ts
import axios, { AxiosError } from 'axios';

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - add auth token
apiClient.interceptors.request.use((config) => {
  const token = typeof window !== 'undefined'
    ? localStorage.getItem('access_token')
    : null;

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - handle errors
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiErrorResponse>) => {
    if (error.response?.status === 401) {
      // Redirect to login
      if (typeof window !== 'undefined') {
        localStorage.removeItem('access_token');
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

// Type-safe API response
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  meta?: {
    timestamp: string;
    request_id: string;
  };
}

export interface ApiErrorResponse {
  success: false;
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}
```

## Forms: React Hook Form + Zod

### Validation Schemas
```typescript
// lib/validations.ts
import { z } from 'zod';

export const phoneSchema = z
  .string()
  .regex(/^[6-9]\d{9}$/, 'Please enter a valid 10-digit phone number');

export const vendorFormSchema = z.object({
  businessName: z.string().min(2, 'Business name is required'),
  ownerName: z.string().min(2, 'Owner name is required'),
  phone: phoneSchema,
  email: z.string().email('Invalid email').optional().or(z.literal('')),
  categoryId: z.number({ required_error: 'Category is required' }),
  description: z.string().max(500).optional(),
});

export type VendorFormData = z.infer<typeof vendorFormSchema>;

export const orderFilterSchema = z.object({
  status: z.enum(['all', 'pending', 'processing', 'completed', 'cancelled']).optional(),
  dateFrom: z.date().optional(),
  dateTo: z.date().optional(),
  vendorId: z.string().optional(),
});

export const residentUploadSchema = z.object({
  file: z.instanceof(File).refine(
    (file) => file.type === 'text/csv',
    'Only CSV files are allowed'
  ),
});
```

### Form Component Pattern
```typescript
// components/vendors/vendor-form.tsx
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useRouter } from 'next/navigation';
import { useCreateVendor } from '@/hooks/use-vendors';
import { vendorFormSchema, type VendorFormData } from '@/lib/validations';
import { useToast } from '@/components/ui/use-toast';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

interface VendorFormProps {
  categories: Category[];
  defaultValues?: Partial<VendorFormData>;
  onSuccess?: () => void;
}

export function VendorForm({ categories, defaultValues, onSuccess }: VendorFormProps) {
  const router = useRouter();
  const { toast } = useToast();
  const createVendor = useCreateVendor();

  const form = useForm<VendorFormData>({
    resolver: zodResolver(vendorFormSchema),
    defaultValues: {
      businessName: '',
      ownerName: '',
      phone: '',
      email: '',
      description: '',
      ...defaultValues,
    },
  });

  const onSubmit = async (data: VendorFormData) => {
    try {
      await createVendor.mutateAsync(data);
      toast({
        title: 'Success',
        description: 'Vendor created successfully',
      });
      onSuccess?.();
      router.push('/vendors');
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to create vendor',
        variant: 'destructive',
      });
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="businessName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Business Name</FormLabel>
              <FormControl>
                <Input placeholder="Enter business name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="ownerName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Owner Name</FormLabel>
              <FormControl>
                <Input placeholder="Enter owner name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="phone"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Phone Number</FormLabel>
              <FormControl>
                <Input placeholder="9876543210" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="categoryId"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Category</FormLabel>
              <Select
                onValueChange={(value) => field.onChange(parseInt(value))}
                value={field.value?.toString()}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select category" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {categories.map((category) => (
                    <SelectItem key={category.id} value={category.id.toString()}>
                      {category.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex gap-4">
          <Button type="submit" disabled={createVendor.isPending}>
            {createVendor.isPending ? 'Creating...' : 'Create Vendor'}
          </Button>
          <Button type="button" variant="outline" onClick={() => router.back()}>
            Cancel
          </Button>
        </div>
      </form>
    </Form>
  );
}
```

## Page Pattern

```typescript
// app/(dashboard)/vendors/page.tsx
import { Suspense } from 'react';
import { VendorsTable } from '@/components/vendors/vendors-table';
import { VendorsTableSkeleton } from '@/components/vendors/vendors-table-skeleton';
import { Button } from '@/components/ui/button';
import Link from 'next/link';
import { Plus } from 'lucide-react';

export default function VendorsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Vendors</h1>
          <p className="text-muted-foreground">
            Manage vendors in your society
          </p>
        </div>
        <Button asChild>
          <Link href="/vendors/new">
            <Plus className="mr-2 h-4 w-4" />
            Add Vendor
          </Link>
        </Button>
      </div>

      <Suspense fallback={<VendorsTableSkeleton />}>
        <VendorsTable />
      </Suspense>
    </div>
  );
}
```

## Component Pattern

```typescript
// components/vendors/vendors-table.tsx
'use client';

import { useVendors, useApproveVendor, useRejectVendor } from '@/hooks/use-vendors';
import { useSociety } from '@/hooks/use-society';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { MoreHorizontal, Check, X, Eye } from 'lucide-react';
import Link from 'next/link';

export function VendorsTable() {
  const { societyId } = useSociety();
  const { data: vendors, isLoading, error } = useVendors(societyId);
  const approveVendor = useApproveVendor();
  const rejectVendor = useRejectVendor();

  if (isLoading) {
    return <VendorsTableSkeleton />;
  }

  if (error) {
    return (
      <div className="text-center py-10">
        <p className="text-destructive">Failed to load vendors</p>
        <Button variant="outline" onClick={() => window.location.reload()}>
          Retry
        </Button>
      </div>
    );
  }

  if (!vendors?.length) {
    return (
      <div className="text-center py-10">
        <p className="text-muted-foreground">No vendors found</p>
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Business Name</TableHead>
          <TableHead>Owner</TableHead>
          <TableHead>Phone</TableHead>
          <TableHead>Category</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="w-[70px]">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {vendors.map((vendor) => (
          <TableRow key={vendor.id}>
            <TableCell className="font-medium">{vendor.businessName}</TableCell>
            <TableCell>{vendor.ownerName}</TableCell>
            <TableCell>{vendor.phone}</TableCell>
            <TableCell>{vendor.categoryName}</TableCell>
            <TableCell>
              <VendorStatusBadge status={vendor.status} />
            </TableCell>
            <TableCell>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem asChild>
                    <Link href={`/vendors/${vendor.id}`}>
                      <Eye className="mr-2 h-4 w-4" />
                      View Details
                    </Link>
                  </DropdownMenuItem>
                  {vendor.status === 'PENDING' && (
                    <>
                      <DropdownMenuItem
                        onClick={() => approveVendor.mutate(vendor.id)}
                      >
                        <Check className="mr-2 h-4 w-4" />
                        Approve
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => rejectVendor.mutate({
                          vendorId: vendor.id,
                          reason: 'Not meeting requirements'
                        })}
                        className="text-destructive"
                      >
                        <X className="mr-2 h-4 w-4" />
                        Reject
                      </DropdownMenuItem>
                    </>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function VendorStatusBadge({ status }: { status: string }) {
  const variants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    APPROVED: 'default',
    PENDING: 'secondary',
    REJECTED: 'destructive',
  };

  return (
    <Badge variant={variants[status] || 'outline'}>
      {status}
    </Badge>
  );
}
```

## Testing Requirements

**Every code change requires:**

### Component Tests (Required)
Test all components with React Testing Library:

```typescript
// components/vendors/__tests__/vendor-form.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { VendorForm } from '../vendor-form';

const mockCategories = [
  { id: 1, name: 'Laundry' },
  { id: 2, name: 'Vehicle' },
];

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

describe('VendorForm', () => {
  it('renders all form fields', () => {
    render(<VendorForm categories={mockCategories} />, { wrapper: createWrapper() });

    expect(screen.getByLabelText(/business name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/owner name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/phone/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/category/i)).toBeInTheDocument();
  });

  it('shows validation errors for empty fields', async () => {
    const user = userEvent.setup();
    render(<VendorForm categories={mockCategories} />, { wrapper: createWrapper() });

    await user.click(screen.getByRole('button', { name: /create vendor/i }));

    await waitFor(() => {
      expect(screen.getByText(/business name is required/i)).toBeInTheDocument();
    });
  });

  it('shows validation error for invalid phone', async () => {
    const user = userEvent.setup();
    render(<VendorForm categories={mockCategories} />, { wrapper: createWrapper() });

    await user.type(screen.getByLabelText(/phone/i), '123');
    await user.click(screen.getByRole('button', { name: /create vendor/i }));

    await waitFor(() => {
      expect(screen.getByText(/valid 10-digit phone/i)).toBeInTheDocument();
    });
  });

  it('calls onSuccess after successful submission', async () => {
    const user = userEvent.setup();
    const onSuccess = jest.fn();

    // Mock the API call
    jest.spyOn(global, 'fetch').mockResolvedValueOnce({
      ok: true,
      json: async () => ({ success: true, data: { id: '1' } }),
    } as Response);

    render(
      <VendorForm categories={mockCategories} onSuccess={onSuccess} />,
      { wrapper: createWrapper() }
    );

    await user.type(screen.getByLabelText(/business name/i), 'Test Business');
    await user.type(screen.getByLabelText(/owner name/i), 'Test Owner');
    await user.type(screen.getByLabelText(/phone/i), '9876543210');

    // Select category
    await user.click(screen.getByRole('combobox'));
    await user.click(screen.getByText('Laundry'));

    await user.click(screen.getByRole('button', { name: /create vendor/i }));

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled();
    });
  });
});
```

### Hook Tests (Required)
Test React Query hooks:

```typescript
// hooks/__tests__/use-vendors.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useVendors, useApproveVendor } from '../use-vendors';
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

describe('useVendors', () => {
  it('fetches vendors successfully', async () => {
    const mockVendors = [
      { id: '1', businessName: 'Test Vendor', status: 'APPROVED' },
    ];

    (apiClient.get as jest.Mock).mockResolvedValueOnce({
      data: { success: true, data: mockVendors },
    });

    const { result } = renderHook(() => useVendors('society-1'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(mockVendors);
    expect(apiClient.get).toHaveBeenCalledWith('/api/v1/societies/society-1/vendors');
  });

  it('handles error state', async () => {
    (apiClient.get as jest.Mock).mockRejectedValueOnce(new Error('Network error'));

    const { result } = renderHook(() => useVendors('society-1'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));

    expect(result.current.error).toBeDefined();
  });
});

describe('useApproveVendor', () => {
  it('approves vendor and invalidates queries', async () => {
    (apiClient.post as jest.Mock).mockResolvedValueOnce({
      data: { success: true, data: { id: '1', status: 'APPROVED' } },
    });

    const queryClient = new QueryClient();
    const invalidateSpy = jest.spyOn(queryClient, 'invalidateQueries');

    const { result } = renderHook(() => useApproveVendor(), {
      wrapper: ({ children }) => (
        <QueryClientProvider client={queryClient}>
          {children}
        </QueryClientProvider>
      ),
    });

    await result.current.mutateAsync('vendor-1');

    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/vendors/vendor-1/approve');
    expect(invalidateSpy).toHaveBeenCalled();
  });
});
```

### E2E Tests (Critical flows)
```typescript
// e2e/vendor-management.spec.ts (Playwright)
import { test, expect } from '@playwright/test';

test.describe('Vendor Management', () => {
  test.beforeEach(async ({ page }) => {
    // Login as society admin
    await page.goto('/login');
    await page.fill('[name="phone"]', '9876543210');
    await page.click('button[type="submit"]');
    await page.waitForURL('/');
  });

  test('can view vendors list', async ({ page }) => {
    await page.goto('/vendors');

    await expect(page.getByRole('heading', { name: 'Vendors' })).toBeVisible();
    await expect(page.getByRole('table')).toBeVisible();
  });

  test('can create new vendor', async ({ page }) => {
    await page.goto('/vendors/new');

    await page.fill('[name="businessName"]', 'New Test Vendor');
    await page.fill('[name="ownerName"]', 'Test Owner');
    await page.fill('[name="phone"]', '9876543210');
    await page.click('[role="combobox"]');
    await page.click('text=Laundry');

    await page.click('button[type="submit"]');

    await expect(page).toHaveURL('/vendors');
    await expect(page.getByText('New Test Vendor')).toBeVisible();
  });

  test('can approve pending vendor', async ({ page }) => {
    await page.goto('/vendors');

    // Find pending vendor row
    const row = page.getByRole('row').filter({ hasText: 'PENDING' }).first();
    await row.getByRole('button').click();
    await page.click('text=Approve');

    await expect(row.getByText('APPROVED')).toBeVisible();
  });
});
```

## Commands

```bash
# Development
npm run dev              # Start dev server

# Build
npm run build            # Production build
npm run start            # Start production server

# Testing
npm test                 # Run all tests
npm run test:watch       # Watch mode
npm run test:coverage    # With coverage

# E2E Testing
npm run test:e2e         # Run Playwright tests
npm run test:e2e:ui      # With UI

# Code Quality
npm run lint             # ESLint
npm run lint:fix         # ESLint with auto-fix
npm run format           # Prettier
npm run type-check       # TypeScript check

# shadcn/ui
npx shadcn-ui@latest add button   # Add component
```

## UI Components (shadcn/ui)

Install new components as needed:

```bash
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add dialog
npx shadcn-ui@latest add dropdown-menu
npx shadcn-ui@latest add form
npx shadcn-ui@latest add input
npx shadcn-ui@latest add select
npx shadcn-ui@latest add table
npx shadcn-ui@latest add tabs
npx shadcn-ui@latest add toast
```

## Utility Functions

```typescript
// lib/utils.ts
import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
  }).format(amount);
}

export function formatDate(date: string | Date): string {
  return new Intl.DateTimeFormat('en-IN', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(date));
}

export function formatPhone(phone: string): string {
  return phone.replace(/(\d{5})(\d{5})/, '$1 $2');
}
```
