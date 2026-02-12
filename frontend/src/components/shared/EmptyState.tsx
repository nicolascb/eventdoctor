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
        <Card className="border-2 border-dashed">
            <CardContent className="flex flex-col items-center justify-center p-12 text-muted-foreground">
                <div className="h-16 w-16 rounded-full bg-muted/50 flex items-center justify-center mb-4">
                    {icon ?? <Search className="h-8 w-8" />}
                </div>
                <p className="font-semibold text-lg">{title}</p>
                {description && (
                    <p className="text-sm mt-1">{description}</p>
                )}
                {action && <div className="mt-4">{action}</div>}
            </CardContent>
        </Card>
    );
}
