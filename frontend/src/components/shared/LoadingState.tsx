import { Card, CardContent } from '@/components/ui/card';
import { Loader2 } from 'lucide-react';

interface LoadingStateProps {
    message?: string;
}

export function LoadingState({ message = 'Loading data...' }: LoadingStateProps) {
    return (
        <Card className="animate-in">
            <CardContent className="flex items-center justify-center py-16">
                <div className="flex flex-col items-center gap-3 text-muted-foreground">
                    <Loader2 className="h-8 w-8 animate-spin" />
                    <span className="text-sm">{message}</span>
                </div>
            </CardContent>
        </Card>
    );
}
