package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Server struct {
	db     *sql.DB
	router http.Handler
}

func New(db *sql.DB) *Server {
	router := NewRouter(db)
	return &Server{db: db, router: router}
}

func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
