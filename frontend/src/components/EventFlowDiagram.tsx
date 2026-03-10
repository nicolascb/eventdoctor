import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import type { TopicView } from "@/types";
import {
    Check,
    ChevronDown,
    ChevronUp,
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

interface TopicFlowData {
    topic: string;
    events: string[];
    producers: string[];
    consumers: string[];
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
/*  TopicFlowRow – a topic group                              */
/* ────────────────────────────────────────────────────────── */

function TopicFlowRow({
    flow,
    index,
}: {
    flow: TopicFlowData;
    index: number;
}) {
    const containerRef = useRef<HTMLDivElement>(null);
    const topicRef = useRef<HTMLDivElement>(null);
    const prodRefs = useRef<Map<string, HTMLDivElement>>(new Map());
    const consRefs = useRef<Map<string, HTMLDivElement>>(new Map());
    const [lines, setLines] = useState<ConnLine[]>([]);
    const [expanded, setExpanded] = useState(false);

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
        const topicEl = topicRef.current;
        if (!container || !topicEl) return;

        const cr = container.getBoundingClientRect();
        const tr = topicEl.getBoundingClientRect();

        // Topic block connection points
        const ty = tr.top - cr.top + tr.height / 2;

        const result: ConnLine[] = [];
        let idx = 0;

        // producer → topic block
        prodRefs.current.forEach((el, name) => {
            const r = el.getBoundingClientRect();
            const sx = r.right - cr.left;
            const sy = r.top - cr.top + r.height / 2;
            const tLeft = tr.left - cr.left;
            const dx = (tLeft - sx) * 0.45;
            const path = `M ${sx} ${sy} C ${sx + dx} ${sy}, ${tLeft - dx} ${ty}, ${tLeft} ${ty}`;
            result.push({
                id: `p-${name}-${idx}`,
                path,
                color: "var(--muted-foreground)",
                delay: idx * 0.3,
            });
            idx++;
        });

        // topic block → consumer
        consRefs.current.forEach((el, name) => {
            const r = el.getBoundingClientRect();
            const cx = r.left - cr.left;
            const cy = r.top - cr.top + r.height / 2;
            const tRight = tr.right - cr.left;
            const dx = (cx - tRight) * 0.45;
            const path = `M ${tRight} ${ty} C ${tRight + dx} ${ty}, ${cx - dx} ${cy}, ${cx} ${cy}`;
            result.push({
                id: `c-${name}-${idx}`,
                path,
                color: "var(--muted-foreground)",
                delay: idx * 0.3,
            });
            idx++;
        });

        setLines(result);
    }, []);

    useLayoutEffect(() => {
        const delay = index * 100 + 400;
        const timer = setTimeout(computeLines, delay);
        return () => clearTimeout(timer);
    }, [computeLines, index]);

    // Recalculate on resize
    useLayoutEffect(() => {
        const container = containerRef.current;
        if (!container) return;
        const ro = new ResizeObserver(() => computeLines());
        ro.observe(container);
        return () => ro.disconnect();
    }, [computeLines]);

    const baseDelay = index * 0.12;
    const MAX_EVENTS_VISIBLE = 4;
    const hasMoreEvents = flow.events.length > MAX_EVENTS_VISIBLE;
    const visibleEvents = expanded ? flow.events : flow.events.slice(0, MAX_EVENTS_VISIBLE);

    return (
        <div
            ref={containerRef}
            className="flow-row-enter rounded-lg border bg-card/60 p-4 relative"
            style={{ animationDelay: `${index * 0.12}s` }}
        >
            {/* SVG overlay */}
            <svg
                className="pointer-events-none absolute inset-0 h-full w-full"
                style={{ zIndex: 1 }}
            >
                {lines.map((l) => (
                    <g key={l.id}>
                        {/* Animated flowing dashed line */}
                        <path
                            d={l.path}
                            fill="none"
                            stroke={l.color}
                            strokeWidth={1.5}
                            strokeDasharray="6 4"
                            opacity={0.4}
                        >
                            <animate
                                attributeName="stroke-dashoffset"
                                values="10;0"
                                dur="1.2s"
                                repeatCount="indefinite"
                                calcMode="linear"
                            />
                        </path>

                        {/* Moving glowing particle (represents event payload) */}
                        <g opacity={0}>
                            {/* Inner dot */}
                            <circle r="3.5" className="fill-primary" />
                            {/* Outer glow aura */}
                            <circle r="8" className="fill-primary" opacity="0.25">
                                <animate
                                    attributeName="r"
                                    values="6;9;6"
                                    dur="1.5s"
                                    repeatCount="indefinite"
                                />
                            </circle>

                            {/* Ease in/out motion path */}
                            <animateMotion
                                path={l.path}
                                dur="2.5s"
                                repeatCount="indefinite"
                                begin={`${l.delay}s`}
                                keyPoints="0;1"
                                keyTimes="0;1"
                                calcMode="spline"
                                keySplines="0.4 0 0.2 1"
                            />
                            <animate
                                attributeName="opacity"
                                values="0;1;1;0"
                                keyTimes="0;0.1;0.9;1"
                                dur="2.5s"
                                repeatCount="indefinite"
                                begin={`${l.delay}s`}
                            />
                        </g>
                    </g>
                ))}
            </svg>

            <div className="relative z-10 w-full">
                {/* Column labels */}
                <div className="flex items-center mb-4 text-[10px] uppercase tracking-wider text-muted-foreground font-medium">
                    <div style={{ width: "25%" }} className="text-right pr-4">
                        Producers
                    </div>
                    <div style={{ width: "50%" }} className="text-center">
                        Topic & Events
                    </div>
                    <div style={{ width: "25%" }} className="pl-4">
                        Consumers
                    </div>
                </div>

                <div className="flex items-stretch gap-0 w-full min-h-[100px]">
                    {/* Producers column */}
                    <div
                        className="flex flex-col items-end justify-center gap-2 shrink-0 py-4"
                        style={{ width: "25%" }}
                    >
                        {flow.producers.map((p, i) => (
                            <div
                                key={p}
                                ref={setProdRef(p)}
                                className="flow-node-enter flex items-center gap-1.5 rounded-full border bg-card px-3 py-1.5 text-xs font-medium shadow-sm transition-all hover:border-primary/50 hover:shadow-md"
                                style={{ animationDelay: `${baseDelay + i * 0.08}s` }}
                            >
                                <Database className="h-3.5 w-3.5 text-muted-foreground" />
                                <span className="truncate max-w-[120px]">{p}</span>
                            </div>
                        ))}
                        {flow.producers.length === 0 && (
                            <span className="text-xs text-muted-foreground italic px-3 py-1.5">—</span>
                        )}
                    </div>

                    {/* Center: Topic Block */}
                    <div
                        className="flex flex-col items-center justify-center shrink-0 px-6 py-2"
                        style={{ width: "50%" }}
                    >
                        <div
                            ref={topicRef}
                            className="flow-node-enter flex flex-col w-full max-w-[340px] rounded-xl border-2 border-primary/20 bg-card shadow-sm overflow-hidden transition-all duration-300 hover:border-primary/40"
                            style={{ animationDelay: `${baseDelay + 0.15}s` }}
                        >
                            {/* Topic Header */}
                            <div className="bg-primary/5 px-4 py-2.5 flex items-center gap-2 border-b border-primary/10">
                                <Layers className="h-4 w-4 text-primary" />
                                <span className="text-sm font-semibold text-foreground truncate">{flow.topic}</span>
                            </div>

                            {/* Events list */}
                            <div className="p-3.5 flex flex-wrap gap-1.5 justify-center bg-card/30">
                                {visibleEvents.map((e) => (
                                    <Badge
                                        key={e}
                                        variant="secondary"
                                        className="font-mono text-[10px] flex items-center gap-1 bg-secondary/40 hover:bg-secondary/70 border-secondary-foreground/10 px-2 font-medium"
                                    >
                                        <Zap className="h-2.5 w-2.5 opacity-60" />
                                        {e}
                                    </Badge>
                                ))}
                                {visibleEvents.length === 0 && (
                                    <span className="text-xs text-muted-foreground italic">No events mapped</span>
                                )}
                            </div>

                            {/* Expand toggle */}
                            {hasMoreEvents && (
                                <button
                                    className="w-full py-2 text-[10px] font-medium text-muted-foreground bg-muted/20 hover:bg-muted/50 hover:text-foreground transition-colors flex items-center justify-center gap-1 border-t border-border focus:outline-none"
                                    onClick={() => setExpanded(!expanded)}
                                >
                                    {expanded ? (
                                        <>Hide {flow.events.length - MAX_EVENTS_VISIBLE} events <ChevronUp className="h-3 w-3" /></>
                                    ) : (
                                        <>Show {flow.events.length - MAX_EVENTS_VISIBLE} more events <ChevronDown className="h-3 w-3" /></>
                                    )}
                                </button>
                            )}
                        </div>
                    </div>

                    {/* Consumers column */}
                    <div
                        className="flex flex-col items-start justify-center gap-2 shrink-0 py-4"
                        style={{ width: "25%" }}
                    >
                        {flow.consumers.map((c, i) => (
                            <div
                                key={c}
                                ref={setConsRef(c)}
                                className="flow-node-enter flex items-center gap-1.5 rounded-full border bg-card px-3 py-1.5 text-xs font-medium shadow-sm transition-all hover:border-primary/50 hover:shadow-md"
                                style={{ animationDelay: `${baseDelay + 0.2 + i * 0.08}s` }}
                            >
                                <Network className="h-3.5 w-3.5 text-muted-foreground" />
                                <span className="truncate max-w-[120px]">{c}</span>
                            </div>
                        ))}
                        {flow.consumers.length === 0 && (
                            <span className="text-xs text-muted-foreground italic px-3 py-1.5">—</span>
                        )}
                    </div>
                </div>
            </div>
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

    const flows: TopicFlowData[] = useMemo(() => {
        const visible =
            selected.length === 0
                ? topics
                : topics.filter((t) => selected.includes(t.topic));

        if (visible.length === 0) return [];

        return visible.map((t) => {
            const events = t.events.map(e => e.name);
            const producers = [...new Set(t.producers.map(p => p.service))];
            const consumers = [...new Set(t.consumers.map(c => c.service))];

            return { topic: t.topic, events, producers, consumers };
        });
    }, [topics, selected]);

    return (
        <Card>
            <CardHeader className="pb-3 border-b border-border/50 bg-muted/5">
                <div className="flex flex-col sm:flex-row sm:items-center gap-3 justify-between">
                    <CardTitle className="text-base font-semibold flex items-center gap-2">
                        <Zap className="h-4 w-4 text-primary" />
                        Architecture Flow
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
            <CardContent className="space-y-6 pt-6">
                {((!hideSelector && selected.length === 0) || (hideSelector && topics.length === 0)) ? (
                    <div className="h-full min-h-[500px] w-full rounded-2xl border-2 border-dashed border-border/60 flex flex-col items-center justify-center p-12 text-center bg-gradient-to-b from-muted/5 to-muted/20">
                        <div className="relative mb-6">
                            <div className="absolute inset-0 animate-ping rounded-full bg-primary/10 opacity-75" />
                            <div className="relative rounded-full bg-background border shadow-sm p-6">
                                <Layers className="h-10 w-10 text-primary/60" />
                            </div>
                        </div>
                        <h3 className="text-xl font-semibold text-foreground tracking-tight mb-2">Visualize Everything</h3>
                        <p className="max-w-[420px] text-sm text-muted-foreground leading-relaxed">
                            No topics selected yet. Pick one or multiple topics from the list on the left to instantly map and visualize their flow between producers and consumers.
                        </p>
                    </div>
                ) : flows.length === 0 ? (
                    <p className="text-center text-sm text-muted-foreground py-8">
                        No topics match the current selection.
                    </p>
                ) : (
                    flows.map((f, i) => (
                        <TopicFlowRow key={f.topic} flow={f} index={i} />
                    ))
                )}
            </CardContent>
        </Card>
    );
}
