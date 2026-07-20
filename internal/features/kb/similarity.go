package kb

import "math"

// CosineSimilarity returns the cosine similarity between two vectors,
// a value between -1 and 1 where 1 means identical direction. Returns
// 0 if either vector is empty, of mismatched length, or has zero
// magnitude (e.g. a chunk that was never embedded).
func CosineSimilarity(a, b []float32) float64 {

	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	if len(a) != len(b) {
		return 0
	}

	var dotProduct float64
	var magnitudeA float64
	var magnitudeB float64

	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		magnitudeA += float64(a[i]) * float64(a[i])
		magnitudeB += float64(b[i]) * float64(b[i])
	}

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
}
