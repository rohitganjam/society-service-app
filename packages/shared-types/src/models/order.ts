export type OrderStatus =
  | 'BOOKING_CREATED'
  | 'PICKUP_SCHEDULED'
  | 'PICKUP_IN_PROGRESS'
  | 'COUNT_APPROVAL_PENDING'
  | 'PICKED_UP'
  | 'PROCESSING_IN_PROGRESS'
  | 'READY_FOR_DELIVERY'
  | 'OUT_FOR_DELIVERY'
  | 'DELIVERED'
  | 'COMPLETED'
  | 'CANCELLED';

export interface Order {
  order_id: string;
  order_number: string;
  resident_id: string;
  vendor_id: string;
  society_id: number;
  status: OrderStatus;
  has_multiple_services: boolean;
  estimated_price: number;
  final_price?: number;
  pickup_datetime: string;
  expected_delivery_date: string;
  actual_delivery_date?: string;
  pickup_address: string;
  delivery_preference: 'SINGLE' | 'PARTIAL';
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  item_id: number;
  order_id: string;
  service_id: number;
  item_name: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  created_at: string;
}

export interface OrderServiceStatus {
  status_id: number;
  order_id: string;
  service_id: number;
  item_count: number;
  total_amount: number;
  status: OrderStatus;
  expected_delivery_date: string;
  actual_delivery_date?: string;
  created_at: string;
  updated_at: string;
}
