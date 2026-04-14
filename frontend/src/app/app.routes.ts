import { Routes } from '@angular/router';

import { authGuard } from './core/guards/auth.guard';
import { guestGuard } from './core/guards/guest.guard';
import { DashboardPageComponent } from './features/dashboard/dashboard-page.component';
import { HomepagePageComponent } from './features/homepage/homepage-page.component';
import { LoginPageComponent } from './features/auth/login/login-page.component';
import { RegisterPageComponent } from './features/auth/register/register-page.component';
import { StockPageComponent } from './features/dashboard/stock/stock-page.component';
import { InvoicesPageComponent } from './features/dashboard/invoices/invoices-page.component';
import { HistoryProductsPageComponent } from './features/history/products/history-products-page.component';
import { HistoryInvoicesPageComponent } from './features/history/invoices/history-invoices-page.component';
import { InsightsPageComponent } from './features/insights/insights-page.component';

export const routes: Routes = [
  {
    path: '',
    component: HomepagePageComponent,
  },
  {
    path: 'login',
    component: LoginPageComponent,
    canActivate: [guestGuard],
  },
  {
    path: 'register',
    component: RegisterPageComponent,
    canActivate: [guestGuard],
  },
  {
    path: 'dashboard',
    component: DashboardPageComponent,
    canActivate: [authGuard],
    children: [
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'stock',
      },
      {
        path: 'stock',
        component: StockPageComponent,
      },
      {
        path: 'invoices',
        component: InvoicesPageComponent,
      },
      {
        path: 'history-products',
        component: HistoryProductsPageComponent,
      },
      {
        path: 'history-invoices',
        component: HistoryInvoicesPageComponent,
      },
      {
        path: 'insights',
        component: InsightsPageComponent,
      },
    ],
  },
  {
    path: '**',
    redirectTo: '',
  },
];
