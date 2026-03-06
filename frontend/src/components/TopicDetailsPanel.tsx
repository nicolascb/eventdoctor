import {
    Sheet,
    SheetContent,
    SheetHeader,
    SheetTitle,
    SheetDescription,
} from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { TopicView } from "@/types";
import { EventDetails } from "./EventDetails";
// Removed Dialog and DialogTrigger import
import {
    Layers,
    Database,
    Network,
    Zap,
    FileJson,
    Code,
    AlertCircle,
    XCircle,
    ArrowRight,
    Server
} from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { useState } from "react"; // Added useState

interface TopicDetailsPanelProps {
    topic: TopicView | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function TopicDetailsPanel({ topic, open, onOpenChange }: TopicDetailsPanelProps) {
    const [selectedEventId, setSelectedEventId] = useState<string | null>(null); // Added state for selected event

    if (!topic) return null;

    const getEventStatus = (eventName: string) => {
        const producers = topic.producers.filter((p) => p.event === eventName);
        const consumers = topic.consumers.filter((c) => c.event === eventName);

        if (producers.length === 0) return "orphaned";
        if (consumers.length === 0) return "unconsumed";
        return "active";
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className="w-full sm:max-w-xl overflow-y-auto p-0 flex flex-col">
                <div className="p-6 pb-4 border-b border-border bg-muted/20">
                    <SheetHeader>
                        <div className="flex items-center gap-3 mb-2">
                            <div className="rounded-md bg-primary/10 p-2 text-primary">
                                <Layers className="h-5 w-5" />
                            </div>
                            <div>
                                <SheetTitle className="text-xl font-mono">{topic.topic}</SheetTitle>
                                <SheetDescription>Topic Details</SheetDescription>
                            </div>
                        </div>
                    </SheetHeader>

                    {/* Quick Stats Banner */}
                    <div className="flex items-center gap-4 mt-6">
                        <div className="flex items-center gap-2 text-sm text-muted-foreground bg-background rounded-md border px-3 py-1.5 flex-1 justify-center">
                            <Database className="h-4 w-4" />
                            <span className="font-medium">{topic.producers.length}</span> Producers
                        </div>
                        <ArrowRight className="h-4 w-4 text-muted-foreground opacity-50 flex-shrink-0" />
                        <div className="flex items-center gap-2 text-sm text-primary font-medium bg-primary/10 rounded-md px-3 py-1.5 flex-1 justify-center">
                            <Zap className="h-4 w-4" />
                            <span>{topic.events.length}</span> Events
                        </div>
                        <ArrowRight className="h-4 w-4 text-muted-foreground opacity-50 flex-shrink-0" />
                        <div className="flex items-center gap-2 text-sm text-muted-foreground bg-background rounded-md border px-3 py-1.5 flex-1 justify-center">
                            <Network className="h-4 w-4" />
                            <span className="font-medium">{topic.consumers.length}</span> Consumers
                        </div>
                    </div>
                </div>

                <div className="p-6 flex-1 bg-background space-y-8">
                    {/* Events List */}
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <h3 className="text-sm font-semibold flex items-center gap-2">
                                <Zap className="h-4 w-4 text-muted-foreground" />
                                Registered Events
                            </h3>
                            <Badge variant="secondary">{topic.events.length}</Badge>
                        </div>

                        <div className="space-y-3">
                            {topic.events.map((event) => {
                                const eventProducers = topic.producers.filter((p) => p.event === event.name);
                                const eventConsumers = topic.consumers.filter((c) => c.event === event.name);
                                const status = getEventStatus(event.name);

                                let statusIcon, statusColor;
                                if (status === "orphaned") {
                                    statusIcon = <XCircle className="h-4 w-4" />;
                                    statusColor = "text-destructive border-destructive/20 bg-destructive/5";
                                } else if (status === "unconsumed") {
                                    statusIcon = <AlertCircle className="h-4 w-4" />;
                                    statusColor = "text-yellow-600 border-yellow-600/20 bg-yellow-600/5";
                                } else {
                                    statusIcon = <Zap className="h-4 w-4" />;
                                    statusColor = "text-primary border-primary/20 bg-primary/5";
                                }

                                return (
                                    <div key={event.name}> {/* Wrapped in a div to contain both the card and the Sheet */}
                                        <div className={`p-4 rounded-lg border flex flex-col gap-4 ${statusColor}`}>
                                            <div className="flex items-start justify-between gap-4">
                                                <div className="space-y-1">
                                                    <div className="flex items-center gap-2">
                                                        {statusIcon}
                                                        <span className="font-mono font-medium">{event.name}</span>
                                                        {event.version && (
                                                            <Badge variant="outline" className="text-[10px] bg-background">
                                                                v{event.version}
                                                            </Badge>
                                                        )}
                                                    </div>
                                                    {event.description && (
                                                        <p className="text-xs text-muted-foreground opacity-90">{event.description}</p>
                                                    )}
                                                </div>
                                                <div className="flex items-center gap-1 shrink-0">
                                                    {event.schema_url && (
                                                        <Button variant="ghost" size="sm" className="h-7 w-7 p-0 bg-background/50 hover:bg-background" asChild>
                                                            <a href={event.schema_url} target="_blank" rel="noopener noreferrer" title="Schema URL">
                                                                <Code className="h-3.5 w-3.5" />
                                                            </a>
                                                        </Button>
                                                    )}
                                                    <Button
                                                        variant="ghost"
                                                        size="sm"
                                                        className="h-7 w-7 p-0 bg-background/50 hover:bg-background"
                                                        title="Event Details"
                                                        onClick={() => setSelectedEventId(event.name)} // Set selected event
                                                    >
                                                        <FileJson className="h-3.5 w-3.5" />
                                                    </Button>
                                                </div>
                                            </div>

                                            <Separator className="bg-border/50" />

                                            {/* Sub-relationship (Who publishes/consumes this specific event) */}
                                            <div className="grid grid-cols-2 gap-4">
                                                <div className="space-y-2">
                                                    <span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">Producers</span>
                                                    {eventProducers.length === 0 ? (
                                                        <p className="text-xs text-destructive italic">No producers found</p>
                                                    ) : (
                                                        <ul className="space-y-1">
                                                            {eventProducers.map((p) => (
                                                                <li key={`${p.service}-${p.repository}`} className="text-xs flex items-center gap-1.5 text-foreground/80">
                                                                    <Server className="h-3 w-3 text-muted-foreground" />
                                                                    <span className="truncate">{p.service}</span>
                                                                </li>
                                                            ))}
                                                        </ul>
                                                    )}
                                                </div>

                                                <div className="space-y-2">
                                                    <span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">Consumers</span>
                                                    {eventConsumers.length === 0 ? (
                                                        <p className="text-xs text-muted-foreground italic">No consumers found</p>
                                                    ) : (
                                                        <ul className="space-y-1">
                                                            {eventConsumers.map((c) => (
                                                                <li key={`${c.service}-${c.group}`} className="text-xs flex flex-col gap-0.5 text-foreground/80">
                                                                    <div className="flex items-center gap-1.5">
                                                                        <Network className="h-3 w-3 text-muted-foreground" />
                                                                        <span className="truncate font-medium">{c.service}</span>
                                                                    </div>
                                                                    <span className="text-[10px] text-muted-foreground pl-4.5 truncate">Group: {c.group}</span>
                                                                </li>
                                                            ))}
                                                        </ul>
                                                    )}
                                                </div>
                                            </div>
                                        </div>
                                        {/* Event Details Sheet */}
                                        <Sheet open={selectedEventId === event.name} onOpenChange={(open) => {
                                            if (!open) setSelectedEventId(null);
                                        }}>
                                            {selectedEventId === event.name && (
                                                <EventDetails eventId={event.id} />
                                            )}
                                        </Sheet>
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                </div>
            </SheetContent>
        </Sheet>
    );
}
