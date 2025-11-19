// API Request/Response Types

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
  meta?: {
    timestamp: string;
    request_id: string;
  };
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface CreateOrderDTO {
  vendor_id: string;
  society_id: number;
  category_id: number;
  items: OrderItemDTO[];
  pickup_datetime: string;
  pickup_address: string;
  delivery_preference?: 'SINGLE' | 'PARTIAL';
}

export interface OrderItemDTO {
  service_id: number;
  item_name: string;
  quantity: number;
  unit_price: number;
}
