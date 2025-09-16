package server

import (
    "net/http"
    "strings"
)

// PUT /crypto/{symbol}/refresh
func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/crypto/"), "/")
    if len(parts) != 2 || parts[1] != "refresh" {
        writeErr(w, http.StatusNotFound, "not found")
        return
    }
    sym := strings.TrimSpace(parts[0])
    if sym == "" {
        writeErr(w, http.StatusBadRequest, "symbol required")
        return
    }
    c, err := s.repo.RefreshPrice(sym)
    if err != nil {
        writeMappedError(w, err, nil)
        return
    }
    writeCrypto(w, http.StatusOK, toCryptoView(c))
}
