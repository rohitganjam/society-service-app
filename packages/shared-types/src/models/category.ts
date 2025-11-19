export interface ParentCategory {
  category_id: number;
  category_key: string;
  category_name: string;
  category_icon?: string;
  is_live: boolean;
  sort_order: number;
  created_at: string;
}

export interface ServiceCategory {
  service_id: number;
  parent_category_id: number;
  service_key: string;
  service_name: string;
  service_description?: string;
  estimated_duration_hours?: number;
  created_at: string;
}
