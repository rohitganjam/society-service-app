#!/bin/bash

# Seed database with initial data

echo "Seeding database..."

supabase db reset --linked

echo "✅ Database seeded successfully"
