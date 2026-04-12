import { computed, inject, Injectable, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import { CreateInvoiceRequest, Invoice } from '../../../core/models/invoice.models';
import { InvoiceService } from '../../../core/services/invoice.service';
import { DashboardNotificationsStore } from '../notifications/dashboard-notifications.store';
import { DashboardUiStore } from '../dashboard-ui.store';
import { StockStore } from '../stock/stock.store';

@Injectable()
export class InvoicesStore {
  private readonly invoiceService = inject(InvoiceService);
  private readonly ui = inject(DashboardUiStore);
  private readonly notifications = inject(DashboardNotificationsStore);
  private readonly stockStore = inject(StockStore);

  readonly invoices = signal<Invoice[]>([]);
  readonly loadingInvoices = signal(false);
  readonly savingInvoice = signal(false);
  readonly updatingInvoice = signal(false);
  readonly deletingInvoice = signal(false);
  readonly closingInvoice = signal(false);

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
      const message = error?.error?.error?.message ?? 'Falha ao criar nota.';
      this.ui.setError(message);

      const stockNotification = this.buildStockNotification(payload.number, String(message));
      if (stockNotification) {
        this.notifications.addNotification(stockNotification);
      }
    } finally {
      this.savingInvoice.set(false);
    }
  }

  async updateInvoice(originalNumber: number, payload: CreateInvoiceRequest): Promise<boolean> {
    try {
      this.updatingInvoice.set(true);
      this.ui.clearError();

      await firstValueFrom(this.invoiceService.update(originalNumber, payload));
      this.ui.showNotice('Nota fiscal atualizada com sucesso.');
      await this.loadInvoices();
      return true;
    } catch (error: any) {
      this.ui.setError(error?.error?.error?.message ?? 'Falha ao atualizar nota.');
      return false;
    } finally {
      this.updatingInvoice.set(false);
    }
  }

  async deleteInvoice(number: number): Promise<boolean> {
    try {
      this.deletingInvoice.set(true);
      this.ui.clearError();

      await firstValueFrom(this.invoiceService.delete(number));
      this.ui.showNotice('Nota fiscal removida com sucesso.');
      await this.loadInvoices();
      return true;
    } catch (error: any) {
      this.ui.setError(error?.error?.error?.message ?? 'Falha ao remover nota.');
      return false;
    } finally {
      this.deletingInvoice.set(false);
    }
  }

  async closeInvoice(number: number): Promise<void> {
    try {
      this.closingInvoice.set(true);
      this.ui.clearError();

      await firstValueFrom(this.invoiceService.close(number));
      this.ui.showNotice('Nota fiscal fechada com sucesso.');
      await this.loadInvoices();
    } catch (error: any) {
      const message = error?.error?.error?.message ?? 'Falha ao fechar nota.';
      this.ui.setError(message);

      const stockNotification = this.buildStockNotification(number, String(message));
      if (stockNotification) {
        this.notifications.addNotification(stockNotification);
      }
    } finally {
      this.closingInvoice.set(false);
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

  private buildStockNotification(invoiceNumber: number, message: string): string | null {
    const normalizedMessage = message.toLowerCase();
    if (!normalizedMessage.includes('stock')) {
      return null;
    }

    const productCodeMatch = message.match(/product\s+([a-z0-9_-]+)/i);
    const productCode = productCodeMatch?.[1];
    if (!productCode) {
      return `NF ${invoiceNumber}: ha produto em falta no estoque.`;
    }

    const product = this.stockStore.products().find((item) => item.code.toLowerCase() === productCode.toLowerCase());
    if (product) {
      return `NF ${invoiceNumber}: ${product.description} (${product.code}) esta em falta no estoque.`;
    }

    return `NF ${invoiceNumber}: produto ${productCode} esta em falta no estoque.`;
  }
}
