import { api } from '@/lib/api';
import type { Producer } from '@/types';
import { useAsync } from './useAsync';

export function useProducers() {
    const { data, loading, error, refetch } = useAsync<Producer[]>(() => api.getProducers(), []);

    const producers = data ?? [];

    const totalEvents = producers.reduce((acc, p) => acc + p.events.length, 0);
    const ownersCount = producers.filter(p => p.owner).length;
    const activeWriters = producers.filter(p => p.writes).length;

    return {
        producers,
        loading,
        error,
        refetch,
        totalEvents,
        ownersCount,
        activeWriters,
    };
}
