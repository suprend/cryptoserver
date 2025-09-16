package server

import (
    "net/http"
    "strings"
)

// GET /crypto/{symbol}/stats
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
    if !strings.HasSuffix(r.URL.Path, "/stats") {
        writeErr(w, http.StatusNotFound, "not found")
        return
    }
    sym := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/stats"), "/crypto/")
    sym = strings.TrimSpace(sym)
    if sym == "" || strings.Contains(sym, "/") {
        writeErr(w, http.StatusNotFound, "not found")
        return
    }
    // Get current price and history length
    c, err := s.repo.Get(sym)
    if err != nil {
        writeMappedError(w, err, map[int]string{http.StatusNotFound: "not found"})
        return
    }
    st, err := s.repo.Stats(sym)
    if err != nil {
        writeMappedError(w, err, map[int]string{http.StatusNotFound: "not found"})
        return
    }
    // Build response with domain stats directly (tags match external contract)
    resp := map[string]any{
        "symbol":        c.Symbol,
        "current_price": c.CurrentPrice,
        "stats":         st,
    }
    writeJSON(w, http.StatusOK, resp)
}
