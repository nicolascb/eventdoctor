import type { ConsumerView, EventsListView, EventView, OverviewResponse, ProducerDetailView, ProducersListView, TopicListView, UndocumentedConsumersView } from '@/types';

// In Docker, __API_URL__ is replaced at container startup with the actual API_URL env var.
// In development, VITE_API_URL takes precedence; the placeholder stays as-is and the fallback kicks in.
const PLACEHOLDER = '__API_URL__';
const API_BASE_URL =
    import.meta.env.VITE_API_URL ??
    (PLACEHOLDER.startsWith('__') ? '/v1' : PLACEHOLDER);

export class ApiError extends Error {
    status: number;
    body: string;

    constructor(
        message: string,
        status: number,
        body: string,
    ) {
        super(message);
        this.name = 'ApiError';
        this.status = status;
        this.body = body;
    }
}

class ApiClient {
    private baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    private async request<T>(path: string, options?: RequestInit): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        const response = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options?.headers,
            },
        });

        if (!response.ok) {
            const errorBody = await response.text().catch(() => 'Unknown error');
            throw new ApiError(
                `HTTP ${response.status}: ${response.statusText}`,
                response.status,
                errorBody,
            );
        }

        return response.json();
    }

    getProducers(page?: number, pageSize?: number, search?: string): Promise<ProducersListView> {
        const params = new URLSearchParams();
        if (page) params.set('page', String(page));
        if (pageSize) params.set('page_size', String(pageSize));
        if (search) params.set('search', search);
        const qs = params.toString();
        return this.request<ProducersListView>(`/producers${qs ? `?${qs}` : ''}`);
    }

    getProducerDetail(serviceId: number, topicId: number): Promise<ProducerDetailView> {
        return this.request<ProducerDetailView>(`/producers/${serviceId}/${topicId}`);
    }

    getEvents(page?: number, pageSize?: number, search?: string): Promise<EventsListView> {
        const params = new URLSearchParams();
        if (page) params.set('page', String(page));
        if (pageSize) params.set('page_size', String(pageSize));
        if (search) params.set('search', search);
        const qs = params.toString();
        return this.request<EventsListView>(`/events${qs ? `?${qs}` : ''}`);
    }

    getEvent(id: number): Promise<EventView> {
        return this.request<EventView>(`/events/${id}`);
    }

    getConsumers(page?: number, pageSize?: number, search?: string): Promise<ConsumerView> {
        const params = new URLSearchParams();
        if (page) params.set('page', String(page));
        if (pageSize) params.set('page_size', String(pageSize));
        if (search) params.set('search', search);
        const qs = params.toString();
        return this.request<ConsumerView>(`/consumers${qs ? `?${qs}` : ''}`);
    }

    getUndocumentedConsumers(): Promise<UndocumentedConsumersView> {
        return this.request<UndocumentedConsumersView>('/consumers?undocumented_only=true');
    }

    getTopics(page?: number, pageSize?: number): Promise<TopicListView> {
        const params = new URLSearchParams();
        if (page) params.set('page', String(page));
        if (pageSize) params.set('page_size', String(pageSize));
        const qs = params.toString();
        return this.request<TopicListView>(`/topics${qs ? `?${qs}` : ''}`);
    }

    getTopic(name: string): Promise<TopicListView> {
        return this.request<TopicListView>(`/topics/${name}`);
    }

    getOverview(): Promise<OverviewResponse> {
        return this.request<OverviewResponse>('/overview');
    }

    uploadConfig(config: string): Promise<{ message: string }> {
        return this.request<{ message: string }>('/config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-yaml',
            },
            body: config,
        });
    }
}

export const api = new ApiClient(API_BASE_URL);