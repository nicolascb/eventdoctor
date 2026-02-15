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
import { api } from "@/lib/api";
import type { ProducerDetailView, ProducerEventEntry, ProducerServiceItem } from "@/types";
import { Skeleton } from "@/components/ui/skeleton";
import { AlertCircle, FileJson, Layers, Loader2, ExternalLink, Zap, ChevronLeft, ChevronRight, Server, Workflow } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useState } from "react";
import { SearchInput } from "@/components/shared/SearchInput";
import { ProducerDetails } from "./ProducerDetails";

export function ProducersView() {
    const {
        producers,
        loading,
        error,
        pagination,
        page,
        setPage,
        setSearch,
        search,
    } = useProducers();

    // Detail dialog state
    const [selectedProducer, setSelectedProducer] = useState<ProducerServiceItem | null>(null);
    const [detail, setDetail] = useState<ProducerDetailView | null>(null);
    const [detailLoading, setDetailLoading] = useState(false);
    const [selectedEvent, setSelectedEvent] = useState<ProducerEventEntry | null>(null);

    const handleTopicClick = async (serviceID: number, topicID: number) => {
        setDetailLoading(true);
        setDetail(null);
        setSelectedEvent(null);
        try {
            const result = await api.getProducerDetail(serviceID, topicID);
            setDetail(result);
        } catch {
            setDetail(null);
        } finally {
            setDetailLoading(false);
        }
    };

    const handleCloseDetail = () => {
        setDetail(null);
        setSelectedEvent(null);
    };

    if (error) {
        return (
            <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
            </Alert>
        );
    }

    return (
        <div className="space-y-4">
            <SearchInput
                value={search}
                onChange={setSearch}
                placeholder="Search producers..."
            />
            <div className="rounded-lg border border-border bg-card overflow-hidden">
                <Table>
                    <TableHeader>
                        <TableRow className="hover:bg-transparent">
                            <TableHead className="text-[11px] font-medium">Service</TableHead>
                            <TableHead className="text-[11px] font-medium text-left">Topics</TableHead>
                            <TableHead className="text-[11px] font-medium text-center">Events</TableHead>
                            <TableHead className="text-[11px] font-medium text-right"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <TableRow key={i}>
                                    <TableCell><Skeleton className="h-4 w-[200px]" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-[150px]" /></TableCell>
                                    <TableCell><Skeleton className="h-4 w-[50px] ml-auto" /></TableCell>
                                </TableRow>
                            ))
                        ) : producers.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={3} className="h-24 text-center">
                                    No producers found.
                                </TableCell>
                            </TableRow>
                        ) : (
                            producers.map((service) => (
                                <TableRow
                                    key={service.service_id}
                                    className="cursor-pointer group"
                                    onClick={() => { setSelectedProducer(service); setDetail(null); }}
                                >
                                    <TableCell className="py-3">
                                        <div className="flex flex-col gap-0.5">
                                            <div className="flex items-center gap-2">
                                                <Server className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                                <span className="font-medium text-sm">{service.service}</span>
                                            </div>
                                            <div className="text-xs text-muted-foreground pl-5.5 truncate max-w-[250px]">
                                                {service.repository}
                                            </div>
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-3 text-center">
                                        <div className="flex flex-wrap gap-1 justify-start">
                                            {service.topics.slice(0, 3).map((topic) => (
                                                <Badge
                                                    key={topic.topic_id}
                                                    variant="secondary"
                                                    className="font-normal text-[10px]"
                                                >
                                                    {topic.topic}
                                                </Badge>
                                            ))}
                                            {service.topics.length > 3 && (
                                                <Badge
                                                    variant="outline"
                                                    className="text-[10px] font-normal text-muted-foreground"
                                                >
                                                    +{service.topics.length - 3} topics
                                                </Badge>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-3 text-center">
                                        <Badge variant="secondary" className="text-[10px] font-normal tabular-nums">
                                            {service.topics.reduce((acc, t) => acc + t.event_count, 0)}
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
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>

            {/* Pagination Controls */}
            {pagination && pagination.total_pages > 1 && (
                <div className="flex items-center justify-between px-1">
                    <span className="text-xs text-muted-foreground">
                        Page {pagination.page} of {pagination.total_pages} ({pagination.total} services)
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

            {/* Producer List Dialog */}
            <Dialog open={!!selectedProducer && !detail && !detailLoading} onOpenChange={(open) => { if (!open) setSelectedProducer(null); }}>
                {selectedProducer && (
                    <DialogContent className="max-w-xl">
                        <ProducerDetails
                            producer={selectedProducer}
                            onTopicClick={(serviceId, topicId) => handleTopicClick(serviceId, topicId)}
                        />
                    </DialogContent>
                )}
            </Dialog>

            {/* Topic Detail Dialog (shows events) */}
            <Dialog open={detailLoading || (!!detail && !selectedEvent)} onOpenChange={(open) => { if (!open) handleCloseDetail(); }}>
                <DialogContent className="max-w-xl">
                    {detailLoading ? (
                        <div className="flex items-center justify-center py-8">
                            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                        </div>
                    ) : detail && (
                        <>
                            <DialogHeader>
                                <div className="flex items-center gap-2">
                                    <Layers className="h-4 w-4 text-muted-foreground" />
                                    <DialogTitle className="text-base font-semibold font-mono">
                                        {detail.topic}
                                    </DialogTitle>
                                    {detail.owner && (
                                        <Badge variant="outline" className="text-[10px] font-normal">Owner</Badge>
                                    )}
                                </div>
                                {detail.description && (
                                    <DialogDescription>
                                        {detail.description}
                                    </DialogDescription>
                                )}
                            </DialogHeader>

                            {/* Metadata rows */}
                            <div className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm py-1">
                                <span className="text-muted-foreground text-xs font-medium">Service</span>
                                <span className="text-xs font-medium">{detail.service}</span>
                                <span className="text-muted-foreground text-xs font-medium">Repository</span>
                                {detail.repository ? (
                                    <a
                                        href={detail.repository}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 transition-colors truncate"
                                    >
                                        <ExternalLink className="h-3 w-3 flex-shrink-0" />
                                        <span className="truncate">{detail.repository}</span>
                                    </a>
                                ) : (
                                    <span className="text-xs text-muted-foreground italic">Not specified</span>
                                )}
                            </div>

                            {/* Events */}
                            {detail.events.length > 0 && (
                                <>
                                    <Separator />
                                    <div className="space-y-2">
                                        <div className="flex items-center justify-between">
                                            <h4 className="text-xs font-medium text-muted-foreground">
                                                Published Events
                                            </h4>
                                            <span className="text-xs font-mono text-muted-foreground">
                                                {detail.events.length}
                                            </span>
                                        </div>

                                        <div className="space-y-1.5">
                                            {detail.events.map((ev) => (
                                                <button
                                                    key={`${ev.name}-${ev.version || ''}`}
                                                    type="button"
                                                    className="w-full text-left rounded-md border border-border px-3 py-2.5 hover:bg-accent/30 transition-colors cursor-pointer group"
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
                                                        <div className="flex items-center gap-2 flex-shrink-0 group-hover:text-primary transition-colors">
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
                        </>
                    )}
                </DialogContent>
            </Dialog>

            {/* Event Detail Dialog */}
            <Dialog open={!!selectedEvent} onOpenChange={(open) => { if (!open) setSelectedEvent(null); }}>
                {selectedEvent && detail && (
                    <DialogContent className="max-w-3xl">
                        <DialogHeader>
                            <div className="flex items-center gap-2 mb-1">
                                <Zap className="h-4 w-4 text-muted-foreground" />
                                <DialogTitle className="font-mono text-base flex items-center gap-2">
                                    <span className="text-muted-foreground font-normal">{detail.topic}</span>
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
                                                <TableRow className="bg-muted/50 hover:bg-muted/50">
                                                    <TableHead className="h-8 text-[11px] font-medium">Header Name</TableHead>
                                                    <TableHead className="h-8 text-[11px] font-medium">Description</TableHead>
                                                </TableRow>
                                            </TableHeader>
                                            <TableBody>
                                                {selectedEvent.headers.map((header) => (
                                                    <TableRow key={header.name} className="hover:bg-transparent">
                                                        <TableCell className="py-2 text-xs font-mono">{header.name}</TableCell>
                                                        <TableCell className="py-2 text-xs text-muted-foreground">{header.description || '-'}</TableCell>
                                                    </TableRow>
                                                ))}
                                            </TableBody>
                                        </Table>
                                    </div>
                                </div>
                            </>
                        )}
                    </DialogContent>
                )}
            </Dialog>
        </div>
    );
}