export interface InvoiceItem {
  product_code: string;
  quantity: number;
}

export interface Invoice {
  owner_id?: string;
  number: number;
  status: string;
  items: InvoiceItem[];
  created_at?: string;
}

export interface CreateInvoiceRequest {
  number: number;
  items: InvoiceItem[];
}
