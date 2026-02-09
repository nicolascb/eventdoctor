import { EmptyState, PageHeader, SearchInput, StatCard } from "@/components/shared";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import type { Consumer } from "@/types";
import { Activity, Box, Layers, Server, Workflow } from "lucide-react";
import { useMemo, useState } from "react";

interface ConsumersViewProps {
    consumers: Consumer[];
}

export function ConsumersView({ consumers }: ConsumersViewProps) {
    const [searchQuery, setSearchQuery] = useState("");

    const filteredConsumers = useMemo(() => {
        if (!searchQuery) return consumers;

        return consumers.filter(consumer =>
            consumer.group.toLowerCase().includes(searchQuery.toLowerCase()) ||
            consumer.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            consumer.topics.some(t => t.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                t.events.some(e => e.name?.toLowerCase().includes(searchQuery.toLowerCase()))
            )
        );
    }, [consumers, searchQuery]);

    const totalTopics = useMemo(() => {
        return consumers.reduce((acc, c) => acc + c.topics.length, 0);
    }, [consumers]);

    const totalEvents = useMemo(() => {
        return consumers.reduce((acc, c) =>
            acc + c.topics.reduce((sum, t) => sum + t.events.length, 0), 0
        );
    }, [consumers]);

    return (
        <div className="space-y-8 animate-in">
            {/* Header Section */}
            <div className="relative">
                <div className="flex flex-col gap-6">
                    <PageHeader
                        title="Consumers"
                        badge={`${consumers.length} services`}
                    />

                    {/* Stats Cards */}
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                        <StatCard
                            label="Consumer Groups"
                            value={consumers.length}
                            description="Active consumer services"
                            icon={<Layers className="h-6 w-6 text-primary" />}
                        />
                        <StatCard
                            label="Subscribed Topics"
                            value={totalTopics}
                            description="Total topic subscriptions"
                            icon={<Server className="h-6 w-6 text-accent" />}
                            iconClassName="bg-accent/10"
                        />
                        <StatCard
                            label="Event Handlers"
                            value={totalEvents}
                            description="Registered event handlers"
                            icon={<Activity className="h-6 w-6 text-primary" />}
                        />
                    </div>

                    {/* Search Bar */}
                    <SearchInput
                        value={searchQuery}
                        onChange={setSearchQuery}
                        placeholder="Search consumers, topics, or events..."
                        resultCount={searchQuery ? filteredConsumers.length : undefined}
                        totalCount={searchQuery ? consumers.length : undefined}
                    />
                </div>
            </div>

            {/* Consumers Grid */}
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {filteredConsumers.length === 0 ? (
                    <div className="col-span-full">
                        <EmptyState
                            title="No consumers found"
                            description="Try adjusting your search criteria"
                        />
                    </div>
                ) : (
                    filteredConsumers.map((consumer, index) => (
                        <Card
                            key={index}
                            className="event-card flex flex-col h-full hover:shadow-lg transition-all border-2"
                            style={{
                                animationDelay: `${index * 0.05}s`,
                                animationFillMode: 'backwards'
                            }}
                        >
                            <CardHeader className="pb-4">
                                <div className="flex items-start gap-3 mb-2">
                                    <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                                        <Layers className="h-5 w-5 text-primary" />
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <CardTitle className="text-lg font-semibold text-foreground mb-2 truncate">
                                            {consumer.group}
                                        </CardTitle>
                                        <CardDescription className="text-sm line-clamp-2 leading-relaxed">
                                            {consumer.description || "No description provided"}
                                        </CardDescription>
                                    </div>
                                </div>
                            </CardHeader>

                            <CardContent className="flex-1 flex flex-col gap-4 pt-0">
                                <Separator />

                                <div className="space-y-3">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-2 text-sm font-semibold text-foreground">
                                            <Workflow className="h-4 w-4 text-muted-foreground" />
                                            <span>Subscriptions</span>
                                        </div>
                                        <Badge variant="secondary" className="font-mono text-xs">
                                            {consumer.topics.length} {consumer.topics.length === 1 ? 'topic' : 'topics'}
                                        </Badge>
                                    </div>

                                    <ScrollArea className="h-[220px]">
                                        <div className="space-y-4 pr-4">
                                            {consumer.topics.map((topic, topicIndex) => (
                                                <div key={topicIndex} className="space-y-2">
                                                    <div className="flex items-center gap-2 group">
                                                        <Box className="h-3.5 w-3.5 text-primary" />
                                                        <code className="text-xs font-mono font-semibold bg-muted/50 px-2 py-1 rounded-md border border-border flex-1 truncate">
                                                            {topic.name}
                                                        </code>
                                                    </div>

                                                    <div className="flex flex-wrap gap-1.5 pl-6">
                                                        {topic.events.slice(0, 3).map((ev, eventIndex) => (
                                                            <Badge
                                                                key={eventIndex}
                                                                variant="outline"
                                                                className="text-[10px] font-medium border-primary/30 bg-primary/5 hover:bg-primary/10 transition-colors"
                                                            >
                                                                <Activity className="h-3 w-3 mr-1" />
                                                                {ev.name || String(ev)}
                                                            </Badge>
                                                        ))}
                                                        {topic.events.length > 3 && (
                                                            <Badge
                                                                variant="secondary"
                                                                className="text-[10px]"
                                                            >
                                                                +{topic.events.length - 3} more
                                                            </Badge>
                                                        )}
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </ScrollArea>
                                </div>
                            </CardContent>
                        </Card>
                    ))
                )}
            </div>
        </div>
    );
}