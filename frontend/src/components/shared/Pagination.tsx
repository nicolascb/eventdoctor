import { Button } from '@/components/ui/button';
import type { Pagination as PaginationType } from '@/types';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface PaginationProps {
    pagination: PaginationType;
    page: number;
    onPageChange: (page: number) => void;
    label?: string;
}

export function Pagination({ pagination, page, onPageChange, label = 'items' }: PaginationProps) {
    if (pagination.total_pages <= 1) return null;

    return (
        <div className="flex items-center justify-between px-1">
            <span className="text-xs text-muted-foreground">
                Page {pagination.page} of {pagination.total_pages} ({pagination.total} {label})
            </span>
            <div className="flex items-center gap-1">
                <Button
                    variant="outline"
                    size="sm"
                    className="h-7 px-2 text-xs"
                    disabled={page <= 1}
                    onClick={() => onPageChange(page - 1)}
                >
                    <ChevronLeft className="h-3.5 w-3.5 mr-1" />
                    Previous
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    className="h-7 px-2 text-xs"
                    disabled={page >= pagination.total_pages}
                    onClick={() => onPageChange(page + 1)}
                >
                    Next
                    <ChevronRight className="h-3.5 w-3.5 ml-1" />
                </Button>
            </div>
        </div>
    );
}
