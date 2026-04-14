import { computed, inject, Injectable, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import { AIInsights } from '../../core/models/ai.models';
import { AIService } from '../../core/services/ai.service';

@Injectable()
export class DashboardInsightsStore {
  private readonly aiService = inject(AIService);

  readonly insights = signal<AIInsights | null>(null);
  readonly loading = signal(false);
  readonly error = signal('');
  readonly initialized = signal(false);

  readonly hasInsights = computed(() => this.insights() !== null);
  readonly formattedGeneratedAt = computed(() => {
    const value = this.insights()?.generated_at;
    if (!value) {
      return '';
    }

    return new Intl.DateTimeFormat('pt-BR', {
      dateStyle: 'short',
      timeStyle: 'short',
    }).format(new Date(value));
  });

  async refresh(): Promise<void> {
    try {
      this.loading.set(true);
      this.error.set('');

      const response = await firstValueFrom(this.aiService.getInsights());
      this.insights.set(response.data);
    } catch (error: any) {
      this.insights.set(null);
      this.error.set(error?.error?.error?.message ?? 'Falha ao gerar insights com IA.');
    } finally {
      this.loading.set(false);
      this.initialized.set(true);
    }
  }
}
