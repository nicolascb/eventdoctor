import { api } from '@/lib/api';
import type { TopicWithEvents } from '@/types';
import { useAsync } from './useAsync';

export function useEvents() {
    const { data, loading, error, refetch } = useAsync<TopicWithEvents[]>(() => api.getEvents(), []);

    const topics = data ?? [];
    const totalEvents = topics.reduce((acc, t) => acc + t.events.length, 0);

    return {
        topics,
        loading,
        error,
        refetch,
        totalEvents,
    };
}
