package geckoclient

import "testing"

// TestGetNameKnownSymbol checks that a known symbol returns a non-empty name.
func TestGetNameKnownSymbol(t *testing.T) {
	name, err := GetName("btc")
	if err != nil {
		t.Fatalf("GetName(\"btc\") error = %v", err)
	}
	if name == "" {
		t.Errorf("GetName(\"btc\") returned empty name")
	}
}

// TestGetNameUnknownSymbol ensures GetName returns an error for an unknown symbol.
func TestGetNameUnknownSymbol(t *testing.T) {
	if _, err := GetName("unknownsymbol123"); err == nil {
		t.Errorf("GetName(\"unknownsymbol123\") expected error, got nil")
	}
}

// TestGetPriceKnownSymbol checks that a known symbol returns a positive price.
func TestGetPriceKnownSymbol(t *testing.T) {
	price, err := GetPrice("btc")
	if err != nil {
		t.Fatalf("GetPrice(\"btc\") error = %v", err)
	}
	if price <= 0 {
		t.Errorf("GetPrice(\"btc\") = %f; want > 0", price)
	}
}

// TestGetPriceUnknownSymbol ensures GetPrice returns an error for an unknown symbol.
func TestGetPriceUnknownSymbol(t *testing.T) {
	if _, err := GetPrice("unknownsymbol123"); err == nil {
		t.Errorf("GetPrice(\"unknownsymbol123\") expected error, got nil")
	}
}
