package server

import (
	"net/http"
)

// GET /crypto
func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	items, err := s.repo.List()
	if err != nil {
		writeMappedError(w, err, nil)
		return
	}
	views := make([]CryptoView, 0, len(items))
	for _, c := range items {
		views = append(views, toCryptoView(c))
	}
	writeJSON(w, http.StatusOK, map[string]any{"cryptos": views})
}
