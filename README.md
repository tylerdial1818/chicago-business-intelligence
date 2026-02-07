# Chicago Business Intelligence

> Strategic planning and analytics platform for urban development using real-time data from the City of Chicago.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Next.js](https://img.shields.io/badge/Next.js-14-black?style=flat&logo=next.js)](https://nextjs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-blue?style=flat&logo=typescript)](https://www.typescriptlang.org)

## 🎯 Overview

Chicago Business Intelligence is a comprehensive analytics platform that combines transportation patterns, COVID-19 metrics, building permits, and socioeconomic indicators to drive data-informed decision-making for Chicago's neighborhoods and strategic planning initiatives.

The platform provides six core intelligence modules designed for city planners, business analysts, and policy makers:

1. **COVID-19 Alert System** - Track COVID metrics by zip code correlated with taxi trip patterns
2. **Airport Traffic Analysis** - Monitor trips from O'Hare and Midway to neighborhoods
3. **Community Vulnerability Index** - Identify high CCVI neighborhoods for targeted support
4. **Traffic Forecasting** - Predict taxi trip volumes for infrastructure planning
5. **Investment Opportunities** - Target neighborhoods for strategic investment
6. **Small Business Loan Eligibility** - Identify zip codes qualifying for emergency loans

## 🏗️ Architecture

### Tech Stack

**Backend**
- **Language**: Go 1.21+
- **Database**: PostgreSQL 14+
- **APIs**: RESTful JSON endpoints
- **Deployment**: Docker + Google Cloud Run

**Frontend**
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS
- **Charts**: Recharts
- **Deployment**: Vercel

**Data Collection**
- **Source**: City of Chicago Data Portal (SODA API)
- **Datasets**: Taxi Trips, Building Permits, COVID-19, Unemployment, CCVI

### System Design

```
┌─────────────────┐
│   City of       │
│   Chicago       │
│   Data Portal   │
└────────┬────────┘
         │
         │ SODA API
         │
         ▼
┌─────────────────┐       ┌──────────────────┐
│   Go Backend    │◄─────►│   PostgreSQL     │
│   (Data Lake)   │       │   (Data Storage) │
└────────┬────────┘       └──────────────────┘
         │
         │ REST API
         │
         ▼
┌─────────────────┐
│   Next.js       │
│   Dashboard     │
│   (Frontend)    │
└─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+ (or use Docker)

### Local Development

1. **Clone the repository**

```bash
git clone https://github.com/tylerdial1818/chicago-business-intelligence.git
cd chicago-business-intelligence
```

2. **Start PostgreSQL**

```bash
docker run --name chicago-db -e POSTGRES_PASSWORD=root -e POSTGRES_DB=chicago_business_intelligence -p 5432:5432 -d postgres:14
```

3. **Run the backend**

```bash
# Install dependencies
go mod download

# Run the server
go run main.go

# Server will start on http://localhost:8080
```

4. **Run the frontend**

```bash
cd frontend
npm install
npm run dev

# Frontend will start on http://localhost:3000
```

5. **Access the application**

- Frontend: http://localhost:3000
- API Health Check: http://localhost:8080/health

## 📊 API Endpoints

### Health & Status

```
GET /health
```

Returns API health status.

### COVID-19 Analytics

```
GET /api/covid-alerts?zip={zipCode}
```

Returns COVID-19 weekly metrics and alert levels correlated with taxi trip patterns.

**Response:**
```json
{
  "zip_code": "60601",
  "data": [
    {
      "week_start": "2024-01-29",
      "cases_weekly": "45",
      "case_rate": "120.5",
      "taxi_trips": 234,
      "alert_level": "MEDIUM"
    }
  ]
}
```

### Airport Traffic

```
GET /api/airport-traffic
```

Returns taxi trip patterns from O'Hare (60666) and Midway (60638) airports.

### Community Vulnerability

```
GET /api/high-ccvi
```

Returns neighborhoods with HIGH CCVI category and associated mobility patterns.

### Investment Targets

```
GET /api/investment-targets
```

Returns top 5 neighborhoods ranked by unemployment and poverty rates.

### Small Business Loans

```
GET /api/small-business-loans
```

Returns zip codes eligible for emergency loan programs.

### Traffic Patterns

```
GET /api/traffic-patterns?zip={zipCode}
```

Returns historical taxi trip volumes for the specified zip code.

### Zip Codes

```
GET /api/zip-codes
```

Returns list of all zip codes with available data.

## 📈 Data Sources

All data is sourced from the [City of Chicago Data Portal](https://data.cityofchicago.org/):

- **Transportation**: [Taxi Trips](https://data.cityofchicago.org/Transportation/Taxi-Trips/wrvz-psew) & [TNP Trips](https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips/m6dm-c72p)
- **Buildings**: [Building Permits](https://data.cityofchicago.org/Buildings/Building-Permits/ydr8-5enu)
- **Health**: [COVID-19 Cases](https://data.cityofchicago.org/Health-Human-Services/COVID-19-Cases-Tests-and-Deaths-by-ZIP-Code/yhhz-zm2v) & [CCVI](https://data.cityofchicago.org/Health-Human-Services/COVID-19-Community-Vulnerability-Index-CCVI/xhc6-88s9)
- **Socioeconomic**: [Public Health Statistics](https://data.cityofchicago.org/Health-Human-Services/Public-Health-Statistics-Selected-public-health-in/iqnk-2tcu)

## 🎨 Features

### Real-Time Analytics

- Live COVID-19 alert system with three-tier classification (LOW/MEDIUM/HIGH)
- Taxi trip pattern analysis and correlation with health metrics
- Airport traffic monitoring to 77+ Chicago zip codes

### Strategic Planning

- Investment opportunity identification based on unemployment and poverty rates
- Building permit analysis for infrastructure planning
- Small business loan eligibility targeting

### Interactive Dashboard

- Dynamic data visualizations with charts and graphs
- Zip code filtering and selection
- Responsive design for desktop and mobile
- Real-time data updates

## 🔐 Environment Variables

Create a `.env` file in the backend directory:

```bash
DATABASE_URL=postgresql://postgres:root@localhost:5432/chicago_business_intelligence
PORT=8080
GOOGLE_GEOCODING_API_KEY=your_api_key_here
```

## 🐳 Docker Deployment

```bash
# Build the backend container
docker build -t chicago-intelligence-backend .

# Run with Docker Compose
docker-compose up -d
```

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## 👤 Author

**Tyler Dial**  
Dialed Intelligence LLC  
[GitHub](https://github.com/tylerdial1818) | [Website](https://dialedintelligence.com)

## 🙏 Acknowledgments

- City of Chicago for providing open data access
- Northwestern University MSDS Program
- Built with modern web technologies and best practices

---

**Note**: This project demonstrates real-world application of data engineering, cloud-native architecture, and full-stack development skills for portfolio purposes.
