import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { AlertTriangle, RefreshCw } from 'lucide-react';

interface ErrorStateProps {
    title?: string;
    message: string;
    details?: Record<string, string | null>;
    onRetry?: () => void;
}

export function ErrorState({
    title = 'Connection Error',
    message,
    details,
    onRetry,
}: ErrorStateProps) {
    return (
        <Card className="border-destructive/20 animate-in">
            <CardHeader>
                <CardTitle className="flex items-center gap-2 text-destructive">
                    <AlertTriangle className="h-5 w-5" />
                    {title}
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
                <p className="text-sm text-muted-foreground">{message}</p>

                {details && (
                    <div className="space-y-1">
                        {Object.entries(details).map(
                            ([key, value]) =>
                                value && (
                                    <p key={key} className="text-sm text-destructive">
                                        {key}: {value}
                                    </p>
                                ),
                        )}
                    </div>
                )}

                {onRetry && (
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={onRetry}
                        className="gap-2 mt-2"
                    >
                        <RefreshCw className="h-4 w-4" />
                        Try again
                    </Button>
                )}
            </CardContent>
        </Card>
    );
}
