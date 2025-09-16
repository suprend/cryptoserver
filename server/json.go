package server

import (
    "encoding/json"
    "net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(status)
    if v == nil {
        _, _ = w.Write([]byte("{}"))
        return
    }
    _ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]any{"error": msg})
}

// writeCrypto wraps a CryptoView into {"crypto": ...} envelope.
func writeCrypto(w http.ResponseWriter, status int, v CryptoView) {
    writeJSON(w, status, map[string]any{"crypto": v})
}
