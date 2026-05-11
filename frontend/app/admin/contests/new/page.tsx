'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import api from '@/api';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

export default function NewContestPage() {
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        eligible_to: '',
        registration_status: 'open',
        registration_start_time: '',
        registration_end_time: '',
        start_time: '',
        end_time: '',
    });
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            await api.post('/admin/contest', {
                ...formData,
                registration_start_time: new Date(formData.registration_start_time).getTime(),
                registration_end_time: new Date(formData.registration_end_time).getTime(),
                start_time: new Date(formData.start_time).getTime(),
                end_time: new Date(formData.end_time).getTime(),
            });
            router.push('/admin/contests');
        } catch (err) {
            console.error('Failed to create contest', err);
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setFormData({ ...formData, [e.target.id]: e.target.value });
    };

    const handleStatusChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setFormData({ ...formData, registration_status: e.target.value });
    };

    return (
        <div className="container mx-auto max-w-2xl p-6">
            <Card>
                <CardHeader>
                    <CardTitle className="text-2xl font-black">Create New Contest</CardTitle>
                    <CardDescription>Set up a scored coding round.</CardDescription>
                </CardHeader>
                <form onSubmit={handleSubmit}>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="name">Contest Name</Label>
                            <Input id="name" placeholder="Weekly Rated Round 42" value={formData.name} onChange={handleChange} required />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="description">Description (Markdown supported)</Label>
                            <Textarea id="description" placeholder="Short overview of the contest..." value={formData.description} onChange={handleChange} />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="eligible_to">Eligible Semesters / Years (e.g. 4,6,8)</Label>
                            <Input id="eligible_to" value={formData.eligible_to} onChange={handleChange} />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="registration_status">Registration Status</Label>
                            <select
                                id="registration_status"
                                className="h-8 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm"
                                value={formData.registration_status}
                                onChange={handleStatusChange}
                            >
                                <option value="open">Open</option>
                                <option value="closed">Closed</option>
                                <option value="invite-only">Invite only</option>
                            </select>
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="registration_start_time">Reg Start</Label>
                                <Input id="registration_start_time" type="datetime-local" value={formData.registration_start_time} onChange={handleChange} required />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="registration_end_time">Reg End</Label>
                                <Input id="registration_end_time" type="datetime-local" value={formData.registration_end_time} onChange={handleChange} required />
                            </div>
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="start_time">Contest Start</Label>
                                <Input id="start_time" type="datetime-local" value={formData.start_time} onChange={handleChange} required />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="end_time">Contest End</Label>
                                <Input id="end_time" type="datetime-local" value={formData.end_time} onChange={handleChange} required />
                            </div>
                        </div>
                    </CardContent>
                    <CardFooter>
                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? 'Creating...' : 'Create Contest'}
                        </Button>
                    </CardFooter>
                </form>
            </Card>
        </div>
    );
}
