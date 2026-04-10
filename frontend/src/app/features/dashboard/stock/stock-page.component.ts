import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { Product } from '../../../core/models/product.models';
import { DashboardUiStore } from '../dashboard-ui.store';
import { StockStore } from './stock.store';

@Component({
  selector: 'app-stock-page',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './stock-page.component.html',
  styleUrl: './stock-page.component.scss',
})
export class StockPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly ui = inject(DashboardUiStore);
  readonly store = inject(StockStore);

  readonly products = this.store.products;
  readonly loadingProducts = this.store.loadingProducts;
  readonly savingProduct = this.store.savingProduct;

  readonly productForm = this.fb.nonNullable.group({
    code: ['', [Validators.required]],
    description: ['', [Validators.required]],
    stock: [0, [Validators.required, Validators.min(0)]],
  });

  async submitProduct(): Promise<void> {
    if (this.productForm.invalid) {
      this.productForm.markAllAsTouched();
      return;
    }

    await this.store.createProduct(this.productForm.getRawValue());
    if (!this.ui.pageError()) {
      this.productForm.reset({ code: '', description: '', stock: 0 });
    }
  }

  trackByProductCode(_: number, product: Product): string {
    return product.code;
  }
}
