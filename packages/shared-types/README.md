# Shared Types Package

TypeScript types shared between backend and frontend applications.

## Purpose

This package provides:
- TypeScript interfaces for all data models
- API request/response types
- Database types (generated from Supabase)

## Usage

### In Backend

```typescript
import { Order, CreateOrderDTO, ApiResponse } from '@society-platform/shared-types';

const response: ApiResponse<Order> = {
  success: true,
  data: order
};
```

### In Admin Web

```typescript
import { Order, Vendor } from '@society-platform/shared-types';

interface Props {
  order: Order;
  vendor: Vendor;
}
```

## Development

Build types:
```bash
npm run build
```

Watch mode:
```bash
npm run dev
```

## Generating Database Types

To generate types from Supabase schema:

```bash
supabase gen types typescript --project-id <your-project-id> > src/database.ts
```
