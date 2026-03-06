import { EmptyState, Pagination, SearchInput, StatCard } from "@/components/shared";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { useConsumers } from "@/hooks/useConsumers";
import type { Consumer } from "@/types";
import { Activity, Layers, Server, ChevronRight } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { useState } from "react";
import { ConsumerDetails } from "./ConsumerDetails";
import { EventDetails } from "./EventDetails";
import { Sheet } from "@/components/ui/sheet";

export function ConsumersView() {
    const {
        consumers,
        loading,
        pagination,
        page,
        setPage,
        totalTopics,
        totalEvents,
        search,
        setSearch,
    } = useConsumers();

    const [selectedConsumer, setSelectedConsumer] = useState<Consumer | null>(null);
    const [selectedEventId, setSelectedEventId] = useState<number | null>(null);

    const totalConsumerGroups = pagination?.total ?? consumers.length;

    return (
        <div className="space-y-6 animate-in">
            {/* Stats Cards */}
            {!loading && (
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <StatCard
                        label="Consumer Groups"
                        value={totalConsumerGroups}
                        description="Active consumer services"
                        icon={<Layers className="h-4 w-4 text-muted-foreground" />}
                    />
                    <StatCard
                        label="Subscribed Topics"
                        value={totalTopics}
                        description="Total topic subscriptions"
                        icon={<Server className="h-4 w-4 text-muted-foreground" />}
                    />
                    <StatCard
                        label="Event Handlers"
                        value={totalEvents}
                        description="Registered event handlers"
                        icon={<Activity className="h-4 w-4 text-muted-foreground" />}
                    />
                </div>
            )}

            <SearchInput
                value={search}
                onChange={setSearch}
                placeholder="Search consumers, topics, or events..."
                resultCount={search && !loading ? (pagination?.total ?? consumers.length) : undefined}
            />

            <div className="rounded-xl border border-border/60 bg-card/50 shadow-sm overflow-hidden">
                <Table>
                    <TableHeader className="bg-muted/30">
                        <TableRow className="hover:bg-transparent border-b border-border/60">
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider pl-4">Group</TableHead>
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider">Service</TableHead>
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center">Topics</TableHead>
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center">Events</TableHead>
                            <TableHead className="h-11 text-right pr-4"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    <TableCell className="pl-4"><Skeleton className="h-12 w-[250px] rounded-md" /></TableCell>
                                    <TableCell><Skeleton className="h-6 w-[150px] rounded-full" /></TableCell>
                                    <TableCell className="text-center"><Skeleton className="h-6 w-[50px] mx-auto rounded-md" /></TableCell>
                                    <TableCell className="text-center"><Skeleton className="h-6 w-[50px] mx-auto rounded-md" /></TableCell>
                                    <TableCell className="pr-4"><Skeleton className="h-8 w-8 ml-auto rounded-full" /></TableCell>
                                </TableRow>
                            ))
                        ) : consumers.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="h-32 text-center p-0">
                                    <div className="py-8">
                                        <EmptyState
                                            title="No consumers found"
                                            description="Try adjusting your search criteria"
                                        />
                                    </div>
                                </TableCell>
                            </TableRow>
                        ) : (
                            consumers.map((consumer) => {
                                const consumerEventCount = consumer.topics.reduce((sum, t) => sum + t.events.length, 0);
                                return (
                                    <TableRow
                                        key={consumer.group}
                                        className="cursor-pointer group hover:bg-muted/30 transition-colors"
                                        onClick={() => { setSelectedConsumer(consumer); }}
                                    >
                                        <TableCell className="py-4 pl-4">
                                            <div className="flex items-center gap-4">
                                                <div className="h-10 w-10 shrink-0 rounded-lg bg-blue-500/10 flex items-center justify-center border border-blue-500/20 group-hover:bg-blue-500/20 group-hover:border-blue-500/30 transition-colors">
                                                    <Layers className="h-5 w-5 text-blue-500/80" />
                                                </div>
                                                <div className="flex flex-col gap-1">
                                                    <span className="font-semibold text-foreground/90 group-hover:text-primary transition-colors">{consumer.group}</span>
                                                    {consumer.description && (
                                                        <span className="text-xs text-muted-foreground line-clamp-1 max-w-[300px]">
                                                            {consumer.description}
                                                        </span>
                                                    )}
                                                </div>
                                            </div>
                                        </TableCell>
                                        <TableCell className="py-4">
                                            <div className="flex items-center gap-2">
                                                <Server className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                                                <span className="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors">{consumer.service}</span>
                                            </div>
                                        </TableCell>
                                        <TableCell className="py-4 text-center">
                                            <Badge variant="secondary" className="font-medium text-[11px] bg-secondary/50 hover:bg-secondary transition-colors tabular-nums">
                                                {consumer.topics.length} topics
                                            </Badge>
                                        </TableCell>
                                        <TableCell className="py-4 text-center">
                                            <div className="inline-flex items-center justify-center px-2.5 py-1 rounded-md bg-muted/40 border border-border/50 group-hover:bg-muted/60 transition-colors">
                                                <Activity className="h-3.5 w-3.5 text-muted-foreground mr-1.5" />
                                                <span className="text-xs font-semibold tabular-nums text-foreground/80">
                                                    {consumerEventCount}
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
                                );
                            })
                        )}
                    </TableBody>
                </Table>
            </div>

            {(pagination || (!loading && consumers.length > 0)) && (
                <Pagination
                    pagination={pagination!}
                    page={page}
                    onPageChange={setPage}
                    label="groups"
                />
            )}

            {/* Consumer Detail Dialog */}
            <ConsumerDetails
                consumer={selectedConsumer}
                open={!!selectedConsumer}
                onOpenChange={(open) => { if (!open) setSelectedConsumer(null); }}
                onEventClick={(ev) => setSelectedEventId(ev.id)}
            />

            {/* Event Details Dialog */}
            <Sheet open={!!selectedEventId} onOpenChange={(open) => { if (!open) setSelectedEventId(null); }}>
                {selectedEventId && (
                    <EventDetails eventId={selectedEventId} />
                )}
            </Sheet>
        </div>
    );
}
