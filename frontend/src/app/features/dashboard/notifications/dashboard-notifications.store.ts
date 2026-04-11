import { Injectable, signal } from '@angular/core';

export interface DashboardNotification {
  id: string;
  message: string;
  createdAt: string;
}

@Injectable()
export class DashboardNotificationsStore {
  readonly notifications = signal<DashboardNotification[]>([]);

  addNotification(message: string): void {
    const notification: DashboardNotification = {
      id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
      message,
      createdAt: new Date().toISOString(),
    };

    this.notifications.update((current) => [notification, ...current].slice(0, 20));
  }

  clearAll(): void {
    this.notifications.set([]);
  }
}
