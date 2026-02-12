import { EmptyState, PageHeader, SearchInput, StatCard } from "@/components/shared";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import type { Producer } from "@/types";
import { Activity, Box, Database, FileText, Zap } from "lucide-react";
import { useMemo, useState } from "react";

interface ProducersViewProps {
    producers: Producer[];
}

export function ProducersView({ producers }: ProducersViewProps) {
    const [searchQuery, setSearchQuery] = useState("");

    const filteredProducers = useMemo(() => {
        if (!searchQuery) return producers;

        return producers.filter(producer =>
            producer.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            producer.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            producer.topic.toLowerCase().includes(searchQuery.toLowerCase()) ||
            producer.events.some(e => e.name?.toLowerCase().includes(searchQuery.toLowerCase()))
        );
    }, [producers, searchQuery]);

    const totalEvents = useMemo(() => {
        return producers.reduce((acc, p) => acc + p.events.length, 0);
    }, [producers]);

    const activeProducers = useMemo(() => {
        return producers.filter(p => p.writes).length;
    }, [producers]);

    return (
        <div className="space-y-8 animate-in">
            {/* Header Section */}
            <div className="relative">
                <div className="flex flex-col gap-6">
                    <PageHeader
                        title="Producers"
                        badge={`${producers.length} services`}
                    />

                    {/* Stats Cards */}
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                        <StatCard
                            label="Total Producers"
                            value={producers.length}
                            description="Event publishing services"
                            icon={<Database className="h-6 w-6 text-primary" />}
                        />
                        <StatCard
                            label="Active Writers"
                            value={activeProducers}
                            description="Currently writing events"
                            icon={<Activity className="h-6 w-6 text-accent" />}
                            iconClassName="bg-accent/10"
                        />
                        <StatCard
                            label="Event Types"
                            value={totalEvents}
                            description="Published event types"
                            icon={<Zap className="h-6 w-6 text-primary" />}
                        />
                    </div>

                    {/* Search Bar */}
                    <SearchInput
                        value={searchQuery}
                        onChange={setSearchQuery}
                        placeholder="Search producers, topics, or events..."
                        resultCount={searchQuery ? filteredProducers.length : undefined}
                        totalCount={searchQuery ? producers.length : undefined}
                    />
                </div>
            </div>

            {/* Producers Grid */}
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {filteredProducers.length === 0 ? (
                    <div className="col-span-full">
                        <EmptyState
                            title="No producers found"
                            description="Try adjusting your search criteria"
                        />
                    </div>
                ) : (
                    filteredProducers.map((producer, index) => (
                        <Card
                            key={`${producer.name}-${producer.topic}`}
                            className="event-card flex flex-col h-full hover:shadow-lg transition-all border-2"
                            style={{
                                animationDelay: `${index * 0.05}s`,
                                animationFillMode: 'backwards'
                            }}
                        >
                            <CardHeader className="pb-4">
                                <div className="flex items-start justify-between gap-3 mb-2">
                                    <div className="flex items-start gap-3 flex-1 min-w-0">
                                        <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                                            <Database className="h-5 w-5 text-primary" />
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <CardTitle className="text-lg font-semibold text-foreground mb-1 truncate">
                                                {producer.name}
                                            </CardTitle>
                                            <div className="flex items-center gap-2 mt-2">
                                                <Box className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                                <code className="text-xs font-mono bg-muted/50 px-2 py-1 rounded-md border border-border truncate flex-1">
                                                    {producer.topic}
                                                </code>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="flex flex-col gap-1.5 flex-shrink-0">
                                        {producer.owner && (
                                            <Badge variant="secondary" className="text-[10px] gap-1">
                                                <Database className="h-3 w-3" />
                                                Owner
                                            </Badge>
                                        )}
                                        {producer.writes && (
                                            <Badge variant="default" className="text-[10px] gap-1">
                                                <Activity className="h-3 w-3" />
                                                Active
                                            </Badge>
                                        )}
                                    </div>
                                </div>
                            </CardHeader>

                            <CardContent className="flex-1 flex flex-col gap-4 pt-0">
                                {producer.description && (
                                    <>
                                        <Separator />
                                        <div className="space-y-2">
                                            <div className="flex items-center gap-2">
                                                <FileText className="h-4 w-4 text-muted-foreground" />
                                                <span className="text-sm font-semibold text-foreground">Description</span>
                                            </div>
                                            <p className="text-sm text-muted-foreground leading-relaxed pl-6">
                                                {producer.description}
                                            </p>
                                        </div>
                                    </>
                                )}

                                {(producer.events && producer.events.length > 0) && (
                                    <>
                                        <Separator />
                                        <div className="space-y-3">
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-2">
                                                    <Zap className="h-4 w-4 text-muted-foreground" />
                                                    <span className="text-sm font-semibold text-foreground">
                                                        Events
                                                    </span>
                                                </div>
                                                <Badge variant="secondary" className="font-mono text-xs">
                                                    {producer.events.length}
                                                </Badge>
                                            </div>
                                            <div className="flex flex-wrap gap-2 pl-6">
                                                {producer.events.slice(0, 4).map((ev, eventIndex) => (
                                                    <Badge
                                                        key={`${producer.name}-${producer.topic}-${ev.name || eventIndex}${ev.version ? `-${ev.version}` : ''}`}
                                                        variant="outline"
                                                        className="text-[10px] font-medium border-primary/30 bg-primary/5 hover:bg-primary/10 transition-colors"
                                                    >
                                                        {ev.name || String(ev)}
                                                    </Badge>
                                                ))}
                                                {producer.events.length > 4 && (
                                                    <Badge
                                                        variant="secondary"
                                                        className="text-[10px]"
                                                    >
                                                        +{producer.events.length - 4} more
                                                    </Badge>
                                                )}
                                            </div>
                                        </div>
                                    </>
                                )}
                            </CardContent>
                        </Card>
                    ))
                )}
            </div>
        </div>
    );
}