import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';

import { DashboardInsightsStore } from '../dashboard/dashboard-insights.store';
import { InvoicesStore } from '../dashboard/invoices/invoices.store';
import { StockStore } from '../dashboard/stock/stock.store';

@Component({
  selector: 'app-insights-page',
  imports: [CommonModule],
  templateUrl: './insights-page.component.html',
  styleUrl: './insights-page.component.scss',
  providers: [DashboardInsightsStore],
})
export class InsightsPageComponent {
  readonly insightsStore = inject(DashboardInsightsStore);
  readonly stockStore = inject(StockStore);
  readonly invoicesStore = inject(InvoicesStore);

  constructor() {
    if (!this.insightsStore.initialized()) {
      void this.insightsStore.refresh();
    }
  }
}
