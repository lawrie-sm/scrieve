package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lawrie-sm/scrieve/internal/data"
)

type service struct {
	db  *data.DB
	mux *http.ServeMux
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func newService() *service {
	s := &service{}
	s.setupRoutes()

	s.db = &data.DB{}
	s.db.Connect()
	s.db.Setup()

	log.SetFlags(log.Lshortfile)
	return s
}

func run() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	service := newService()
	port := os.Getenv("PORT")
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	server := http.Server{
		Addr:    addr,
		Handler: service,
	}

	log.Println("Listening on", port)
	return server.ListenAndServe()
}

func main() {
	err := run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
