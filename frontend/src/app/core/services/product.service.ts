import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { apiConfig } from '../config/api.config';
import { Envelope } from '../models/auth.models';
import { CreateProductRequest, Product } from '../models/product.models';

@Injectable({ providedIn: 'root' })
export class ProductService {
  private readonly http = inject(HttpClient);

  list(): Observable<Envelope<Product[]>> {
    return this.http.get<Envelope<Product[]>>(`${apiConfig.stockBaseUrl}/api/v1/products`);
  }

  create(payload: CreateProductRequest): Observable<Envelope<Product>> {
    return this.http.post<Envelope<Product>>(`${apiConfig.stockBaseUrl}/api/v1/products`, payload);
  }

  update(originalCode: string, payload: CreateProductRequest): Observable<Envelope<Product>> {
    return this.http.put<Envelope<Product>>(`${apiConfig.stockBaseUrl}/api/v1/products/${originalCode}`, payload);
  }

  delete(code: string): Observable<Envelope<{ message: string }>> {
    return this.http.delete<Envelope<{ message: string }>>(`${apiConfig.stockBaseUrl}/api/v1/products/${code}`);
  }
}
