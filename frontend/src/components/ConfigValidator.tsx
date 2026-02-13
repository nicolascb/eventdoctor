import { StatCard } from "@/components/shared";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { AlertCircle, CheckCircle, ClipboardPaste, FileCode, Loader2, RotateCcw, Shield, XCircle } from "lucide-react";
import { useCallback, useState } from 'react';

interface ValidationResult {
    valid: boolean;
    message: string;
    details?: {
        producers?: number;
        consumers?: number;
        events?: number;
        errors?: string[];
    };
}

export function ConfigValidator() {
    const [config, setConfig] = useState('');
    const [validating, setValidating] = useState(false);
    const [result, setResult] = useState<ValidationResult | null>(null);

    const validateConfig = async () => {
        if (!config.trim()) {
            setResult({
                valid: false,
                message: 'Configuration is required for validation'
            });
            return;
        }

        setValidating(true);
        setResult(null);

        await new Promise(resolve => setTimeout(resolve, 600));

        try {
            const lines = config.split('\n');
            let producerCount = 0;
            let consumerCount = 0;
            let eventCount = 0;
            const errors: string[] = [];

            lines.forEach(line => {
                if (line.trim().startsWith('- topic:')) producerCount++;
                if (line.trim().startsWith('- group:')) consumerCount++;
                if (line.trim().startsWith('- name:')) eventCount++;
            });

            if (!config.includes('version:')) {
                errors.push('Missing required field: version');
            }
            if (!config.includes('service:')) {
                errors.push('Missing required field: service');
            }
            if (producerCount === 0) {
                errors.push('No producers defined in configuration');
            }
            if (consumerCount === 0) {
                errors.push('No consumers defined in configuration');
            }
            if (eventCount === 0) {
                errors.push('No events defined in configuration');
            }

            const isValid = errors.length === 0;

            setResult({
                valid: isValid,
                message: isValid
                    ? `Configuration is valid — ${producerCount} producers, ${consumerCount} consumers, ${eventCount} events detected.`
                    : `Found ${errors.length} validation ${errors.length === 1 ? 'error' : 'errors'}`,
                details: {
                    producers: producerCount,
                    consumers: consumerCount,
                    events: eventCount,
                    errors: errors.length > 0 ? errors : undefined
                }
            });
        } catch (error) {
            setResult({
                valid: false,
                message: 'Failed to parse configuration. Please check YAML syntax.',
                details: {
                    errors: [error instanceof Error ? error.message : 'Unknown parsing error']
                }
            });
        } finally {
            setValidating(false);
        }
    };

    const sampleConfig = `version: "2.1.7"
service: "my-event-service"

config:
  servers:
    - environment: "production"
      url: "kafka://localhost:9092"
  repository: "https://github.com/my-org/event-schemas"

producers:
  - topic: "events.user.lifecycle"
    title: "User Lifecycle Events"
    owner: true
    writes: true
    description: "User creation, update, and deletion events"
    events:
      - "user.created"
      - "user.updated"
      - "user.deleted"

  - topic: "events.system.monitoring"
    title: "System Monitoring"
    owner: true
    writes: true
    description: "System health and performance metrics"
    events:
      - "system.health.check"
      - "performance.metrics"

consumers:
  - group: "analytics-processor"
    description: "Process user events for analytics"
    topics:
      - name: "events.user.lifecycle"
        events:
          - "user.created"
          - "user.updated"

  - group: "monitoring-service"
    description: "Monitor system health and alerts"
    topics:
      - name: "events.system.monitoring"
        events:
          - "system.health.check"

events:
  - name: "user.created"
    version: "1.0"
    description: "Triggered when a new user is created"
    schema_url: "https://schemas.example.com/user/created/v1.0"
    headers:
      source: "string"
      timestamp: "datetime"
    properties:
      user_id: "uuid"
      email: "string"
      created_at: "datetime"

  - name: "system.health.check"
    version: "1.0"
    description: "System health monitoring event"
    headers:
      service_name: "string"
      timestamp: "datetime"
    properties:
      status: "string"
      cpu_usage: "number"
      memory_usage: "number"`;

    const handleLoadSample = () => {
        setConfig(sampleConfig);
        setResult(null);
    };

    const handleClear = useCallback(() => {
        setConfig('');
        setResult(null);
    }, []);

    return (
        <div className="space-y-6 animate-in">
            {/* Stats Cards — only shown after validation */}
            {result?.details && result.valid && (
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <StatCard
                        label="Producers"
                        value={result.details.producers ?? 0}
                        description="Detected in config"
                        icon={<FileCode className="h-4 w-4 text-muted-foreground" />}
                    />
                    <StatCard
                        label="Consumers"
                        value={result.details.consumers ?? 0}
                        description="Detected in config"
                        icon={<Shield className="h-4 w-4 text-muted-foreground" />}
                    />
                    <StatCard
                        label="Events"
                        value={result.details.events ?? 0}
                        description="Detected in config"
                        icon={<Shield className="h-4 w-4 text-muted-foreground" />}
                    />
                </div>
            )}

            {/* Main content */}
            <div className="grid gap-4 lg:grid-cols-5">
                {/* Input panel — takes 3/5 */}
                <Card className="lg:col-span-3">
                    <CardHeader>
                        <CardTitle className="text-sm font-medium">Configuration Input</CardTitle>
                        <CardDescription className="text-xs">
                            Paste your YAML configuration below to validate
                        </CardDescription>
                    </CardHeader>

                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="config" className="text-xs font-medium">
                                YAML Configuration
                            </Label>
                            <Textarea
                                id="config"
                                name="config"
                                placeholder={`# Paste your EventDoctor configuration here...\n# Example:\n# version: "2.1.7"\n# service: my-service\n# producers:\n#   - topic: events.user\n#     ...`}
                                value={config}
                                onChange={(e) => {
                                    setConfig(e.target.value);
                                    if (result) setResult(null);
                                }}
                                className="min-h-[320px] font-mono text-xs resize-y"
                                spellCheck={false}
                                autoComplete="off"
                            />
                        </div>

                        <div className="flex gap-3">
                            <Button
                                onClick={validateConfig}
                                disabled={validating || !config.trim()}
                                className="flex-1 gap-2"
                            >
                                {validating ? (
                                    <>
                                        <Loader2 className="h-4 w-4 animate-spin" />
                                        Validating...
                                    </>
                                ) : (
                                    <>
                                        <CheckCircle className="h-4 w-4" />
                                        Validate
                                    </>
                                )}
                            </Button>

                            <Button
                                variant="outline"
                                onClick={handleLoadSample}
                                className="gap-2"
                            >
                                <ClipboardPaste className="h-4 w-4" />
                                Load Sample
                            </Button>

                            {config && (
                                <Button
                                    variant="ghost"
                                    onClick={handleClear}
                                    className="gap-2 text-muted-foreground"
                                >
                                    <RotateCcw className="h-4 w-4" />
                                    Clear
                                </Button>
                            )}
                        </div>

                        {/* Validation Result */}
                        {result && (
                            <Card className={`${result.valid
                                ? 'border-foreground/20 bg-muted/50'
                                : 'border-destructive/30 bg-destructive/5'
                                }`}>
                                <CardContent className="pt-4">
                                    <div className="flex items-start gap-3">
                                        {result.valid ? (
                                            <CheckCircle className="h-5 w-5 text-foreground flex-shrink-0 mt-0.5" />
                                        ) : (
                                            <XCircle className="h-5 w-5 text-destructive flex-shrink-0 mt-0.5" />
                                        )}
                                        <div className="flex-1">
                                            <h4 className="font-medium text-sm mb-0.5">
                                                {result.valid ? 'Validation Passed' : 'Validation Failed'}
                                            </h4>
                                            <p className="text-xs text-muted-foreground">
                                                {result.message}
                                            </p>
                                        </div>
                                    </div>

                                    {result.details?.errors && result.details.errors.length > 0 && (
                                        <div className="space-y-1.5 mt-3 ml-8">
                                            {result.details.errors.map((error, index) => (
                                                <div key={index} className="flex items-start gap-2 text-xs">
                                                    <AlertCircle className="h-3.5 w-3.5 text-destructive mt-0.5 flex-shrink-0" />
                                                    <span className="text-muted-foreground">{error}</span>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </CardContent>
                            </Card>
                        )}
                    </CardContent>
                </Card>

                {/* Reference panel — takes 2/5 */}
                <Card className="lg:col-span-2">
                    <CardHeader>
                        <CardTitle className="text-sm font-medium">Reference</CardTitle>
                        <CardDescription className="text-xs">
                            Required structure and fields
                        </CardDescription>
                    </CardHeader>

                    <CardContent className="space-y-4">
                        <Card className="border">
                            <CardContent className="pt-4">
                                <div className="flex items-center gap-2 mb-3">
                                    <FileCode className="h-3.5 w-3.5 text-muted-foreground" />
                                    <span className="text-xs font-medium">Required Fields</span>
                                </div>
                                <div className="space-y-2 text-xs text-muted-foreground">
                                    <div className="flex items-start gap-2">
                                        <code className="text-[10px] bg-muted px-1.5 py-0.5 rounded font-mono mt-0.5">version</code>
                                        <span>Spec version (e.g. "2.1.7")</span>
                                    </div>
                                    <div className="flex items-start gap-2">
                                        <code className="text-[10px] bg-muted px-1.5 py-0.5 rounded font-mono mt-0.5">service</code>
                                        <span>Service name identifier</span>
                                    </div>
                                    <div className="flex items-start gap-2">
                                        <code className="text-[10px] bg-muted px-1.5 py-0.5 rounded font-mono mt-0.5">producers</code>
                                        <span>At least one producer with topic</span>
                                    </div>
                                    <div className="flex items-start gap-2">
                                        <code className="text-[10px] bg-muted px-1.5 py-0.5 rounded font-mono mt-0.5">consumers</code>
                                        <span>At least one consumer group</span>
                                    </div>
                                    <div className="flex items-start gap-2">
                                        <code className="text-[10px] bg-muted px-1.5 py-0.5 rounded font-mono mt-0.5">events</code>
                                        <span>At least one event type definition</span>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        <Card className="border">
                            <CardContent className="pt-4">
                                <div className="flex items-center gap-2 mb-3">
                                    <FileCode className="h-3.5 w-3.5 text-muted-foreground" />
                                    <span className="text-xs font-medium">Sections Guide</span>
                                </div>
                                <div className="space-y-2 text-xs text-muted-foreground">
                                    <div className="flex items-start gap-2">
                                        <span className="text-[10px] font-bold bg-muted rounded px-1.5 py-0.5 mt-0.5">P</span>
                                        <div>
                                            <span className="font-medium text-foreground">Producers:</span> Services that publish events to topics
                                        </div>
                                    </div>
                                    <div className="flex items-start gap-2">
                                        <span className="text-[10px] font-bold bg-muted rounded px-1.5 py-0.5 mt-0.5">C</span>
                                        <div>
                                            <span className="font-medium text-foreground">Consumers:</span> Services that subscribe to topics
                                        </div>
                                    </div>
                                    <div className="flex items-start gap-2">
                                        <span className="text-[10px] font-bold bg-muted rounded px-1.5 py-0.5 mt-0.5">E</span>
                                        <div>
                                            <span className="font-medium text-foreground">Events:</span> Event type schemas and metadata
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
