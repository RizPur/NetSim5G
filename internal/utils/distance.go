package utils

import "math"

func CalculateDistance(X1 float64, Y1 float64, X2 float64, Y2 float64) float64 {
	return math.Sqrt((X1-X2)*(X1-X2) + (Y1-Y2)*(Y1-Y2))
}
