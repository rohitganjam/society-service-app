# Repository Structure Guide

**Version:** 2.0
**Last Updated:** November 18, 2025
**Architecture:** Flutter Mobile + Node.js Backend + Next.js Admin
**Purpose:** Complete monorepo structure for multi-category service platform

---

## Table of Contents

1. [Repository Strategy](#1-repository-strategy)
2. [Architecture Overview](#2-architecture-overview)
3. [Complete Folder Structure](#3-complete-folder-structure)
4. [Backend API (Node.js)](#4-backend-api-nodejs)
5. [Mobile Apps (Flutter)](#5-mobile-apps-flutter)
6. [Society Admin Dashboard (Next.js)](#6-society-admin-dashboard-nextjs)
7. [Platform Admin Dashboard (Next.js)](#7-platform-admin-dashboard-nextjs)
8. [Shared Packages](#8-shared-packages)
9. [Database & Migrations](#9-database--migrations)
10. [Development Workflow](#10-development-workflow)
11. [Deployment Strategy](#11-deployment-strategy)

---

## 1. Repository Strategy

### ✅ Recommended: Monorepo

**Why Monorepo:**
- Atomic commits across mobile, backend, and admin
- Shared TypeScript types between backend and admin
- Single CI/CD pipeline
- Easier dependency management
- Simplified versioning
- Perfect for team of 1-5 developers

**Technology Stack:**
```
├── Mobile Apps: Flutter (iOS + Android)
├── Backend API: Node.js + Express (Vercel)
├── Society Admin Dashboard: Next.js 14 (Vercel)
├── Platform Admin Dashboard: Next.js 14 (Vercel)
├── Database: PostgreSQL (Supabase)
└── Edge Functions: Deno (Supabase)
```

---

## 2. Architecture Overview

### 2.1 High-Level Structure

```
multi-service-platform/              # Root monorepo
│
├── backend/                         # Node.js API server
├── apps/
│   ├── resident-app/                # Flutter - Resident app
│   ├── vendor-app/                  # Flutter - Vendor app
│   ├── society-admin-web/           # Next.js - Society Admin dashboard
│   └── platform-admin-web/          # Next.js - Platform Admin dashboard
│
├── packages/
│   └── shared-types/                # TypeScript types
│
├── supabase/
│   ├── migrations/                  # Database schemas
│   └── functions/                   # Edge functions
│
└── docs/                            # Documentation
```

### 2.2 Data Flow

```
Flutter Apps (UI Only)
        ↓ API Calls (Dio/HTTP)
Node.js Backend (All Business Logic)
        ↓ SQL Queries
Supabase PostgreSQL (Data Storage)
        ↓ Triggers
Supabase Edge Functions (Webhooks/Cron)
```

**Key Principle:** Mobile apps NEVER talk to database directly. All operations go through Backend API.

---

## 3. Complete Folder Structure

```
multi-service-platform/
│
├── .github/
│   └── workflows/
│       ├── backend-deploy.yml        # Deploy backend to Vercel
│       ├── admin-deploy.yml          # Deploy admin to Vercel
│       ├── test-backend.yml          # Backend tests
│       └── deploy-functions.yml      # Deploy Supabase functions
│
├── backend/                          # NODE.JS BACKEND API
│   ├── api/                          # Vercel serverless functions
│   │   ├── index.ts                  # Main entry point
│   │   └── v1/
│   │       ├── auth/
│   │       │   ├── login.ts          # POST /api/v1/auth/login
│   │       │   ├── verify-otp.ts     # POST /api/v1/auth/verify-otp
│   │       │   └── refresh.ts        # POST /api/v1/auth/refresh
│   │       │
│   │       ├── categories/
│   │       │   ├── list.ts           # GET /api/v1/categories
│   │       │   └── services.ts       # GET /api/v1/categories/:id/services
│   │       │
│   │       ├── residents/
│   │       │   ├── me.ts             # GET/PATCH /api/v1/residents/me
│   │       │   └── orders.ts         # GET /api/v1/residents/me/orders
│   │       │
│   │       ├── vendors/
│   │       │   ├── register.ts       # POST /api/v1/vendors/register
│   │       │   ├── get.ts            # GET /api/v1/vendors/:id
│   │       │   ├── services.ts       # GET/POST /api/v1/vendors/:id/services
│   │       │   ├── rate-card.ts      # GET/POST /api/v1/vendors/:id/rate-card
│   │       │   ├── dashboard.ts      # GET /api/v1/vendors/me/dashboard
│   │       │   └── analytics.ts      # GET /api/v1/vendors/me/analytics
│   │       │
│   │       ├── orders/
│   │       │   ├── create.ts         # POST /api/v1/orders
│   │       │   ├── get.ts            # GET /api/v1/orders/:id
│   │       │   ├── list.ts           # GET /api/v1/orders
│   │       │   ├── update-status.ts  # PATCH /api/v1/orders/:id/status
│   │       │   ├── service-status.ts # PATCH /api/v1/orders/:id/service-status
│   │       │   ├── approve-count.ts  # POST /api/v1/orders/:id/approve-count
│   │       │   ├── workflow.ts       # GET /api/v1/orders/:id/workflow
│   │       │   └── cancel.ts         # POST /api/v1/orders/:id/cancel
│   │       │
│   │       ├── payments/
│   │       │   ├── create.ts         # POST /api/v1/payments
│   │       │   ├── verify.ts         # POST /api/v1/payments/verify
│   │       │   └── list.ts           # GET /api/v1/payments
│   │       │
│   │       ├── societies/
│   │       │   ├── list.ts           # GET /api/v1/societies
│   │       │   ├── get.ts            # GET /api/v1/societies/:id
│   │       │   ├── create.ts         # POST /api/v1/societies
│   │       │   └── update.ts         # PATCH /api/v1/societies/:id
│   │       │
│   │       ├── subscriptions/
│   │       │   ├── list.ts           # GET /api/v1/subscriptions
│   │       │   ├── get.ts            # GET /api/v1/subscriptions/:id
│   │       │   ├── invoice.ts        # POST /api/v1/subscriptions/:id/invoice
│   │       │   └── update-status.ts  # PATCH /api/v1/subscriptions/:id/status
│   │       │
│   │       └── admin/
│   │           ├── dashboard.ts      # GET /api/v1/admin/dashboard
│   │           ├── vendors.ts        # GET /api/v1/admin/vendors
│   │           ├── approve-vendor.ts # POST /api/v1/admin/vendors/:id/approve
│   │           └── disputes.ts       # GET /api/v1/admin/disputes
│   │
│   ├── src/
│   │   ├── middleware/
│   │   │   ├── auth.ts               # JWT verification
│   │   │   ├── validate.ts           # Request validation (Zod)
│   │   │   ├── error-handler.ts      # Global error handling
│   │   │   └── rate-limit.ts         # Rate limiting
│   │   │
│   │   ├── services/
│   │   │   ├── order-service.ts      # Order business logic
│   │   │   ├── pricing-service.ts    # Price calculation
│   │   │   ├── workflow-service.ts   # Service workflow management
│   │   │   ├── payment-service.ts    # Razorpay integration
│   │   │   ├── notification-service.ts # FCM notifications
│   │   │   └── analytics-service.ts  # Analytics & reporting
│   │   │
│   │   ├── repositories/
│   │   │   ├── order-repository.ts   # Database operations
│   │   │   ├── vendor-repository.ts
│   │   │   ├── user-repository.ts
│   │   │   ├── category-repository.ts
│   │   │   └── workflow-repository.ts
│   │   │
│   │   ├── models/
│   │   │   ├── order.model.ts
│   │   │   ├── vendor.model.ts
│   │   │   ├── user.model.ts
│   │   │   └── category.model.ts
│   │   │
│   │   ├── utils/
│   │   │   ├── supabase.ts           # Supabase client
│   │   │   ├── logger.ts             # Winston logger
│   │   │   ├── validators.ts         # Validation helpers
│   │   │   └── helpers.ts
│   │   │
│   │   └── types/
│   │       ├── api.types.ts          # API request/response types
│   │       ├── database.types.ts     # Generated from Supabase
│   │       └── index.ts
│   │
│   ├── tests/
│   │   ├── orders/
│   │   │   ├── create-order.test.ts
│   │   │   └── workflow.test.ts
│   │   ├── vendors/
│   │   └── payments/
│   │
│   ├── vercel.json                   # Vercel configuration
│   ├── package.json
│   ├── tsconfig.json
│   └── README.md
│
├── apps/
│   │
│   ├── resident-app/                 # FLUTTER - RESIDENT APP
│   │   ├── android/                  # Android native
│   │   ├── ios/                      # iOS native
│   │   │
│   │   ├── lib/
│   │   │   ├── core/
│   │   │   │   ├── api/
│   │   │   │   │   ├── api_client.dart        # Dio HTTP client
│   │   │   │   │   ├── interceptors.dart      # Auth, logging
│   │   │   │   │   └── endpoints.dart         # API endpoints
│   │   │   │   │
│   │   │   │   ├── config/
│   │   │   │   │   ├── env.dart               # Environment config
│   │   │   │   │   └── constants.dart
│   │   │   │   │
│   │   │   │   ├── models/
│   │   │   │   │   ├── order.dart
│   │   │   │   │   ├── vendor.dart
│   │   │   │   │   ├── category.dart
│   │   │   │   │   ├── service.dart
│   │   │   │   │   └── workflow.dart
│   │   │   │   │
│   │   │   │   ├── providers/
│   │   │   │   │   ├── auth_provider.dart
│   │   │   │   │   └── app_provider.dart
│   │   │   │   │
│   │   │   │   └── router/
│   │   │   │       └── app_router.dart        # go_router config
│   │   │   │
│   │   │   ├── features/
│   │   │   │   │
│   │   │   │   ├── auth/
│   │   │   │   │   ├── data/
│   │   │   │   │   │   └── auth_repository.dart
│   │   │   │   │   ├── domain/
│   │   │   │   │   │   └── auth_model.dart
│   │   │   │   │   └── presentation/
│   │   │   │   │       ├── screens/
│   │   │   │   │       │   ├── login_screen.dart
│   │   │   │   │       │   └── otp_screen.dart
│   │   │   │   │       └── widgets/
│   │   │   │   │
│   │   │   │   ├── home/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           └── home_screen.dart    # Browse categories
│   │   │   │   │
│   │   │   │   ├── categories/
│   │   │   │   │   ├── data/
│   │   │   │   │   │   └── categories_repository.dart
│   │   │   │   │   └── presentation/
│   │   │   │   │       ├── screens/
│   │   │   │   │       │   ├── categories_list_screen.dart
│   │   │   │   │       │   └── category_services_screen.dart
│   │   │   │   │       └── widgets/
│   │   │   │   │           ├── category_card.dart
│   │   │   │   │           └── service_card.dart
│   │   │   │   │
│   │   │   │   ├── vendors/
│   │   │   │   │   ├── data/
│   │   │   │   │   │   └── vendors_repository.dart
│   │   │   │   │   └── presentation/
│   │   │   │   │       ├── screens/
│   │   │   │   │       │   ├── vendors_list_screen.dart
│   │   │   │   │       │   ├── vendor_detail_screen.dart
│   │   │   │   │       │   └── rate_card_screen.dart
│   │   │   │   │       └── widgets/
│   │   │   │   │           ├── vendor_card.dart
│   │   │   │   │           └── rate_card_item.dart
│   │   │   │   │
│   │   │   │   ├── orders/
│   │   │   │   │   ├── data/
│   │   │   │   │   │   └── orders_repository.dart  # API calls only
│   │   │   │   │   ├── domain/
│   │   │   │   │   │   └── order_model.dart
│   │   │   │   │   └── presentation/
│   │   │   │   │       ├── screens/
│   │   │   │   │       │   ├── create_order_screen.dart
│   │   │   │   │       │   ├── order_tracking_screen.dart
│   │   │   │   │       │   ├── order_history_screen.dart
│   │   │   │   │       │   └── count_approval_screen.dart
│   │   │   │   │       └── widgets/
│   │   │   │   │           ├── order_card.dart
│   │   │   │   │           ├── order_timeline.dart
│   │   │   │   │           ├── workflow_stepper.dart
│   │   │   │   │           └── service_progress_card.dart
│   │   │   │   │
│   │   │   │   ├── payments/
│   │   │   │   │   ├── data/
│   │   │   │   │   │   └── payments_repository.dart
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           ├── payment_screen.dart
│   │   │   │   │           └── payment_history_screen.dart
│   │   │   │   │
│   │   │   │   ├── disputes/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           ├── report_issue_screen.dart
│   │   │   │   │           └── dispute_details_screen.dart
│   │   │   │   │
│   │   │   │   ├── ratings/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           └── rate_service_screen.dart
│   │   │   │   │
│   │   │   │   └── profile/
│   │   │   │       └── presentation/
│   │   │   │           └── screens/
│   │   │   │               ├── profile_screen.dart
│   │   │   │               └── notifications_screen.dart
│   │   │   │
│   │   │   ├── shared/
│   │   │   │   └── widgets/
│   │   │   │       ├── app_button.dart
│   │   │   │       ├── app_text_field.dart
│   │   │   │       ├── loading_indicator.dart
│   │   │   │       └── error_widget.dart
│   │   │   │
│   │   │   └── main.dart
│   │   │
│   │   ├── assets/
│   │   │   ├── images/
│   │   │   └── icons/
│   │   │
│   │   ├── test/
│   │   ├── pubspec.yaml
│   │   └── README.md
│   │
│   ├── vendor-app/                   # FLUTTER - VENDOR APP
│   │   ├── android/
│   │   ├── ios/
│   │   │
│   │   ├── lib/
│   │   │   ├── core/                 # Same structure as resident-app
│   │   │   │
│   │   │   ├── features/
│   │   │   │   ├── auth/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           └── registration_screen.dart
│   │   │   │   │
│   │   │   │   ├── dashboard/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           └── dashboard_screen.dart  # Today's tasks
│   │   │   │   │
│   │   │   │   ├── services/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           └── services_setup_screen.dart
│   │   │   │   │
│   │   │   │   ├── orders/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           ├── orders_list_screen.dart
│   │   │   │   │           ├── order_detail_screen.dart
│   │   │   │   │           ├── pickup_screen.dart
│   │   │   │   │           ├── update_count_screen.dart
│   │   │   │   │           ├── workflow_update_screen.dart
│   │   │   │   │           └── delivery_screen.dart
│   │   │   │   │
│   │   │   │   ├── rate_card/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           ├── rate_card_setup_screen.dart
│   │   │   │   │           └── edit_rate_card_screen.dart
│   │   │   │   │
│   │   │   │   ├── settlements/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           ├── settlements_screen.dart
│   │   │   │   │           └── payment_history_screen.dart
│   │   │   │   │
│   │   │   │   ├── disputes/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           └── respond_dispute_screen.dart
│   │   │   │   │
│   │   │   │   ├── availability/
│   │   │   │   │   └── presentation/
│   │   │   │   │       └── screens/
│   │   │   │   │           └── manage_availability_screen.dart
│   │   │   │   │
│   │   │   │   └── profile/
│   │   │   │       └── presentation/
│   │   │   │           └── screens/
│   │   │   │               ├── profile_screen.dart
│   │   │   │               ├── business_details_screen.dart
│   │   │   │               └── analytics_screen.dart
│   │   │   │
│   │   │   └── main.dart
│   │   │
│   │   ├── pubspec.yaml
│   │   └── README.md
│   │
│   ├── society-admin-web/            # NEXT.JS - SOCIETY ADMIN DASHBOARD
│   │   ├── public/
│   │   ├── src/
│   │   │   ├── app/                  # Next.js App Router
│   │   │   │   ├── (auth)/
│   │   │   │   │   └── login/
│   │   │   │   │       └── page.tsx
│   │   │   │   │
│   │   │   │   ├── (dashboard)/
│   │   │   │   │   ├── layout.tsx
│   │   │   │   │   ├── page.tsx      # Society admin dashboard home
│   │   │   │   │   │
│   │   │   │   │   ├── vendors/
│   │   │   │   │   │   ├── page.tsx           # Pending vendor approvals
│   │   │   │   │   │   ├── active/
│   │   │   │   │   │   │   └── page.tsx       # Active vendors
│   │   │   │   │   │   ├── [id]/
│   │   │   │   │   │   │   └── page.tsx       # Vendor details
│   │   │   │   │   │   └── approve/
│   │   │   │   │   │       └── page.tsx       # Approve vendor
│   │   │   │   │   │
│   │   │   │   │   ├── residents/
│   │   │   │   │   │   ├── page.tsx           # Resident roster
│   │   │   │   │   │   ├── upload/
│   │   │   │   │   │   │   └── page.tsx       # Upload roster CSV
│   │   │   │   │   │   └── [id]/
│   │   │   │   │   │       └── page.tsx       # Resident details
│   │   │   │   │   │
│   │   │   │   │   ├── orders/
│   │   │   │   │   │   ├── page.tsx           # All society orders
│   │   │   │   │   │   ├── [id]/
│   │   │   │   │   │   │   └── page.tsx       # Order details
│   │   │   │   │   │   └── workflow/
│   │   │   │   │   │       └── page.tsx       # Workflow tracking
│   │   │   │   │   │
│   │   │   │   │   ├── disputes/
│   │   │   │   │   │   ├── page.tsx           # Escalated disputes
│   │   │   │   │   │   └── [id]/
│   │   │   │   │   │       └── page.tsx       # Resolve dispute
│   │   │   │   │   │
│   │   │   │   │   ├── subscription/
│   │   │   │   │   │   ├── page.tsx           # Subscription status
│   │   │   │   │   │   └── invoices/
│   │   │   │   │   │       └── page.tsx       # Billing history
│   │   │   │   │   │
│   │   │   │   │   ├── analytics/
│   │   │   │   │   │   └── page.tsx           # Society analytics
│   │   │   │   │   │
│   │   │   │   │   └── settings/
│   │   │   │   │       └── page.tsx           # Society settings
│   │   │   │   │
│   │   │   │   └── layout.tsx
│   │   │   │
│   │   │   ├── components/
│   │   │   │   ├── ui/               # shadcn components
│   │   │   │   │
│   │   │   │   ├── vendors/
│   │   │   │   │   ├── vendor-approval-card.tsx
│   │   │   │   │   ├── vendor-table.tsx
│   │   │   │   │   └── vendor-rate-card.tsx
│   │   │   │   │
│   │   │   │   ├── residents/
│   │   │   │   │   ├── resident-table.tsx
│   │   │   │   │   └── roster-upload.tsx
│   │   │   │   │
│   │   │   │   ├── orders/
│   │   │   │   │   ├── order-table.tsx
│   │   │   │   │   └── workflow-viewer.tsx
│   │   │   │   │
│   │   │   │   ├── disputes/
│   │   │   │   │   ├── dispute-table.tsx
│   │   │   │   │   └── resolution-form.tsx
│   │   │   │   │
│   │   │   │   └── analytics/
│   │   │   │       ├── completion-stats.tsx
│   │   │   │       └── vendor-performance.tsx
│   │   │   │
│   │   │   ├── lib/
│   │   │   │   ├── api-client.ts     # Axios wrapper
│   │   │   │   └── utils.ts
│   │   │   │
│   │   │   ├── hooks/
│   │   │   │   ├── use-vendors.ts
│   │   │   │   ├── use-residents.ts
│   │   │   │   ├── use-orders.ts
│   │   │   │   └── use-disputes.ts
│   │   │   │
│   │   │   └── types/
│   │   │       └── index.ts
│   │   │
│   │   ├── package.json
│   │   ├── next.config.js
│   │   ├── tailwind.config.ts
│   │   └── README.md
│   │
│   └── platform-admin-web/           # NEXT.JS - PLATFORM ADMIN DASHBOARD
│       ├── public/
│       ├── src/
│       │   ├── app/                  # Next.js App Router
│       │   │   ├── (auth)/
│       │   │   │   └── login/
│       │   │   │       └── page.tsx
│       │   │   │
│       │   │   ├── (dashboard)/
│       │   │   │   ├── layout.tsx
│       │   │   │   ├── page.tsx      # Platform overview
│       │   │   │   │
│       │   │   │   ├── societies/
│       │   │   │   │   ├── page.tsx           # All societies
│       │   │   │   │   ├── [id]/
│       │   │   │   │   │   └── page.tsx       # Society details
│       │   │   │   │   └── new/
│       │   │   │   │       └── page.tsx       # Onboard new society
│       │   │   │   │
│       │   │   │   ├── subscriptions/
│       │   │   │   │   ├── page.tsx           # All subscriptions
│       │   │   │   │   ├── invoices/
│       │   │   │   │   │   └── page.tsx       # Invoice management
│       │   │   │   │   └── overdue/
│       │   │   │   │       └── page.tsx       # Overdue payments
│       │   │   │   │
│       │   │   │   ├── categories/
│       │   │   │   │   ├── page.tsx           # Manage categories
│       │   │   │   │   ├── [id]/
│       │   │   │   │   │   ├── page.tsx       # Category details
│       │   │   │   │   │   └── workflows/
│       │   │   │   │   │       └── page.tsx   # Workflow config
│       │   │   │   │   └── activate/
│       │   │   │   │       └── page.tsx       # Activate new category
│       │   │   │   │
│       │   │   │   ├── vendors/
│       │   │   │   │   ├── page.tsx           # All platform vendors
│       │   │   │   │   └── [id]/
│       │   │   │   │       └── page.tsx       # Vendor analytics
│       │   │   │   │
│       │   │   │   ├── orders/
│       │   │   │   │   ├── page.tsx           # Platform-wide orders
│       │   │   │   │   └── workflow-analytics/
│       │   │   │   │       └── page.tsx       # Workflow bottlenecks
│       │   │   │   │
│       │   │   │   ├── disputes/
│       │   │   │   │   ├── page.tsx           # Critical escalations
│       │   │   │   │   └── [id]/
│       │   │   │   │       └── page.tsx       # Dispute resolution
│       │   │   │   │
│       │   │   │   └── analytics/
│       │   │   │       ├── page.tsx           # Platform metrics
│       │   │   │       ├── revenue/
│       │   │   │       │   └── page.tsx       # Revenue analytics
│       │   │   │       └── growth/
│       │   │   │           └── page.tsx       # Growth metrics
│       │   │   │
│       │   │   └── layout.tsx
│       │   │
│       │   ├── components/
│       │   │   ├── ui/               # shadcn components
│       │   │   │   ├── button.tsx
│       │   │   │   ├── card.tsx
│       │   │   │   └── ...
│       │   │   │
│       │   │   ├── societies/
│       │   │   │   ├── society-table.tsx
│       │   │   │   ├── society-form.tsx
│       │   │   │   └── onboarding-wizard.tsx
│       │   │   │
│       │   │   ├── subscriptions/
│       │   │   │   ├── subscription-table.tsx
│       │   │   │   └── invoice-generator.tsx
│       │   │   │
│       │   │   ├── categories/
│       │   │   │   ├── category-list.tsx
│       │   │   │   ├── workflow-editor.tsx
│       │   │   │   └── activation-form.tsx
│       │   │   │
│       │   │   ├── vendors/
│       │   │   │   └── vendor-analytics.tsx
│       │   │   │
│       │   │   ├── orders/
│       │   │   │   ├── order-table.tsx
│       │   │   │   └── workflow-tracker.tsx
│       │   │   │
│       │   │   └── analytics/
│       │   │       ├── revenue-chart.tsx
│       │   │       ├── growth-metrics.tsx
│       │   │       └── workflow-bottlenecks.tsx
│       │   │
│       │   ├── lib/
│       │   │   ├── api-client.ts     # Axios wrapper (calls backend)
│       │   │   └── utils.ts
│       │   │
│       │   ├── hooks/
│       │   │   ├── use-societies.ts
│       │   │   ├── use-subscriptions.ts
│       │   │   ├── use-vendors.ts
│       │   │   ├── use-orders.ts
│       │   │   └── use-categories.ts
│       │   │
│       │   └── types/
│       │       └── index.ts
│       │
│       ├── package.json
│       ├── next.config.js
│       ├── tailwind.config.ts
│       └── README.md
│
├── packages/
│   └── shared-types/                 # SHARED TYPESCRIPT TYPES
│       ├── src/
│       │   ├── database.ts           # Generated from Supabase
│       │   ├── api.ts                # API request/response types
│       │   ├── models/
│       │   │   ├── user.ts
│       │   │   ├── order.ts
│       │   │   ├── vendor.ts
│       │   │   ├── category.ts
│       │   │   ├── workflow.ts
│       │   │   └── payment.ts
│       │   └── index.ts
│       │
│       ├── package.json
│       └── tsconfig.json
│
├── supabase/
│   ├── config.toml
│   │
│   ├── migrations/
│   │   ├── 20250101000000_initial_schema.sql
│   │   ├── 20250101000001_add_categories.sql
│   │   ├── 20250101000002_add_workflows.sql
│   │   ├── 20250101000003_add_orders.sql
│   │   ├── 20250101000004_add_payments.sql
│   │   ├── 20250101000005_add_subscriptions.sql
│   │   ├── 20250101000006_add_rls_policies.sql
│   │   └── ...
│   │
│   ├── functions/
│   │   ├── send-notification/
│   │   │   └── index.ts
│   │   │
│   │   ├── razorpay-webhook/
│   │   │   └── index.ts
│   │   │
│   │   ├── generate-invoices/
│   │   │   └── index.ts
│   │   │
│   │   ├── send-sms/
│   │   │   └── index.ts
│   │   │
│   │   └── _shared/
│   │       └── supabase-client.ts
│   │
│   └── seed/
│       ├── 01_categories.sql
│       ├── 02_workflows.sql
│       ├── 03_societies.sql
│       └── 04_test_users.sql
│
├── docs/
│   ├── API.md
│   ├── SETUP.md
│   ├── DEPLOYMENT.md
│   ├── ARCHITECTURE.md
│   ├── DATABASE_SCHEMA.md
│   ├── TECH_STACK.md
│   └── APP_FUNCTIONALITY_SUMMARY.md
│
├── scripts/
│   ├── generate-types.sh            # Generate TS types from Supabase
│   ├── seed-database.sh
│   └── deploy.sh
│
├── .gitignore
├── .prettierrc
├── .eslintrc.json
├── package.json                     # Root package.json (workspaces)
├── pnpm-workspace.yaml
└── README.md
```

---

## 4. Backend API (Node.js)

### 4.1 Structure Principles

**Clean Architecture:**
- **API Layer** (`api/v1/`): HTTP endpoints (request/response)
- **Service Layer** (`src/services/`): Business logic
- **Repository Layer** (`src/repositories/`): Database operations
- **Middleware** (`src/middleware/`): Auth, validation, error handling

**No business logic in API routes** - all logic in services.

### 4.2 Example Endpoint

**File: `backend/api/v1/orders/create.ts`**

```typescript
import { Request, Response } from 'express';
import { z } from 'zod';
import { OrderService } from '@/services/order-service';
import { authenticate } from '@/middleware/auth';
import { validate } from '@/middleware/validate';

const createOrderSchema = z.object({
  vendor_id: z.string().uuid(),
  society_id: z.number().int(),
  category_id: z.number().int(),
  items: z.array(z.object({
    service_id: z.number().int(),
    item_name: z.string(),
    quantity: z.number().int().positive(),
    unit_price: z.number().positive()
  })).min(1),
  pickup_datetime: z.string().datetime(),
  pickup_address: z.string()
});

export default authenticate(
  validate(createOrderSchema),
  async (req: Request, res: Response) => {
    const userId = req.user.id;
    const orderData = req.body;

    // All business logic in service layer
    const result = await OrderService.createOrder(userId, orderData);

    res.status(201).json({
      success: true,
      data: result
    });
  }
);
```

### 4.3 Service Layer Example

**File: `backend/src/services/order-service.ts`**

```typescript
export class OrderService {
  static async createOrder(residentId: string, orderData: CreateOrderDTO) {
    // 1. Validate vendor
    const vendor = await VendorRepository.getById(orderData.vendor_id);
    if (!vendor.is_available) {
      throw new Error('Vendor not available');
    }

    // 2. Calculate pricing
    const pricing = await PricingService.calculateTotal(orderData.items);

    // 3. Calculate delivery estimate
    const deliveryDate = await WorkflowService.calculateDeliveryDate(
      orderData.items
    );

    // 4. Create order
    const order = await OrderRepository.create({
      resident_id: residentId,
      vendor_id: orderData.vendor_id,
      estimated_price: pricing.total,
      expected_delivery_date: deliveryDate,
      // ...
    });

    // 5. Initialize workflow tracking
    await WorkflowService.initializeTracking(order.order_id, orderData.items);

    // 6. Send notification
    await NotificationService.sendOrderNotification(order.order_id);

    return {
      order_id: order.order_id,
      total: pricing.total,
      delivery_date: deliveryDate
    };
  }
}
```

### 4.4 Key Files

| File | Purpose |
|------|---------|
| `api/v1/*/` | HTTP endpoints |
| `src/services/` | Business logic |
| `src/repositories/` | Database queries |
| `src/middleware/auth.ts` | JWT verification |
| `src/middleware/validate.ts` | Zod validation |
| `src/utils/supabase.ts` | Supabase client |
| `vercel.json` | Vercel config |

---

## 5. Mobile Apps (Flutter)

### 5.1 Architecture Pattern

**Clean Architecture + Repository Pattern**

```
Presentation Layer (UI)
        ↓
Domain Layer (Models)
        ↓
Data Layer (Repositories - API calls ONLY)
        ↓
Backend API
```

### 5.2 Key Structure

**Features-First Organization:**

Each feature has:
- `data/` - Repositories (API calls)
- `domain/` - Models
- `presentation/` - Screens & Widgets

**Example: Orders Feature**

```dart
// lib/features/orders/data/orders_repository.dart
class OrdersRepository {
  final ApiClient _apiClient;

  // Only API calls, no business logic
  Future<Order> createOrder(CreateOrderRequest request) async {
    final response = await _apiClient.post(
      '/orders',
      data: request.toJson(),
    );
    return Order.fromJson(response);
  }
}

// lib/features/orders/presentation/screens/create_order_screen.dart
class CreateOrderScreen extends ConsumerWidget {
  Future<void> _createOrder() async {
    // Call repository -> backend API
    final order = await ref.read(ordersRepositoryProvider)
        .createOrder(orderRequest);

    // Navigate to tracking
    context.go('/orders/${order.orderId}');
  }
}
```

### 5.3 Key Packages

```yaml
dependencies:
  flutter_riverpod: ^2.4.0      # State management
  go_router: ^12.0.0            # Navigation
  dio: ^5.4.0                   # HTTP client
  freezed: ^2.4.0               # Immutable models
  flutter_secure_storage: ^9.0.0 # Token storage
  razorpay_flutter: ^1.3.0      # Payments
  firebase_messaging: ^14.7.0   # Notifications
  cached_network_image: ^3.3.0  # Image caching
```

### 5.4 No Business Logic in Flutter

❌ **Don't do:**
```dart
// DON'T calculate prices in Flutter
final total = items.fold(0, (sum, item) => sum + item.price);
```

✅ **Do:**
```dart
// Let backend calculate and return total
final order = await api.createOrder(items); // Backend returns total
final total = order.total;
```

---

## 6. Society Admin Dashboard (Next.js)

### 6.1 Purpose

**Society Admin dashboard** for managing a single society's vendors, residents, orders, and disputes.

**Key Features:**
- Approve/reject vendor registrations for their society
- Upload resident rosters (CSV) for instant verification
- Monitor all orders within the society
- View workflow progress for all orders
- Resolve escalated disputes
- View society-level analytics
- Manage subscription billing status

### 6.2 App Router Structure

```
app/
├── (auth)/
│   └── login/page.tsx           # /login
│
├── (dashboard)/
│   ├── layout.tsx               # Society admin layout
│   ├── page.tsx                 # /dashboard (Overview)
│   │
│   ├── vendors/
│   │   ├── page.tsx             # /dashboard/vendors (Pending approvals)
│   │   ├── active/page.tsx      # /dashboard/vendors/active
│   │   ├── [id]/page.tsx        # /dashboard/vendors/:id
│   │   └── approve/page.tsx     # /dashboard/vendors/approve
│   │
│   ├── residents/
│   │   ├── page.tsx             # /dashboard/residents (Roster)
│   │   ├── upload/page.tsx      # /dashboard/residents/upload
│   │   └── [id]/page.tsx        # /dashboard/residents/:id
│   │
│   ├── orders/
│   │   ├── page.tsx             # /dashboard/orders (All society orders)
│   │   ├── [id]/page.tsx        # /dashboard/orders/:id
│   │   └── workflow/page.tsx    # /dashboard/orders/workflow
│   │
│   ├── disputes/
│   │   ├── page.tsx             # /dashboard/disputes
│   │   └── [id]/page.tsx        # /dashboard/disputes/:id (Resolve)
│   │
│   ├── subscription/
│   │   ├── page.tsx             # /dashboard/subscription
│   │   └── invoices/page.tsx    # /dashboard/subscription/invoices
│   │
│   ├── analytics/
│   │   └── page.tsx             # /dashboard/analytics
│   │
│   └── settings/
│       └── page.tsx             # /dashboard/settings
```

### 6.3 API Client Pattern

```typescript
// src/lib/api-client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL, // Backend API
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
```

### 6.4 Example: Vendor Approval Page

```typescript
// src/app/(dashboard)/vendors/page.tsx
'use client';

import { useQuery, useMutation } from '@tanstack/react-query';
import apiClient from '@/lib/api-client';
import { VendorApprovalCard } from '@/components/vendors/vendor-approval-card';

export default function VendorsPage() {
  const { data: pendingVendors } = useQuery({
    queryKey: ['vendors', 'pending'],
    queryFn: async () => {
      const { data } = await apiClient.get('/api/v1/society-admin/vendors/pending');
      return data;
    }
  });

  const approveMutation = useMutation({
    mutationFn: (vendorId: string) =>
      apiClient.post(`/api/v1/society-admin/vendors/${vendorId}/approve`),
    onSuccess: () => {
      queryClient.invalidateQueries(['vendors']);
    }
  });

  return (
    <div>
      <h1>Pending Vendor Approvals</h1>
      {pendingVendors?.data.map(vendor => (
        <VendorApprovalCard
          key={vendor.id}
          vendor={vendor}
          onApprove={() => approveMutation.mutate(vendor.id)}
        />
      ))}
    </div>
  );
}
```

### 6.5 Key Components

- `vendor-approval-card.tsx` - Vendor approval UI
- `vendor-table.tsx` - Active vendors list
- `resident-table.tsx` - Resident roster
- `roster-upload.tsx` - CSV upload for residents
- `order-table.tsx` - Society orders
- `workflow-viewer.tsx` - View order workflow progress
- `dispute-table.tsx` - Escalated disputes
- `resolution-form.tsx` - Resolve disputes
- `completion-stats.tsx` - Society analytics
- `vendor-performance.tsx` - Vendor ratings/stats

---

## 7. Platform Admin Dashboard (Next.js)

### 7.1 Purpose

**Platform Admin dashboard** for managing the entire multi-service platform across all societies.

**Key Features:**
- Onboard new societies to the platform
- Manage platform-wide subscriptions and billing
- Monitor overdue payments and suspend societies if needed
- **Manage service categories and workflows:**
  - Add new parent categories (Laundry, Vehicle, Home, Personal)
  - Define workflow steps per service type
  - Activate/deactivate categories for launch
- View platform-wide order metrics
- Analyze workflow bottlenecks across all orders
- Handle critical escalations
- Platform-wide revenue and growth analytics

### 7.2 App Router Structure

```
app/
├── (auth)/
│   └── login/page.tsx           # /login
│
├── (dashboard)/
│   ├── layout.tsx               # Platform admin layout
│   ├── page.tsx                 # /dashboard (Platform overview)
│   │
│   ├── societies/
│   │   ├── page.tsx             # /dashboard/societies (All societies)
│   │   ├── [id]/page.tsx        # /dashboard/societies/:id
│   │   └── new/page.tsx         # /dashboard/societies/new (Onboard)
│   │
│   ├── subscriptions/
│   │   ├── page.tsx             # /dashboard/subscriptions
│   │   ├── invoices/page.tsx    # /dashboard/subscriptions/invoices
│   │   └── overdue/page.tsx     # /dashboard/subscriptions/overdue
│   │
│   ├── categories/
│   │   ├── page.tsx             # /dashboard/categories (Manage)
│   │   ├── [id]/
│   │   │   ├── page.tsx         # /dashboard/categories/:id
│   │   │   └── workflows/
│   │   │       └── page.tsx     # /dashboard/categories/:id/workflows
│   │   └── activate/
│   │       └── page.tsx         # /dashboard/categories/activate
│   │
│   ├── vendors/
│   │   ├── page.tsx             # /dashboard/vendors (Platform-wide)
│   │   └── [id]/page.tsx        # /dashboard/vendors/:id (Analytics)
│   │
│   ├── orders/
│   │   ├── page.tsx             # /dashboard/orders (All orders)
│   │   └── workflow-analytics/
│   │       └── page.tsx         # /dashboard/orders/workflow-analytics
│   │
│   ├── disputes/
│   │   ├── page.tsx             # /dashboard/disputes
│   │   └── [id]/page.tsx        # /dashboard/disputes/:id
│   │
│   └── analytics/
│       ├── page.tsx             # /dashboard/analytics (Platform metrics)
│       ├── revenue/page.tsx     # /dashboard/analytics/revenue
│       └── growth/page.tsx      # /dashboard/analytics/growth
```

### 7.3 Example: Category Management Page

```typescript
// src/app/(dashboard)/categories/page.tsx
'use client';

import { useQuery, useMutation } from '@tanstack/react-query';
import apiClient from '@/lib/api-client';
import { CategoryList } from '@/components/categories/category-list';

export default function CategoriesPage() {
  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: async () => {
      const { data } = await apiClient.get('/api/v1/admin/categories');
      return data;
    }
  });

  const activateMutation = useMutation({
    mutationFn: (categoryId: number) =>
      apiClient.post(`/api/v1/admin/categories/${categoryId}/activate`),
    onSuccess: () => {
      queryClient.invalidateQueries(['categories']);
    }
  });

  return (
    <div>
      <h1>Manage Service Categories</h1>
      <p>Activate/deactivate categories to launch new services</p>

      {categories?.data.map(category => (
        <div key={category.id}>
          <h3>{category.name}</h3>
          <p>Status: {category.is_live ? 'ACTIVE' : 'INACTIVE'}</p>
          {!category.is_live && (
            <button onClick={() => activateMutation.mutate(category.id)}>
              Activate Category
            </button>
          )}
        </div>
      ))}
    </div>
  );
}
```

### 7.4 Key Components

- `society-table.tsx` - All societies list
- `society-form.tsx` - Onboard new society
- `onboarding-wizard.tsx` - Multi-step society onboarding
- `subscription-table.tsx` - All subscriptions
- `invoice-generator.tsx` - Generate invoices
- `category-list.tsx` - All categories
- `workflow-editor.tsx` - Define workflow steps per service
- `activation-form.tsx` - Activate new category
- `vendor-analytics.tsx` - Platform-wide vendor analytics
- `order-table.tsx` - All platform orders
- `workflow-tracker.tsx` - Workflow progress tracker
- `revenue-chart.tsx` - Revenue analytics
- `growth-metrics.tsx` - Growth metrics
- `workflow-bottlenecks.tsx` - Identify workflow bottlenecks

### 7.5 Key Dependencies

```json
{
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.2.0",
    "@tanstack/react-query": "^5.8.0",
    "axios": "^1.6.0",
    "shadcn/ui": "latest",
    "tailwindcss": "^3.3.0",
    "recharts": "^2.10.0"
  }
}
```

---

## 8. Shared Packages

### 8.1 Shared Types Package

**Purpose:** Share TypeScript types between backend and both admin dashboards

**File: `packages/shared-types/src/models/order.ts`**

```typescript
export interface Order {
  order_id: string;
  order_number: string;
  resident_id: string;
  vendor_id: string;
  society_id: number;
  status: OrderStatus;
  has_multiple_services: boolean;
  estimated_price: number;
  final_price?: number;
  pickup_datetime: string;
  expected_delivery_date: string;
  created_at: string;
}

export type OrderStatus =
  | 'BOOKING_CREATED'
  | 'PICKUP_SCHEDULED'
  | 'PICKUP_IN_PROGRESS'
  | 'COUNT_APPROVAL_PENDING'
  | 'PICKED_UP'
  | 'PROCESSING_IN_PROGRESS'
  | 'READY_FOR_DELIVERY'
  | 'OUT_FOR_DELIVERY'
  | 'DELIVERED'
  | 'COMPLETED'
  | 'CANCELLED';

export interface CreateOrderDTO {
  vendor_id: string;
  society_id: number;
  category_id: number;
  items: OrderItem[];
  pickup_datetime: string;
  pickup_address: string;
}
```

**Usage:**

```typescript
// backend/src/services/order-service.ts
import { Order, CreateOrderDTO } from '@laundry-platform/shared-types';

export class OrderService {
  static async createOrder(data: CreateOrderDTO): Promise<Order> {
    // ...
  }
}

// admin-web/src/hooks/use-orders.ts
import { Order } from '@laundry-platform/shared-types';

export const useOrders = () => {
  return useQuery<Order[]>({
    // ...
  });
};
```

---

## 9. Database & Migrations

### 9.1 Supabase Migrations

**Migration files in chronological order:**

```sql
-- supabase/migrations/20250101000000_initial_schema.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  phone VARCHAR(15) UNIQUE NOT NULL,
  user_type VARCHAR(20) NOT NULL,
  -- ...
);

-- supabase/migrations/20250101000001_add_categories.sql
CREATE TABLE parent_categories (
  category_id SERIAL PRIMARY KEY,
  category_key VARCHAR(50) UNIQUE NOT NULL,
  category_name VARCHAR(100) NOT NULL,
  is_live BOOLEAN DEFAULT false,
  -- ...
);

CREATE TABLE service_categories (
  service_id SERIAL PRIMARY KEY,
  parent_category_id INTEGER REFERENCES parent_categories(category_id),
  service_key VARCHAR(50) NOT NULL,
  service_name VARCHAR(100) NOT NULL,
  -- ...
);

-- supabase/migrations/20250101000002_add_workflows.sql
CREATE TABLE service_workflow_templates (
  template_id SERIAL PRIMARY KEY,
  service_id INTEGER REFERENCES service_categories(service_id),
  template_name VARCHAR(100) NOT NULL,
  -- ...
);

CREATE TABLE workflow_steps (
  step_id SERIAL PRIMARY KEY,
  template_id INTEGER REFERENCES service_workflow_templates(template_id),
  step_name VARCHAR(100) NOT NULL,
  step_order INTEGER NOT NULL,
  -- ...
);
```

### 9.2 Seed Data

**File: `supabase/seed/01_categories.sql`**

```sql
-- Seed parent categories
INSERT INTO parent_categories (category_key, category_name, is_live) VALUES
  ('LAUNDRY', 'Laundry Services', true),
  ('VEHICLE', 'Vehicle Services', false),
  ('HOME', 'Home Services', false),
  ('PERSONAL', 'Personal Care', false);

-- Seed service categories
DO $$
DECLARE
  laundry_id INTEGER;
BEGIN
  SELECT category_id INTO laundry_id FROM parent_categories WHERE category_key = 'LAUNDRY';

  INSERT INTO service_categories (parent_category_id, service_key, service_name) VALUES
    (laundry_id, 'IRONING', 'Ironing Only'),
    (laundry_id, 'WASHING_IRONING', 'Washing + Ironing'),
    (laundry_id, 'DRY_CLEANING', 'Dry Cleaning');
END $$;
```

### 9.3 Type Generation

**Script: `scripts/generate-types.sh`**

```bash
#!/bin/bash

# Generate TypeScript types from Supabase schema
supabase gen types typescript \
  --project-id $SUPABASE_PROJECT_ID \
  --schema public \
  > packages/shared-types/src/database.ts

echo "Types generated successfully!"
```

**Run:**
```bash
chmod +x scripts/generate-types.sh
./scripts/generate-types.sh
```

---

## 10. Development Workflow

### 10.1 Initial Setup

```bash
# 1. Clone repository
git clone https://github.com/yourorg/multi-service-platform.git
cd multi-service-platform

# 2. Install dependencies
npm install              # Root + backend
cd apps/society-admin-web && npm install
cd ../platform-admin-web && npm install
cd ../resident-app && flutter pub get
cd ../vendor-app && flutter pub get

# 3. Setup environment variables
cp backend/.env.example backend/.env
cp apps/society-admin-web/.env.local.example apps/society-admin-web/.env.local
cp apps/platform-admin-web/.env.local.example apps/platform-admin-web/.env.local

# 4. Start Supabase locally
supabase start

# 5. Run migrations
supabase db reset

# 6. Generate types
./scripts/generate-types.sh

# 7. Start development servers
# Terminal 1: Backend
cd backend && npm run dev

# Terminal 2: Society Admin
cd apps/society-admin-web && npm run dev

# Terminal 3: Platform Admin
cd apps/platform-admin-web && npm run dev

# Terminal 4: Resident app
cd apps/resident-app && flutter run

# Terminal 5: Vendor app
cd apps/vendor-app && flutter run
```

### 10.2 Development URLs

```
Backend API: http://localhost:3000/api/v1
Society Admin Web: http://localhost:3001
Platform Admin Web: http://localhost:3002
Supabase Studio: http://localhost:54323
PostgreSQL: postgresql://postgres:postgres@localhost:54322/postgres
```

### 10.3 Common Commands

```bash
# Backend
npm run dev              # Start dev server
npm test                 # Run tests
npm run lint             # Lint code

# Flutter apps
flutter run              # Run app
flutter test             # Run tests
flutter build apk        # Build Android
flutter build ios        # Build iOS

# Society Admin / Platform Admin
npm run dev              # Start dev server
npm run build            # Build for production
npm run lint             # Lint code

# Database
supabase db reset        # Reset database
supabase db push         # Push migrations
supabase functions deploy # Deploy edge functions

# Types
./scripts/generate-types.sh  # Generate types from DB
```

---

## 11. Deployment Strategy

### 11.1 Backend API (Vercel)

**File: `backend/vercel.json`**

```json
{
  "version": 2,
  "builds": [
    {
      "src": "api/**/*.ts",
      "use": "@vercel/node"
    }
  ],
  "routes": [
    {
      "src": "/api/v1/(.*)",
      "dest": "/api/$1"
    }
  ],
  "env": {
    "SUPABASE_URL": "@supabase-url",
    "SUPABASE_SERVICE_KEY": "@supabase-service-key",
    "JWT_SECRET": "@jwt-secret"
  }
}
```

**Deploy:**
```bash
cd backend
vercel --prod
```

### 11.2 Society Admin Web (Vercel)

**Deploy:**
```bash
cd apps/society-admin-web
vercel --prod
```

### 11.3 Platform Admin Web (Vercel)

**Deploy:**
```bash
cd apps/platform-admin-web
vercel --prod
```

### 11.4 Mobile Apps

**Resident App:**
```bash
cd apps/resident-app

# iOS
flutter build ios --release
# Upload to App Store Connect

# Android
flutter build appbundle --release
# Upload to Play Console
```

**Vendor App:**
```bash
cd apps/vendor-app

# iOS
flutter build ios --release
# Upload to App Store Connect

# Android
flutter build appbundle --release
# Upload to Play Console
```

### 11.5 Supabase

**Migrations:**
```bash
supabase db push --linked
```

**Edge Functions:**
```bash
supabase functions deploy send-notification
supabase functions deploy razorpay-webhook
supabase functions deploy generate-invoices
```

### 11.6 CI/CD Pipeline

**File: `.github/workflows/deploy.yml`**

```yaml
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
      - run: cd backend && npm install
      - run: cd backend && npm test
      - uses: amondnet/vercel-action@v25
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}

  deploy-society-admin:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: cd apps/society-admin-web && npm install
      - run: cd apps/society-admin-web && npm run build
      - uses: amondnet/vercel-action@v25

  deploy-platform-admin:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: cd apps/platform-admin-web && npm install
      - run: cd apps/platform-admin-web && npm run build
      - uses: amondnet/vercel-action@v25

  deploy-functions:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: supabase/setup-cli@v1
      - run: supabase functions deploy
```

---

## Summary

### Repository Overview

```
multi-service-platform/
├── backend/                 # Node.js API (Vercel)
├── apps/
│   ├── resident-app/        # Flutter (iOS + Android)
│   ├── vendor-app/          # Flutter (iOS + Android)
│   ├── society-admin-web/   # Next.js (Vercel)
│   └── platform-admin-web/  # Next.js (Vercel)
├── packages/
│   └── shared-types/        # TypeScript types
├── supabase/
│   ├── migrations/          # Database schemas
│   └── functions/           # Edge functions
└── docs/                    # Documentation
```

### Key Principles

✅ **Monorepo** - Single repository for all code
✅ **Clean Architecture** - Clear separation of concerns
✅ **API-First** - All operations through Backend API
✅ **Type Safety** - TypeScript everywhere
✅ **Serverless** - Zero server management
✅ **Scalable** - Built for multi-category expansion

### Development Flow

```
Code → Git Push → GitHub Actions → Deploy
  ├─→ Vercel (Backend API)
  ├─→ Vercel (Society Admin Web)
  ├─→ Vercel (Platform Admin Web)
  ├─→ Supabase (Edge Functions)
  └─→ App Stores (Mobile Apps)
```

**Zero downtime deployments. Automatic rollbacks on failure.**

---

**End of Document**
