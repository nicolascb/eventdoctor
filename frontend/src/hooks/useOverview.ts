import { api } from '@/lib/api';
import type { OverviewResponse } from '@/types';
import { useAsync } from './useAsync';

const emptyOverview: OverviewResponse = {
    total_topics: 0,
    total_events: 0,
    total_producers: 0,
    total_consumers: 0,
};

export function useOverview() {
    const { data, loading, error, refetch } = useAsync<OverviewResponse>(() => api.getOverview(), []);

    return {
        overview: data ?? emptyOverview,
        loading,
        error,
        refetch,
    };
}
