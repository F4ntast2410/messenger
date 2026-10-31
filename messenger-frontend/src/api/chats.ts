import { request } from "./client";
import type { ChatListItem, Message } from "./types";



export function getChats(): Promise<ChatListItem[]> {
    return request<ChatListItem[]>('/chats')
}

export function getMessages(chatId: string, limit = 30, beforeId?: string): Promise<Message[]> {
    const params = new URLSearchParams({ limit: String(limit) });
    if (beforeId) {
        params.set('before_id', beforeId);
    }
    return request<Message[]>(`/chats/${chatId}/messages?${params}`);
}