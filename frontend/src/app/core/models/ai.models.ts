export interface AIInsights {
  generated_at: string;
  model: string;
  content: string;
  product_count: number;
  invoice_count: number;
  open_invoice_count: number;
  low_stock_count: number;
  out_of_stock_count: number;
}
