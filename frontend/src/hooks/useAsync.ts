import { useCallback, useEffect, useRef, useState, type DependencyList } from 'react';

export function useAsync<T>(asyncFn: () => Promise<T>, deps: DependencyList = []) {
    const [data, setData] = useState<T | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const asyncFnRef = useRef(asyncFn);
    asyncFnRef.current = asyncFn;

    const execute = useCallback(async () => {
        setLoading(true);
        setError(null);

        try {
            const result = await asyncFnRef.current();
            setData(result);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
            setData(null);
        } finally {
            setLoading(false);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, deps);

    useEffect(() => {
        execute();
    }, [execute]);

    return { data, loading, error, refetch: execute };
}
