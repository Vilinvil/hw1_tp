package calc

import "math"

const tolerance = 1e-8

func equalFloat(num1, num2 float64) bool {
	return math.Abs(num1-num2) < tolerance
}
