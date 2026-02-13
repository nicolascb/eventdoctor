import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
    Database,
    Eye,
    Layers,
    Moon,
    Network,
    PanelLeft,
    RefreshCw,
    Stethoscope,
    Sun,
    Zap,
} from "lucide-react";

export type NavItem = "overview" | "producers" | "topics" | "consumers" | "validator";

interface SidebarProps {
    activeItem: NavItem;
    onNavigate: (item: NavItem) => void;
    counts: {
        producers: number;
        topics: number;
        consumers: number;
    };
    isLoading: boolean;
    onRefresh: () => void;
    collapsed: boolean;
    onToggleCollapse: () => void;
    theme: "light" | "dark";
    onToggleTheme: () => void;
}

const navItems = [
    { id: "overview" as const, label: "Overview", icon: Eye },
    { id: "producers" as const, label: "Producers", icon: Database, countKey: "producers" as const },
    { id: "topics" as const, label: "Topics", icon: Zap, countKey: "topics" as const },
    { id: "consumers" as const, label: "Consumers", icon: Network, countKey: "consumers" as const },
];

const toolItems = [
    { id: "validator" as const, label: "Validator", icon: Stethoscope },
];

export function Sidebar({
    activeItem,
    onNavigate,
    counts,
    isLoading,
    onRefresh,
    collapsed,
    onToggleCollapse,
    theme,
    onToggleTheme,
}: SidebarProps) {
    return (
        <aside
            className={`fixed left-0 top-0 z-40 h-screen border-r bg-sidebar transition-all duration-200 flex flex-col ${collapsed ? "w-[60px]" : "w-[240px]"
                }`}
            style={{ borderColor: "var(--sidebar-border)" }}
        >
            {/* Logo */}
            <div className="flex h-14 items-center px-4 border-b" style={{ borderColor: "var(--sidebar-border)" }}>
                <div className="flex items-center gap-2.5 min-w-0">
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                        <Layers className="h-4 w-4" />
                    </div>
                    {!collapsed && (
                        <span className="text-sm font-semibold tracking-tight truncate" style={{ color: "var(--sidebar-accent-foreground)" }}>
                            EventDoctor
                        </span>
                    )}
                </div>
            </div>

            {/* Navigation */}
            <nav className="flex-1 overflow-y-auto px-3 py-4 space-y-1">
                {!collapsed && (
                    <p className="px-3 mb-2 text-[11px] font-medium uppercase tracking-wider" style={{ color: "var(--sidebar-muted)" }}>
                        Dashboard
                    </p>
                )}
                {navItems.map((item) => {
                    const isActive = activeItem === item.id;
                    const count = item.countKey ? counts[item.countKey] : undefined;
                    return (
                        <button
                            key={item.id}
                            type="button"
                            onClick={() => onNavigate(item.id)}
                            className={`sidebar-nav-item w-full ${isActive ? "active" : ""}`}
                            title={collapsed ? item.label : undefined}
                        >
                            <item.icon className="h-4 w-4 shrink-0" />
                            {!collapsed && (
                                <>
                                    <span className="flex-1 text-left truncate">{item.label}</span>
                                    {count !== undefined && count > 0 && (
                                        <Badge
                                            variant="secondary"
                                            className="ml-auto h-5 min-w-[20px] justify-center px-1.5 text-[10px] font-medium"
                                        >
                                            {count}
                                        </Badge>
                                    )}
                                </>
                            )}
                        </button>
                    );
                })}

                <Separator className="my-3" />

                {!collapsed && (
                    <p className="px-3 mb-2 text-[11px] font-medium uppercase tracking-wider" style={{ color: "var(--sidebar-muted)" }}>
                        Tools
                    </p>
                )}
                {toolItems.map((item) => {
                    const isActive = activeItem === item.id;
                    return (
                        <button
                            key={item.id}
                            type="button"
                            onClick={() => onNavigate(item.id)}
                            className={`sidebar-nav-item w-full ${isActive ? "active" : ""}`}
                            title={collapsed ? item.label : undefined}
                        >
                            <item.icon className="h-4 w-4 shrink-0" />
                            {!collapsed && (
                                <span className="flex-1 text-left truncate">{item.label}</span>
                            )}
                        </button>
                    );
                })}
            </nav>

            {/* Bottom actions */}
            <div className="border-t px-3 py-3 space-y-1" style={{ borderColor: "var(--sidebar-border)" }}>
                <Button
                    variant="ghost"
                    size="sm"
                    onClick={onRefresh}
                    disabled={isLoading}
                    className={`w-full justify-start gap-3 text-sm font-medium ${collapsed ? "px-2" : "px-3"}`}
                    title={collapsed ? "Refresh data" : undefined}
                    style={{ color: "var(--sidebar-foreground)" }}
                >
                    <RefreshCw className={`h-4 w-4 shrink-0 ${isLoading ? "animate-spin" : ""}`} />
                    {!collapsed && <span>Refresh</span>}
                </Button>

                <Button
                    variant="ghost"
                    size="sm"
                    onClick={onToggleTheme}
                    className={`w-full justify-start gap-3 text-sm font-medium ${collapsed ? "px-2" : "px-3"}`}
                    title={collapsed ? (theme === "dark" ? "Light mode" : "Dark mode") : undefined}
                    style={{ color: "var(--sidebar-foreground)" }}
                >
                    {theme === "dark" ? (
                        <Sun className="h-4 w-4 shrink-0" />
                    ) : (
                        <Moon className="h-4 w-4 shrink-0" />
                    )}
                    {!collapsed && <span>{theme === "dark" ? "Light mode" : "Dark mode"}</span>}
                </Button>

                <Button
                    variant="ghost"
                    size="sm"
                    onClick={onToggleCollapse}
                    className={`w-full justify-start gap-3 text-sm font-medium ${collapsed ? "px-2" : "px-3"}`}
                    title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
                    style={{ color: "var(--sidebar-foreground)" }}
                >
                    <PanelLeft className={`h-4 w-4 shrink-0 transition-transform ${collapsed ? "rotate-180" : ""}`} />
                    {!collapsed && <span>Collapse</span>}
                </Button>
            </div>
        </aside>
    );
}
