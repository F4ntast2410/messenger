import { request } from './client';
import type { User } from './types';

export function login(email: string, password: string): Promise<void> {
  return request('/users/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export function register(email: string, username: string, name: string, password: string): Promise<void> {
  return request('/users/register', {
    method: 'POST',
    body: JSON.stringify({ email, username, name, password }),
  });
}

export function logout(): Promise<void> {
  return request('/users/logout', { method: 'POST' });
}

export function getMe(): Promise<User> {
  return request<User>('/users/me');
}