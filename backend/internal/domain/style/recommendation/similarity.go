package recommendation

import "math"

// CosineSimilarity calculates cosine similarity between two TF-IDF vectors.
func CosineSimilarity(vec1, vec2 map[string]float64) float64 {
	mag1 := VectorMagnitude(vec1)
	mag2 := VectorMagnitude(vec2)
	if mag1 == 0 || mag2 == 0 {
		return 0
	}

	score := DotProduct(vec1, vec2) / (mag1 * mag2)
	if score < 0 {
		return 0
	}
	if score > 1 {
		return 1
	}
	return score
}

// VectorMagnitude calculates the magnitude of a vector.
func VectorMagnitude(vec map[string]float64) float64 {
	sum := 0.0
	for _, value := range vec {
		sum += value * value
	}
	return math.Sqrt(sum)
}

// DotProduct calculates dot product of two vectors.
func DotProduct(vec1, vec2 map[string]float64) float64 {
	sum := 0.0
	if len(vec1) > len(vec2) {
		vec1, vec2 = vec2, vec1
	}
	for term, value1 := range vec1 {
		if value2, ok := vec2[term]; ok {
			sum += value1 * value2
		}
	}
	return sum
}
