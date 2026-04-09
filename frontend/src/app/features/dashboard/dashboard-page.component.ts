import { Component, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormArray, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { firstValueFrom } from 'rxjs';

import { AuthService } from '../../core/services/auth.service';
import { ProductService } from '../../core/services/product.service';
import { InvoiceService } from '../../core/services/invoice.service';
import { Product } from '../../core/models/product.models';
import { Invoice } from '../../core/models/invoice.models';

@Component({
  selector: 'app-dashboard-page',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
})
export class DashboardPageComponent {
  private readonly fb = inject(FormBuilder);
  private readonly authService = inject(AuthService);
  private readonly productService = inject(ProductService);
  private readonly invoiceService = inject(InvoiceService);
  private readonly router = inject(Router);

  readonly user = this.authService.user;
  readonly products = signal<Product[]>([]);
  readonly invoices = signal<Invoice[]>([]);
  readonly pageError = signal('');
  readonly pageNotice = signal('');
  readonly loadingProducts = signal(false);
  readonly loadingInvoices = signal(false);
  readonly savingProduct = signal(false);
  readonly savingInvoice = signal(false);

  readonly productCount = computed(() => this.products().length);
  readonly invoiceCount = computed(() => this.invoices().length);
  readonly openInvoiceCount = computed(() => this.invoices().filter((invoice) => invoice.status === 'OPEN').length);
  readonly hasProducts = computed(() => this.products().length > 0);

  readonly productForm = this.fb.nonNullable.group({
    code: ['', [Validators.required]],
    description: ['', [Validators.required]],
    stock: [0, [Validators.required, Validators.min(0)]],
  });

  readonly invoiceForm = this.fb.nonNullable.group({
    number: [1, [Validators.required, Validators.min(1)]],
    items: this.fb.array([this.createInvoiceItemGroup()]),
  });

  constructor() {
    void this.refreshAll();
  }

  get invoiceItems(): FormArray {
    return this.invoiceForm.controls.items;
  }

  addInvoiceItem(): void {
    this.invoiceItems.push(this.createInvoiceItemGroup());
  }

  removeInvoiceItem(index: number): void {
    if (this.invoiceItems.length === 1) {
      return;
    }

    this.invoiceItems.removeAt(index);
  }

  async refreshAll(): Promise<void> {
    this.pageError.set('');
    await Promise.all([this.loadProducts(), this.loadInvoices()]);
  }

  async submitProduct(): Promise<void> {
    if (this.productForm.invalid) {
      this.productForm.markAllAsTouched();
      return;
    }

    try {
      this.savingProduct.set(true);
      this.pageNotice.set('');
      await firstValueFrom(this.productService.create(this.productForm.getRawValue()));
      this.productForm.reset({ code: '', description: '', stock: 0 });
      this.pageNotice.set('Produto cadastrado com sucesso.');
      await this.loadProducts();
    } catch (error: any) {
      this.pageError.set(error?.error?.error?.message ?? 'Falha ao cadastrar produto.');
    } finally {
      this.savingProduct.set(false);
    }
  }

  async submitInvoice(): Promise<void> {
    if (!this.hasProducts()) {
      this.pageError.set('Cadastre pelo menos um produto antes de criar uma nota.');
      return;
    }

    if (this.invoiceForm.invalid) {
      this.invoiceForm.markAllAsTouched();
      return;
    }

    try {
      this.savingInvoice.set(true);
      this.pageNotice.set('');
      await firstValueFrom(this.invoiceService.create(this.invoiceForm.getRawValue()));
      this.invoiceForm.reset({
        number: this.invoiceCount() + 1,
        items: [{ product_code: '', quantity: 1 }],
      });
      while (this.invoiceItems.length > 1) {
        this.invoiceItems.removeAt(this.invoiceItems.length - 1);
      }
      this.pageNotice.set('Nota fiscal criada com sucesso.');
      await this.loadInvoices();
    } catch (error: any) {
      this.pageError.set(error?.error?.error?.message ?? 'Falha ao criar nota.');
    } finally {
      this.savingInvoice.set(false);
    }
  }

  logout(): void {
    this.authService.logout();
    void this.router.navigate(['/']);
  }

  trackByProductCode(_: number, product: Product): string {
    return product.code;
  }

  trackByInvoiceNumber(_: number, invoice: Invoice): number {
    return invoice.number;
  }

  private async loadProducts(): Promise<void> {
    try {
      this.loadingProducts.set(true);
      const response = await firstValueFrom(this.productService.list());
      this.products.set(response.data);
    } catch (error: any) {
      this.pageError.set(error?.error?.error?.message ?? 'Falha ao carregar produtos.');
    } finally {
      this.loadingProducts.set(false);
    }
  }

  private async loadInvoices(): Promise<void> {
    try {
      this.loadingInvoices.set(true);
      const response = await firstValueFrom(this.invoiceService.list());
      this.invoices.set(response.data);
    } catch (error: any) {
      this.pageError.set(error?.error?.error?.message ?? 'Falha ao carregar notas.');
    } finally {
      this.loadingInvoices.set(false);
    }
  }

  private createInvoiceItemGroup() {
    return this.fb.nonNullable.group({
      product_code: ['', [Validators.required]],
      quantity: [1, [Validators.required, Validators.min(1)]],
    });
  }
}
