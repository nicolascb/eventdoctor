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
import { useProducers } from "@/hooks/useProducers";
import type { Event, GroupedProducer, Producer, ProducerTopic } from "@/types";
import { Database, ExternalLink, FileJson, Layers, Server, Zap } from "lucide-react";
import { useMemo, useState } from "react";

function groupProducersByService(producers: Producer[]): GroupedProducer[] {
    const map = new Map<string, GroupedProducer>();

    for (const p of producers) {
        let group = map.get(p.service);
        if (!group) {
            group = {
                service: p.service,
                repository: p.repository,
                topics: [],
            };
            map.set(p.service, group);
        }

        group.topics.push({
            topic: p.topic,
            description: p.description,
            owner: p.owner,
            writes: p.writes,
            events: p.events,
        });
    }

    return Array.from(map.values());
}

export function ProducersView() {
    const { producers, loading } = useProducers();
    const [searchQuery, setSearchQuery] = useState("");
    const [selectedGroup, setSelectedGroup] = useState<GroupedProducer | null>(null);
    const [selectedTopic, setSelectedTopic] = useState<ProducerTopic | null>(null);
    const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);

    const grouped = useMemo(() => groupProducersByService(producers), [producers]);

    const filteredGroups = useMemo(() => {
        if (!searchQuery) return grouped;
        const q = searchQuery.toLowerCase();
        return grouped
            .map(group => {
                const matchesService = group.service.toLowerCase().includes(q);
                if (matchesService) return group;

                const filteredTopics = group.topics.filter(t =>
                    t.topic.toLowerCase().includes(q) ||
                    t.description?.toLowerCase().includes(q) ||
                    t.events.some(e => e.name?.toLowerCase().includes(q))
                );

                if (filteredTopics.length === 0) return null;
                return { ...group, topics: filteredTopics };
            })
            .filter((g): g is GroupedProducer => g !== null);
    }, [grouped, searchQuery]);

    const totalEvents = useMemo(() => {
        return producers.reduce((acc, p) => acc + p.events.length, 0);
    }, [producers]);

    const uniqueServices = useMemo(() => grouped.length, [grouped]);

    const totalTopics = useMemo(() => {
        return new Set(producers.map(p => p.topic)).size;
    }, [producers]);

    const totalEventsForGroup = (group: GroupedProducer) =>
        group.topics.reduce((acc, t) => acc + t.events.length, 0);

    if (loading) return null;

    return (
        <div className="space-y-6 animate-in">
            {/* Stats Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <StatCard
                    label="Services"
                    value={uniqueServices}
                    description="Unique publishing services"
                    icon={<Database className="h-4 w-4 text-muted-foreground" />}
                />
                <StatCard
                    label="Topics"
                    value={totalTopics}
                    description="Topics with active producers"
                    icon={<Layers className="h-4 w-4 text-muted-foreground" />}
                />
                <StatCard
                    label="Event Types"
                    value={totalEvents}
                    description="Published event types"
                    icon={<Zap className="h-4 w-4 text-muted-foreground" />}
                />
            </div>

            {/* Search Bar */}
            <SearchInput
                value={searchQuery}
                onChange={setSearchQuery}
                placeholder="Search producers, topics, or events..."
                resultCount={searchQuery ? filteredGroups.length : undefined}
                totalCount={searchQuery ? grouped.length : undefined}
            />

            {/* Producers Table */}
            <div className="rounded-lg border border-border bg-card overflow-hidden">
                {filteredGroups.length === 0 ? (
                    <div className="p-6">
                        <EmptyState
                            title="No producers found"
                            description="Try adjusting your search criteria"
                        />
                    </div>
                ) : (
                    <Table>
                        <TableHeader>
                            <TableRow className="hover:bg-transparent">
                                <TableHead className="text-[11px] font-medium">Service</TableHead>
                                <TableHead className="text-[11px] font-medium">Topics</TableHead>
                                <TableHead className="text-[11px] font-medium text-center">Events</TableHead>
                                <TableHead className="text-[11px] font-medium text-right"></TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {filteredGroups.map((group) => (
                                <TableRow
                                    key={group.service}
                                    className="cursor-pointer group"
                                    onClick={() => { setSelectedGroup(group); setSelectedTopic(null); setSelectedEvent(null); }}
                                >
                                    <TableCell className="py-3">
                                        <div className="flex items-center gap-2">
                                            <Server className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                            <span className="font-medium text-sm">{group.service}</span>
                                            {group.topics.some(t => t.owner) && (
                                                <Badge variant="outline" className="text-[10px] font-normal">
                                                    Owner
                                                </Badge>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-3">
                                        <div className="flex items-center gap-1.5">
                                            {group.topics.slice(0, 2).map((t) => (
                                                <code
                                                    key={t.topic}
                                                    className="text-xs font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded truncate max-w-[160px]"
                                                >
                                                    {t.topic}
                                                </code>
                                            ))}
                                            {group.topics.length > 2 && (
                                                <Badge variant="secondary" className="text-[10px] font-normal tabular-nums flex-shrink-0">
                                                    +{group.topics.length - 2}
                                                </Badge>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-3 text-center">
                                        <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                            {totalEventsForGroup(group)}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="py-3 text-right">
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            className="h-7 w-7 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                                        >
                                            <FileJson className="h-3.5 w-3.5" />
                                        </Button>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                )}
            </div>

            {/* Producer Detail Dialog (service level) */}
            <Dialog open={!!selectedGroup && !selectedTopic && !selectedEvent} onOpenChange={(open) => { if (!open) setSelectedGroup(null); }}>
                {selectedGroup && (
                    <DialogContent className="max-w-xl">
                        <DialogHeader>
                            <div className="flex items-center gap-2">
                                <Server className="h-4 w-4 text-muted-foreground" />
                                <DialogTitle className="text-base font-semibold">
                                    {selectedGroup.service}
                                </DialogTitle>
                                {selectedGroup.topics.some(t => t.owner) && (
                                    <Badge variant="outline" className="text-[10px] font-normal">Owner</Badge>
                                )}
                            </div>
                            <DialogDescription>
                                Produces to {selectedGroup.topics.length} {selectedGroup.topics.length === 1 ? "topic" : "topics"} with {totalEventsForGroup(selectedGroup)} event {totalEventsForGroup(selectedGroup) === 1 ? "type" : "types"}
                            </DialogDescription>
                        </DialogHeader>

                        {/* Metadata rows */}
                        <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-1">
                            <span className="text-muted-foreground text-xs font-medium">Repository</span>
                            {selectedGroup.repository ? (
                                <a
                                    href={selectedGroup.repository}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                                >
                                    <ExternalLink className="h-3 w-3 flex-shrink-0" />
                                    <span className="truncate">{selectedGroup.repository}</span>
                                </a>
                            ) : (
                                <span className="text-xs text-muted-foreground italic">Not specified</span>
                            )}
                        </div>

                        {/* Topics list */}
                        <Separator />
                        <div className="space-y-2">
                            <div className="flex items-center justify-between">
                                <h4 className="text-xs font-medium text-muted-foreground">
                                    Topics
                                </h4>
                                <span className="text-xs font-mono text-muted-foreground">
                                    {selectedGroup.topics.length}
                                </span>
                            </div>

                            <div className="space-y-1.5">
                                {selectedGroup.topics.map((topic) => (
                                    <button
                                        key={topic.topic}
                                        type="button"
                                        className="w-full text-left rounded-md border border-border px-3 py-2.5 hover:bg-accent/30 transition-colors cursor-pointer"
                                        onClick={() => setSelectedTopic(topic)}
                                    >
                                        <div className="flex items-center justify-between gap-3">
                                            <div className="flex items-center gap-2 min-w-0">
                                                <Layers className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                                <code className="font-mono text-xs font-medium">{topic.topic}</code>
                                                {topic.owner && (
                                                    <Badge variant="outline" className="text-[10px] font-normal flex-shrink-0">
                                                        Owner
                                                    </Badge>
                                                )}
                                            </div>
                                            <Badge variant="secondary" className="text-[10px] font-normal tabular-nums flex-shrink-0">
                                                {topic.events.length} {topic.events.length === 1 ? "event" : "events"}
                                            </Badge>
                                        </div>
                                        {topic.description && (
                                            <p className="text-xs text-muted-foreground mt-1 pl-5.5 line-clamp-2">
                                                {topic.description}
                                            </p>
                                        )}
                                    </button>
                                ))}
                            </div>
                        </div>
                    </DialogContent>
                )}
            </Dialog>

            {/* Topic Detail Dialog (events within a topic) */}
            <Dialog open={!!selectedTopic && !selectedEvent} onOpenChange={(open) => { if (!open) setSelectedTopic(null); }}>
                {selectedTopic && selectedGroup && (
                    <DialogContent className="max-w-xl">
                        <DialogHeader>
                            <div className="flex items-center gap-2">
                                <Layers className="h-4 w-4 text-muted-foreground" />
                                <DialogTitle className="text-base font-semibold font-mono">
                                    {selectedTopic.topic}
                                </DialogTitle>
                                {selectedTopic.owner && (
                                    <Badge variant="outline" className="text-[10px] font-normal">Owner</Badge>
                                )}
                            </div>
                            {selectedTopic.description && (
                                <DialogDescription>
                                    {selectedTopic.description}
                                </DialogDescription>
                            )}
                        </DialogHeader>

                        {/* Metadata rows */}
                        <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-1">
                            <span className="text-muted-foreground text-xs font-medium">Service</span>
                            <span className="text-xs font-medium">{selectedGroup.service}</span>
                        </div>

                        {/* Events */}
                        {selectedTopic.events.length > 0 && (
                            <>
                                <Separator />
                                <div className="space-y-2">
                                    <div className="flex items-center justify-between">
                                        <h4 className="text-xs font-medium text-muted-foreground">
                                            Published Events
                                        </h4>
                                        <span className="text-xs font-mono text-muted-foreground">
                                            {selectedTopic.events.length}
                                        </span>
                                    </div>

                                    <div className="space-y-1.5">
                                        {selectedTopic.events.map((ev) => (
                                            <button
                                                key={`${ev.name}-${ev.version || ''}`}
                                                type="button"
                                                className="w-full text-left rounded-md border border-border px-3 py-2.5 hover:bg-accent/30 transition-colors cursor-pointer"
                                                onClick={() => setSelectedEvent(ev)}
                                            >
                                                <div className="flex items-center justify-between gap-3">
                                                    <div className="flex items-center gap-2 min-w-0">
                                                        <span className="font-mono text-xs font-medium">{ev.name}</span>
                                                        {ev.version && (
                                                            <Badge variant="outline" className="font-mono text-[10px] font-normal flex-shrink-0">
                                                                v{ev.version}
                                                            </Badge>
                                                        )}
                                                    </div>
                                                    <div className="flex items-center gap-2 flex-shrink-0">
                                                        {ev.schema_url && (
                                                            <div onClick={(e) => e.stopPropagation()}>
                                                                <a
                                                                    href={ev.schema_url}
                                                                    target="_blank"
                                                                    rel="noopener noreferrer"
                                                                    className="text-muted-foreground hover:text-foreground transition-colors"
                                                                >
                                                                    <ExternalLink className="h-3 w-3" />
                                                                </a>
                                                            </div>
                                                        )}
                                                        <FileJson className="h-3.5 w-3.5 text-muted-foreground" />
                                                    </div>
                                                </div>
                                                {ev.description && (
                                                    <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                                                        {ev.description}
                                                    </p>
                                                )}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            </>
                        )}

                        <Separator />
                        <div>
                            <Button variant="outline" size="sm" className="text-xs" onClick={() => setSelectedTopic(null)}>
                                ← Back to Service
                            </Button>
                        </div>
                    </DialogContent>
                )}
            </Dialog>

            {/* Event Detail Dialog */}
            <Dialog open={!!selectedEvent} onOpenChange={(open) => { if (!open) setSelectedEvent(null); }}>
                {selectedEvent && selectedTopic && (
                    <DialogContent className="max-w-3xl">
                        <DialogHeader>
                            <div className="flex items-center gap-2 mb-1">
                                <Zap className="h-4 w-4 text-muted-foreground" />
                                <DialogTitle className="font-mono text-base flex items-center gap-2">
                                    <span className="text-muted-foreground font-normal">{selectedTopic.topic}</span>
                                    <span className="text-muted-foreground/30">/</span>
                                    <span className="font-semibold">{selectedEvent.name}</span>
                                    {selectedEvent.version && (
                                        <Badge variant="outline" className="font-mono text-[10px] font-normal ml-1">
                                            v{selectedEvent.version}
                                        </Badge>
                                    )}
                                </DialogTitle>
                            </div>
                            {selectedEvent.description && (
                                <DialogDescription className="text-sm">
                                    {selectedEvent.description}
                                </DialogDescription>
                            )}
                        </DialogHeader>

                        {/* Schema metadata */}
                        <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-xs py-1">
                            <span className="text-muted-foreground font-medium">Headers</span>
                            <span className="font-mono">{selectedEvent.headers?.length ?? 0}</span>
                            {selectedEvent.properties && (
                                <>
                                    <span className="text-muted-foreground font-medium">Properties</span>
                                    <span className="font-mono">{Object.keys(selectedEvent.properties).length}</span>
                                </>
                            )}
                            {selectedEvent.schema_url && (
                                <>
                                    <span className="text-muted-foreground font-medium">Schema</span>
                                    <a
                                        href={selectedEvent.schema_url}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors"
                                    >
                                        <ExternalLink className="h-3 w-3" />
                                        <span className="truncate">{selectedEvent.schema_url}</span>
                                    </a>
                                </>
                            )}
                        </div>

                        {/* Headers */}
                        {selectedEvent.headers && selectedEvent.headers.length > 0 && (
                            <>
                                <Separator />
                                <div className="space-y-2">
                                    <h4 className="text-xs font-medium text-muted-foreground">Headers</h4>
                                    <div className="rounded-md border border-border overflow-hidden">
                                        <Table>
                                            <TableHeader>
                                                <TableRow className="hover:bg-transparent">
                                                    <TableHead className="text-[11px] font-medium">Name</TableHead>
                                                    <TableHead className="text-[11px] font-medium">Description</TableHead>
                                                </TableRow>
                                            </TableHeader>
                                            <TableBody>
                                                {selectedEvent.headers.map((header) => (
                                                    <TableRow key={header.name}>
                                                        <TableCell className="font-mono text-xs py-2">{header.name}</TableCell>
                                                        <TableCell className="text-xs text-muted-foreground py-2">{header.description || '—'}</TableCell>
                                                    </TableRow>
                                                ))}
                                            </TableBody>
                                        </Table>
                                    </div>
                                </div>
                            </>
                        )}

                        {/* Properties */}
                        {selectedEvent.properties && (
                            <>
                                <Separator />
                                <div className="space-y-2">
                                    <h4 className="text-xs font-medium text-muted-foreground">Properties Schema</h4>
                                    <div className="bg-muted p-3 rounded-md border border-border overflow-x-auto">
                                        <pre className="text-xs font-mono leading-relaxed">
                                            {JSON.stringify(selectedEvent.properties, null, 2)}
                                        </pre>
                                    </div>
                                </div>
                            </>
                        )}

                        <Separator />
                        <div>
                            <Button variant="outline" size="sm" className="text-xs" onClick={() => setSelectedEvent(null)}>
                                ← Back to Topic
                            </Button>
                        </div>
                    </DialogContent>
                )}
            </Dialog>
        </div>
    );
}