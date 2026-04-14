export interface AIRecommendation {
  name: string;
  category: string;
  reason: string;
  market_signal: string;
  stock_relation: string;
}

export interface AISource {
  title: string;
  uri: string;
}

export interface AIInsights {
  generated_at: string;
  model: string;
  overview: string;
  alerts: string[];
  actions: string[];
  billing_notes: string[];
  buy_recommendations: AIRecommendation[];
  search_queries: string[];
  sources: AISource[];
  product_count: number;
  invoice_count: number;
  open_invoice_count: number;
  low_stock_count: number;
  out_of_stock_count: number;
}
