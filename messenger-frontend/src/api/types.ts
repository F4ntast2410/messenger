export interface User {
  id: string;
  name: string;
  username: string;
  email: string;
  created_at: string;
}

export interface ChatListItem {
  chat_id: string;
  chat_title: string | null;
  chat_type: string;
  msg_text: string | null;
  msg_sent_at: string | null;
  msg_updated_at: string | null;
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