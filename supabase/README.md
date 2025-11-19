# Supabase Configuration

This directory contains Supabase-related files including database migrations, edge functions, and seed data.

## Structure

```
supabase/
├── config.toml          # Supabase project configuration
├── migrations/          # Database schema migrations
├── functions/           # Edge functions (Deno)
│   ├── send-notification/
│   ├── razorpay-webhook/
│   ├── generate-invoices/
│   └── send-sms/
└── seed/               # Database seed files
```

## Setup

1. Install Supabase CLI:
```bash
brew install supabase/tap/supabase
# or
npm install -g supabase
```

2. Login to Supabase:
```bash
supabase login
```

3. Link to your project:
```bash
supabase link --project-ref <your-project-id>
```

4. Start local Supabase:
```bash
supabase start
```

## Migrations

Create a new migration:
```bash
supabase migration new <migration_name>
```

Apply migrations:
```bash
supabase db push
```

Reset database (local):
```bash
supabase db reset
```

## Edge Functions

Deploy all functions:
```bash
supabase functions deploy
```

Deploy specific function:
```bash
supabase functions deploy send-notification
```

Test function locally:
```bash
supabase functions serve send-notification
```

## Generate Types

Generate TypeScript types from database schema:
```bash
supabase gen types typescript --project-id <project-id> > ../packages/shared-types/src/database.ts
```

## Local Development

Access local Supabase services:
- Studio: http://localhost:54323
- API: http://localhost:54321
- Database: postgresql://postgres:postgres@localhost:54322/postgres
