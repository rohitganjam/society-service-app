# Backend API

Node.js + Express backend API for the Society Service Platform.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Copy `.env.example` to `.env` and fill in your credentials:
```bash
cp .env.example .env
```

3. Run development server:
```bash
npm run dev
```

## API Routes

Base URL: `/api/v1`

### Authentication
- `POST /auth/login` - Send OTP
- `POST /auth/verify-otp` - Verify OTP and get JWT
- `POST /auth/refresh` - Refresh JWT token

### Categories
- `GET /categories` - List all categories
- `GET /categories/:id/services` - List services in category

### Orders
- `POST /orders` - Create new order
- `GET /orders/:id` - Get order details
- `GET /orders` - List orders
- `PATCH /orders/:id/status` - Update order status

### Vendors
- `POST /vendors/register` - Register new vendor
- `GET /vendors/:id` - Get vendor details
- `GET /vendors/:id/rate-card` - Get vendor rate card

### Payments
- `POST /payments` - Create payment
- `POST /payments/verify` - Verify payment

## Testing

```bash
npm test
```

## Deployment

Deploy to Vercel:
```bash
vercel --prod
```
