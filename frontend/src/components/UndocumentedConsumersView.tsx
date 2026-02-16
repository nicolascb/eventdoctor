import { EmptyState, SearchInput, StatCard } from "@/components/shared";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { useUndocumentedConsumers } from "@/hooks/useUndocumentedConsumers";
import { AlertTriangle, Layers } from "lucide-react";
import { useMemo, useState } from "react";

export function UndocumentedConsumersView() {
    const { undocumentedGroups, loading } = useUndocumentedConsumers();
    const [searchQuery, setSearchQuery] = useState("");

    const filteredGroups = useMemo(() => {
        if (!searchQuery) return undocumentedGroups;
        const q = searchQuery.toLowerCase();
        return undocumentedGroups.filter(g =>
            g.group.toLowerCase().includes(q) ||
            g.topic.toLowerCase().includes(q)
        );
    }, [undocumentedGroups, searchQuery]);

    if (loading) return null;

    return (
        <div className="space-y-4 animate-in">
            {/* Stats */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <StatCard
                    label="Undocumented Groups"
                    value={undocumentedGroups.length}
                    description="Consumer groups without documentation"
                    icon={<AlertTriangle className="h-4 w-4 text-muted-foreground" />}
                />
                <StatCard
                    label="Unique Topics"
                    value={new Set(undocumentedGroups.map(g => g.topic)).size}
                    description="Topics with undocumented consumers"
                    icon={<Layers className="h-4 w-4 text-muted-foreground" />}
                />
            </div>

            {/* Search */}
            <SearchInput
                value={searchQuery}
                onChange={setSearchQuery}
                placeholder="Search groups or topics..."
                resultCount={searchQuery ? filteredGroups.length : undefined}
                totalCount={searchQuery ? undocumentedGroups.length : undefined}
            />

            {/* Table */}
            <div className="rounded-lg border border-border bg-card overflow-hidden">
                {filteredGroups.length === 0 ? (
                    <div className="p-6">
                        <EmptyState
                            title="No undocumented consumers found"
                            description={searchQuery ? "Try adjusting your search criteria" : "All consumer groups are properly documented"}
                        />
                    </div>
                ) : (
                    <Table>
                        <TableHeader>
                            <TableRow className="hover:bg-transparent">
                                <TableHead className="text-[11px] font-medium text-left">Group</TableHead>
                                <TableHead className="text-[11px] font-medium text-left">Topic</TableHead>
                                <TableHead className="text-[11px] font-medium text-center">First Seen</TableHead>
                                <TableHead className="text-[11px] font-medium text-center">Last Updated</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {filteredGroups.map((g) => (
                                <TableRow key={`${g.group}-${g.topic}`}>
                                    <TableCell className="py-3">
                                        <div className="flex items-center gap-2">
                                            <AlertTriangle className="h-3.5 w-3.5 text-amber-500 flex-shrink-0" />
                                            <span className="font-medium text-sm">{g.group}</span>
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-3">
                                        <code className="text-xs font-mono">{g.topic}</code>
                                    </TableCell>
                                    <TableCell className="py-3 text-center">
                                        <span className="text-xs text-muted-foreground">{g.created_at}</span>
                                    </TableCell>
                                    <TableCell className="py-3 text-center">
                                        <span className="text-xs text-muted-foreground">{g.updated_at}</span>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                )}
            </div>
        </div>
    );
}
