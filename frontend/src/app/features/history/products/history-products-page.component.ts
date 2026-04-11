import { CommonModule } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { Product } from '../../../core/models/product.models';
import { StockStore } from '../../dashboard/stock/stock.store';

type ProductSort = 'name' | 'stock';

@Component({
  selector: 'app-history-products-page',
  imports: [CommonModule, FormsModule],
  templateUrl: './history-products-page.component.html',
  styleUrl: './history-products-page.component.scss',
})
export class HistoryProductsPageComponent {
  readonly store = inject(StockStore);

  readonly searchTerm = signal('');
  readonly sortBy = signal<ProductSort>('name');

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

  trackByProductCode(_: number, product: Product): string {
    return product.code;
  }
}
