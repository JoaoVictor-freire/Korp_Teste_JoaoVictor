import { CommonModule } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { CreateProductRequest, Product } from '../../../core/models/product.models';
import { StockStore } from '../../dashboard/stock/stock.store';

type ProductSort = 'name' | 'stock';
const PAGE_SIZE = 20;

@Component({
  selector: 'app-history-products-page',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './history-products-page.component.html',
  styleUrl: './history-products-page.component.scss',
})
export class HistoryProductsPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly store = inject(StockStore);

  readonly searchTerm = signal('');
  readonly sortBy = signal<ProductSort>('name');
  readonly currentPage = signal(1);
  readonly editingProductCode = signal<string | null>(null);
  readonly deletingProductCode = signal<string | null>(null);
  readonly editForm = this.fb.nonNullable.group({
    code: ['', [Validators.required]],
    description: ['', [Validators.required]],
    stock: [0, [Validators.required, Validators.min(0)]],
  });

  readonly filteredProducts = computed(() => {
    const query = this.searchTerm().trim().toLowerCase();
    const currentSort = this.sortBy();

    const filtered = this.store.products().filter((product) => {
      if (!query) {
        return true;
      }

      return product.description.toLowerCase().includes(query) || product.code.toLowerCase().includes(query);
    });

    return [...filtered].sort((left: Product, right: Product) => {
      if (currentSort === 'stock') {
        return right.stock - left.stock;
      }

      return left.description.localeCompare(right.description);
    });
  });

  readonly totalPages = computed(() => {
    const total = Math.ceil(this.filteredProducts().length / PAGE_SIZE);
    return Math.max(1, total);
  });

  readonly paginatedProducts = computed(() => {
    const safePage = Math.min(this.currentPage(), this.totalPages());
    const start = (safePage - 1) * PAGE_SIZE;
    return this.filteredProducts().slice(start, start + PAGE_SIZE);
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

  updateSearchTerm(value: string): void {
    this.searchTerm.set(value);
    this.currentPage.set(1);
  }

  updateSortBy(value: ProductSort): void {
    this.sortBy.set(value);
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

  trackByProductCode(_: number, product: Product): string {
    return product.code;
  }

  openEditModal(product: Product): void {
    this.editingProductCode.set(product.code);
    this.editForm.setValue({
      code: product.code,
      description: product.description,
      stock: product.stock,
    });
  }

  closeEditModal(): void {
    this.editingProductCode.set(null);
    this.editForm.reset({
      code: '',
      description: '',
      stock: 0,
    });
  }

  openDeleteModal(product: Product): void {
    this.deletingProductCode.set(product.code);
  }

  closeDeleteModal(): void {
    this.deletingProductCode.set(null);
  }

  editingProduct(): Product | null {
    const code = this.editingProductCode();
    if (!code) {
      return null;
    }

    return this.store.products().find((product) => product.code === code) ?? null;
  }

  deletingProduct(): Product | null {
    const code = this.deletingProductCode();
    if (!code) {
      return null;
    }

    return this.store.products().find((product) => product.code === code) ?? null;
  }

  async submitEdit(): Promise<void> {
    const originalCode = this.editingProductCode();
    if (!originalCode) {
      return;
    }

    if (this.editForm.invalid) {
      this.editForm.markAllAsTouched();
      return;
    }

    const payload: CreateProductRequest = this.editForm.getRawValue();
    const success = await this.store.updateProduct(originalCode, payload);
    if (success) {
      this.closeEditModal();
    }
  }

  async confirmDelete(): Promise<void> {
    const product = this.deletingProduct();
    if (!product) {
      return;
    }

    const success = await this.store.deleteProduct(product.code);
    if (success) {
      this.closeDeleteModal();
    }
  }
}
