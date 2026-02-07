'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { api } from '@/src/lib/api';
import { StatCard } from '@/src/components/StatCard';
import { LoadingState } from '@/src/components/LoadingState';
import {
  Database,
  AlertTriangle,
  Plane,
  Shield,
  TrendingUp,
  Building2
} from 'lucide-react';

export default function HomePage() {
  const [health, setHealth] = useState<any>(null);
  const [pipelineStatus, setPipelineStatus] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.health(),
      api.pipelineStatus(),
    ])
      .then(([healthData, pipelineData]) => {
        setHealth(healthData);
        setPipelineStatus(pipelineData);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <LoadingState />;

  const successfulRuns = pipelineStatus?.runs?.filter((r: any) => r.status === 'SUCCESS').length || 0;

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-slate-50">Chicago Business Intelligence</h1>
        <p className="text-slate-400 mt-2">
          Data engineering platform for strategic planning and analytics
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          label="Database Status"
          value={health?.database === 'connected' ? 'Connected' : 'Disconnected'}
          icon={<Database className="w-8 h-8" />}
        />
        <StatCard
          label="Pipeline Runs"
          value={pipelineStatus?.total_runs || 0}
          icon={<TrendingUp className="w-8 h-8" />}
        />
        <StatCard
          label="Successful Runs"
          value={successfulRuns}
          icon={<Database className="w-8 h-8" />}
        />
        <StatCard
          label="Datasets Loaded"
          value="7"
          icon={<Building2 className="w-8 h-8" />}
        />
      </div>

      {/* Feature Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <FeatureCard
          title="COVID-19 Alerts"
          description="Track COVID metrics and taxi trip correlations by zip code with alert levels"
          icon={<AlertTriangle className="w-6 h-6" />}
          href="/covid-alerts"
        />
        <FeatureCard
          title="Airport Traffic"
          description="Monitor trips from O'Hare and Midway to neighborhoods across Chicago"
          icon={<Plane className="w-6 h-6" />}
          href="/airport-traffic"
        />
        <FeatureCard
          title="CCVI Analysis"
          description="Analyze mobility patterns in high vulnerability communities"
          icon={<Shield className="w-6 h-6" />}
          href="/ccvi"
        />
        <FeatureCard
          title="Traffic Forecast"
          description="Forecast daily, weekly, and monthly trip volumes by zip code"
          icon={<TrendingUp className="w-6 h-6" />}
          href="/traffic-forecast"
        />
        <FeatureCard
          title="Investment Targets"
          description="Identify top neighborhoods for development investment"
          icon={<Building2 className="w-6 h-6" />}
          href="/investments"
        />
        <FeatureCard
          title="Small Business Loans"
          description="Find eligible zip codes for emergency loan programs"
          icon={<Building2 className="w-6 h-6" />}
          href="/small-business-loans"
        />
      </div>
    </div>
  );
}

function FeatureCard({ title, description, icon, href }: {
  title: string;
  description: string;
  icon: React.ReactNode;
  href: string;
}) {
  return (
    <Link href={href}>
      <div className="bg-slate-800 border border-slate-700 rounded-lg p-6 hover:border-blue-500 transition-colors cursor-pointer">
        <div className="flex items-center mb-3">
          <div className="text-blue-500 mr-3">
            {icon}
          </div>
          <h3 className="text-lg font-semibold text-slate-50">{title}</h3>
        </div>
        <p className="text-sm text-slate-400">{description}</p>
      </div>
    </Link>
  );
}
