import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { getMe, logout as apiLogout } from '../api/auth';
import { ApiError } from '../api/client';
import type { User } from '../api/types';

interface AuthContextValue {
    user: User | null;
    isLoading: boolean;
    logout: () => Promise<void>;
    refetchUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    async function checkAuth() {
        try {
            const me = await getMe();
            setUser(me);
        } catch (error) {
            if (error instanceof ApiError && error.status === 401) {
                setUser(null);
            } else {
                console.error('Failed to check auth:', error);
                setUser(null);
            }
        } finally {
            setIsLoading(false);
        }
    }

    useEffect(() => {
        checkAuth();
    }, []);

    async function logout() {
        await apiLogout();
        setUser(null);
    }

    const value: AuthContextValue = {
        user,
        isLoading,
        logout,
        refetchUser: checkAuth,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
    const ctx = useContext(AuthContext);
    if (!ctx) {
        throw new Error('useAuth must be used within AuthProvider');
    }
    return ctx;
}