# Society Service Platform

A comprehensive multi-service platform for society management, starting with laundry services and expandable to vehicle, home, and personal care services.

## Architecture

This is a **monorepo** containing all components of the platform:

- **Backend API** (Go Gin API)
- **Mobile Apps** (Flutter - Resident & Vendor)
- **Admin Dashboards** (Next.js - Society & Platform Admin)
- **Shared Packages** (TypeScript types)
- **Database** (Supabase PostgreSQL)

## Repository Structure

```
society-service-platform/
├── backend/                    # GO Gin API (Vercel)
├── apps/
│   ├── resident-app/          # Flutter - Resident mobile app
│   ├── vendor-app/            # Flutter - Vendor mobile app
│   ├── society-admin-web/     # Next.js - Society admin dashboard
│   └── platform-admin-web/    # Next.js - Platform admin dashboard
├── packages/
│   └── shared-types/          # Shared TypeScript types
├── supabase/
│   ├── migrations/            # Database migrations
│   ├── functions/             # Edge functions
│   └── seed/                  # Seed data
├── scripts/                   # Utility scripts
└── docs/                      # Documentation
    └── service_app_docs/      # Project documentation
```

## Quick Start

### Prerequisites

- Node.js 20 LTS
- Flutter SDK (for mobile apps)
- Supabase CLI
- Vercel CLI (for deployment)

### Installation

1. **Clone repository:**
```bash
git clone https://github.com/yourorg/society-service-platform.git
cd society-service-platform
```

2. **Install root dependencies:**
```bash
npm install
```

3. **Setup Backend:**
```bash
cd backend
npm install
cp .env.example .env
# Edit .env with your credentials
npm run dev
```

4. **Setup Society Admin Dashboard:**
```bash
cd apps/society-admin-web
npm install
cp .env.local.example .env.local
npm run dev
```

5. **Setup Platform Admin Dashboard:**
```bash
cd apps/platform-admin-web
npm install
cp .env.local.example .env.local
npm run dev
```

6. **Setup Flutter Apps:**
```bash
# Resident App
cd apps/resident-app
flutter pub get
flutter run

# Vendor App
cd apps/vendor-app
flutter pub get
flutter run
```

7. **Setup Supabase (Local):**
```bash
supabase start
supabase db reset
```

## Development URLs

- Backend API: http://localhost:3000/api/v1
- Society Admin: http://localhost:3001
- Platform Admin: http://localhost:3002
- Supabase Studio: http://localhost:54323

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Backend API | Go 1.23 + Gin |
| Mobile Apps | Flutter + Dart |
| Admin Dashboards | Next.js 14 + React + TypeScript |
| Database | PostgreSQL (Supabase) |
| Edge Functions | Deno (Supabase) |
| Hosting | Vercel (Backend + Web), App Stores (Mobile) |
| Payments | Razorpay |
| Notifications | Firebase Cloud Messaging |

## Key Features

### For Residents
- Browse service categories
- View vendor rate cards
- Create multi-service orders
- Track order workflow in real-time
- Make payments (Cash/UPI)
- Rate and review vendors

### For Vendors
- Registration and onboarding
- Manage rate cards for different societies
- Accept and process orders
- Update workflow status
- Track settlements
- View analytics

### For Society Admins
- Approve/reject vendors
- Upload resident rosters
- Monitor orders
- Resolve disputes
- View analytics

### For Platform Admins
- Onboard new societies
- Manage subscriptions
- Configure service categories
- Define workflows
- Platform-wide analytics

## Documentation

Detailed documentation is available in the `docs/` folder:

- [Architecture](docs/service_app_docs/TECH_STACK.md)
- [Database Schema](docs/service_app_docs/DATABASE_SCHEMA.md)
- [API Documentation](docs/service_app_docs/API.md)
- [Repository Structure](docs/service_app_docs/REPOSITORY_STRUCTURE.md)
- [Functionality Summary](docs/service_app_docs/APP_FUNCTIONALITY_SUMMARY.md)

## Deployment

### Backend API (Vercel)
```bash
cd backend
vercel --prod
```

### Admin Dashboards (Vercel)
```bash
cd apps/society-admin-web
vercel --prod

cd apps/platform-admin-web
vercel --prod
```

### Mobile Apps
```bash
# iOS
flutter build ios --release

# Android
flutter build appbundle --release
```

### Supabase Functions
```bash
supabase functions deploy
```

## Scripts

Available utility scripts in `scripts/`:

- `generate-types.sh` - Generate TypeScript types from Supabase
- `seed-database.sh` - Seed local database
- `deploy.sh` - Deploy all services

## Claude Code Agent System

This project includes a comprehensive Claude Code agent system for AI-assisted development. The system provides specialized agents for each layer of the stack, slash commands for common tasks, and enforces testing requirements.

### Quick Start with Claude Code

1. **Install Claude Code CLI** (if not already installed)
2. **Navigate to project root** and start Claude Code
3. **Use slash commands** or reference agents in conversation

### Available Slash Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `/api` | Create new API endpoint | `/api POST /orders/{id}/rating - Submit order rating` |
| `/component` | Scaffold component with tests | `/component flutter:resident:orders:RatingCard` |
| `/test` | Generate tests for a file | `/test backend/internal/services/order_service.go` |
| `/feature` | Full-stack feature implementation | `/feature Add vendor rating system` |
| `/migrate` | Create database migration | `/migrate create_ratings_table` |

### Specialized Agents

The system includes four specialized agents in `.claude/commands/`:

| Agent | Scope | Use For |
|-------|-------|---------|
| **backend-api** | `backend/**` | Go API development, handlers, services, repositories |
| **mobile** | `apps/*_app/**` | Flutter development, Riverpod, GoRouter, Freezed |
| **frontend-web** | `apps/*-web/**` | Next.js dashboards, React Query, shadcn/ui |
| **orchestrator** | All layers | Full-stack features, API contracts, cross-platform work |

### Usage Examples

**Creating a new API endpoint:**
```
/api POST /vendors/{id}/rating - Create vendor rating endpoint
```

**Creating a Flutter widget:**
```
/component flutter:resident:orders:OrderTrackingWidget
```

**Creating a React component:**
```
/component web:society:vendors:VendorApprovalTable
```

**Full-stack feature (uses orchestrator):**
```
Using the orchestrator approach, implement a vendor rating system that:
- Allows residents to rate vendors after order completion
- Shows ratings on vendor profiles
- Provides analytics for society admins
```

**Reference an agent directly:**
```
Following the backend-api agent guidelines, create the order
cancellation endpoint with proper error handling
```

### Testing Requirements

All agents enforce comprehensive testing:

- **Backend**: Unit tests for services, handler tests for endpoints
- **Mobile**: Widget tests for screens, provider tests for state
- **Web**: Component tests, hook tests, E2E for critical flows

Tests are automatically included when using slash commands.

### Project Context

The `CLAUDE.md` file at the project root contains:
- Architecture overview
- Coding conventions per platform
- API response/error format standards
- Testing requirements
- Common patterns and commands

### Agent Files Reference

```
.claude/
├── settings.json           # Hooks and permissions
└── commands/
    ├── backend-api.md      # Go backend agent
    ├── mobile.md           # Flutter mobile agent
    ├── frontend-web.md     # Next.js web agent
    ├── orchestrator.md     # Full-stack coordinator
    ├── api.md              # /api command
    ├── component.md        # /component command
    ├── test.md             # /test command
    ├── feature.md          # /feature command
    └── migrate.md          # /migrate command
```

### Hooks Configuration

The system includes automatic hooks for:
- **Auto-formatting** Go and Dart files after edits
- **Pre-approved commands** for common operations (make, flutter test, npm test)

---

## Contributing

1. Create a feature branch
2. Make your changes (use Claude Code agents for assistance)
3. Ensure all tests pass:
   ```bash
   # Backend
   cd backend && make test

   # Mobile
   cd apps/resident_app && flutter test
   cd apps/vendor_app && flutter test

   # Web
   cd apps/society-admin-web && npm test
   cd apps/platform-admin-web && npm test
   ```
4. Commit with descriptive message
5. Create pull request

## License

MIT

## Support

For issues and questions:
- GitHub Issues: [repository-url]/issues
- Documentation: `docs/` folder
