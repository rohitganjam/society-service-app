export interface VendorRateCard {
  rate_card_id: number;
  vendor_id: string;
  service_id: number;
  item_name: string;
  price: number;
  unit: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface VendorService {
  vendor_service_id: number;
  vendor_id: string;
  service_id: number;
  is_offered: boolean;
  created_at: string;
}
