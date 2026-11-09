import { request } from "./client";
import type { User } from "./types";

export function getUserById(userId: string): Promise<User> {
  return request<User>(`/users/${userId}`)
}

export function getUsersByIDs(userIDs: string[]): Promise<User[]> {
  return request<User[]>('/users/batch', {
    method: 'POST',
    body: JSON.stringify({ user_ids: userIDs }),
  });
} 