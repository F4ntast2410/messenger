import { useEffect, useState } from "react";
import type { Message } from "../api/types";
import { getMessages } from "../api/chats";
import { useParams } from "react-router-dom";


export function ChatPage() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const { chatId } = useParams();
    const limit = 30
    useEffect(() => {
        if (!chatId) return;
        (async () => {
            setIsLoading(true);
            try {
                const messageList = await getMessages(chatId, limit);
                setMessages(messageList);
            } catch(err) {
                console.error(err)
                setError('Ошибка при загрузке сообщений');
            } finally {
                setIsLoading(false);
            }

        })();

    }, [chatId])
    return (
        <div>
            <p>{messages.length}</p>
        </div>
    );
}