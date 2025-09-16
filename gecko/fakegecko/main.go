package main

import (
	"encoding/json"
	"cryptoserver/gecko/geckocoins"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Coin is the minimal shape expected from /api/v3/coins/list
// We only need ID to filter unknown ids later.
func main() {
	addr := "127.0.0.1:5050"
	listPath := "crypto_list.json"

	// Load and cache coins list bytes. Also build a set of known ids (lowercase).
	coinsBytes, idSet := loadCoins(listPath)

	mux := http.NewServeMux()

	// GET /coins/list — return file verbatim
	mux.HandleFunc("/coins/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(coinsBytes)
	})

	// GET /simple/price?ids=...&vs_currencies=...
	mux.HandleFunc("/simple/price", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		ids := splitCSV(q.Get("ids"))
		vcs := splitCSV(q.Get("vs_currencies"))
		if len(ids) == 0 {
			writeJSONError(w, http.StatusBadRequest, "missing ids")
			return
		}
		if len(vcs) == 0 {
			writeJSONError(w, http.StatusBadRequest, "missing vs_currencies")
			return
		}

		// Lowercase for response keys and filtering.
		for i := range ids {
			ids[i] = strings.ToLower(ids[i])
		}
		for i := range vcs {
			vcs[i] = strings.ToLower(vcs[i])
		}

		// Build response: omit unknown ids; include requested currencies.
		resp := make(map[string]map[string]float64, len(ids))
		for _, id := range ids {
			if _, ok := idSet[id]; !ok {
				continue
			}
			inner := make(map[string]float64, len(vcs))
			for _, c := range vcs {
				inner[c] = randomPrice(0.000001, 1_000_000)
			}
			resp[id] = inner
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(resp)
	})

	log.Printf("Fake CoinGecko server: http://%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func loadCoins(path string) ([]byte, map[string]struct{}) {
	if path == "" {
		log.Fatal("-list path must not be empty")
	}
	p := filepath.Clean(path)
	b, err := os.ReadFile(p)
	if err != nil {
		log.Fatalf("read coins list: %v", err)
	}
	var coins []geckocoins.CoinInfo
	if err := json.Unmarshal(b, &coins); err != nil {
		log.Fatalf("parse coins list JSON: %v", err)
	}
	ids := make(map[string]struct{}, len(coins))
	for _, c := range coins {
		id := strings.ToLower(strings.TrimSpace(c.ID))
		if id != "" {
			ids[id] = struct{}{}
		}
	}
	return b, ids
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// randomPrice returns a uniform random float64 in [min,max] rounded to 6 decimals.
func randomPrice(min, max float64) float64 {
	if max <= min {
		return min
	}
	u := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() // [0,1)
	return math.Round((min+u*(max-min))*1e6) / 1e6
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": msg})
}
