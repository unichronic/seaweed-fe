'use client';

import React, { createContext, useContext, useEffect } from 'react';
import { onAuthStateChanged, User } from 'firebase/auth';
import { auth } from '@/firebase';
import { useAuthStore } from '@/store';

const AuthContext = createContext({});

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const { setUser, setLoading } = useAuthStore();

    useEffect(() => {
        if (process.env.NEXT_PUBLIC_AUTH_BYPASS === 'true') {
            setUser({
                uid: 'local-user-123',
                email: 'local@example.com',
                getIdToken: async () => 'dummy-token',
            } as unknown as User);
            setLoading(false);
            return;
        }

        if (auth) {
            const unsubscribe = onAuthStateChanged(auth, (user) => {
                setUser(user);
                setLoading(false);
            });

            return () => unsubscribe();
        }
    }, [setUser, setLoading]);

    return <AuthContext.Provider value={{}}>{children}</AuthContext.Provider>;
};

export const useAuth = () => useContext(AuthContext);
