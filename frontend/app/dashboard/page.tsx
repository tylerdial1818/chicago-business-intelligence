'use client'

import { useState, useEffect } from 'react'
import { BarChart, Bar, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import { AlertTriangle, TrendingUp, Building2, Map, DollarSign, Activity } from 'lucide-react'

// Mock data for demonstration - will be replaced with API calls
const mockCOVIDData = [
  { date: '2024-01-01', cases: 45, trips: 230, alert: 'LOW' },
  { date: '2024-01-08', cases: 67, trips: 210, alert: 'MEDIUM' },
  { date: '2024-01-15', cases: 89, trips: 195, alert: 'MEDIUM' },
  { date: '2024-01-22', cases: 120, trips: 180, alert: 'HIGH' },
  { date: '2024-01-29', cases: 98, trips: 185, alert: 'MEDIUM' },
]

const mockAirportData = [
  { zip: '60601', trips: 450, caseRate: 45.2 },
  { zip: '60602', trips: 380, caseRate: 67.8 },
  { zip: '60603', trips: 320, caseRate: 34.1 },
  { zip: '60604', trips: 290, caseRate: 89.3 },
  { zip: '60605', trips: 265, caseRate: 23.7 },
]

const mockInvestmentTargets = [
  { name: 'Englewood', unemployment: '18.5%', poverty: '42.3%', permits: 12 },
  { name: 'West Englewood', unemployment: '16.8%', poverty: '38.7%', permits: 8 },
  { name: 'Austin', unemployment: '15.2%', poverty: '35.4%', permits: 15 },
  { name: 'North Lawndale', unemployment: '14.9%', poverty: '34.1%', permits: 10 },
  { name: 'South Shore', unemployment: '13.7%', poverty: '31.8%', permits: 18 },
]

export default function DashboardPage() {
  const [selectedZip, setSelectedZip] = useState('60601')

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="mx-auto max-w-7xl px-6 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                Chicago Business Intelligence Dashboard
              </h1>
              <p className="mt-1 text-sm text-gray-600">
                Real-time analytics for strategic planning
              </p>
            </div>
            <div className="flex items-center gap-2 text-sm">
              <div className="h-2 w-2 rounded-full bg-green-500 animate-pulse"></div>
              <span className="text-gray-600">Live Data</span>
            </div>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-6 py-8">
        {/* Key Metrics */}
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4 mb-8">
          <MetricCard
            icon={<AlertTriangle className="h-6 w-6" />}
            title="Active COVID Alerts"
            value="3"
            change="+12% from last week"
            color="red"
          />
          <MetricCard
            icon={<Activity className="h-6 w-6" />}
            title="Taxi Trips Today"
            value="12,450"
            change="+8% from yesterday"
            color="blue"
          />
          <MetricCard
            icon={<Building2 className="h-6 w-6" />}
            title="Building Permits"
            value="87"
            change="This month"
            color="green"
          />
          <MetricCard
            icon={<Map className="h-6 w-6" />}
            title="Monitored Zip Codes"
            value="77"
            change="Citywide coverage"
            color="purple"
          />
        </div>

        {/* COVID Alert System */}
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2 mb-8">
          <ChartCard title="COVID-19 Alerts by Zip Code" icon={<AlertTriangle />}>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Select Zip Code
              </label>
              <select
                value={selectedZip}
                onChange={(e) => setSelectedZip(e.target.value)}
                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              >
                <option value="60601">60601 - Loop</option>
                <option value="60602">60602 - Near North Side</option>
                <option value="60603">60603 - Grant Park</option>
              </select>
            </div>
            <ResponsiveContainer width="100%" height={250}>
              <LineChart data={mockCOVIDData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{fontSize: 12}} />
                <YAxis />
                <Tooltip />
                <Legend />
                <Line type="monotone" dataKey="cases" stroke="#EF4444" name="Weekly Cases" />
                <Line type="monotone" dataKey="trips" stroke="#3B82F6" name="Taxi Trips" />
              </LineChart>
            </ResponsiveContainer>
            <div className="mt-4 flex gap-2">
              <Badge color="red">HIGH</Badge>
              <Badge color="yellow">MEDIUM</Badge>
              <Badge color="green">LOW</Badge>
            </div>
          </ChartCard>

          <ChartCard title="Airport Traffic by Destination" icon={<Map />}>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={mockAirportData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="zip" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="trips" fill="#3B82F6" name="Trip Count" />
                <Bar dataKey="caseRate" fill="#EF4444" name="Case Rate" />
              </BarChart>
            </ResponsiveContainer>
            <p className="mt-4 text-sm text-gray-600">
              Trips from O'Hare (60666) and Midway (60638) airports to city zip codes
            </p>
          </ChartCard>
        </div>

        {/* Investment Opportunities */}
        <ChartCard title="Top 5 Neighborhoods for Investment" icon={<Building2 />}>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Neighborhood
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Unemployment
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Poverty Rate
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Building Permits
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Incentive
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {mockInvestmentTargets.map((target, idx) => (
                  <tr key={idx} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {target.name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {target.unemployment}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {target.poverty}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {target.permits}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <Badge color="green">Permit Fee Waiver</Badge>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="mt-4 p-4 bg-blue-50 rounded-lg">
            <h4 className="text-sm font-semibold text-blue-900 mb-2">
              Investment Strategy
            </h4>
            <p className="text-sm text-blue-700">
              These neighborhoods qualify for infrastructure investment and building permit fee waivers 
              to encourage business development in high-need areas.
            </p>
          </div>
        </ChartCard>

        {/* Small Business Loan Program */}
        <div className="mt-8">
          <ChartCard title="Small Business Emergency Loan Eligibility" icon={<DollarSign />}>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className="text-sm font-semibold text-gray-900 mb-3">
                  Illinois Small Business Emergency Loan Fund Delta
                </h4>
                <p className="text-sm text-gray-600 mb-4">
                  Low-interest loans up to $250,000 for small businesses in qualifying zip codes.
                </p>
                <div className="space-y-3">
                  <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <span className="text-sm font-medium">Permit Type</span>
                    <Badge color="blue">NEW CONSTRUCTION</Badge>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <span className="text-sm font-medium">Income Threshold</span>
                    <span className="text-sm text-gray-700">&lt; $30,000</span>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <span className="text-sm font-medium">Max Loan Amount</span>
                    <span className="text-sm text-gray-700">$250,000</span>
                  </div>
                </div>
              </div>
              <div className="bg-gradient-to-br from-green-50 to-emerald-50 p-6 rounded-lg">
                <h4 className="text-sm font-semibold text-green-900 mb-2">
                  Program Impact
                </h4>
                <p className="text-sm text-green-700 mb-4">
                  Helping small businesses compete with larger players like Amazon and Walmart 
                  for warehouse and commercial spaces.
                </p>
                <div className="mt-4 pt-4 border-t border-green-200">
                  <div className="text-3xl font-bold text-green-900">23</div>
                  <div className="text-sm text-green-600">Eligible Zip Codes</div>
                </div>
              </div>
            </div>
          </ChartCard>
        </div>
      </main>
    </div>
  )
}

interface MetricCardProps {
  icon: React.ReactNode
  title: string
  value: string
  change: string
  color: 'red' | 'blue' | 'green' | 'purple'
}

function MetricCard({ icon, title, value, change, color }: MetricCardProps) {
  const colorClasses = {
    red: 'bg-red-50 text-red-600',
    blue: 'bg-blue-50 text-blue-600',
    green: 'bg-green-50 text-green-600',
    purple: 'bg-purple-50 text-purple-600',
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between">
        <div className={`rounded-lg p-3 ${colorClasses[color]}`}>
          {icon}
        </div>
      </div>
      <div className="mt-4">
        <p className="text-sm text-gray-600">{title}</p>
        <p className="text-3xl font-bold text-gray-900 mt-2">{value}</p>
        <p className="text-sm text-gray-500 mt-1">{change}</p>
      </div>
    </div>
  )
}

interface ChartCardProps {
  title: string
  icon: React.ReactNode
  children: React.ReactNode
}

function ChartCard({ title, icon, children }: ChartCardProps) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center gap-3 mb-6">
        <div className="text-blue-600">{icon}</div>
        <h2 className="text-xl font-bold text-gray-900">{title}</h2>
      </div>
      {children}
    </div>
  )
}

function Badge({ color, children }: { color: 'red' | 'yellow' | 'green' | 'blue', children: React.ReactNode }) {
  const colorClasses = {
    red: 'bg-red-100 text-red-800',
    yellow: 'bg-yellow-100 text-yellow-800',
    green: 'bg-green-100 text-green-800',
    blue: 'bg-blue-100 text-blue-800',
  }

  return (
    <span className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ${colorClasses[color]}`}>
      {children}
    </span>
  )
}
