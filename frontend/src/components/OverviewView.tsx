import { PageHeader, StatCard } from "@/components/shared";
import type { OverviewResponse } from "@/types";
import { Database, Layers, Network, Zap } from "lucide-react";

interface OverviewViewProps {
    overview: OverviewResponse;
}

export function OverviewView({ overview }: OverviewViewProps) {
    return (
        <div className="space-y-8 animate-in">
            <div className="relative">
                <div className="flex flex-col gap-6">
                    <PageHeader
                        title="Overview"
                        badge={`${overview.total_topics} topics · ${overview.total_events} events`}
                    />

                    <div className="grid grid-cols-1 sm:grid-cols-4 gap-4">
                        <StatCard
                            label="Topics"
                            value={overview.total_topics}
                            description="Message channels"
                            icon={<Layers className="h-6 w-6 text-primary" />}
                        />
                        <StatCard
                            label="Events"
                            value={overview.total_events}
                            description="Event types"
                            icon={<Zap className="h-6 w-6 text-accent" />}
                            iconClassName="bg-accent/10"
                        />
                        <StatCard
                            label="Producers"
                            value={overview.total_producers}
                            description="Publishing services"
                            icon={<Database className="h-6 w-6 text-primary" />}
                        />
                        <StatCard
                            label="Consumers"
                            value={overview.total_consumers}
                            description="Consuming services"
                            icon={<Network className="h-6 w-6 text-accent" />}
                            iconClassName="bg-accent/10"
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
