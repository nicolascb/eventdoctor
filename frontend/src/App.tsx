import { ConfigValidator } from '@/components/ConfigValidator';
import { ConsumersView } from '@/components/ConsumersView';
import { OverviewView } from '@/components/OverviewView';
import { ProducersView } from '@/components/ProducersView';
import { ErrorState, LoadingState } from '@/components/shared';
import { TopicsView } from '@/components/TopicsView';
import { Badge } from "@/components/ui/badge";
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useConsumers } from '@/hooks/useConsumers';
import { useEvents } from '@/hooks/useEvents';
import { useOverview } from '@/hooks/useOverview';
import { useProducers } from '@/hooks/useProducers';
import { Database, Eye, Network, RefreshCw, Stethoscope, Zap } from "lucide-react";

// Get the API URL from environment variable
const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8087/v1';

function App() {
  const { producers, loading: producersLoading, error: producersError, refetch: refetchProducers } = useProducers();
  const { topics, loading: eventsLoading, error: eventsError, refetch: refetchEvents } = useEvents();
  const { consumers, loading: consumersLoading, error: consumersError, refetch: refetchConsumers } = useConsumers();
  const { overview, loading: overviewLoading, error: overviewError, refetch: refetchOverview } = useOverview();

  const isLoading = producersLoading || eventsLoading || consumersLoading || overviewLoading;
  const hasError = producersError || eventsError || consumersError || overviewError;

  const handleRetry = () => {
    refetchProducers();
    refetchEvents();
    refetchConsumers();
    refetchOverview();
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Sticky header */}
      <header className="sticky top-0 z-50 border-b border-border bg-card/95 backdrop-blur supports-[backdrop-filter]:bg-card/60">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="p-2.5 bg-primary rounded-lg">
                <Stethoscope className="h-5 w-5 text-primary-foreground" />
              </div>
              <div>
                <h1 className="text-xl font-bold text-foreground leading-none">EventDoctor</h1>
                <p className="text-xs text-muted-foreground mt-0.5">Event-driven architecture monitoring</p>
              </div>
            </div>

            <div className="flex items-center gap-3">
              {!isLoading && !hasError && (
                <div className="hidden sm:flex gap-2">
                  <Badge variant="secondary" className="h-7 px-2.5 text-xs">
                    {producers.length} Producers
                  </Badge>
                  <Badge variant="secondary" className="h-7 px-2.5 text-xs">
                    {topics.length} Topics
                  </Badge>
                  <Badge variant="secondary" className="h-7 px-2.5 text-xs">
                    {consumers.length} Consumers
                  </Badge>
                </div>
              )}
              <Button
                variant="outline"
                size="sm"
                onClick={handleRetry}
                disabled={isLoading}
                className="gap-2"
              >
                <RefreshCw className={`h-3.5 w-3.5 ${isLoading ? 'animate-spin' : ''}`} />
                <span className="hidden sm:inline">Refresh</span>
              </Button>
            </div>
          </div>
        </div>
      </header>

      <div className="container mx-auto px-6 py-8">
        {/* Loading State */}
        {isLoading && <LoadingState message="Loading EventDoctor data..." />}

        {/* Error State */}
        {hasError && (
          <ErrorState
            message={`Make sure the EventDoctor API is running at ${API_URL}.`}
            details={{
              Producers: producersError,
              Events: eventsError,
              Consumers: consumersError,
            }}
            onRetry={handleRetry}
          />
        )}

        {/* Main Interface */}
        {!isLoading && !hasError && (
          <div className="animate-in">
            <Tabs defaultValue="overview" className="space-y-8">
              <TabsList className="grid w-full grid-cols-5">
                <TabsTrigger value="overview" className="flex items-center gap-2">
                  <Eye className="h-4 w-4" />
                  Overview
                </TabsTrigger>
                <TabsTrigger value="producers" className="flex items-center gap-2">
                  <Database className="h-4 w-4" />
                  Producers
                  {producers.length > 0 && (
                    <Badge variant="secondary" className="ml-1 h-5 px-1.5 text-[10px] font-medium">
                      {producers.length}
                    </Badge>
                  )}
                </TabsTrigger>
                <TabsTrigger value="events" className="flex items-center gap-2">
                  <Zap className="h-4 w-4" />
                  Topics
                  {topics.length > 0 && (
                    <Badge variant="secondary" className="ml-1 h-5 px-1.5 text-[10px] font-medium">
                      {topics.length}
                    </Badge>
                  )}
                </TabsTrigger>
                <TabsTrigger value="consumers" className="flex items-center gap-2">
                  <Network className="h-4 w-4" />
                  Consumers
                  {consumers.length > 0 && (
                    <Badge variant="secondary" className="ml-1 h-5 px-1.5 text-[10px] font-medium">
                      {consumers.length}
                    </Badge>
                  )}
                </TabsTrigger>
                <TabsTrigger value="validator" className="flex items-center gap-2">
                  <Stethoscope className="h-4 w-4" />
                  Validator
                </TabsTrigger>
              </TabsList>

              <TabsContent value="overview" className="space-y-8">
                <OverviewView overview={overview} />
              </TabsContent>

              <TabsContent value="producers">
                <ProducersView producers={producers} />
              </TabsContent>

              <TabsContent value="events">
                <TopicsView topics={topics} producers={producers} consumers={consumers} />
              </TabsContent>

              <TabsContent value="consumers">
                <ConsumersView consumers={consumers} />
              </TabsContent>

              <TabsContent value="validator">
                <ConfigValidator />
              </TabsContent>
            </Tabs>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
