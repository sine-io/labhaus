package recommendation

import (
	"math"
	"strings"
	"unicode"
)

// Document represents a text document for TF-IDF calculation.
type Document struct {
	ID      string
	Content string
	Tokens  []string
}

// TFIDFCalculator calculates TF-IDF scores for documents.
type TFIDFCalculator struct {
	documents      []*Document
	idf            map[string]float64
	vocabularySize int
}

// NewTFIDFCalculator creates a new TF-IDF calculator.
func NewTFIDFCalculator(documents []*Document) *TFIDFCalculator {
	c := &TFIDFCalculator{
		documents: documents,
		idf:       make(map[string]float64),
	}

	documentFrequency := make(map[string]int)
	for _, doc := range documents {
		ensureTokens(doc)
		seen := make(map[string]bool)
		for _, token := range doc.Tokens {
			if token == "" || seen[token] {
				continue
			}
			seen[token] = true
			documentFrequency[token]++
		}
	}

	c.vocabularySize = len(documentFrequency)
	totalDocuments := float64(len(documents))
	for term, frequency := range documentFrequency {
		c.idf[term] = math.Log(totalDocuments / float64(frequency))
	}

	return c
}

// Calculate computes TF-IDF scores for a document.
func (c *TFIDFCalculator) Calculate(doc *Document) map[string]float64 {
	ensureTokens(doc)

	vector := make(map[string]float64)
	seen := make(map[string]bool)
	for _, token := range doc.Tokens {
		if token == "" || seen[token] {
			continue
		}
		seen[token] = true
		vector[token] = c.TF(token, doc) * c.idfForCalculation(token)
	}

	return vector
}

// TF calculates term frequency.
func (c *TFIDFCalculator) TF(term string, doc *Document) float64 {
	ensureTokens(doc)
	if len(doc.Tokens) == 0 {
		return 0
	}

	count := 0
	for _, token := range doc.Tokens {
		if token == term {
			count++
		}
	}

	return float64(count) / float64(len(doc.Tokens))
}

// IDF calculates inverse document frequency.
func (c *TFIDFCalculator) IDF(term string) float64 {
	return c.idf[term]
}

// Tokenize splits text into lowercase alphanumeric tokens.
func Tokenize(text string) []string {
	var builder strings.Builder
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			continue
		}
		builder.WriteRune(' ')
	}

	return strings.Fields(builder.String())
}

func ensureTokens(doc *Document) {
	if doc == nil || len(doc.Tokens) > 0 {
		return
	}
	doc.Tokens = Tokenize(doc.Content)
}

func (c *TFIDFCalculator) idfForCalculation(term string) float64 {
	if idf, ok := c.idf[term]; ok {
		return idf
	}
	if len(c.documents) == 0 {
		return 0
	}
	return math.Log(float64(len(c.documents)) / 1.0)
}
