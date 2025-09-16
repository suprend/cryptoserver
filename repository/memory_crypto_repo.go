package repository

import (
    "errors"
    "fmt"
    "cryptoserver/gecko/geckoclient"
    "slices"
    "strings"
    "sync"
    "time"
)

// MemoryCryptoRepo хранит криптовалюты в памяти.
type MemoryCryptoRepo struct {
	data map[string]Crypto
	mu   sync.Mutex
}

func NewMemoryCryptoRepo() *MemoryCryptoRepo {
	return &MemoryCryptoRepo{
		data: make(map[string]Crypto),
	}
}

func (r *MemoryCryptoRepo) Create(symbol string) (Crypto, error) {
    symbol = strings.ToLower(strings.TrimSpace(symbol))
    if symbol == "" {
        return Crypto{}, ErrInvalidSymbol
    }

	r.mu.Lock()
	if _, exists := r.data[symbol]; exists {
		r.mu.Unlock()
		return Crypto{}, ErrAlreadyExists
	}
	r.mu.Unlock()

    name, err := geckoclient.GetName(symbol)
    if err != nil {
        switch {
        case errors.Is(err, geckoclient.ErrNotFound):
            return Crypto{}, fmt.Errorf("%w: %v", ErrInvalidSymbol, err)
        case errors.Is(err, geckoclient.ErrServiceUnavailable), errors.Is(err, geckoclient.ErrBadResponse):
            return Crypto{}, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
        default:
            return Crypto{}, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
        }
    }
    price, err := geckoclient.GetPrice(symbol)
    if err != nil {
        switch {
        case errors.Is(err, geckoclient.ErrNotFound):
            return Crypto{}, fmt.Errorf("%w: %v", ErrPriceUnavailable, err)
        case errors.Is(err, geckoclient.ErrServiceUnavailable), errors.Is(err, geckoclient.ErrBadResponse):
            return Crypto{}, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
        default:
            return Crypto{}, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
        }
    }

	now := time.Now()

	c := Crypto{
		Symbol:       symbol,
		Name:         name,
		CurrentPrice: price,
		LastUpdated:  now,
		History: []PriceRecord{
			{Price: price, Timestamp: now},
		},
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[symbol]; exists {
		return Crypto{}, ErrAlreadyExists
	}

	r.data[symbol] = c
	return c.Copy(), nil
}

func (r *MemoryCryptoRepo) List() ([]Crypto, error) {
	r.mu.Lock()
	result := make([]Crypto, 0, len(r.data))
	for _, c := range r.data {
		result = append(result, c.Copy())
	}
	r.mu.Unlock()

	return result, nil
}

func (r *MemoryCryptoRepo) Get(symbol string) (Crypto, error) {
    symbol = strings.ToLower(strings.TrimSpace(symbol))
    if symbol == "" {
        return Crypto{}, ErrInvalidSymbol
    }

	r.mu.Lock()
	defer r.mu.Unlock()

	if c, exists := r.data[symbol]; exists {
		return c.Copy(), nil
	}
	return Crypto{}, ErrNotFound
}

func (r *MemoryCryptoRepo) Delete(symbol string) error {
    symbol = strings.ToLower(strings.TrimSpace(symbol))
    if symbol == "" {
        return ErrInvalidSymbol
    }

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[symbol]; !exists {
		return ErrNotFound
	}
	delete(r.data, symbol)
	return nil
}

func (r *MemoryCryptoRepo) RefreshPrice(symbol string) (Crypto, error) {
    symbol = strings.ToLower(strings.TrimSpace(symbol))
    if symbol == "" {
        return Crypto{}, ErrInvalidSymbol
    }

	r.mu.Lock()
	_, exists := r.data[symbol]
	r.mu.Unlock()

	if !exists {
		return Crypto{}, ErrNotFound
	}
    price, err := geckoclient.GetPrice(symbol)
    if err != nil {
        switch {
        case errors.Is(err, geckoclient.ErrNotFound):
            return Crypto{}, fmt.Errorf("%w: %v", ErrPriceUnavailable, err)
        case errors.Is(err, geckoclient.ErrServiceUnavailable), errors.Is(err, geckoclient.ErrBadResponse):
            return Crypto{}, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
        default:
            return Crypto{}, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
        }
    }
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

    c, exists := r.data[symbol]
    if !exists {
        return Crypto{}, ErrNotFound
    }

	c.CurrentPrice = price
	c.LastUpdated = now

	c.History = append(c.History, PriceRecord{Price: price, Timestamp: now})
	if len(c.History) > 100 {
		c.History = c.History[len(c.History)-100:]
		c.History = slices.Clone(c.History)
	}
	r.data[symbol] = c

	return c.Copy(), nil
}

func (r *MemoryCryptoRepo) History(symbol string) ([]PriceRecord, error) {
    symbol = strings.ToLower(strings.TrimSpace(symbol))
    if symbol == "" {
        return nil, ErrInvalidSymbol
    }

	r.mu.Lock()
	defer r.mu.Unlock()

    c, exists := r.data[symbol]
    if !exists {
        return nil, ErrNotFound
    }

	return slices.Clone(c.History), nil
}

func (r *MemoryCryptoRepo) Stats(symbol string) (PriceStats, error) {
    symbol = strings.ToLower(strings.TrimSpace(symbol))
    if symbol == "" {
        return PriceStats{}, ErrInvalidSymbol
    }

	r.mu.Lock()
    c, exists := r.data[symbol]
    if !exists {
        r.mu.Unlock()
        return PriceStats{}, ErrNotFound
    }
	c = c.Copy()
	r.mu.Unlock()

	h := c.History
	if len(h) == 0 {
		return PriceStats{}, nil
	}

	minP, maxP := h[0].Price, h[0].Price
	sum := 0.0
	for _, rec := range h {
		p := rec.Price
		minP = min(minP, p)
		maxP = max(maxP, p)
		sum += p
	}
	first := h[0].Price
	last := h[len(h)-1].Price
	change := last - first
	var pct float64
	if first != 0 {
		pct = change / first * 100
	}

    return PriceStats{
        MinPrice:       minP,
        MaxPrice:       maxP,
        AvgPrice:       sum / float64(len(h)),
        PriceChange:    change,
        PriceChangePct: pct,
        RecordsCount:   len(h),
    }, nil
}
