package geckoclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"cryptoserver/gecko/geckocoins"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var tickerMap map[string]geckocoins.CoinInfo

// Base URL is configurable via env COINGECKO_BASE_URL.
// If not set, auto-detects local fake server, else falls back to real API.
var baseUrl string

var (
    ErrNotFound          = errors.New("not found")
    ErrServiceUnavailable = errors.New("service unavailable")
    ErrBadResponse       = errors.New("bad response")
)

func init() {
	baseUrl = os.Getenv("COINGECKO_BASE_URL")
	if baseUrl == "" {
		if isLocalFakeAlive() {
			baseUrl = "http://127.0.0.1:5050"
		} else {
			baseUrl = "https://api.coingecko.com/api/v3"
		}
	}
	if err := loadAllCryptoNames(); err != nil {
		log.Fatalf("geckoclient initialization failed: %v", err)
	}
}

func isLocalFakeAlive() bool {
	client := &http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get("http://127.0.0.1:5050/coins/list")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil || len(b) == 0 {
		return false
	}
	trimmed := bytes.TrimLeft(b, " \t\r\n")
	return len(trimmed) > 0 && trimmed[0] == '['
}

func loadAllCryptoNames() error {
	//log.Println("LoadAllCryptoNames: starting")
	url := baseUrl + "/coins/list"
	//log.Printf("LoadAllCryptoNames: GET %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	//log.Printf("LoadAllCryptoNames: received %d bytes", len(body))

	// ensure response is JSON array
	trimmed := bytes.TrimLeft(body, " \t\r\n")
	if len(trimmed) == 0 || trimmed[0] != '[' {
		return fmt.Errorf("%w: unexpected response format", ErrBadResponse)
	}

	var coins []geckocoins.CoinInfo
	if err := json.Unmarshal(body, &coins); err != nil {
		return fmt.Errorf("%w: %v", ErrBadResponse, err)
	}
	// detect symbol collisions
	symbolCount := make(map[string]int)
	// reset tickerMap and use lowercase keys
	tickerMap = make(map[string]geckocoins.CoinInfo)
	for _, coin := range coins {
		key := strings.ToLower(coin.Symbol)
		symbolCount[key]++
		if symbolCount[key] > 1 {
			//log.Printf("LoadAllCryptoNames: collision for symbol %q: duplicate id %q (first id %q)", key, coin.ID, tickerMap[key].ID)
			continue
		}
		tickerMap[key] = coin
	}
	//log.Printf("LoadAllCryptoNames: mapped %d symbols", len(tickerMap))
	return nil
}

func GetPrice(symbol string) (float64, error) {
	//log.Printf("GetPrice: called with symbol %q", symbol)
	key := strings.ToLower(symbol)
	info, ok := tickerMap[key]
	var id string
	if ok {
		id = info.ID
	} else {
		id = key
	}
	url := fmt.Sprintf(baseUrl+"/simple/price?ids=%s&vs_currencies=usd", id)
	//log.Printf("GetPrice: requesting URL %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	//log.Printf("GetPrice: response body length %d bytes", len(body))
	var data map[string]map[string]float64
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("%s", body)
		return 0, fmt.Errorf("%w: %v", ErrBadResponse, err)
	}
	if priceData, ok := data[id]; ok {
		if price, ok := priceData["usd"]; ok {
			//log.Printf("GetPrice: price for %q is %f", symbol, price)
			return price, nil
		}
	}
	return 0, fmt.Errorf("%w: price not found for %s", ErrNotFound, symbol)
}

func GetName(symbol string) (string, error) {
	//log.Printf("GetName: called with symbol %q", symbol)
	key := strings.ToLower(symbol)
	info, ok := tickerMap[key]
	if ok {
		return info.Name, nil
	}
	return "", fmt.Errorf("%w: name not found for %s", ErrNotFound, symbol)
}
