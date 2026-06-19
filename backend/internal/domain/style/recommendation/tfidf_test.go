package recommendation

import (
	"math"
	"testing"
)

func TestTokenize(t *testing.T) {
	got := Tokenize("Modern, Clean! UI design 2026.")
	want := []string{"modern", "clean", "ui", "design", "2026"}

	if len(got) != len(want) {
		t.Fatalf("Tokenize() len = %d, want %d; got %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Tokenize()[%d] = %q, want %q; got %#v", i, got[i], want[i], got)
		}
	}
}

func TestTFIDFCalculator_TF(t *testing.T) {
	doc := &Document{ID: "1", Tokens: []string{"modern", "clean", "modern", "ui"}}
	calc := NewTFIDFCalculator([]*Document{doc})

	if got := calc.TF("modern", doc); math.Abs(got-0.5) > 1e-9 {
		t.Fatalf("TF(modern) = %v, want 0.5", got)
	}
	if got := calc.TF("ui", doc); math.Abs(got-0.25) > 1e-9 {
		t.Fatalf("TF(ui) = %v, want 0.25", got)
	}
	if got := calc.TF("missing", doc); got != 0 {
		t.Fatalf("TF(missing) = %v, want 0", got)
	}
}

func TestTFIDFCalculator_IDF(t *testing.T) {
	docs := []*Document{
		{ID: "1", Tokens: []string{"modern", "clean"}},
		{ID: "2", Tokens: []string{"modern", "retro"}},
		{ID: "3", Tokens: []string{"nature", "forest"}},
	}
	calc := NewTFIDFCalculator(docs)

	if got := calc.IDF("modern"); math.Abs(got-math.Log(3.0/2.0)) > 1e-9 {
		t.Fatalf("IDF(modern) = %v, want %v", got, math.Log(3.0/2.0))
	}
	if got := calc.IDF("forest"); math.Abs(got-math.Log(3.0/1.0)) > 1e-9 {
		t.Fatalf("IDF(forest) = %v, want %v", got, math.Log(3.0/1.0))
	}
	if got := calc.IDF("missing"); got != 0 {
		t.Fatalf("IDF(missing) = %v, want 0", got)
	}
}

func TestTFIDFCalculator_Calculate(t *testing.T) {
	docs := []*Document{
		{ID: "1", Tokens: []string{"modern", "clean"}},
		{ID: "2", Tokens: []string{"modern", "retro"}},
	}
	calc := NewTFIDFCalculator(docs)
	doc := &Document{ID: "3", Tokens: []string{"modern", "modern", "ui"}}

	got := calc.Calculate(doc)

	wantModern := 2.0 / 3.0 * math.Log(2.0/2.0)
	wantUI := 1.0 / 3.0 * math.Log(2.0/1.0)

	if len(got) != 2 {
		t.Fatalf("Calculate() len = %d, want 2; got %#v", len(got), got)
	}
	if math.Abs(got["modern"]-wantModern) > 1e-9 {
		t.Fatalf("Calculate()[modern] = %v, want %v", got["modern"], wantModern)
	}
	if math.Abs(got["ui"]-wantUI) > 1e-9 {
		t.Fatalf("Calculate()[ui] = %v, want %v", got["ui"], wantUI)
	}
}
