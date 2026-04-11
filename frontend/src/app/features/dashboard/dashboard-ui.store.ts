import { DestroyRef, inject, Injectable, signal } from '@angular/core';

import { ToastService } from '../../core/services/toast.service';

@Injectable()
export class DashboardUiStore {
  private readonly destroyRef = inject(DestroyRef);
  private readonly toastService = inject(ToastService);

  readonly pageError = signal('');
  readonly pageNotice = signal('');

  private noticeTimer: ReturnType<typeof setTimeout> | null = null;

  constructor() {
    this.destroyRef.onDestroy(() => {
      if (this.noticeTimer) {
        clearTimeout(this.noticeTimer);
      }
    });
  }

  clearError(): void {
    this.pageError.set('');
  }

  setError(message: string): void {
    this.pageError.set(message);
    this.toastService.showError(message);

    if (this.noticeTimer) {
      clearTimeout(this.noticeTimer);
    }

    this.noticeTimer = setTimeout(() => {
      this.pageError.set('');
      this.noticeTimer = null;
    }, 3000);
  }

  showNotice(message: string): void {
    this.pageNotice.set(message);
    this.toastService.showSuccess(message);

    if (this.noticeTimer) {
      clearTimeout(this.noticeTimer);
    }

    this.noticeTimer = setTimeout(() => {
      this.pageNotice.set('');
      this.noticeTimer = null;
    }, 3000);
  }
}
