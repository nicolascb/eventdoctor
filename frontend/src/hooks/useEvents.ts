import { api } from '@/lib/api';
import type { EventsListView, Pagination } from '@/types';
import { useCallback, useEffect, useRef, useState } from 'react';

const DEFAULT_PAGE_SIZE = 15;
const DEBOUNCE_MS = 500;
const MIN_SEARCH_LENGTH = 3;

export function useEvents() {
    const [data, setData] = useState<EventsListView | null>(null);
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

    const fetchEvents = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const result = await api.getEvents(page, pageSize, debouncedSearch || undefined);
            setData(result);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            setData(null);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, debouncedSearch]);

    useEffect(() => {
        fetchEvents();
    }, [fetchEvents]);

    const events = data?.events ?? [];
    const pagination: Pagination | undefined = data?.pagination;

    return {
        events,
        loading,
        error,
        refetch: fetchEvents,
        pagination,
        page,
        pageSize,
        setPage,
        setPageSize,
        search,
        setSearch,
    };
}
