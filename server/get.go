package server

import (
    "net/http"
    "strings"
)

// GET /crypto/{symbol}
func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
    sym := strings.TrimPrefix(r.URL.Path, "/crypto/")
    sym = strings.TrimSpace(sym)
    if sym == "" || strings.Contains(sym, "/") {
        writeErr(w, http.StatusNotFound, "not found")
        return
    }
    c, err := s.repo.Get(sym)
    if err != nil {
        writeMappedError(w, err, map[int]string{http.StatusNotFound: "not found"})
        return
    }
    // Do not include history in this view
    writeJSON(w, http.StatusOK, toCryptoView(c))
}
