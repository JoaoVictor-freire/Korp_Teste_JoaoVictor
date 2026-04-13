import { computed, inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';

import { apiConfig } from '../config/api.config';
import { AuthPayload, Envelope, User } from '../models/auth.models';

interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

const TOKEN_KEY = 'korp.token';
const USER_KEY = 'korp.user';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly tokenState = signal<string | null>(this.readStoredToken());
  private readonly userState = signal<User | null>(this.readStoredUser());

  readonly token = computed(() => this.tokenState());
  readonly user = computed(() => this.userState());
  readonly isAuthenticated = computed(() => this.isTokenValid(this.tokenState()));

  register(payload: RegisterRequest): Observable<Envelope<AuthPayload>> {
    return this.http
      .post<Envelope<AuthPayload>>(`${apiConfig.stockBaseUrl}/api/v1/auth/register`, payload)
      .pipe(tap((response) => this.persistSession(response.data)));
  }

  login(payload: LoginRequest): Observable<Envelope<AuthPayload>> {
    return this.http
      .post<Envelope<AuthPayload>>(`${apiConfig.stockBaseUrl}/api/v1/auth/login`, payload)
      .pipe(tap((response) => this.persistSession(response.data)));
  }

  logout(): void {
    this.clearSession();
  }

  private persistSession(payload: AuthPayload): void {
    this.tokenState.set(payload.token);
    this.userState.set(payload.user);
    localStorage.setItem(TOKEN_KEY, payload.token);
    localStorage.setItem(USER_KEY, JSON.stringify(payload.user));
  }

  private readStoredToken(): string | null {
    const token = localStorage.getItem(TOKEN_KEY);
    if (!this.isTokenValid(token)) {
      this.clearStoredSession();
      return null;
    }

    return token;
  }

  private readStoredUser(): User | null {
    if (!this.isTokenValid(localStorage.getItem(TOKEN_KEY))) {
      this.clearStoredSession();
      return null;
    }

    const raw = localStorage.getItem(USER_KEY);
    if (!raw) {
      return null;
    }

    try {
      return JSON.parse(raw) as User;
    } catch {
      this.clearStoredSession();
      return null;
    }
  }

  private clearSession(): void {
    this.tokenState.set(null);
    this.userState.set(null);
    this.clearStoredSession();
  }

  private clearStoredSession(): void {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  }

  private isTokenValid(token: string | null): boolean {
    if (!token) {
      return false;
    }

    const payload = this.decodeJwtPayload(token);
    if (!payload) {
      return false;
    }

    if (typeof payload.exp !== 'number') {
      return false;
    }

    return payload.exp * 1000 > Date.now();
  }

  private decodeJwtPayload(token: string): { exp?: number } | null {
    const parts = token.split('.');
    if (parts.length !== 3) {
      return null;
    }

    try {
      const normalized = parts[1].replace(/-/g, '+').replace(/_/g, '/');
      const padded = normalized.padEnd(normalized.length + ((4 - (normalized.length % 4)) % 4), '=');
      return JSON.parse(atob(padded)) as { exp?: number };
    } catch {
      return null;
    }
  }
}
