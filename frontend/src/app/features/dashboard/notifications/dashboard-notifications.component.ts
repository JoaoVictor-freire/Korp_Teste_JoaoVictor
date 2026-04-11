import { CommonModule } from '@angular/common';
import { Component, computed, HostListener, inject, signal } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';

import { DashboardNotificationsStore } from './dashboard-notifications.store';

@Component({
  selector: 'app-dashboard-notifications',
  imports: [CommonModule, MatIconModule],
  templateUrl: './dashboard-notifications.component.html',
  styleUrl: './dashboard-notifications.component.scss',
})
export class DashboardNotificationsComponent {
  readonly store = inject(DashboardNotificationsStore);
  readonly open = signal(false);
  readonly count = computed(() => this.store.notifications().length);

  toggle(): void {
    this.open.update((value) => !value);
  }

  close(): void {
    this.open.set(false);
  }

  clearAll(): void {
    this.store.clearAll();
    this.close();
  }

  formatTimestamp(value: string): string {
    const parsedDate = new Date(value);
    if (Number.isNaN(parsedDate.getTime())) {
      return '';
    }

    return new Intl.DateTimeFormat('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(parsedDate);
  }

  @HostListener('document:keydown.escape')
  onEscape(): void {
    this.close();
  }
}
