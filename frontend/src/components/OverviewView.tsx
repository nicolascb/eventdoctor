import { EventFlowDiagram } from "@/components/EventFlowDiagram";
import { StatCard } from "@/components/shared";
import type { Consumer, OverviewResponse, Producer, TopicWithEvents } from "@/types";
import { Database, Layers, Network, Zap } from "lucide-react";

interface OverviewViewProps {
    overview: OverviewResponse;
    topics: TopicWithEvents[];
    producers: Producer[];
    consumers: Consumer[];
}

export function OverviewView({ overview, topics, producers, consumers }: OverviewViewProps) {
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

            <EventFlowDiagram topics={topics} producers={producers} consumers={consumers} />
        </div>
    );
}
