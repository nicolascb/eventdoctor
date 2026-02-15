import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { Consumer, Topic } from "@/types";
import { ExternalLink, Layers, Zap } from "lucide-react";

interface ConsumerDetailsProps {
    consumer: Consumer;
    onTopicClick?: (topic: Topic) => void;
}

export function ConsumerDetails({ consumer, onTopicClick }: ConsumerDetailsProps) {
    return (
        <div className="space-y-4">
            {/* Header */}
            <div>
                <div className="flex items-center gap-2">
                    <Layers className="h-4 w-4 text-muted-foreground" />
                    <h3 className="text-base font-semibold">
                        {consumer.group}
                    </h3>
                </div>
                {consumer.description && (
                    <p className="text-sm text-foreground/80 mt-1">
                        {consumer.description}
                    </p>
                )}
            </div>

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
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <h4 className="text-xs font-medium text-muted-foreground">
                                Topic Subscriptions
                            </h4>
                            <span className="text-xs font-mono text-muted-foreground">
                                {consumer.topics.length}
                            </span>
                        </div>

                        <div className="space-y-1.5">
                            {consumer.topics.map((topic) => {
                                const Content = (
                                    <>
                                        <div className="flex items-center justify-between gap-3">
                                            <code className="text-xs font-mono font-medium">{topic.name}</code>
                                            <div className="flex items-center gap-2 flex-shrink-0">
                                                <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                                    {topic.events.length} {topic.events.length === 1 ? 'event' : 'events'}
                                                </Badge>
                                                <Zap className="h-3.5 w-3.5 text-muted-foreground" />
                                            </div>
                                        </div>
                                        {topic.events.length > 0 && (
                                            <div className="flex flex-wrap gap-1 mt-1.5">
                                                {topic.events.slice(0, 4).map((ev) => (
                                                    <Badge
                                                        key={ev.name}
                                                        variant="outline"
                                                        className="text-[10px] font-normal font-mono"
                                                    >
                                                        {ev.name}
                                                        {ev.version && <span className="ml-0.5 text-muted-foreground">v{ev.version}</span>}
                                                    </Badge>
                                                ))}
                                                {topic.events.length > 4 && (
                                                    <Badge variant="secondary" className="text-[10px] font-normal">
                                                        +{topic.events.length - 4}
                                                    </Badge>
                                                )}
                                            </div>
                                        )}
                                    </>
                                );

                                if (onTopicClick) {
                                    return (
                                        <button
                                            key={topic.name}
                                            type="button"
                                            className="w-full text-left rounded-md border border-border px-3 py-2.5 hover:bg-accent/30 transition-colors cursor-pointer"
                                            onClick={() => onTopicClick(topic)}
                                        >
                                            {Content}
                                        </button>
                                    );
                                }

                                return (
                                    <div
                                        key={topic.name}
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
