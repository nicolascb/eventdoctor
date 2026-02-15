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
    service: string;
    repository: string;
    topic: string;
    description: string;
    owner: boolean;
    writes: boolean;
    events: Event[];
}

export interface ProducerTopic {
    topic: string;
    description: string;
    owner: boolean;
    writes: boolean;
    events: Event[];
}

export interface GroupedProducer {
    service: string;
    repository: string;
    topics: ProducerTopic[];
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
    service: string;
    repository: string;
    group: string;
    description?: string;
    topics: Topic[];
}

export interface UndocumentedGroup {
    topic: string;
    group: string;
    created_at: string;
    updated_at: string;
}

export interface Pagination {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
}

export interface ConsumerView {
    consumers: Consumer[];
    pagination?: Pagination;
}

export interface UndocumentedConsumersView {
    groups_undocumented: UndocumentedGroup[];
}

// ─── Topic View Types ──────────────────────────────────────

export interface TopicProducerEntry {
    service: string;
    repository: string;
    event: string;
    writes: boolean;
    owner: boolean;
}

export interface TopicConsumerEntry {
    service: string;
    repository: string;
    event: string;
    group: string;
    version?: string;
}

export interface TopicView {
    topic: string;
    owner_service?: string;
    events: Event[];
    producers: TopicProducerEntry[];
    consumers: TopicConsumerEntry[];
}

export interface TopicListView {
    topics: TopicView[];
    count_events: number;
    count_unconsumed: number;
    count_orphaned: number;
    pagination?: Pagination;
}

// ─── Overview Types ─────────────────────────────────────────

export interface OverviewResponse {
    total_topics: number;
    total_events: number;
    total_producers: number;
    total_consumers: number;
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