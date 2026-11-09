import type { ChatListItem, Message } from "../api/types";

export function applyMessage(
    prev: ChatListItem[],
    msg: Message,
    isUnread: boolean,
): ChatListItem[] {
    const existing = prev.find((c) => c.chat_id === msg.chat_id);
    if (!existing) return prev;

    const updated: ChatListItem = {
        ...existing,
        msg_text: msg.content,
        msg_sent_at: msg.created_at,
        msg_updated_at: msg.updated_at,
        msg_sender_id: msg.sender_id,
        unread_msg_cnt: existing.unread_msg_cnt + (isUnread ? 1 : 0),
    };

    return [updated, ...prev.filter((c) => c.chat_id !== msg.chat_id)];
}