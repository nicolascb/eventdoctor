import { SearchInput, StatCard } from "@/components/shared";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { Consumer, Producer, TopicWithEvents } from "@/types";
import { AlertCircle, ArrowRight, Box, CheckCircle2, ChevronDown, ChevronRight, Code, FileJson, Filter, Layers, Search as SearchIcon, Workflow, XCircle, Zap } from "lucide-react";
import { useCallback, useMemo, useState } from "react";

interface TopicsViewProps {
    topics: TopicWithEvents[];
    producers: Producer[];
    consumers: Consumer[];
}

type FilterMode = 'all' | 'orphaned' | 'unconsumed' | 'active';

export function TopicsView({ topics, producers, consumers }: TopicsViewProps) {
    const [searchQuery, setSearchQuery] = useState("");
    const [openTopics, setOpenTopics] = useState<Record<string, boolean>>({});
    const [filterMode, setFilterMode] = useState<FilterMode>('all');
    const [expandAll, setExpandAll] = useState(false);

    const getEventProducers = useCallback((topic: string, eventType: string) => {
        return producers.filter(p => p.topic === topic && p.events.some(e => e.name === eventType));
    }, [producers]);

    const getEventConsumers = useCallback((topic: string, eventType: string) => {
        return consumers.filter(c =>
            c.topics.some(t => t.name === topic && t.events.some(e => e.name === eventType))
        );
    }, [consumers]);

    const toggleTopic = (topicName: string) => {
        setOpenTopics(prev => ({
            ...prev,
            [topicName]: !prev[topicName]
        }));
    };

    const getEventStatus = useCallback((topic: string, eventType: string) => {
        const eventProducers = getEventProducers(topic, eventType);
        const eventConsumers = getEventConsumers(topic, eventType);

        if (eventProducers.length === 0) return 'orphaned';
        if (eventConsumers.length === 0) return 'unconsumed';
        return 'active';
    }, [getEventProducers, getEventConsumers]);

    const totalEvents = useMemo(() => {
        return topics.reduce((acc, t) => acc + t.events.length, 0);
    }, [topics]);

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
            }).filter((t): t is TopicWithEvents => t !== null);
        }

        // Apply status filter
        if (filterMode !== 'all') {
            result = result.map(topic => {
                const filteredEvents = topic.events.filter(event => {
                    const status = getEventStatus(topic.topic, event.name);
                    return status === filterMode;
                });

                if (filteredEvents.length > 0) {
                    return { ...topic, events: filteredEvents };
                }
                return null;
            }).filter((t): t is TopicWithEvents => t !== null);
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

    const statsData = useMemo(() => {
        let orphanedCount = 0;
        let unconsumedCount = 0;
        let activeCount = 0;

        topics.forEach(topic => {
            topic.events.forEach(event => {
                const status = getEventStatus(topic.topic, event.name);
                if (status === 'orphaned') orphanedCount++;
                else if (status === 'unconsumed') unconsumedCount++;
                else activeCount++;
            });
        });

        return { orphanedCount, unconsumedCount, activeCount };
    }, [topics, getEventStatus]);


    return (
        <div className="space-y-6 animate-in">
            {/* Stats Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <StatCard
                    label="Active Events"
                    value={statsData.activeCount}
                    description="With producers & consumers"
                    icon={<CheckCircle2 className="h-4 w-4 text-muted-foreground" />}
                    onClick={() => setFilterMode(filterMode === 'active' ? 'all' : 'active')}
                    active={filterMode === 'active'}
                />
                <StatCard
                    label="Unconsumed"
                    value={statsData.unconsumedCount}
                    description="No active consumers"
                    icon={<AlertCircle className="h-4 w-4 text-muted-foreground" />}
                    onClick={() => setFilterMode(filterMode === 'unconsumed' ? 'all' : 'unconsumed')}
                    active={filterMode === 'unconsumed'}
                />
                <StatCard
                    label="Orphaned"
                    value={statsData.orphanedCount}
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
                        totalCount={searchQuery || filterMode !== 'all' ? totalEvents : undefined}
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
                            e => getEventStatus(topicData.topic, e.name) === 'orphaned'
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
                                            const eventProducers = getEventProducers(topicData.topic, event.name);
                                            const eventConsumers = getEventConsumers(topicData.topic, event.name);
                                            const status = getEventStatus(topicData.topic, event.name);

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
                                                    <DialogContent className="max-w-4xl max-h-[85vh] overflow-y-auto">
                                                        <DialogHeader>
                                                            <div className="flex items-center gap-2 mb-1">
                                                                <span className={statusColor}>{statusIcon}</span>
                                                                <DialogTitle className="font-mono flex items-center gap-2 text-base">
                                                                    <span className="text-muted-foreground font-normal">{topicData.topic}</span>
                                                                    <span className="text-muted-foreground/30">/</span>
                                                                    <span className="font-semibold">{event.name}</span>
                                                                    {event.version && (
                                                                        <Badge variant="outline" className="font-mono text-[10px] font-normal ml-1">
                                                                            v{event.version}
                                                                        </Badge>
                                                                    )}
                                                                </DialogTitle>
                                                            </div>
                                                            <DialogDescription className="text-sm">
                                                                {event.description}
                                                            </DialogDescription>
                                                        </DialogHeader>

                                                        <div className="py-4">
                                                            <Tabs defaultValue="structure" className="w-full">
                                                                <TabsList className="mb-4 grid w-full grid-cols-2">
                                                                    <TabsTrigger value="structure" className="flex items-center gap-2 text-xs">
                                                                        <Code className="h-3.5 w-3.5" />
                                                                        Schema & Structure
                                                                    </TabsTrigger>
                                                                    <TabsTrigger value="flow" className="flex items-center gap-2 text-xs">
                                                                        <Workflow className="h-3.5 w-3.5" />
                                                                        Service Flow
                                                                    </TabsTrigger>
                                                                </TabsList>

                                                                <TabsContent value="structure" className="space-y-4">
                                                                    <div className="grid grid-cols-2 gap-4">
                                                                        <Card>
                                                                            <CardHeader className="pb-2">
                                                                                <CardTitle className="text-xs font-medium flex items-center gap-2">
                                                                                    <FileJson className="h-3.5 w-3.5 text-muted-foreground" />
                                                                                    Schema Overview
                                                                                </CardTitle>
                                                                            </CardHeader>
                                                                            <CardContent className="space-y-2">
                                                                                <div className="flex justify-between items-center py-1.5 border-b border-border text-xs">
                                                                                    <span className="text-muted-foreground">Headers</span>
                                                                                    <span className="font-mono">
                                                                                        {event.headers ? event.headers.length : 0}
                                                                                    </span>
                                                                                </div>
                                                                                {event.properties && (
                                                                                    <div className="flex justify-between items-center py-1.5 border-b border-border text-xs">
                                                                                        <span className="text-muted-foreground">Properties</span>
                                                                                        <span className="font-mono">
                                                                                            {Object.keys(event.properties).length}
                                                                                        </span>
                                                                                    </div>
                                                                                )}
                                                                                {event.schema_url && (
                                                                                    <div className="pt-1">
                                                                                        <Button
                                                                                            variant="outline"
                                                                                            size="sm"
                                                                                            className="w-full gap-2 text-xs"
                                                                                            asChild
                                                                                        >
                                                                                            <a
                                                                                                href={event.schema_url}
                                                                                                target="_blank"
                                                                                                rel="noopener noreferrer"
                                                                                            >
                                                                                                <Code className="h-3.5 w-3.5" />
                                                                                                View External Schema
                                                                                            </a>
                                                                                        </Button>
                                                                                    </div>
                                                                                )}
                                                                            </CardContent>
                                                                        </Card>
                                                                        <Card>
                                                                            <CardHeader className="pb-2">
                                                                                <CardTitle className="text-xs font-medium flex items-center gap-2">
                                                                                    <Workflow className="h-3.5 w-3.5 text-muted-foreground" />
                                                                                    Usage Statistics
                                                                                </CardTitle>
                                                                            </CardHeader>
                                                                            <CardContent className="space-y-2">
                                                                                <div className="flex justify-between items-center py-1.5 border-b border-border text-xs">
                                                                                    <span className="text-muted-foreground">Producers</span>
                                                                                    <span className="font-mono">
                                                                                        {eventProducers.length}
                                                                                    </span>
                                                                                </div>
                                                                                <div className="flex justify-between items-center py-1.5 border-b border-border text-xs">
                                                                                    <span className="text-muted-foreground">Consumers</span>
                                                                                    <span className="font-mono">
                                                                                        {eventConsumers.length}
                                                                                    </span>
                                                                                </div>
                                                                                <div className="flex justify-between items-center py-1.5 text-xs">
                                                                                    <span className="text-muted-foreground">Status</span>
                                                                                    <Badge variant="secondary" className="text-[10px] font-normal gap-1">
                                                                                        {statusIcon}
                                                                                        {status.charAt(0).toUpperCase() + status.slice(1)}
                                                                                    </Badge>
                                                                                </div>
                                                                            </CardContent>
                                                                        </Card>
                                                                    </div>

                                                                    {event.properties && (
                                                                        <div className="space-y-2">
                                                                            <h4 className="text-xs font-medium flex items-center gap-2">
                                                                                <Code className="h-3.5 w-3.5 text-muted-foreground" />
                                                                                Properties Schema
                                                                            </h4>
                                                                            <div className="bg-muted p-4 rounded-lg border border-border overflow-x-auto">
                                                                                <pre className="text-xs font-mono leading-relaxed">
                                                                                    {JSON.stringify(event.properties, null, 2)}
                                                                                </pre>
                                                                            </div>
                                                                        </div>
                                                                    )}

                                                                    {event.headers && event.headers.length > 0 && (
                                                                        <div className="space-y-2">
                                                                            <h4 className="text-xs font-medium flex items-center gap-2">
                                                                                <FileJson className="h-3.5 w-3.5 text-muted-foreground" />
                                                                                Headers
                                                                            </h4>
                                                                            <div className="rounded-lg border border-border overflow-hidden">
                                                                                <Table>
                                                                                    <TableHeader>
                                                                                        <TableRow className="hover:bg-transparent">
                                                                                            <TableHead className="text-[11px] font-medium">Name</TableHead>
                                                                                            <TableHead className="text-[11px] font-medium">Description</TableHead>
                                                                                        </TableRow>
                                                                                    </TableHeader>
                                                                                    <TableBody>
                                                                                        {event.headers.map((header) => (
                                                                                            <TableRow key={header.name}>
                                                                                                <TableCell className="font-mono text-xs py-2">
                                                                                                    {header.name}
                                                                                                </TableCell>
                                                                                                <TableCell className="text-xs text-muted-foreground py-2">
                                                                                                    {header.description || '—'}
                                                                                                </TableCell>
                                                                                            </TableRow>
                                                                                        ))}
                                                                                    </TableBody>
                                                                                </Table>
                                                                            </div>
                                                                        </div>
                                                                    )}
                                                                </TabsContent>

                                                                <TabsContent value="flow" className="space-y-4">
                                                                    <div className="grid md:grid-cols-2 gap-4">
                                                                        <div className="space-y-3">
                                                                            <div className="flex items-center gap-2 pb-2 border-b border-border">
                                                                                <Box className="h-3.5 w-3.5 text-muted-foreground" />
                                                                                <h4 className="text-xs font-medium">Producers ({eventProducers.length})</h4>
                                                                            </div>
                                                                            {eventProducers.length > 0 ? (
                                                                                <div className="space-y-2">
                                                                                    {eventProducers.map(p => (
                                                                                        <div key={`${p.service}-${p.repository}-${p.topic}`} className="p-3 bg-muted/50 rounded-lg border border-border flex flex-col gap-1">
                                                                                            <span className="font-medium text-sm">{p.service}</span>
                                                                                            <code className="text-xs text-muted-foreground font-mono">
                                                                                                {p.topic}
                                                                                            </code>
                                                                                        </div>
                                                                                    ))}
                                                                                </div>
                                                                            ) : (
                                                                                <div className="text-center p-6 bg-muted/30 rounded-lg border border-dashed border-border">
                                                                                    <XCircle className="h-6 w-6 mx-auto mb-2 text-destructive" />
                                                                                    <p className="text-xs font-medium">No producers registered</p>
                                                                                </div>
                                                                            )}
                                                                        </div>

                                                                        <div className="space-y-3">
                                                                            <div className="flex items-center gap-2 pb-2 border-b border-border">
                                                                                <Layers className="h-3.5 w-3.5 text-muted-foreground" />
                                                                                <h4 className="text-xs font-medium">Consumers ({eventConsumers.length})</h4>
                                                                            </div>
                                                                            {eventConsumers.length > 0 ? (
                                                                                <div className="space-y-2">
                                                                                    {eventConsumers.map(c => (
                                                                                        <div key={`${c.service}-${c.repository}-${c.group}`} className="p-3 bg-muted/50 rounded-lg border border-border flex flex-col gap-1">
                                                                                            <span className="font-medium text-sm">{c.service}</span>
                                                                                            <span className="text-xs text-muted-foreground">{c.group}</span>
                                                                                        </div>
                                                                                    ))}
                                                                                </div>
                                                                            ) : (
                                                                                <div className="text-center p-6 bg-muted/30 rounded-lg border border-dashed border-border">
                                                                                    <AlertCircle className="h-6 w-6 mx-auto mb-2 text-muted-foreground" />
                                                                                    <p className="text-xs font-medium">No consumers registered</p>
                                                                                </div>
                                                                            )}
                                                                        </div>
                                                                    </div>
                                                                </TabsContent>
                                                            </Tabs>
                                                        </div>
                                                    </DialogContent>
                                                </Dialog>
                                            );
                                        })}
                                    </div>
                                )}
                            </div>
                        );
                    })
                )}
            </div>
        </div>
    );
}
