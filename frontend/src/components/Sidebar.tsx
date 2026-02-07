'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  LayoutDashboard,
  AlertTriangle,
  Plane,
  Shield,
  TrendingUp,
  Building2,
  DollarSign,
  Activity
} from 'lucide-react';

const navigation = [
  { name: 'Overview', href: '/', icon: LayoutDashboard },
  { name: 'COVID Alerts', href: '/covid-alerts', icon: AlertTriangle },
  { name: 'Airport Traffic', href: '/airport-traffic', icon: Plane },
  { name: 'CCVI Analysis', href: '/ccvi', icon: Shield },
  { name: 'Traffic Forecast', href: '/traffic-forecast', icon: TrendingUp },
  { name: 'Investments', href: '/investments', icon: Building2 },
  { name: 'Business Loans', href: '/small-business-loans', icon: DollarSign },
  { name: 'Pipeline Status', href: '/pipeline', icon: Activity },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <div className="flex flex-col w-64 bg-slate-800 border-r border-slate-700">
      <div className="flex items-center h-16 px-6 border-b border-slate-700">
        <h1 className="text-xl font-bold text-slate-50">Chicago BI</h1>
      </div>

      <nav className="flex-1 px-4 py-6 space-y-1">
        {navigation.map((item) => {
          const isActive = pathname === item.href;
          const Icon = item.icon;

          return (
            <Link
              key={item.name}
              href={item.href}
              className={`
                flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-colors
                ${isActive
                  ? 'bg-blue-500 text-white'
                  : 'text-slate-400 hover:text-slate-50 hover:bg-slate-700'
                }
              `}
            >
              <Icon className="w-5 h-5 mr-3" />
              {item.name}
            </Link>
          );
        })}
      </nav>

      <div className="p-4 border-t border-slate-700">
        <p className="text-xs text-slate-500 text-center">
          Chicago Business Intelligence v1.0
        </p>
      </div>
    </div>
  );
}
