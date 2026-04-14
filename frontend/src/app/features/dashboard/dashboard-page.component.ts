import { Component, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import { MatIconModule } from '@angular/material/icon';
import { NavigationEnd, Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { filter, map, startWith } from 'rxjs';

import { AuthService } from '../../core/services/auth.service';
import { DashboardNotificationsComponent } from './notifications/dashboard-notifications.component';
import { DashboardNotificationsStore } from './notifications/dashboard-notifications.store';
import { DashboardInsightsStore } from './dashboard-insights.store';
import { DashboardUiStore } from './dashboard-ui.store';
import { InvoicesStore } from './invoices/invoices.store';
import { StockStore } from './stock/stock.store';

@Component({
  selector: 'app-dashboard-page',
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive, MatIconModule, DashboardNotificationsComponent],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
  providers: [DashboardUiStore, DashboardNotificationsStore, StockStore, InvoicesStore, DashboardInsightsStore],
})
export class DashboardPageComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  readonly user = this.authService.user;
  readonly ui = inject(DashboardUiStore);
  readonly stockStore = inject(StockStore);
  readonly invoicesStore = inject(InvoicesStore);
  readonly insightsStore = inject(DashboardInsightsStore);

  readonly sidebarOpen = signal(false);
  private readonly currentUrl = toSignal(
    this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      map(() => this.router.url),
      startWith(this.router.url),
    ),
    { initialValue: this.router.url },
  );
  readonly showMetrics = computed(() => !this.currentUrl().includes('/dashboard/history-'));

  constructor() {
    void this.refreshAll();
  }

  async refreshAll(): Promise<void> {
    this.ui.clearError();
    await Promise.all([this.stockStore.refresh(), this.invoicesStore.refresh(), this.insightsStore.refresh()]);
  }

  closeSidebar(): void {
    this.sidebarOpen.set(false);
  }

  toggleSidebar(): void {
    this.sidebarOpen.update((open) => !open);
  }

  logout(): void {
    this.authService.logout();
    void this.router.navigate(['/']);
  }
}
