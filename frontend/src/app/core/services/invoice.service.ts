import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { apiConfig } from '../config/api.config';
import { Envelope } from '../models/auth.models';
import { CreateInvoiceRequest, Invoice } from '../models/invoice.models';

@Injectable({ providedIn: 'root' })
export class InvoiceService {
  private readonly http = inject(HttpClient);

  list(): Observable<Envelope<Invoice[]>> {
    return this.http.get<Envelope<Invoice[]>>(`${apiConfig.billingBaseUrl}/api/v1/invoices`);
  }

  create(payload: CreateInvoiceRequest): Observable<Envelope<Invoice>> {
    return this.http.post<Envelope<Invoice>>(`${apiConfig.billingBaseUrl}/api/v1/invoices`, payload);
  }
}
