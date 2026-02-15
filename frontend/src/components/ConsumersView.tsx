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
import { useConsumers } from "@/hooks/useConsumers";
import type { Consumer, Topic } from "@/types";
import { Activity, ChevronLeft, ChevronRight, ExternalLink, Layers, Server, Workflow, Zap } from "lucide-react";
import { useState } from "react";

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
    const [selectedTopic, setSelectedTopic] = useState<Topic | null>(null);

    const totalConsumerGroups = pagination?.total ?? consumers.length;

    return (
        <div className="space-y-4 animate-in">
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

            {/* Search Bar */}
            <SearchInput
                value={search}
                onChange={setSearch}
                placeholder="Search consumers, topics, or events..."
                resultCount={search ? (pagination?.total ?? consumers.length) : undefined}
                totalCount={search ? undefined : undefined}
            />

            {!loading && (
                <>
                    {/* Consumers Table */}
                    <div className="rounded-lg border border-border bg-card overflow-hidden">
                        {consumers.length === 0 ? (
                            <div className="p-6">
                                <EmptyState
                                    title="No consumers found"
                                    description="Try adjusting your search criteria"
                                />
                            </div>
                        ) : (
                            <Table>
                                <TableHeader>
                                    <TableRow className="hover:bg-transparent">
                                        <TableHead className="text-[11px] font-medium">Group</TableHead>
                                        <TableHead className="text-[11px] font-medium">Service</TableHead>
                                        <TableHead className="text-[11px] font-medium text-center">Topics</TableHead>
                                        <TableHead className="text-[11px] font-medium text-center">Events</TableHead>
                                        <TableHead className="text-[11px] font-medium text-right"></TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {consumers.map((consumer) => {
                                        const consumerEventCount = consumer.topics.reduce((sum, t) => sum + t.events.length, 0);
                                        return (
                                            <TableRow
                                                key={consumer.group}
                                                className="cursor-pointer group"
                                                onClick={() => { setSelectedConsumer(consumer); setSelectedTopic(null); }}
                                            >
                                                <TableCell className="py-3">
                                                    <div className="flex flex-col gap-0.5">
                                                        <div className="flex items-center gap-2">
                                                            <Layers className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                                            <span className="font-medium text-sm">{consumer.group}</span>
                                                        </div>
                                                        {consumer.description && (
                                                            <span className="text-xs text-muted-foreground pl-5.5 line-clamp-2">
                                                                {consumer.description}
                                                            </span>
                                                        )}
                                                    </div>
                                                </TableCell>
                                                <TableCell className="py-3">
                                                    <div className="flex items-center gap-2">
                                                        <Server className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                                        <span className="text-sm">{consumer.service}</span>
                                                    </div>
                                                </TableCell>
                                                <TableCell className="py-3 text-center">
                                                    <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                                        {consumer.topics.length}
                                                    </Badge>
                                                </TableCell>
                                                <TableCell className="py-3 text-center">
                                                    <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                                        {consumerEventCount}
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
                                        );
                                    })}
                                </TableBody>
                            </Table>
                        )}
                    </div>

                    {/* Pagination Controls */}
                    {pagination && pagination.total_pages > 1 && (
                        <div className="flex items-center justify-between px-1">
                            <span className="text-xs text-muted-foreground">
                                Page {pagination.page} of {pagination.total_pages} ({pagination.total} groups)
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

                    {/* Consumer Detail Dialog */}
                    <Dialog open={!!selectedConsumer && !selectedTopic} onOpenChange={(open) => { if (!open) setSelectedConsumer(null); }}>
                        {selectedConsumer && (
                            <DialogContent className="max-w-xl">
                                <DialogHeader>
                                    <div className="flex items-center gap-2">
                                        <Layers className="h-4 w-4 text-muted-foreground" />
                                        <DialogTitle className="text-base font-semibold">
                                            {selectedConsumer.group}
                                        </DialogTitle>
                                    </div>
                                    {selectedConsumer.description && (
                                        <DialogDescription>
                                            {selectedConsumer.description}
                                        </DialogDescription>
                                    )}
                                </DialogHeader>

                                {/* Metadata rows */}
                                <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-1">
                                    <span className="text-muted-foreground text-xs font-medium">Service</span>
                                    <span className="text-xs font-medium">{selectedConsumer.service}</span>
                                    <span className="text-muted-foreground text-xs font-medium">Repository</span>
                                    {selectedConsumer.repository ? (
                                        <a
                                            href={selectedConsumer.repository}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                                        >
                                            <ExternalLink className="h-3 w-3 flex-shrink-0" />
                                            <span className="truncate">{selectedConsumer.repository}</span>
                                        </a>
                                    ) : (
                                        <span className="text-xs text-muted-foreground italic">Not specified</span>
                                    )}
                                </div>

                                {/* Topic Subscriptions */}
                                {selectedConsumer.topics.length > 0 && (
                                    <>
                                        <Separator />
                                        <div className="space-y-2">
                                            <div className="flex items-center justify-between">
                                                <h4 className="text-xs font-medium text-muted-foreground">
                                                    Topic Subscriptions
                                                </h4>
                                                <span className="text-xs font-mono text-muted-foreground">
                                                    {selectedConsumer.topics.length}
                                                </span>
                                            </div>

                                            <div className="space-y-1.5">
                                                {selectedConsumer.topics.map((topic) => (
                                                    <button
                                                        key={topic.name}
                                                        type="button"
                                                        className="w-full text-left rounded-md border border-border px-3 py-2.5 hover:bg-accent/30 transition-colors cursor-pointer"
                                                        onClick={() => setSelectedTopic(topic)}
                                                    >
                                                        <div className="flex items-center justify-between gap-3">
                                                            <code className="text-xs font-mono font-medium">{topic.name}</code>
                                                            <div className="flex items-center gap-2 flex-shrink-0">
                                                                <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                                                    {topic.events.length} {topic.events.length === 1 ? 'event' : 'events'}
                                                                </Badge>
                                                                <Zap className="h-3.5 w-3.5 text-muted-foreground" />
                                                            </div>
                                                        </div>
                                                        {topic.events.length > 0 && (
                                                            <div className="flex flex-wrap gap-1 mt-1.5">
                                                                {topic.events.slice(0, 4).map((ev) => (
                                                                    <Badge
                                                                        key={ev.name}
                                                                        variant="outline"
                                                                        className="text-[10px] font-normal font-mono"
                                                                    >
                                                                        {ev.name}
                                                                        {ev.version && <span className="ml-0.5 text-muted-foreground">v{ev.version}</span>}
                                                                    </Badge>
                                                                ))}
                                                                {topic.events.length > 4 && (
                                                                    <Badge variant="secondary" className="text-[10px] font-normal">
                                                                        +{topic.events.length - 4}
                                                                    </Badge>
                                                                )}
                                                            </div>
                                                        )}
                                                    </button>
                                                ))}
                                            </div>
                                        </div>
                                    </>
                                )}
                            </DialogContent>
                        )}
                    </Dialog>

                    {/* Topic Detail Dialog */}
                    <Dialog open={!!selectedTopic} onOpenChange={(open) => { if (!open) setSelectedTopic(null); }}>
                        {selectedTopic && selectedConsumer && (
                            <DialogContent className="max-w-xl">
                                <DialogHeader>
                                    <div className="flex items-center gap-2">
                                        <Zap className="h-4 w-4 text-muted-foreground" />
                                        <DialogTitle className="font-mono text-base flex items-center gap-2">
                                            <span className="text-muted-foreground font-normal">{selectedConsumer.group}</span>
                                            <span className="text-muted-foreground/30">/</span>
                                            <span className="font-semibold">{selectedTopic.name}</span>
                                        </DialogTitle>
                                    </div>
                                    <DialogDescription>
                                        Events consumed from this topic by {selectedConsumer.service}
                                    </DialogDescription>
                                </DialogHeader>

                                <div className="space-y-1.5">
                                    {selectedTopic.events.map((ev) => (
                                        <div
                                            key={`${ev.name}-${ev.version || ''}`}
                                            className="flex items-center justify-between rounded-md border border-border px-3 py-2.5"
                                        >
                                            <span className="font-mono text-xs font-medium">{ev.name}</span>
                                            {ev.version ? (
                                                <Badge variant="outline" className="font-mono text-[10px] font-normal">
                                                    v{ev.version}
                                                </Badge>
                                            ) : (
                                                <span className="text-xs text-muted-foreground">—</span>
                                            )}
                                        </div>
                                    ))}
                                </div>

                                <Separator />
                                <div>
                                    <Button variant="outline" size="sm" className="text-xs" onClick={() => setSelectedTopic(null)}>
                                        ← Back to Consumer
                                    </Button>
                                </div>
                            </DialogContent>
                        )}
                    </Dialog>
                </>
            )}
        </div>
    );
}
