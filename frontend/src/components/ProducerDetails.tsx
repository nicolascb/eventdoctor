import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { ProducerDetailView, ProducerEventEntry, ProducerServiceItem } from "@/types";
import { ExternalLink, Layers, Server, ChevronDown, Loader2, FileJson } from "lucide-react";
import { useState } from "react";
import { api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ServiceDetailsDialog } from "./ServiceDetailsDialog";

interface ProducerDetailsProps {
    producer?: ProducerServiceItem | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onEventClick: (event: ProducerEventEntry) => void;
}

export function ProducerDetails({ producer, open, onOpenChange, onEventClick }: ProducerDetailsProps) {
    const [expandedTopicId, setExpandedTopicId] = useState<number | null>(null);
    const [topicDetail, setTopicDetail] = useState<ProducerDetailView | null>(null);
    const [loadingDetail, setLoadingDetail] = useState(false);

    // Early return if dialog is closed or no producer selected
    if (!open || !producer) {
        return null;
    }

    const handleTopicClick = async (topicId: number) => {
        if (expandedTopicId === topicId) {
            setExpandedTopicId(null);
            return;
        }

        setExpandedTopicId(topicId);
        setLoadingDetail(true);
        setTopicDetail(null);

        try {
            const detail = await api.getProducerDetail(producer.service_id, topicId);
            setTopicDetail(detail);
        } catch (error) {
            console.error("Failed to fetch topic details:", error);
        } finally {
            setLoadingDetail(false);
        }
    };

    return (
        <ServiceDetailsDialog
            open={open}
            onOpenChange={onOpenChange}
            icon={<Server className="h-4 w-4 text-muted-foreground" />}
            title={producer.service}
            description={producer.repository}
        >
            {/* Metadata rows */}
            <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-1">
                <span className="text-muted-foreground text-xs font-medium">Repository</span>
                {producer.repository ? (
                    <a
                        href={producer.repository}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                    >
                        <ExternalLink className="h-3 w-3 flex-shrink-0" />
                        <span className="truncate">{producer.repository}</span>
                    </a>
                ) : (
                    <span className="text-xs text-muted-foreground italic">Not specified</span>
                )}
            </div>

            {/* Topics List */}
            {producer.topics.length > 0 && (
                <>
                    <Separator />
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <h4 className="text-sm font-medium text-foreground">
                                Published Topics
                            </h4>
                            <Badge variant="outline" className="font-mono text-xs">
                                {producer.topics.length}
                            </Badge>
                        </div>

                        <div className="space-y-2">
                            {producer.topics.map((topic) => {
                                const isExpanded = expandedTopicId === topic.topic_id;
                                return (
                                    <div
                                        key={topic.topic_id}
                                        className="rounded-lg border border-border bg-card overflow-hidden"
                                    >
                                        <button
                                            type="button"
                                            className={cn(
                                                "w-full flex items-center justify-between px-4 py-3 hover:bg-accent/30 transition-colors",
                                                isExpanded && "bg-accent/30"
                                            )}
                                            onClick={() => handleTopicClick(topic.topic_id)}
                                        >
                                            <div className="flex items-center gap-3">
                                                <div className="p-2 bg-secondary/50 rounded-md">
                                                    <Layers className="h-4 w-4 text-secondary-foreground" />
                                                </div>
                                                <div className="flex flex-col items-start gap-1">
                                                    <div className="flex items-center gap-2">
                                                        <code className="text-sm font-mono font-medium">{topic.topic}</code>
                                                        {topic.owner && (
                                                            <span className="text-yellow-600 dark:text-yellow-400 text-[10px]" title="Owner">★</span>
                                                        )}
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="flex items-center gap-3">
                                                <Badge variant="secondary" className="text-xs font-normal tabular-nums">
                                                    {topic.event_count} {topic.event_count === 1 ? 'event' : 'events'}
                                                </Badge>
                                                <ChevronDown className={cn(
                                                    "h-4 w-4 text-muted-foreground transition-transform duration-200",
                                                    isExpanded && "transform rotate-180"
                                                )} />
                                            </div>
                                        </button>

                                        {isExpanded && (
                                            <div className="border-t border-border bg-accent/5 px-4 pt-4 pb-2 space-y-3 animate-in slide-in-from-top-2 duration-200">
                                                {loadingDetail ? (
                                                    <div className="flex justify-center py-4">
                                                        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                                                    </div>
                                                ) : topicDetail ? (
                                                    <>
                                                        {topicDetail.description && (
                                                            <p className="text-sm text-muted-foreground">
                                                                {topicDetail.description}
                                                            </p>
                                                        )}

                                                        <div className="space-y-2">
                                                            <h5 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Events</h5>
                                                            <div className="grid gap-2 max-h-[300px] overflow-y-auto pr-1 scrollbar-thin scrollbar-thumb-border scrollbar-track-transparent">
                                                                {topicDetail.events.map((ev) => (
                                                                    <div
                                                                        key={`${ev.name}-${ev.version}`}
                                                                        className="flex items-start justify-between p-3 rounded-md border border-border bg-background hover:bg-accent/40 transition-colors cursor-pointer group"
                                                                        onClick={() => onEventClick(ev)}
                                                                    >
                                                                        <div className="space-y-1">
                                                                            <div className="flex items-center gap-2">
                                                                                <span className="font-mono text-sm font-medium">{ev.name}</span>
                                                                                {ev.version && (
                                                                                    <Badge variant="outline" className="text-[10px] h-5">v{ev.version}</Badge>
                                                                                )}
                                                                            </div>
                                                                            {ev.description && (
                                                                                <p className="text-xs text-muted-foreground line-clamp-1">{ev.description}</p>
                                                                            )}
                                                                        </div>
                                                                        <div className="flex items-center gap-2">
                                                                            {ev.schema_url && (
                                                                                <a
                                                                                    href={ev.schema_url}
                                                                                    target="_blank"
                                                                                    rel="noopener noreferrer"
                                                                                    className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded-md transition-colors"
                                                                                    onClick={(e) => e.stopPropagation()}
                                                                                    title="View Schema"
                                                                                >
                                                                                    <ExternalLink className="h-3.5 w-3.5" />
                                                                                </a>
                                                                            )}
                                                                            <Button variant="ghost" size="icon" className="h-7 w-7">
                                                                                <FileJson className="h-4 w-4" />
                                                                            </Button>
                                                                        </div>
                                                                    </div>
                                                                ))}
                                                            </div>
                                                        </div>
                                                    </>
                                                ) : (
                                                    <div className="text-center py-4 text-sm text-muted-foreground">
                                                        Failed to load details.
                                                    </div>
                                                )}
                                            </div>
                                        )
                                        }
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
