import { api } from '@/lib/api';
import type { ConsumerView, Pagination } from '@/types';
import { useCallback, useEffect, useState } from 'react';

const DEFAULT_PAGE_SIZE = 15;

export function useConsumers() {
    const [data, setData] = useState<ConsumerView | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);

    const fetchConsumers = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const result = await api.getConsumers(page, pageSize);
            setData(result);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            setData(null);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize]);

    useEffect(() => {
        fetchConsumers();
    }, [fetchConsumers]);

    const consumers = data?.consumers ?? [];
    const pagination: Pagination | undefined = data?.pagination;

    const totalTopics = consumers.reduce((acc, c) => acc + c.topics.length, 0);
    const totalEvents = consumers.reduce(
        (acc, c) => acc + c.topics.reduce((sum, t) => sum + t.events.length, 0),
        0,
    );

    return {
        consumers,
        loading,
        error,
        refetch: fetchConsumers,
        totalTopics,
        totalEvents,
        pagination,
        page,
        pageSize,
        setPage,
        setPageSize,
    };
}
