import { api } from '@/lib/api';
import type { TopicListView } from '@/types';
import { useCallback, useEffect, useState, useRef } from 'react';

const DEFAULT_PAGE_SIZE = 10;
const DEBOUNCE_MS = 500;
const MIN_SEARCH_LENGTH = 1;

export function useTopics() {
    const [data, setData] = useState<TopicListView | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
    const [search, setSearchRaw] = useState('');
    const [debouncedSearch, setDebouncedSearch] = useState('');
    const debounceRef = useRef<ReturnType<typeof setTimeout>>(undefined);

    const setSearch = useCallback((value: string) => {
        setSearchRaw(value);
        setPage(1);

        if (debounceRef.current) clearTimeout(debounceRef.current);
        debounceRef.current = setTimeout(() => {
            setDebouncedSearch(value.length >= MIN_SEARCH_LENGTH ? value : '');
        }, DEBOUNCE_MS);
    }, []);

    useEffect(() => {
        return () => {
            if (debounceRef.current) clearTimeout(debounceRef.current);
        };
    }, []);

    const fetchTopics = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const result = await api.getTopics(page, pageSize, debouncedSearch || undefined);
            setData(result);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            setData(null);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, debouncedSearch]);

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
        search,
        setSearch,
    };
}
