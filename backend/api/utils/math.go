package utils

import "math/rand/v2"

func Float64Between(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
