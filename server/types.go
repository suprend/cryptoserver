package server

import "cryptoserver/repository"

type Server struct {
    repo repository.CryptoRepository
}

func New(repo repository.CryptoRepository) *Server {
    return &Server{repo: repo}
}

