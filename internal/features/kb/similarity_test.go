package kb

import "testing"

func TestCosineSimilarity_IdenticalVectors(t *testing.T) {
	a := []float32{1, 2, 3}
	b := []float32{1, 2, 3}

	got := CosineSimilarity(a, b)

	if got < 0.999 || got > 1.001 {
		t.Fatalf("expected ~1.0, got %v", got)
	}
}

func TestCosineSimilarity_OrthogonalVectors(t *testing.T) {
	a := []float32{1, 0}
	b := []float32{0, 1}

	got := CosineSimilarity(a, b)

	if got < -0.001 || got > 0.001 {
		t.Fatalf("expected ~0.0, got %v", got)
	}
}

func TestCosineSimilarity_OppositeVectors(t *testing.T) {
	a := []float32{1, 0}
	b := []float32{-1, 0}

	got := CosineSimilarity(a, b)

	if got < -1.001 || got > -0.999 {
		t.Fatalf("expected ~-1.0, got %v", got)
	}
}

func TestCosineSimilarity_EmptyVectors(t *testing.T) {
	if got := CosineSimilarity(nil, []float32{1, 2}); got != 0 {
		t.Fatalf("expected 0 for empty vector, got %v", got)
	}
}

func TestCosineSimilarity_MismatchedLength(t *testing.T) {
	a := []float32{1, 2, 3}
	b := []float32{1, 2}

	if got := CosineSimilarity(a, b); got != 0 {
		t.Fatalf("expected 0 for mismatched length, got %v", got)
	}
}

func TestCosineSimilarity_ZeroMagnitude(t *testing.T) {
	a := []float32{0, 0, 0}
	b := []float32{1, 2, 3}

	if got := CosineSimilarity(a, b); got != 0 {
		t.Fatalf("expected 0 for zero-magnitude vector, got %v", got)
	}
}
