import { CommonModule } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';

import { Invoice } from '../../../core/models/invoice.models';
import { InvoicesStore } from '../../dashboard/invoices/invoices.store';
import { StockStore } from '../../dashboard/stock/stock.store';

type HistoryFilter = 'ALL' | 'OPEN' | 'CLOSED';

@Component({
  selector: 'app-history-invoices-page',
  imports: [CommonModule],
  templateUrl: './history-invoices-page.component.html',
  styleUrl: './history-invoices-page.component.scss',
})
export class HistoryInvoicesPageComponent {
  readonly invoicesStore = inject(InvoicesStore);
  readonly stockStore = inject(StockStore);

  readonly filter = signal<HistoryFilter>('ALL');
  readonly selectedInvoiceNumber = signal<number | null>(null);
  readonly modalOpen = signal(false);

  readonly filteredInvoices = computed(() => {
    const currentFilter = this.filter();
    if (currentFilter === 'ALL') {
      return this.invoicesStore.invoices();
    }

    return this.invoicesStore.invoices().filter((invoice) => invoice.status === currentFilter);
  });

  readonly selectedInvoice = computed(() => {
    const invoiceNumber = this.selectedInvoiceNumber();
    if (invoiceNumber === null) {
      return null;
    }

    return this.invoicesStore.invoices().find((invoice) => invoice.number === invoiceNumber) ?? null;
  });

  setFilter(filter: HistoryFilter): void {
    this.filter.set(filter);
  }

  openInvoice(invoiceNumber: number): void {
    this.selectedInvoiceNumber.set(invoiceNumber);
    this.modalOpen.set(true);
  }

  closeModal(): void {
    this.modalOpen.set(false);
  }

  async closeInvoice(): Promise<void> {
    const invoice = this.selectedInvoice();
    if (!invoice || invoice.status === 'CLOSED') {
      return;
    }

    await this.invoicesStore.closeInvoice(invoice.number);
    const refreshedInvoice = this.invoicesStore.invoices().find((item) => item.number === invoice.number);
    if (refreshedInvoice?.status === 'CLOSED') {
      this.modalOpen.set(false);
    }
  }

  productLabel(code: string): string {
    const product = this.stockStore.products().find((item) => item.code === code);
    if (!product) {
      return code;
    }

    return `${product.code} - ${product.description}`;
  }

  trackByInvoiceNumber(_: number, invoice: Invoice): number {
    return invoice.number;
  }
}
