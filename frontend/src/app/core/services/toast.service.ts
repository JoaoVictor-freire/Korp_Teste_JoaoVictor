import { Injectable, signal } from '@angular/core';

type ToastKind = 'success' | 'error';

export interface ToastMessage {
  id: number;
  kind: ToastKind;
  text: string;
}

@Injectable({ providedIn: 'root' })
export class ToastService {
  readonly toasts = signal<ToastMessage[]>([]);

  private nextId = 1;

  showSuccess(text: string): void {
    this.push('success', text);
  }

  showError(text: string): void {
    this.push('error', text);
  }

  dismiss(id: number): void {
    this.toasts.update((current) => current.filter((toast) => toast.id !== id));
  }

  private push(kind: ToastKind, text: string): void {
    const id = this.nextId++;
    this.toasts.update((current) => [...current, { id, kind, text }]);

    setTimeout(() => {
      this.dismiss(id);
    }, 3000);
  }
}
