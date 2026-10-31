import { useEffect, useState } from "react";
import { Link, Outlet } from "react-router-dom";
import { getChats } from "../api/chats";
import type { ChatListItem } from "../api/types";

export function ChatsPage() {
    const [chats, setChats] = useState<ChatListItem[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    useEffect(() => {
        (async () => {
            setIsLoading(true);
            try {
                const chatList = await getChats();
                setChats(chatList);
            } catch(err) {
                setError('Ошибка при загрузке списка чатов');
            } finally {
                setIsLoading(false);
            }

        })();

    }, [])
    return (
        <div style={{ display: "flex", height: "100vh" }}>
            <aside style={{ width: 300, borderRight: "1px solid #ccc", overflowY: "auto" }}>
                {isLoading && <p>Загрузка...</p>}
                {error && <p>{error}</p>}
                {!isLoading && !error && chats.length === 0 && <p>Чатов пока нет</p>}

                {chats.map((chat) => (
                    <Link key={chat.chat_id} to={`/chats/${chat.chat_id}`}
                        style={{ display: "block", padding: 12, textDecoration: "none", color: "inherit" }}>
                        <strong>{chat.chat_title ?? "Личный чат"}</strong>
                        <div>{chat.msg_text ?? "Сообщений пока нет"}</div>
                        {chat.unread_msg_cnt > 0 && <span>{chat.unread_msg_cnt}</span>}
                    </Link>
                ))}
            </aside>

            <main style={{ flex: 1 }}>
                <Outlet />
            </main>
        </div>
    );
}