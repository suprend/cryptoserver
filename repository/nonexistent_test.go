package repository

import "testing"

// TestNonexistentSymbolOperations groups tests for non-existent symbol operations.
func TestNonexistentSymbolOperations(t *testing.T) {
	repo := NewMemoryCryptoRepo()
	if _, err := repo.Create("missing"); err == nil {
		t.Errorf("expected error on Create for missing symbol, got nil")
	}
	if _, err := repo.Get("missing"); err == nil {
		t.Errorf("expected error on Get for missing symbol, got nil")
	}
	if err := repo.Delete("missing"); err == nil {
		t.Errorf("expected error on Delete for missing symbol, got nil")
	}
	if _, err := repo.Stats("missing"); err == nil {
		t.Errorf("expected error on Stats for missing symbol, got nil")
	}
	if _, err := repo.RefreshPrice("missing"); err == nil {
		t.Errorf("expected error on RefreshPrice for missing symbol, got nil")
	}
	if _, err := repo.History("missing"); err == nil {
		t.Errorf("expected error on History for missing symbol, got nil")
	}
}
