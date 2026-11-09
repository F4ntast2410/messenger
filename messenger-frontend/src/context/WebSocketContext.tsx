import { createContext, useCallback, useContext, useEffect, useMemo, useRef, type ReactNode } from "react";
import type { ChatEvent } from "../api/types";

type Handler = (event: ChatEvent) => void;

interface OutgoingEvent {
    type: string;
    chat_id: string;
    payload?: unknown;
}

interface WebSocketContextValue {
    subscribe: (handler: Handler) => () => void;
    send: (event: OutgoingEvent) => void;
}

const WebSocketContext = createContext<WebSocketContextValue | null>(null);

export function WebSocketProvider({ chatIds, children }: { chatIds: string[]; children: ReactNode }) {
    const chatIdsKey = chatIds.join(",");
    const socketRef = useRef<WebSocket | null>(null);
    const handlersRef = useRef<Set<Handler>>(new Set());

    useEffect(() => {
        if (!chatIdsKey) return;
        let closedByUs = false;
        let retryTimer: number | undefined;

        function connect() {
            const protocol = location.protocol === "https:" ? "wss" : "ws";
            const ws = new WebSocket(`${protocol}://${location.host}/api/v1/ws?chat_ids=${chatIdsKey}`);
            socketRef.current = ws;

            ws.onmessage = (e) => {
                const event: ChatEvent = JSON.parse(e.data);
                handlersRef.current.forEach((handler) => handler(event));
            };

            ws.onclose = () => {
                if (closedByUs) return;
                retryTimer = window.setTimeout(connect, 2000);
            };
            
            
        }

        connect();
        
        return () => {
            closedByUs = true;
            clearTimeout(retryTimer);
            socketRef.current?.close();
        };
    }, [chatIdsKey]);

    const subscribe = useCallback((handler: Handler) => {
        handlersRef.current.add(handler);
        return () => {
            handlersRef.current.delete(handler);
        };
    }, []);

    const send = useCallback((event: OutgoingEvent) => {
        if (socketRef.current?.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify(event));
        }
    }, []);

    const value = useMemo(() => ({ subscribe, send }), [subscribe, send]);

    return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
}

export function useSocket() {
    const ctx = useContext(WebSocketContext);
    if (!ctx) throw new Error("useSocket must be used within WebSocketProvider");
    return ctx;
}