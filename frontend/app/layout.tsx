import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import Link from 'next/link'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Chicago Business Intelligence - Strategic Planning & Analytics',
  description: 'Real-time data analytics combining transportation patterns, COVID-19 metrics, building permits, and socioeconomic indicators for Chicago strategic planning.',
  keywords: ['Chicago', 'Business Intelligence', 'Analytics', 'COVID-19', 'Urban Planning', 'Data Science'],
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <nav className="bg-white border-b border-gray-200">
          <div className="mx-auto max-w-7xl px-6 py-4">
            <div className="flex items-center justify-between">
              <Link href="/" className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-lg bg-gradient-to-br from-blue-600 to-blue-800 flex items-center justify-center">
                  <span className="text-white font-bold text-xl">CB</span>
                </div>
                <span className="text-xl font-bold text-gray-900">
                  Chicago Business Intelligence
                </span>
              </Link>
              <div className="flex gap-6">
                <Link 
                  href="/" 
                  className="text-gray-600 hover:text-gray-900 font-medium transition"
                >
                  Home
                </Link>
                <Link 
                  href="/dashboard" 
                  className="text-gray-600 hover:text-gray-900 font-medium transition"
                >
                  Dashboard
                </Link>
                <a 
                  href="https://github.com/tylerdial1818/chicago-business-intelligence"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-gray-600 hover:text-gray-900 font-medium transition"
                >
                  GitHub
                </a>
              </div>
            </div>
          </div>
        </nav>
        {children}
        <footer className="bg-gray-900 text-white py-12 mt-16">
          <div className="mx-auto max-w-7xl px-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div>
                <h3 className="text-lg font-semibold mb-4">Chicago Business Intelligence</h3>
                <p className="text-gray-400 text-sm">
                  Strategic planning and analytics for urban development using real-time data from the City of Chicago.
                </p>
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-4">Data Sources</h3>
                <ul className="space-y-2 text-sm text-gray-400">
                  <li>City of Chicago Data Portal</li>
                  <li>Transportation Network Providers</li>
                  <li>Building Permits Database</li>
                  <li>COVID-19 Health Metrics</li>
                </ul>
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-4">Built By</h3>
                <p className="text-gray-400 text-sm">
                  Tyler Dial<br />
                  Dialed Intelligence LLC<br />
                  <a href="https://github.com/tylerdial1818" className="text-blue-400 hover:text-blue-300 transition">
                    GitHub Profile →
                  </a>
                </p>
              </div>
            </div>
            <div className="mt-8 pt-8 border-t border-gray-800 text-center text-sm text-gray-400">
              <p>© 2026 Dialed Intelligence LLC. Built for strategic planning and analytics.</p>
            </div>
          </div>
        </footer>
      </body>
    </html>
  )
}
