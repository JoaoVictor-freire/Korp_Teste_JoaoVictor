import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { FormArray, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { Invoice } from '../../../core/models/invoice.models';
import { DashboardUiStore } from '../dashboard-ui.store';
import { StockStore } from '../stock/stock.store';
import { InvoicesStore } from './invoices.store';

@Component({
  selector: 'app-invoices-page',
  imports: [CommonModule, ReactiveFormsModule, MatIconModule],
  templateUrl: './invoices-page.component.html',
  styleUrl: './invoices-page.component.scss',
})
export class InvoicesPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly ui = inject(DashboardUiStore);
  readonly stockStore = inject(StockStore);
  readonly invoicesStore = inject(InvoicesStore);

  readonly products = this.stockStore.products;
  readonly hasProducts = this.stockStore.hasProducts;

  readonly filteredInvoices = this.invoicesStore.filteredInvoices;
  readonly selectedInvoice = this.invoicesStore.selectedInvoice;
  readonly loadingInvoices = this.invoicesStore.loadingInvoices;
  readonly savingInvoice = this.invoicesStore.savingInvoice;

  readonly invoiceForm = this.fb.nonNullable.group({
    number: [1, [Validators.required, Validators.min(1)]],
    items: this.fb.array([this.createInvoiceItemGroup()]),
  });

  get invoiceItems(): FormArray {
    return this.invoiceForm.controls.items;
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

  async submitInvoice(): Promise<void> {
    if (!this.hasProducts()) {
      this.ui.setError('Cadastre pelo menos um produto antes de criar uma nota.');
      return;
    }

    if (this.invoiceForm.invalid) {
      this.invoiceForm.markAllAsTouched();
      return;
    }

    if (this.hasDuplicateProducts()) {
      this.ui.setError('Nao e permitido adicionar o mesmo produto duas vezes na mesma nota.');
      return;
    }

    const payload = this.invoiceForm.getRawValue();
    await this.invoicesStore.createInvoice(payload);

    if (this.ui.pageError()) {
      return;
    }

    this.invoiceForm.reset({
      number: this.invoicesStore.invoiceCount() + 1,
      items: [{ product_code: '', quantity: 1 }],
    });
    while (this.invoiceItems.length > 1) {
      this.invoiceItems.removeAt(this.invoiceItems.length - 1);
    }
  }

  selectInvoice(invoiceNumber: number): void {
    this.invoicesStore.selectInvoice(invoiceNumber);
  }

  trackByInvoiceNumber(_: number, invoice: Invoice): number {
    return invoice.number;
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
    this.ui.setError('Nao e permitido adicionar o mesmo produto duas vezes na mesma nota.');
  }

  canAddInvoiceItem(): boolean {
    return this.invoiceItems.length < this.products().length;
  }

  productLabel(code: string): string {
    const product = this.products().find((item) => item.code === code);
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

  private createInvoiceItemGroup() {
    return this.fb.nonNullable.group({
      product_code: ['', [Validators.required]],
      quantity: [1, [Validators.required, Validators.min(1)]],
    });
  }

  private hasDuplicateProducts(): boolean {
    const selectedCodes = this.invoiceItems.controls
      .map((control) => control.get('product_code')?.value)
      .filter((value): value is string => !!value);

    return new Set(selectedCodes).size !== selectedCodes.length;
  }
}
