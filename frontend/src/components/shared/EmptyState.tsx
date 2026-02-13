import { Card, CardContent } from '@/components/ui/card';
import { Search } from 'lucide-react';
import type { ReactNode } from 'react';

interface EmptyStateProps {
    icon?: ReactNode;
    title: string;
    description?: string;
    action?: ReactNode;
}

export function EmptyState({
    icon,
    title,
    description,
    action,
}: EmptyStateProps) {
    return (
        <Card className="border border-dashed">
            <CardContent className="flex flex-col items-center justify-center p-10 text-muted-foreground">
                <div className="mb-3">
                    {icon ?? <Search className="h-6 w-6" />}
                </div>
                <p className="font-medium text-sm">{title}</p>
                {description && (
                    <p className="text-xs mt-1">{description}</p>
                )}
                {action && <div className="mt-3">{action}</div>}
            </CardContent>
        </Card>
    );
}
