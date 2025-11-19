export type PaymentMethod = 'CASH' | 'UPI' | 'CARD' | 'NET_BANKING';
export type PaymentStatus = 'PENDING' | 'COMPLETED' | 'FAILED' | 'REFUNDED';

export interface Payment {
  payment_id: number;
  order_id: string;
  amount: number;
  payment_method: PaymentMethod;
  payment_status: PaymentStatus;
  razorpay_order_id?: string;
  razorpay_payment_id?: string;
  razorpay_signature?: string;
  paid_at?: string;
  created_at: string;
  updated_at: string;
}
