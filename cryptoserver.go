package main

import (
    "errors"
    "fmt"
    "cryptoserver/repository"
    "cryptoserver/server"
    "log"
    "net/http"
    "os"
    "strconv"
)

func envPort() (int, error) {
    p := os.Getenv("PORT")
    if p == "" {
        return 8080, nil
    }
    n, err := strconv.Atoi(p)
	if err != nil || n <= 0 || n > 65535 {
		return 0, errors.New("invalid PORT")
	}
	return n, nil
}

func main() {
    s := server.New(repository.NewMemoryCryptoRepo())
    p, err := envPort()
    if err != nil {
        log.Fatal(err)
    }
    addr := fmt.Sprintf(":%d", p)
    log.Printf("listening on %s", addr)
    if err := http.ListenAndServe(addr, s); err != nil {
        log.Fatal(err)
    }
}
