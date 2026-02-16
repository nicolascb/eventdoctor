import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
    Dialog,
} from "@/components/ui/dialog";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { useProducers } from "@/hooks/useProducers";
import type { ProducerEventEntry, ProducerServiceItem } from "@/types";
import { Skeleton } from "@/components/ui/skeleton";
import { AlertCircle, ChevronLeft, ChevronRight, Server, Workflow } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useState } from "react";
import { SearchInput } from "@/components/shared/SearchInput";
import { ProducerDetails } from "./ProducerDetails";
import { EventDetails } from "./EventDetails";

export function ProducersView() {
    const {
        producers,
        loading,
        error,
        pagination,
        page,
        setPage,
        setSearch,
        search,
    } = useProducers();

    // Detail dialog state
    const [selectedProducer, setSelectedProducer] = useState<ProducerServiceItem | null>(null);
    const [selectedEvent, setSelectedEvent] = useState<ProducerEventEntry | null>(null);

    if (error) {
        return (
            <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
            </Alert>
        );
    }

    return (
        <div className="space-y-4">
            <SearchInput
                value={search}
                onChange={setSearch}
                placeholder="Search producers..."
            />
            <div className="rounded-lg border border-border bg-card overflow-hidden">
                <Table>
                    <TableHeader>
                        <TableRow className="hover:bg-transparent">
                            <TableHead className="text-[11px] font-medium">Service</TableHead>
                            <TableHead className="text-[11px] font-medium text-left">Topics</TableHead>
                            <TableHead className="text-[11px] font-medium text-center">Events</TableHead>
                            <TableHead className="text-[11px] font-medium text-right"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    <TableCell><Skeleton className="h-4 w-[200px]" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-[150px]" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-[50px] ml-auto" /></TableCell>
                                </TableRow>
                            ))
                        ) : producers.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={3} className="h-24 text-center">
                                    No producers found.
                                </TableCell>
                            </TableRow>
                        ) : (
                            producers.map((service) => (
                                <TableRow
                                    key={service.service_id}
                                    className="cursor-pointer group"
                                    onClick={() => { setSelectedProducer(service); }}
                                >
                                    <TableCell className="py-3">
                                        <div className="flex flex-col gap-0.5">
                                            <div className="flex items-center gap-2">
                                                <Server className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                                <span className="font-medium text-sm">{service.service}</span>
                                            </div>
                                            <div className="text-xs text-muted-foreground pl-5.5 truncate max-w-[250px]">
                                                {service.repository}
                                            </div>
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-3 text-center">
                                        <div className="flex flex-wrap gap-1 justify-start">
                                            {service.topics.slice(0, 3).map((topic) => (
                                                <Badge
                                                    key={topic.topic_id}
                                                    variant="secondary"
                                                    className="font-normal text-[10px]"
                                                >
                                                    {topic.topic}
                                                </Badge>
                                            ))}
                                            {service.topics.length > 3 && (
                                                <Badge
                                                    variant="outline"
                                                    className="text-[10px] font-normal text-muted-foreground"
                                                >
                                                    +{service.topics.length - 3} topics
                                                </Badge>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-3 text-center">
                                        <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                            {service.topics.reduce((acc, t) => acc + t.event_count, 0)}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="py-3 text-right">
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            className="h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                                        >
                                            <Workflow className="h-3.5 w-3.5" />
                                        </Button>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>

            {/* Pagination Controls */}
            {pagination && pagination.total_pages > 1 && (
                <div className="flex items-center justify-between px-1">
                    <span className="text-xs text-muted-foreground">
                        Page {pagination.page} of {pagination.total_pages} ({pagination.total} services)
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

            {/* Producer Details Dialog */}
            <ProducerDetails
                producer={selectedProducer}
                open={!!selectedProducer}
                onOpenChange={(open) => { if (!open) setSelectedProducer(null); }}
                onEventClick={(event) => setSelectedEvent(event)}
            />

            {/* Event Details Dialog */}
            <Dialog open={!!selectedEvent} onOpenChange={(open) => { if (!open) setSelectedEvent(null); }}>
                {selectedEvent && (
                    <EventDetails eventId={selectedEvent.id} />
                )}
            </Dialog>
        </div>
    );
}