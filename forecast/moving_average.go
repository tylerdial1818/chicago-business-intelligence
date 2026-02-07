package forecast

import (
	"math"
)

// ForecastResult represents a single forecast point
type ForecastResult struct {
	Period     string  `json:"period"`
	Predicted  float64 `json:"predicted"`
	LowerBound float64 `json:"lower_bound"`
	UpperBound float64 `json:"upper_bound"`
	Historical bool    `json:"historical"`
}

// MovingAverageForecast takes historical trip counts and forecasts forward
// Uses weighted moving average with trend component
func MovingAverageForecast(historical []float64, periodsAhead int, windowSize int) []ForecastResult {
	if len(historical) == 0 {
		return []ForecastResult{}
	}

	if windowSize > len(historical) {
		windowSize = len(historical)
	}

	if windowSize < 2 {
		windowSize = 2
	}

	results := make([]ForecastResult, 0, len(historical)+periodsAhead)

	// Add historical data as results
	for _, value := range historical {
		results = append(results, ForecastResult{
			Period:     "",
			Predicted:  value,
			LowerBound: value,
			UpperBound: value,
			Historical: true,
		})
	}

	// Calculate trend from the last windowSize points
	trend := calculateTrend(historical, windowSize)

	// Calculate standard deviation for confidence bounds
	stdDev := calculateStdDev(historical, windowSize)

	// Generate forecasts
	lastValue := historical[len(historical)-1]

	for i := 1; i <= periodsAhead; i++ {
		// Weighted moving average with trend
		predicted := lastValue + (trend * float64(i))

		// Widen confidence bounds as we go further into the future
		confidenceMultiplier := 1.0 + (float64(i) * 0.2)
		lowerBound := predicted - (stdDev * confidenceMultiplier)
		upperBound := predicted + (stdDev * confidenceMultiplier)

		// Ensure non-negative predictions (trip counts can't be negative)
		if lowerBound < 0 {
			lowerBound = 0
		}
		if predicted < 0 {
			predicted = 0
		}

		results = append(results, ForecastResult{
			Period:     "",
			Predicted:  predicted,
			LowerBound: lowerBound,
			UpperBound: upperBound,
			Historical: false,
		})
	}

	return results
}

// calculateTrend calculates the linear trend from the last n points
func calculateTrend(data []float64, windowSize int) float64 {
	if len(data) < windowSize {
		windowSize = len(data)
	}

	if windowSize < 2 {
		return 0
	}

	// Use last windowSize points
	window := data[len(data)-windowSize:]

	// Calculate simple linear regression slope
	n := float64(len(window))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, y := range window {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Slope = (n * sumXY - sumX * sumY) / (n * sumX2 - sumX * sumX)
	numerator := (n * sumXY) - (sumX * sumY)
	denominator := (n * sumX2) - (sumX * sumX)

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// calculateStdDev calculates the standard deviation of the last n points
func calculateStdDev(data []float64, windowSize int) float64 {
	if len(data) < windowSize {
		windowSize = len(data)
	}

	if windowSize < 2 {
		return 0
	}

	// Use last windowSize points
	window := data[len(data)-windowSize:]

	// Calculate mean
	mean := 0.0
	for _, value := range window {
		mean += value
	}
	mean /= float64(len(window))

	// Calculate variance
	variance := 0.0
	for _, value := range window {
		diff := value - mean
		variance += diff * diff
	}
	variance /= float64(len(window))

	// Standard deviation
	return math.Sqrt(variance)
}

// SimpleMovingAverage calculates a simple moving average
func SimpleMovingAverage(data []float64, windowSize int) float64 {
	if len(data) == 0 {
		return 0
	}

	if windowSize > len(data) {
		windowSize = len(data)
	}

	sum := 0.0
	for i := len(data) - windowSize; i < len(data); i++ {
		sum += data[i]
	}

	return sum / float64(windowSize)
}
