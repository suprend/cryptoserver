package server

import (
    "net/http"
    "strings"
)

// DELETE /crypto/{symbol}
func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
    // Ensure no extra segments after symbol
    sym := strings.TrimPrefix(r.URL.Path, "/crypto/")
    if sym == "" || strings.Contains(sym, "/") {
        writeErr(w, http.StatusNotFound, "not found")
        return
    }
    if err := s.repo.Delete(sym); err != nil {
        writeMappedError(w, err, map[int]string{http.StatusNotFound: "not found"})
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{})
}

