import { Card, CardContent } from '@/components/ui/card';
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
    iconClassName = 'bg-primary/10',
    onClick,
    active,
}: StatCardProps) {
    return (
        <Card
            className={`${onClick ? 'cursor-pointer hover:shadow-lg' : ''} transition-all ${active ? 'ring-2 ring-primary' : ''}`}
            onClick={onClick}
        >
            <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                    <div>
                        <p className="text-sm font-medium text-muted-foreground">{label}</p>
                        <p className="text-3xl font-bold mt-1">{value}</p>
                    </div>
                    <div className={`h-12 w-12 rounded-full ${iconClassName} flex items-center justify-center`}>
                        {icon}
                    </div>
                </div>
                <p className="text-xs text-muted-foreground mt-3">{description}</p>
            </CardContent>
        </Card>
    );
}
