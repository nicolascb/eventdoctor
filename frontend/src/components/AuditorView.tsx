import { UndocumentedConsumersView } from "./UndocumentedConsumersView";
import { AlertTriangle, ShieldAlert } from "lucide-react";

export function AuditorView() {
    return (
        <div className="space-y-6 animate-in fade-in">
            <div className="flex flex-col gap-2">
                <h2 className="text-2xl font-bold tracking-tight flex items-center gap-2">
                    <ShieldAlert className="h-6 w-6 text-amber-500" />
                    Governance Auditor
                </h2>
                <p className="text-muted-foreground">
                    Review and resolve infrastructure anomalies, undocumented events, and ghost consumers to maintain a healthy event catalog.
                </p>
            </div>
            
            <div className="rounded-lg border border-amber-500/20 bg-amber-500/5 p-6 space-y-4">
                <div className="flex items-center gap-2 text-amber-700 dark:text-amber-400 font-medium">
                    <AlertTriangle className="h-5 w-5" />
                    Undocumented Consumers
                </div>
                <p className="text-sm text-amber-700/80 dark:text-amber-400/80 mb-4">
                    The following consumer groups were detected consuming from Kafka topics, but are not properly documented in your <code>eventdoctor.yaml</code> contract.
                </p>
                <div className="bg-background rounded-md border shadow-sm">
                    <UndocumentedConsumersView />
                </div>
            </div>
        </div>
    );
}
