import { PageHeader, SearchInput, StatCard } from "@/components/shared";
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
        <div className="space-y-8 animate-in">
            {/* Header Section */}
            <div className="relative">
                <div className="flex flex-col gap-6">
                    <PageHeader
                        title="Topics"
                        badge={`${topics.length} topics · ${totalEvents} events`}
                    />

                    {/* Stats Cards */}
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                        <StatCard
                            label="Active Events"
                            value={statsData.activeCount}
                            description="With producers & consumers"
                            icon={<CheckCircle2 className="h-6 w-6 text-primary" />}
                            onClick={() => setFilterMode(filterMode === 'active' ? 'all' : 'active')}
                            active={filterMode === 'active'}
                        />
                        <StatCard
                            label="Unconsumed"
                            value={statsData.unconsumedCount}
                            description="No active consumers"
                            icon={<AlertCircle className="h-6 w-6 text-accent" />}
                            iconClassName="bg-accent/10"
                            onClick={() => setFilterMode(filterMode === 'unconsumed' ? 'all' : 'unconsumed')}
                            active={filterMode === 'unconsumed'}
                        />
                        <StatCard
                            label="Orphaned"
                            value={statsData.orphanedCount}
                            description="No producers defined"
                            icon={<XCircle className="h-6 w-6 text-destructive" />}
                            iconClassName="bg-destructive/10"
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
                                {expandAll ? 'Collapse All' : 'Expand All'}
                            </Button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Events Layout */}
            <div className="rounded-xl border border-border bg-card shadow-sm overflow-hidden">
                {filteredTopics.length === 0 ? (
                    <div className="flex flex-col items-center justify-center p-10 text-muted-foreground">
                        <div className="h-16 w-16 rounded-full bg-muted/50 flex items-center justify-center mb-4">
                            <SearchIcon className="h-8 w-8" />
                        </div>
                        <p className="font-semibold text-lg">No events found</p>
                        <p className="text-sm mt-1">Try adjusting your search or filter criteria</p>
                    </div>
                ) : (
                    filteredTopics.map((topicData, topicIndex) => (
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
                                className="w-full flex items-center justify-between px-4 py-4 md:px-6 bg-muted/20 hover:bg-muted/40 transition-colors"
                            >
                                <div className="flex items-center gap-3 font-medium text-left">
                                    <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
                                        {openTopics[topicData.topic] ? (
                                            <ChevronDown className="h-5 w-5 text-primary" />
                                        ) : (
                                            <ChevronRight className="h-5 w-5 text-primary" />
                                        )}
                                    </div>
                                    <Box className="h-5 w-5 text-primary" />
                                    <span className="font-mono text-base font-semibold">{topicData.topic}</span>
                                    <Badge variant="secondary" className="ml-2 text-xs h-6 px-2.5 font-medium">
                                        {topicData.events.length} {topicData.events.length === 1 ? 'event' : 'events'}
                                    </Badge>
                                </div>
                                <div className="hidden sm:flex items-center gap-2 text-xs text-muted-foreground">
                                    <span className="flex items-center gap-1"><Box className="h-3 w-3" /> Producers</span>
                                    <ArrowRight className="h-3 w-3" />
                                    <span className="flex items-center gap-1"><Layers className="h-3 w-3" /> Consumers</span>
                                </div>
                            </button>

                            {openTopics[topicData.topic] && (
                                <div className="px-4 pb-6 md:px-6 space-y-3 bg-card">
                                    {topicData.events.map((event, eventIndex) => {
                                        const eventProducers = getEventProducers(topicData.topic, event.name);
                                        const eventConsumers = getEventConsumers(topicData.topic, event.name);
                                        const status = getEventStatus(topicData.topic, event.name);

                                        let statusIcon;
                                        let statusColor;
                                        let statusBg;

                                        if (status === 'orphaned') {
                                            statusIcon = <XCircle className="h-3.5 w-3.5" />;
                                            statusColor = 'text-destructive';
                                            statusBg = 'bg-destructive/10';
                                        } else if (status === 'unconsumed') {
                                            statusIcon = <AlertCircle className="h-3.5 w-3.5" />;
                                            statusColor = 'text-accent';
                                            statusBg = 'bg-accent/10';
                                        } else {
                                            statusIcon = <Zap className="h-3.5 w-3.5" />;
                                            statusColor = 'text-primary';
                                            statusBg = 'bg-primary/10';
                                        }

                                        return (
                                            <Dialog key={`${topicData.topic}-${event.name}`}>
                                                <DialogTrigger asChild>
                                                    <button
                                                        type="button"
                                                        className="group relative rounded-xl border border-border bg-background/60 p-4 md:p-5 shadow-xs hover:border-primary/50 transition-colors cursor-pointer text-left w-full"
                                                        style={{
                                                            animationDelay: `${(topicIndex * 0.05) + (eventIndex * 0.03)}s`,
                                                            animationFillMode: 'backwards'
                                                        }}
                                                    >
                                                        <div className="flex flex-col gap-3">
                                                            <div className="flex items-start justify-between gap-4">
                                                                <div className="flex flex-col gap-2">
                                                                    <div className="flex items-center gap-2">
                                                                        <div className={`h-8 w-8 rounded-lg ${statusBg} flex items-center justify-center ${statusColor}`}>
                                                                            {statusIcon}
                                                                        </div>
                                                                        <div className="flex flex-wrap items-center gap-2">
                                                                            <span className="font-mono text-sm font-semibold">{event.name}</span>
                                                                            {event.version && (
                                                                                <Badge variant="outline" className="font-mono text-xs font-medium px-2.5 py-1">
                                                                                    v{event.version}
                                                                                </Badge>
                                                                            )}
                                                                        </div>
                                                                    </div>
                                                                    <span className="text-sm text-muted-foreground leading-relaxed">
                                                                        {event.description}
                                                                    </span>
                                                                </div>

                                                                <div className="flex items-center gap-2">
                                                                    <Button
                                                                        variant="ghost"
                                                                        size="sm"
                                                                        className="h-9 w-9 p-0 opacity-60 group-hover:opacity-100 transition-opacity hover:bg-primary/10"
                                                                        title="View Details"
                                                                    >
                                                                        <FileJson className="h-4 w-4" />
                                                                    </Button>

                                                                    {event.schema_url && (
                                                                        <div onClick={(e) => e.stopPropagation()}>
                                                                            <Button
                                                                                variant="ghost"
                                                                                size="sm"
                                                                                className="h-9 w-9 p-0 hover:bg-primary/10"
                                                                                asChild
                                                                                title="Open Schema URL"
                                                                            >
                                                                                <a
                                                                                    href={event.schema_url}
                                                                                    target="_blank"
                                                                                    rel="noopener noreferrer"
                                                                                >
                                                                                    <Code className="h-4 w-4" />
                                                                                </a>
                                                                            </Button>
                                                                        </div>
                                                                    )}
                                                                </div>
                                                            </div>

                                                            <div className="grid gap-4 md:grid-cols-[1fr_auto_1fr] items-start pt-1">
                                                                <div className="flex flex-col gap-2">
                                                                    <div className="flex items-center gap-2 text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                                                                        <Box className="h-3.5 w-3.5 text-primary" />
                                                                        Producers
                                                                    </div>
                                                                    {eventProducers.length === 0 ? (
                                                                        <div className="flex items-center gap-2 px-3 py-2 rounded-md bg-destructive/5 border border-destructive/20">
                                                                            <XCircle className="h-3.5 w-3.5 text-destructive" />
                                                                            <span className="text-xs text-destructive font-medium italic">No producers</span>
                                                                        </div>
                                                                    ) : (
                                                                        <div className="flex flex-wrap gap-2">
                                                                            {eventProducers.slice(0, 3).map(p => (
                                                                                <div key={p.service} className="flex items-center gap-2 px-3 py-1.5 rounded-md bg-muted/30 border border-border hover:bg-muted/50 transition-colors">
                                                                                    <span className="text-xs text-foreground font-medium truncate max-w-[140px]">
                                                                                        {p.service}
                                                                                    </span>
                                                                                </div>
                                                                            ))}
                                                                            {eventProducers.length > 3 && (
                                                                                <Badge variant="secondary" className="text-[10px] h-5 px-2">
                                                                                    +{eventProducers.length - 3}
                                                                                </Badge>
                                                                            )}
                                                                        </div>
                                                                    )}
                                                                </div>

                                                                <div className="hidden md:flex items-center justify-center">
                                                                    <ArrowRight className="h-4 w-4 text-primary" />
                                                                </div>

                                                                <div className="flex flex-col gap-2">
                                                                    <div className="flex items-center gap-2 text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                                                                        <Layers className="h-3.5 w-3.5 text-primary" />
                                                                        Consumers
                                                                    </div>
                                                                    {eventConsumers.length === 0 ? (
                                                                        <div className="flex items-center gap-2 px-3 py-2 rounded-md bg-accent/5 border border-accent/20">
                                                                            <AlertCircle className="h-3.5 w-3.5 text-accent" />
                                                                            <span className="text-xs text-accent font-medium italic">No consumers</span>
                                                                        </div>
                                                                    ) : (
                                                                        <div className="flex flex-wrap gap-2">
                                                                            {eventConsumers.slice(0, 3).map(c => (
                                                                                <div key={c.service} className="flex items-center gap-2 px-3 py-1.5 rounded-md bg-muted/30 border border-border hover:bg-muted/50 transition-colors">
                                                                                    <span className="text-xs text-foreground font-medium truncate max-w-[140px]">
                                                                                        {c.service}
                                                                                    </span>
                                                                                </div>
                                                                            ))}
                                                                            {eventConsumers.length > 3 && (
                                                                                <Badge variant="secondary" className="text-[10px] h-5 px-2">
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
                                                        <div className="flex items-center gap-3 mb-2">
                                                            <div className={`h-10 w-10 rounded-lg ${statusBg} flex items-center justify-center ${statusColor}`}>
                                                                {statusIcon}
                                                            </div>
                                                            <div>
                                                                <DialogTitle className="font-mono flex items-center gap-2 text-lg">
                                                                    <span className="text-muted-foreground font-normal">{topicData.topic}</span>
                                                                    <span className="text-muted-foreground/30">/</span>
                                                                    <span className="font-semibold">{event.name}</span>
                                                                    {event.version && (
                                                                        <Badge variant="outline" className="font-mono text-xs ml-1">
                                                                            v{event.version}
                                                                        </Badge>
                                                                    )}
                                                                </DialogTitle>
                                                                <DialogDescription className="text-sm mt-1.5">
                                                                    {event.description}
                                                                </DialogDescription>
                                                            </div>
                                                        </div>
                                                    </DialogHeader>

                                                    <div className="py-6">
                                                        <Tabs defaultValue="structure" className="w-full">
                                                            <TabsList className="mb-6 grid w-full grid-cols-2">
                                                                <TabsTrigger value="structure" className="flex items-center gap-2">
                                                                    <Code className="h-4 w-4" />
                                                                    Schema & Structure
                                                                </TabsTrigger>
                                                                <TabsTrigger value="flow" className="flex items-center gap-2">
                                                                    <Workflow className="h-4 w-4" />
                                                                    Service Flow
                                                                </TabsTrigger>
                                                            </TabsList>

                                                            <TabsContent value="structure" className="space-y-6">
                                                                <div className="grid grid-cols-2 gap-4">
                                                                    <Card className="border-2">
                                                                        <CardHeader className="pb-3">
                                                                            <CardTitle className="text-sm font-semibold flex items-center gap-2">
                                                                                <FileJson className="h-4 w-4 text-primary" />
                                                                                Schema Overview
                                                                            </CardTitle>
                                                                        </CardHeader>
                                                                        <CardContent className="space-y-3">
                                                                            <div className="flex justify-between items-center py-2 border-b border-border">
                                                                                <span className="text-sm text-muted-foreground">Headers</span>
                                                                                <Badge variant="secondary" className="font-mono">
                                                                                    {event.headers ? event.headers.length : 0}
                                                                                </Badge>
                                                                            </div>
                                                                            {event.properties && (
                                                                                <div className="flex justify-between items-center py-2 border-b border-border">
                                                                                    <span className="text-sm text-muted-foreground">Properties</span>
                                                                                    <Badge variant="secondary" className="font-mono">
                                                                                        {Object.keys(event.properties).length}
                                                                                    </Badge>
                                                                                </div>
                                                                            )}
                                                                            {event.schema_url && (
                                                                                <div className="pt-2">
                                                                                    <Button
                                                                                        variant="outline"
                                                                                        size="sm"
                                                                                        className="w-full gap-2"
                                                                                        asChild
                                                                                    >
                                                                                        <a
                                                                                            href={event.schema_url}
                                                                                            target="_blank"
                                                                                            rel="noopener noreferrer"
                                                                                        >
                                                                                            <Code className="h-4 w-4" />
                                                                                            View External Schema
                                                                                        </a>
                                                                                    </Button>
                                                                                </div>
                                                                            )}
                                                                        </CardContent>
                                                                    </Card>
                                                                    <Card className="border-2">
                                                                        <CardHeader className="pb-3">
                                                                            <CardTitle className="text-sm font-semibold flex items-center gap-2">
                                                                                <Workflow className="h-4 w-4 text-primary" />
                                                                                Usage Statistics
                                                                            </CardTitle>
                                                                        </CardHeader>
                                                                        <CardContent className="space-y-3">
                                                                            <div className="flex justify-between items-center py-2 border-b border-border">
                                                                                <span className="text-sm text-muted-foreground">Producers</span>
                                                                                <Badge variant="secondary" className="font-mono">
                                                                                    {eventProducers.length}
                                                                                </Badge>
                                                                            </div>
                                                                            <div className="flex justify-between items-center py-2 border-b border-border">
                                                                                <span className="text-sm text-muted-foreground">Consumers</span>
                                                                                <Badge variant="secondary" className="font-mono">
                                                                                    {eventConsumers.length}
                                                                                </Badge>
                                                                            </div>
                                                                            <div className="flex justify-between items-center py-2">
                                                                                <span className="text-sm text-muted-foreground">Status</span>
                                                                                <Badge
                                                                                    variant={status === 'active' ? 'default' : 'secondary'}
                                                                                    className={`${status !== 'active' ? statusColor : ''} gap-1`}
                                                                                >
                                                                                    {statusIcon}
                                                                                    {status.charAt(0).toUpperCase() + status.slice(1)}
                                                                                </Badge>
                                                                            </div>
                                                                        </CardContent>
                                                                    </Card>
                                                                </div>

                                                                {event.properties && (
                                                                    <div className="space-y-3">
                                                                        <h4 className="text-sm font-semibold flex items-center gap-2">
                                                                            <Code className="h-4 w-4 text-muted-foreground" />
                                                                            Properties Schema
                                                                        </h4>
                                                                        <div className="bg-muted/50 p-5 rounded-lg border-2 border-border overflow-x-auto">
                                                                            <pre className="text-xs font-mono leading-relaxed">
                                                                                {JSON.stringify(event.properties, null, 2)}
                                                                            </pre>
                                                                        </div>
                                                                    </div>
                                                                )}

                                                                {event.headers && event.headers.length > 0 && (
                                                                    <div className="space-y-3">
                                                                        <h4 className="text-sm font-semibold flex items-center gap-2">
                                                                            <FileJson className="h-4 w-4 text-muted-foreground" />
                                                                            Headers
                                                                        </h4>
                                                                        <div className="rounded-lg border-2 border-border overflow-hidden">
                                                                            <Table>
                                                                                <TableHeader>
                                                                                    <TableRow className="hover:bg-transparent">
                                                                                        <TableHead className="text-xs font-semibold">Name</TableHead>
                                                                                        <TableHead className="text-xs font-semibold">Description</TableHead>
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

                                                            <TabsContent value="flow" className="space-y-6">
                                                                <div className="grid md:grid-cols-2 gap-6">
                                                                    <div className="space-y-4">
                                                                        <div className="flex items-center gap-2 pb-3 border-b-2 border-border">
                                                                            <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
                                                                                <Box className="h-4 w-4 text-primary" />
                                                                            </div>
                                                                            <h4 className="text-sm font-semibold">Producers ({eventProducers.length})</h4>
                                                                        </div>
                                                                        {eventProducers.length > 0 ? (
                                                                            <div className="space-y-2">
                                                                                {eventProducers.map(p => (
                                                                                    <div key={`${p.service}-${p.topic}`} className="p-4 bg-card rounded-lg border-2 border-border hover:border-primary/50 transition-colors flex flex-col gap-2">
                                                                                        <span className="font-semibold text-sm">{p.service}</span>
                                                                                        <span className="text-xs text-muted-foreground font-mono bg-muted/50 px-2 py-1 rounded inline-block">
                                                                                            {p.topic}
                                                                                        </span>
                                                                                    </div>
                                                                                ))}
                                                                            </div>
                                                                        ) : (
                                                                            <div className="text-sm text-muted-foreground italic p-6 text-center bg-muted/20 rounded-lg border-2 border-dashed border-border">
                                                                                <XCircle className="h-8 w-8 mx-auto mb-2 text-destructive" />
                                                                                <p className="font-medium">No producers registered</p>
                                                                                <p className="text-xs mt-1">This event has no active producers</p>
                                                                            </div>
                                                                        )}
                                                                    </div>

                                                                    <div className="space-y-4">
                                                                        <div className="flex items-center gap-2 pb-3 border-b-2 border-border">
                                                                            <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
                                                                                <Layers className="h-4 w-4 text-primary" />
                                                                            </div>
                                                                            <h4 className="text-sm font-semibold">Consumers ({eventConsumers.length})</h4>
                                                                        </div>
                                                                        {eventConsumers.length > 0 ? (
                                                                            <div className="space-y-2">
                                                                                {eventConsumers.map(c => (
                                                                                    <div key={c.group} className="p-4 bg-card rounded-lg border-2 border-border hover:border-primary/50 transition-colors flex flex-col gap-2">
                                                                                        <span className="font-semibold text-sm">{c.service}</span>
                                                                                        <span className="text-xs text-muted-foreground">{c.group}</span>
                                                                                    </div>
                                                                                ))}
                                                                            </div>
                                                                        ) : (
                                                                            <div className="text-sm text-muted-foreground italic p-6 text-center bg-muted/20 rounded-lg border-2 border-dashed border-border">
                                                                                <AlertCircle className="h-8 w-8 mx-auto mb-2 text-accent" />
                                                                                <p className="font-medium">No consumers registered</p>
                                                                                <p className="text-xs mt-1">This event is not being consumed</p>
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
                    ))
                )}
            </div>
        </div>
    );
}
