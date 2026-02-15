import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { ProducerServiceItem } from "@/types";
import { ExternalLink, Layers, Server } from "lucide-react";

interface ProducerDetailsProps {
    producer: ProducerServiceItem;
    onTopicClick?: (serviceId: number, topicId: number) => void;
}

export function ProducerDetails({ producer, onTopicClick }: ProducerDetailsProps) {
    return (
        <div className="space-y-4">
            <div className="flex items-center gap-2">
                <Server className="h-4 w-4 text-muted-foreground" />
                <h3 className="text-base font-semibold">
                    {producer.service}
                </h3>
            </div>

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
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <h4 className="text-xs font-medium text-muted-foreground">
                                Published Topics
                            </h4>
                            <span className="text-xs font-mono text-muted-foreground">
                                {producer.topics.length}
                            </span>
                        </div>

                        <div className="space-y-1.5">
                            {producer.topics.map((topic) => {
                                const Content = (
                                    <div className="flex items-center justify-between gap-3">
                                        <div className="flex items-center gap-2">
                                            <code className="text-xs font-mono font-medium">{topic.topic}</code>
                                            {topic.owner && <span className="text-yellow-600 dark:text-yellow-400 text-[10px]">★</span>}
                                        </div>
                                        <div className="flex items-center gap-2 flex-shrink-0">
                                            <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                                {topic.event_count} {topic.event_count === 1 ? 'event' : 'events'}
                                            </Badge>
                                            <Layers className="h-3.5 w-3.5 text-muted-foreground" />
                                        </div>
                                    </div>
                                );

                                if (onTopicClick) {
                                    return (
                                        <button
                                            key={topic.topic_id}
                                            type="button"
                                            className="w-full text-left rounded-md border border-border px-3 py-2.5 hover:bg-accent/30 transition-colors cursor-pointer"
                                            onClick={() => onTopicClick(producer.service_id, topic.topic_id)}
                                        >
                                            {Content}
                                        </button>
                                    );
                                }

                                return (
                                    <div
                                        key={topic.topic_id}
                                        className="rounded-md border border-border px-3 py-2.5 bg-card"
                                    >
                                        {Content}
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                </>
            )}
        </div>
    );
}
