import { api } from '@/lib/api';
import type { Consumer } from '@/types';
import { useAsync } from './useAsync';

export function useConsumers() {
    const { data, loading, error, refetch } = useAsync<Consumer[]>(() => api.getConsumers(), []);

    const consumers = data ?? [];

    const totalTopics = consumers.reduce((acc, c) => acc + c.topics.length, 0);
    const totalEvents = consumers.reduce(
        (acc, c) => acc + c.topics.reduce((sum, t) => sum + t.events.length, 0),
        0,
    );

    return {
        consumers,
        loading,
        error,
        refetch,
        totalTopics,
        totalEvents,
    };
}
