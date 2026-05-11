import Image from 'next/image';
import Link from 'next/link';
import { ArrowRight, Code2, Medal, Play, Terminal, Trophy, Users, Zap } from 'lucide-react';
import { Button } from '@/components/ui/button';

export default function LandingPage() {
    return (
        <div className="pb-12">
            <section className="relative overflow-hidden border-b bg-background">
                <div className="mx-auto grid max-w-7xl grid-cols-1 gap-10 px-6 py-12 lg:min-h-[calc(100vh-64px)] lg:grid-cols-[0.95fr_1.05fr] lg:items-center lg:py-16">
                    <div className="max-w-2xl space-y-8">
                        <div className="inline-flex items-center gap-3 rounded-lg border bg-card px-3 py-2 text-sm font-bold">
                            <span className="flex size-7 items-center justify-center rounded-md bg-accent text-accent-foreground">
                                <Zap className="size-4" />
                            </span>
                            Live coding rounds, instant verdicts
                        </div>
                        <div className="space-y-5">
                            <h1 className="max-w-3xl text-4xl font-black leading-[0.98] sm:text-5xl md:text-7xl">
                                Code contests with a workshop edge.
                            </h1>
                            <p className="max-w-xl text-lg leading-8 text-muted-foreground md:text-xl">
                                Seaweed Arena gives teams a focused place to launch rounds, publish problems, judge submissions, and crown the leaderboard.
                            </p>
                        </div>
                        <div className="flex flex-col gap-3 sm:flex-row">
                            <Link href="/contests">
                                <Button size="lg" className="h-11 w-full gap-2 px-5 font-bold sm:w-auto">
                                    Enter Contests <ArrowRight className="size-4" />
                                </Button>
                            </Link>
                            <Link href="/auth/signup">
                                <Button size="lg" variant="outline" className="h-11 w-full gap-2 px-5 font-bold sm:w-auto">
                                    Join Arena
                                </Button>
                            </Link>
                        </div>
                        <Image
                            src="/seaweed-arena-badge.svg"
                            alt="Seaweed Arena badge"
                            width={144}
                            height={144}
                            priority
                            className="mx-auto size-32 drop-shadow-xl md:hidden"
                        />
                    </div>

                    <div className="relative hidden md:block lg:min-h-[520px]">
                        <Image
                            src="/seaweed-arena-badge.svg"
                            alt="Seaweed Arena badge"
                            width={260}
                            height={260}
                            priority
                            className="relative z-10 mx-auto mb-4 size-36 drop-shadow-xl sm:size-44 lg:absolute lg:right-10 lg:top-0 lg:mx-0 lg:mb-0 lg:size-56"
                        />
                        <div className="relative mx-auto max-w-2xl rounded-lg border-4 border-foreground bg-card p-4 shadow-[12px_12px_0_hsl(216_36%_13%)] lg:absolute lg:inset-x-0 lg:bottom-0">
                            <div className="mb-4 flex items-center justify-between border-b border-foreground/20 pb-3">
                                <div className="flex items-center gap-2 font-black">
                                    <Terminal className="size-5 text-primary" />
                                    Round Control
                                </div>
                                <div className="rounded-md bg-accent px-3 py-1 text-xs font-black text-accent-foreground">LIVE</div>
                            </div>
                            <div className="grid gap-4 md:grid-cols-[1fr_0.8fr]">
                                <div className="rounded-lg bg-foreground p-4 text-background">
                                    <div className="mb-3 flex items-center justify-between text-xs font-bold text-background/70">
                                        <span>PROBLEM A</span>
                                        <span>1200 pts</span>
                                    </div>
                                    <pre className="overflow-hidden text-sm leading-7">
                                        <code>{'fn solve(grid) {\n  push(queue, start)\n  while queue.len > 0 {\n    step()\n  }\n}'}</code>
                                    </pre>
                                    <div className="mt-5 grid grid-cols-3 gap-2 text-center text-xs font-bold">
                                        <div className="rounded-md bg-background/10 px-2 py-2">PY</div>
                                        <div className="rounded-md bg-background/10 px-2 py-2">CPP</div>
                                        <div className="rounded-md bg-background/10 px-2 py-2">JAVA</div>
                                    </div>
                                </div>
                                <div className="space-y-3">
                                    <ScoreRow rank="01" name="Riya" score="840" tone="primary" />
                                    <ScoreRow rank="02" name="Dev" score="720" tone="secondary" />
                                    <ScoreRow rank="03" name="Mira" score="690" tone="accent" />
                                    <div className="rounded-lg border bg-background p-4">
                                        <div className="mb-3 flex items-center gap-2 text-sm font-black">
                                            <Play className="size-4 text-primary" />
                                            Judge queue
                                        </div>
                                        <div className="space-y-2 text-xs font-bold text-muted-foreground">
                                            <div className="flex justify-between"><span>Accepted</span><span>128</span></div>
                                            <div className="flex justify-between"><span>Running</span><span>09</span></div>
                                            <div className="flex justify-between"><span>Penalty</span><span>14m</span></div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <section className="mx-auto grid max-w-7xl grid-cols-1 gap-4 px-6 py-12 md:grid-cols-3">
                <FeatureCard
                    icon={<Code2 className="size-6" />}
                    title="Problem Workshop"
                    description="Create statements, scores, test cases, and language rules for each round."
                />
                <FeatureCard
                    icon={<Trophy className="size-6" />}
                    title="Judge Pipeline"
                    description="Submissions move through verdicts and update standings without manual scoring."
                />
                <FeatureCard
                    icon={<Users className="size-6" />}
                    title="Contest Arena"
                    description="Participants register, solve, submit, and watch the board move in real time."
                />
            </section>

            <section className="mx-auto max-w-7xl px-6">
                <div className="grid gap-4 border-t border-foreground/20 pt-8 md:grid-cols-[1fr_auto] md:items-center">
                    <div>
                        <h2 className="text-3xl font-black">Run the next round.</h2>
                        <p className="mt-2 max-w-2xl text-muted-foreground">Open a contest, publish the problem set, and let the scoreboard settle it.</p>
                    </div>
                    <Link href="/admin/contests">
                        <Button variant="secondary" className="h-11 w-full gap-2 px-5 font-bold md:w-auto">
                            Manage Contests <Medal className="size-4" />
                        </Button>
                    </Link>
                </div>
            </section>
        </div>
    );
}

function ScoreRow({ rank, name, score, tone }: { rank: string; name: string; score: string; tone: 'primary' | 'secondary' | 'accent' }) {
    const toneClass = {
        primary: 'bg-primary text-primary-foreground',
        secondary: 'bg-secondary text-secondary-foreground',
        accent: 'bg-accent text-accent-foreground',
    }[tone];

    return (
        <div className="flex items-center justify-between rounded-lg border bg-background p-3">
            <div className="flex items-center gap-3">
                <span className={`flex size-8 items-center justify-center rounded-md text-sm font-black ${toneClass}`}>{rank}</span>
                <span className="font-bold">{name}</span>
            </div>
            <span className="font-mono text-lg font-black">{score}</span>
        </div>
    );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
    return (
        <div className="rounded-lg border bg-card p-5 shadow-[6px_6px_0_hsl(216_36%_13%_/_0.12)]">
            <div className="mb-5 flex size-11 items-center justify-center rounded-lg bg-foreground text-background">
                {icon}
            </div>
            <h3 className="text-xl font-black">{title}</h3>
            <p className="mt-2 text-sm leading-6 text-muted-foreground">{description}</p>
        </div>
    );
}
