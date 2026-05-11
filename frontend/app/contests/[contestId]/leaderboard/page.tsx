'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import api from '@/api';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Card, CardContent } from '@/components/ui/card';
import { Trophy, Medal, Award } from 'lucide-react';

interface Ranking {
    user_id: string;
    name: string;
    usn: string;
    department: string;
    score: number;
    correct_attempts: number;
    incorrect_attempts: number;
    rank: number;
}

export default function LeaderboardPage() {
    const { contestId } = useParams<{ contestId: string }>();
    const [rankings, setRankings] = useState<Ranking[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchRankings = async () => {
            try {
                const res = await api.get(`/contests/${contestId}/leaderboard`);
                setRankings(res.data);
            } catch (err) {
                console.error('Failed to fetch rankings', err);
            } finally {
                setLoading(false);
            }
        };
        fetchRankings();

        const interval = setInterval(fetchRankings, 10000);
        return () => clearInterval(interval);
    }, [contestId]);

    const getRankIcon = (rank: number) => {
        if (rank === 1) return <Trophy className="w-5 h-5 text-yellow-500" />;
        if (rank === 2) return <Medal className="w-5 h-5 text-gray-400" />;
        if (rank === 3) return <Award className="w-5 h-5 text-amber-600" />;
        return null;
    };

    if (loading) return <div className="p-12 text-center">Loading leaderboard...</div>;

    return (
        <div className="container mx-auto p-6 space-y-8">
            <div className="space-y-2 border-b border-foreground/20 pb-6 text-center">
                <h1 className="text-4xl font-black text-primary">Leaderboard</h1>
                <p className="text-sm font-bold text-muted-foreground">Real-time standings</p>
            </div>

            <Card>
                <CardContent className="p-0">
                    <Table>
                        <TableHeader className="bg-muted/50">
                            <TableRow>
                                <TableHead className="w-20 text-center font-bold">Rank</TableHead>
                                <TableHead className="font-bold">Participant</TableHead>
                                <TableHead className="font-bold hidden md:table-cell">USN</TableHead>
                                <TableHead className="font-bold hidden sm:table-cell">Department</TableHead>
                                <TableHead className="text-right font-bold hidden lg:table-cell">Attempts</TableHead>
                                <TableHead className="text-right font-bold pr-8">Score</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {rankings.map((r) => (
                                <TableRow key={r.user_id} className="hover:bg-muted/30 transition-colors">
                                    <TableCell className="text-center font-medium">
                                        <div className="flex items-center justify-center gap-2">
                                            {getRankIcon(r.rank)}
                                            {r.rank}
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        <div className="font-bold">{r.name}</div>
                                        <div className="text-xs text-muted-foreground md:hidden">{r.usn}</div>
                                    </TableCell>
                                    <TableCell className="hidden md:table-cell">{r.usn}</TableCell>
                                    <TableCell className="hidden sm:table-cell">{r.department}</TableCell>
                                    <TableCell className="text-right hidden lg:table-cell">{r.correct_attempts}/{r.incorrect_attempts}</TableCell>
                                    <TableCell className="text-right pr-8 font-mono font-bold text-primary">
                                        {r.score}
                                    </TableCell>
                                </TableRow>
                            ))}
                            {rankings.length === 0 && (
                                <TableRow>
                                    <TableCell colSpan={6} className="text-center py-12 text-muted-foreground">
                                        No submissions yet. Be the first to solve!
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
