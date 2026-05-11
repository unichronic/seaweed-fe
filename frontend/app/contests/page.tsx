'use client';

import { useEffect, useState } from 'react';
import api from '@/api';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import Link from 'next/link';
import { Calendar, Trophy, Users } from 'lucide-react';

interface Contest {
    id: string;
    name: string;
    description: string;
    start_time: number;
    end_time: number;
}

export default function ContestListPage() {
    const [contests, setContests] = useState<Contest[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchContests = async () => {
            try {
                const res = await api.get('/contests/list');
                setContests(res.data);
            } catch (err) {
                console.error('Failed to fetch contests', err);
            } finally {
                setLoading(false);
            }
        };
        fetchContests();
    }, []);

    if (loading) return <div className="flex justify-center p-12 font-bold">Loading contest board...</div>;

    return (
        <div className="container mx-auto space-y-8 p-6">
            <div className="flex justify-between items-center border-b border-foreground/20 pb-6">
                <div>
                    <h1 className="text-4xl font-black">Contest Board</h1>
                    <p className="text-muted-foreground">Pick a round, register, and start solving when the timer opens.</p>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {contests.length > 0 ? contests.map((contest) => (
                    <Card key={contest.id} className="flex flex-col border-2 border-foreground/15 transition-shadow hover:shadow-[8px_8px_0_hsl(216_36%_13%_/_0.12)]">
                        <CardHeader>
                            <CardTitle className="text-xl font-black">{contest.name}</CardTitle>
                            <CardDescription className="line-clamp-2">{contest.description}</CardDescription>
                        </CardHeader>
                        <CardContent className="flex-1 space-y-4">
                            <div className="flex items-center text-sm text-muted-foreground gap-2">
                                <Calendar className="w-4 h-4" />
                                <span>Starts: {new Date(contest.start_time).toLocaleString()}</span>
                            </div>
                            <div className="flex items-center text-sm text-muted-foreground gap-2">
                                <Users className="w-4 h-4" />
                                <span>Registration Open</span>
                            </div>
                        </CardContent>
                        <CardFooter className="pt-0">
                            <Link href={`/contests/${contest.id}`} className="w-full">
                                <Button className="w-full font-bold">Open Round</Button>
                            </Link>
                        </CardFooter>
                    </Card>
                )) : (
                    <div className="col-span-full rounded-lg border-2 border-dashed py-12 text-center">
                        <Trophy className="w-12 h-12 text-muted-foreground mx-auto mb-4 opacity-20" />
                        <p className="text-muted-foreground">No contests are open right now.</p>
                    </div>
                )}
            </div>
        </div>
    );
}
