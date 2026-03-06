import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
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
    SheetContent,
    SheetHeader,
    SheetTitle,
    SheetDescription,
} from "@/components/ui/sheet";
import type { EventView } from "@/types";
import { ExternalLink, Zap, Server, Workflow, FileJson, Loader2, AlertCircle, Network, Code } from "lucide-react";
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
            <SheetContent className="sm:max-w-xl w-[90vw] overflow-y-auto">
                <SheetHeader>
                    <SheetTitle>
                        <div className="flex items-center gap-2">
                            <Zap className="h-4 w-4 text-muted-foreground" />
                            <h3 className="text-base font-semibold font-mono">
                                Loading...
                            </h3>
                        </div>
                    </SheetTitle>
                </SheetHeader>
                <div className="flex h-40 items-center justify-center">
                    <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                </div>
            </SheetContent>
        );
    }

    if (error) {
        return (
            <SheetContent className="sm:max-w-xl w-[90vw] overflow-y-auto">
                <SheetHeader>
                    <SheetTitle>
                        <div className="flex items-center gap-2">
                            <Zap className="h-4 w-4 text-muted-foreground" />
                            <h3 className="text-base font-semibold font-mono">
                                Error
                            </h3>
                        </div>
                    </SheetTitle>
                </SheetHeader>
                <Alert variant="destructive" className="mt-4">
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>Error</AlertTitle>
                    <AlertDescription>{error}</AlertDescription>
                </Alert>
            </SheetContent>
        );
    }

    if (!event) {
        return (
            <SheetContent className="sm:max-w-xl w-[90vw] overflow-y-auto">
                <div className="flex h-40 items-center justify-center text-muted-foreground">
                    Event not found
                </div>
            </SheetContent>
        );
    }

    return (
        <SheetContent className="sm:max-w-xl w-[90vw] overflow-y-auto p-0 flex flex-col">
            <div className="p-6 pb-4 border-b border-border bg-muted/20">
                <SheetHeader>
                    <div className="flex items-center gap-3 mb-2">
                        <div className="rounded-md bg-orange-500/10 p-2 text-orange-500">
                            <Zap className="h-5 w-5" />
                        </div>
                        <div className="flex-1 text-left">
                            <div className="flex flex-row items-center gap-2">
                                <SheetTitle className="text-xl font-mono">{event.name}</SheetTitle>
                                {event.version && (
                                    <Badge variant="outline" className="text-[10px] bg-background">
                                        v{event.version}
                                    </Badge>
                                )}
                            </div>
                            {event.description ? (
                                <SheetDescription className="mt-1">{event.description}</SheetDescription>
                            ) : (
                                <SheetDescription className="mt-1">Event Details</SheetDescription>
                            )}
                        </div>
                    </div>
                </SheetHeader>

                {/* Quick Stats Banner */}
                <div className="flex items-center gap-4 mt-6">
                    <div className="flex items-center gap-2 text-sm text-foreground/80 bg-background rounded-md border px-3 py-1.5 flex-1 justify-center">
                        <Network className="h-4 w-4 text-muted-foreground" />
                        <span className="font-mono">{event.topic}</span>
                    </div>
                    {event.schema_url && (
                        <Button variant="outline" size="sm" className="flex-1 h-9 bg-background" asChild>
                            <a href={event.schema_url} target="_blank" rel="noopener noreferrer" className="gap-2">
                                <Code className="h-4 w-4 text-muted-foreground" />
                                <span>Schema</span>
                                <ExternalLink className="h-3 w-3 ml-auto opacity-50" />
                            </a>
                        </Button>
                    )}
                </div>
            </div>

            <div className="p-6 flex-1 bg-background space-y-8">
                {/* Producers */}
                <div className="space-y-4">
                    <div className="flex items-center justify-between">
                        <h4 className="text-sm font-semibold flex items-center gap-2">
                            <Server className="h-4 w-4 text-muted-foreground" />
                            Producers
                        </h4>
                        <Badge variant="secondary">{event.producers.length}</Badge>
                    </div>

                    {event.producers.length > 0 ? (
                        <div className="rounded-lg border border-border/60 overflow-hidden bg-card/50 shadow-sm">
                            <Table>
                                <TableHeader className="bg-muted/30">
                                    <TableRow className="hover:bg-transparent border-b border-border/60">
                                        <TableHead className="h-9 text-xs font-semibold text-muted-foreground uppercase tracking-wider">Service</TableHead>
                                        <TableHead className="h-9 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-right"></TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {event.producers.map((producer) => (
                                        <TableRow key={`${producer.service}-${producer.repository}`} className="hover:bg-muted/30 transition-colors">
                                            <TableCell className="py-3">
                                                <div className="flex flex-col gap-0.5">
                                                    <span className="text-sm font-medium text-foreground/90">{producer.service}</span>
                                                    <span className="text-[11px] text-muted-foreground truncate max-w-[250px]">
                                                        {producer.repository}
                                                    </span>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-3 text-right">
                                                {producer.owner === "true" && (
                                                    <Badge variant="outline" className="text-[10px] h-5 font-normal border-yellow-500/30 text-yellow-600 bg-yellow-500/5">
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

                <Separator className="bg-border/50" />

                {/* Consumers */}
                <div className="space-y-4">
                    <div className="flex items-center justify-between">
                        <h4 className="text-sm font-semibold flex items-center gap-2">
                            <Workflow className="h-4 w-4 text-muted-foreground" />
                            Consumers
                        </h4>
                        <Badge variant="secondary">{event.consumers.length}</Badge>
                    </div>

                    {event.consumers.length > 0 ? (
                        <div className="rounded-lg border border-border/60 overflow-hidden bg-card/50 shadow-sm">
                            <Table>
                                <TableHeader className="bg-muted/30">
                                    <TableRow className="hover:bg-transparent border-b border-border/60">
                                        <TableHead className="h-9 text-xs font-semibold text-muted-foreground uppercase tracking-wider">Service</TableHead>
                                        <TableHead className="h-9 text-xs font-semibold text-muted-foreground uppercase tracking-wider">Consumer Group</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {event.consumers.map((consumer) => (
                                        <TableRow key={`${consumer.service}-${consumer.id}`} className="hover:bg-muted/30 transition-colors">
                                            <TableCell className="py-3">
                                                <div className="flex flex-col gap-0.5">
                                                    <span className="text-sm font-medium text-foreground/90">{consumer.service}</span>
                                                    <span className="text-[11px] text-muted-foreground truncate max-w-[250px]">
                                                        {consumer.repository}
                                                    </span>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-3">
                                                <span className="text-xs font-mono text-muted-foreground bg-muted/30 px-2 py-1 rounded-md border border-border/50">
                                                    {consumer.id}
                                                </span>
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
                        <Separator className="bg-border/50" />
                        <div className="space-y-4 pb-6">
                            <div className="flex items-center justify-between">
                                <h4 className="text-sm font-semibold flex items-center gap-2">
                                    <FileJson className="h-4 w-4 text-muted-foreground" />
                                    Headers
                                </h4>
                            </div>
                            <div className="rounded-lg border border-border/60 overflow-hidden bg-card/50 shadow-sm">
                                <Table>
                                    <TableHeader className="bg-muted/30">
                                        <TableRow className="hover:bg-transparent border-b border-border/60">
                                            <TableHead className="h-9 text-xs font-semibold text-muted-foreground uppercase tracking-wider">Header Name</TableHead>
                                            <TableHead className="h-9 text-xs font-semibold text-muted-foreground uppercase tracking-wider">Description</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {event.headers.map((header) => (
                                            <TableRow key={header.name} className="hover:bg-muted/30 transition-colors">
                                                <TableCell className="py-3 text-xs font-mono font-medium text-foreground/80">{header.name}</TableCell>
                                                <TableCell className="py-3 text-xs text-muted-foreground">{header.description || '-'}</TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </SheetContent>
    );
}
