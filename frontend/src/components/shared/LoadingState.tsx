import { Card, CardContent } from '@/components/ui/card';
import { Loader2 } from 'lucide-react';

interface LoadingStateProps {
    message?: string;
}

export function LoadingState({ message = 'Loading data...' }: LoadingStateProps) {
    return (
        <Card className="animate-in">
            <CardContent className="flex items-center justify-center py-12">
                <div className="flex flex-col items-center gap-3 text-muted-foreground">
                    <Loader2 className="h-5 w-5 animate-spin" />
                    <span className="text-xs">{message}</span>
                </div>
            </CardContent>
        </Card>
    );
}
