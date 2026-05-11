'use client';

import { useState } from 'react';
import { createUserWithEmailAndPassword } from 'firebase/auth';
import { auth } from '@/firebase';
import { useRouter } from 'next/navigation';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import Link from 'next/link';
import axios from 'axios';

export default function SignupPage() {
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        usn: '',
        department: '',
        joining_year: '',
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const handleSignup = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        try {
            if (process.env.NEXT_PUBLIC_AUTH_BYPASS === 'true') {
                await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/users/create`, {
                    name: formData.name,
                    email: formData.email,
                    usn: formData.usn,
                    department: formData.department,
                    joining_year: parseInt(formData.joining_year, 10),
                }, {
                    headers: { Authorization: 'Bearer dummy-token' }
                });
                router.push('/contests');
                return;
            }

            const userCredential = await createUserWithEmailAndPassword(auth, formData.email, formData.password);
            const user = userCredential.user;
            const idToken = await user.getIdToken();

            await axios.post(`${process.env.NEXT_PUBLIC_API_URL}/users/create`, {
                name: formData.name,
                email: formData.email,
                usn: formData.usn,
                department: formData.department,
                joining_year: parseInt(formData.joining_year, 10),
            }, {
                headers: { Authorization: `Bearer ${idToken}` }
            });

            router.push('/contests');
        } catch (err: unknown) {
            if (axios.isAxiosError(err)) {
                setError(err.response?.data?.message || err.message);
            } else {
                setError(err instanceof Error ? err.message : 'Failed to create account');
            }
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [e.target.id]: e.target.value });
    };

    return (
        <div className="flex items-center justify-center min-h-[calc(100vh-64px)] p-4 bg-muted/30">
            <Card className="w-full max-w-lg">
                <CardHeader className="space-y-1">
                    <CardTitle className="text-2xl font-bold">Sign Up</CardTitle>
                    <CardDescription>
                        Create an account to start participating in contests
                    </CardDescription>
                </CardHeader>
                <form onSubmit={handleSignup}>
                    <CardContent className="space-y-4">
                        {error && <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-md">{error}</div>}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="name">Full Name</Label>
                                <Input id="name" placeholder="John Doe" value={formData.name} onChange={handleChange} required />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="email">Email</Label>
                                <Input id="email" type="email" placeholder="john@example.com" value={formData.email} onChange={handleChange} required />
                            </div>
                        </div>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="password">Password</Label>
                                <Input id="password" type="password" value={formData.password} onChange={handleChange} required />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="usn">USN / Serial Number</Label>
                                <Input id="usn" placeholder="1XY22CS001" value={formData.usn} onChange={handleChange} required />
                            </div>
                        </div>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="department">Department</Label>
                                <Input id="department" placeholder="Computer Science" value={formData.department} onChange={handleChange} required />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="joining_year">Joining Year</Label>
                                <Input id="joining_year" type="number" placeholder="2022" value={formData.joining_year} onChange={handleChange} required />
                            </div>
                        </div>
                    </CardContent>
                    <CardFooter className="flex flex-col gap-4">
                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? 'Creating account...' : 'Create Account'}
                        </Button>
                        <div className="text-sm text-center text-muted-foreground">
                            Already have an account?{' '}
                            <Link href="/auth/login" className="underline hover:text-primary transition-colors">
                                Login
                            </Link>
                        </div>
                    </CardFooter>
                </form>
            </Card>
        </div>
    );
}
