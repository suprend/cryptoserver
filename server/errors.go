package server

import (
    "errors"
    "net/http"

    "cryptoserver/repository"
)

// mapRepoError converts repository/domain errors into HTTP status + default message.
func mapRepoError(err error) (int, string) {
    switch {
    case errors.Is(err, repository.ErrInvalidSymbol):
        return http.StatusBadRequest, err.Error()
    case errors.Is(err, repository.ErrAlreadyExists):
        return http.StatusConflict, err.Error()
    case errors.Is(err, repository.ErrNotFound):
        return http.StatusNotFound, "not found"
    case errors.Is(err, repository.ErrNameUnavailable), errors.Is(err, repository.ErrPriceUnavailable):
        return http.StatusBadGateway, err.Error()
    case errors.Is(err, repository.ErrServiceUnavailable):
        return http.StatusServiceUnavailable, err.Error()
    default:
        return http.StatusInternalServerError, err.Error()
    }
}

// writeMappedError writes a JSON error using mapping rules and optional per-status overrides for message.
func writeMappedError(w http.ResponseWriter, err error, overrides map[int]string) {
    status, msg := mapRepoError(err)
    if overrides != nil {
        if m, ok := overrides[status]; ok && m != "" {
            msg = m
        }
    }
    writeErr(w, status, msg)
}

