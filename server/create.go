package server

import (
    "encoding/json"
    "net/http"
    "strings"
)

// POST /crypto {symbol}
func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
    var req struct{ Symbol string `json:"symbol"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeErr(w, http.StatusBadRequest, "invalid json")
        return
    }
    sym := strings.TrimSpace(req.Symbol)
    if sym == "" {
        writeErr(w, http.StatusBadRequest, "symbol required")
        return
    }
    c, err := s.repo.Create(sym)
    if err != nil {
        writeMappedError(w, err, nil)
        return
    }
    writeCrypto(w, http.StatusCreated, toCryptoView(c))
}
