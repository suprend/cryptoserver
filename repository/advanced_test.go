package repository

import (
	"testing"
)

// TestRefreshPriceSuccess ensures RefreshPrice updates price, timestamp, and history.
func TestMemoryCryptoRepo_RefreshPriceSuccess(t *testing.T) {
	repo := NewMemoryCryptoRepo()

	// Create initial entry
	c, err := repo.Create("btc")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	// Initial invariants
	if len(c.History) != 1 {
		t.Fatalf("expected initial History length 1, got %d", len(c.History))
	}
	if c.CurrentPrice != c.History[0].Price {
		t.Errorf("expected CurrentPrice == last history price, got %v vs %v", c.CurrentPrice, c.History[0].Price)
	}
	if !c.LastUpdated.Equal(c.History[0].Timestamp) {
		t.Errorf("expected LastUpdated == last history timestamp, got %v vs %v", c.LastUpdated, c.History[0].Timestamp)
	}
	initialLastUpdated := c.LastUpdated
	initialHistoryLen := len(c.History)

	// Refresh price
	updated, err := repo.RefreshPrice("btc")
	if err != nil {
		t.Fatalf("RefreshPrice failed: %v", err)
	}

	// Timestamp should be advanced
	if !updated.LastUpdated.After(initialLastUpdated) {
		t.Errorf("expected LastUpdated after %v, got %v", initialLastUpdated, updated.LastUpdated)
	}
	// History should grow by one
	if len(updated.History) != initialHistoryLen+1 {
		t.Errorf("expected History length %d, got %d", initialHistoryLen+1, len(updated.History))
	}
	// Invariants after refresh
	n := len(updated.History)
	last := updated.History[n-1]
	prev := updated.History[n-2]
	if updated.CurrentPrice != last.Price {
		t.Errorf("expected CurrentPrice == last history price, got %v vs %v", updated.CurrentPrice, last.Price)
	}
	if !updated.LastUpdated.Equal(last.Timestamp) {
		t.Errorf("expected LastUpdated == last history timestamp, got %v vs %v", updated.LastUpdated, last.Timestamp)
	}
	if !prev.Timestamp.Before(last.Timestamp) {
		t.Errorf("expected history to be strictly increasing by timestamp, got prev=%v, last=%v", prev.Timestamp, last.Timestamp)
	}
}

// TestMultipleCreateRefresh tests creating multiple cryptos and refreshing them.
func TestMultipleCreateRefresh(t *testing.T) {
	repo := NewMemoryCryptoRepo()
	symbols := []string{"eth", "btc"}

	// Create initial entries
	for _, sym := range symbols {
		c, err := repo.Create(sym)
		if err != nil {
			t.Fatalf("Create failed for %s: %v", sym, err)
		}
		if len(c.History) != 1 {
			t.Errorf("expected initial History length 1 for %s, got %d", sym, len(c.History))
		}
		if c.CurrentPrice != c.History[0].Price {
			t.Errorf("%s: expected CurrentPrice == last history price, got %v vs %v", sym, c.CurrentPrice, c.History[0].Price)
		}
		if !c.LastUpdated.Equal(c.History[0].Timestamp) {
			t.Errorf("%s: expected LastUpdated == last history timestamp, got %v vs %v", sym, c.LastUpdated, c.History[0].Timestamp)
		}
	}

	// List
	list, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != len(symbols) {
		t.Errorf("expected list length %d, got %d", len(symbols), len(list))
	}

	// Refresh and verify each
	for _, sym := range symbols {
		before, err := repo.Get(sym)
		if err != nil {
			t.Fatalf("Get failed for %s: %v", sym, err)
		}
		prevUpdated := before.LastUpdated
		prevHistoryLen := len(before.History)

		updated, err := repo.RefreshPrice(sym)
		if err != nil {
			t.Fatalf("RefreshPrice failed for %s: %v", sym, err)
		}
		if !updated.LastUpdated.After(prevUpdated) {
			t.Errorf("expected LastUpdated after %v for %s, got %v", prevUpdated, sym, updated.LastUpdated)
		}
		if len(updated.History) != prevHistoryLen+1 {
			t.Errorf("expected History length %d for %s, got %d", prevHistoryLen+1, sym, len(updated.History))
		}
		// Invariants after refresh
		n := len(updated.History)
		last := updated.History[n-1]
		if updated.CurrentPrice != last.Price {
			t.Errorf("%s: expected CurrentPrice == last history price, got %v vs %v", sym, updated.CurrentPrice, last.Price)
		}
		if !updated.LastUpdated.Equal(last.Timestamp) {
			t.Errorf("%s: expected LastUpdated == last history timestamp, got %v vs %v", sym, updated.LastUpdated, last.Timestamp)
		}
		if n >= 2 {
			prevRec := updated.History[n-2]
			if !prevRec.Timestamp.Before(last.Timestamp) {
				t.Errorf("%s: expected history to be strictly increasing by timestamp, got prev=%v, last=%v", sym, prevRec.Timestamp, last.Timestamp)
			}
		}
	}
}
