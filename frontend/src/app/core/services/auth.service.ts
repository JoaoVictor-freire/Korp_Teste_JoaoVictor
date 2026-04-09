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
  private readonly tokenState = signal<string | null>(localStorage.getItem(TOKEN_KEY));
  private readonly userState = signal<User | null>(this.readUser());

  readonly token = computed(() => this.tokenState());
  readonly user = computed(() => this.userState());
  readonly isAuthenticated = computed(() => !!this.tokenState());

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
    this.tokenState.set(null);
    this.userState.set(null);
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  }

  private persistSession(payload: AuthPayload): void {
    this.tokenState.set(payload.token);
    this.userState.set(payload.user);
    localStorage.setItem(TOKEN_KEY, payload.token);
    localStorage.setItem(USER_KEY, JSON.stringify(payload.user));
  }

  private readUser(): User | null {
    const raw = localStorage.getItem(USER_KEY);
    if (!raw) {
      return null;
    }

    try {
      return JSON.parse(raw) as User;
    } catch {
      localStorage.removeItem(USER_KEY);
      return null;
    }
  }
}
