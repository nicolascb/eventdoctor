import { ConsumersView } from "@/components/ConsumersView";
import { UndocumentedConsumersView } from "@/components/UndocumentedConsumersView";
import { AlertTriangle, Network } from "lucide-react";
import { useState } from "react";

type ConsumerTab = "all" | "undocumented";

const tabs = [
    { id: "all" as const, label: "All Consumers", icon: Network },
    { id: "undocumented" as const, label: "Undocumented", icon: AlertTriangle },
];

export function ConsumersPage() {
    const [activeTab, setActiveTab] = useState<ConsumerTab>("all");

    return (
        <div className="space-y-4">
            {/* Sub-navigation tabs */}
            <div className="flex items-center gap-1 border-b border-border">
                {tabs.map((tab) => (
                    <button
                        key={tab.id}
                        type="button"
                        onClick={() => setActiveTab(tab.id)}
                        className={`inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium transition-colors border-b-2 -mb-px ${
                            activeTab === tab.id
                                ? "border-primary text-foreground"
                                : "border-transparent text-muted-foreground hover:text-foreground"
                        }`}
                    >
                        <tab.icon className="h-3.5 w-3.5" />
                        {tab.label}
                    </button>
                ))}
            </div>

            {/* Tab content */}
            {activeTab === "all" && <ConsumersView />}
            {activeTab === "undocumented" && <UndocumentedConsumersView />}
        </div>
    );
}
