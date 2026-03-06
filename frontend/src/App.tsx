import { ConfigValidator } from '@/components/ConfigValidator';
import { EventsView } from '@/components/EventsView';
import { OverviewView } from '@/components/OverviewView';
import { ProducersView } from '@/components/ProducersView';
import { ErrorState, LoadingState } from '@/components/shared';
import { Sidebar, type NavItem } from '@/components/Sidebar';
import { TopicsView } from '@/components/TopicsView';
import { ConsumersView } from '@/components/ConsumersView';
import { Separator } from '@/components/ui/separator';
import { useOverview } from '@/hooks/useOverview';
import { useTheme } from '@/hooks/useTheme';
import { useState } from 'react';

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8087/v1';

const pageConfig: Record<NavItem, { title: string; description: string }> = {
  overview: { title: 'Overview', description: 'Monitor your event-driven architecture at a glance.' },
  producers: { title: 'Producers', description: 'Services publishing events to topics.' },
  events: { title: 'Events', description: 'Registered event types in the system.' },
  topics: { title: 'Topics', description: 'Message channels and their event types.' },
  consumers: { title: 'Consumers', description: 'Services consuming events from topics.' },
  validator: { title: 'Validator', description: 'Validate your EventDoctor YAML configuration.' },
  auditor: { title: 'Auditor', description: 'Auditor consumer status.' },
};

function App() {
  const [activeView, setActiveView] = useState<NavItem>('overview');
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const { theme, toggleTheme } = useTheme();

  const { overview, loading: overviewLoading, error: overviewError, refetch: refetchOverview } = useOverview();

  const renderContent = () => {
    if (overviewLoading) {
      return <LoadingState message="Loading EventDoctor data..." />;
    }

    if (overviewError) {
      return (
        <ErrorState
          message={`Make sure the EventDoctor API is running at ${API_URL}.`}
          details={{ Overview: overviewError }}
          onRetry={refetchOverview}
        />
      );
    }

    switch (activeView) {
      case 'overview':
        return <OverviewView overview={overview} />;
      case 'producers':
        return <ProducersView />;
      case 'events':
        return <EventsView />;
      case 'topics':
        return <TopicsView />;
      case 'consumers':
        return <ConsumersView />;
      case 'validator':
        return <ConfigValidator />;
      case 'auditor':
        return <div>Auditor view coming soon...</div>;
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <Sidebar
        activeItem={activeView}
        onNavigate={setActiveView}
        counts={{
          producers: overview.total_producers,
          events: overview.total_events,
          topics: overview.total_topics,
          consumers: overview.total_consumers,
        }}
        isLoading={overviewLoading}
        onRefresh={refetchOverview}
        collapsed={sidebarCollapsed}
        onToggleCollapse={() => setSidebarCollapsed(!sidebarCollapsed)}
        theme={theme}
        onToggleTheme={toggleTheme}
      />

      {/* Main content area */}
      <main
        className="transition-all duration-200"
        style={{ marginLeft: sidebarCollapsed ? 60 : 240 }}
      >
        {/* Top bar */}
        <header className="sticky top-0 z-30 flex h-14 items-center gap-4 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 px-6">
          <div className="flex flex-1 items-center gap-3">
            <h1 className="text-sm font-semibold text-foreground">
              {pageConfig[activeView].title}
            </h1>
            <Separator orientation="vertical" className="h-4" />
            <p className="text-sm text-muted-foreground hidden sm:block">
              {pageConfig[activeView].description}
            </p>
          </div>
        </header>

        {/* Page content */}
        <div className="p-6">
          <div className="animate-in">
            {renderContent()}
          </div>
        </div>
      </main>
    </div>
  );
}

export default App
