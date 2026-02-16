import { SearchInput, StatCard } from "@/components/shared";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

import {
    Dialog,
    DialogTrigger,
} from "@/components/ui/dialog";


import { useTopics } from "@/hooks/useTopics";
import type { TopicConsumerEntry, TopicProducerEntry, TopicView } from "@/types";
import { AlertCircle, ArrowRight, Box, CheckCircle2, ChevronDown, ChevronLeft, ChevronRight, Code, FileJson, Filter, Layers, Search as SearchIcon, XCircle, Zap } from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import { EventDetails } from "./EventDetails";

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
    const [openTopics, setOpenTopics] = useState<Record<string, boolean>>({});
    const [filterMode, setFilterMode] = useState<FilterMode>('all');
    const [expandAll, setExpandAll] = useState(false);

    const getEventProducers = useCallback((topicData: TopicView, eventName: string): TopicProducerEntry[] => {
        return topicData.producers.filter(p => p.event === eventName);
    }, []);

    const getEventConsumers = useCallback((topicData: TopicView, eventName: string): TopicConsumerEntry[] => {
        return topicData.consumers.filter(c => c.event === eventName);
    }, []);

    const toggleTopic = (topicName: string) => {
        setOpenTopics(prev => ({
            ...prev,
            [topicName]: !prev[topicName]
        }));
    };

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

    const toggleAllTopics = () => {
        const newExpandAll = !expandAll;
        setExpandAll(newExpandAll);
        const newOpenTopics: Record<string, boolean> = {};
        filteredTopics.forEach(topic => {
            newOpenTopics[topic.topic] = newExpandAll;
        });
        setOpenTopics(newOpenTopics);
    };

    if (loading) return null;

    return (
        <div className="space-y-6 animate-in">
            {/* Stats Cards */}
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
                    <Button
                        variant="outline"
                        size="default"
                        onClick={toggleAllTopics}
                        className="gap-2"
                    >
                        {expandAll ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                        {expandAll ? 'Collapse' : 'Expand'}
                    </Button>
                </div>
            </div>

            {/* Events Layout */}
            <div className="rounded-lg border border-border bg-card overflow-hidden">
                {filteredTopics.length === 0 ? (
                    <div className="flex flex-col items-center justify-center p-10 text-muted-foreground">
                        <SearchIcon className="h-8 w-8 mb-3" />
                        <p className="font-medium text-sm">No events found</p>
                        <p className="text-xs mt-1">Try adjusting your search or filter criteria</p>
                    </div>
                ) : (
                    filteredTopics.map((topicData, topicIndex) => {
                        const orphanedInTopic = topicData.events.filter(
                            e => getEventStatus(topicData, e.name) === 'orphaned'
                        ).length;

                        return (
                            <div
                                key={topicData.topic}
                                className="border-b border-border last:border-b-0"
                                style={{
                                    animationDelay: `${topicIndex * 0.05}s`,
                                    animationFillMode: 'backwards'
                                }}
                            >
                                <button
                                    type="button"
                                    onClick={() => toggleTopic(topicData.topic)}
                                    className="w-full flex items-center justify-between px-4 py-3 md:px-5 hover:bg-accent/30 transition-colors"
                                >
                                    <div className="flex items-center gap-2.5 font-medium text-left">
                                        {openTopics[topicData.topic] ? (
                                            <ChevronDown className="h-4 w-4 text-muted-foreground" />
                                        ) : (
                                            <ChevronRight className="h-4 w-4 text-muted-foreground" />
                                        )}
                                        <Box className="h-4 w-4 text-muted-foreground" />
                                        <span className="font-mono text-sm font-medium">{topicData.topic}</span>
                                        <Badge variant="secondary" className="ml-1 text-[10px] font-normal">
                                            {topicData.events.length} {topicData.events.length === 1 ? 'event' : 'events'}
                                        </Badge>
                                        {orphanedInTopic > 0 && (
                                            <Badge variant="destructive" className="ml-0.5 text-[10px] font-normal gap-1">
                                                <XCircle className="h-3 w-3" />
                                                Missing producer
                                            </Badge>
                                        )}
                                    </div>
                                    <div className="hidden sm:flex items-center gap-2 text-xs text-muted-foreground">
                                        <span className="flex items-center gap-1"><Box className="h-3 w-3" /> Producers</span>
                                        <ArrowRight className="h-3 w-3" />
                                        <span className="flex items-center gap-1"><Layers className="h-3 w-3" /> Consumers</span>
                                    </div>
                                </button>

                                {openTopics[topicData.topic] && (
                                    <div className="px-4 pb-4 md:px-5 space-y-2 bg-muted/20">
                                        {topicData.events.map((event, eventIndex) => {
                                            const eventProducers = getEventProducers(topicData, event.name);
                                            const eventConsumers = getEventConsumers(topicData, event.name);
                                            const status = getEventStatus(topicData, event.name);

                                            let statusIcon;
                                            let statusColor;

                                            if (status === 'orphaned') {
                                                statusIcon = <XCircle className="h-3.5 w-3.5" />;
                                                statusColor = 'text-destructive';
                                            } else if (status === 'unconsumed') {
                                                statusIcon = <AlertCircle className="h-3.5 w-3.5" />;
                                                statusColor = 'text-muted-foreground';
                                            } else {
                                                statusIcon = <Zap className="h-3.5 w-3.5" />;
                                                statusColor = 'text-foreground';
                                            }

                                            return (
                                                <Dialog key={`${topicData.topic}-${event.name}`}>
                                                    <DialogTrigger asChild>
                                                        <button
                                                            type="button"
                                                            className="group rounded-lg border border-border bg-card p-4 hover:bg-accent/30 transition-colors cursor-pointer text-left w-full"
                                                            style={{
                                                                animationDelay: `${(topicIndex * 0.05) + (eventIndex * 0.03)}s`,
                                                                animationFillMode: 'backwards'
                                                            }}
                                                        >
                                                            <div className="flex flex-col gap-2.5">
                                                                <div className="flex items-start justify-between gap-4">
                                                                    <div className="flex flex-col gap-1.5">
                                                                        <div className="flex items-center gap-2">
                                                                            <span className={statusColor}>{statusIcon}</span>
                                                                            <div className="flex flex-wrap items-center gap-2">
                                                                                <span className="font-mono text-sm font-medium">{event.name}</span>
                                                                                {event.version && (
                                                                                    <Badge variant="outline" className="font-mono text-[10px] font-normal px-2 py-0.5">
                                                                                        v{event.version}
                                                                                    </Badge>
                                                                                )}
                                                                            </div>
                                                                        </div>
                                                                        <span className="text-xs text-muted-foreground">
                                                                            {event.description}
                                                                        </span>
                                                                    </div>

                                                                    <div className="flex items-center gap-1.5">
                                                                        <Button
                                                                            variant="ghost"
                                                                            size="sm"
                                                                            className="h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                                                                            title="View Details"
                                                                        >
                                                                            <FileJson className="h-3.5 w-3.5" />
                                                                        </Button>

                                                                        {event.schema_url && (
                                                                            <div onClick={(e) => e.stopPropagation()}>
                                                                                <Button
                                                                                    variant="ghost"
                                                                                    size="sm"
                                                                                    className="h-7 w-7 p-0"
                                                                                    asChild
                                                                                    title="Open Schema URL"
                                                                                >
                                                                                    <a
                                                                                        href={event.schema_url}
                                                                                        target="_blank"
                                                                                        rel="noopener noreferrer"
                                                                                    >
                                                                                        <Code className="h-3.5 w-3.5" />
                                                                                    </a>
                                                                                </Button>
                                                                            </div>
                                                                        )}
                                                                    </div>
                                                                </div>

                                                                <div className="grid gap-3 md:grid-cols-[1fr_auto_1fr] items-start pt-1">
                                                                    <div className="flex flex-col gap-1.5">
                                                                        <span className="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">
                                                                            Producers
                                                                        </span>
                                                                        {eventProducers.length === 0 ? (
                                                                            <span className="text-xs text-destructive italic">No producers</span>
                                                                        ) : (
                                                                            <div className="flex flex-wrap gap-1.5">
                                                                                {eventProducers.slice(0, 3).map(p => (
                                                                                    <span key={`${p.service}-${p.repository}`} className="text-xs bg-muted px-2 py-0.5 rounded font-medium">
                                                                                        {p.service}
                                                                                    </span>
                                                                                ))}
                                                                                {eventProducers.length > 3 && (
                                                                                    <Badge variant="secondary" className="text-[10px] font-normal">
                                                                                        +{eventProducers.length - 3}
                                                                                    </Badge>
                                                                                )}
                                                                            </div>
                                                                        )}
                                                                    </div>

                                                                    <div className="hidden md:flex items-center justify-center">
                                                                        <ArrowRight className="h-3.5 w-3.5 text-muted-foreground" />
                                                                    </div>

                                                                    <div className="flex flex-col gap-1.5">
                                                                        <span className="text-[10px] font-medium text-muted-foreground uppercase tracking-wider">
                                                                            Consumers
                                                                        </span>
                                                                        {eventConsumers.length === 0 ? (
                                                                            <span className="text-xs text-muted-foreground italic">No consumers</span>
                                                                        ) : (
                                                                            <div className="flex flex-wrap gap-1.5">
                                                                                {eventConsumers.slice(0, 3).map(c => (
                                                                                    <span key={`${c.service}-${c.repository}-${c.group}`} className="text-xs bg-muted px-2 py-0.5 rounded font-medium">
                                                                                        {c.service}
                                                                                    </span>
                                                                                ))}
                                                                                {eventConsumers.length > 3 && (
                                                                                    <Badge variant="secondary" className="text-[10px] font-normal">
                                                                                        +{eventConsumers.length - 3}
                                                                                    </Badge>
                                                                                )}
                                                                            </div>
                                                                        )}
                                                                    </div>
                                                                </div>
                                                            </div>
                                                        </button>
                                                    </DialogTrigger>
                                                    <EventDetails eventId={event.id} />

                                                </Dialog>
                                            );
                                        })}
                                    </div>
                                )
                                }
                            </div>
                        );
                    })
                )}
            </div>

            {/* Pagination Controls */}
            {
                pagination && pagination.total_pages > 1 && (
                    <div className="flex items-center justify-between px-1">
                        <span className="text-xs text-muted-foreground">
                            Page {pagination.page} of {pagination.total_pages} ({pagination.total} topics)
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
                )
            }
        </div >
    );
}
