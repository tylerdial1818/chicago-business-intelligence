export interface HealthStatus {
  status: string;
  database: string;
  version: string;
}

export interface PipelineRun {
  id: number;
  dataset_name: string;
  records_fetched: number;
  records_loaded: number;
  records_rejected: number;
  started_at: string;
  completed_at: string;
  status: string;
  error_message?: string;
}

export interface PipelineStatus {
  runs: PipelineRun[];
  total_runs: number;
  last_update: string;
}

export interface ZipCode {
  zip_code: string;
  neighborhood: string;
  community_area: string;
}

export interface CovidAlert {
  zip_code: string;
  week_start: string;
  week_end: string;
  cases_weekly: number;
  case_rate_weekly: number;
  tests_weekly: number;
  percent_tested_positive: number;
  taxi_trips: number;
  alert_level: 'LOW' | 'MEDIUM' | 'HIGH';
}

export interface AirportTraffic {
  airport: string;
  destination_zip: string;
  destination_neighborhood: string;
  trip_count: number;
  avg_miles: number;
  avg_fare: number;
}

export interface CCVITrip {
  community_area_name: string;
  ccvi_score: number;
  ccvi_category: string;
  zip_code: string;
  trips_from: number;
  trips_to: number;
  total_trips: number;
}

export interface TrafficPattern {
  date: string;
  trip_count: number;
  avg_miles: number;
  avg_fare: number;
}

export interface ForecastPoint {
  period: string;
  predicted: number;
  lower_bound: number;
  upper_bound: number;
  historical: boolean;
}

export interface ForecastResponse {
  zip_code: string;
  period: string;
  historical: ForecastPoint[];
  forecast: ForecastPoint[];
}

export interface InvestmentTarget {
  community_area: number;
  community_area_name: string;
  unemployment_rate: number;
  poverty_rate: number;
  per_capita_income: number;
  hardship_index: number;
  total_permits: number;
  new_construction_permits: number;
}

export interface SmallBusinessLoan {
  zip_code: string;
  community_area_name: string;
  per_capita_income: number;
  unemployment_rate: number;
  poverty_rate: number;
  new_construction_permits: number;
  max_loan_amount: number;
}

export interface BuildingPermit {
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
}
