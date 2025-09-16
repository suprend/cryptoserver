package repository

import (
	"slices"
	"testing"
	"time"
)

// TestHistoryImmutability verifies that repository.History returns a clone
// and mutating the returned slice does not affect the internal stored history.
func TestHistoryImmutability(t *testing.T) {
	repo := NewMemoryCryptoRepo()

	if _, err := repo.Create("btc"); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	// Add a couple of points so history has multiple records
	if _, err := repo.RefreshPrice("btc"); err != nil {
		t.Fatalf("RefreshPrice failed: %v", err)
	}
	if _, err := repo.RefreshPrice("btc"); err != nil {
		t.Fatalf("RefreshPrice failed: %v", err)
	}

	histA, err := repo.History("btc")
	if err != nil {
		t.Fatalf("History failed: %v", err)
	}
	if len(histA) != 3 {
		t.Fatalf("expected lenght history equal to 3, but got %v", len(histA))
	}

	// Keep a baseline copy of what repo returned
	baseline := slices.Clone(histA)

	// Mutate the returned slice deeply and structurally
	histA[0].Price += 123.456
	histA[0].Timestamp = time.Time{}
	histA = append(histA, PriceRecord{Price: -1, Timestamp: time.Unix(1, 0)})

	// Fetch history again from repo and ensure it was not affected
	histB, err := repo.History("btc")
	if err != nil {
		t.Fatalf("History failed: %v", err)
	}

	if len(histB) != len(baseline) {
		t.Fatalf("internal history length changed: got %d want %d", len(histB), len(baseline))
	}

	for i := range baseline {
		if histB[i].Price != baseline[i].Price {
			t.Errorf("record[%d].Price changed internally: got %v want %v", i, histB[i].Price, baseline[i].Price)
		}
		if !histB[i].Timestamp.Equal(baseline[i].Timestamp) {
			t.Errorf("record[%d].Timestamp changed internally: got %v want %v", i, histB[i].Timestamp, baseline[i].Timestamp)
		}
	}
}
