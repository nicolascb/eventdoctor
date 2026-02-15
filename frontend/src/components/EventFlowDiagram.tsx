import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import type { TopicView } from "@/types";
import {
    Check,
    ChevronDown,
    Database,
    Layers,
    Network,
    Search,
    X,
    Zap,
} from "lucide-react";
import {
    useCallback,
    useLayoutEffect,
    useMemo,
    useRef,
    useState,
} from "react";

/* ────────────────────────────────────────────────────────── */
/*  Types                                                     */
/* ────────────────────────────────────────────────────────── */

interface EventFlowDiagramProps {
    topics: TopicView[];
    hideSelector?: boolean;
}

interface FlowEvent {
    name: string;
    producers: string[];
    consumers: string[];
}

interface FlowData {
    topic: string;
    flowEvents: FlowEvent[];
}

interface ConnLine {
    id: string;
    path: string;
    color: string;
    delay: number;
}

/* ────────────────────────────────────────────────────────── */
/*  TopicSelect – searchable multi-select dropdown            */
/* ────────────────────────────────────────────────────────── */

function TopicSelect({
    options,
    selected,
    onChange,
}: {
    options: string[];
    selected: string[];
    onChange: (v: string[]) => void;
}) {
    const [open, setOpen] = useState(false);
    const [search, setSearch] = useState("");
    const ref = useRef<HTMLDivElement>(null);

    // close on outside click
    useLayoutEffect(() => {
        if (!open) return;
        const handler = (e: MouseEvent) => {
            if (ref.current && !ref.current.contains(e.target as Node))
                setOpen(false);
        };
        document.addEventListener("mousedown", handler);
        return () => document.removeEventListener("mousedown", handler);
    }, [open]);

    const filtered = options.filter((o) =>
        o.toLowerCase().includes(search.toLowerCase()),
    );

    const toggle = (name: string) =>
        onChange(
            selected.includes(name)
                ? selected.filter((s) => s !== name)
                : [...selected, name],
        );

    return (
        <div ref={ref} className="relative w-full max-w-sm">
            {/* trigger */}
            <button
                type="button"
                onClick={() => setOpen(!open)}
                className="flex w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm hover:bg-accent transition-colors"
            >
                <span className="truncate text-muted-foreground">
                    {selected.length === 0
                        ? "Select topics…"
                        : `${selected.length} topic${selected.length > 1 ? "s" : ""} selected`}
                </span>
                <ChevronDown className="h-4 w-4 shrink-0 opacity-50" />
            </button>

            {/* selected pills */}
            {selected.length > 0 && (
                <div className="mt-1.5 flex flex-wrap gap-1">
                    {selected.map((s) => (
                        <Badge
                            key={s}
                            variant="secondary"
                            className="gap-1 pr-1 text-xs cursor-pointer"
                            onClick={() => toggle(s)}
                        >
                            {s}
                            <X className="h-3 w-3" />
                        </Badge>
                    ))}
                </div>
            )}

            {/* dropdown */}
            {open && (
                <div className="absolute z-50 mt-1 w-full rounded-md border bg-popover shadow-lg">
                    <div className="flex items-center border-b px-3">
                        <Search className="h-4 w-4 shrink-0 text-muted-foreground" />
                        <Input
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            placeholder="Search topics…"
                            className="border-0 shadow-none focus-visible:ring-0 h-9 text-sm"
                            autoFocus
                        />
                    </div>
                    <div className="max-h-52 overflow-y-auto p-1">
                        {filtered.length === 0 && (
                            <p className="py-4 text-center text-sm text-muted-foreground">
                                No topics found.
                            </p>
                        )}
                        {filtered.map((name) => {
                            const active = selected.includes(name);
                            return (
                                <button
                                    key={name}
                                    type="button"
                                    onClick={() => toggle(name)}
                                    className={`flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm transition-colors ${active
                                        ? "bg-accent text-accent-foreground"
                                        : "hover:bg-accent/50"
                                        }`}
                                >
                                    <span
                                        className={`flex h-4 w-4 items-center justify-center rounded-sm border ${active
                                            ? "bg-primary border-primary text-primary-foreground"
                                            : "border-muted-foreground/30"
                                            }`}
                                    >
                                        {active && <Check className="h-3 w-3" />}
                                    </span>
                                    {name}
                                </button>
                            );
                        })}
                    </div>
                </div>
            )}
        </div>
    );
}

/* ────────────────────────────────────────────────────────── */
/*  EventFlowLine – one event row with SVG connections        */
/* ────────────────────────────────────────────────────────── */

function EventFlowLine({
    event,
    producerNames,
    consumerNames,
    rowIndex,
}: {
    event: string;
    producerNames: string[];
    consumerNames: string[];
    rowIndex: number;
}) {
    const containerRef = useRef<HTMLDivElement>(null);
    const eventRef = useRef<HTMLDivElement>(null);
    const prodRefs = useRef<Map<string, HTMLDivElement>>(new Map());
    const consRefs = useRef<Map<string, HTMLDivElement>>(new Map());
    const [lines, setLines] = useState<ConnLine[]>([]);

    const setProdRef = useCallback(
        (name: string) => (el: HTMLDivElement | null) => {
            if (el) prodRefs.current.set(name, el);
            else prodRefs.current.delete(name);
        },
        [],
    );
    const setConsRef = useCallback(
        (name: string) => (el: HTMLDivElement | null) => {
            if (el) consRefs.current.set(name, el);
            else consRefs.current.delete(name);
        },
        [],
    );

    const computeLines = useCallback(() => {
        const container = containerRef.current;
        const eventEl = eventRef.current;
        if (!container || !eventEl) return;

        const cr = container.getBoundingClientRect();
        const er = eventEl.getBoundingClientRect();
        const ex = er.left - cr.left + er.width / 2;
        const ey = er.top - cr.top + er.height / 2;

        const result: ConnLine[] = [];
        let idx = 0;

        // producer → event
        prodRefs.current.forEach((el, name) => {
            const r = el.getBoundingClientRect();
            const sx = r.right - cr.left;
            const sy = r.top - cr.top + r.height / 2;
            const dx = (ex - sx) * 0.45;
            const path = `M ${sx} ${sy} C ${sx + dx} ${sy}, ${ex - dx} ${ey}, ${ex} ${ey}`;
            result.push({
                id: `p-${name}-${idx}`,
                path,
                color: "#6C757D",
                delay: idx * 0.3,
            });
            idx++;
        });

        // event → consumer
        consRefs.current.forEach((el, name) => {
            const r = el.getBoundingClientRect();
            const tx = r.left - cr.left;
            const ty = r.top - cr.top + r.height / 2;
            const eRight = er.right - cr.left;
            const dx = (tx - eRight) * 0.45;
            const path = `M ${eRight} ${ey} C ${eRight + dx} ${ey}, ${tx - dx} ${ty}, ${tx} ${ty}`;
            result.push({
                id: `c-${name}-${idx}`,
                path,
                color: "#6C757D",
                delay: idx * 0.3,
            });
            idx++;
        });

        setLines(result);
    }, []);

    useLayoutEffect(() => {
        const delay = rowIndex * 100 + 400;
        const timer = setTimeout(computeLines, delay);
        return () => clearTimeout(timer);
    }, [computeLines, rowIndex]);

    // Recalculate on resize
    useLayoutEffect(() => {
        const container = containerRef.current;
        if (!container) return;
        const ro = new ResizeObserver(() => computeLines());
        ro.observe(container);
        return () => ro.disconnect();
    }, [computeLines]);

    const baseDelay = rowIndex * 0.12;

    return (
        <div
            ref={containerRef}
            className="relative flex items-center gap-0 py-2"
            style={{ minHeight: 48 }}
        >
            {/* SVG overlay */}
            <svg
                className="pointer-events-none absolute inset-0 h-full w-full"
                style={{ zIndex: 1 }}
            >
                {lines.map((l) => (
                    <g key={l.id}>
                        {/* line */}
                        <path
                            d={l.path}
                            fill="none"
                            stroke={l.color}
                            strokeWidth={1.5}
                            strokeDasharray="6 4"
                            opacity={0.5}
                        />
                        {/* envelope */}
                        <g opacity={0}>
                            <text
                                fontSize="12"
                                textAnchor="middle"
                                dominantBaseline="central"
                                fill="#6C757D"
                            >
                                ✉
                            </text>
                            <animateMotion
                                path={l.path}
                                dur="4s"
                                repeatCount="indefinite"
                                begin={`${l.delay}s`}
                                keyPoints="0;1"
                                keyTimes="0;1"
                                calcMode="linear"
                            />
                            <animate
                                attributeName="opacity"
                                values="0;1;1;0"
                                keyTimes="0;0.1;0.85;1"
                                dur="4s"
                                repeatCount="indefinite"
                                begin={`${l.delay}s`}
                            />
                        </g>
                    </g>
                ))}
            </svg>

            {/* Producers column */}
            <div
                className="flex flex-col items-end gap-1 shrink-0"
                style={{ width: "28%", zIndex: 2 }}
            >
                {producerNames.map((p, i) => (
                    <div
                        key={p}
                        ref={setProdRef(p)}
                        className="flow-node-enter flex items-center gap-1.5 rounded-full border bg-card px-3 py-1 text-xs font-medium shadow-sm"
                        style={{ animationDelay: `${baseDelay + i * 0.08}s` }}
                    >
                        <Database className="h-3 w-3 text-muted-foreground" />
                        <span className="truncate max-w-[120px]">{p}</span>
                    </div>
                ))}
                {producerNames.length === 0 && (
                    <span className="text-xs text-muted-foreground italic">—</span>
                )}
            </div>

            {/* Center: Event pill */}
            <div
                className="flex justify-center shrink-0"
                style={{ width: "44%", zIndex: 2 }}
            >
                <div
                    ref={eventRef}
                    className="flow-node-enter flex items-center gap-1.5 rounded-full bg-primary text-primary-foreground px-4 py-1.5 text-xs font-semibold shadow-md"
                    style={{ animationDelay: `${baseDelay + 0.15}s` }}
                >
                    <Zap className="h-3 w-3" />
                    <span className="truncate max-w-[180px]">{event}</span>
                </div>
            </div>

            {/* Consumers column */}
            <div
                className="flex flex-col items-start gap-1 shrink-0"
                style={{ width: "28%", zIndex: 2 }}
            >
                {consumerNames.map((c, i) => (
                    <div
                        key={c}
                        ref={setConsRef(c)}
                        className="flow-node-enter flex items-center gap-1.5 rounded-full border bg-card px-3 py-1 text-xs font-medium shadow-sm"
                        style={{ animationDelay: `${baseDelay + 0.2 + i * 0.08}s` }}
                    >
                        <Network className="h-3 w-3 text-muted-foreground" />
                        <span className="truncate max-w-[120px]">{c}</span>
                    </div>
                ))}
                {consumerNames.length === 0 && (
                    <span className="text-xs text-muted-foreground italic">—</span>
                )}
            </div>
        </div>
    );
}

/* ────────────────────────────────────────────────────────── */
/*  FlowRow – a topic group                                   */
/* ────────────────────────────────────────────────────────── */

function FlowRow({
    flow,
    index,
}: {
    flow: FlowData;
    index: number;
}) {
    return (
        <div
            className="flow-row-enter rounded-lg border bg-card/60 p-4"
            style={{ animationDelay: `${index * 0.12}s` }}
        >
            {/* Topic header */}
            <div className="mb-3 flex items-center gap-2">
                <Layers className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm font-semibold">{flow.topic}</span>
                <Badge variant="outline" className="text-[10px] ml-auto">
                    {flow.flowEvents.length} event{flow.flowEvents.length !== 1 ? "s" : ""}
                </Badge>
            </div>

            {/* Column labels */}
            <div className="flex items-center mb-1 text-[10px] uppercase tracking-wider text-muted-foreground font-medium">
                <div style={{ width: "28%" }} className="text-right pr-4">
                    Producers
                </div>
                <div style={{ width: "44%" }} className="text-center">
                    Event
                </div>
                <div style={{ width: "28%" }} className="pl-4">
                    Consumers
                </div>
            </div>

            {/* Event rows */}
            {flow.flowEvents.map((fe, i) => (
                <EventFlowLine
                    key={fe.name}
                    event={fe.name}
                    producerNames={fe.producers}
                    consumerNames={fe.consumers}
                    rowIndex={i}
                />
            ))}
        </div>
    );
}

/* ────────────────────────────────────────────────────────── */
/*  EventFlowDiagram – main export                            */
/* ────────────────────────────────────────────────────────── */

export function EventFlowDiagram({
    topics,
    hideSelector = false,
}: EventFlowDiagramProps) {
    const allTopicNames = useMemo(
        () => topics.map((t) => t.topic).sort(),
        [topics],
    );
    const [selected, setSelected] = useState<string[]>([]);

    const flows: FlowData[] = useMemo(() => {
        const visible =
            selected.length === 0
                ? topics
                : topics.filter((t) => selected.includes(t.topic));

        if (visible.length === 0) return [];

        return visible.map((t) => {
            const flowEvents: FlowEvent[] = t.events.map((e) => {
                const prods = t.producers
                    .filter((p) => p.event === e.name)
                    .map((p) => p.service);
                const cons = t.consumers
                    .filter((c) => c.event === e.name)
                    .map((c) => c.service);
                return {
                    name: e.name,
                    producers: [...new Set(prods)],
                    consumers: [...new Set(cons)],
                };
            });

            return { topic: t.topic, flowEvents };
        });
    }, [topics, selected]);

    if (topics.length === 0) return null;

    return (
        <Card>
            <CardHeader className="pb-3">
                <div className="flex flex-col sm:flex-row sm:items-center gap-3 justify-between">
                    <CardTitle className="text-base font-semibold flex items-center gap-2">
                        <Zap className="h-4 w-4" />
                        Event Flow
                    </CardTitle>
                    {!hideSelector && (
                        <TopicSelect
                            options={allTopicNames}
                            selected={selected}
                            onChange={setSelected}
                        />
                    )}
                </div>
            </CardHeader>
            <CardContent className="space-y-4">
                {(selected.length === 0 && !hideSelector) ? (
                    <p className="text-center text-sm text-muted-foreground py-8">
                        Select a topic above to visualize its event flow.
                    </p>
                ) : flows.length === 0 ? (
                    <p className="text-center text-sm text-muted-foreground py-8">
                        No topics match the current selection.
                    </p>
                ) : (
                    flows.map((f, i) => (
                        <FlowRow key={f.topic} flow={f} index={i} />
                    ))
                )}
            </CardContent>
        </Card>
    );
}
