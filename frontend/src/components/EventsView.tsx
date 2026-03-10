import { Pagination, SearchInput, StatCard } from "@/components/shared";
import { Sheet } from "@/components/ui/sheet";
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
import { Network, Zap } from "lucide-react";
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
            {!loading && !search && (
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
                    <div className="rounded-xl border border-border/60 bg-card/50 shadow-sm overflow-hidden">
                        {events.length === 0 ? (
                            <div className="flex flex-col items-center justify-center p-12 text-muted-foreground bg-muted/10">
                                <div className="h-12 w-12 rounded-full bg-muted/50 flex items-center justify-center mb-4 border border-border/50">
                                    <Zap className="h-6 w-6 text-muted-foreground/70" />
                                </div>
                                <p className="font-semibold text-foreground/80">No events found</p>
                                <p className="text-sm mt-1">No events registered in the system.</p>
                            </div>
                        ) : (
                            <Table>
                                <TableHeader className="bg-muted/30">
                                    <TableRow className="hover:bg-transparent border-b border-border/60">
                                        <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider pl-4">Event</TableHead>
                                        <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider">Topic</TableHead>
                                        <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center">Consumers</TableHead>
                                        <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center">Producers</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {events.map((event) => (
                                        <TableRow
                                            key={event.id}
                                            className="cursor-pointer group hover:bg-muted/30 transition-colors"
                                            onClick={() => setSelectedEvent(event)}
                                        >
                                            <TableCell className="py-4 pl-4">
                                                <div className="flex items-center gap-4">
                                                    <div className="h-10 w-10 shrink-0 rounded-lg bg-orange-500/10 flex items-center justify-center border border-orange-500/20 group-hover:bg-orange-500/20 group-hover:border-orange-500/30 transition-colors">
                                                        <Zap className="h-5 w-5 text-orange-500/80" />
                                                    </div>
                                                    <div className="flex flex-col gap-1.5">
                                                        <div className="flex items-center gap-2">
                                                            <span className="font-semibold text-foreground/90 group-hover:text-primary transition-colors">{event.name}</span>
                                                            {event.version && (
                                                                <span className="text-[10px] font-medium text-muted-foreground border border-border/60 bg-muted/30 px-1.5 py-0.5 rounded-md">v{event.version}</span>
                                                            )}
                                                        </div>
                                                        {event.description && (
                                                            <span className="text-xs text-muted-foreground line-clamp-1 max-w-[300px]">
                                                                {event.description}
                                                            </span>
                                                        )}
                                                    </div>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-4">
                                                <div className="flex items-center gap-2">
                                                    <Network className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                                                    <span className="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors">{event.topic}</span>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-4 text-center">
                                                <div className="inline-flex items-center justify-center px-2.5 py-1 rounded-md bg-muted/40 border border-border/50 group-hover:bg-muted/60 transition-colors">
                                                    <span className="text-xs font-semibold tabular-nums text-foreground/80">
                                                        {event.consumers.length}
                                                    </span>
                                                </div>
                                            </TableCell>
                                            <TableCell className="py-4 text-center">
                                                <div className="inline-flex items-center justify-center px-2.5 py-1 rounded-md bg-muted/40 border border-border/50 group-hover:bg-muted/60 transition-colors">
                                                    <span className="text-xs font-semibold tabular-nums text-foreground/80">
                                                        {event.producers.length}
                                                    </span>
                                                </div>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        )}
                    </div>

                    {pagination && (
                        <Pagination
                            pagination={pagination}
                            page={page}
                            onPageChange={setPage}
                            label="events"
                        />
                    )}

                    {/* Event Detail Dialog */}
                    <Sheet open={!!selectedEvent} onOpenChange={(open) => { if (!open) setSelectedEvent(null); }}>
                        {selectedEvent && (
                            <EventDetails eventId={selectedEvent.id} />
                        )}
                    </Sheet>
                </>
            )}
        </div>
    );
}

