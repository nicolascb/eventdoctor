import { EmptyState, SearchInput, StatCard } from "@/components/shared";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { useEvents } from "@/hooks/useEvents";
import type { EventView } from "@/types";
import { ChevronLeft, ChevronRight, ExternalLink, Layers, Server, Zap } from "lucide-react";
import { useState } from "react";

export function EventsView() {
    const {
        events,
        loading,
        pagination,
        page,
        setPage,
        search,
        setSearch,
    } = useEvents();

    const [selectedEvent, setSelectedEvent] = useState<EventView | null>(null);

    const totalEvents = pagination?.total ?? events.length;

    return (
        <div className="space-y-4 animate-in">
            {/* Stats Cards */}
            {!loading && (
                <div className="grid grid-cols-1 sm:grid-cols-1 gap-4">
                    <StatCard
                        label="Total Events"
                        value={totalEvents}
                        description="Registered event types"
                        icon={<Zap className="h-4 w-4 text-muted-foreground" />}
                    />
                </div>
            )}

            {/* Search Bar */}
            <SearchInput
                value={search}
                onChange={setSearch}
                placeholder="Search events or topics..."
                resultCount={search ? (pagination?.total ?? events.length) : undefined}
                totalCount={search ? undefined : undefined}
            />

            {!loading && (
                <>
                    {/* Events Table */}
                    <div className="rounded-lg border border-border bg-card overflow-hidden">
                        {events.length === 0 ? (
                            <div className="p-6">
                                <EmptyState
                                    title="No events found"
                                    description="No events registered in the system."
                                />
                            </div>
                        ) : (
                            <Table>
                                <TableHeader>
                                    <TableRow className="hover:bg-transparent">
                                        <TableHead className="text-[11px] font-medium">Event</TableHead>
                                        <TableHead className="text-[11px] font-medium">Topic</TableHead>
                                        <TableHead className="text-[11px] font-medium text-center">Consumers</TableHead>
                                        <TableHead className="text-[11px] font-medium text-center">Producers</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {events.map((event) => (
                                        <TableRow
                                            key={event.id}
                                            className="cursor-pointer group"
                                            onClick={() => setSelectedEvent(event)}
                                        >
                                            <TableCell className="py-3">
                                                <div className="flex flex-col gap-0.5">
                                                    <div className="flex items-center gap-2">
                                                        <Zap className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                                        <span className="font-medium text-sm font-mono">{event.name}</span>
                                                        {event.version && (
                                                            <span className="text-xs text-muted-foreground border px-1 rounded">v{event.version}</span>
                                                        )}
                                                    </div>
                                                    {event.description && (
                                                        <span className="text-xs text-muted-foreground pl-5.5 line-clamp-1">
                                                            {event.description}
                                                        </span>
                                                    )}
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-3">
                                                <div className="flex items-center gap-2">
                                                    <NetworkIcon className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                                    <span className="text-sm font-mono">{event.topic}</span>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-3 text-center">
                                                <div className="flex items-center justify-center gap-1">
                                                    <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                                        {event.consumers.length}
                                                    </Badge>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-3 text-center">
                                                <div className="flex items-center justify-center gap-1">
                                                    <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                                        {event.producers.length}
                                                    </Badge>
                                                </div>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        )}
                    </div>

                    {/* Pagination Controls */}
                    {pagination && pagination.total_pages > 1 && (
                        <div className="flex items-center justify-between px-1">
                            <span className="text-xs text-muted-foreground">
                                Page {pagination.page} of {pagination.total_pages} ({pagination.total} events)
                            </span>
                            <div className="flex items-center gap-1">
                                <Button
                                    variant="outline"
                                    size="sm"
                                    className="h-7 px-2 text-xs"
                                    disabled={page <= 1}
                                    onClick={() => setPage(page - 1)}
                                >
                                    <ChevronLeft className="h-3.5 w-3.5 mr-1" />
                                    Previous
                                </Button>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    className="h-7 px-2 text-xs"
                                    disabled={page >= pagination.total_pages}
                                    onClick={() => setPage(page + 1)}
                                >
                                    Next
                                    <ChevronRight className="h-3.5 w-3.5 ml-1" />
                                </Button>
                            </div>
                        </div>
                    )}

                    {/* Event Detail Dialog */}
                    <Dialog open={!!selectedEvent} onOpenChange={(open) => { if (!open) setSelectedEvent(null); }}>
                        {selectedEvent && (
                            <DialogContent className="max-w-xl">
                                <DialogHeader>
                                    <div className="flex items-center gap-2">
                                        <Zap className="h-4 w-4 text-muted-foreground" />
                                        <DialogTitle className="text-base font-semibold font-mono">
                                            {selectedEvent.name}
                                        </DialogTitle>
                                        {selectedEvent.version && (
                                            <Badge variant="outline" className="text-[10px] font-normal font-mono">
                                                v{selectedEvent.version}
                                            </Badge>
                                        )}
                                    </div>
                                    {selectedEvent.description && (
                                        <DialogDescription>
                                            {selectedEvent.description}
                                        </DialogDescription>
                                    )}
                                </DialogHeader>

                                {/* Metadata rows */}
                                <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-1">
                                    <span className="text-muted-foreground text-xs font-medium">Topic</span>
                                    <span className="text-xs font-medium font-mono">{selectedEvent.topic}</span>

                                    <span className="text-muted-foreground text-xs font-medium">Schema</span>
                                    {selectedEvent.schema_url ? (
                                        <a
                                            href={selectedEvent.schema_url}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                                        >
                                            <ExternalLink className="h-3 w-3 flex-shrink-0" />
                                            <span className="truncate">{selectedEvent.schema_url}</span>
                                        </a>
                                    ) : (
                                        <span className="text-xs text-muted-foreground italic">Not specified</span>
                                    )}
                                </div>

                                {/* Headers */}
                                {selectedEvent.headers && selectedEvent.headers.length > 0 && (
                                    <>
                                        <Separator />
                                        <div className="space-y-2">
                                            <div className="flex items-center justify-between">
                                                <h4 className="text-xs font-medium text-muted-foreground">
                                                    Headers
                                                </h4>
                                                <span className="text-xs font-mono text-muted-foreground">
                                                    {selectedEvent.headers.length}
                                                </span>
                                            </div>
                                            <div className="space-y-1.5">
                                                {selectedEvent.headers.map((header) => (
                                                    <div
                                                        key={header.name}
                                                        className="flex items-center justify-between rounded-md border border-border px-3 py-2.5"
                                                    >
                                                        <code className="text-xs font-mono font-medium">{header.name}</code>
                                                        {header.description && (
                                                            <span className="text-xs text-muted-foreground ml-3 truncate">{header.description}</span>
                                                        )}
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    </>
                                )}

                                {/* Producers */}
                                {selectedEvent.producers.length > 0 && (
                                    <>
                                        <Separator />
                                        <div className="space-y-2">
                                            <div className="flex items-center justify-between">
                                                <h4 className="text-xs font-medium text-muted-foreground">
                                                    Producers
                                                </h4>
                                                <span className="text-xs font-mono text-muted-foreground">
                                                    {selectedEvent.producers.length}
                                                </span>
                                            </div>
                                            <div className="space-y-1.5">
                                                {selectedEvent.producers.map((producer) => (
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
                                {selectedEvent.consumers.length > 0 && (
                                    <>
                                        <Separator />
                                        <div className="space-y-2">
                                            <div className="flex items-center justify-between">
                                                <h4 className="text-xs font-medium text-muted-foreground">
                                                    Consumers
                                                </h4>
                                                <span className="text-xs font-mono text-muted-foreground">
                                                    {selectedEvent.consumers.length}
                                                </span>
                                            </div>
                                            <div className="space-y-1.5">
                                                {selectedEvent.consumers.map((consumer) => (
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
                            </DialogContent>
                        )}
                    </Dialog>
                </>
            )}
        </div>
    );
}

function NetworkIcon(props: React.SVGProps<SVGSVGElement>) {
    return (
        <svg
            {...props}
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
        >
            <rect x="16" y="16" width="6" height="6" rx="1" />
            <rect x="2" y="16" width="6" height="6" rx="1" />
            <rect x="9" y="2" width="6" height="6" rx="1" />
            <path d="M5 16v-3a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1v3" />
            <path d="M12 12V8" />
        </svg>
    );
}
