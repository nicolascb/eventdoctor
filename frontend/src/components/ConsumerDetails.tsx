import { Separator } from "@/components/ui/separator";
import type { Consumer, ConsumerEvent } from "@/types";
import { ExternalLink, Layers, Server } from "lucide-react";
import { useState } from "react";
import { DetailPane, EventCard, MasterDetailLayout, TopicMenu, type TopicMenuItem } from "./shared/MasterDetail";
import { ServiceDetailsDialog } from "./ServiceDetailsDialog";

interface ConsumerDetailsProps {
    consumer?: Consumer | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onEventClick?: (ev: ConsumerEvent) => void;
}

export function ConsumerDetails({ consumer, open, onOpenChange, onEventClick }: ConsumerDetailsProps) {
    const [selectedTopicNameState, setSelectedTopicNameState] = useState<string | null>(null);

    if (!open || !consumer) return null;

    const hasTopics = consumer.topics.length > 0;
    const isValidSelection =
        selectedTopicNameState !== null && consumer.topics.some((t) => t.name === selectedTopicNameState);
    const selectedTopicName = isValidSelection ? selectedTopicNameState : hasTopics ? consumer.topics[0].name : null;

    const selectedTopic = consumer.topics.find((t) => t.name === selectedTopicName);

    const topicMenuItems: TopicMenuItem[] = consumer.topics.map((t) => ({
        id: t.name,
        name: t.name,
        eventCount: t.events.length,
    }));

    return (
        <ServiceDetailsDialog
            open={open}
            onOpenChange={onOpenChange}
            icon={<Layers className="h-6 w-6 text-muted-foreground" />}
            title={consumer.group}
            description={consumer.description}
            className="sm:max-w-4xl" // Expanded width for Master-Detail layout
        >
            {/* Metadata rows */}
            <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm px-1">
                <span className="text-muted-foreground text-xs font-semibold uppercase tracking-wider">Service</span>
                <span className="text-xs font-medium text-foreground flex items-center gap-1.5">
                    <Server className="h-3.5 w-3.5 text-muted-foreground" /> {consumer.service}
                </span>

                <span className="text-muted-foreground text-xs font-semibold uppercase tracking-wider">Repository</span>
                {consumer.repository ? (
                    <a
                        href={consumer.repository}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                    >
                        <ExternalLink className="h-3 w-3 flex-shrink-0" />
                        <span className="truncate">{consumer.repository}</span>
                    </a>
                ) : (
                    <span className="text-xs text-muted-foreground italic">Not specified</span>
                )}
            </div>

            <Separator className="bg-border/50" />

            {/* Master-Detail Layout */}
            <MasterDetailLayout>
                <TopicMenu
                    title="Subscriptions"
                    items={topicMenuItems}
                    selectedId={selectedTopicName}
                    onSelect={(id: string | number) => setSelectedTopicNameState(id as string)}
                    emptyMessage="No topics subscribed."
                />

                <DetailPane
                    isSelected={selectedTopicName !== null}
                    title={selectedTopic?.name}
                    eventCount={selectedTopic?.events.length}
                    eventsLabel="Events Consumed"
                    emptySelectionMessage="Select a topic from the list to view its consumed events."
                    emptyEventsMessage="No events found for this subscription."
                >
                    {selectedTopic?.events.map((ev) => (
                        <EventCard
                            key={`${ev.name}-${ev.version}`}
                            name={ev.name}
                            version={ev.version}
                            onClick={onEventClick ? () => onEventClick(ev) : undefined}
                        />
                    ))}
                </DetailPane>
            </MasterDetailLayout>
        </ServiceDetailsDialog>
    );
}
