import { request } from "./client";
import type { ChatListItem, Message  } from "./types";



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

export function sendMessage(content: string, chatId: string, replyToId?: string): Promise<Message> {
  return request<Message>(`/chats/${chatId}/messages`, {
    method: 'POST',
    body: JSON.stringify({ content, reply_to_id: replyToId }),
  });
}

export function markMessagesAsRead(
    chatId: string,
    messageId: string
): Promise<void> {
    return request<void>(`/chats/${chatId}/read`, {
        method: 'PUT',
        body: JSON.stringify({
            message_id: messageId,
        }),
    });
}