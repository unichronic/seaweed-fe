'use client';

import Link from 'next/link';
import { useAuthStore } from '@/store';
import { Button } from '@/components/ui/button';
import { auth } from '@/firebase';
import { signOut } from 'firebase/auth';
import { Braces } from 'lucide-react';

export function Navbar() {
    const { user, setUser } = useAuthStore();

    const handleLogout = async () => {
        if (process.env.NEXT_PUBLIC_AUTH_BYPASS === 'true') {
            setUser(null);
            return;
        }
        await signOut(auth);
    };

    return (
        <nav className="sticky top-0 z-50 flex items-center justify-between border-b bg-background/95 p-4 backdrop-blur supports-[backdrop-filter]:bg-background/80">
            <div className="flex items-center gap-6">
                <Link href="/" className="flex items-center gap-2 text-xl font-black text-foreground">
                    <span className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                        <Braces className="size-5" />
                    </span>
                    Seaweed Arena
                </Link>
                <div className="flex gap-4">
                    <Link href="/contests" className="text-sm font-medium hover:text-primary transition-colors">
                        Contests
                    </Link>
                </div>
            </div>
            <div className="flex items-center gap-4">
                {user ? (
                    <>
                        <span className="text-sm text-muted-foreground hidden sm:inline-block">
                            {user.email}
                        </span>
                        <Button variant="outline" onClick={handleLogout} size="sm">
                            Logout
                        </Button>
                    </>
                ) : (
                    <>
                        <Link href="/auth/login">
                            <Button variant="ghost" size="sm">Login</Button>
                        </Link>
                        <Link href="/auth/signup">
                            <Button size="sm">Join</Button>
                        </Link>
                    </>
                )}
            </div>
        </nav>
    );
}
