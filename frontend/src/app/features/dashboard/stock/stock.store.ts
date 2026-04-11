import { computed, inject, Injectable, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import { ProductService } from '../../../core/services/product.service';
import { CreateProductRequest, Product } from '../../../core/models/product.models';
import { DashboardUiStore } from '../dashboard-ui.store';

@Injectable()
export class StockStore {
  private readonly productService = inject(ProductService);
  private readonly ui = inject(DashboardUiStore);

  readonly products = signal<Product[]>([]);
  readonly loadingProducts = signal(false);
  readonly savingProduct = signal(false);
  readonly updatingProduct = signal(false);
  readonly deletingProduct = signal(false);

  readonly productCount = computed(() => this.products().length);
  readonly hasProducts = computed(() => this.products().length > 0);

  async refresh(): Promise<void> {
    await this.loadProducts();
  }

  async createProduct(payload: CreateProductRequest): Promise<boolean> {
    try {
      this.savingProduct.set(true);
      this.ui.clearError();

      await firstValueFrom(this.productService.create(payload));
      this.ui.showNotice('Produto cadastrado com sucesso.');
      await this.loadProducts();
      return true;
    } catch (error: any) {
      this.ui.setError(error?.error?.error?.message ?? 'Falha ao cadastrar produto.');
      return false;
    } finally {
      this.savingProduct.set(false);
    }
  }

  async updateProduct(originalCode: string, payload: CreateProductRequest): Promise<boolean> {
    try {
      this.updatingProduct.set(true);
      this.ui.clearError();

      await firstValueFrom(this.productService.update(originalCode, payload));
      this.ui.showNotice('Produto atualizado com sucesso.');
      await this.loadProducts();
      return true;
    } catch (error: any) {
      this.ui.setError(error?.error?.error?.message ?? 'Falha ao atualizar produto.');
      return false;
    } finally {
      this.updatingProduct.set(false);
    }
  }

  async deleteProduct(code: string): Promise<boolean> {
    try {
      this.deletingProduct.set(true);
      this.ui.clearError();

      await firstValueFrom(this.productService.delete(code));
      this.ui.showNotice('Produto removido com sucesso.');
      await this.loadProducts();
      return true;
    } catch (error: any) {
      this.ui.setError(error?.error?.error?.message ?? 'Falha ao remover produto.');
      return false;
    } finally {
      this.deletingProduct.set(false);
    }
  }

  private async loadProducts(): Promise<void> {
    try {
      this.loadingProducts.set(true);
      const response = await firstValueFrom(this.productService.list());
      this.products.set(response.data);
    } catch (error: any) {
      this.ui.setError(error?.error?.error?.message ?? 'Falha ao carregar produtos.');
    } finally {
      this.loadingProducts.set(false);
    }
  }
}
