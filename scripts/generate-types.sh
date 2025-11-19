#!/bin/bash

# Generate TypeScript types from Supabase schema

echo "Generating TypeScript types from Supabase..."

supabase gen types typescript \
  --project-id $SUPABASE_PROJECT_ID \
  --schema public \
  > packages/shared-types/src/database.ts

echo "✅ Types generated successfully at packages/shared-types/src/database.ts"
