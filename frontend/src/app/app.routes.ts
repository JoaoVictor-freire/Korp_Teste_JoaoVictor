import { Routes } from '@angular/router';

import { authGuard } from './core/guards/auth.guard';
import { guestGuard } from './core/guards/guest.guard';
import { DashboardPageComponent } from './features/dashboard/dashboard-page.component';
import { HomepagePageComponent } from './features/homepage/homepage-page.component';
import { LoginPageComponent } from './features/auth/login/login-page.component';
import { RegisterPageComponent } from './features/auth/register/register-page.component';

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
  },
  {
    path: '**',
    redirectTo: '',
  },
];
