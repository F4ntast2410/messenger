import { useCallback, useEffect, useRef, useState, type SubmitEvent } from "react";
import type { Message, MessageReadPayload } from "../api/types";
import { getMessages, markMessagesAsRead, sendMessage } from "../api/chats";
import { getUserById } from "../api/users"
import { useOutletContext, useParams } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { MessageBubble } from "../components/MessageBubble";
import { useSocket } from "../context/WebSocketContext";

function isNewerMessage(a: string, b: string): boolean {
    return a > b;
}

export function ChatPage() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [sendError, setSendError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [content, setContent] = useState('');
    const { chatId } = useParams();
    const [typingUser, setTypingUser] = useState<string | null>(null);

    const [visibleMessages, setVisibleMessages] = useState<Set<string>>(
        () => new Set()
    );

    const lastMarkedRead = useRef<string | null>(null);
    const pendingRead = useRef<string | null>(null);
    const messagesContainerRef = useRef<HTMLDivElement | null>(null);
    const [readCursors, setReadCursors] = useState<Record<string, string>>({});

    const typingTimeout = useRef<number | null>(null);
    const lastTypingSent = useRef(0);
    const { subscribe, send } = useSocket();
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
    useEffect(() => {
        setTypingUser(null);
        lastTypingSent.current = 0;

        if (typingTimeout.current !== null) {
            clearTimeout(typingTimeout.current);
            typingTimeout.current = null;
        }

        const unsubscribe = subscribe((event) => {
            if (event.chat_id !== chatId) return;

            if (event.type === "NEW_MESSAGE") {
                const msg = event.payload as Message;

                setMessages((prev) =>
                    prev.some((m) => m.id === msg.id)
                        ? prev
                        : [msg, ...prev]
                );

                setTypingUser(null);
                return;
            }
            if (event.type === "MESSAGE_READ") {
                const payload = event.payload as MessageReadPayload;

                setReadCursors((prev) => {
                    const currentCursor = prev[event.sender_id];

                    if (
                        currentCursor &&
                        payload.message_id <= currentCursor
                    ) {
                        return prev;
                    }

                    return {
                        ...prev,
                        [event.sender_id]: payload.message_id,
                    };
                });

                return;
            }
            (async () => {
                if (event.type === "USER_TYPING") {
                const user = await getUserById(event.sender_id);
                setTypingUser(user.name);

                if (typingTimeout.current !== null) {
                    clearTimeout(typingTimeout.current);
                }

                typingTimeout.current = window.setTimeout(() => {
                    setTypingUser(null);
                    typingTimeout.current = null;
                }, 3000);
            }
            })();
        });

        return () => {
            unsubscribe();

            if (typingTimeout.current !== null) {
                clearTimeout(typingTimeout.current);
                typingTimeout.current = null;
            }
        };
    }, [chatId, subscribe]);
    
    const { user } = useAuth();
    const ordered = [...messages].reverse();
    useEffect(() => {
        if (!chatId || isLoading || ordered.length === 0) return;

        const newestVisible = [...ordered]
            .reverse()
            .find((message) => visibleMessages.has(message.id));

        if (!newestVisible) return;

        const currentCursor = pendingRead.current ?? lastMarkedRead.current;

        if (
            currentCursor &&
            !isNewerMessage(newestVisible.id, currentCursor)
        ) {
            return;
        }

        pendingRead.current = newestVisible.id;

        markMessagesAsRead(chatId, newestVisible.id)
            .then(() => {
                lastMarkedRead.current = newestVisible.id;
            })
            .catch((err) => {
                console.error("Не удалось обновить статус прочтения:", err);
            })
            .finally(() => {
                if (pendingRead.current === newestVisible.id) {
                    pendingRead.current = null;
                }
            });
    }, [chatId, isLoading, messages, visibleMessages]);
    const { onMessageSent } = useOutletContext<{ onMessageSent: (msg: Message) => void }>();

    function handleTyping(value: string) {
        setContent(value);

        if (!chatId || !value.trim()) return;

        const now = Date.now();

        if (now - lastTypingSent.current < 2000) {
            return;
        }

        lastTypingSent.current = now;

        send({
            type: "USER_TYPING",
            chat_id: chatId,
            payload: {},
        });
    }
    async function handleSubmit(e: SubmitEvent<HTMLFormElement>) {
        e.preventDefault();
        if (!chatId) return;
        setSendError(null);
        setIsSubmitting(true);

        try {
        const newMsg = await sendMessage(content, chatId);
        setMessages((prev) => [newMsg, ...prev]);
        onMessageSent(newMsg);
        setContent('');
        lastTypingSent.current = 0;
        } catch (err) {
            setSendError('Что-то пошло не так, попробуй ещё раз');
        } finally {
        setIsSubmitting(false);
        }
    }
    const handleMessageVisible = useCallback((messageId: string) => {
        setVisibleMessages((prev) => {
            if (prev.has(messageId)) return prev;

            const next = new Set(prev);
            next.add(messageId);
            return next;
        });
    }, []);
    return (
        <div style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
            
            <div ref={messagesContainerRef} style={{ flex: 1, overflowY: "auto", padding: 16, minHeight: 24 }}>
                {isLoading && <p>Загрузка...</p>}
                {error && <p>{error}</p>}
                {ordered.map((m) => (
                    <MessageBubble key={m.id} message={m} isOwn={m.sender_id === user?.id} onVisible={handleMessageVisible} rootRef={messagesContainerRef} />
                ))}
            </div>
            {typingUser && (
                <p style={{ margin: "4px 16px" }}>
                    {typingUser} печатает...
                </p>
            )}
            <form onSubmit={handleSubmit}>
                <div>
                <label htmlFor="message-form-text">Введите сообщение...</label>
                <input id="message-form-text" type="text" value={content}
                        onChange={(e) => handleTyping(e.target.value)} required />
                </div>

                {sendError && <p style={{ color: 'red' }}>{sendError}</p>}

                <button type="submit" disabled={isSubmitting}>
                {isSubmitting ? 'Отправка...' : 'Отправить'}
                </button>
            </form>
        </div>
    );
}