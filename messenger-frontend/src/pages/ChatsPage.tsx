import { useEffect, useRef, useState } from "react";
import { Outlet } from "react-router-dom";
import { getChats } from "../api/chats";
import type { ChatListItem, Message } from "../api/types";
import { WebSocketProvider } from "../context/WebSocketContext";
import { ChatSideBarView } from "../components/ChatSideBar";
import { applyMessage } from "../utils/chats";
import { getUsersByIDs } from "../api/users";
import type { User } from "../api/types";

export function ChatsPage() {
    const [chats, setChats] = useState<ChatListItem[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [usersById, setUsersById] = useState<Record<string, User>>({});
    const requestedUserIDs = useRef(new Set<string>());
    function onMessageSent(msg: Message) {
        setChats((prev) => applyMessage(prev, msg, false));
    }
    useEffect(() => {
        const userIDs = [
            ...new Set(
                chats
                    .map((chat) => chat.msg_sender_id)
                    .filter((id): id is string => id !== null)
            ),
        ];

        const missingIDs = userIDs.filter(
            (id) => !requestedUserIDs.current.has(id)
        );

        if (missingIDs.length === 0) return;

        missingIDs.forEach((id) => requestedUserIDs.current.add(id));

        getUsersByIDs(missingIDs)
            .then((users) => {
                setUsersById((prev) => {
                    const updated = { ...prev };

                    users.forEach((user) => {
                        updated[user.id] = user;
                    });

                    return updated;
                });
            })
            .catch((err) => {
                missingIDs.forEach((id) => requestedUserIDs.current.delete(id));
                console.error("Ошибка при загрузке пользователей:", err);
            });
    }, [chats]);
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
            <WebSocketProvider chatIds={chats.map((c) => c.chat_id)}>
                <ChatSideBarView
                    chats={chats}
                    setChats={setChats}
                    isLoading={isLoading}
                    error={error}
                    usersById={usersById}
                />

                <main style={{ flex: 1 }}>
                    <Outlet context={{ onMessageSent }} />
                </main>
            </WebSocketProvider>
        </div>
    );
}