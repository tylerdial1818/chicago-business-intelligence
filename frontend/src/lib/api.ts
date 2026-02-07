const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function fetchAPI<T>(endpoint: string): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`);

  if (!response.ok) {
    throw new Error(`API error: ${response.statusText}`);
  }

  return response.json();
}

export const api = {
  // Health check
  health: () => fetchAPI<{
    status: string;
    database: string;
    version: string;
  }>('/health'),

  // Pipeline status
  pipelineStatus: () => fetchAPI<{
    runs: Array<{
      id: number;
      dataset_name: string;
      records_fetched: number;
      records_loaded: number;
      records_rejected: number;
      started_at: string;
      completed_at: string;
      status: string;
      error_message?: string;
    }>;
    total_runs: number;
    last_update: string;
  }>('/api/pipeline-status'),

  // Zip codes
  zipCodes: () => fetchAPI<Array<{
    zip_code: string;
    neighborhood: string;
    community_area: string;
  }>>('/api/zip-codes'),

  // COVID alerts
  covidAlerts: (zipCode: string) => fetchAPI<Array<{
    zip_code: string;
    week_start: string;
    week_end: string;
    cases_weekly: number;
    case_rate_weekly: number;
    tests_weekly: number;
    percent_tested_positive: number;
    taxi_trips: number;
    alert_level: string;
  }>>(`/api/covid-alerts?zip=${zipCode}`),

  // Airport traffic
  airportTraffic: () => fetchAPI<Array<{
    airport: string;
    destination_zip: string;
    destination_neighborhood: string;
    trip_count: number;
    avg_miles: number;
    avg_fare: number;
  }>>('/api/airport-traffic'),

  // High CCVI
  highCCVI: () => fetchAPI<Array<{
    community_area_name: string;
    ccvi_score: number;
    ccvi_category: string;
    zip_code: string;
    trips_from: number;
    trips_to: number;
    total_trips: number;
  }>>('/api/high-ccvi'),

  // Traffic patterns
  trafficPatterns: (zipCode: string) => fetchAPI<Array<{
    date: string;
    trip_count: number;
    avg_miles: number;
    avg_fare: number;
  }>>(`/api/traffic-patterns?zip=${zipCode}`),

  // Forecast
  forecast: (zipCode: string, period: 'd' | 'w' | 'm' = 'd') => fetchAPI<{
    zip_code: string;
    period: string;
    historical: Array<{
      period: string;
      predicted: number;
      lower_bound: number;
      upper_bound: number;
      historical: boolean;
    }>;
    forecast: Array<{
      period: string;
      predicted: number;
      lower_bound: number;
      upper_bound: number;
      historical: boolean;
    }>;
  }>(`/api/forecast?zip=${zipCode}&period=${period}`),

  // Investment targets
  investmentTargets: () => fetchAPI<Array<{
    community_area: number;
    community_area_name: string;
    unemployment_rate: number;
    poverty_rate: number;
    per_capita_income: number;
    hardship_index: number;
    total_permits: number;
    new_construction_permits: number;
  }>>('/api/investment-targets'),

  // Small business loans
  smallBusinessLoans: () => fetchAPI<Array<{
    zip_code: string;
    community_area_name: string;
    per_capita_income: number;
    unemployment_rate: number;
    poverty_rate: number;
    new_construction_permits: number;
    max_loan_amount: number;
  }>>('/api/small-business-loans'),

  // Building permits
  buildingPermits: (zipCode: string) => fetchAPI<Array<{
    id: string;
    permit_number: string;
    permit_type: string;
    review_type: string;
    application_start_date: string;
    issue_date: string;
    street_address: string;
    work_description: string;
    subtotal_paid: number;
    zip_code: string;
  }>>(`/api/building-permits?zip=${zipCode}`),
};
