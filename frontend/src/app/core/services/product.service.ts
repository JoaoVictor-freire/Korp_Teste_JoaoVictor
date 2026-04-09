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
}
