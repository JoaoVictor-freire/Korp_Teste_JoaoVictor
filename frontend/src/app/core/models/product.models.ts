export interface Product {
  owner_id?: string;
  code: string;
  description: string;
  stock: number;
  created_at?: string;
}

export interface CreateProductRequest {
  code: string;
  description: string;
  stock: number;
}
