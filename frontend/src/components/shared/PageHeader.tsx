import { Badge } from '@/components/ui/badge';
import type { ReactNode } from 'react';

interface PageHeaderProps {
    title: string;
    description?: string;
    badge?: string;
    children?: ReactNode;
}

export function PageHeader({ title, description, badge, children }: PageHeaderProps) {
    return (
        <div className="flex flex-col gap-3">
            <div className="flex items-baseline gap-4">
                <h2 className="text-5xl font-bold tracking-tight text-foreground">
                    {title}
                </h2>
                {badge && (
                    <Badge variant="secondary" className="text-xs px-3 py-1 font-mono">
                        {badge}
                    </Badge>
                )}
            </div>
            {description && (
                <p className="text-lg text-muted-foreground max-w-2xl leading-relaxed">
                    {description}
                </p>
            )}
            {children}
        </div>
    );
}
