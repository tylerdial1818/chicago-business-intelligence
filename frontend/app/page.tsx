import Link from 'next/link'
import { ArrowRight, BarChart3, Map, TrendingUp, Building2, AlertTriangle, DollarSign } from 'lucide-react'

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
      {/* Hero Section */}
      <div className="relative overflow-hidden bg-gradient-to-r from-blue-600 to-blue-800 text-white">
        <div className="absolute inset-0 bg-grid-white/10"></div>
        <div className="relative mx-auto max-w-7xl px-6 py-24 sm:py-32">
          <div className="text-center">
            <h1 className="text-5xl font-bold tracking-tight sm:text-6xl">
              Chicago Business Intelligence
            </h1>
            <p className="mt-6 text-xl leading-8 text-blue-100">
              Strategic Planning & Analytics for Urban Development
            </p>
            <p className="mx-auto mt-4 max-w-2xl text-lg text-blue-200">
              Real-time data analytics combining transportation patterns, COVID-19 metrics, 
              building permits, and socioeconomic indicators to drive informed decision-making 
              for Chicago's neighborhoods.
            </p>
            <div className="mt-10 flex items-center justify-center gap-x-6">
              <Link
                href="/dashboard"
                className="rounded-md bg-white px-6 py-3 text-base font-semibold text-blue-600 shadow-sm hover:bg-blue-50 transition flex items-center gap-2"
              >
                View Dashboard
                <ArrowRight className="h-5 w-5" />
              </Link>
              <a
                href="https://github.com/tylerdial1818/chicago-business-intelligence"
                target="_blank"
                rel="noopener noreferrer"
                className="text-base font-semibold leading-7 text-white hover:text-blue-100 transition"
              >
                View on GitHub <span aria-hidden="true">→</span>
              </a>
            </div>
          </div>
        </div>
      </div>

      {/* Features Grid */}
      <div className="mx-auto max-w-7xl px-6 py-24">
        <div className="text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
            Comprehensive Urban Analytics
          </h2>
          <p className="mt-4 text-lg text-gray-600">
            Six core intelligence modules for data-driven strategic planning
          </p>
        </div>

        <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
          {/* Feature 1: COVID Alerts */}
          <FeatureCard
            icon={<AlertTriangle className="h-8 w-8" />}
            title="COVID-19 Alert System"
            description="Track COVID-19 metrics by zip code correlated with taxi trip patterns to identify high-risk areas and prevent super-spreader events."
            href="/dashboard"
            color="red"
          />

          {/* Feature 2: Airport Traffic */}
          <FeatureCard
            icon={<Map className="h-8 w-8" />}
            title="Airport Traffic Analysis"
            description="Monitor taxi trips from O'Hare and Midway airports to neighborhoods, tracking COVID exposure and traffic patterns."
            href="/dashboard"
            color="blue"
          />

          {/* Feature 3: Vulnerability Index */}
          <FeatureCard
            icon={<TrendingUp className="h-8 w-8" />}
            title="Community Vulnerability"
            description="Identify high CCVI (COVID Community Vulnerability Index) neighborhoods and track mobility patterns for targeted support."
            href="/dashboard"
            color="orange"
          />

          {/* Feature 4: Traffic Forecasting */}
          <FeatureCard
            icon={<BarChart3 className="h-8 w-8" />}
            title="Traffic Forecasting"
            description="Forecast daily, weekly, and monthly taxi trip volumes by zip code for streetscaping and infrastructure planning."
            href="/dashboard"
            color="green"
          />

          {/* Feature 5: Investment Targeting */}
          <FeatureCard
            icon={<Building2 className="h-8 w-8" />}
            title="Investment Opportunities"
            description="Identify top 5 neighborhoods by unemployment and poverty for strategic infrastructure investment and permit fee waivers."
            href="/dashboard"
            color="purple"
          />

          {/* Feature 6: Small Business Loans */}
          <FeatureCard
            icon={<DollarSign className="h-8 w-8" />}
            title="Small Business Loans"
            description="Target low-construction zip codes with low per capita income for small business emergency loan programs."
            href="/dashboard"
            color="indigo"
          />
        </div>
      </div>

      {/* Tech Stack Section */}
      <div className="bg-white py-16">
        <div className="mx-auto max-w-7xl px-6">
          <h3 className="text-center text-2xl font-bold text-gray-900 mb-8">
            Built with Modern Technology
          </h3>
          <div className="grid grid-cols-2 gap-8 md:grid-cols-4 lg:grid-cols-6 text-center">
            <TechBadge name="Next.js" />
            <TechBadge name="TypeScript" />
            <TechBadge name="Go" />
            <TechBadge name="PostgreSQL" />
            <TechBadge name="Docker" />
            <TechBadge name="Google Cloud" />
          </div>
        </div>
      </div>

      {/* Data Sources */}
      <div className="mx-auto max-w-7xl px-6 py-16">
        <h3 className="text-center text-2xl font-bold text-gray-900 mb-8">
          Data Sources
        </h3>
        <p className="text-center text-gray-600 max-w-3xl mx-auto mb-8">
          All data sourced from the City of Chicago Data Portal, including taxi trips, 
          building permits, COVID-19 statistics, unemployment metrics, and community health indicators.
        </p>
        <div className="text-center">
          <a
            href="https://data.cityofchicago.org/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:text-blue-800 font-semibold"
          >
            Explore Chicago Data Portal →
          </a>
        </div>
      </div>
    </div>
  )
}

interface FeatureCardProps {
  icon: React.ReactNode
  title: string
  description: string
  href: string
  color: 'red' | 'blue' | 'orange' | 'green' | 'purple' | 'indigo'
}

function FeatureCard({ icon, title, description, href, color }: FeatureCardProps) {
  const colorClasses = {
    red: 'bg-red-50 text-red-600',
    blue: 'bg-blue-50 text-blue-600',
    orange: 'bg-orange-50 text-orange-600',
    green: 'bg-green-50 text-green-600',
    purple: 'bg-purple-50 text-purple-600',
    indigo: 'bg-indigo-50 text-indigo-600',
  }

  return (
    <Link href={href} className="group">
      <div className="relative overflow-hidden rounded-lg border border-gray-200 bg-white p-6 hover:shadow-lg transition-all duration-200">
        <div className={`inline-flex rounded-lg p-3 ${colorClasses[color]}`}>
          {icon}
        </div>
        <h3 className="mt-4 text-xl font-semibold text-gray-900 group-hover:text-blue-600 transition">
          {title}
        </h3>
        <p className="mt-2 text-gray-600 text-sm leading-relaxed">
          {description}
        </p>
      </div>
    </Link>
  )
}

function TechBadge({ name }: { name: string }) {
  return (
    <div className="flex items-center justify-center">
      <span className="inline-flex items-center rounded-full bg-blue-50 px-4 py-2 text-sm font-medium text-blue-700">
        {name}
      </span>
    </div>
  )
}
