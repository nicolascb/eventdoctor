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
import { Skeleton } from "@/components/ui/skeleton";

import { useTopics } from "@/hooks/useTopics";
import type { TopicConsumerEntry, TopicProducerEntry, TopicView } from "@/types";
import { AlertCircle, Box, CheckCircle2, ChevronRight, Filter, Layers, XCircle, Zap } from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import { TopicDetailsPanel } from "./TopicDetailsPanel";

type FilterMode = 'all' | 'orphaned' | 'unconsumed' | 'active';

export function TopicsView() {
    const {
        topics,
        countEvents,
        countOrphaned,
        countUnconsumed,
        loading,
        pagination,
        page,
        setPage
    } = useTopics();
    const [searchQuery, setSearchQuery] = useState("");
    const [filterMode, setFilterMode] = useState<FilterMode>('all');
    const [selectedTopic, setSelectedTopic] = useState<TopicView | null>(null);

    const getEventProducers = useCallback((topicData: TopicView, eventName: string): TopicProducerEntry[] => {
        return topicData.producers.filter(p => p.event === eventName);
    }, []);

    const getEventConsumers = useCallback((topicData: TopicView, eventName: string): TopicConsumerEntry[] => {
        return topicData.consumers.filter(c => c.event === eventName);
    }, []);

    const getEventStatus = useCallback((topicData: TopicView, eventName: string) => {
        const producers = getEventProducers(topicData, eventName);
        const consumers = getEventConsumers(topicData, eventName);

        if (producers.length === 0) return 'orphaned';
        if (consumers.length === 0) return 'unconsumed';
        return 'active';
    }, [getEventProducers, getEventConsumers]);

    const activeCount = countEvents - countOrphaned - countUnconsumed;

    const filteredTopics = useMemo(() => {
        let result = topics;

        // Apply search filter
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            result = result.map(topic => {
                const matchesTopic = topic.topic.toLowerCase().includes(query);
                const matchingEvents = topic.events.filter(event =>
                    event.name.toLowerCase().includes(query) ||
                    event.description.toLowerCase().includes(query)
                );

                if (matchesTopic || matchingEvents.length > 0) {
                    return {
                        ...topic,
                        events: matchingEvents.length > 0 ? matchingEvents : topic.events
                    };
                }
                return null;
            }).filter((t): t is TopicView => t !== null);
        }

        // Apply status filter
        if (filterMode !== 'all') {
            result = result.map(topic => {
                const filteredEvents = topic.events.filter(event => {
                    const status = getEventStatus(topic, event.name);
                    return status === filterMode;
                });

                if (filteredEvents.length > 0) {
                    return { ...topic, events: filteredEvents };
                }
                return null;
            }).filter((t): t is TopicView => t !== null);
        }

        return result;
    }, [topics, searchQuery, filterMode, getEventStatus]);

    const filteredEventsCount = useMemo(() => {
        return filteredTopics.reduce((acc, t) => acc + t.events.length, 0);
    }, [filteredTopics]);

    return (
        <div className="space-y-6 animate-in">
            {/* Stats Cards */}
            {!loading && (
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <StatCard
                        label="Active Events"
                        value={activeCount}
                        description="With producers & consumers"
                        icon={<CheckCircle2 className="h-4 w-4 text-muted-foreground" />}
                        onClick={() => setFilterMode(filterMode === 'active' ? 'all' : 'active')}
                        active={filterMode === 'active'}
                    />
                    <StatCard
                        label="Unconsumed"
                        value={countUnconsumed}
                        description="No active consumers"
                        icon={<AlertCircle className="h-4 w-4 text-muted-foreground" />}
                        onClick={() => setFilterMode(filterMode === 'unconsumed' ? 'all' : 'unconsumed')}
                        active={filterMode === 'unconsumed'}
                    />
                    <StatCard
                        label="Orphaned"
                        value={countOrphaned}
                        description="No producers defined"
                        icon={<XCircle className="h-4 w-4 text-muted-foreground" />}
                        onClick={() => setFilterMode(filterMode === 'orphaned' ? 'all' : 'orphaned')}
                        active={filterMode === 'orphaned'}
                    />
                </div>
            )}

            {/* Search and Filter Bar */}
            <div className="flex flex-col sm:flex-row gap-3">
                <div className="flex-1">
                    <SearchInput
                        value={searchQuery}
                        onChange={setSearchQuery}
                        placeholder="Search topics, event types, or descriptions..."
                        resultCount={searchQuery || filterMode !== 'all'
                            ? filteredEventsCount
                            : undefined}
                        totalCount={searchQuery || filterMode !== 'all' ? countEvents : undefined}
                    />
                </div>
                <div className="flex gap-2">
                    <Button
                        variant={filterMode !== 'all' ? 'default' : 'outline'}
                        size="default"
                        onClick={() => setFilterMode('all')}
                        className="gap-2"
                    >
                        <Filter className="h-4 w-4" />
                        {filterMode === 'all' ? 'All Events' : `Filter: ${filterMode}`}
                    </Button>
                </div>
            </div>

            {/* Topics Layout (Master List) */}
            <div className="rounded-xl border border-border/60 bg-card/50 shadow-sm overflow-hidden">
                <Table>
                    <TableHeader className="bg-muted/30">
                        <TableRow className="hover:bg-transparent border-b border-border/60">
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider pl-4">Topic</TableHead>
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center">Events</TableHead>
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center">Producers</TableHead>
                            <TableHead className="h-11 text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center">Consumers</TableHead>
                            <TableHead className="h-11 text-right pr-4"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    <TableCell className="pl-4"><Skeleton className="h-12 w-[250px] rounded-md" /></TableCell>
                                    <TableCell className="text-center"><Skeleton className="h-6 w-[60px] mx-auto rounded-md" /></TableCell>
                                    <TableCell className="text-center"><Skeleton className="h-6 w-[50px] mx-auto rounded-md" /></TableCell>
                                    <TableCell className="text-center"><Skeleton className="h-6 w-[50px] mx-auto rounded-md" /></TableCell>
                                    <TableCell className="pr-4"><Skeleton className="h-8 w-8 ml-auto rounded-full" /></TableCell>
                                </TableRow>
                            ))
                        ) : filteredTopics.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="h-32 text-center p-0">
                                    <div className="py-8">
                                        <EmptyState
                                            title="No topics found"
                                            description="Try adjusting your search or filter criteria"
                                        />
                                    </div>
                                </TableCell>
                            </TableRow>
                        ) : (
                            filteredTopics.map((topicData) => {
                                const orphanedInTopic = topicData.events.filter(
                                    e => getEventStatus(topicData, e.name) === 'orphaned'
                                ).length;

                                return (
                                    <TableRow
                                        key={topicData.topic}
                                        className="cursor-pointer group hover:bg-muted/30 transition-colors"
                                        onClick={() => { setSelectedTopic(topicData); }}
                                    >
                                        <TableCell className="py-4 pl-4">
                                            <div className="flex items-center gap-4">
                                                <div className="h-10 w-10 shrink-0 rounded-lg bg-indigo-500/10 flex items-center justify-center border border-indigo-500/20 group-hover:bg-indigo-500/20 group-hover:border-indigo-500/30 transition-colors">
                                                    <Box className="h-5 w-5 text-indigo-500/80" />
                                                </div>
                                                <div className="flex flex-col gap-1.5">
                                                    <span className="font-semibold text-foreground/90 group-hover:text-primary transition-colors">{topicData.topic}</span>
                                                    {orphanedInTopic > 0 && (
                                                        <Badge variant="destructive" className="font-medium text-[10px] bg-destructive/10 text-destructive border-destructive/20 hover:bg-destructive/20 transition-colors gap-1 px-1.5 py-0 w-fit">
                                                            <XCircle className="h-3 w-3" />
                                                            Missing producer
                                                        </Badge>
                                                    )}
                                                </div>
                                            </div>
                                        </TableCell>
                                        <TableCell className="py-4 text-center">
                                            <Badge variant="secondary" className="font-medium text-[11px] bg-secondary/50 hover:bg-secondary transition-colors tabular-nums">
                                                {topicData.events.length} {topicData.events.length === 1 ? 'event' : 'events'}
                                            </Badge>
                                        </TableCell>
                                        <TableCell className="py-4 text-center">
                                            <div className="inline-flex items-center justify-center px-2.5 py-1 rounded-md bg-muted/40 border border-border/50 group-hover:bg-muted/60 transition-colors">
                                                <Layers className="h-3.5 w-3.5 text-muted-foreground mr-1.5" />
                                                <span className="text-xs font-semibold tabular-nums text-foreground/80">
                                                    {topicData.producers.length}
                                                </span>
                                            </div>
                                        </TableCell>
                                        <TableCell className="py-4 text-center">
                                            <div className="inline-flex items-center justify-center px-2.5 py-1 rounded-md bg-muted/40 border border-border/50 group-hover:bg-muted/60 transition-colors">
                                                <Zap className="h-3.5 w-3.5 text-muted-foreground mr-1.5" />
                                                <span className="text-xs font-semibold tabular-nums text-foreground/80">
                                                    {topicData.consumers.length}
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

            {pagination && (
                <Pagination
                    pagination={pagination}
                    page={page}
                    onPageChange={setPage}
                    label="topics"
                />
            )}

            <TopicDetailsPanel
                topic={selectedTopic}
                open={!!selectedTopic}
                onOpenChange={(open) => {
                    if (!open) setSelectedTopic(null);
                }}
            />
        </div>
    );
}
