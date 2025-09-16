package server

import (
    "net/http"
    "strings"
)

// GET /crypto/{symbol}/history
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
    if !strings.HasSuffix(r.URL.Path, "/history") {
        writeErr(w, http.StatusNotFound, "not found")
        return
    }
    sym := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/history"), "/crypto/")
    sym = strings.TrimSpace(sym)
    if sym == "" || strings.Contains(sym, "/") {
        writeErr(w, http.StatusNotFound, "not found")
        return
    }
    hist, err := s.repo.History(sym)
    if err != nil {
        writeMappedError(w, err, map[int]string{http.StatusNotFound: "not found"})
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{"symbol": sym, "history": hist})
}

