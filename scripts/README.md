# Scripts

Utility scripts for development and deployment.

## Available Scripts

### `generate-types.sh`
Generate TypeScript types from Supabase database schema.

```bash
chmod +x scripts/generate-types.sh
export SUPABASE_PROJECT_ID=your-project-id
./scripts/generate-types.sh
```

### `seed-database.sh`
Reset and seed local database with test data.

```bash
chmod +x scripts/seed-database.sh
./scripts/seed-database.sh
```

### `deploy.sh`
Deploy all services (backend, admin dashboards, edge functions).

```bash
chmod +x scripts/deploy.sh
export SUPABASE_PROJECT_REF=your-project-ref
./scripts/deploy.sh
```

## Requirements

- Supabase CLI
- Vercel CLI
- Node.js 20+
- Environment variables configured
