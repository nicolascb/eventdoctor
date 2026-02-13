import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { ReactNode } from 'react';

interface StatCardProps {
    label: string;
    value: number | string;
    description: string;
    icon: ReactNode;
    iconClassName?: string;
    onClick?: () => void;
    active?: boolean;
}

export function StatCard({
    label,
    value,
    description,
    icon,
    onClick,
    active,
}: StatCardProps) {
    return (
        <Card
            className={`${onClick ? 'cursor-pointer hover:bg-accent/50' : ''} transition-colors ${active ? 'ring-2 ring-foreground' : ''}`}
            onClick={onClick}
        >
            <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">{label}</CardTitle>
                {icon}
            </CardHeader>
            <CardContent>
                <div className="text-2xl font-bold stat-number">{value}</div>
                <p className="text-xs text-muted-foreground mt-1">{description}</p>
            </CardContent>
        </Card>
    );
}
