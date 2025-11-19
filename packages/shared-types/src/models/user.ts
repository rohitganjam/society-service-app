export type UserType = 'RESIDENT' | 'VENDOR' | 'SOCIETY_ADMIN' | 'PLATFORM_ADMIN';

export interface User {
  user_id: string;
  phone: string;
  user_type: UserType;
  email?: string;
  full_name?: string;
  is_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface Resident {
  resident_id: string;
  user_id: string;
  society_id: number;
  flat_number: string;
  tower?: string;
  is_active: boolean;
  created_at: string;
}

export interface Vendor {
  vendor_id: string;
  user_id: string;
  business_name: string;
  owner_name: string;
  phone: string;
  email?: string;
  category_id: number;
  society_id: number;
  approval_status: 'PENDING' | 'APPROVED' | 'REJECTED';
  is_available: boolean;
  rating?: number;
  total_orders?: number;
  created_at: string;
}
