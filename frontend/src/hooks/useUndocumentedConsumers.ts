import { api } from '@/lib/api';
import type { UndocumentedConsumersView } from '@/types';
import { useAsync } from './useAsync';

export function useUndocumentedConsumers() {
    const { data, loading, error, refetch } = useAsync<UndocumentedConsumersView>(
        () => api.getUndocumentedConsumers(),
        [],
    );

    const undocumentedGroups = data?.groups_undocumented ?? [];

    return {
        undocumentedGroups,
        loading,
        error,
        refetch,
    };
}
