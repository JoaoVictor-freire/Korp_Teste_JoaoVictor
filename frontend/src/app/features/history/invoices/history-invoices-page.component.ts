import { CommonModule } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';

import { Invoice } from '../../../core/models/invoice.models';
import { InvoicesStore } from '../../dashboard/invoices/invoices.store';
import { StockStore } from '../../dashboard/stock/stock.store';

type HistoryFilter = 'ALL' | 'OPEN' | 'CLOSED';
const PAGE_SIZE = 20;

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
  readonly currentPage = signal(1);
  readonly selectedInvoiceNumber = signal<number | null>(null);
  readonly modalOpen = signal(false);

  readonly filteredInvoices = computed(() => {
    const currentFilter = this.filter();
    if (currentFilter === 'ALL') {
      return this.invoicesStore.invoices();
    }

    return this.invoicesStore.invoices().filter((invoice) => invoice.status === currentFilter);
  });

  readonly totalPages = computed(() => {
    const total = Math.ceil(this.filteredInvoices().length / PAGE_SIZE);
    return Math.max(1, total);
  });

  readonly paginatedInvoices = computed(() => {
    const safePage = Math.min(this.currentPage(), this.totalPages());
    const start = (safePage - 1) * PAGE_SIZE;
    return this.filteredInvoices().slice(start, start + PAGE_SIZE);
  });

  readonly paginationItems = computed(() => {
    const totalPages = this.totalPages();
    const currentPage = Math.min(this.currentPage(), totalPages);

    if (totalPages <= 5) {
      return Array.from({ length: totalPages }, (_, index) => index + 1);
    }

    if (currentPage <= 3) {
      return [1, 2, 3, '...', totalPages];
    }

    if (currentPage >= totalPages - 2) {
      return [1, '...', totalPages - 2, totalPages - 1, totalPages];
    }

    return [1, '...', currentPage, '...', totalPages];
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
    this.currentPage.set(1);
  }

  goToPage(page: number): void {
    if (page < 1 || page > this.totalPages()) {
      return;
    }

    this.currentPage.set(page);
  }

  isPageNumber(item: number | string): item is number {
    return typeof item === 'number';
  }

  goToPreviousPage(): void {
    this.goToPage(this.currentPage() - 1);
  }

  goToNextPage(): void {
    this.goToPage(this.currentPage() + 1);
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

  formatInvoiceDate(value?: string): string {
    if (!value) {
      return 'Data indisponivel';
    }

    const parsedDate = new Date(value);
    if (Number.isNaN(parsedDate.getTime())) {
      return 'Data indisponivel';
    }

    return new Intl.DateTimeFormat('pt-BR').format(parsedDate);
  }

  trackByInvoiceNumber(_: number, invoice: Invoice): number {
    return invoice.number;
  }
}
