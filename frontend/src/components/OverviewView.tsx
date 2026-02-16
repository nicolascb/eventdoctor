import { EventFlowDiagram } from "@/components/EventFlowDiagram";
import { StatCard } from "@/components/shared";
import { Button } from "@/components/ui/button";
import { useTopics } from "@/hooks/useTopics";
import type { OverviewResponse, TopicView } from "@/types";
import { Check, ChevronLeft, ChevronRight, Database, Layers, Network, Zap } from "lucide-react";
import { useState } from "react";

interface OverviewViewProps {
    overview: OverviewResponse;
}

export function OverviewView({ overview }: OverviewViewProps) {
    const { topics, pagination, page, setPage, loading: topicsLoading } = useTopics();
    const [selectedTopics, setSelectedTopics] = useState<TopicView[]>([]);

    const toggleTopic = (topic: TopicView) => {
        setSelectedTopics(prev =>
            prev.find(t => t.topic === topic.topic)
                ? prev.filter(t => t.topic !== topic.topic)
                : [...prev, topic]
        );
    };

    return (
        <div className="space-y-6 animate-in">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <StatCard
                    label="Total Topics"
                    value={overview.total_topics}
                    description="Message channels"
                    icon={<Layers className="h-4 w-4 text-muted-foreground" />}
                />
                <StatCard
                    label="Total Events"
                    value={overview.total_events}
                    description="Event types defined"
                    icon={<Zap className="h-4 w-4 text-muted-foreground" />}
                />
                <StatCard
                    label="Producers"
                    value={overview.total_producers}
                    description="Publishing services"
                    icon={<Database className="h-4 w-4 text-muted-foreground" />}
                />
                <StatCard
                    label="Consumers"
                    value={overview.total_consumers}
                    description="Consuming services"
                    icon={<Network className="h-4 w-4 text-muted-foreground" />}
                />
            </div>

            <div className="grid grid-cols-1 xl:grid-cols-4 gap-6 text-foreground">
                {/* Topic Selector */}
                <div className="xl:col-span-1 space-y-4">
                    <div className="flex items-center justify-between">
                        <h3 className="text-sm font-semibold flex items-center gap-2">
                            <Layers className="h-4 w-4 text-muted-foreground" />
                            Select Topics
                        </h3>
                        {selectedTopics.length > 0 && (
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-7 text-[10px]"
                                onClick={() => setSelectedTopics([])}
                            >
                                Clear ({selectedTopics.length})
                            </Button>
                        )}
                    </div>

                    <div className="rounded-lg border border-border bg-card overflow-hidden">
                        <div className="max-h-[500px] overflow-y-auto p-1 space-y-1">
                            {topicsLoading && topics.length === 0 ? (
                                <div className="p-4 text-center text-xs text-muted-foreground">Loading...</div>
                            ) : (
                                topics.map(topic => {
                                    const isSelected = selectedTopics.some(t => t.topic === topic.topic);
                                    return (
                                        <button
                                            key={topic.topic}
                                            onClick={() => toggleTopic(topic)}
                                            className={`w-full flex items-center justify-between p-2 rounded-md transition-colors text-left ${isSelected
                                                ? "bg-primary/10 text-primary border border-primary/20"
                                                : "hover:bg-accent text-sm"
                                                }`}
                                        >
                                            <div className="flex flex-col gap-0.5 max-w-[180px]">
                                                <span className="font-mono text-xs truncate font-medium">{topic.topic}</span>
                                                <span className="text-[10px] text-muted-foreground">
                                                    {topic.events.length} event{topic.events.length !== 1 ? 's' : ''}
                                                </span>
                                            </div>
                                            {isSelected && <Check className="h-3 w-3" />}
                                        </button>
                                    );
                                })
                            )}
                        </div>

                        {/* Selector Pagination */}
                        {pagination && pagination.total_pages > 1 && (
                            <div className="p-2 border-t border-border flex items-center justify-between bg-muted/30">
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    className="h-7 w-7 p-0"
                                    disabled={page <= 1}
                                    onClick={() => setPage(page - 1)}
                                >
                                    <ChevronLeft className="h-4 w-4" />
                                </Button>
                                <span className="text-[10px] text-muted-foreground font-medium">
                                    {page} / {pagination.total_pages}
                                </span>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    className="h-7 w-7 p-0"
                                    disabled={page >= pagination.total_pages}
                                    onClick={() => setPage(page + 1)}
                                >
                                    <ChevronRight className="h-4 w-4" />
                                </Button>
                            </div>
                        )}
                    </div>
                </div>

                {/* Diagram */}
                <div className="xl:col-span-3">
                    {selectedTopics.length === 0 ? (
                        <div className="h-full min-h-[400px] rounded-xl border border-dashed border-border flex flex-col items-center justify-center p-12 text-center bg-muted/5">
                            <div className="rounded-full bg-muted/20 p-4 mb-4">
                                <Zap className="h-8 w-8 text-muted-foreground/40" />
                            </div>
                            <h3 className="text-lg font-medium text-foreground mb-2">Visualizer</h3>
                            <p className="max-w-[300px] text-sm text-muted-foreground">
                                Select topics from the list on the left to visualize the event flow between producers and consumers.
                            </p>
                        </div>
                    ) : (
                        <EventFlowDiagram topics={selectedTopics} hideSelector />
                    )}
                </div>
            </div>
        </div>
    );
}
