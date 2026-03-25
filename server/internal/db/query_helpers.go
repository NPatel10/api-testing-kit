package db

import "math"

func normalizeLimit(limit int32, fallback int32, max int32) int32 {
	if fallback <= 0 {
		fallback = 1
	}
	if max < fallback {
		max = fallback
	}
	if limit <= 0 {
		return fallback
	}
	return int32(math.Min(float64(limit), float64(max)))
}
