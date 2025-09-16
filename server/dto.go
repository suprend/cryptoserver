package server

import (
    "cryptoserver/repository"
    "time"
)

// CryptoView is the transport shape for crypto without history.
type CryptoView struct {
    Symbol       string    `json:"symbol"`
    Name         string    `json:"name"`
    CurrentPrice float64   `json:"current_price"`
    LastUpdated  time.Time `json:"last_updated"`
}

func toCryptoView(c repository.Crypto) CryptoView {
    return CryptoView{
        Symbol:       c.Symbol,
        Name:         c.Name,
        CurrentPrice: c.CurrentPrice,
        LastUpdated:  c.LastUpdated,
    }
}

