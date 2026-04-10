import { computed, inject, Injectable, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import { CreateInvoiceRequest, Invoice } from '../../../core/models/invoice.models';
import { InvoiceService } from '../../../core/services/invoice.service';
import { DashboardUiStore } from '../dashboard-ui.store';

@Injectable()
export class InvoicesStore {
  private readonly invoiceService = inject(InvoiceService);
  private readonly ui = inject(DashboardUiStore);

  readonly invoices = signal<Invoice[]>([]);
  readonly loadingInvoices = signal(false);
  readonly savingInvoice = signal(false);

  // Creation screen: show only open invoices. History screen can override later.
  readonly invoiceStatusFilter = signal<'ALL' | 'OPEN' | 'CLOSED'>('OPEN');
  readonly selectedInvoiceNumber = signal<number | null>(null);

  readonly invoiceCount = computed(() => this.invoices().length);
  readonly openInvoiceCount = computed(() => this.invoices().filter((invoice) => invoice.status === 'OPEN').length);

  readonly filteredInvoices = computed(() => {
    const filter = this.invoiceStatusFilter();
    if (filter === 'ALL') {
      return this.invoices();
    }
    return this.invoices().filter((invoice) => invoice.status === filter);
  });

  readonly selectedInvoice = computed(() => {
    const selectedNumber = this.selectedInvoiceNumber();
    const availableInvoices = this.filteredInvoices();
    if (!availableInvoices.length) {
      return null;
    }
    if (selectedNumber === null) {
      return availableInvoices[0];
    }
    return availableInvoices.find((invoice) => invoice.number === selectedNumber) ?? availableInvoices[0];
  });

  async refresh(): Promise<void> {
    await this.loadInvoices();
  }

  async createInvoice(payload: CreateInvoiceRequest): Promise<void> {
    try {
      this.savingInvoice.set(true);
      this.ui.clearError();

      await firstValueFrom(this.invoiceService.create(payload));
      this.ui.showNotice('Nota fiscal criada com sucesso.');
      await this.loadInvoices();
    } catch (error: any) {
      this.ui.setError(error?.error?.error?.message ?? 'Falha ao criar nota.');
    } finally {
      this.savingInvoice.set(false);
    }
  }

  setInvoiceFilter(filter: 'ALL' | 'OPEN' | 'CLOSED'): void {
    this.invoiceStatusFilter.set(filter);
    this.syncSelectedInvoice();
  }

  selectInvoice(invoiceNumber: number): void {
    this.selectedInvoiceNumber.set(invoiceNumber);
  }

  private async loadInvoices(): Promise<void> {
    try {
      this.loadingInvoices.set(true);
      const response = await firstValueFrom(this.invoiceService.list());
      this.invoices.set(response.data);
      this.syncSelectedInvoice();
    } catch (error: any) {
      this.ui.setError(error?.error?.error?.message ?? 'Falha ao carregar notas.');
    } finally {
      this.loadingInvoices.set(false);
    }
  }

  private syncSelectedInvoice(): void {
    const selectedNumber = this.selectedInvoiceNumber();
    const availableInvoices = this.filteredInvoices();

    if (!availableInvoices.length) {
      this.selectedInvoiceNumber.set(null);
      return;
    }

    if (selectedNumber === null || !availableInvoices.some((invoice) => invoice.number === selectedNumber)) {
      this.selectedInvoiceNumber.set(availableInvoices[0].number);
    }
  }
}
