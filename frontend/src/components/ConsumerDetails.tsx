import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { Consumer, ConsumerEvent } from "@/types";
import { ExternalLink, Layers, ChevronDown, FileJson } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ServiceDetailsDialog } from "./ServiceDetailsDialog";

interface ConsumerDetailsProps {
    consumer?: Consumer | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onEventClick?: (ev: ConsumerEvent) => void;
}

export function ConsumerDetails({ consumer, open, onOpenChange, onEventClick }: ConsumerDetailsProps) {
    const [expandedTopic, setExpandedTopic] = useState<string | null>(null);

    // Early return if dialog is closed or no consumer selected
    if (!open || !consumer) {
        return null;
    }

    const handleTopicClick = (topicName: string) => {
        if (expandedTopic === topicName) {
            setExpandedTopic(null);
        } else {
            setExpandedTopic(topicName);
        }
    };

    return (
        <ServiceDetailsDialog
            open={open}
            onOpenChange={onOpenChange}
            icon={<Layers className="h-4 w-4 text-muted-foreground" />}
            title={consumer.group}
            description={consumer.description}
        >
            {/* Metadata rows */}
            <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-1">
                <span className="text-muted-foreground text-xs font-medium">Service</span>
                <span className="text-xs font-medium">{consumer.service}</span>
                <span className="text-muted-foreground text-xs font-medium">Repository</span>
                {consumer.repository ? (
                    <a
                        href={consumer.repository}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                    >
                        <ExternalLink className="h-3 w-3 flex-shrink-0" />
                        <span className="truncate">{consumer.repository}</span>
                    </a>
                ) : (
                    <span className="text-xs text-muted-foreground italic">Not specified</span>
                )}
            </div>

            {/* Topic Subscriptions */}
            {consumer.topics.length > 0 && (
                <>
                    <Separator />
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <h4 className="text-sm font-medium text-foreground">
                                Topic Subscriptions
                            </h4>
                            <Badge variant="outline" className="font-mono text-xs">
                                {consumer.topics.length}
                            </Badge>
                        </div>

                        <div className="space-y-2">
                            {consumer.topics.map((topic) => {
                                const isExpanded = expandedTopic === topic.name;
                                return (
                                    <div
                                        key={topic.name}
                                        className="rounded-lg border border-border bg-card overflow-hidden"
                                    >
                                        <button
                                            type="button"
                                            className={cn(
                                                "w-full flex items-center justify-between px-4 py-3 hover:bg-accent/30 transition-colors",
                                                isExpanded && "bg-accent/30"
                                            )}
                                            onClick={() => handleTopicClick(topic.name)}
                                        >
                                            <div className="flex items-center gap-3">
                                                <div className="p-2 bg-secondary/50 rounded-md">
                                                    <Layers className="h-4 w-4 text-secondary-foreground" />
                                                </div>
                                                <div className="flex flex-col items-start gap-1">
                                                    <div className="flex items-center gap-2">
                                                        <code className="text-sm font-mono font-medium">{topic.name}</code>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="flex items-center gap-3">
                                                <Badge variant="secondary" className="text-xs font-normal tabular-nums">
                                                    {topic.events.length} {topic.events.length === 1 ? 'event' : 'events'}
                                                </Badge>
                                                <ChevronDown className={cn(
                                                    "h-4 w-4 text-muted-foreground transition-transform duration-200",
                                                    isExpanded && "transform rotate-180"
                                                )} />
                                            </div>
                                        </button>

                                        {isExpanded && (
                                            <div className="border-t border-border bg-accent/5 px-4 pt-4 pb-2 space-y-3 animate-in slide-in-from-top-2 duration-200">
                                                <div className="space-y-2">
                                                    <h5 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Events</h5>
                                                    <div className="grid gap-2 max-h-[300px] overflow-y-auto pr-1 scrollbar-thin scrollbar-thumb-border scrollbar-track-transparent">
                                                        {topic.events.map((ev) => (
                                                            <div
                                                                key={`${ev.name}-${ev.version}`}
                                                                className={cn(
                                                                    "flex items-start justify-between p-3 rounded-md border border-border bg-background transition-colors group",
                                                                    onEventClick ? "hover:bg-accent/40 cursor-pointer" : ""
                                                                )}
                                                                onClick={() => onEventClick?.(ev)}
                                                            >
                                                                <div className="space-y-1">
                                                                    <div className="flex items-center gap-2">
                                                                        <span className="font-mono text-sm font-medium">{ev.name}</span>
                                                                        {ev.version && (
                                                                            <Badge variant="outline" className="text-[10px] h-5">v{ev.version}</Badge>
                                                                        )}
                                                                    </div>
                                                                </div>
                                                                <div className="flex items-center gap-2">
                                                                    <Button variant="ghost" size="icon" className="h-7 w-7">
                                                                        <FileJson className="h-4 w-4" />
                                                                    </Button>
                                                                </div>
                                                            </div>
                                                        ))}
                                                    </div>
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                </>
            )}
        </ServiceDetailsDialog>
    );
}
