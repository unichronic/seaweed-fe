'use client';

import { useEffect, useState } from 'react';
import api from '@/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import Link from 'next/link';
import { Plus, Settings, Eye } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

interface Contest {
    id: string;
    name: string;
    start_time: number;
    end_time: number;
}

export default function AdminContestListPage() {
    const [contests, setContests] = useState<Contest[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchContests = async () => {
            try {
                const res = await api.get('/admin/contests/list'); 
                setContests(res.data);
            } catch (err) {
                console.error('Failed to fetch admin contests', err);
            } finally {
                setLoading(false);
            }
        };
        fetchContests();
    }, []);

    return (
        <div className="container mx-auto p-6 space-y-8">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold">Admin Dashboard</h1>
                    <p className="text-muted-foreground">Manage contests and problem sets.</p>
                </div>
                <Link href="/admin/contests/new">
                    <Button className="gap-2">
                        <Plus className="w-4 h-4" /> Create Contest
                    </Button>
                </Link>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>All Contests</CardTitle>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Contest Name</TableHead>
                                <TableHead>Duration</TableHead>
                                <TableHead className="text-right">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {contests.map((c) => (
                                <TableRow key={c.id}>
                                    <TableCell className="font-medium">{c.name}</TableCell>
                                    <TableCell>
                                        {new Date(c.start_time).toLocaleDateString()} - {new Date(c.end_time).toLocaleDateString()}
                                    </TableCell>
                                    <TableCell className="text-right space-x-2">
                                        <Link href={`/contests/${c.id}`}>
                                            <Button variant="ghost" size="sm" className="gap-1">
                                                <Eye className="w-4 h-4" /> View
                                            </Button>
                                        </Link>
                                        <Link href={`/admin/contests/${c.id}/problems/new`}>
                                            <Button variant="outline" size="sm" className="gap-1">
                                                <Plus className="w-4 h-4" /> Add Problem
                                            </Button>
                                        </Link>
                                        <Button variant="ghost" size="sm">
                                            <Settings className="w-4 h-4" />
                                        </Button>
                                    </TableCell>
                                </TableRow>
                            ))}
                            {contests.length === 0 && !loading && (
                                <TableRow>
                                    <TableCell colSpan={3} className="text-center py-12 text-muted-foreground">
                                        No contests created yet.
                                    </TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>
        </div>
    );
}
