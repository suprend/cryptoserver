package repository

import (
	"math"
	"testing"
)

// absolute + relative tolerance check for floats
func almostEqual(a, b, eps float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	d := math.Abs(a - b)
	if d <= eps {
		return true
	}
	denom := math.Max(math.Abs(a)+math.Abs(b), 1.0)
	return d/denom <= eps
}

// compute expected aggregates directly from history
func deriveStats(h []PriceRecord) (min, max, avg, changePct float64, okChange bool) {
	if len(h) == 0 {
		return 0, 0, 0, 0, false
	}
	min, max = h[0].Price, h[0].Price
	sum := 0.0
	for _, r := range h {
		if r.Price < min {
			min = r.Price
		}
		if r.Price > max {
			max = r.Price
		}
		sum += r.Price
	}
	avg = sum / float64(len(h))
	first := h[0].Price
	last := h[len(h)-1].Price
	if first == 0 {
		return min, max, avg, 0, false // undefined percent change; skip check
	}
	changePct = (last - first) / first * 100
	return min, max, avg, changePct, true
}

func TestStats_SinglePoint(t *testing.T) {
	repo := NewMemoryCryptoRepo()
	c, err := repo.Create("btc")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	s, err := repo.Stats("btc")
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}

	min, max, avg, pct, ok := deriveStats(c.History)
	const eps = 1e-9
	if !almostEqual(s.MinPrice, min, eps) {
		t.Errorf("MinPrice got %v want %v", s.MinPrice, min)
	}
	if !almostEqual(s.MaxPrice, max, eps) {
		t.Errorf("MaxPrice got %v want %v", s.MaxPrice, max)
	}
	if !almostEqual(s.AvgPrice, avg, eps) {
		t.Errorf("AvgPrice got %v want %v", s.AvgPrice, avg)
	}
	if ok && !almostEqual(s.PriceChangePct, pct, eps) {
		t.Errorf("PriceChangePct got %v want %v", s.PriceChangePct, pct)
	}
}

func TestStats_MultiPoint(t *testing.T) {
	repo := NewMemoryCryptoRepo()
	if _, err := repo.Create("eth"); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// add several points to history
	for i := 0; i < 4; i++ { // initial + 4 refreshes => >=5 points
		if _, err := repo.RefreshPrice("eth"); err != nil {
			t.Fatalf("RefreshPrice failed: %v", err)
		}
	}

	c, err := repo.Get("eth")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if len(c.History) < 2 {
		t.Fatalf("need >=2 points, got %d", len(c.History))
	}

	s, err := repo.Stats("eth")
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}

	min, max, avg, pct, ok := deriveStats(c.History)
	const eps = 1e-9
	if !almostEqual(s.MinPrice, min, eps) {
		t.Errorf("MinPrice got %v want %v", s.MinPrice, min)
	}
	if !almostEqual(s.MaxPrice, max, eps) {
		t.Errorf("MaxPrice got %v want %v", s.MaxPrice, max)
	}
	if !almostEqual(s.AvgPrice, avg, eps) {
		t.Errorf("AvgPrice got %v want %v", s.AvgPrice, avg)
	}
	if ok && !almostEqual(s.PriceChangePct, pct, eps) {
		t.Errorf("PriceChangePct got %v want %v", s.PriceChangePct, pct)
	}
}
