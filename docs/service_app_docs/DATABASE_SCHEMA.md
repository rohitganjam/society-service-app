# Database Schema Documentation

**Version:** 3.0
**Date:** November 20, 2025
**Database:** PostgreSQL 15 (Supabase)
**Architecture:** Multi-category service platform with hierarchical structure, custom workflows, and multi-society support

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
- ✅ **Multi-society support (one user, multiple residences)**
- ✅ **Independent house support (apartments + layouts)**
- ✅ **Society roster for instant verification**
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
- Primary keys: `INTEGER GENERATED ALWAYS AS IDENTITY` for auto-increment integers (SQL standard, PostgreSQL 10+), `UUID` for distributed records
- Legacy: `SERIAL` is shorthand but `GENERATED ALWAYS AS IDENTITY` is preferred for better control
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

**Key Updates in Version 3.0:**
- ✅ Multi-society support: One user can have residences in multiple societies
- ✅ Independent house support: Societies can be APARTMENT or LAYOUT type
- ✅ **Hierarchical society structure**: Buildings/blocks for apartments, phases for layouts
- ✅ **Multi-level vendor assignments**: Assign vendors to entire society, specific buildings, or specific phases
- ✅ **Smart vendor filtering**: Default vendor visibility based on resident's building/phase, with override option
- ✅ Society roster table: Pre-approved residents for instant verification
- ✅ Residents table redesigned: Allows multiple residences per user with is_primary and is_active flags
- ✅ Support for multiple households per house: Different floors in same house number

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
│         USER MANAGEMENT, MULTI-SOCIETY & UNIFIED 4-LEVEL HIERARCHY       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────┐                              ┌──────────────┐         │
│  │   users      │ 1                      ∞     │  residents   │         │
│  │ ─────────────│◄─────────────────────────────┤ ─────────────│         │
│  │ • user_id (PK)                              │ • resident_id│         │
│  │ • phone      │                              │   (PK)       │         │
│  │ • user_type  │                              │ • user_id(FK)│         │
│  │ • is_verified│                              │ • society_id │         │
│  └──────────────┘                              │   (FK)       │         │
│         │                                       │ • group_id   │         │
│         │                                       │   (FK)       │         │
│         │                                       │ • unit_type  │         │
│         │                                       │ • flat_number│         │
│         │ 1                                     │ • house_num  │         │
│         │                                       │ • floor      │         │
│         ▼ ∞                                     │ • is_primary │         │
│  ┌──────────────┐                              │ • is_active  │         │
│  │   vendors    │                              │ • verification         │
│  │ ─────────────│                              │   _status    │         │
│  │ • vendor_id  │                              └──────┬───────┘         │
│  │   (PK, FK)   │                                     │                 │
│  │ • business_  │                                     │ ∞               │
│  │   name       │                                     │                 │
│  │ • store_addr │                                     │ 1               │
│  └──────┬───────┘                                     ▼                 │
│         │                                       ┌──────────────┐         │
│         │ 1                                     │  societies   │         │
│         │                                       │ ─────────────│         │
│         ▼ ∞                                     │ • society_id │         │
│  ┌──────────────┐                              │   (PK)       │         │
│  │ vendor_      │                              │ • name       │         │
│  │  services    │                              │ • society_   │         │
│  │ ─────────────│                              │   type       │         │
│  │ • vendor_id  │                              │ • total_flats│         │
│  │ • service_id │                              │ • total_     │         │
│  │   (FK)       │                              │   houses     │         │
│  │ • turnaround │                              └───┬───┬──────┘         │
│  └──────┬───────┘                                  │   │                │
│         │                                      1   │   │ 1              │
│         │ 1                                   ┌────┘   └────┐           │
│         │                                     │             │           │
│         ▼ ∞                                   ▼ ∞           ▼ ∞         │
│  ┌──────────────┐            ┌────────────────┐     ┌───────────────┐  │
│  │ rate_cards   │            │ society_       │     │ society_groups│  │
│  │ ─────────────│            │  roster        │     │ ──────────────│  │
│  │ • rate_card_id            │ ─────────────  │     │ • group_id(PK)│  │
│  │ • vendor_id  │            │ • roster_id    │     │ • society_id  │  │
│  │ • society_id │            │ • society_id   │     │   (FK)        │  │
│  │   (FK)       │            │   (FK)         │     │ • group_name  │  │
│  │ • is_published            │ • group_id     │     │ • group_type  │  │
│  └──────┬───────┘            │   (FK)         │     │ • group_code  │  │
│         │                     │ • phone        │     │ • total_units │  │
│         │ 1                   │ • unit_type    │     │ • total_floors│  │
│         │                     │ • flat_number  │     └───────┬───────┘  │
│         ▼ ∞                   │ • house_num    │             │          │
│  ┌──────────────┐            │ • floor        │             │ ∞        │
│  │ rate_card_   │            └────────────────┘             │          │
│  │   items      │                                      1    ▼          │
│  │ ─────────────│                              ┌──────────────────┐    │
│  │ • item_id    │                              │ vendor_service_  │    │
│  │ • rate_card_id                              │   areas          │    │
│  │   (FK)       │      ┌───────────────────────┤ ─────────────────│    │
│  │ • service_id │      │                       │ • assignment_id  │    │
│  │ • item_name  │      │ ∞                     │   (PK)           │    │
│  │ • price      │      │                       │ • vendor_id (FK) │    │
│  └──────────────┘      │                       │ • society_id(FK) │    │
│                         │                       │ • assignment_    │    │
│  ┌──────────────┐      │                       │   type           │    │
│  │  vendors     │──────┘                       │ • group_id       │    │
│  │ ─────────────│ 1                            │   (FK, nullable) │    │
│  │ • vendor_id  │                              │ • is_active      │    │
│  │   (PK)       │                              └──────────────────┘    │
│  └──────────────┘                                                      │
│                                                                          │
│  **Unified 4-Level Hierarchy:**                                         │
│  - Society → Groups (Buildings/Phases) → Units (Flats/Houses) → Floors  │
│  - Single society_groups table for both apartments and layouts          │
│  - group_type: BUILDING, TOWER, BLOCK, WING, PHASE, SECTION, ZONE      │
│                                                                          │
│  **Vendor Assignment (Simplified):**                                    │
│  - SOCIETY: Vendor serves entire society (all groups)                   │
│  - GROUP: Vendor assigned to specific group(s) - works for both         │
│    buildings (apartments) and phases (layouts)                          │
│  - Residents linked to groups via group_id                              │
│  - Default vendor filtering based on resident's group                   │
│  - Residents can override to view all vendors in society                │
│                                                                          │
│  **Notes:**                                                              │
│  - residents table supports multi-society (one user, multiple residences)│
│  - Only one is_primary per user, only one is_active per user (context)  │
│  - Supports FLAT and HOUSE unit types uniformly via group_id            │
│  - Multiple households per unit supported via different floor values    │
│  - society_roster includes group_id for instant verification            │
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

**Purpose:** Housing societies/complexes (apartments and independent house layouts)

```sql
CREATE TABLE societies (
  society_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  society_type VARCHAR(20) DEFAULT 'APARTMENT'
    CHECK (society_type IN ('APARTMENT', 'LAYOUT')),
  address TEXT NOT NULL,
  city VARCHAR(100) NOT NULL,
  state VARCHAR(100) NOT NULL,
  pincode VARCHAR(10) NOT NULL,

  -- Contact
  contact_person VARCHAR(255),
  contact_phone VARCHAR(15),
  contact_email VARCHAR(255),

  -- Stats - for apartments
  total_flats INTEGER,
  occupied_flats INTEGER,

  -- Stats - for layouts/independent houses
  total_houses INTEGER,
  occupied_houses INTEGER,

  -- Status
  status VARCHAR(20) DEFAULT 'PENDING'
    CHECK (status IN ('PENDING', 'ACTIVE', 'SUSPENDED', 'INACTIVE')),
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  created_by UUID REFERENCES users(user_id),

  -- Constraints: Must have either flats or houses based on type
  CHECK (
    (society_type = 'APARTMENT' AND total_flats IS NOT NULL) OR
    (society_type = 'LAYOUT' AND total_houses IS NOT NULL)
  )
);

-- Indexes
CREATE INDEX idx_societies_status ON societies(status);
CREATE INDEX idx_societies_city ON societies(city);
CREATE INDEX idx_societies_pincode ON societies(pincode);
CREATE INDEX idx_societies_type ON societies(society_type);
CREATE INDEX idx_societies_active ON societies(is_active) WHERE is_active = true;

-- Full-text search (includes pincode for better search)
CREATE INDEX idx_societies_search ON societies USING gin(
  to_tsvector('english', name || ' ' || address || ' ' || pincode)
);
```

**Notes:**
- `society_type`: 'APARTMENT' for multi-unit buildings, 'LAYOUT' for independent houses
- For apartments: `total_flats` and `occupied_flats` are used
- For layouts: `total_houses` and `occupied_houses` are used
- Constraint ensures appropriate stats are populated based on society type

---

### 4.3 Residents Table

**Purpose:** Resident-society relationships (supports multi-society membership and independent houses)

```sql
CREATE TABLE residents (
  resident_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  society_id INTEGER NOT NULL REFERENCES societies(society_id) ON DELETE CASCADE,

  -- Unit type
  unit_type VARCHAR(10) NOT NULL CHECK (unit_type IN ('FLAT', 'HOUSE')),

  -- For apartments (FLAT)
  flat_number VARCHAR(20),
  tower VARCHAR(10),

  -- For independent houses (HOUSE)
  house_number VARCHAR(20),
  street VARCHAR(100),

  -- Common fields
  floor INTEGER,
  notes TEXT,

  -- Multi-society support
  is_primary BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT false,

  -- Preferences
  preferred_pickup_time TIME,
  default_pickup_address TEXT,

  -- Status
  verification_status VARCHAR(20) DEFAULT 'PENDING'
    CHECK (verification_status IN ('PENDING', 'VERIFIED', 'REJECTED')),
  rejection_reason TEXT,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  verified_at TIMESTAMP,
  verified_by UUID REFERENCES users(user_id),

  -- Constraints
  CHECK (
    (unit_type = 'FLAT' AND flat_number IS NOT NULL) OR
    (unit_type = 'HOUSE' AND house_number IS NOT NULL)
  ),

  -- Allow multiple households in same unit (different floors)
  -- Remove the old UNIQUE constraint, add these instead:
  UNIQUE(society_id, unit_type, flat_number, tower, floor),
  UNIQUE(society_id, unit_type, house_number, street, floor)
);

-- Indexes
CREATE INDEX idx_residents_user ON residents(user_id);
CREATE INDEX idx_residents_society ON residents(society_id);
CREATE INDEX idx_residents_status ON residents(verification_status);
CREATE INDEX idx_residents_primary ON residents(user_id, is_primary) WHERE is_primary = true;
CREATE INDEX idx_residents_active ON residents(user_id, is_active) WHERE is_active = true;
CREATE INDEX idx_residents_unit_type ON residents(unit_type);

-- Composite index for flat lookup
CREATE INDEX idx_residents_flat_lookup ON residents(society_id, flat_number, tower)
  WHERE unit_type = 'FLAT';

-- Composite index for house lookup
CREATE INDEX idx_residents_house_lookup ON residents(society_id, house_number, street)
  WHERE unit_type = 'HOUSE';

-- Most common query: get user's active verified residence
CREATE INDEX idx_residents_user_active_verified ON residents(user_id)
  WHERE is_active = true AND verification_status = 'VERIFIED';
```

**Notes:**
- Changed from `resident_id UUID PRIMARY KEY` to `resident_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY` to allow multiple residences per user
- `user_id`: Links to users table (one user can have multiple residences)
- `unit_type`: 'FLAT' for apartments, 'HOUSE' for independent houses
- `is_primary`: User's main residence (only one can be true per user)
- `is_active`: Currently selected society context (only one can be true per user)
- Multiple households in same house: Same `house_number` but different `floor` values
- UNIQUE constraints allow multiple floors in same flat/house number
- Triggers (defined later) ensure only one primary and one active per user

---

### 4.4 Society Roster Table

**Purpose:** Pre-approved resident lists for instant verification during onboarding

```sql
CREATE TABLE society_roster (
  roster_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  society_id INTEGER NOT NULL REFERENCES societies(society_id) ON DELETE CASCADE,
  phone VARCHAR(15) NOT NULL,
  resident_name VARCHAR(255),

  -- Unit information
  unit_type VARCHAR(10) NOT NULL CHECK (unit_type IN ('FLAT', 'HOUSE')),

  -- For apartments
  flat_number VARCHAR(20),
  tower VARCHAR(10),

  -- For independent houses
  house_number VARCHAR(20),
  street VARCHAR(100),

  -- Common
  floor INTEGER,
  notes TEXT,

  -- Status
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  added_at TIMESTAMP DEFAULT NOW(),
  added_by UUID REFERENCES users(user_id),
  updated_at TIMESTAMP DEFAULT NOW(),

  -- Constraints
  CHECK (
    (unit_type = 'FLAT' AND flat_number IS NOT NULL) OR
    (unit_type = 'HOUSE' AND house_number IS NOT NULL)
  )
);

-- Indexes
CREATE INDEX idx_roster_phone ON society_roster(phone);
CREATE INDEX idx_roster_society ON society_roster(society_id);
CREATE INDEX idx_roster_active ON society_roster(is_active) WHERE is_active = true;

-- Most common query: check if phone exists in roster
CREATE INDEX idx_roster_phone_active ON society_roster(phone, is_active)
  WHERE is_active = true;

-- Lookup by phone and society
CREATE INDEX idx_roster_lookup ON society_roster(phone, society_id)
  WHERE is_active = true;
```

**Notes:**
- Uploaded by society admins to pre-approve residents
- One phone can appear multiple times (multiple societies or multiple floors in same society)
- Enables instant verification during resident onboarding
- `is_active`: Allows soft-delete of roster entries without removing data
- No UNIQUE constraint on phone - allows same phone in multiple societies/units

**Example Use Cases:**
1. Family with multiple properties: Same phone in roster for 2 different societies
2. Multi-floor household: Same phone, same house_number, different floors
3. Joint ownership: Multiple phones for same flat/house

---

### 4.5 Vendors Table

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
  id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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

### 4.6 Society Groups Table (Unified 4-Level Hierarchy)

**Purpose:** Unified grouping structure supporting 4-level hierarchy: Society → Groups → Units → Floors

```sql
CREATE TABLE society_groups (
  group_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  society_id INTEGER NOT NULL REFERENCES societies(society_id) ON DELETE CASCADE,
  group_name VARCHAR(100) NOT NULL,
  group_code VARCHAR(20),
  group_type VARCHAR(20) NOT NULL
    CHECK (group_type IN ('BUILDING', 'BLOCK', 'TOWER', 'WING', 'PHASE', 'SECTION', 'ZONE')),
  description TEXT,

  -- Stats (applicable based on society_type)
  total_units INTEGER,      -- Total flats for apartments OR total houses for layouts
  total_floors INTEGER,     -- For buildings/towers

  -- Status
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  created_by UUID REFERENCES users(user_id),

  UNIQUE(society_id, group_name)
);

-- Indexes
CREATE INDEX idx_groups_society ON society_groups(society_id);
CREATE INDEX idx_groups_type ON society_groups(group_type);
CREATE INDEX idx_groups_active ON society_groups(is_active) WHERE is_active = true;
```

**Notes:**
- **Unified table** for both apartment and layout societies
- **Flexible naming:** Supports "Building", "Phase", "Tower", "Block", "Wing", "Section", "Zone"
- **For Apartments:** `group_type` = 'BUILDING', 'BLOCK', 'TOWER', or 'WING'
- **For Layouts:** `group_type` = 'PHASE', 'SECTION', or 'ZONE'
- `total_units`: Number of flats (apartments) OR houses (layouts) in this group
- `total_floors`: Only applicable for multi-story buildings

**4-Level Hierarchy Examples:**

**Apartments:**
```
Society → Building A → Flat A-101 → Floor 1
Society → Building A → Flat A-101 → Floor 2
Society → Tower B → Flat B-205 → (no floors - single household)
```

**Layouts:**
```
Society → Phase 1 → House #101 → Ground Floor
Society → Phase 1 → House #101 → First Floor
Society → Phase 2 → House #205 → (no floors - single household)
```

**Example Data:**
```sql
-- Apartment society
INSERT INTO society_groups (society_id, group_name, group_code, group_type, total_units, total_floors)
VALUES
  (1, 'Building A', 'A', 'BUILDING', 60, 15),
  (1, 'Tower B', 'B', 'TOWER', 80, 20);

-- Layout society
INSERT INTO society_groups (society_id, group_name, group_code, group_type, total_units)
VALUES
  (2, 'Phase 1', 'P1', 'PHASE', 50),
  (2, 'East Section', 'ES', 'SECTION', 35);
```

---

### 4.7 Vendor Service Area Assignments Table

**Purpose:** Define which groups (buildings/phases) a vendor serves (simplified with unified groups)

```sql
CREATE TABLE vendor_service_areas (
  assignment_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  vendor_id UUID NOT NULL REFERENCES vendors(vendor_id) ON DELETE CASCADE,
  society_id INTEGER NOT NULL REFERENCES societies(society_id) ON DELETE CASCADE,

  -- Assignment level
  assignment_type VARCHAR(20) NOT NULL
    CHECK (assignment_type IN ('SOCIETY', 'GROUP')),

  -- Reference ID (nullable when assignment_type = 'SOCIETY')
  group_id INTEGER REFERENCES society_groups(group_id) ON DELETE CASCADE,

  -- Status
  is_active BOOLEAN DEFAULT true,

  -- Metadata
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  assigned_by UUID REFERENCES users(user_id),

  -- Constraints
  CHECK (
    (assignment_type = 'SOCIETY' AND group_id IS NULL) OR
    (assignment_type = 'GROUP' AND group_id IS NOT NULL)
  ),

  -- Prevent duplicate assignments
  UNIQUE(vendor_id, society_id, assignment_type, group_id)
);

-- Indexes
CREATE INDEX idx_vendor_areas_vendor ON vendor_service_areas(vendor_id);
CREATE INDEX idx_vendor_areas_society ON vendor_service_areas(society_id);
CREATE INDEX idx_vendor_areas_group ON vendor_service_areas(group_id) WHERE group_id IS NOT NULL;
CREATE INDEX idx_vendor_areas_type ON vendor_service_areas(assignment_type);
CREATE INDEX idx_vendor_areas_active ON vendor_service_areas(is_active) WHERE is_active = true;

-- Composite index for vendor lookup in a society
CREATE INDEX idx_vendor_areas_lookup ON vendor_service_areas(society_id, vendor_id, is_active)
  WHERE is_active = true;
```

**Notes:**
- **Simplified Assignment Types:**
  - `SOCIETY`: Vendor serves entire society (all groups)
  - `GROUP`: Vendor assigned to specific group(s) - works for both buildings and phases
- One vendor can have multiple group assignments (e.g., Building A + Building B)
- Used for default vendor filtering in resident app
- Residents can override and view all vendors if needed

**Example Assignments:**
```sql
-- Vendor serves entire society
INSERT INTO vendor_service_areas (vendor_id, society_id, assignment_type)
VALUES ('vendor-uuid', 1, 'SOCIETY');

-- Vendor serves specific groups (works for both buildings and phases)
INSERT INTO vendor_service_areas (vendor_id, society_id, assignment_type, group_id)
VALUES
  ('vendor-uuid', 1, 'GROUP', 1),  -- Building A
  ('vendor-uuid', 1, 'GROUP', 2);  -- Building B

-- Vendor serves specific phases in layout
INSERT INTO vendor_service_areas (vendor_id, society_id, assignment_type, group_id)
VALUES ('vendor-uuid', 2, 'GROUP', 5);  -- Phase 1
```

---

### 4.8 Updated Residents Table for Group References

**Purpose:** Link residents to their groups (buildings or phases) in the 4-level hierarchy

Update the residents table schema to include group reference:

```sql
-- Add new column to residents table
ALTER TABLE residents
ADD COLUMN group_id INTEGER REFERENCES society_groups(group_id) ON DELETE SET NULL;

-- Add check constraint to ensure group_id is set for both FLAT and HOUSE
ALTER TABLE residents
ADD CONSTRAINT residents_group_check CHECK (
  group_id IS NOT NULL
);

-- Add index
CREATE INDEX idx_residents_group ON residents(group_id);
```

**Notes:**
- **4-Level Hierarchy:** Society → Group → Unit (flat/house) → Floor (optional)
- `group_id`: References the building (for apartments) or phase (for layouts)
- `flat_number` or `house_number`: The unit within the group
- `floor`: Optional - for multi-floor households
- Works uniformly for both apartment and layout societies

**Example Data:**
```sql
-- Apartment resident
INSERT INTO residents (user_id, society_id, group_id, unit_type, flat_number, floor)
VALUES ('user-uuid', 1, 1, 'FLAT', 'A-101', 1);
-- Represents: Society 1 → Building A (group_id: 1) → Flat A-101 → Floor 1

-- Layout resident
INSERT INTO residents (user_id, society_id, group_id, unit_type, house_number, floor)
VALUES ('user-uuid', 2, 5, 'HOUSE', '101', 0);
-- Represents: Society 2 → Phase 1 (group_id: 5) → House 101 → Ground Floor
```

---

### 4.9 Updated Society Roster Table for Group References

**Purpose:** Include group information in pre-approved roster (4-level hierarchy)

Update the society_roster table schema:

```sql
-- Add new column to society_roster table
ALTER TABLE society_roster
ADD COLUMN group_id INTEGER REFERENCES society_groups(group_id) ON DELETE SET NULL;

-- Add check constraint
ALTER TABLE society_roster
ADD CONSTRAINT roster_group_check CHECK (
  group_id IS NOT NULL
);

-- Add index
CREATE INDEX idx_roster_group ON society_roster(group_id);
```

**Notes:**
- **Unified approach:** Works for both apartments and layouts
- `group_id`: References building (apartments) or phase (layouts)
- Used for instant resident verification during onboarding
- Matches the same structure as the residents table

---

## 5. Category & Service Tables

### 5.1 Parent Categories Table

**Purpose:** Top-level service categories (Laundry, Vehicle, Home, Personal)

```sql
CREATE TABLE parent_categories (
  category_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  service_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  rate_card_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  item_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  template_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  step_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  settlement_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  subscription_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  invoice_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  dispute_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  rating_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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
  notification_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
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

### 12.5 Check Resident in Roster

```sql
CREATE OR REPLACE FUNCTION check_resident_in_roster(p_phone VARCHAR(15))
RETURNS TABLE(
  society_id INTEGER,
  society_name VARCHAR(255),
  society_type VARCHAR(20),
  address TEXT,
  city VARCHAR(100),
  unit_type VARCHAR(10),
  flat_number VARCHAR(20),
  tower VARCHAR(10),
  house_number VARCHAR(20),
  street VARCHAR(100),
  floor INTEGER,
  suggested_name VARCHAR(255),
  notes TEXT
) AS $$
BEGIN
  RETURN QUERY
  SELECT
    s.society_id,
    s.name as society_name,
    s.society_type,
    s.address,
    s.city,
    sr.unit_type,
    sr.flat_number,
    sr.tower,
    sr.house_number,
    sr.street,
    sr.floor,
    sr.resident_name as suggested_name,
    sr.notes
  FROM society_roster sr
  JOIN societies s ON sr.society_id = s.society_id
  WHERE sr.phone = p_phone
    AND sr.is_active = true
    AND s.is_active = true
    AND s.status = 'ACTIVE'
  ORDER BY
    -- Primary residence first (if marked in roster notes)
    CASE WHEN sr.notes ILIKE '%primary%' THEN 0 ELSE 1 END,
    sr.added_at ASC;
END;
$$ LANGUAGE plpgsql;
```

**Usage:**
```sql
SELECT * FROM check_resident_in_roster('+919876543210');
```

### 12.6 Set Active Society

```sql
CREATE OR REPLACE FUNCTION set_active_society(
  p_user_id UUID,
  p_society_id INTEGER
)
RETURNS BOOLEAN AS $$
DECLARE
  v_count INTEGER;
BEGIN
  -- Check if user has verified residence in this society
  SELECT COUNT(*) INTO v_count
  FROM residents
  WHERE user_id = p_user_id
    AND society_id = p_society_id
    AND verification_status = 'VERIFIED';

  IF v_count = 0 THEN
    RAISE EXCEPTION 'User not verified in this society';
  END IF;

  -- Deactivate all other residences for this user
  UPDATE residents
  SET is_active = false,
      updated_at = NOW()
  WHERE user_id = p_user_id
    AND society_id != p_society_id;

  -- Activate the selected society
  UPDATE residents
  SET is_active = true,
      updated_at = NOW()
  WHERE user_id = p_user_id
    AND society_id = p_society_id;

  RETURN true;
END;
$$ LANGUAGE plpgsql;
```

**Usage:**
```sql
SELECT set_active_society('user-uuid', 3);
```

### 12.7 Get User Active Society

```sql
CREATE OR REPLACE FUNCTION get_user_active_society(p_user_id UUID)
RETURNS TABLE(
  resident_id INTEGER,
  society_id INTEGER,
  society_name VARCHAR(255),
  society_type VARCHAR(20),
  unit_type VARCHAR(10),
  flat_number VARCHAR(20),
  house_number VARCHAR(20),
  floor INTEGER
) AS $$
BEGIN
  RETURN QUERY
  SELECT
    r.resident_id,
    r.society_id,
    s.name as society_name,
    s.society_type,
    r.unit_type,
    r.flat_number,
    r.house_number,
    r.floor
  FROM residents r
  JOIN societies s ON r.society_id = s.society_id
  WHERE r.user_id = p_user_id
    AND r.is_active = true
    AND r.verification_status = 'VERIFIED'
  LIMIT 1;
END;
$$ LANGUAGE plpgsql;
```

**Usage:**
```sql
SELECT * FROM get_user_active_society('user-uuid');
```

### 12.8 Get User All Residences

```sql
CREATE OR REPLACE FUNCTION get_user_all_residences(p_user_id UUID)
RETURNS TABLE(
  resident_id INTEGER,
  society_id INTEGER,
  society_name VARCHAR(255),
  society_type VARCHAR(20),
  address TEXT,
  city VARCHAR(100),
  unit_type VARCHAR(10),
  flat_number VARCHAR(20),
  tower VARCHAR(10),
  house_number VARCHAR(20),
  street VARCHAR(100),
  floor INTEGER,
  notes TEXT,
  is_primary BOOLEAN,
  is_active BOOLEAN,
  verification_status VARCHAR(20),
  verified_at TIMESTAMP
) AS $$
BEGIN
  RETURN QUERY
  SELECT
    r.resident_id,
    r.society_id,
    s.name as society_name,
    s.society_type,
    s.address,
    s.city,
    r.unit_type,
    r.flat_number,
    r.tower,
    r.house_number,
    r.street,
    r.floor,
    r.notes,
    r.is_primary,
    r.is_active,
    r.verification_status,
    r.verified_at
  FROM residents r
  JOIN societies s ON r.society_id = s.society_id
  WHERE r.user_id = p_user_id
  ORDER BY
    r.is_primary DESC,
    r.is_active DESC,
    r.verified_at DESC NULLS LAST;
END;
$$ LANGUAGE plpgsql;
```

**Usage:**
```sql
SELECT * FROM get_user_all_residences('user-uuid');
```

---

## 13. Triggers for Multi-Society Business Logic

### 13.1 Ensure Only One Primary Residence

```sql
CREATE OR REPLACE FUNCTION enforce_single_primary_residence()
RETURNS TRIGGER AS $$
BEGIN
  -- If setting this residence as primary
  IF NEW.is_primary = true THEN
    -- Remove primary flag from all other residences for this user
    UPDATE residents
    SET is_primary = false,
        updated_at = NOW()
    WHERE user_id = NEW.user_id
      AND resident_id != COALESCE(NEW.resident_id, -1);
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_enforce_single_primary
  BEFORE INSERT OR UPDATE ON residents
  FOR EACH ROW
  EXECUTE FUNCTION enforce_single_primary_residence();
```

### 13.2 Ensure Only One Active Society

```sql
CREATE OR REPLACE FUNCTION enforce_single_active_society()
RETURNS TRIGGER AS $$
BEGIN
  -- If setting this residence as active
  IF NEW.is_active = true THEN
    -- Deactivate all other residences for this user
    UPDATE residents
    SET is_active = false,
        updated_at = NOW()
    WHERE user_id = NEW.user_id
      AND resident_id != COALESCE(NEW.resident_id, -1);
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_enforce_single_active
  BEFORE INSERT OR UPDATE ON residents
  FOR EACH ROW
  EXECUTE FUNCTION enforce_single_active_society();
```

### 13.3 Auto-set First Residence Flags

```sql
CREATE OR REPLACE FUNCTION set_first_residence_flags()
RETURNS TRIGGER AS $$
DECLARE
  v_count INTEGER;
BEGIN
  -- Count existing residences for this user
  SELECT COUNT(*) INTO v_count
  FROM residents
  WHERE user_id = NEW.user_id;

  -- If this is the first residence, make it primary and active
  IF v_count = 0 THEN
    NEW.is_primary := true;
    NEW.is_active := true;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_first_residence_flags
  BEFORE INSERT ON residents
  FOR EACH ROW
  EXECUTE FUNCTION set_first_residence_flags();
```

---

## 14. Row Level Security (RLS)

### 14.1 Enable RLS

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

### 14.2 RLS Policies

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

## 15. Sample Data

### 15.1 Complete Sample Dataset

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

## 16. Data Flow Examples & Storage Patterns

This section provides comprehensive examples of how data flows through the system and how relationships are stored across tables.

---

### 16.1 User Registration & Authentication Flow

#### 16.1.1 Resident Registration (From Roster)

**Scenario:** Ramesh Kumar registers for Maple Gardens apartment A-404. He's already in the society roster.

**Step 1: Check Roster**
```sql
-- API: POST /api/v1/onboarding/resident/check-roster
SELECT
  sr.society_id,
  s.society_name,
  s.society_type,
  sr.unit_type,
  sr.flat_number,
  sr.house_number,
  sr.floor
FROM society_roster sr
JOIN societies s ON s.society_id = sr.society_id
WHERE sr.phone = '+919876543210'
  AND sr.is_active = true;

-- Result: Found in roster for Maple Gardens, A-404
```

**Step 2: Register User**
```sql
-- API: POST /api/v1/onboarding/resident/register
-- Creates records in TWO tables:

-- 1. users table (user account)
INSERT INTO users (user_id, phone, full_name, email, user_type, is_verified)
VALUES (
  'uuid-ramesh-001',
  '+919876543210',
  'Ramesh Kumar',
  'ramesh@example.com',
  'RESIDENT',
  true  -- Auto-verified because from roster
);

-- 2. residents table (residence information)
INSERT INTO residents (
  user_id, society_id, unit_type, flat_number, tower, floor,
  is_primary, is_active, verification_status
)
VALUES (
  'uuid-ramesh-001',
  1,  -- Maple Gardens
  'FLAT',
  'A-404',
  'A',
  4,
  true,   -- First residence is primary
  true,   -- First residence is active
  'VERIFIED'  -- Auto-verified from roster
);
```

**Storage Pattern:**
```
users table:
┌──────────────────┬────────────────┬──────────────┬───────────────┬─────────────┐
│ user_id          │ phone          │ full_name    │ user_type     │ is_verified │
├──────────────────┼────────────────┼──────────────┼───────────────┼─────────────┤
│ uuid-ramesh-001  │ +919876543210  │ Ramesh Kumar │ RESIDENT      │ true        │
└──────────────────┴────────────────┴──────────────┴───────────────┴─────────────┘

residents table:
┌───────────┬──────────────────┬────────────┬─────────────┬────────────┬────────────┬──────────┐
│resident_id│ user_id          │ society_id │ flat_number │ is_primary │ is_active  │ verified │
├───────────┼──────────────────┼────────────┼─────────────┼────────────┼────────────┼──────────┤
│ 1         │ uuid-ramesh-001  │ 1          │ A-404       │ true       │ true       │ VERIFIED │
└───────────┴──────────────────┴────────────┴─────────────┴────────────┴────────────┴──────────┘
```

---

#### 16.1.2 Multi-Society Registration

**Scenario:** Ramesh later adds a weekend home in Palm Residency

**API Call:** `POST /api/v1/residents/add-residence`

```sql
-- Check roster for second society
SELECT * FROM society_roster
WHERE phone = '+919876543210'
  AND society_id = 2;  -- Palm Residency

-- Add second residence (same user_id, different society)
INSERT INTO residents (
  user_id, society_id, unit_type, flat_number,
  is_primary, is_active, verification_status
)
VALUES (
  'uuid-ramesh-001',  -- Same user_id
  2,  -- Palm Residency
  'FLAT',
  '201',
  false,  -- Not primary
  false,  -- Not active (Maple Gardens is still active)
  'VERIFIED'
);
```

**Storage Pattern (Multi-Society):**
```
users table (unchanged):
┌──────────────────┬────────────────┬──────────────┐
│ user_id          │ phone          │ full_name    │
├──────────────────┼────────────────┼──────────────┤
│ uuid-ramesh-001  │ +919876543210  │ Ramesh Kumar │
└──────────────────┴────────────────┴──────────────┘

residents table (now has 2 rows for same user):
┌───────────┬──────────────────┬────────────┬──────────────────┬─────────────┬────────────┬──────────┐
│resident_id│ user_id          │ society_id │ flat_number      │ is_primary  │ is_active  │ verified │
├───────────┼──────────────────┼────────────┼──────────────────┼─────────────┼────────────┼──────────┤
│ 1         │ uuid-ramesh-001  │ 1          │ A-404 (Maple)    │ true        │ true       │ VERIFIED │
│ 2         │ uuid-ramesh-001  │ 2          │ 201 (Palm)       │ false       │ false      │ VERIFIED │
└───────────┴──────────────────┴────────────┴──────────────────┴─────────────┴────────────┴──────────┘
```

**Key Point:** One `user_id` can have multiple rows in `residents` table (one per society)

---

#### 16.1.3 Switching Active Society

**Scenario:** Ramesh switches to his Palm Residency home for the weekend

**API Call:** `POST /api/v1/residents/switch-society`

```sql
-- Function handles the switch with triggers
SELECT set_active_society('uuid-ramesh-001', 2);

-- What happens internally:
-- 1. Set all residences for this user to is_active = false
UPDATE residents
SET is_active = false
WHERE user_id = 'uuid-ramesh-001';

-- 2. Set selected society to is_active = true
UPDATE residents
SET is_active = true
WHERE user_id = 'uuid-ramesh-001'
  AND society_id = 2;
```

**Result:**
```
residents table (after switch):
┌───────────┬──────────────────┬────────────┬─────────────┬────────────┬──────────┐
│resident_id│ user_id          │ society_id │ flat_number │ is_active  │ verified │
├───────────┼──────────────────┼────────────┼─────────────┼────────────┼──────────┤
│ 1         │ uuid-ramesh-001  │ 1          │ A-404       │ false      │ VERIFIED │
│ 2         │ uuid-ramesh-001  │ 2          │ 201         │ true ✓     │ VERIFIED │
└───────────┴──────────────────┴────────────┴─────────────┴────────────┴──────────┘
```

---

### 16.2 Vendor Registration & Profile Setup Flow

#### 16.2.1 Initial Vendor Registration

**Scenario:** Priya Sharma registers "Perfect Press" laundry business

**API Call:** `POST /api/v1/onboarding/vendor/register`

```sql
-- Step 1: Create user account
INSERT INTO users (user_id, phone, full_name, email, user_type, is_verified)
VALUES (
  'uuid-priya-001',
  '+919876543211',
  'Priya Sharma',
  'priya@perfectpress.com',
  'VENDOR',
  false  -- Pending platform approval
);

-- Step 2: Create vendor business profile
INSERT INTO vendors (
  vendor_id,          -- Same as user_id (1:1 relationship)
  business_name,
  store_address,
  id_proof_type,
  id_proof_number,
  gst_number,
  pan_number,
  approval_status
)
VALUES (
  'uuid-priya-001',   -- References users.user_id
  'Perfect Press',
  '789 Market Street, Koramangala',
  'AADHAAR',
  '1234-5678-9012',
  '29ABCDE1234F1Z5',
  'ABCDE1234F',
  'PENDING'
);
```

**Storage Pattern:**
```
users table:
┌──────────────────┬────────────────┬──────────────┬─────────────────────────┬───────────┐
│ user_id          │ phone          │ full_name    │ email                   │ user_type │
├──────────────────┼────────────────┼──────────────┼─────────────────────────┼───────────┤
│ uuid-priya-001   │ +919876543211  │ Priya Sharma │ priya@perfectpress.com  │ VENDOR    │
└──────────────────┴────────────────┴──────────────┴─────────────────────────┴───────────┘

vendors table:
┌──────────────────┬──────────────┬──────────────────┬─────────────────┬──────────────┐
│ vendor_id (PK,FK)│ business_name│ store_address    │ gst_number      │ approval_    │
│                  │              │                  │                 │ status       │
├──────────────────┼──────────────┼──────────────────┼─────────────────┼──────────────┤
│ uuid-priya-001   │ Perfect Press│ 789 Market St... │ 29ABCDE1234F1Z5 │ PENDING      │
└──────────────────┴──────────────┴──────────────────┴─────────────────┴──────────────┘
```

**Key Relationship:** `vendors.vendor_id` is both PRIMARY KEY and FOREIGN KEY to `users.user_id`

---

#### 16.2.2 Adding Bank Details

**API Call:** `PUT /api/v1/onboarding/vendor/{vendor_id}/bank-details`

```sql
UPDATE vendors
SET
  bank_account_number = '1234567890123',
  bank_ifsc_code = 'SBIN0001234',
  bank_account_holder = 'Priya Sharma',
  bank_name = 'State Bank of India',
  branch_name = 'Koramangala Branch',
  updated_at = NOW()
WHERE vendor_id = 'uuid-priya-001';
```

---

#### 16.2.3 Selecting Services Offered

**API Call:** `POST /api/v1/onboarding/vendor/{vendor_id}/services`

**Scenario:** Perfect Press offers Ironing, Washing+Ironing, and Dry Cleaning

```sql
-- Get service IDs first
SELECT service_id, service_key, service_name
FROM service_categories
WHERE service_key IN ('IRONING', 'WASHING_IRONING', 'DRY_CLEANING');

-- Result:
-- service_id: 1, service_key: 'IRONING'
-- service_id: 2, service_key: 'WASHING_IRONING'
-- service_id: 3, service_key: 'DRY_CLEANING'

-- Create vendor_services entries
INSERT INTO vendor_services (vendor_id, service_id, turnaround_hours, is_active)
VALUES
  ('uuid-priya-001', 1, 24, true),   -- Ironing: 24 hours
  ('uuid-priya-001', 2, 48, true),   -- Washing+Ironing: 48 hours
  ('uuid-priya-001', 3, 120, true);  -- Dry Cleaning: 120 hours
```

**Storage Pattern:**
```
vendor_services table:
┌────┬──────────────────┬────────────┬───────────────────┬──────────────────┬───────────┐
│ id │ vendor_id        │ service_id │ turnaround_hours  │ service_name     │ is_active │
├────┼──────────────────┼────────────┼───────────────────┼──────────────────┼───────────┤
│ 1  │ uuid-priya-001   │ 1          │ 24                │ Ironing          │ true      │
│ 2  │ uuid-priya-001   │ 2          │ 48                │ Wash+Iron        │ true      │
│ 3  │ uuid-priya-001   │ 3          │ 120               │ Dry Cleaning     │ true      │
└────┴──────────────────┴────────────┴───────────────────┴──────────────────┴───────────┘
```

**This creates the linkage:** Vendor → Services Offered

---

#### 16.2.4 Requesting Society Access

**API Call:** `POST /api/v1/onboarding/vendor/{vendor_id}/societies`

**Scenario:** Perfect Press requests to serve Maple Gardens and Palm Residency

```sql
-- Request access to multiple societies
INSERT INTO vendor_societies (vendor_id, society_id, approval_status)
VALUES
  ('uuid-priya-001', 1, 'PENDING'),  -- Maple Gardens
  ('uuid-priya-001', 2, 'PENDING');  -- Palm Residency
```

**Storage Pattern:**
```
vendor_societies table:
┌────┬──────────────────┬────────────┬──────────────────┬────────────────┬───────────────┐
│ id │ vendor_id        │ society_id │ society_name     │ approval_status│ approved_at   │
├────┼──────────────────┼────────────┼──────────────────┼────────────────┼───────────────┤
│ 1  │ uuid-priya-001   │ 1          │ Maple Gardens    │ PENDING        │ NULL          │
│ 2  │ uuid-priya-001   │ 2          │ Palm Residency   │ PENDING        │ NULL          │
└────┴──────────────────┴────────────┴──────────────────┴────────────────┴───────────────┘
```

**This creates the linkage:** Vendor → Societies They Want to Serve

---

### 16.3 Vendor-Society-Service Complete Linkage

#### 16.3.1 Three-Way Relationship Overview

After Perfect Press completes onboarding, here's the complete data structure:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          PERFECT PRESS (uuid-priya-001)                  │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                ┌───────────────────┼───────────────────┐
                │                   │                   │
                ▼                   ▼                   ▼

   [1] SERVICES OFFERED    [2] SOCIETIES SERVED    [3] RATE CARDS
   (vendor_services)       (vendor_societies)      (rate_cards + items)

┌─────────────────────┐  ┌──────────────────────┐  ┌────────────────────────┐
│ • Ironing (24h)     │  │ • Maple Gardens      │  │ Maple Gardens Pricing: │
│ • Wash+Iron (48h)   │  │   (APPROVED)         │  │ ├─ Shirt (Iron): ₹15   │
│ • Dry Clean (120h)  │  │ • Palm Residency     │  │ ├─ Shirt (W+I): ₹25    │
│                     │  │   (APPROVED)         │  │ ├─ Trouser (Iron): ₹20 │
│                     │  │                      │  │ └─ Saree (DC): ₹150    │
│                     │  │                      │  │                        │
│                     │  │                      │  │ Palm Residency Pricing:│
│                     │  │                      │  │ ├─ Shirt (Iron): ₹18   │
│                     │  │                      │  │ ├─ Shirt (W+I): ₹28    │
│                     │  │                      │  │ └─ Trouser (Iron): ₹22 │
└─────────────────────┘  └──────────────────────┘  └────────────────────────┘
```

---

#### 16.3.2 Creating Rate Cards Per Society

**After Society Approval:** When Maple Gardens approves Perfect Press

**API Call:** `POST /api/v1/vendors/{vendor_id}/rate-cards`

**Step 1: Create Rate Card Container**
```sql
INSERT INTO rate_cards (vendor_id, society_id, is_active, is_published)
VALUES ('uuid-priya-001', 1, true, false)  -- Not published yet
RETURNING rate_card_id;

-- Returns: rate_card_id = 1
```

**Step 2: Add Rate Card Items for Each Service**
```sql
-- For Ironing service (service_id = 1)
INSERT INTO rate_card_items (rate_card_id, service_id, item_name, price_per_piece, display_order)
VALUES
  (1, 1, 'Shirt', 15.00, 1),
  (1, 1, 'Trouser', 20.00, 2),
  (1, 1, 'T-Shirt', 12.00, 3),
  (1, 1, 'Saree', 50.00, 4);

-- For Washing+Ironing service (service_id = 2)
INSERT INTO rate_card_items (rate_card_id, service_id, item_name, price_per_piece, display_order)
VALUES
  (1, 2, 'Shirt', 25.00, 1),
  (1, 2, 'Trouser', 30.00, 2),
  (1, 2, 'T-Shirt', 22.00, 3),
  (1, 2, 'Jeans', 40.00, 4);

-- For Dry Cleaning service (service_id = 3)
INSERT INTO rate_card_items (rate_card_id, service_id, item_name, price_per_piece, display_order)
VALUES
  (1, 3, 'Suit 2-piece', 250.00, 1),
  (1, 3, 'Saree', 150.00, 2),
  (1, 3, 'Jacket', 180.00, 3),
  (1, 3, 'Coat', 200.00, 4);
```

**Step 3: Publish Rate Card**
```sql
UPDATE rate_cards
SET is_published = true, published_at = NOW()
WHERE rate_card_id = 1;
```

**Complete Storage Pattern:**
```
rate_cards table:
┌──────────────┬──────────────────┬────────────┬──────────────────┬─────────────┬──────────────┐
│ rate_card_id │ vendor_id        │ society_id │ society_name     │ is_published│ published_at │
├──────────────┼──────────────────┼────────────┼──────────────────┼─────────────┼──────────────┤
│ 1            │ uuid-priya-001   │ 1          │ Maple Gardens    │ true        │ 2025-11-20   │
│ 2            │ uuid-priya-001   │ 2          │ Palm Residency   │ true        │ 2025-11-21   │
└──────────────┴──────────────────┴────────────┴──────────────────┴─────────────┴──────────────┘

rate_card_items table (Maple Gardens rate card):
┌─────────┬──────────────┬────────────┬─────────────┬───────────────────┬────────────────┐
│ item_id │ rate_card_id │ service_id │ item_name   │ price_per_piece   │ service_type   │
├─────────┼──────────────┼────────────┼─────────────┼───────────────────┼────────────────┤
│ 1       │ 1            │ 1          │ Shirt       │ 15.00             │ Ironing        │
│ 2       │ 1            │ 1          │ Trouser     │ 20.00             │ Ironing        │
│ 3       │ 1            │ 2          │ Shirt       │ 25.00             │ Wash+Iron      │
│ 4       │ 1            │ 2          │ Trouser     │ 30.00             │ Wash+Iron      │
│ 5       │ 1            │ 3          │ Saree       │ 150.00            │ Dry Cleaning   │
│ 6       │ 1            │ 3          │ Suit 2-piece│ 250.00            │ Dry Cleaning   │
└─────────┴──────────────┴────────────┴─────────────┴───────────────────┴────────────────┘

rate_card_items table (Palm Residency rate card):
┌─────────┬──────────────┬────────────┬─────────────┬───────────────────┬────────────────┐
│ item_id │ rate_card_id │ service_id │ item_name   │ price_per_piece   │ service_type   │
├─────────┼──────────────┼────────────┼─────────────┼───────────────────┼────────────────┤
│ 7       │ 2            │ 1          │ Shirt       │ 18.00  (higher!)  │ Ironing        │
│ 8       │ 2            │ 1          │ Trouser     │ 22.00  (higher!)  │ Ironing        │
│ 9       │ 2            │ 2          │ Shirt       │ 28.00  (higher!)  │ Wash+Iron      │
└─────────┴──────────────┴────────────┴─────────────┴───────────────────┴────────────────┘
```

**Key Insight:** Same vendor can have different pricing for different societies!

---

### 16.4 Vendor Listing & Discovery Flow

#### 16.4.1 How Residents Discover Vendors

**Scenario:** Ramesh (at Maple Gardens) searches for laundry vendors

**API Call:** `GET /api/v1/vendors/search?society_id=1&service_type=LAUNDRY`

**Query Execution:**
```sql
-- Complex query that joins multiple tables
SELECT DISTINCT
  v.vendor_id,
  v.business_name,
  v.store_address,
  v.avg_rating,
  v.completed_orders,
  vs_mapping.approval_status as society_approval,
  rc.rate_card_id,
  rc.is_published
FROM vendors v
-- Vendor must be approved to serve this society
INNER JOIN vendor_societies vs_mapping
  ON vs_mapping.vendor_id = v.vendor_id
  AND vs_mapping.society_id = 1  -- Maple Gardens
  AND vs_mapping.approval_status = 'APPROVED'
-- Vendor must offer laundry services
INNER JOIN vendor_services vs
  ON vs.vendor_id = v.vendor_id
  AND vs.is_active = true
INNER JOIN service_categories sc
  ON sc.service_id = vs.service_id
INNER JOIN parent_categories pc
  ON pc.category_id = sc.parent_category_id
  AND pc.category_key = 'LAUNDRY'
-- Vendor must have published rate card for this society
LEFT JOIN rate_cards rc
  ON rc.vendor_id = v.vendor_id
  AND rc.society_id = 1
  AND rc.is_published = true
WHERE v.approval_status = 'APPROVED'
  AND v.is_available = true
ORDER BY v.avg_rating DESC, v.completed_orders DESC;
```

**What This Query Checks:**
1. ✅ Vendor approved by platform (`v.approval_status = 'APPROVED'`)
2. ✅ Vendor approved by Maple Gardens (`vs_mapping.approval_status = 'APPROVED'`)
3. ✅ Vendor offers laundry services (`vendor_services` + `service_categories`)
4. ✅ Vendor has published rate card for Maple Gardens (`rate_cards.is_published = true`)
5. ✅ Vendor is currently available (`v.is_available = true`)

**Result:**
```json
{
  "vendors": [
    {
      "vendor_id": "uuid-priya-001",
      "business_name": "Perfect Press",
      "store_address": "789 Market Street, Koramangala",
      "avg_rating": 4.8,
      "completed_orders": 1247,
      "services_offered": ["Ironing", "Washing+Ironing", "Dry Cleaning"],
      "has_rate_card": true,
      "society_approval": "APPROVED"
    }
  ]
}
```

---

#### 16.4.2 Viewing Vendor Rate Card

**API Call:** `GET /api/v1/vendors/{vendor_id}/rate-card?society_id=1`

**Query Execution:**
```sql
-- Get rate card items grouped by service
SELECT
  rc.rate_card_id,
  sc.service_id,
  sc.service_name,
  sc.service_key,
  pc.category_name,
  vs.turnaround_hours,
  json_agg(
    json_build_object(
      'item_id', rci.item_id,
      'item_name', rci.item_name,
      'description', rci.description,
      'price', rci.price_per_piece,
      'display_order', rci.display_order
    ) ORDER BY rci.display_order
  ) as items
FROM rate_cards rc
INNER JOIN rate_card_items rci
  ON rci.rate_card_id = rc.rate_card_id
  AND rci.is_active = true
INNER JOIN service_categories sc
  ON sc.service_id = rci.service_id
INNER JOIN parent_categories pc
  ON pc.category_id = sc.parent_category_id
INNER JOIN vendor_services vs
  ON vs.vendor_id = rc.vendor_id
  AND vs.service_id = sc.service_id
WHERE rc.vendor_id = 'uuid-priya-001'
  AND rc.society_id = 1
  AND rc.is_published = true
GROUP BY rc.rate_card_id, sc.service_id, sc.service_name,
         sc.service_key, pc.category_name, vs.turnaround_hours
ORDER BY sc.display_order;
```

**Result:**
```json
{
  "rate_card_id": 1,
  "vendor_name": "Perfect Press",
  "society_name": "Maple Gardens",
  "services": [
    {
      "service_id": 1,
      "service_name": "Ironing Only",
      "category": "Laundry Services",
      "turnaround_hours": 24,
      "items": [
        {"item_name": "Shirt", "price": 15.00},
        {"item_name": "Trouser", "price": 20.00},
        {"item_name": "T-Shirt", "price": 12.00},
        {"item_name": "Saree", "price": 50.00}
      ]
    },
    {
      "service_id": 2,
      "service_name": "Washing + Ironing",
      "category": "Laundry Services",
      "turnaround_hours": 48,
      "items": [
        {"item_name": "Shirt", "price": 25.00},
        {"item_name": "Trouser", "price": 30.00},
        {"item_name": "T-Shirt", "price": 22.00},
        {"item_name": "Jeans", "price": 40.00}
      ]
    },
    {
      "service_id": 3,
      "service_name": "Dry Cleaning",
      "category": "Laundry Services",
      "turnaround_hours": 120,
      "items": [
        {"item_name": "Suit 2-piece", "price": 250.00},
        {"item_name": "Saree", "price": 150.00},
        {"item_name": "Jacket", "price": 180.00},
        {"item_name": "Coat", "price": 200.00}
      ]
    }
  ]
}
```

---

### 16.5 Profile Update Flow

#### 16.5.1 Updating User Email (With Verification)

**Scenario:** Ramesh wants to add/update his email address

**Step 1: Request Email Update**

**API Call:** `POST /api/v1/users/{user_id}/update-email`

```sql
-- Create pending verification record
INSERT INTO email_verifications (
  verification_id,
  user_id,
  new_email,
  otp_code,
  otp_expires_at,
  is_verified
)
VALUES (
  'uuid-verify-001',
  'uuid-ramesh-001',
  'ramesh.new@example.com',
  '123456',  -- OTP sent via email
  NOW() + INTERVAL '10 minutes',
  false
);

-- NOTE: Email is NOT updated in users table yet!
```

**Step 2: User Receives OTP and Verifies**

**API Call:** `POST /api/v1/users/{user_id}/verify-email`

```sql
-- Verify OTP
SELECT * FROM email_verifications
WHERE verification_id = 'uuid-verify-001'
  AND otp_code = '123456'
  AND otp_expires_at > NOW()
  AND is_verified = false;

-- If valid, update users table
UPDATE users
SET
  email = 'ramesh.new@example.com',
  email_verified = true,
  updated_at = NOW()
WHERE user_id = 'uuid-ramesh-001';

-- Mark verification as complete
UPDATE email_verifications
SET is_verified = true, verified_at = NOW()
WHERE verification_id = 'uuid-verify-001';
```

**Storage Pattern:**
```
BEFORE verification:
users table:
┌──────────────────┬──────────────────────┬────────────────┐
│ user_id          │ email                │ email_verified │
├──────────────────┼──────────────────────┼────────────────┤
│ uuid-ramesh-001  │ ramesh@example.com   │ true           │
└──────────────────┴──────────────────────┴────────────────┘

email_verifications table (temporary):
┌──────────────────┬──────────────────┬─────────────────────────┬──────────┬─────────────┐
│ verification_id  │ user_id          │ new_email               │ otp_code │ is_verified │
├──────────────────┼──────────────────┼─────────────────────────┼──────────┼─────────────┤
│ uuid-verify-001  │ uuid-ramesh-001  │ ramesh.new@example.com  │ 123456   │ false       │
└──────────────────┴──────────────────┴─────────────────────────┴──────────┴─────────────┘

AFTER verification:
users table:
┌──────────────────┬─────────────────────────┬────────────────┐
│ user_id          │ email                   │ email_verified │
├──────────────────┼─────────────────────────┼────────────────┤
│ uuid-ramesh-001  │ ramesh.new@example.com  │ true           │  ← Updated!
└──────────────────┴─────────────────────────┴────────────────┘

email_verifications table:
┌──────────────────┬──────────────────┬─────────────────────────┬──────────┬─────────────┐
│ verification_id  │ user_id          │ new_email               │ otp_code │ is_verified │
├──────────────────┼──────────────────┼─────────────────────────┼──────────┼─────────────┤
│ uuid-verify-001  │ uuid-ramesh-001  │ ramesh.new@example.com  │ 123456   │ true ✓      │
└──────────────────┴──────────────────┴─────────────────────────┴──────────┴─────────────┘
```

---

#### 16.5.2 Updating Phone Number (With OTP)

**Same pattern as email, but uses phone_verifications table and SMS OTP**

**API Calls:**
1. `POST /api/v1/users/{user_id}/update-phone` → Sends OTP to new number
2. `POST /api/v1/users/{user_id}/verify-phone` → Verifies OTP and updates

---

### 16.6 Complete Data Relationship Summary

#### All Tables Involved in Vendor Listing

```
┌──────────────────────────────────────────────────────────────────────┐
│                     VENDOR LISTING QUERY FLOW                         │
└──────────────────────────────────────────────────────────────────────┘

1. users table (vendor account)
   │
   ▼
2. vendors table (business profile)
   │
   ├──> vendor_services ──> service_categories ──> parent_categories
   │    (What services?)     (Service details)     (Laundry, Vehicle, etc.)
   │
   ├──> vendor_societies ──> societies
   │    (Which societies?)    (Society details)
   │
   └──> rate_cards ──> rate_card_items
        (Pricing per society) (Items + prices per service)
```

#### Query Checks Table Chain

```
For: "Show me laundry vendors in my society"

residents table → Get user's active society
    ↓
vendor_societies → Vendors approved for this society
    ↓
vendors → Check vendor is approved & available
    ↓
vendor_services → Check vendor offers requested service
    ↓
service_categories → Verify service is laundry
    ↓
rate_cards → Check vendor has published pricing
    ↓
RESULT: List of vendors meeting ALL criteria
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
