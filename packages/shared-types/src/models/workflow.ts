export interface ServiceWorkflowTemplate {
  template_id: number;
  service_id: number;
  template_name: string;
  created_at: string;
}

export interface WorkflowStep {
  step_id: number;
  template_id: number;
  step_name: string;
  step_order: number;
  is_customer_facing: boolean;
  estimated_duration_hours?: number;
  created_at: string;
}
