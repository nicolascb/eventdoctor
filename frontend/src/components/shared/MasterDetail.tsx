import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Box, CheckCircle2, ExternalLink, FileJson, Loader2 } from "lucide-react";
import React from "react";

// --- EVENT CARD ---

export interface EventCardProps {
    name: string;
    version?: string;
    description?: string;
    schemaUrl?: string;
    onClick?: () => void;
}

export function EventCard({ name, version, description, schemaUrl, onClick }: EventCardProps) {
    return (
        <div
            className={`flex flex-col sm:flex-row sm:items-center justify-between p-4 rounded-lg border bg-card transition-all gap-4 group ${onClick ? "hover:bg-muted/30 cursor-pointer hover:border-border/80 hover:shadow-sm" : ""
                }`}
            onClick={onClick}
            role={onClick ? "button" : undefined}
            tabIndex={onClick ? 0 : undefined}
            onKeyDown={(e) => {
                if (e.key === "Enter" && onClick) onClick();
            }}
        >
            <div className="space-y-1.5 flex-1 min-w-0">
                <div className="flex items-center gap-2">
                    <span className="font-mono text-sm font-medium text-foreground truncate">{name}</span>
                    {version && (
                        <Badge variant="secondary" className="text-[10px] h-5 px-2 font-mono shrink-0">
                            v{version}
                        </Badge>
                    )}
                </div>
                {description && (
                    <p className="text-xs text-muted-foreground/80 line-clamp-2 pr-4">{description}</p>
                )}
            </div>
            <div className="flex items-center gap-1.5 self-start sm:self-center shrink-0">
                {schemaUrl && (
                    <Button
                        variant="outline"
                        size="icon"
                        className="h-8 w-8 bg-background opacity-0 group-hover:opacity-100 transition-opacity"
                        asChild
                    >
                        <a
                            href={schemaUrl}
                            target="_blank"
                            rel="noopener noreferrer"
                            onClick={(e) => e.stopPropagation()}
                            title="Open Schema URL"
                        >
                            <ExternalLink className="h-3.5 w-3.5 text-muted-foreground" />
                        </a>
                    </Button>
                )}
                <Button
                    variant="secondary"
                    size="icon"
                    className="h-8 w-8 opacity-70 group-hover:opacity-100 transition-opacity"
                    title="View Payload Details"
                >
                    <FileJson className="h-4 w-4" />
                </Button>
            </div>
        </div>
    );
}

// --- MASTER DETAIL LAYOUT ---

export function MasterDetailLayout({ children }: { children: React.ReactNode }) {
    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 relative min-h-[400px]">
            {children}
        </div>
    );
}

// --- TOPIC MENU (LEFT SIDE) ---

export interface TopicMenuItem {
    id: string | number;
    name: string;
    eventCount: number;
    isOwner?: boolean;
}

interface TopicMenuProps {
    title: string;
    icon?: React.ReactNode;
    items: TopicMenuItem[];
    selectedId: string | number | null;
    onSelect: (id: string | number) => void;
    emptyMessage?: string;
}

export function TopicMenu({ title, icon, items, selectedId, onSelect, emptyMessage = "No topics found." }: TopicMenuProps) {
    return (
        <div className="md:col-span-1 border-r border-border/50 pr-6 space-y-4 h-full">
            <div className="flex items-center gap-2">
                {icon || <Box className="h-4 w-4 text-muted-foreground" />}
                <h4 className="text-sm font-semibold text-foreground">{title}</h4>
                <Badge variant="secondary" className="font-mono text-[10px] ml-auto">
                    {items.length}
                </Badge>
            </div>
            <div className="space-y-2">
                {items.length === 0 ? (
                    <p className="text-xs text-muted-foreground italic">{emptyMessage}</p>
                ) : (
                    items.map((item) => {
                        const isSelected = selectedId === item.id;
                        return (
                            <button
                                key={item.id}
                                type="button"
                                onClick={() => onSelect(item.id)}
                                className={`w-full flex items-center justify-between p-3 rounded-md transition-colors border text-left flex-shrink-0 ${isSelected
                                        ? "bg-primary/5 border-primary/20 ring-1 ring-primary/20"
                                        : "bg-card border-transparent hover:bg-muted/50 hover:border-border"
                                    }`}
                            >
                                <div className="flex flex-col gap-1.5 overflow-hidden w-full">
                                    <div className="flex items-center gap-2">
                                        <Box
                                            className={`h-3.5 w-3.5 flex-shrink-0 ${isSelected ? "text-primary" : "text-muted-foreground"
                                                }`}
                                        />
                                        <span
                                            className={`font-mono text-xs font-medium truncate ${isSelected ? "text-foreground" : "text-foreground/80"
                                                }`}
                                        >
                                            {item.name}
                                        </span>
                                    </div>
                                    <div className="flex items-center gap-2 pl-5.5">
                                        {item.isOwner && (
                                            <span
                                                className="text-yellow-600 dark:text-yellow-400 text-[9px] flex items-center gap-1"
                                                title="Owner"
                                            >
                                                <CheckCircle2 className="h-2.5 w-2.5" /> Owner
                                            </span>
                                        )}
                                        <Badge
                                            variant="outline"
                                            className={`text-[9px] px-1.5 py-0 h-4 ${isSelected ? "bg-primary/10 border-primary/20" : "bg-background"
                                                }`}
                                        >
                                            {item.eventCount} events
                                        </Badge>
                                    </div>
                                </div>
                            </button>
                        );
                    })
                )}
            </div>
        </div>
    );
}

// --- DETAIL PANE (RIGHT SIDE) ---

interface DetailPaneProps {
    isSelected: boolean;
    isLoading?: boolean;
    hasFailed?: boolean;
    title?: string;
    description?: string;
    eventCount?: number;
    eventsLabel?: string;
    emptySelectionMessage?: string;
    emptyEventsMessage?: string;
    failedMessage?: string;
    children?: React.ReactNode;
}

export function DetailPane({
    isSelected,
    isLoading,
    hasFailed,
    title,
    description,
    eventCount,
    eventsLabel = "Events",
    emptySelectionMessage = "Select an item to view details.",
    emptyEventsMessage = "No events found.",
    failedMessage = "Failed to load details.",
    children,
}: DetailPaneProps) {
    if (!isSelected) {
        return (
            <div className="md:col-span-2 space-y-4 h-full">
                <div className="h-full flex flex-col items-center justify-center text-center p-8 bg-muted/10 rounded-lg border border-dashed border-border/50 min-h-[300px]">
                    <Box className="h-8 w-8 text-muted-foreground/30 mb-3" />
                    <p className="text-sm text-muted-foreground">{emptySelectionMessage}</p>
                </div>
            </div>
        );
    }

    if (isLoading) {
        return (
            <div className="md:col-span-2 space-y-4 h-full">
                <div className="h-full flex flex-col items-center justify-center py-12 min-h-[300px]">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground mb-4" />
                    <p className="text-sm text-muted-foreground animate-pulse">Loading details...</p>
                </div>
            </div>
        );
    }

    if (hasFailed) {
        return (
            <div className="md:col-span-2 space-y-4">
                <div className="text-center py-8 text-sm text-destructive bg-destructive/5 rounded-md border border-destructive/20 mt-4">
                    {failedMessage}
                </div>
            </div>
        );
    }

    return (
        <div className="md:col-span-2 space-y-4">
            <div className="animate-in fade-in slide-in-from-right-4 duration-300 space-y-6">
                <div className="space-y-2">
                    <h3 className="text-lg font-mono font-medium">{title}</h3>
                    {description && (
                        <p className="text-sm text-muted-foreground bg-muted/30 p-3 rounded-md border border-border/50">
                            {description}
                        </p>
                    )}
                </div>

                <div className="space-y-3">
                    <div className="flex items-center justify-between border-b pb-2">
                        <h5 className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
                            {eventsLabel}
                        </h5>
                        {eventCount !== undefined && (
                            <Badge variant="secondary" className="text-[10px]">
                                {eventCount}
                            </Badge>
                        )}
                    </div>

                    <div className="grid gap-3 max-h-[400px] overflow-y-auto pr-2 pb-4">
                        {eventCount === 0 ? (
                            <p className="text-sm text-muted-foreground italic text-center py-8">
                                {emptyEventsMessage}
                            </p>
                        ) : (
                            children
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
