import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

import { AuthService } from '../../core/services/auth.service';
import { DashboardUiStore } from './dashboard-ui.store';
import { InvoicesStore } from './invoices/invoices.store';
import { StockStore } from './stock/stock.store';

@Component({
  selector: 'app-dashboard-page',
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
  providers: [DashboardUiStore, StockStore, InvoicesStore],
})
export class DashboardPageComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  readonly user = this.authService.user;
  readonly ui = inject(DashboardUiStore);
  readonly stockStore = inject(StockStore);
  readonly invoicesStore = inject(InvoicesStore);

  readonly sidebarOpen = signal(false);

  constructor() {
    void this.refreshAll();
  }

  async refreshAll(): Promise<void> {
    this.ui.clearError();
    await Promise.all([this.stockStore.refresh(), this.invoicesStore.refresh()]);
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
