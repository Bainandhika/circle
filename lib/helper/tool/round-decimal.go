package tool

import "math"

func Round(value float64, precision int) float64 {
	// Calculate the factor to multiply with based on the precision.
	factor := math.Pow(10, float64(precision))
	// Multiply the value by the factor, round it, then divide by the factor.
	return math.Round(value*factor) / factor
}
