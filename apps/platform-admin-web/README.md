# Platform Admin Dashboard

Next.js web dashboard for platform administrators to manage the entire multi-service platform.

## Features

- Society onboarding and management
- Subscription and billing management
- Service category and workflow configuration
- Platform-wide vendor and order monitoring
- Revenue and growth analytics
- Dispute escalation management

## Setup

1. Install dependencies:
```bash
npm install
```

2. Copy `.env.local.example` to `.env.local` and configure:
```bash
cp .env.local.example .env.local
```

3. Run development server:
```bash
npm run dev
```

Dashboard will be available at http://localhost:3002

## Build

```bash
npm run build
npm start
```

## Deploy

Deploy to Vercel:
```bash
vercel --prod
```

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **UI**: Tailwind CSS + shadcn/ui
- **State Management**: React Query
- **Forms**: React Hook Form + Zod
- **Charts**: Recharts
- **API Client**: Axios
