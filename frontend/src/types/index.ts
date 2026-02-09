// ─── Event Types ────────────────────────────────────────────

export interface EventHeader {
    name: string;
    description: string;
}

export interface Event {
    name: string;
    version?: string;
    description: string;
    schema_url?: string;
    headers?: EventHeader[];
    properties?: Record<string, unknown>;
}

export interface TopicWithEvents {
    topic: string;
    events: Event[];
}

// ─── Producer Types ─────────────────────────────────────────

export interface Producer {
    topic: string;
    name: string;
    owner: boolean;
    writes: boolean;
    description: string;
    events: Event[];
}

// ─── Consumer Types ─────────────────────────────────────────

export interface ConsumerEvent {
    name: string;
    version?: string;
}

export interface Topic {
    name: string;
    events: ConsumerEvent[];
}

export interface Consumer {
    group: string;
    description: string;
    topics: Topic[];
}

// ─── Config Types ───────────────────────────────────────────

export interface Server {
    environment: string;
    url: string;
}

export interface Config {
    servers: Server[];
    repository: string;
}

export interface EventDoctorSpec {
    version: string;
    service: string;
    config: Config;
    producers: Producer[];
    consumers: Consumer[];
}

// ─── Async State Types ──────────────────────────────────────

export type AsyncState<T> =
    | { status: 'idle' }
    | { status: 'loading' }
    | { status: 'success'; data: T }
    | { status: 'error'; error: string };

export interface ApiResponse<T> {
    data: T;
    message?: string;
}

export interface ApiErrorResponse {
    message: string;
    code?: string;
    status: number;
}