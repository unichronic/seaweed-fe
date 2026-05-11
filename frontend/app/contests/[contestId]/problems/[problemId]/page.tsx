'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { useParams } from 'next/navigation';
import api from '@/api';
import { Editor } from '@monaco-editor/react';
import { Button } from '@/components/ui/button';
import { MarkdownRenderer } from '@/components/MarkdownRenderer';
import { Loader2, CheckCircle2, XCircle, Clock } from 'lucide-react';

interface Problem {
    id: string;
    name: string;
    description: string;
    score: number;
}

interface Submission {
    id: string;
    language: string;
    status: string;
    created_at: number;
}

export default function ProblemPage() {
    const { contestId, problemId } = useParams<{ contestId: string; problemId: string }>();
    const [problem, setProblem] = useState<Problem | null>(null);
    const [submissions, setSubmissions] = useState<Submission[]>([]);
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('python');
    const [submitting, setSubmitting] = useState(false);
    const [status, setStatus] = useState<string | null>(null);
    const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);

    const fetchHistory = useCallback(async () => {
        const res = await api.get('/submission/list', { params: { problem_id: problemId } });
        setSubmissions(res.data);
    }, [problemId]);

    useEffect(() => {
        const fetchProblem = async () => {
            try {
                const res = await api.get(`/contests/${contestId}/problems/${problemId}`);
                setProblem(res.data);
                await fetchHistory();
            } catch (err) {
                console.error('Failed to fetch problem', err);
            }
        };
        fetchProblem();
        return () => {
            if (pollRef.current) {
                clearInterval(pollRef.current);
            }
        };
    }, [contestId, problemId, fetchHistory]);

    const handleSubmit = async () => {
        setSubmitting(true);
        setStatus('pending');
        try {
            const res = await api.post('/submission/submit', {
                contest_id: contestId,
                problem_id: problemId,
                language,
                code,
            });
            pollStatus(res.data.submission_id);
            await fetchHistory();
        } catch (err) {
            console.error('Failed to submit', err);
            setStatus('failed_to_process');
            setSubmitting(false);
        }
    };

    const pollStatus = (submissionId: string) => {
        if (pollRef.current) {
            clearInterval(pollRef.current);
        }
        pollRef.current = setInterval(async () => {
            try {
                const res = await api.get(`/submission/${submissionId}/status`);
                const nextStatus = res.data.status;
                setStatus(nextStatus);
                if (nextStatus !== 'pending') {
                    setSubmitting(false);
                    if (pollRef.current) {
                        clearInterval(pollRef.current);
                        pollRef.current = null;
                    }
                    await fetchHistory();
                }
            } catch (err) {
                console.error('Failed to poll submission status', err);
                setSubmitting(false);
                if (pollRef.current) {
                    clearInterval(pollRef.current);
                    pollRef.current = null;
                }
            }
        }, 2000);
    };

    if (!problem) return <div className="p-12 text-center">Loading problem...</div>;

    return (
        <div className="flex flex-col h-[calc(100vh-64px)]">
            <div className="grid grid-cols-1 lg:grid-cols-2 flex-1 overflow-hidden">
                <div className="p-6 overflow-y-auto border-r bg-background">
                    <div className="max-w-3xl mx-auto space-y-8">
                        <div className="flex justify-between items-start gap-4">
                            <div>
                                <h1 className="text-3xl font-bold">{problem.name}</h1>
                                <p className="text-sm text-muted-foreground mt-1">Score: {problem.score} points</p>
                            </div>
                            {status && (
                                <div className={`flex items-center gap-2 px-3 py-1 rounded-full text-sm font-medium ${status === 'accepted' ? 'bg-green-100 text-green-700' :
                                    status === 'pending' ? 'bg-yellow-100 text-yellow-700' : 'bg-red-100 text-red-700'
                                    }`}>
                                    {status === 'pending' ? <Loader2 className="w-4 h-4 animate-spin" /> :
                                        status === 'accepted' ? <CheckCircle2 className="w-4 h-4" /> : <XCircle className="w-4 h-4" />}
                                    {status.toUpperCase()}
                                </div>
                            )}
                        </div>
                        <MarkdownRenderer content={problem.description} />

                        <div className="space-y-3">
                            <h2 className="text-lg font-bold">Submission History</h2>
                            <div className="rounded-lg border divide-y">
                                {submissions.map((sub) => (
                                    <div key={sub.id} className="flex items-center justify-between gap-4 p-3 text-sm">
                                        <div className="flex items-center gap-2 min-w-0">
                                            <Clock className="w-4 h-4 text-muted-foreground" />
                                            <span className="font-mono truncate">{sub.id}</span>
                                        </div>
                                        <div className="flex items-center gap-4 shrink-0">
                                            <span className="text-muted-foreground">{sub.language}</span>
                                            <span className="font-medium">{sub.status}</span>
                                        </div>
                                    </div>
                                ))}
                                {submissions.length === 0 && (
                                    <div className="p-4 text-sm text-muted-foreground text-center">No submissions yet.</div>
                                )}
                            </div>
                        </div>
                    </div>
                </div>

                <div className="flex flex-col bg-[#1e1e1e]">
                    <div className="p-2 border-b border-white/10 flex justify-between items-center bg-[#252526]">
                        <select
                            className="bg-transparent text-sm text-white border border-white/20 rounded px-2 py-1 outline-none"
                            value={language}
                            onChange={(e) => setLanguage(e.target.value)}
                        >
                            <option value="python">Python 3</option>
                            <option value="cpp">C++ 17</option>
                            <option value="java">Java 11</option>
                        </select>
                        <Button
                            size="sm"
                            onClick={handleSubmit}
                            disabled={submitting || code.trim() === ''}
                            className="px-6"
                        >
                            {submitting ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : null}
                            Submit Code
                        </Button>
                    </div>
                    <div className="flex-1">
                        <Editor
                            height="100%"
                            theme="vs-dark"
                            language={language === 'cpp' ? 'cpp' : language}
                            value={code}
                            onChange={(value) => setCode(value || '')}
                            options={{
                                minimap: { enabled: false },
                                fontSize: 14,
                                fontFamily: 'JetBrains Mono, Menlo, Monaco, Courier New, monospace',
                                scrollBeyondLastLine: false,
                                automaticLayout: true,
                            }}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
