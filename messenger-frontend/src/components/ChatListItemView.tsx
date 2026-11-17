import { Link } from "react-router-dom";
import type { ChatListItem, User } from "../api/types";

interface ChatListItemProps {
    chat: ChatListItem;
    usersById: Record<string, User>;
    userId: string | undefined;
}

export function ChatListItemView({ chat, usersById, userId }: ChatListItemProps) {

    return (
        <Link to={`/chats/${chat.chat_id}`}
              style={{ display: "block", padding: 12, textDecoration: "none", color: "inherit" }}>
            <strong>{chat.chat_title ?? "Личный чат"}</strong>
            <div>{chat.msg_text ?? "Сообщений пока нет"}</div>
            {chat.msg_sender_id === userId &&
                chat.msg_id &&
                chat.last_read_message_id &&
                (chat.last_read_message_id >= chat.msg_id) && (
                    <span style={{
                        color: "#2697E8",
                        fontSize: 12
                    }}>
                        ✓✓
                    </span>
                )}
            {chat.msg_sender_id && (
                <div>
                    {usersById[chat.msg_sender_id]?.name ?? "Загрузка пользователя..."}
                </div>
            )}
            {chat.unread_msg_cnt > 0 && <span>{chat.unread_msg_cnt}</span>}
        </Link>
    );
}