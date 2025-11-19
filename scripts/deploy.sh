#!/bin/bash

# Deploy all services

echo "🚀 Deploying all services..."

# Deploy backend to Vercel
echo "Deploying backend API..."
cd backend && vercel --prod && cd ..

# Deploy society admin web
echo "Deploying society admin dashboard..."
cd apps/society-admin-web && vercel --prod && cd ../..

# Deploy platform admin web
echo "Deploying platform admin dashboard..."
cd apps/platform-admin-web && vercel --prod && cd ../..

# Deploy Supabase functions
echo "Deploying Supabase edge functions..."
supabase functions deploy --project-ref $SUPABASE_PROJECT_REF

echo "✅ All services deployed successfully"
