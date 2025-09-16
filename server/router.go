package server

import (
    "net/http"
    "strings"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch {
    case r.Method == http.MethodGet && r.URL.Path == "/crypto":
        s.handleList(w, r)
        return
    case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/crypto/") && strings.HasSuffix(r.URL.Path, "/history"):
        s.handleHistory(w, r)
        return
    case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/crypto/") && strings.HasSuffix(r.URL.Path, "/stats"):
        s.handleStats(w, r)
        return
    case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/crypto/") && strings.HasSuffix(r.URL.Path, "/refresh"):
        s.handleRefresh(w, r)
        return
    case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/crypto/"):
        s.handleDelete(w, r)
        return
    case r.Method == http.MethodPost && r.URL.Path == "/crypto":
        s.handleCreate(w, r)
        return
    case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/crypto/"):
        s.handleGet(w, r)
        return
    }
    writeErr(w, http.StatusNotFound, "not found")
}
