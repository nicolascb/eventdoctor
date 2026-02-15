import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { EventView } from "@/types";
import { ExternalLink, Layers, Server, Zap } from "lucide-react";

interface EventDetailsProps {
    event: EventView;
}

export function EventDetails({ event }: EventDetailsProps) {
    return (
        <div className="space-y-4">
            {/* Header / Basic Info */}
            <div className="space-y-2">
                <div className="flex items-center gap-2">
                    <Zap className="h-4 w-4 text-muted-foreground" />
                    <h3 className="text-base font-semibold font-mono">
                        {event.name}
                    </h3>
                    {event.version && (
                        <Badge variant="outline" className="text-[10px] font-normal font-mono">
                            v{event.version}
                        </Badge>
                    )}
                </div>
                {event.description && (
                    <p className="text-sm text-foreground/80">
                        {event.description}
                    </p>
                )}
            </div>

            {/* Metadata rows */}
            <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-1">
                <span className="text-muted-foreground text-xs font-medium">Topic</span>
                <span className="text-xs font-medium font-mono">{event.topic}</span>

                <span className="text-muted-foreground text-xs font-medium">Schema</span>
                {event.schema_url ? (
                    <a
                        href={event.schema_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                    >
                        <ExternalLink className="h-3 w-3 flex-shrink-0" />
                        <span className="truncate">{event.schema_url}</span>
                    </a>
                ) : (
                    <span className="text-xs text-muted-foreground italic">Not specified</span>
                )}
            </div>

            {/* Headers */}
            {event.headers && event.headers.length > 0 && (
                <>
                    <Separator />
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <h4 className="text-xs font-medium text-muted-foreground">
                                Headers
                            </h4>
                            <span className="text-xs font-mono text-muted-foreground">
                                {event.headers.length}
                            </span>
                        </div>
                        <div className="space-y-1.5">
                            {event.headers.map((header) => (
                                <div
                                    key={header.name}
                                    className="flex items-center justify-between rounded-md border border-border px-3 py-2.5"
                                >
                                    <code className="text-xs font-mono font-medium">{header.name}</code>
                                    {header.description && (
                                        <span className="text-xs text-muted-foreground ml-3 truncate text-right">
                                            {header.description}
                                        </span>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                </>
            )}

            {/* Producers */}
            {event.producers && event.producers.length > 0 && (
                <>
                    <Separator />
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <h4 className="text-xs font-medium text-muted-foreground">
                                Producers
                            </h4>
                            <span className="text-xs font-mono text-muted-foreground">
                                {event.producers.length}
                            </span>
                        </div>
                        <div className="space-y-1.5">
                            {event.producers.map((producer) => (
                                <div
                                    key={producer.id}
                                    className="rounded-md border border-border px-3 py-2.5"
                                >
                                    <div className="flex items-center justify-between gap-3">
                                        <div className="flex items-center gap-2">
                                            <Server className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                            <span className="text-xs font-medium">{producer.service}</span>
                                        </div>
                                        {producer.owner && (
                                            <Badge variant="secondary" className="text-[10px] font-normal">owner</Badge>
                                        )}
                                    </div>
                                    {producer.repository && (
                                        <a
                                            href={producer.repository}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors mt-1 truncate"
                                        >
                                            <ExternalLink className="h-3 w-3 flex-shrink-0" />
                                            <span className="truncate">{producer.repository}</span>
                                        </a>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                </>
            )}

            {/* Consumers */}
            {event.consumers && event.consumers.length > 0 && (
                <>
                    <Separator />
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <h4 className="text-xs font-medium text-muted-foreground">
                                Consumers
                            </h4>
                            <span className="text-xs font-mono text-muted-foreground">
                                {event.consumers.length}
                            </span>
                        </div>
                        <div className="space-y-1.5">
                            {event.consumers.map((consumer) => (
                                <div
                                    key={consumer.id}
                                    className="rounded-md border border-border px-3 py-2.5"
                                >
                                    <div className="flex items-center gap-2">
                                        <Layers className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                        <span className="text-xs font-medium">{consumer.service}</span>
                                    </div>
                                    {consumer.repository && (
                                        <a
                                            href={consumer.repository}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors mt-1 truncate"
                                        >
                                            <ExternalLink className="h-3 w-3 flex-shrink-0" />
                                            <span className="truncate">{consumer.repository}</span>
                                        </a>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                </>
            )}
        </div>
    );
}
