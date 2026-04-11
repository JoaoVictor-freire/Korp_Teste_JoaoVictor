import { CommonModule } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { FormArray, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { CreateInvoiceRequest, Invoice } from '../../../core/models/invoice.models';
import { InvoicesStore } from '../../dashboard/invoices/invoices.store';
import { StockStore } from '../../dashboard/stock/stock.store';

type HistoryFilter = 'ALL' | 'OPEN' | 'CLOSED';
const PAGE_SIZE = 20;

@Component({
  selector: 'app-history-invoices-page',
  imports: [CommonModule, ReactiveFormsModule, MatIconModule],
  templateUrl: './history-invoices-page.component.html',
  styleUrl: './history-invoices-page.component.scss',
})
export class HistoryInvoicesPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly invoicesStore = inject(InvoicesStore);
  readonly stockStore = inject(StockStore);

  readonly filter = signal<HistoryFilter>('ALL');
  readonly currentPage = signal(1);
  readonly selectedInvoiceNumber = signal<number | null>(null);
  readonly emittingInvoiceNumber = signal<number | null>(null);
  readonly modalOpen = signal(false);
  readonly editingInvoiceNumber = signal<number | null>(null);
  readonly deletingInvoiceNumber = signal<number | null>(null);
  readonly editForm = this.fb.nonNullable.group({
    number: [1, [Validators.required, Validators.min(1)]],
    items: this.fb.array([this.createInvoiceItemGroup()]),
  });

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

  readonly products = this.stockStore.products;

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

  async emitInvoice(invoice: Invoice): Promise<void> {
    if (invoice.status === 'CLOSED') {
      return;
    }

    this.emittingInvoiceNumber.set(invoice.number);
    try {
      await this.invoicesStore.closeInvoice(invoice.number);

      const refreshedInvoice = this.invoicesStore.invoices().find((item) => item.number === invoice.number);
      if (this.selectedInvoiceNumber() === invoice.number && refreshedInvoice?.status === 'CLOSED') {
        this.modalOpen.set(false);
      }
    } finally {
      this.emittingInvoiceNumber.set(null);
    }
  }

  openEditModal(invoice: Invoice): void {
    if (invoice.status === 'CLOSED') {
      return;
    }

    this.editingInvoiceNumber.set(invoice.number);
    this.editForm.controls.number.setValue(invoice.number);

    while (this.invoiceItems.length > 0) {
      this.invoiceItems.removeAt(0);
    }

    for (const item of invoice.items) {
      this.invoiceItems.push(
        this.fb.nonNullable.group({
          product_code: [item.product_code, [Validators.required]],
          quantity: [item.quantity, [Validators.required, Validators.min(1)]],
        }),
      );
    }
  }

  closeEditModal(): void {
    this.editingInvoiceNumber.set(null);
    this.resetEditForm();
  }

  openDeleteModal(invoice: Invoice): void {
    if (invoice.status === 'CLOSED') {
      return;
    }

    this.deletingInvoiceNumber.set(invoice.number);
  }

  closeDeleteModal(): void {
    this.deletingInvoiceNumber.set(null);
  }

  get invoiceItems(): FormArray {
    return this.editForm.controls.items;
  }

  productLabel(code: string): string {
    const product = this.stockStore.products().find((item) => item.code === code);
    if (!product) {
      return code;
    }

    return `${product.code} - ${product.description}`;
  }

  editingInvoice(): Invoice | null {
    const invoiceNumber = this.editingInvoiceNumber();
    if (invoiceNumber === null) {
      return null;
    }

    return this.invoicesStore.invoices().find((invoice) => invoice.number === invoiceNumber) ?? null;
  }

  deletingInvoice(): Invoice | null {
    const invoiceNumber = this.deletingInvoiceNumber();
    if (invoiceNumber === null) {
      return null;
    }

    return this.invoicesStore.invoices().find((invoice) => invoice.number === invoiceNumber) ?? null;
  }

  addInvoiceItem(): void {
    if (!this.canAddInvoiceItem()) {
      return;
    }

    this.invoiceItems.push(this.createInvoiceItemGroup());
  }

  removeInvoiceItem(index: number): void {
    if (this.invoiceItems.length === 1) {
      return;
    }

    this.invoiceItems.removeAt(index);
  }

  canAddInvoiceItem(): boolean {
    return this.invoiceItems.length < this.products().length;
  }

  isProductSelectedInAnotherRow(productCode: string, currentIndex: number): boolean {
    if (!productCode) {
      return false;
    }

    return this.invoiceItems.controls.some((control, index) => {
      if (index === currentIndex) {
        return false;
      }

      return control.get('product_code')?.value === productCode;
    });
  }

  handleProductSelection(rowIndex: number): void {
    const selectedCode = this.invoiceItems.at(rowIndex)?.get('product_code')?.value;
    if (!selectedCode) {
      return;
    }

    if (!this.isProductSelectedInAnotherRow(selectedCode, rowIndex)) {
      return;
    }

    this.invoiceItems.at(rowIndex)?.get('product_code')?.setValue('');
  }

  async submitEdit(): Promise<void> {
    const originalNumber = this.editingInvoiceNumber();
    if (originalNumber === null) {
      return;
    }

    if (this.editForm.invalid) {
      this.editForm.markAllAsTouched();
      return;
    }

    if (this.hasDuplicateProducts()) {
      return;
    }

    const payload: CreateInvoiceRequest = this.editForm.getRawValue();
    const success = await this.invoicesStore.updateInvoice(originalNumber, payload);
    if (success) {
      this.closeEditModal();
    }
  }

  async confirmDelete(): Promise<void> {
    const invoice = this.deletingInvoice();
    if (!invoice) {
      return;
    }

    const success = await this.invoicesStore.deleteInvoice(invoice.number);
    if (success) {
      this.closeDeleteModal();
    }
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

  private createInvoiceItemGroup() {
    return this.fb.nonNullable.group({
      product_code: ['', [Validators.required]],
      quantity: [1, [Validators.required, Validators.min(1)]],
    });
  }

  private resetEditForm(): void {
    this.editForm.reset({
      number: 1,
      items: [{ product_code: '', quantity: 1 }],
    });

    while (this.invoiceItems.length > 1) {
      this.invoiceItems.removeAt(this.invoiceItems.length - 1);
    }
  }

  private hasDuplicateProducts(): boolean {
    const selectedCodes = this.invoiceItems.controls
      .map((control) => control.get('product_code')?.value)
      .filter((value): value is string => !!value);

    return new Set(selectedCodes).size !== selectedCodes.length;
  }
}
