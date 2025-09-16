package repository

import (
	"errors"
	"slices"
	"time"
)

type PriceRecord struct {
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type Crypto struct {
	Symbol       string        `json:"symbol"`
	Name         string        `json:"name"`
	CurrentPrice float64       `json:"current_price"`
	LastUpdated  time.Time     `json:"last_updated"`
	History      []PriceRecord `json:"history"`
}

type PriceStats struct {
    MinPrice       float64 `json:"min_price"`
    MaxPrice       float64 `json:"max_price"`
    AvgPrice       float64 `json:"avg_price"`
    PriceChange    float64 `json:"price_change"`     // last - first
    PriceChangePct float64 `json:"price_change_percent"` // (last-first)/first*100
    RecordsCount   int     `json:"records_count"`
}

type CryptoRepository interface {
	Create(symbol string) (Crypto, error)
	Get(symbol string) (Crypto, error)
	List() ([]Crypto, error)
	Delete(symbol string) error
	RefreshPrice(symbol string) (Crypto, error)
	History(symbol string) ([]PriceRecord, error)
	Stats(symbol string) (PriceStats, error)
}

func (c Crypto) Copy() Crypto {
	out := c
	out.History = slices.Clone(c.History)
	return out
}

var (
    ErrAlreadyExists = errors.New("crypto already exists")
    ErrNotFound      = errors.New("crypto not found")
    ErrInvalidSymbol = errors.New("invalid symbol")
    ErrNameUnavailable  = errors.New("name unavailable")
    ErrPriceUnavailable = errors.New("price unavailable")
    ErrServiceUnavailable = errors.New("service unavailable")
)
