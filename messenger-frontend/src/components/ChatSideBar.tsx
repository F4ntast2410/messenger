import { useEffect, type Dispatch, type SetStateAction } from "react";
import type { ChatListItem, Message, MessageReadPayload, User } from "../api/types";
import { ChatListItemView } from "./ChatListItemView";
import { useSocket } from "../context/WebSocketContext";
import { useParams } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { applyMessage } from "../utils/chats";


interface ChatSideBarProps {
    chats: ChatListItem[];
    setChats: Dispatch<SetStateAction<ChatListItem[]>>;
    isLoading: boolean;
    error: string | null;
    usersById: Record<string, User>;
}

export function ChatSideBarView({chats, setChats, isLoading, error, usersById}: ChatSideBarProps) {
    const { subscribe } = useSocket();
    const { chatId } = useParams();
    const { user } = useAuth();
    useEffect(() => {
    return subscribe((event) => {
        if (event.type === "NEW_MESSAGE") {
            const msg = event.payload as Message;
            const isUnread = event.sender_id !== user?.id;
            setChats((prev) => applyMessage(prev, msg, isUnread));
            return;
        }
        if (event.type === "MESSAGE_READ") {
            const payload = event.payload as MessageReadPayload;
            const msgId = payload.message_id;
            const senderId = event.sender_id;
            if (senderId !== user?.id) {
                setChats((prev) => {
                    const existing = prev.find((c) => c.chat_id === event.chat_id);
                    if (!existing) return prev;
                    if (existing.last_read_message_id && existing.last_read_message_id >= msgId) return prev;
                    const updated: ChatListItem = {
                        ...existing,
                        last_read_message_id: msgId,
                    };

                    return [updated, ...prev.filter((c) => c.chat_id !== event.chat_id)];
                });
                return;
            }
        } 
    });
}, [subscribe, chatId, user?.id]);
    return (
        <aside style={{ width: 300, borderRight: "1px solid #ccc", overflowY: "auto" }}>
            {isLoading && <p>Загрузка...</p>}
            {error && <p>{error}</p>}
            {!isLoading && !error && chats.length === 0 && <p>Чатов пока нет</p>}

            {chats.map((chat) => (
                <ChatListItemView key={chat.chat_id} chat={chat} usersById={usersById} userId={user?.id}/>
            ))}
        </aside>
    );
}