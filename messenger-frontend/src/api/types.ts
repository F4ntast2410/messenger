export interface User {
  id: string;
  name: string;
  username: string;
  email: string | null;
  created_at: string | null;
}

export interface ChatListItem {
  chat_id: string;
  chat_title: string | null;
  chat_type: string;
  msg_text: string | null;
  msg_sent_at: string | null;
  msg_updated_at: string | null;
  msg_sender_id: string | null;
  unread_msg_cnt: number;
}

export interface SearchUserResult {
  id: string;
  name: string;
  username: string;
}

export interface Message {
  id: string;
  chat_id: string;
  sender_id: string;
  content: string;
  created_at: string;
  updated_at: string | null;
  reply_to_id: string | null;
  is_deleted: boolean;
}

export type EventType = 'NEW_MESSAGE' | 'USER_TYPING' | 'MESSAGE_READ';

export interface ChatEvent {
    type: EventType;
    chat_id: string;
    sender_id: string;
    payload: unknown;
    created_at: string;
}
export interface MessageReadPayload {
    message_id: string;
}