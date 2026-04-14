import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { apiConfig } from '../config/api.config';
import { Envelope } from '../models/auth.models';
import { AIInsights } from '../models/ai.models';

@Injectable({ providedIn: 'root' })
export class AIService {
  private readonly http = inject(HttpClient);

  getInsights(): Observable<Envelope<AIInsights>> {
    return this.http.get<Envelope<AIInsights>>(`${apiConfig.stockBaseUrl}/api/v1/ai/insights`);
  }
}
