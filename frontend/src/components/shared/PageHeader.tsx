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
        <div className="flex flex-col gap-1">
            <div className="flex items-center gap-3">
                <h2 className="text-lg font-semibold tracking-tight text-foreground">
                    {title}
                </h2>
                {badge && (
                    <Badge variant="secondary" className="text-[11px] px-2 py-0.5 font-mono font-normal">
                        {badge}
                    </Badge>
                )}
            </div>
            {description && (
                <p className="text-sm text-muted-foreground">
                    {description}
                </p>
            )}
            {children}
        </div>
    );
}
