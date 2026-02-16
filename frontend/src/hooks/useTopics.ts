import { api } from '@/lib/api';
import type { TopicListView } from '@/types';
import { useCallback, useEffect, useState } from 'react';

const DEFAULT_PAGE_SIZE = 10;

export function useTopics() {
    const [data, setData] = useState<TopicListView | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);

    const fetchTopics = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const result = await api.getTopics(page, pageSize);
            setData(result);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            setData(null);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize]);

    useEffect(() => {
        fetchTopics();
    }, [fetchTopics]);

    return {
        topics: data?.topics ?? [],
        countEvents: data?.count_events ?? 0,
        countUnconsumed: data?.count_unconsumed ?? 0,
        countOrphaned: data?.count_orphaned ?? 0,
        pagination: data?.pagination,
        loading,
        error,
        refetch: fetchTopics,
        page,
        pageSize,
        setPage,
        setPageSize,
    };
}
