// Database types (will be generated from Supabase)
// Run: supabase gen types typescript --project-id <project-id>

export type Json = string | number | boolean | null | { [key: string]: Json | undefined } | Json[];

// Placeholder for Supabase generated types
export interface Database {
  public: {
    Tables: {
      users: any;
      orders: any;
      vendors: any;
      categories: any;
      // ... other tables
    };
  };
}
