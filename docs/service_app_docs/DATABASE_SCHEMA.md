# Database Schema Documentation

**Version:** 3.0
**Date:** November 17, 2025
**Database:** PostgreSQL 15 (Supabase)
**Architecture:** Multi-category service platform with hierarchical structure and custom workflows

---

## Table of Contents

1. [Overview](#1-overview)
2. [Schema Design Principles](#2-schema-design-principles)
3. [Entity Relationship Diagram](#3-entity-relationship-diagram)
4. [Core Tables](#4-core-tables)
5. [Category & Service Tables](#5-category--service-tables)
6. [Workflow Configuration Tables](#6-workflow-configuration-tables)
7. [Order Management Tables](#7-order-management-tables)
8. [Payment & Settlement Tables](#8-payment--settlement-tables)
9. [Subscription & Billing Tables](#9-subscription--billing-tables)
10. [Support & Communication Tables](#10-support--communication-tables)
11. [Indexes & Performance](#11-indexes--performance)
12. [Database Functions](#12-database-functions)
13. [Row Level Security (RLS)](#13-row-level-security-rls)
14. [Sample Data](#14-sample-data)

---

## 1. Overview

### 1.1 Database Summary

**Total Tables:** 30+ core tables
**Database Size (estimated):**
- 100 societies: ~2GB
- 500 societies: ~8GB
- 1000 societies: ~15GB

**Key Features:**
- ✅ Hierarchical category structure (multi-service support)
- ✅ Mixed-service orders (multiple service types per order)
- ✅ **Custom workflow configuration per service type**
- ✅ **Workflow step tracking for each order service**
- ✅ Society subscription billing
- ✅ Multi-tenancy with data isolation
- ✅ Audit trails for critical operations
- ✅ Soft deletes where applicable

### 1.2 Database Extensions

```sql
-- Enable required PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";      -- UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";       -- Encryption
CREATE EXTENSION IF NOT EXISTS "pg_trgm";        -- Full-text search
CREATE EXTENSION IF NOT EXISTS "btree_gin";      -- GIN indexes
```

---

## 2. Schema Design Principles

### 2.1 Naming Conventions

**Tables:**
- Plural nouns: `users`, `orders`, `societies`, `vendors`
- Snake_case: `order_items`, `rate_card_items`, `workflow_steps`
- Clear and descriptive
- **Vendor** terminology (not "laundry_person" or "LP")

**Columns:**
- Snake_case: `created_at`, `order_id`, `service_type`
- ID columns: `{table_name}_id` (e.g., `order_id`, `vendor_id`)
- Boolean columns: `is_active`, `has_multiple_services`
- Timestamp columns: `{action}_at` (e.g., `created_at`, `delivered_at`)

**Constraints:**
- Primary keys: `pk_{table_name}`
- Foreign keys: `fk_{table_name}_{referenced_table}`
- Unique constraints: `unique_{table_name}_{columns}`
- Indexes: `idx_{table_name}_{columns}`

### 2.2 Data Types

**IDs:**
- Primary keys: `SERIAL` for auto-increment integers, `UUID` for distributed records
- Foreign keys: Match referenced column type

**Text:**
- Short strings (< 50 chars): `VARCHAR(n)`
- Medium strings (50-255): `VARCHAR(255)`
- Long text: `TEXT`

**Numbers:**
- Currency: `DECIMAL(10,2)` (supports up to ₹99,999,999.99)
- Quantities: `INTEGER`
- Percentages: `DECIMAL(5,2)` (supports 0.00 to 999.99%)

**Dates/Times:**
- Timestamps: `TIMESTAMP WITH TIME ZONE` (stores UTC)
- Dates only: `DATE`

**Enums:**
- Use `CHECK` constraints for better flexibility
- Alternatively, reference tables for dynamic values

### 2.3 Audit & Metadata

**Standard columns on most tables:**
```sql
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP DEFAULT NOW(),
created_by UUID REFERENCES users(user_id),  -- Optional
deleted_at TIMESTAMP                        -- Soft delete
```

---

## 3. Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         CATEGORY & WORKFLOW STRUCTURE                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────────┐           ┌────────────────────────────────────┐ │
│  │ parent_categories│ 1     ∞   │ service_categories                 │ │
│  │ ─────────────────│◄──────────┤ ──────────────────────────────────│ │
│  │ • category_id (PK)           │ • service_id (PK)                  │ │
│  │ • category_key   │           │ • parent_category_id (FK)          │ │
│  │ • category_name  │           │ • service_key                      │ │
│  │ • is_live        │           │ • service_name                     │ │
│  └──────────────────┘           │ • default_turnaround_hours         │ │
│                                  │ • pricing_model                    │ │
│                                  └────────────┬───────────────────────┘ │
│                                               │                          │
│                                               │ 1                        │
│                                               │                          │
│                                               ▼ ∞                        │
│                                  ┌────────────────────────────────────┐ │
│                                  │ service_workflow_templates         │ │
│                                  │ ──────────────────────────────────│ │
│                                  │ • template_id (PK)                 │ │
│                                  │ • service_id (FK)                  │ │
│                                  │ • template_name                    │ │
│                                  └────────────┬───────────────────────┘ │
│                                               │                          │
│                                               │ 1                        │
│                                               │                          │
│                                               ▼ ∞                        │
│                                  ┌────────────────────────────────────┐ │
│                                  │ workflow_steps                     │ │
│                                  │ ──────────────────────────────────│ │
│                                  │ • step_id (PK)                     │ │
│                                  │ • template_id (FK)                 │ │
│                                  │ • step_name                        │ │
│                                  │ • step_order                       │ │
│                                  │ • is_required                      │ │
│                                  │ • requires_photo                   │ │
│                                  └────────────────────────────────────┘ │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                         USER MANAGEMENT                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────┐         ┌──────────────┐         ┌──────────────┐   │
│  │   users      │         │  residents   │         │   vendors    │   │
│  │ ─────────────│         │ ─────────────│         │ ─────────────│   │
│  │ • user_id (PK)    ∞   │ • resident_id│    ∞   │ • vendor_id  │   │
│  │ • phone      │◄───1────┤   (PK, FK)   │◄───1────┤   (PK, FK)   │   │
│  │ • user_type  │         │ • flat_no    │         │ • business_  │   │
│  │ • is_verified│         └──────────────┘         │   name       │   │
│  └──────────────┘                                   │ • store_addr │   │
│         │                                           └──────────────┘   │
│         │                                                    │          │
│         │                                                    │          │
│         │                                                    ▼          │
│         │                                           ┌──────────────┐   │
│         │                                           │ vendor_      │   │
│         │                                           │  services    │   │
│         │                                           │ ─────────────│   │
│         │                                           │ • vendor_id  │   │
│         │                                           │ • service_id │   │
│         │                                           │   (FK)       │   │
│         │                                           └──────────────┘   │
│         │                                                    │          │
│         │                                                    │          │
│         ▼                                                    ▼          │
│  ┌──────────────┐                                  ┌──────────────┐   │
│  │  societies   │                                  │ rate_cards   │   │
│  │ ─────────────│                                  │ ─────────────│   │
│  │ • society_id │    1                        ∞   │ • rate_card_id   │
│  │   (PK)       │◄─────────────────────────────────┤ • vendor_id  │   │
│  │ • name       │                                  │ • society_id │   │
│  │ • address    │                                  └──────────────┘   │
│  └──────────────┘                                           │          │
│                                                              │          │
│                                                              ▼          │
│                                                     ┌──────────────┐   │
│                                                     │ rate_card_   │   │
│                                                     │   items      │   │
│                                                     │ ─────────────│   │
│                                                     │ • item_id    │   │
│                                                     │ • service_id │   │
│                                                     │ • item_name  │   │
│                                                     │ • price      │   │
│                                                     └──────────────┘   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                         ORDER & WORKFLOW TRACKING                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────┐                                                       │
│  │   orders     │                                                       │
│  │ ─────────────│                                                       │
│  │ • order_id   │                                                       │
│  │   (PK)       │                                                       │
│  │ • resident_id│                                                       │
│  │ • vendor_id  │                                                       │
│  │ • society_id │                                                       │
│  │ • status     │                                                       │
│  │ • has_multiple                                                       │
│  │   _services  │                                                       │
│  └──────┬───────┘                                                       │
│         │                                                                │
│         ├──────────────────┬────────────────────┐                      │
│         │                  │                    │                      │
│         ▼                  ▼                    ▼                      │
│  ┌──────────────┐  ┌──────────────┐   ┌──────────────────┐           │
│  │ order_items  │  │ order_service│   │ order_status_log │           │
│  │ ─────────────│  │   _status    │   │ ─────────────────│           │
│  │ • order_id   │  │ ─────────────│   │ • order_id       │           │
│  │ • service_id │  │ • order_id   │   │ • status         │           │
│  │ • item_name  │  │ • service_id │   │ • changed_at     │           │
│  │ • quantity   │  │ • status     │   │ • changed_by     │           │
│  │ • unit_price │  │ • item_count │   └──────────────────┘           │
│  └──────────────┘  │ • current_   │                                    │
│                     │   step_id    │                                    │
│                     └──────┬───────┘                                    │
│                            │                                            │
│                            │ 1                                          │
│                            │                                            │
│                            ▼ ∞                                          │
│                     ┌──────────────────┐                               │
│                     │ order_workflow_  │                               │
│                     │   progress       │                               │
│                     │ ─────────────────│                               │
│                     │ • order_id       │                               │
│                     │ • service_id     │                               │
│                     │ • step_id        │                               │
│                     │ • status         │                               │
│                     │ • completed_at   │                               │
│                     │ • photos         │                               │
│                     └──────────────────┘                               │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                      PAYMENTS & SETTLEMENTS                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────┐         ┌──────────────┐                             │
│  │   orders     │    1    │  payments    │                             │
│  │ ─────────────│◄────∞───┤ ─────────────│                             │
│  │ • order_id   │         │ • payment_id │                             │
│  │   (PK)       │         │   (PK)       │                             │
│  └──────────────┘         │ • order_id   │                             │
│                            │   (FK)       │                             │
│  ┌──────────────┐         │ • amount     │                             │
│  │   vendors    │         │ • status     │                             │
│  │ ─────────────│         │ • payment_   │                             │
│  │ • vendor_id  │         │   method     │                             │
│  │   (PK)       │         │ • razorpay_  │                             │
│  └──────┬───────┘         │   order_id   │                             │
│         │                 └──────────────┘                             │
│         │ 1                                                             │
│         │                                                               │
│         ▼ ∞                                                             │
│  ┌──────────────┐                                                       │
│  │ settlements  │                                                       │
│  │ ─────────────│                                                       │
│  │ • settlement │                                                       │
│  │   _id (PK)   │                                                       │
│  │ • vendor_id  │                                                       │
│  │   (FK)       │                                                       │
│  │ • period_    │                                                       │
│  │   start      │                                                       │
│  │ • period_end │                                                       │
│  │ • gross_amt  │                                                       │
│  │ • net_amount │                                                       │
│  │ • status     │                                                       │
│  └──────────────┘                                                       │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                    SUBSCRIPTIONS & BILLING                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────────┐                                                   │
│  │   societies      │                                                   │
│  │ ─────────────────│                                                   │
│  │ • society_id (PK)│                                                   │
│  └──────┬───────────┘                                                   │
│         │ 1                                                             │
│         │                                                               │
│         ▼ 1                                                             │
│  ┌──────────────────────┐                                               │
│  │ society_subscriptions│                                               │
│  │ ─────────────────────│                                               │
│  │ • subscription_id    │                                               │
│  │   (PK)               │                                               │
│  │ • society_id (FK)    │                                               │
│  │ • tier               │                                               │
│  │ • status             │                                               │
│  │ • monthly_fee        │                                               │
│  │ • next_billing_date  │                                               │
│  └──────┬───────────────┘                                               │
│         │ 1                                                             │
│         │                                                               │
│         ▼ ∞                                                             │
│  ┌──────────────────────┐                                               │
│  │ subscription_invoices│                                               │
│  │ ─────────────────────│                                               │
│  │ • invoice_id (PK)    │                                               │
│  │ • subscription_id    │                                               │
│  │   (FK)               │                                               │
│  │ • society_id (FK)    │                                               │
│  │ • amount             │                                               │
│  │ • status             │                                               │
│  │ • due_date           │                                               │
│  │ • razorpay_order_id  │                                               │
│  └──────────────────────┘                                               │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                    SUPPORT & COMMUNICATION                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────┐         ┌──────────────┐                             │
│  │   orders     │    1    │  disputes    │                             │
│  │ ─────────────│◄────∞───┤ ─────────────│                             │
│  │ • order_id   │         │ • dispute_id │                             │
│  │   (PK)       │         │   (PK)       │                             │
│  └──────────────┘         │ • order_id   │                             │
│                            │   (FK)       │                             │
│  ┌──────────────┐         │ • service_id │                             │
│  │   orders     │         │ • issue_type │                             │
│  │ ─────────────│         │ • status     │                             │
│  │ • order_id   │         │ • priority   │                             │
│  │   (PK)       │         └──────────────┘                             │
│  └──────┬───────┘                                                       │
│         │ 1                                                             │
│         │                 ┌──────────────┐                             │
│         ▼ 1               │   ratings    │                             │
│  ┌──────────────┐         │ ─────────────│                             │
│  │   ratings    │         │ • rating_id  │                             │
│  │ ─────────────│         │   (PK)       │                             │
│  │ • rating_id  │         │ • order_id   │                             │
│  │   (PK)       │         │   (FK)       │                             │
│  │ • order_id   │         │ • vendor_id  │                             │
│  │   (FK)       │         │ • service_id │                             │
│  │ • vendor_id  │         │ • rating     │                             │
│  │ • rating     │         │ • review     │                             │
│  │ • review     │         │ • tags       │                             │
│  └──────────────┘         └──────────────┘                             │
│                                                                          │
│  ┌──────────────┐                                                       │
│  │   users      │    1                                                  │
│  │ ─────────────│◄────┐                                                │
│  │ • user_id    │     │                                                │
│  │   (PK)       │     │ ∞                                               │
│  └──────────────┘     │                                                │
│                        │                                                │
│                 ┌──────────────┐                                        │
│                 │notifications │                                        │
│                 │ ─────────────│                                        │
│                 │ • notification                                        │
│                 │   _id (PK)   │                                        │
│                 │ • user_id    │                                        │
│                 │   (FK)       │                                        │
│                 │ • title      │                                        │
│                 │ • body       │                                        │
│                 │ • is_read    │                                        │
│                 │ • sent_via   │                                        │
│                 └──────────────┘                                        │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Core Tables

### 4.1 Users Table

**Purpose:** Central authentication and user management

```sql
CREATE TABLE users (
  user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  phone VARCHAR(15) UNIQUE NOT NULL,
  email VARCHAR(255),
  full_name VARCHAR(255),
  user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('RESIDENT', 'VENDOR', 'ADMIN', 'SOCIETY_ADMIN')),

  -- Profile
  profile_photo_url TEXT,

  -- Status
  is_verified BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,

  -- Notifications
  fcm_token TEXT,
  notification_enabled BOOLEAN DEFAULT true,

  -- Authentication
  last_login_at TIMESTAMP,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_type ON users(user_type);
CREATE INDEX idx_users_verified ON users(is_verified) WHERE is_verified = true;
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = true;
```

**Sample Data:**
```sql
INSERT INTO users (phone, full_name, user_type, is_verified) VALUES
  ('+919876543210', 'Ramesh Kumar', 'RESIDENT', true),
  ('+919876543211', 'Priya Sharma', 'VENDOR', true),
  ('+919876543212', 'Admin User', 'ADMIN', true);
```

---

### 4.2 Societies Table

**Purpose:** Housing societies/complexes

```sql
CREATE TABLE societies (
  society_id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  address TEXT NOT NULL,
  city VARCHAR(100) NOT NULL,
  state VARCHAR(100) NOT NULL,
  pincode VARCHAR(10) NOT NULL,

  -- Contact
  contact_person VARCHAR(255),
  contact_phone VARCHAR(15),
  contact_email VARCHAR(255),

  -- Stats
  total_flats INTEGER,
  occupied_flats INTEGER,

  -- Status
  status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'ACTIVE', 'SUSPENDED', 'INACTIVE')),
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  created_by UUID REFERENCES users(user_id)
);

-- Indexes
CREATE INDEX idx_societies_status ON societies(status);
CREATE INDEX idx_societies_city ON societies(city);
CREATE INDEX idx_societies_pincode ON societies(pincode);
CREATE INDEX idx_societies_active ON societies(is_active) WHERE is_active = true;

-- Full-text search
CREATE INDEX idx_societies_search ON societies USING gin(to_tsvector('english', name || ' ' || address));
```

---

### 4.3 Residents Table

**Purpose:** Resident-specific information

```sql
CREATE TABLE residents (
  resident_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
  society_id INTEGER REFERENCES societies(society_id),
  flat_number VARCHAR(20) NOT NULL,
  tower VARCHAR(10),
  floor INTEGER,

  -- Secondary societies (for multi-property residents)
  additional_societies JSONB DEFAULT '[]'::jsonb,

  -- Preferences
  preferred_pickup_time TIME,
  default_pickup_address TEXT,

  -- Status
  verification_status VARCHAR(20) DEFAULT 'PENDING' CHECK (verification_status IN ('PENDING', 'VERIFIED', 'REJECTED')),

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  UNIQUE(society_id, flat_number)
);

-- Indexes
CREATE INDEX idx_residents_society ON residents(society_id);
CREATE INDEX idx_residents_status ON residents(verification_status);
```

---

### 4.4 Vendors Table

**Purpose:** Service provider information (formerly laundry_persons)

```sql
CREATE TABLE vendors (
  vendor_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
  business_name VARCHAR(255),
  store_address TEXT NOT NULL,
  store_photo_url TEXT,

  -- Identity verification
  id_proof_type VARCHAR(50),
  id_proof_number VARCHAR(100),
  id_proof_photo_url TEXT,

  -- Business details
  gst_number VARCHAR(20),
  pan_number VARCHAR(20),

  -- Bank details for settlements
  bank_account_number VARCHAR(50),
  bank_ifsc_code VARCHAR(20),
  bank_account_holder VARCHAR(255),

  -- Stats
  total_orders INTEGER DEFAULT 0,
  completed_orders INTEGER DEFAULT 0,
  avg_rating DECIMAL(3,2) DEFAULT 0,

  -- Status
  approval_status VARCHAR(20) DEFAULT 'PENDING' CHECK (approval_status IN ('PENDING', 'APPROVED', 'REJECTED', 'SUSPENDED')),
  is_available BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  approved_at TIMESTAMP,
  approved_by UUID REFERENCES users(user_id)
);

-- Indexes
CREATE INDEX idx_vendors_status ON vendors(approval_status);
CREATE INDEX idx_vendors_available ON vendors(is_available) WHERE is_available = true;
CREATE INDEX idx_vendors_rating ON vendors(avg_rating DESC);

-- Full-text search
CREATE INDEX idx_vendors_search ON vendors USING gin(
  to_tsvector('english', COALESCE(business_name, '') || ' ' || store_address)
);
```

---

### 4.5 Vendor-Society Mapping

**Purpose:** Many-to-many relationship between vendors and societies

```sql
CREATE TABLE vendor_societies (
  id SERIAL PRIMARY KEY,
  vendor_id UUID REFERENCES vendors(vendor_id) ON DELETE CASCADE,
  society_id INTEGER REFERENCES societies(society_id) ON DELETE CASCADE,

  -- Approval per society
  approval_status VARCHAR(20) DEFAULT 'PENDING' CHECK (approval_status IN ('PENDING', 'APPROVED', 'REJECTED')),
  approved_at TIMESTAMP,
  approved_by UUID REFERENCES users(user_id),
  rejection_reason TEXT,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  UNIQUE(vendor_id, society_id)
);

-- Indexes
CREATE INDEX idx_vendor_societies_vendor ON vendor_societies(vendor_id);
CREATE INDEX idx_vendor_societies_society ON vendor_societies(society_id);
CREATE INDEX idx_vendor_societies_status ON vendor_societies(approval_status);
```

---

## 5. Category & Service Tables

### 5.1 Parent Categories Table

**Purpose:** Top-level service categories (Laundry, Vehicle, Home, Personal)

```sql
CREATE TABLE parent_categories (
  category_id SERIAL PRIMARY KEY,
  category_key VARCHAR(50) UNIQUE NOT NULL,
  category_name VARCHAR(100) NOT NULL,
  description TEXT,
  icon_url TEXT,
  color_hex VARCHAR(7),
  display_order INTEGER DEFAULT 0,

  -- Status
  is_active BOOLEAN DEFAULT true,
  is_live BOOLEAN DEFAULT false,  -- Controls visibility to users

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_parent_categories_active ON parent_categories(is_active) WHERE is_active = true;
CREATE INDEX idx_parent_categories_live ON parent_categories(is_live) WHERE is_live = true;
CREATE INDEX idx_parent_categories_order ON parent_categories(display_order);
```

**Initial Data:**
```sql
INSERT INTO parent_categories (category_key, category_name, description, color_hex, display_order, is_live) VALUES
  ('LAUNDRY', 'Laundry Services', 'Professional laundry and garment care services', '#3B82F6', 1, true),
  ('VEHICLE', 'Vehicle Services', 'Car and bike washing and detailing services', '#10B981', 2, false),
  ('HOME', 'Home Services', 'Gardening, plumbing, and home maintenance', '#F59E0B', 3, false),
  ('PERSONAL', 'Personal Care', 'Barber, salon, and spa services', '#EC4899', 4, false);
```

---

### 5.2 Service Categories Table

**Purpose:** Specific services under each parent category

```sql
CREATE TABLE service_categories (
  service_id SERIAL PRIMARY KEY,
  parent_category_id INTEGER REFERENCES parent_categories(category_id) ON DELETE CASCADE,
  service_key VARCHAR(50) NOT NULL,
  service_name VARCHAR(100) NOT NULL,
  description TEXT,
  icon_url TEXT,

  -- Service configuration
  default_turnaround_hours INTEGER DEFAULT 24,
  pricing_model VARCHAR(20) DEFAULT 'PER_ITEM' CHECK (pricing_model IN ('PER_ITEM', 'PER_SERVICE', 'HOURLY')),

  -- Display
  display_order INTEGER DEFAULT 0,
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  UNIQUE(parent_category_id, service_key)
);

-- Indexes
CREATE INDEX idx_service_categories_parent ON service_categories(parent_category_id);
CREATE INDEX idx_service_categories_active ON service_categories(is_active) WHERE is_active = true;
CREATE INDEX idx_service_categories_key ON service_categories(service_key);
```

**Initial Data:**
```sql
-- Laundry services
DO $$
DECLARE
  laundry_cat_id INTEGER;
  vehicle_cat_id INTEGER;
  home_cat_id INTEGER;
BEGIN
  SELECT category_id INTO laundry_cat_id FROM parent_categories WHERE category_key = 'LAUNDRY';
  SELECT category_id INTO vehicle_cat_id FROM parent_categories WHERE category_key = 'VEHICLE';
  SELECT category_id INTO home_cat_id FROM parent_categories WHERE category_key = 'HOME';

  -- Laundry services
  INSERT INTO service_categories (parent_category_id, service_key, service_name, description, default_turnaround_hours, display_order) VALUES
    (laundry_cat_id, 'IRONING', 'Ironing Only', 'Press and iron clothes', 24, 1),
    (laundry_cat_id, 'WASHING_IRONING', 'Washing + Ironing', 'Wash and iron clothes', 48, 2),
    (laundry_cat_id, 'DRY_CLEANING', 'Dry Cleaning', 'Professional dry cleaning', 120, 3),
    (laundry_cat_id, 'WASHING_ONLY', 'Washing Only', 'Wash clothes without ironing', 36, 4);

  -- Vehicle services
  INSERT INTO service_categories (parent_category_id, service_key, service_name, description, default_turnaround_hours, pricing_model, display_order) VALUES
    (vehicle_cat_id, 'CAR_WASH', 'Car Wash', 'Exterior and interior car cleaning', 2, 'PER_SERVICE', 1),
    (vehicle_cat_id, 'BIKE_WASH', 'Bike Wash', 'Motorcycle cleaning and polish', 1, 'PER_SERVICE', 2),
    (vehicle_cat_id, 'CAR_DETAILING', 'Car Detailing', 'Complete car detailing and restoration', 24, 'PER_SERVICE', 3);

  -- Home services
  INSERT INTO service_categories (parent_category_id, service_key, service_name, description, default_turnaround_hours, pricing_model, display_order) VALUES
    (home_cat_id, 'GARDENING', 'Gardening', 'Lawn mowing, trimming, and maintenance', 3, 'HOURLY', 1),
    (home_cat_id, 'PLUMBING', 'Plumbing', 'Pipe fixing, leakage repair', 2, 'PER_SERVICE', 2),
    (home_cat_id, 'ELECTRICAL', 'Electrical', 'Wiring, fixture installation', 2, 'PER_SERVICE', 3);
END $$;
```

---

### 5.3 Vendor Services Table

**Purpose:** Services offered by each vendor

```sql
CREATE TABLE vendor_services (
  id SERIAL PRIMARY KEY,
  vendor_id UUID REFERENCES vendors(vendor_id) ON DELETE CASCADE,
  service_id INTEGER REFERENCES service_categories(service_id) ON DELETE CASCADE,

  -- Configuration
  is_active BOOLEAN DEFAULT true,
  turnaround_hours INTEGER DEFAULT 24,  -- Can override default

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  UNIQUE(vendor_id, service_id)
);

-- Indexes
CREATE INDEX idx_vendor_services_vendor ON vendor_services(vendor_id);
CREATE INDEX idx_vendor_services_service ON vendor_services(service_id);
CREATE INDEX idx_vendor_services_active ON vendor_services(is_active) WHERE is_active = true;
```

---

### 5.4 Rate Cards Table

**Purpose:** Pricing structure per vendor per society

```sql
CREATE TABLE rate_cards (
  rate_card_id SERIAL PRIMARY KEY,
  vendor_id UUID REFERENCES vendors(vendor_id) ON DELETE CASCADE,
  society_id INTEGER REFERENCES societies(society_id) ON DELETE CASCADE,

  -- Status
  is_active BOOLEAN DEFAULT true,
  is_published BOOLEAN DEFAULT false,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  published_at TIMESTAMP,

  UNIQUE(vendor_id, society_id)
);

-- Indexes
CREATE INDEX idx_rate_cards_vendor ON rate_cards(vendor_id);
CREATE INDEX idx_rate_cards_society ON rate_cards(society_id);
CREATE INDEX idx_rate_cards_active ON rate_cards(is_active) WHERE is_active = true;
```

---

### 5.5 Rate Card Items Table

**Purpose:** Individual items with pricing per service type

```sql
CREATE TABLE rate_card_items (
  item_id SERIAL PRIMARY KEY,
  rate_card_id INTEGER REFERENCES rate_cards(rate_card_id) ON DELETE CASCADE,
  service_id INTEGER REFERENCES service_categories(service_id) ON DELETE CASCADE,

  -- Item details
  item_name VARCHAR(100) NOT NULL,
  description TEXT,
  price_per_piece DECIMAL(6,2) NOT NULL CHECK (price_per_piece >= 0),

  -- Display
  display_order INTEGER DEFAULT 0,
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_rate_items_card ON rate_card_items(rate_card_id);
CREATE INDEX idx_rate_items_service ON rate_card_items(service_id);
CREATE INDEX idx_rate_items_active ON rate_card_items(is_active) WHERE is_active = true;
```

---

## 6. Workflow Configuration Tables

### 6.1 Service Workflow Templates Table

**Purpose:** Define workflow templates for each service type

```sql
CREATE TABLE service_workflow_templates (
  template_id SERIAL PRIMARY KEY,
  service_id INTEGER REFERENCES service_categories(service_id) ON DELETE CASCADE,
  template_name VARCHAR(100) NOT NULL,
  description TEXT,

  -- Configuration
  is_default BOOLEAN DEFAULT true,
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  UNIQUE(service_id, template_name)
);

-- Indexes
CREATE INDEX idx_workflow_templates_service ON service_workflow_templates(service_id);
CREATE INDEX idx_workflow_templates_default ON service_workflow_templates(service_id, is_default) WHERE is_default = true;
```

**Purpose:**
- Each service type can have a default workflow
- Vendors can customize workflows (future feature)
- Admin can create multiple workflow templates per service

---

### 6.2 Workflow Steps Table

**Purpose:** Define individual steps in each workflow

```sql
CREATE TABLE workflow_steps (
  step_id SERIAL PRIMARY KEY,
  template_id INTEGER REFERENCES service_workflow_templates(template_id) ON DELETE CASCADE,
  step_name VARCHAR(100) NOT NULL,
  step_key VARCHAR(50) NOT NULL,  -- e.g., 'pickup', 'count', 'iron', 'quality_check'
  description TEXT,

  -- Order and requirements
  step_order INTEGER NOT NULL,
  is_required BOOLEAN DEFAULT true,
  requires_photo BOOLEAN DEFAULT false,
  requires_signature BOOLEAN DEFAULT false,
  requires_notes BOOLEAN DEFAULT false,

  -- Time tracking
  estimated_duration_minutes INTEGER,

  -- Status mapping
  order_status_on_complete VARCHAR(50),  -- Maps to order_status enum

  -- Display
  icon VARCHAR(50),
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  UNIQUE(template_id, step_order)
);

-- Indexes
CREATE INDEX idx_workflow_steps_template ON workflow_steps(template_id, step_order);
CREATE INDEX idx_workflow_steps_key ON workflow_steps(step_key);
CREATE INDEX idx_workflow_steps_active ON workflow_steps(is_active) WHERE is_active = true;
```

**Example step keys:**
- `pickup` - Pick up items
- `count` - Count and verify items
- `wash` - Washing process
- `iron` - Ironing process
- `dry_clean` - Dry cleaning process
- `quality_check` - Quality inspection
- `pack` - Package items
- `ready` - Ready for delivery
- `deliver` - Delivery to customer

---

### 6.3 Initial Workflow Data

```sql
-- Ironing workflow
DO $$
DECLARE
  laundry_cat_id INTEGER;
  ironing_service_id INTEGER;
  ironing_template_id INTEGER;
BEGIN
  SELECT category_id INTO laundry_cat_id FROM parent_categories WHERE category_key = 'LAUNDRY';
  SELECT service_id INTO ironing_service_id FROM service_categories WHERE service_key = 'IRONING';

  -- Create workflow template
  INSERT INTO service_workflow_templates (service_id, template_name, description, is_default)
  VALUES (ironing_service_id, 'Standard Ironing Workflow', 'Standard workflow for ironing service', true)
  RETURNING template_id INTO ironing_template_id;

  -- Add workflow steps
  INSERT INTO workflow_steps (template_id, step_name, step_key, step_order, is_required, requires_photo, estimated_duration_minutes, order_status_on_complete) VALUES
    (ironing_template_id, 'Pickup Items', 'pickup', 1, true, false, 15, 'PICKUP_IN_PROGRESS'),
    (ironing_template_id, 'Count Items', 'count', 2, true, true, 10, 'COUNT_APPROVAL_PENDING'),
    (ironing_template_id, 'Iron Items', 'iron', 3, true, false, 60, 'PROCESSING_IN_PROGRESS'),
    (ironing_template_id, 'Quality Check', 'quality_check', 4, true, false, 10, 'READY_FOR_DELIVERY'),
    (ironing_template_id, 'Deliver Items', 'deliver', 5, true, true, 15, 'DELIVERED');
END $$;

-- Dry Cleaning workflow
DO $$
DECLARE
  dry_clean_service_id INTEGER;
  dry_clean_template_id INTEGER;
BEGIN
  SELECT service_id INTO dry_clean_service_id FROM service_categories WHERE service_key = 'DRY_CLEANING';

  -- Create workflow template
  INSERT INTO service_workflow_templates (service_id, template_name, description, is_default)
  VALUES (dry_clean_service_id, 'Standard Dry Cleaning Workflow', 'Standard workflow for dry cleaning service', true)
  RETURNING template_id INTO dry_clean_template_id;

  -- Add workflow steps
  INSERT INTO workflow_steps (template_id, step_name, step_key, step_order, is_required, requires_photo, estimated_duration_minutes, order_status_on_complete) VALUES
    (dry_clean_template_id, 'Pickup Items', 'pickup', 1, true, false, 15, 'PICKUP_IN_PROGRESS'),
    (dry_clean_template_id, 'Count Items', 'count', 2, true, true, 10, 'COUNT_APPROVAL_PENDING'),
    (dry_clean_template_id, 'Pre-Treatment', 'pre_treatment', 3, true, false, 30, 'PROCESSING_IN_PROGRESS'),
    (dry_clean_template_id, 'Dry Clean', 'dry_clean', 4, true, false, 120, 'PROCESSING_IN_PROGRESS'),
    (dry_clean_template_id, 'Quality Check', 'quality_check', 5, true, false, 15, 'PROCESSING_IN_PROGRESS'),
    (dry_clean_template_id, 'Press & Finish', 'press_finish', 6, true, false, 30, 'READY_FOR_DELIVERY'),
    (dry_clean_template_id, 'Deliver Items', 'deliver', 7, true, true, 15, 'DELIVERED');
END $$;

-- Car Wash workflow
DO $$
DECLARE
  car_wash_service_id INTEGER;
  car_wash_template_id INTEGER;
BEGIN
  SELECT service_id INTO car_wash_service_id FROM service_categories WHERE service_key = 'CAR_WASH';

  -- Create workflow template
  INSERT INTO service_workflow_templates (service_id, template_name, description, is_default)
  VALUES (car_wash_service_id, 'Standard Car Wash Workflow', 'Standard workflow for car washing service', true)
  RETURNING template_id INTO car_wash_template_id;

  -- Add workflow steps
  INSERT INTO workflow_steps (template_id, step_name, step_key, step_order, is_required, requires_photo, estimated_duration_minutes, order_status_on_complete) VALUES
    (car_wash_template_id, 'Arrive at Location', 'arrive', 1, true, false, 0, 'PICKUP_IN_PROGRESS'),
    (car_wash_template_id, 'Initial Inspection', 'inspect', 2, true, true, 5, 'PROCESSING_IN_PROGRESS'),
    (car_wash_template_id, 'Exterior Wash', 'exterior_wash', 3, true, false, 20, 'PROCESSING_IN_PROGRESS'),
    (car_wash_template_id, 'Interior Vacuum', 'interior_vacuum', 4, false, false, 15, 'PROCESSING_IN_PROGRESS'),
    (car_wash_template_id, 'Polish & Wax', 'polish', 5, false, false, 20, 'PROCESSING_IN_PROGRESS'),
    (car_wash_template_id, 'Final Check', 'final_check', 6, true, true, 5, 'COMPLETED');
END $$;

-- Gardening workflow
DO $$
DECLARE
  gardening_service_id INTEGER;
  gardening_template_id INTEGER;
BEGIN
  SELECT service_id INTO gardening_service_id FROM service_categories WHERE service_key = 'GARDENING';

  -- Create workflow template
  INSERT INTO service_workflow_templates (service_id, template_name, description, is_default)
  VALUES (gardening_service_id, 'Standard Gardening Workflow', 'Standard workflow for gardening service', true)
  RETURNING template_id INTO gardening_template_id;

  -- Add workflow steps
  INSERT INTO workflow_steps (template_id, step_name, step_key, step_order, is_required, requires_photo, estimated_duration_minutes, order_status_on_complete) VALUES
    (gardening_template_id, 'Arrive & Assess', 'arrive', 1, true, true, 10, 'PICKUP_IN_PROGRESS'),
    (gardening_template_id, 'Trim Plants', 'trim', 2, false, false, 30, 'PROCESSING_IN_PROGRESS'),
    (gardening_template_id, 'Mow Lawn', 'mow', 3, false, false, 30, 'PROCESSING_IN_PROGRESS'),
    (gardening_template_id, 'Weed Removal', 'weed', 4, false, false, 20, 'PROCESSING_IN_PROGRESS'),
    (gardening_template_id, 'Clean Up', 'cleanup', 5, true, false, 15, 'PROCESSING_IN_PROGRESS'),
    (gardening_template_id, 'Final Inspection', 'final_check', 6, true, true, 5, 'COMPLETED');
END $$;

-- Plumbing workflow
DO $$
DECLARE
  plumbing_service_id INTEGER;
  plumbing_template_id INTEGER;
BEGIN
  SELECT service_id INTO plumbing_service_id FROM service_categories WHERE service_key = 'PLUMBING';

  -- Create workflow template
  INSERT INTO service_workflow_templates (service_id, template_name, description, is_default)
  VALUES (plumbing_service_id, 'Standard Plumbing Workflow', 'Standard workflow for plumbing service', true)
  RETURNING template_id INTO plumbing_template_id;

  -- Add workflow steps
  INSERT INTO workflow_steps (template_id, step_name, step_key, step_order, is_required, requires_photo, estimated_duration_minutes, order_status_on_complete) VALUES
    (plumbing_template_id, 'Arrive & Inspect', 'arrive', 1, true, true, 10, 'PICKUP_IN_PROGRESS'),
    (plumbing_template_id, 'Diagnose Issue', 'diagnose', 2, true, false, 15, 'PROCESSING_IN_PROGRESS'),
    (plumbing_template_id, 'Perform Repair', 'repair', 3, true, false, 60, 'PROCESSING_IN_PROGRESS'),
    (plumbing_template_id, 'Test System', 'test', 4, true, false, 10, 'PROCESSING_IN_PROGRESS'),
    (plumbing_template_id, 'Final Check', 'final_check', 5, true, true, 5, 'COMPLETED');
END $$;
```

---

## 7. Order Management Tables

### 7.1 Orders Table

**Purpose:** Main order records

```sql
CREATE TYPE order_status AS ENUM (
  'BOOKING_CREATED',
  'PICKUP_SCHEDULED',
  'PICKUP_IN_PROGRESS',
  'COUNT_APPROVAL_PENDING',
  'PICKED_UP',
  'PROCESSING_IN_PROGRESS',
  'READY_FOR_DELIVERY',
  'OUT_FOR_DELIVERY',
  'DELIVERED',
  'COMPLETED',
  'CANCELLED',
  'DISPUTED',
  'ON_HOLD'
);

CREATE TABLE orders (
  order_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  order_number VARCHAR(50) UNIQUE NOT NULL,

  -- Participants
  resident_id UUID REFERENCES residents(resident_id) ON DELETE RESTRICT,
  vendor_id UUID REFERENCES vendors(vendor_id) ON DELETE RESTRICT,
  society_id INTEGER REFERENCES societies(society_id) ON DELETE RESTRICT,

  -- Status
  status order_status DEFAULT 'BOOKING_CREATED',

  -- Multi-service flag
  has_multiple_services BOOLEAN DEFAULT false,

  -- Pricing
  estimated_price DECIMAL(10,2) NOT NULL,
  final_price DECIMAL(10,2),
  discount_amount DECIMAL(10,2) DEFAULT 0,

  -- Counts
  estimated_item_count INTEGER,
  actual_item_count INTEGER,
  count_difference INTEGER,
  count_approved_by_resident BOOLEAN,

  -- Scheduling
  pickup_datetime TIMESTAMP NOT NULL,
  pickup_address TEXT NOT NULL,
  expected_delivery_date DATE,
  actual_delivery_date DATE,

  -- Photos
  pickup_photos JSONB DEFAULT '[]'::jsonb,
  delivery_photos JSONB DEFAULT '[]'::jsonb,

  -- Notes
  resident_notes TEXT,
  vendor_notes TEXT,
  admin_notes TEXT,

  -- Cancellation
  cancellation_reason TEXT,
  cancelled_by UUID REFERENCES users(user_id),
  cancelled_at TIMESTAMP,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_orders_resident ON orders(resident_id, created_at DESC);
CREATE INDEX idx_orders_vendor ON orders(vendor_id, status, created_at DESC);
CREATE INDEX idx_orders_society ON orders(society_id, created_at DESC);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_pickup ON orders(pickup_datetime);
CREATE INDEX idx_orders_delivery ON orders(expected_delivery_date);
CREATE INDEX idx_orders_number ON orders(order_number);
```

**Trigger: Generate Order Number**

```sql
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TRIGGER AS $$
BEGIN
  NEW.order_number := 'ORD' || TO_CHAR(NOW(), 'YYYYMMDD') || LPAD(NEXTVAL('order_number_seq')::TEXT, 6, '0');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE SEQUENCE order_number_seq;

CREATE TRIGGER trigger_generate_order_number
  BEFORE INSERT ON orders
  FOR EACH ROW
  EXECUTE FUNCTION generate_order_number();
```

---

### 7.2 Order Items Table

**Purpose:** Individual items in each order

```sql
CREATE TABLE order_items (
  id SERIAL PRIMARY KEY,
  order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
  rate_card_item_id INTEGER REFERENCES rate_card_items(item_id),
  service_id INTEGER REFERENCES service_categories(service_id) ON DELETE RESTRICT,

  -- Item details
  item_name VARCHAR(100) NOT NULL,
  quantity INTEGER NOT NULL CHECK (quantity > 0),
  unit_price DECIMAL(6,2) NOT NULL,
  total_price DECIMAL(10,2) NOT NULL,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_service ON order_items(service_id);
```

---

### 7.3 Order Service Status Table

**Purpose:** Track progress of each service type within an order

```sql
CREATE TABLE order_service_status (
  id SERIAL PRIMARY KEY,
  order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
  service_id INTEGER REFERENCES service_categories(service_id) ON DELETE RESTRICT,
  template_id INTEGER REFERENCES service_workflow_templates(template_id),

  -- Aggregates
  item_count INTEGER NOT NULL,
  total_amount DECIMAL(10,2) NOT NULL,

  -- Workflow tracking
  current_step_id INTEGER REFERENCES workflow_steps(step_id),
  current_step_order INTEGER DEFAULT 1,

  -- Status tracking per service
  status order_status DEFAULT 'PICKED_UP',
  processing_started_at TIMESTAMP,
  ready_at TIMESTAMP,
  delivered_at TIMESTAMP,

  -- Expected delivery for this service
  expected_delivery_date DATE,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  UNIQUE(order_id, service_id)
);

-- Indexes
CREATE INDEX idx_order_service_status_order ON order_service_status(order_id);
CREATE INDEX idx_order_service_status_service ON order_service_status(service_id);
CREATE INDEX idx_order_service_status_status ON order_service_status(status);
CREATE INDEX idx_order_service_status_step ON order_service_status(current_step_id);
```

---

### 7.4 Order Workflow Progress Table

**Purpose:** Track completion of individual workflow steps for each service in an order

```sql
CREATE TABLE order_workflow_progress (
  id SERIAL PRIMARY KEY,
  order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
  service_id INTEGER REFERENCES service_categories(service_id) ON DELETE RESTRICT,
  step_id INTEGER REFERENCES workflow_steps(step_id) ON DELETE RESTRICT,

  -- Status
  status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'SKIPPED', 'FAILED')),

  -- Execution details
  started_at TIMESTAMP,
  completed_at TIMESTAMP,
  duration_minutes INTEGER,

  -- Data captured during step
  photos JSONB DEFAULT '[]'::jsonb,
  signature_url TEXT,
  notes TEXT,
  metadata JSONB DEFAULT '{}'::jsonb,

  -- Who completed it
  completed_by UUID REFERENCES users(user_id),

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  UNIQUE(order_id, service_id, step_id)
);

-- Indexes
CREATE INDEX idx_workflow_progress_order ON order_workflow_progress(order_id, service_id);
CREATE INDEX idx_workflow_progress_step ON order_workflow_progress(step_id);
CREATE INDEX idx_workflow_progress_status ON order_workflow_progress(status);
CREATE INDEX idx_workflow_progress_completed ON order_workflow_progress(completed_at DESC);
```

**Purpose:**
- Each row represents one workflow step for one service in an order
- Example: Order #123 has Ironing service → 5 rows (one per workflow step)
- Tracks when each step started, completed, and who did it
- Stores photos/signatures/notes captured during each step

---

### 7.5 Order Status Log Table

**Purpose:** Audit trail of status changes

```sql
CREATE TABLE order_status_log (
  id SERIAL PRIMARY KEY,
  order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
  service_id INTEGER REFERENCES service_categories(service_id),  -- NULL if overall order status

  -- Status change
  from_status order_status,
  to_status order_status NOT NULL,

  -- Who changed it
  changed_by UUID REFERENCES users(user_id),
  changed_by_role VARCHAR(20),

  -- Notes
  notes TEXT,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_order_status_log_order ON order_status_log(order_id, created_at DESC);
CREATE INDEX idx_order_status_log_service ON order_status_log(service_id);
CREATE INDEX idx_order_status_log_status ON order_status_log(to_status);
```

---

## 8. Payment & Settlement Tables

### 8.1 Payments Table

**Purpose:** Track all payments (resident to vendor)

```sql
CREATE TABLE payments (
  payment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  order_id UUID REFERENCES orders(order_id) ON DELETE RESTRICT,

  -- Amount
  amount DECIMAL(10,2) NOT NULL,

  -- Method
  payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('UPI', 'CASH', 'CARD', 'OTHER')),

  -- Status
  status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED', 'REFUNDED')),

  -- Payment details
  razorpay_order_id VARCHAR(100),
  razorpay_payment_id VARCHAR(100),
  razorpay_signature VARCHAR(255),

  -- UPI details
  upi_transaction_id VARCHAR(100),
  upi_vpa VARCHAR(100),

  -- Timing
  paid_at TIMESTAMP,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_method ON payments(payment_method);
CREATE INDEX idx_payments_razorpay_order ON payments(razorpay_order_id);
```

---

### 8.2 Settlements Table

**Purpose:** Track vendor earnings and payouts

```sql
CREATE TABLE settlements (
  settlement_id SERIAL PRIMARY KEY,
  vendor_id UUID REFERENCES vendors(vendor_id) ON DELETE CASCADE,

  -- Period
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,

  -- Amounts
  total_orders INTEGER,
  gross_amount DECIMAL(10,2) NOT NULL,
  platform_fee DECIMAL(10,2) DEFAULT 0,  -- Currently 0
  net_amount DECIMAL(10,2) NOT NULL,

  -- Status
  status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PROCESSING', 'PAID', 'FAILED')),

  -- Payout details
  payout_method VARCHAR(20),
  bank_reference_number VARCHAR(100),
  paid_at TIMESTAMP,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_settlements_vendor ON settlements(vendor_id, period_end DESC);
CREATE INDEX idx_settlements_status ON settlements(status);
CREATE INDEX idx_settlements_period ON settlements(period_start, period_end);
```

---

## 9. Subscription & Billing Tables

### 9.1 Society Subscriptions Table

**Purpose:** Society subscription plans and status

```sql
CREATE TYPE subscription_tier AS ENUM ('STARTER', 'GROWTH', 'ENTERPRISE');
CREATE TYPE subscription_status AS ENUM ('TRIAL', 'ACTIVE', 'SUSPENDED', 'CANCELLED', 'EXPIRED');

CREATE TABLE society_subscriptions (
  subscription_id SERIAL PRIMARY KEY,
  society_id INTEGER UNIQUE REFERENCES societies(society_id) ON DELETE CASCADE,

  -- Plan
  tier subscription_tier NOT NULL,
  monthly_fee DECIMAL(10,2) NOT NULL,

  -- Status
  status subscription_status DEFAULT 'TRIAL',

  -- Trial
  is_trial BOOLEAN DEFAULT true,
  trial_start_date DATE,
  trial_end_date DATE,

  -- Billing
  billing_start_date DATE NOT NULL,
  current_period_start DATE NOT NULL,
  current_period_end DATE NOT NULL,
  next_billing_date DATE NOT NULL,

  -- Cancellation
  cancelled_at TIMESTAMP,
  cancellation_reason TEXT,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_subscriptions_society ON society_subscriptions(society_id);
CREATE INDEX idx_subscriptions_status ON society_subscriptions(status);
CREATE INDEX idx_subscriptions_tier ON society_subscriptions(tier);
CREATE INDEX idx_subscriptions_next_billing ON society_subscriptions(next_billing_date);
```

---

### 9.2 Subscription Invoices Table

**Purpose:** Monthly invoices for society subscriptions

```sql
CREATE TABLE subscription_invoices (
  invoice_id SERIAL PRIMARY KEY,
  subscription_id INTEGER REFERENCES society_subscriptions(subscription_id) ON DELETE CASCADE,
  society_id INTEGER REFERENCES societies(society_id) ON DELETE CASCADE,
  invoice_number VARCHAR(50) UNIQUE NOT NULL,

  -- Amount
  amount DECIMAL(10,2) NOT NULL,
  tax_amount DECIMAL(10,2) DEFAULT 0,
  total_amount DECIMAL(10,2) NOT NULL,

  -- Period
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,

  -- Status
  status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PAID', 'OVERDUE', 'CANCELLED')),

  -- Payment
  due_date DATE NOT NULL,
  paid_at TIMESTAMP,
  payment_method VARCHAR(20),
  payment_reference VARCHAR(100),

  -- Razorpay details
  razorpay_order_id VARCHAR(100),
  razorpay_payment_id VARCHAR(100),

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_invoices_subscription ON subscription_invoices(subscription_id);
CREATE INDEX idx_invoices_society ON subscription_invoices(society_id);
CREATE INDEX idx_invoices_status ON subscription_invoices(status);
CREATE INDEX idx_invoices_due_date ON subscription_invoices(due_date);
CREATE INDEX idx_invoices_number ON subscription_invoices(invoice_number);
```

---

## 10. Support & Communication Tables

### 10.1 Disputes Table

**Purpose:** Handle order issues and disputes

```sql
CREATE TABLE disputes (
  dispute_id SERIAL PRIMARY KEY,
  order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
  service_id INTEGER REFERENCES service_categories(service_id),  -- Which service had the issue
  raised_by UUID REFERENCES users(user_id),

  -- Issue
  issue_type VARCHAR(50) NOT NULL CHECK (issue_type IN (
    'ITEM_MISSING',
    'ITEM_DAMAGED',
    'QUALITY_ISSUE',
    'DELAY',
    'WRONG_COUNT',
    'PAYMENT_ISSUE',
    'OTHER'
  )),
  description TEXT NOT NULL,

  -- Evidence
  photos JSONB DEFAULT '[]'::jsonb,

  -- Status
  status VARCHAR(20) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED', 'ESCALATED')),
  priority VARCHAR(20) DEFAULT 'MEDIUM' CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')),

  -- Resolution
  resolution_notes TEXT,
  resolved_by UUID REFERENCES users(user_id),
  resolved_at TIMESTAMP,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_disputes_order ON disputes(order_id);
CREATE INDEX idx_disputes_service ON disputes(service_id);
CREATE INDEX idx_disputes_raised_by ON disputes(raised_by);
CREATE INDEX idx_disputes_status ON disputes(status);
CREATE INDEX idx_disputes_priority ON disputes(priority);
CREATE INDEX idx_disputes_created ON disputes(created_at DESC);
```

---

### 10.2 Ratings & Reviews Table

**Purpose:** Customer feedback

```sql
CREATE TABLE ratings (
  rating_id SERIAL PRIMARY KEY,
  order_id UUID UNIQUE REFERENCES orders(order_id) ON DELETE CASCADE,
  resident_id UUID REFERENCES residents(resident_id),
  vendor_id UUID REFERENCES vendors(vendor_id),
  service_id INTEGER REFERENCES service_categories(service_id),

  -- Rating
  rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
  review TEXT,

  -- Tags
  tags VARCHAR(50)[],  -- e.g., {'punctual', 'quality', 'friendly'}

  -- Response
  vendor_response TEXT,
  vendor_responded_at TIMESTAMP,

  -- Moderation
  is_published BOOLEAN DEFAULT true,
  moderation_notes TEXT,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_ratings_vendor ON ratings(vendor_id, created_at DESC);
CREATE INDEX idx_ratings_service ON ratings(service_id);
CREATE INDEX idx_ratings_rating ON ratings(rating);
CREATE INDEX idx_ratings_published ON ratings(is_published) WHERE is_published = true;
```

**Trigger: Update Vendor Average Rating**

```sql
CREATE OR REPLACE FUNCTION update_vendor_avg_rating()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE vendors
  SET avg_rating = (
    SELECT AVG(rating)::DECIMAL(3,2)
    FROM ratings
    WHERE vendor_id = NEW.vendor_id
      AND is_published = true
  )
  WHERE vendor_id = NEW.vendor_id;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_vendor_avg_rating
  AFTER INSERT OR UPDATE ON ratings
  FOR EACH ROW
  EXECUTE FUNCTION update_vendor_avg_rating();
```

---

### 10.3 Notifications Table

**Purpose:** Notification history and status

```sql
CREATE TABLE notifications (
  notification_id SERIAL PRIMARY KEY,
  user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,

  -- Notification details
  title VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  notification_type VARCHAR(50) NOT NULL,

  -- Link/Action
  action_type VARCHAR(50),  -- e.g., 'ORDER_DETAILS', 'PAYMENT'
  action_data JSONB,         -- e.g., {"order_id": "123"}

  -- Status
  is_read BOOLEAN DEFAULT false,
  read_at TIMESTAMP,

  -- Delivery
  sent_via VARCHAR(20)[],  -- e.g., {'PUSH', 'SMS', 'EMAIL'}
  delivery_status JSONB,   -- Status per channel

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_notifications_user ON notifications(user_id, created_at DESC);
CREATE INDEX idx_notifications_read ON notifications(user_id, is_read) WHERE is_read = false;
CREATE INDEX idx_notifications_type ON notifications(notification_type);
```

---

## 11. Indexes & Performance

### 11.1 Critical Indexes Summary

```sql
-- User lookups
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_fcm_token ON users(fcm_token) WHERE fcm_token IS NOT NULL;

-- Order queries (most critical)
CREATE INDEX idx_orders_resident_status ON orders(resident_id, status, created_at DESC);
CREATE INDEX idx_orders_vendor_status ON orders(vendor_id, status, pickup_datetime);
CREATE INDEX idx_orders_society_date ON orders(society_id, created_at DESC);
CREATE INDEX idx_orders_status_pickup ON orders(status, pickup_datetime) WHERE status IN ('PICKUP_SCHEDULED', 'PICKUP_IN_PROGRESS');

-- Service lookups
CREATE INDEX idx_service_categories_parent_active ON service_categories(parent_category_id, is_active) WHERE is_active = true;
CREATE INDEX idx_vendor_services_composite ON vendor_services(vendor_id, service_id, is_active);

-- Workflow queries
CREATE INDEX idx_workflow_steps_template_order ON workflow_steps(template_id, step_order);
CREATE INDEX idx_workflow_progress_order_service ON order_workflow_progress(order_id, service_id, step_id);

-- Rate card searches
CREATE INDEX idx_rate_cards_society_active ON rate_cards(society_id, is_active) WHERE is_active = true AND is_published = true;

-- Payment tracking
CREATE INDEX idx_payments_order_status ON payments(order_id, status);

-- Subscription billing
CREATE INDEX idx_subscriptions_billing ON society_subscriptions(next_billing_date, status) WHERE status = 'ACTIVE';
CREATE INDEX idx_invoices_overdue ON subscription_invoices(due_date, status) WHERE status = 'PENDING';
```

---

## 12. Database Functions

### 12.1 Get Service Workflow Steps

```sql
CREATE OR REPLACE FUNCTION get_service_workflow_steps(p_service_id INTEGER)
RETURNS TABLE(
  step_id INTEGER,
  step_name VARCHAR(100),
  step_key VARCHAR(50),
  step_order INTEGER,
  is_required BOOLEAN,
  requires_photo BOOLEAN,
  estimated_duration_minutes INTEGER
) AS $$
BEGIN
  RETURN QUERY
  SELECT
    ws.step_id,
    ws.step_name,
    ws.step_key,
    ws.step_order,
    ws.is_required,
    ws.requires_photo,
    ws.estimated_duration_minutes
  FROM workflow_steps ws
  JOIN service_workflow_templates swt ON ws.template_id = swt.template_id
  WHERE swt.service_id = p_service_id
    AND swt.is_default = true
    AND ws.is_active = true
  ORDER BY ws.step_order;
END;
$$ LANGUAGE plpgsql;
```

### 12.2 Initialize Order Workflow

```sql
CREATE OR REPLACE FUNCTION initialize_order_workflow(
  p_order_id UUID,
  p_service_id INTEGER
)
RETURNS VOID AS $$
DECLARE
  v_template_id INTEGER;
  v_step RECORD;
BEGIN
  -- Get default workflow template for service
  SELECT template_id INTO v_template_id
  FROM service_workflow_templates
  WHERE service_id = p_service_id
    AND is_default = true
    AND is_active = true
  LIMIT 1;

  IF v_template_id IS NULL THEN
    RAISE EXCEPTION 'No workflow template found for service_id %', p_service_id;
  END IF;

  -- Update order_service_status with template
  UPDATE order_service_status
  SET template_id = v_template_id
  WHERE order_id = p_order_id
    AND service_id = p_service_id;

  -- Create workflow progress entries for all steps
  FOR v_step IN
    SELECT step_id, step_order
    FROM workflow_steps
    WHERE template_id = v_template_id
      AND is_active = true
    ORDER BY step_order
  LOOP
    INSERT INTO order_workflow_progress (
      order_id,
      service_id,
      step_id,
      status
    ) VALUES (
      p_order_id,
      p_service_id,
      v_step.step_id,
      CASE WHEN v_step.step_order = 1 THEN 'PENDING' ELSE 'PENDING' END
    )
    ON CONFLICT (order_id, service_id, step_id) DO NOTHING;
  END LOOP;
END;
$$ LANGUAGE plpgsql;
```

### 12.3 Complete Workflow Step

```sql
CREATE OR REPLACE FUNCTION complete_workflow_step(
  p_order_id UUID,
  p_service_id INTEGER,
  p_step_id INTEGER,
  p_completed_by UUID,
  p_photos JSONB DEFAULT NULL,
  p_notes TEXT DEFAULT NULL
)
RETURNS JSONB AS $$
DECLARE
  v_next_step_id INTEGER;
  v_next_step_order INTEGER;
  v_order_status_to_set order_status;
  v_result JSONB;
BEGIN
  -- Mark current step as completed
  UPDATE order_workflow_progress
  SET
    status = 'COMPLETED',
    completed_at = NOW(),
    duration_minutes = EXTRACT(EPOCH FROM (NOW() - started_at)) / 60,
    completed_by = p_completed_by,
    photos = COALESCE(p_photos, photos),
    notes = COALESCE(p_notes, notes),
    updated_at = NOW()
  WHERE order_id = p_order_id
    AND service_id = p_service_id
    AND step_id = p_step_id;

  -- Get order status to set from workflow step
  SELECT order_status_on_complete INTO v_order_status_to_set
  FROM workflow_steps
  WHERE step_id = p_step_id;

  -- Update order_service_status with new order status
  IF v_order_status_to_set IS NOT NULL THEN
    UPDATE order_service_status
    SET
      status = v_order_status_to_set::order_status,
      updated_at = NOW()
    WHERE order_id = p_order_id
      AND service_id = p_service_id;
  END IF;

  -- Get next step
  SELECT ws.step_id, ws.step_order
  INTO v_next_step_id, v_next_step_order
  FROM workflow_steps ws
  JOIN order_service_status oss ON ws.template_id = oss.template_id
  WHERE oss.order_id = p_order_id
    AND oss.service_id = p_service_id
    AND ws.step_order > (SELECT step_order FROM workflow_steps WHERE step_id = p_step_id)
    AND ws.is_active = true
  ORDER BY ws.step_order
  LIMIT 1;

  -- Update current step in order_service_status
  IF v_next_step_id IS NOT NULL THEN
    UPDATE order_service_status
    SET
      current_step_id = v_next_step_id,
      current_step_order = v_next_step_order,
      updated_at = NOW()
    WHERE order_id = p_order_id
      AND service_id = p_service_id;

    -- Mark next step as IN_PROGRESS
    UPDATE order_workflow_progress
    SET
      status = 'IN_PROGRESS',
      started_at = NOW(),
      updated_at = NOW()
    WHERE order_id = p_order_id
      AND service_id = p_service_id
      AND step_id = v_next_step_id;
  END IF;

  -- Build result
  v_result := jsonb_build_object(
    'completed_step_id', p_step_id,
    'next_step_id', v_next_step_id,
    'order_status', v_order_status_to_set,
    'is_final_step', (v_next_step_id IS NULL)
  );

  RETURN v_result;
END;
$$ LANGUAGE plpgsql;
```

### 12.4 Generate Monthly Invoices

```sql
CREATE OR REPLACE FUNCTION generate_monthly_invoices()
RETURNS INTEGER AS $$
DECLARE
  v_count INTEGER := 0;
  v_subscription RECORD;
BEGIN
  FOR v_subscription IN
    SELECT *
    FROM society_subscriptions
    WHERE next_billing_date <= CURRENT_DATE
      AND status IN ('TRIAL', 'ACTIVE')
  LOOP
    -- Generate invoice
    INSERT INTO subscription_invoices (
      subscription_id,
      society_id,
      amount,
      total_amount,
      period_start,
      period_end,
      due_date
    ) VALUES (
      v_subscription.subscription_id,
      v_subscription.society_id,
      v_subscription.monthly_fee,
      v_subscription.monthly_fee,
      v_subscription.current_period_start,
      v_subscription.current_period_end,
      CURRENT_DATE + INTERVAL '7 days'
    );

    -- Update subscription
    UPDATE society_subscriptions
    SET
      current_period_start = current_period_end + INTERVAL '1 day',
      current_period_end = current_period_end + INTERVAL '1 month',
      next_billing_date = next_billing_date + INTERVAL '1 month',
      updated_at = NOW()
    WHERE subscription_id = v_subscription.subscription_id;

    v_count := v_count + 1;
  END LOOP;

  RETURN v_count;
END;
$$ LANGUAGE plpgsql;
```

---

## 13. Row Level Security (RLS)

### 13.1 Enable RLS

```sql
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE residents ENABLE ROW LEVEL SECURITY;
ALTER TABLE vendors ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_workflow_progress ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE ratings ENABLE ROW LEVEL SECURITY;
ALTER TABLE disputes ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;
```

### 13.2 RLS Policies

**Users:**
```sql
CREATE POLICY "Users can view own profile"
ON users FOR SELECT
USING (auth.uid() = user_id);
```

**Orders - Resident access:**
```sql
CREATE POLICY "Residents view own orders"
ON orders FOR SELECT
USING (auth.uid() = resident_id);
```

**Orders - Vendor access:**
```sql
CREATE POLICY "Vendors view assigned orders"
ON orders FOR SELECT
USING (auth.uid() = vendor_id);
```

**Workflow Progress - Vendor access:**
```sql
CREATE POLICY "Vendors view/update workflow for their orders"
ON order_workflow_progress FOR ALL
USING (
  EXISTS (
    SELECT 1 FROM orders
    WHERE orders.order_id = order_workflow_progress.order_id
      AND orders.vendor_id = auth.uid()
  )
);
```

---

## 14. Sample Data

### 14.1 Complete Sample Dataset

```sql
-- Create sample users
INSERT INTO users (phone, full_name, user_type, is_verified) VALUES
  ('+919876543210', 'Ramesh Kumar', 'RESIDENT', true),
  ('+919876543211', 'Priya Sharma', 'VENDOR', true),
  ('+919876543212', 'Amit Verma', 'RESIDENT', true),
  ('+919876543213', 'Perfect Press Owner', 'VENDOR', true),
  ('+919876543214', 'Admin User', 'ADMIN', true);

-- Create societies
INSERT INTO societies (name, address, city, state, pincode, total_flats, status) VALUES
  ('Maple Gardens', '123 MG Road, Koramangala', 'Bangalore', 'Karnataka', '560034', 250, 'ACTIVE'),
  ('Palm Residency', '456 Anna Salai, T Nagar', 'Chennai', 'Tamil Nadu', '600017', 180, 'ACTIVE');

-- Create residents
INSERT INTO residents (resident_id, society_id, flat_number, verification_status)
SELECT user_id, 1, 'A-404', 'VERIFIED'
FROM users WHERE phone = '+919876543210';

-- Create vendors
INSERT INTO vendors (vendor_id, business_name, store_address, approval_status)
SELECT user_id, 'Perfect Press', '789 Market Street, Koramangala', 'APPROVED'
FROM users WHERE phone = '+919876543211';

-- Map vendors to societies
INSERT INTO vendor_societies (vendor_id, society_id, approval_status)
SELECT v.vendor_id, 1, 'APPROVED'
FROM vendors v WHERE v.business_name = 'Perfect Press';

-- Vendor offers services
INSERT INTO vendor_services (vendor_id, service_id, turnaround_hours)
SELECT
  v.vendor_id,
  sc.service_id,
  sc.default_turnaround_hours
FROM vendors v
CROSS JOIN service_categories sc
WHERE v.business_name = 'Perfect Press'
  AND sc.service_key IN ('IRONING', 'WASHING_IRONING');

-- Initialize workflows for offered services
-- (Workflows already created in section 6.3)
```

---

## Summary

### Database Statistics

**Total Tables:** 30+ core tables
**Total Indexes:** 100+ optimized indexes
**Total Functions:** 15+ helper functions
**Total Triggers:** 5+ automation triggers
**RLS Policies:** 20+ security policies

### Key Features Implemented

✅ Hierarchical category structure (multi-service platform)
✅ Mixed-service orders (multiple service types per order)
✅ **Custom workflow configuration per service type**
✅ **Workflow step tracking for each service in each order**
✅ **Order_workflow_progress table tracks individual step completion**
✅ Society subscription billing with automated invoicing
✅ Multi-tenancy with RLS policies
✅ Comprehensive audit trails
✅ Vendor terminology (no more LP_)
✅ Full-text search capabilities

### Workflow System

**5 Service Types with Custom Workflows:**
- Ironing: 5 steps (Pickup → Count → Iron → QC → Deliver)
- Dry Cleaning: 7 steps (Pickup → Count → Pre-treat → Dry Clean → QC → Press → Deliver)
- Car Wash: 6 steps (Arrive → Inspect → Wash → Vacuum → Polish → Final Check)
- Gardening: 6 steps (Arrive → Trim → Mow → Weed → Cleanup → Inspect)
- Plumbing: 5 steps (Arrive → Diagnose → Repair → Test → Final Check)

**Each workflow step can have:**
- Required/optional flag
- Photo requirement
- Signature requirement
- Estimated duration
- Order status mapping

---

**End of Document**
