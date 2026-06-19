package recommendation

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	vec1 := map[string]float64{"modern": 0.8, "clean": 0.6}
	vec2 := map[string]float64{"modern": 0.4, "clean": 0.3}

	if got := CosineSimilarity(vec1, vec2); math.Abs(got-1.0) > 1e-9 {
		t.Fatalf("CosineSimilarity() = %v, want 1", got)
	}

	vec3 := map[string]float64{"retro": 1}
	if got := CosineSimilarity(vec1, vec3); got != 0 {
		t.Fatalf("CosineSimilarity() = %v, want 0", got)
	}
}

func TestVectorMagnitude(t *testing.T) {
	vec := map[string]float64{"a": 3, "b": 4}

	if got := VectorMagnitude(vec); math.Abs(got-5.0) > 1e-9 {
		t.Fatalf("VectorMagnitude() = %v, want 5", got)
	}
}

func TestDotProduct(t *testing.T) {
	vec1 := map[string]float64{"a": 2, "b": 3}
	vec2 := map[string]float64{"a": 4, "c": 5}

	if got := DotProduct(vec1, vec2); math.Abs(got-8.0) > 1e-9 {
		t.Fatalf("DotProduct() = %v, want 8", got)
	}
}
