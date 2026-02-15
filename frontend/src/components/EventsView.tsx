import { EmptyState, SearchInput, StatCard } from "@/components/shared";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
} from "@/components/ui/dialog";
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
import { ChevronLeft, ChevronRight, Zap } from "lucide-react";
import { useState } from "react";
import { EventDetails } from "./EventDetails";

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
                                <EventDetails event={selectedEvent} />
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
