'use client';

import { useCallback, useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import api from '@/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Calendar, Trophy, ListChecks, ArrowRight, UserPlus, UserMinus } from 'lucide-react';

interface Problem {
    id: string;
    name: string;
    score: number;
}

interface Contest {
    id: string;
    name: string;
    description: string;
    registration_start_time: number;
    registration_end_time: number;
    registration_status: string;
    start_time: number;
    end_time: number;
    finalized: boolean;
    registered: boolean;
}

export default function ContestDetailsPage() {
    const { contestId } = useParams<{ contestId: string }>();
    const [contest, setContest] = useState<Contest | null>(null);
    const [problems, setProblems] = useState<Problem[]>([]);
    const [loading, setLoading] = useState(true);
    const [now, setNow] = useState(() => Date.now());
    const [registrationLoading, setRegistrationLoading] = useState(false);

    const fetchProblems = useCallback(async () => {
        const pRes = await api.get(`/contests/${contestId}/problems`);
        setProblems(pRes.data);
    }, [contestId]);

    const fetchDetails = useCallback(async () => {
        try {
            const cRes = await api.get(`/contests/${contestId}`);
            setContest(cRes.data);
            if (cRes.data.registered) {
                await fetchProblems();
            } else {
                setProblems([]);
            }
        } catch (err) {
            console.error('Failed to fetch contest details', err);
        } finally {
            setLoading(false);
        }
    }, [contestId, fetchProblems]);

    useEffect(() => {
        const load = setTimeout(() => {
            void fetchDetails();
        }, 0);
        const timer = setInterval(() => setNow(Date.now()), 30000);
        return () => {
            clearTimeout(load);
            clearInterval(timer);
        };
    }, [fetchDetails]);

    const handleRegistration = async (action: 'register' | 'unregister') => {
        setRegistrationLoading(true);
        try {
            await api.post(`/contests/${contestId}/registration`, { action });
            await fetchDetails();
        } catch (err) {
            console.error('Failed to update registration', err);
        } finally {
            setRegistrationLoading(false);
        }
    };

    if (loading) return <div className="p-12 text-center font-bold">Loading contest details...</div>;
    if (!contest) return <div className="p-12 text-center">Contest not found</div>;

    const isActive = now >= contest.start_time && now <= contest.end_time;
    const registrationOpen = contest.registration_status === 'open' && now >= contest.registration_start_time && now <= contest.registration_end_time;

    return (
        <div className="container mx-auto space-y-12 p-6">
            <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
                <div className="space-y-4">
                    <h1 className="text-4xl font-black">{contest.name}</h1>
                    <p className="text-xl text-muted-foreground max-w-3xl">{contest.description}</p>
                    <div className="flex flex-wrap gap-6 text-sm">
                        <div className="flex items-center gap-2">
                            <Calendar className="w-5 h-5 text-primary" />
                            <span className="font-medium">
                                {new Date(contest.start_time).toLocaleString()} - {new Date(contest.end_time).toLocaleString()}
                            </span>
                        </div>
                        <div className="flex items-center gap-2">
                            <Trophy className="w-5 h-5 text-primary" />
                            <span className="font-medium">{problems.length} Problems</span>
                        </div>
                    </div>
                </div>
                <Card className="w-full lg:w-80">
                    <CardHeader>
                        <CardTitle className="font-black">Registration</CardTitle>
                        <CardDescription>{contest.registered ? 'You are registered for this contest.' : 'Register during the registration window.'}</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="text-sm text-muted-foreground">
                            Window: {new Date(contest.registration_start_time).toLocaleString()} - {new Date(contest.registration_end_time).toLocaleString()}
                        </div>
                        {contest.registered ? (
                                <Button variant="outline" className="w-full gap-2 font-bold" disabled={registrationLoading || isActive} onClick={() => handleRegistration('unregister')}>
                                <UserMinus className="w-4 h-4" /> Unregister
                            </Button>
                        ) : (
                                <Button className="w-full gap-2 font-bold" disabled={registrationLoading || !registrationOpen} onClick={() => handleRegistration('register')}>
                                <UserPlus className="w-4 h-4" /> Register
                            </Button>
                        )}
                        <Link href={`/contests/${contestId}/leaderboard`} className="w-full inline-block">
                            <Button variant="secondary" className="w-full font-bold">View Leaderboard</Button>
                        </Link>
                    </CardContent>
                </Card>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 space-y-6">
                    <h2 className="text-2xl font-bold flex items-center gap-2">
                        <ListChecks className="w-6 h-6 text-primary" />
                        Challenge Problems
                    </h2>
                    {!contest.registered && (
                        <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground">
                            Register for the contest to access the problem list during the contest window.
                        </div>
                    )}
                    {contest.registered && !isActive && (
                        <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground">
                            Problems are visible during the contest window.
                        </div>
                    )}
                    {contest.registered && isActive && (
                        <div className="grid gap-4">
                            {problems.map((problem) => (
                                <Card key={problem.id} className="border-2 border-foreground/15 transition-shadow hover:shadow-[6px_6px_0_hsl(216_36%_13%_/_0.12)]">
                                    <CardContent className="p-6 flex justify-between items-center">
                                        <div>
                                            <h3 className="text-lg font-bold">{problem.name}</h3>
                                            <p className="text-sm text-muted-foreground">{problem.score} points</p>
                                        </div>
                                        <Link href={`/contests/${contestId}/problems/${problem.id}`}>
                                            <Button variant="outline" size="sm" className="gap-2 font-bold">
                                                Solve <ArrowRight className="w-4 h-4" />
                                            </Button>
                                        </Link>
                                    </CardContent>
                                </Card>
                            ))}
                            {problems.length === 0 && (
                                <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground">
                                    No problems have been added yet.
                                </div>
                            )}
                        </div>
                    )}
                </div>

                <div className="space-y-6">
                    <Card>
                        <CardHeader>
                            <CardTitle className="font-black">Status</CardTitle>
                            <CardDescription>{contest.finalized ? 'Final leaderboard is frozen.' : isActive ? 'Contest is live.' : 'Contest is not live.'}</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-2 text-sm text-muted-foreground">
                                <div>Registration: {contest.registration_status}</div>
                                <div>Registered: {contest.registered ? 'Yes' : 'No'}</div>
                                <div>Finalized: {contest.finalized ? 'Yes' : 'No'}</div>
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    );
}
