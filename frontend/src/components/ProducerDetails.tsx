import { Separator } from "@/components/ui/separator";
import { api } from "@/lib/api";
import type { ProducerDetailView, ProducerEventEntry, ProducerServiceItem } from "@/types";
import { ExternalLink, Layers, Server } from "lucide-react";
import { useEffect, useState } from "react";
import { DetailPane, EventCard, MasterDetailLayout, TopicMenu, type TopicMenuItem } from "./shared/MasterDetail";
import { ServiceDetailsDialog } from "./ServiceDetailsDialog";

interface ProducerDetailsProps {
    producer?: ProducerServiceItem | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onEventClick: (event: ProducerEventEntry) => void;
}

export function ProducerDetails({ producer, open, onOpenChange, onEventClick }: ProducerDetailsProps) {
    const [selectedTopicIdState, setSelectedTopicIdState] = useState<number | null>(null);
    const [topicDetail, setTopicDetail] = useState<ProducerDetailView | null>(null);
    const [loadingDetail, setLoadingDetail] = useState(false);
    const [hasFailed, setHasFailed] = useState(false);

    const hasTopics = producer && producer.topics.length > 0;
    const isValidSelection =
        selectedTopicIdState !== null && producer && producer.topics.some((t) => t.topic_id === selectedTopicIdState);
    const selectedTopicId = isValidSelection ? selectedTopicIdState : hasTopics ? producer!.topics[0].topic_id : null;

    // Fetch details when selected topic changes
    useEffect(() => {
        let mounted = true;

        if (selectedTopicId && producer?.service_id) {
            // eslint-disable-next-line react-hooks/set-state-in-effect
            setTopicDetail(null);
             
            setLoadingDetail(true);
             
            setHasFailed(false);
            api.getProducerDetail(producer.service_id, selectedTopicId)
                .then((detail) => {
                    if (mounted) setTopicDetail(detail);
                })
                .catch((error) => {
                    console.error("Failed to fetch topic details:", error);
                    if (mounted) {
                        setTopicDetail(null);
                        setHasFailed(true);
                    }
                })
                .finally(() => {
                    if (mounted) setLoadingDetail(false);
                });
        }
        return () => {
            mounted = false;
        };
    }, [selectedTopicId, producer?.service_id]);

    if (!open || !producer) return null;

    const topicMenuItems: TopicMenuItem[] = producer.topics.map((t) => ({
        id: t.topic_id,
        name: t.topic,
        eventCount: t.event_count,
        isOwner: t.owner,
    }));

    return (
        <ServiceDetailsDialog
            open={open}
            onOpenChange={onOpenChange}
            icon={<Server className="h-6 w-6 text-muted-foreground" />}
            title={producer.service}
            description={producer.repository}
            className="sm:max-w-4xl" // Expanded width for Master-Detail layout
        >
            {/* Metadata rows */}
            <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm px-1">
                <span className="text-muted-foreground text-xs font-semibold uppercase tracking-wider">Repository</span>
                {producer.repository ? (
                    <a
                        href={producer.repository}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                    >
                        <ExternalLink className="h-3 w-3 flex-shrink-0" />
                        <span className="truncate">{producer.repository}</span>
                    </a>
                ) : (
                    <span className="text-xs text-muted-foreground italic">Not specified</span>
                )}
            </div>

            <Separator className="bg-border/50" />

            {/* Master-Detail Layout */}
            <MasterDetailLayout>
                <TopicMenu
                    title="Published Topics"
                    icon={<Layers className="h-4 w-4 text-muted-foreground" />}
                    items={topicMenuItems}
                    selectedId={selectedTopicId}
                    onSelect={(id: string | number) => setSelectedTopicIdState(id as number)}
                    emptyMessage="No topics published."
                />

                <DetailPane
                    isSelected={selectedTopicId !== null}
                    isLoading={loadingDetail}
                    hasFailed={hasFailed}
                    title={topicDetail?.topic}
                    description={topicDetail?.description}
                    eventCount={topicDetail?.events.length}
                    eventsLabel="Events Produced"
                    emptySelectionMessage="Select a topic from the list to view its events."
                    emptyEventsMessage="No events found for this topic."
                >
                    {topicDetail?.events.map((ev) => (
                        <EventCard
                            key={`${ev.name}-${ev.version}`}
                            name={ev.name}
                            version={ev.version}
                            description={ev.description}
                            schemaUrl={ev.schema_url}
                            onClick={() => onEventClick(ev)}
                        />
                    ))}
                </DetailPane>
            </MasterDetailLayout>
        </ServiceDetailsDialog>
    );
}
