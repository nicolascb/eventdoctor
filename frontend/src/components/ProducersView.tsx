import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Dialog } from "@/components/ui/dialog";
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
import { AlertCircle, Server, Layers, Activity, Github, ChevronRight } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useState } from "react";
import { EmptyState, Pagination, SearchInput, StatCard } from "@/components/shared";
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
        totalTopics,
        totalEvents,
    } = useProducers();

    // Detail dialog state
    const [selectedProducer, setSelectedProducer] = useState<ProducerServiceItem | null>(null);
    const [selectedEvent, setSelectedEvent] = useState<ProducerEventEntry | null>(null);

    const totalProducers = pagination?.total ?? producers.length;

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
        <div className="space-y-6 animate-in">
            {/* Stats Cards */}
            {!loading && !search && (
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <StatCard
                        label="Producer Services"
                        value={totalProducers}
                        description="Active producer services"
                        icon={<Server className="h-4 w-4 text-muted-foreground" />}
                    />
                    <StatCard
                        label="Published Topics"
                        value={totalTopics}
                        description="Total topics published"
                        icon={<Layers className="h-4 w-4 text-muted-foreground" />}
                    />
                    <StatCard
                        label="Event Types"
                        value={totalEvents}
                        description="Registered event structures"
                        icon={<Activity className="h-4 w-4 text-muted-foreground" />}
                    />
                </div>
            )}

            <SearchInput
                value={search}
                onChange={setSearch}
                placeholder="Search producers..."
                resultCount={search && !loading ? (pagination?.total ?? producers.length) : undefined}
            />
            <div className="rounded-xl border border-border/60 bg-card/50 shadow-sm overflow-hidden">
                <Table>
                    <TableHeader className="bg-muted/30">
                        <TableRow className="hover:bg-transparent border-b border-border/60">
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider pl-4">Service</TableHead>
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-left">Topics</TableHead>
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center">Events</TableHead>
                            <TableHead className="h-11 text-right pr-4"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    <TableCell className="pl-4"><Skeleton className="h-12 w-[250px] rounded-md" /></TableCell>
                                    <TableCell><Skeleton className="h-6 w-[200px] rounded-full" /></TableCell>
                                    <TableCell className="text-center"><Skeleton className="h-6 w-[60px] mx-auto rounded-md" /></TableCell>
                                    <TableCell className="pr-4"><Skeleton className="h-8 w-8 ml-auto rounded-full" /></TableCell>
                                </TableRow>
                            ))
                        ) : producers.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={4} className="h-32 text-center p-0">
                                    <div className="py-8">
                                        <EmptyState
                                            title="No producers found"
                                            description="Try adjusting your search criteria"
                                        />
                                    </div>
                                </TableCell>
                            </TableRow>
                        ) : (
                            producers.map((service) => (
                                <TableRow
                                    key={service.service_id}
                                    className="cursor-pointer group hover:bg-muted/30 transition-colors"
                                    onClick={() => { setSelectedProducer(service); }}
                                >
                                    <TableCell className="py-4 pl-4">
                                        <div className="flex items-center gap-4">
                                            <div className="h-10 w-10 shrink-0 rounded-lg bg-green-500/10 flex items-center justify-center border border-green-500/20 group-hover:bg-green-500/20 group-hover:border-green-500/30 transition-colors">
                                                <Server className="h-5 w-5 text-green-500/80" />
                                            </div>
                                            <div className="flex flex-col gap-1.5">
                                                <span className="font-semibold text-foreground/90 group-hover:text-primary transition-colors">{service.service}</span>
                                                <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                                    <Github className="h-3.5 w-3.5 opacity-70" />
                                                    <a
                                                        href={service.repository}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="hover:text-foreground hover:underline truncate max-w-[250px] transition-colors"
                                                        onClick={(e) => e.stopPropagation()}
                                                    >
                                                        {service.repository.replace('https://github.com/', '')}
                                                    </a>
                                                </div>
                                            </div>
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-4 text-left">
                                        <div className="flex flex-wrap gap-1.5 justify-start max-w-[400px]">
                                            {service.topics.slice(0, 3).map((topic) => (
                                                <Badge
                                                    key={topic.topic_id}
                                                    variant="secondary"
                                                    className="font-medium text-[11px] bg-secondary/50 hover:bg-secondary transition-colors"
                                                >
                                                    {topic.topic}
                                                </Badge>
                                            ))}
                                            {service.topics.length > 3 && (
                                                <Badge
                                                    variant="outline"
                                                    className="text-[11px] font-medium text-muted-foreground border-border/50 bg-background/50"
                                                >
                                                    +{service.topics.length - 3} topics
                                                </Badge>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-4 text-center">
                                        <div className="inline-flex items-center justify-center px-2.5 py-1 rounded-md bg-muted/40 border border-border/50 group-hover:bg-muted/60 transition-colors">
                                            <Activity className="h-3.5 w-3.5 text-muted-foreground mr-1.5" />
                                            <span className="text-xs font-semibold tabular-nums text-foreground/80">
                                                {service.topics.reduce((acc, t) => acc + t.event_count, 0)}
                                            </span>
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-4 text-right pr-4">
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            className="h-8 w-8 opacity-0 group-hover:opacity-100 transition-all duration-200 text-muted-foreground hover:text-foreground hover:bg-muted"
                                        >
                                            <ChevronRight className="h-4 w-4" />
                                        </Button>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>

            {(pagination || (!loading && producers.length > 0)) && (
                <Pagination
                    pagination={pagination!}
                    page={page}
                    onPageChange={setPage}
                    label="services"
                />
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
                    <EventDetails eventId={selectedEvent.id} mode="dialog" />
                )}
            </Dialog>
        </div>
    );
}