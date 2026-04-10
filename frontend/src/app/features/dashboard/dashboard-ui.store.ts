import { DestroyRef, inject, Injectable, signal } from '@angular/core';

@Injectable()
export class DashboardUiStore {
  private readonly destroyRef = inject(DestroyRef);

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
  }

  showNotice(message: string): void {
    this.pageNotice.set(message);

    if (this.noticeTimer) {
      clearTimeout(this.noticeTimer);
    }

    this.noticeTimer = setTimeout(() => {
      this.pageNotice.set('');
      this.noticeTimer = null;
    }, 3000);
  }
}

