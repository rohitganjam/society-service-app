# Society Admin Dashboard

Next.js web dashboard for society administrators to manage vendors, residents, and orders.

## Features

- Vendor approval and management
- Resident roster management (CSV upload)
- Order monitoring and tracking
- Dispute resolution
- Subscription billing management
- Analytics and reports

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

Dashboard will be available at http://localhost:3001

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
