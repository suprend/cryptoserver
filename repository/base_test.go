package repository

import "testing"

func TestMemoryCryptoRepo_CRUD(t *testing.T) {
	var repo CryptoRepository = NewMemoryCryptoRepo()

	// Test Create
	c, err := repo.Create("btc")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if c.Symbol != "btc" {
		t.Errorf("expected symbol 'btc', got '%s'", c.Symbol)
	}
	if c.Name == "" {
		t.Errorf("expected non-empty Name, got empty")
	}

	// Test Get
	g, err := repo.Get("btc")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !g.LastUpdated.Equal(c.LastUpdated) || g.Name != c.Name || g.Symbol != c.Symbol || g.CurrentPrice != c.CurrentPrice {
		t.Errorf("Get returned a different Crypto: got %+v, want %+v", g, c)
	}

	// Test List
	list, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected list length 1, got %d", len(list))
	}

	// Test Delete
	if err := repo.Delete("btc"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := repo.Get("btc"); err == nil {
		t.Errorf("expected error after Delete, got nil")
	}
}

// TestDuplicateCreate ensures creating the same symbol twice returns an error.
func TestMemoryCryptoRepo_DuplicateCreate(t *testing.T) {
	repo := NewMemoryCryptoRepo()
	if _, err := repo.Create("btc"); err != nil {
		t.Fatalf("first Create failed: %v", err)
	}
	if _, err := repo.Create("btc"); err == nil {
		t.Errorf("expected error on duplicate Create, got nil")
	}

	// Test List
	list, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected list length 1, got %d", len(list))
	}
}
