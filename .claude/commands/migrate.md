# Database Migration

Create migration for: $ARGUMENTS

## Process

### Step 1: Generate Migration File

**Location**: `backend/migrations/{timestamp}_{description}.sql`

**Naming Convention**: `YYYYMMDDHHMMSS_{description}.sql`

Example: `20250115120000_create_ratings_table.sql`

### Step 2: Write Migration

#### Creating a New Table

```sql
-- +migrate Up
-- Description: Create {table_name} table for {purpose}

CREATE TABLE {table_name} (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign keys
    {foreign_key}_id UUID NOT NULL REFERENCES {ref_table}(id) ON DELETE CASCADE,

    -- Data columns
    {column_name} {TYPE} {CONSTRAINTS},

    -- Common columns
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE  -- For soft deletes
);

-- Indexes for common queries
CREATE INDEX idx_{table_name}_{foreign_key}_id ON {table_name}({foreign_key}_id);
CREATE INDEX idx_{table_name}_status ON {table_name}(status);
CREATE INDEX idx_{table_name}_created_at ON {table_name}(created_at DESC);

-- Partial index for active records only
CREATE INDEX idx_{table_name}_active ON {table_name}(id)
WHERE deleted_at IS NULL;

-- Unique constraint (if needed)
CREATE UNIQUE INDEX idx_{table_name}_unique_{columns}
ON {table_name}({column1}, {column2})
WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE TRIGGER update_{table_name}_updated_at
    BEFORE UPDATE ON {table_name}
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_{table_name}_updated_at ON {table_name};
DROP TABLE IF EXISTS {table_name};
```

#### Adding Columns to Existing Table

```sql
-- +migrate Up
-- Description: Add {column} to {table_name}

ALTER TABLE {table_name}
ADD COLUMN {column_name} {TYPE} {CONSTRAINTS};

-- Add index if needed
CREATE INDEX idx_{table_name}_{column_name} ON {table_name}({column_name});

-- +migrate Down
DROP INDEX IF EXISTS idx_{table_name}_{column_name};
ALTER TABLE {table_name}
DROP COLUMN IF EXISTS {column_name};
```

#### Adding Foreign Key

```sql
-- +migrate Up
-- Description: Add foreign key from {table} to {ref_table}

ALTER TABLE {table_name}
ADD COLUMN {ref_table}_id UUID REFERENCES {ref_table}(id);

CREATE INDEX idx_{table_name}_{ref_table}_id ON {table_name}({ref_table}_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_{table_name}_{ref_table}_id;
ALTER TABLE {table_name}
DROP COLUMN IF EXISTS {ref_table}_id;
```

#### Creating Enum Type

```sql
-- +migrate Up
-- Description: Create {enum_name} type

CREATE TYPE {enum_name} AS ENUM (
    'VALUE_ONE',
    'VALUE_TWO',
    'VALUE_THREE'
);

ALTER TABLE {table_name}
ADD COLUMN {column_name} {enum_name} NOT NULL DEFAULT 'VALUE_ONE';

-- +migrate Down
ALTER TABLE {table_name}
DROP COLUMN IF EXISTS {column_name};

DROP TYPE IF EXISTS {enum_name};
```

#### Adding Check Constraint

```sql
-- +migrate Up
-- Description: Add check constraint for {purpose}

ALTER TABLE {table_name}
ADD CONSTRAINT chk_{table_name}_{constraint_name}
CHECK ({column} >= 0 AND {column} <= 100);

-- +migrate Down
ALTER TABLE {table_name}
DROP CONSTRAINT IF EXISTS chk_{table_name}_{constraint_name};
```

#### Creating Junction Table (Many-to-Many)

```sql
-- +migrate Up
-- Description: Create junction table for {table1} and {table2}

CREATE TABLE {table1}_{table2} (
    {table1}_id UUID NOT NULL REFERENCES {table1}(id) ON DELETE CASCADE,
    {table2}_id UUID NOT NULL REFERENCES {table2}(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY ({table1}_id, {table2}_id)
);

CREATE INDEX idx_{table1}_{table2}_{table1}_id ON {table1}_{table2}({table1}_id);
CREATE INDEX idx_{table1}_{table2}_{table2}_id ON {table1}_{table2}({table2}_id);

-- +migrate Down
DROP TABLE IF EXISTS {table1}_{table2};
```

### Step 3: Common Patterns

#### Soft Delete Support

```sql
-- Add to table creation
deleted_at TIMESTAMP WITH TIME ZONE,

-- Partial index for active records
CREATE INDEX idx_{table}_active ON {table}(id) WHERE deleted_at IS NULL;

-- Query pattern
SELECT * FROM {table} WHERE deleted_at IS NULL;
```

#### Audit Columns

```sql
-- Add to table creation
created_by UUID REFERENCES users(id),
updated_by UUID REFERENCES users(id),
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
```

#### ltree for Hierarchy (Society Structure)

```sql
-- Enable extension (one-time)
CREATE EXTENSION IF NOT EXISTS ltree;

-- Add path column
path LTREE NOT NULL,

-- Indexes for ltree
CREATE INDEX idx_{table}_path_gist ON {table} USING GIST (path);
CREATE INDEX idx_{table}_path_btree ON {table} USING BTREE (path);

-- Query examples:
-- Get all children: WHERE path <@ 'parent.path'
-- Get ancestors: WHERE path @> 'child.path'
```

#### JSON/JSONB for Flexible Data

```sql
-- Add JSONB column
metadata JSONB DEFAULT '{}',

-- Index for JSONB queries
CREATE INDEX idx_{table}_metadata ON {table} USING GIN (metadata);

-- Query examples:
-- WHERE metadata->>'key' = 'value'
-- WHERE metadata @> '{"key": "value"}'
```

### Step 4: Data Type Reference

| Type | Use Case | Example |
|------|----------|---------|
| `UUID` | Primary/foreign keys | `id UUID PRIMARY KEY` |
| `VARCHAR(n)` | Short strings | `name VARCHAR(100)` |
| `TEXT` | Long strings | `description TEXT` |
| `INTEGER` | Whole numbers | `quantity INTEGER` |
| `DECIMAL(p,s)` | Money/precise | `price DECIMAL(10,2)` |
| `BOOLEAN` | True/false | `is_active BOOLEAN` |
| `TIMESTAMP WITH TIME ZONE` | Dates/times | `created_at TIMESTAMPTZ` |
| `JSONB` | Flexible data | `metadata JSONB` |
| `LTREE` | Hierarchies | `path LTREE` |
| `ENUM` | Fixed options | `status order_status` |

### Step 5: Apply Migration

```bash
# Apply pending migrations
cd backend && make migrate-up

# Rollback last migration
cd backend && make migrate-down

# Check migration status
cd backend && make migrate-status
```

### Step 6: Generate Types (Optional)

After migration, regenerate types for frontend:

```bash
# Generate TypeScript types from database
./scripts/generate-types.sh

# Or if using SQLC
cd backend && make sqlc
```

### Step 7: Update Application Code

After migration:

1. **Update Go models** in `backend/internal/models/`
2. **Update TypeScript types** in `packages/shared-types/`
3. **Update Dart models** in `apps/*/lib/features/*/domain/`

---

## Example Migrations

### Ratings Table

```sql
-- 20250115120000_create_ratings_table.sql

-- +migrate Up
-- Description: Create ratings table for order ratings

CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    vendor_id UUID NOT NULL REFERENCES vendors(id),
    resident_id UUID NOT NULL REFERENCES residents(id),
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Each order can only be rated once
CREATE UNIQUE INDEX idx_ratings_order_id ON ratings(order_id);

-- For vendor rating aggregation
CREATE INDEX idx_ratings_vendor_id ON ratings(vendor_id);

-- For resident rating history
CREATE INDEX idx_ratings_resident_id ON ratings(resident_id);

-- +migrate Down
DROP TABLE IF EXISTS ratings;
```

### Add Status to Orders

```sql
-- 20250115130000_add_workflow_status_to_orders.sql

-- +migrate Up
-- Description: Add workflow tracking columns to orders

ALTER TABLE orders
ADD COLUMN current_workflow_step INTEGER DEFAULT 0,
ADD COLUMN workflow_started_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN workflow_completed_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX idx_orders_workflow_status ON orders(status, current_workflow_step);

-- +migrate Down
ALTER TABLE orders
DROP COLUMN IF EXISTS current_workflow_step,
DROP COLUMN IF EXISTS workflow_started_at,
DROP COLUMN IF EXISTS workflow_completed_at;

DROP INDEX IF EXISTS idx_orders_workflow_status;
```
