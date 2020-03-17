package format

import "math"

// Duration returns formated value of GET request fetchURL method.
// Rounding precission is set by prec arg only if diff arg is less
// than GET request timeout.
func Duration(diff, prec float64) float64 {
	var result float64
	switch {
	case diff > 5.0:
		result = 5.0
	case diff < 5.0:
		result = math.Floor(diff*prec+0.5) / prec
	}
	return result
}
