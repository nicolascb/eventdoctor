import { Button } from '@/components/ui/button';
import { AlertTriangle, RefreshCw } from 'lucide-react';
import { Component, type ErrorInfo, type ReactNode } from 'react';

interface ErrorBoundaryProps {
    children: ReactNode;
    fallback?: ReactNode;
}

interface ErrorBoundaryState {
    hasError: boolean;
    error?: Error;
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
    state: ErrorBoundaryState = { hasError: false };

    static getDerivedStateFromError(error: Error): ErrorBoundaryState {
        return { hasError: true, error };
    }

    componentDidCatch(error: Error, info: ErrorInfo) {
        console.error('ErrorBoundary caught:', error, info);
    }

    handleReset = () => {
        this.setState({ hasError: false, error: undefined });
    };

    render() {
        if (this.state.hasError) {
            if (this.props.fallback) {
                return this.props.fallback;
            }

            return (
                <div className="flex min-h-screen items-center justify-center bg-background">
                    <div className="text-center max-w-md mx-auto px-6">
                        <div className="mb-4">
                            <AlertTriangle className="h-6 w-6 text-destructive mx-auto" />
                        </div>
                        <h2 className="text-base font-semibold text-foreground mb-2">
                            Something went wrong
                        </h2>
                        <p className="text-sm text-muted-foreground mb-6">
                            {this.state.error?.message ?? 'An unexpected error occurred'}
                        </p>
                        <div className="flex gap-3 justify-center">
                            <Button
                                variant="outline"
                                onClick={this.handleReset}
                                className="gap-2"
                            >
                                <RefreshCw className="h-4 w-4" />
                                Try again
                            </Button>
                            <Button
                                variant="default"
                                onClick={() => window.location.reload()}
                            >
                                Reload page
                            </Button>
                        </div>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}
