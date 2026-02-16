import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
} from "@/components/ui/dialog";
import type { EventView } from "@/types";
import { ExternalLink, Zap, Server, Workflow, FileJson, Loader2, AlertCircle } from "lucide-react";
import { useState, useEffect } from "react";
import { api } from "@/lib/api";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";

interface EventDetailsProps {
    eventId: number;
}

export function EventDetails({ eventId }: EventDetailsProps) {
    const [event, setEvent] = useState<EventView | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        let mounted = true;

        async function fetchEvent() {
            setLoading(true);
            setError(null);
            try {
                const data = await api.getEvent(eventId);
                if (mounted) {
                    setEvent(data);
                }
            } catch (err) {
                if (mounted) {
                    setError(err instanceof Error ? err.message : "Failed to load event details");
                }
            } finally {
                if (mounted) {
                    setLoading(false);
                }
            }
        }

        if (eventId) {
            fetchEvent();
        }

        return () => {
            mounted = false;
        };
    }, [eventId]);

    if (loading) {
        return (
            <DialogContent className="max-w-xl">
                <DialogHeader>
                    <DialogTitle>
                        <div className="flex items-center gap-2">
                            <Zap className="h-4 w-4 text-muted-foreground" />
                            <h3 className="text-base font-semibold font-mono">
                                Loading...
                            </h3>
                        </div>
                    </DialogTitle>
                </DialogHeader>
                <div className="flex h-40 items-center justify-center">
                    <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                </div>
            </DialogContent>
        );
    }

    if (error) {
        return (
            <DialogContent className="max-w-xl">
                <DialogHeader>
                    <DialogTitle>
                        <div className="flex items-center gap-2">
                            <Zap className="h-4 w-4 text-muted-foreground" />
                            <h3 className="text-base font-semibold font-mono">
                                Error
                            </h3>
                        </div>
                    </DialogTitle>
                </DialogHeader>
                <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>Error</AlertTitle>
                    <AlertDescription>{error}</AlertDescription>
                </Alert>
            </DialogContent>
        );
    }

    if (!event) {
        return (
            <DialogContent className="max-w-xl">
                <div className="flex h-40 items-center justify-center text-muted-foreground">
                    Event not found
                </div>
            </DialogContent>
        );
    }

    return (
        <DialogContent className="max-w-xl">
            <DialogHeader>
                <div className="flex items-center justify-between pr-8">
                    <DialogTitle>
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
                    </DialogTitle>
                </div>
                {event.description && (
                    <DialogDescription>
                        {event.description}
                    </DialogDescription>
                )}
            </DialogHeader>

            <div className="space-y-4">
                {/* Metadata Grid */}
                <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-2">
                    <span className="text-muted-foreground text-xs font-medium">Topic</span>
                    <span className="text-xs font-mono">{event.topic}</span>

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

                <Separator />

                {/* Producers */}
                <div className="space-y-3">
                    <div className="flex items-center justify-between">
                        <h4 className="text-sm font-medium flex items-center gap-2">
                            <Server className="h-3.5 w-3.5 text-muted-foreground" />
                            Producers
                        </h4>
                        <Badge variant="secondary" className="font-mono text-xs font-normal">
                            {event.producers.length}
                        </Badge>
                    </div>

                    {event.producers.length > 0 ? (
                        <div className="rounded-md border border-border overflow-hidden">
                            <Table>
                                <TableHeader>
                                    <TableRow className="bg-muted/50 hover:bg-muted/50">
                                        <TableHead className="h-8 text-[11px] font-medium">Service</TableHead>
                                        <TableHead className="h-8 text-[11px] font-medium text-right"></TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {event.producers.map((producer) => (
                                        <TableRow key={`${producer.service}-${producer.repository}`} className="hover:bg-transparent">
                                            <TableCell className="py-2">
                                                <div className="flex flex-col gap-0.5">
                                                    <span className="text-xs font-medium">{producer.service}</span>
                                                    <span className="text-[10px] text-muted-foreground truncate max-w-[200px]">
                                                        {producer.repository}
                                                    </span>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-2 text-right">
                                                {producer.owner === "true" && (
                                                    <Badge variant="outline" className="text-[10px] h-5 font-normal border-yellow-500/50 text-yellow-600 dark:text-yellow-400">
                                                        Owner
                                                    </Badge>
                                                )}
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>
                    ) : (
                        <div className="text-sm text-muted-foreground italic px-2">
                            No producers found.
                        </div>
                    )}
                </div>

                <Separator />

                {/* Consumers */}
                <div className="space-y-3">
                    <div className="flex items-center justify-between">
                        <h4 className="text-sm font-medium flex items-center gap-2">
                            <Workflow className="h-3.5 w-3.5 text-muted-foreground" />
                            Consumers
                        </h4>
                        <Badge variant="secondary" className="font-mono text-xs font-normal">
                            {event.consumers.length}
                        </Badge>
                    </div>

                    {event.consumers.length > 0 ? (
                        <div className="rounded-md border border-border overflow-hidden">
                            <Table>
                                <TableHeader>
                                    <TableRow className="bg-muted/50 hover:bg-muted/50">
                                        <TableHead className="h-8 text-[11px] font-medium">Service</TableHead>
                                        <TableHead className="h-8 text-[11px] font-medium">Consumer Group</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {event.consumers.map((consumer) => (
                                        <TableRow key={`${consumer.service}-${consumer.id}`} className="hover:bg-transparent">
                                            <TableCell className="py-2">
                                                <div className="flex flex-col gap-0.5">
                                                    <span className="text-xs font-medium">{consumer.service}</span>
                                                    <span className="text-[10px] text-muted-foreground truncate max-w-[200px]">
                                                        {consumer.repository}
                                                    </span>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-2 text-xs font-mono text-muted-foreground">
                                                {consumer.id}
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>
                    ) : (
                        <div className="text-sm text-muted-foreground italic px-2">
                            No consumers found.
                        </div>
                    )}
                </div>

                {/* Headers */}
                {event.headers && event.headers.length > 0 && (
                    <>
                        <Separator />
                        <div className="space-y-3">
                            <div className="flex items-center justify-between">
                                <h4 className="text-sm font-medium flex items-center gap-2">
                                    <FileJson className="h-3.5 w-3.5 text-muted-foreground" />
                                    Headers
                                </h4>
                            </div>
                            <div className="rounded-md border border-border overflow-hidden">
                                <Table>
                                    <TableHeader>
                                        <TableRow className="bg-muted/50 hover:bg-muted/50">
                                            <TableHead className="h-8 text-[11px] font-medium">Header Name</TableHead>
                                            <TableHead className="h-8 text-[11px] font-medium">Description</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {event.headers.map((header) => (
                                            <TableRow key={header.name} className="hover:bg-transparent">
                                                <TableCell className="py-2 text-xs font-mono">{header.name}</TableCell>
                                                <TableCell className="py-2 text-xs text-muted-foreground">{header.description || '-'}</TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </DialogContent>
    );
}
