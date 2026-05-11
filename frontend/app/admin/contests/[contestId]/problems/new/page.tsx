'use client';

import { useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import api from '@/api';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

export default function NewProblemPage() {
    const { contestId } = useParams();
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        score: '10',
        test_cases: '[\n  {\n    "id": "sample-1",\n    "stdin": "",\n    "expected_output": ""\n  }\n]',
    });
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            const testCases = JSON.parse(formData.test_cases);
            await api.post(`/admin/${contestId}/problem`, {
                name: formData.name,
                description: formData.description,
                score: parseInt(formData.score),
                test_cases: testCases,
            });
            router.push(`/contests/${contestId}`);
        } catch (err) {
            console.error('Failed to add problem', err);
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setFormData({ ...formData, [e.target.id]: e.target.value });
    };

    return (
        <div className="container mx-auto p-6 max-w-3xl">
            <Card>
                <CardHeader>
                    <CardTitle>Add New Problem</CardTitle>
                    <CardDescription>Adding a problem to contest: {contestId}</CardDescription>
                </CardHeader>
                <form onSubmit={handleSubmit}>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="name">Problem Name</Label>
                            <Input id="name" placeholder="Two Sum" value={formData.name} onChange={handleChange} required />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="score">Base Score</Label>
                            <Input id="score" type="number" value={formData.score} onChange={handleChange} required />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="description">Problem Statement (supports Markdown + LaTeX)</Label>
                            <Textarea
                                id="description"
                                placeholder="Describe the problem, input/output format, and constraints..."
                                className="min-h-[300px] font-mono"
                                value={formData.description}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="test_cases">Judge0 Test Cases (JSON)</Label>
                            <Textarea
                                id="test_cases"
                                className="min-h-[180px] font-mono"
                                value={formData.test_cases}
                                onChange={handleChange}
                                required
                            />
                        </div>
                    </CardContent>
                    <CardFooter>
                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? 'Adding Problem...' : 'Add Problem'}
                        </Button>
                    </CardFooter>
                </form>
            </Card>
        </div>
    );
}
