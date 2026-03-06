import { EventFlowDiagram } from "@/components/EventFlowDiagram";
import { StatCard } from "@/components/shared";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useTopics } from "@/hooks/useTopics";
import type { OverviewResponse, TopicView } from "@/types";
import { Check, ChevronLeft, ChevronRight, Database, Layers, Network, Search, X, Zap } from "lucide-react";
import { useState, useMemo } from "react";

interface OverviewViewProps {
    overview: OverviewResponse;
}

export function OverviewView({ overview }: OverviewViewProps) {
    const { topics, pagination, page, setPage, loading: topicsLoading } = useTopics();
    const [selectedTopics, setSelectedTopics] = useState<TopicView[]>([]);
    const [searchQuery, setSearchQuery] = useState("");

    const filteredTopics = useMemo(() => {
        if (!searchQuery) return topics;
        const lower = searchQuery.toLowerCase();
        return topics.filter(t => t.topic.toLowerCase().includes(lower));
    }, [topics, searchQuery]);

    const toggleTopic = (topic: TopicView) => {
        setSelectedTopics(prev =>
            prev.find(t => t.topic === topic.topic)
                ? prev.filter(t => t.topic !== topic.topic)
                : [...prev, topic]
        );
    };

    const toggleAllVisible = () => {
        const allVisibleSelected = filteredTopics.every(t => selectedTopics.some(st => st.topic === t.topic));

        if (allVisibleSelected) {
            setSelectedTopics(prev => prev.filter(st => !filteredTopics.some(ft => ft.topic === st.topic)));
        } else {
            setSelectedTopics(prev => {
                const newSelections = [...prev];
                filteredTopics.forEach(ft => {
                    if (!newSelections.some(st => st.topic === ft.topic)) {
                        newSelections.push(ft);
                    }
                });
                return newSelections;
            });
        }
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
                            Topics
                        </h3>
                        {selectedTopics.length > 0 && (
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-7 text-[10px]"
                                onClick={() => setSelectedTopics([])}
                            >
                                Clear All ({selectedTopics.length})
                            </Button>
                        )}
                    </div>

                    <div className="rounded-lg border border-border bg-card overflow-hidden flex flex-col shadow-sm">
                        <div className="p-2.5 border-b border-border bg-muted/20">
                            <div className="relative">
                                <Search className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground/70" />
                                <Input
                                    type="text"
                                    placeholder="Filter topics..."
                                    className="h-8 w-full pl-8 pr-8 text-xs bg-background shadow-none border-border/60 focus-visible:ring-1 focus-visible:ring-primary/50"
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                />
                                {searchQuery && (
                                    <button
                                        onClick={() => setSearchQuery('')}
                                        className="absolute right-1.5 top-1.5 p-1 rounded-sm text-muted-foreground/70 hover:text-foreground hover:bg-muted/50 transition-colors"
                                    >
                                        <X className="h-3.5 w-3.5" />
                                    </button>
                                )}
                            </div>
                            <div className="flex justify-between items-center mt-3 px-1">
                                <span className="text-[10px] font-medium text-muted-foreground">
                                    {filteredTopics.length} visible
                                </span>
                                {filteredTopics.length > 0 && (
                                    <button
                                        onClick={toggleAllVisible}
                                        className="flex items-center gap-1.5 text-[10px] font-medium text-muted-foreground hover:text-foreground transition-colors group outline-none"
                                    >
                                        {filteredTopics.every(t => selectedTopics.some(st => st.topic === t.topic)) ? (
                                            <>
                                                <span className="flex items-center justify-center h-3.5 w-3.5 rounded-[4px] bg-primary text-primary-foreground">
                                                    <Check className="h-2.5 w-2.5" />
                                                </span>
                                                Deselect All
                                            </>
                                        ) : (
                                            <>
                                                <span className="flex items-center justify-center h-3.5 w-3.5 rounded-[4px] border border-muted-foreground/30 text-transparent group-hover:border-primary/50">
                                                    <Check className="h-2.5 w-2.5" />
                                                </span>
                                                Select All
                                            </>
                                        )}
                                    </button>
                                )}
                            </div>
                        </div>

                        <div className="max-h-[500px] overflow-y-auto p-1.5 space-y-0.5">
                            {topicsLoading && topics.length === 0 ? (
                                <div className="p-6 text-center text-xs text-muted-foreground">Loading...</div>
                            ) : filteredTopics.length === 0 ? (
                                <div className="p-8 text-center text-xs text-muted-foreground flex flex-col items-center justify-center gap-3">
                                    <Search className="h-6 w-6 text-muted-foreground/30" />
                                    <p>No topics matching <br /> <span className="font-semibold text-foreground/70">"{searchQuery}"</span></p>
                                </div>
                            ) : (
                                filteredTopics.map(topic => {
                                    const isSelected = selectedTopics.some(t => t.topic === topic.topic);
                                    return (
                                        <button
                                            key={topic.topic}
                                            onClick={() => toggleTopic(topic)}
                                            className={`w-full flex items-center gap-3 p-2 rounded-md transition-all text-left outline-none focus-visible:ring-1 focus-visible:ring-ring ${isSelected
                                                ? "bg-primary/5 hover:bg-primary/10"
                                                : "hover:bg-accent/70"
                                                }`}
                                        >
                                            <div className={`flex shrink-0 items-center justify-center h-4 w-4 rounded-[4px] border transition-colors ${isSelected
                                                ? "bg-primary border-primary text-primary-foreground"
                                                : "border-input bg-transparent group-hover:border-primary/40"
                                                }`}>
                                                {isSelected && <Check className="h-3 w-3" />}
                                            </div>

                                            <div className="flex flex-col min-w-0 flex-1 gap-0.5">
                                                <span className={`font-mono text-xs truncate transition-colors ${isSelected ? 'font-semibold text-primary' : 'font-medium text-foreground'}`}>
                                                    {topic.topic}
                                                </span>
                                                <span className="text-[10px] text-muted-foreground flex items-center gap-1.5">
                                                    <span className={`inline-block w-1.5 h-1.5 rounded-full ${topic.events.length > 0 ? "bg-emerald-500/60" : "bg-muted-foreground/30"}`} />
                                                    {topic.events.length} event{topic.events.length !== 1 ? 's' : ''}
                                                </span>
                                            </div>
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
                    <EventFlowDiagram topics={selectedTopics} hideSelector />
                </div>
            </div>
        </div>
    );
}
