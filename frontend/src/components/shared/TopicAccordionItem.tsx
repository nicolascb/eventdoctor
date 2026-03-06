import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { ChevronDown, Layers } from 'lucide-react';
import type { ReactNode } from 'react';

interface TopicAccordionItemProps {
    topicName: string;
    eventCount: number;
    isExpanded: boolean;
    onToggle: () => void;
    titleSuffix?: ReactNode;
    children: ReactNode;
}

export function TopicAccordionItem({
    topicName,
    eventCount,
    isExpanded,
    onToggle,
    titleSuffix,
    children,
}: TopicAccordionItemProps) {
    return (
        <div className="rounded-lg border border-border bg-card overflow-hidden">
            <button
                type="button"
                className={cn(
                    'w-full flex items-center justify-between px-4 py-3 hover:bg-accent/30 transition-colors',
                    isExpanded && 'bg-accent/30',
                )}
                onClick={onToggle}
            >
                <div className="flex items-center gap-3">
                    <div className="p-2 bg-secondary/50 rounded-md">
                        <Layers className="h-4 w-4 text-secondary-foreground" />
                    </div>
                    <code className="text-sm font-mono font-medium">{topicName}</code>
                    {titleSuffix}
                </div>
                <div className="flex items-center gap-3">
                    <Badge variant="secondary" className="text-xs font-normal tabular-nums">
                        {eventCount} {eventCount === 1 ? 'event' : 'events'}
                    </Badge>
                    <ChevronDown
                        className={cn(
                            'h-4 w-4 text-muted-foreground transition-transform duration-200',
                            isExpanded && 'rotate-180',
                        )}
                    />
                </div>
            </button>

            {isExpanded && (
                <div className="border-t border-border bg-accent/5 px-4 pt-4 pb-2 space-y-3 animate-in slide-in-from-top-2 duration-200">
                    {children}
                </div>
            )}
        </div>
    );
}
