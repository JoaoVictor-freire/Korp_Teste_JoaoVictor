export interface User {
  id: string;
  name: string;
  email: string;
  created_at?: string;
}

export interface AuthPayload {
  user: User;
  token: string;
}

export interface Envelope<T> {
  data: T;
}
