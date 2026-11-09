import { useEffect, useRef, type RefObject } from "react";
import type { Message } from "../api/types";

interface MessageBubbleProps {
    message: Message;
    isOwn: boolean;
    onVisible: (messageId: string) => void;
    rootRef: RefObject<HTMLDivElement | null>;
}

export function MessageBubble({
    message,
    isOwn,
    onVisible,
    rootRef
}: MessageBubbleProps) {
    const bubbleRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        if (isOwn) return;

        const element = bubbleRef.current;
        if (!element) return;

        const observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) {
                    onVisible(message.id);
                    observer.disconnect();
                }
            },
            {
                root: rootRef.current,
                threshold: 0.5,
            }
        );

        observer.observe(element);

        return () => observer.disconnect();
    }, [message.id, isOwn, onVisible]);

    return (
        <div
            ref={bubbleRef}
            style={{
                textAlign: isOwn ? "right" : "left",
                marginBottom: 15
            }}
        >
            <span style={{
                background: "#eee",
                padding: "6px 10px",
                borderRadius: 8
            }}>
                {message.is_deleted
                    ? "Сообщение удалено"
                    : message.content}
            </span>
        </div>
    );
}