import type { Consumer, OverviewResponse, Producer, TopicWithEvents } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8087/v1';

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

    getProducers(): Promise<Producer[]> {
        return this.request<Producer[]>('/producers');
    }

    getEvents(): Promise<TopicWithEvents[]> {
        return this.request<TopicWithEvents[]>('/events');
    }

    getConsumers(): Promise<Consumer[]> {
        return this.request<Consumer[]>('/consumers');
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