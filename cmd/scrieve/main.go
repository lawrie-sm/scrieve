package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lawrie-sm/scrieve/internal/data"
	"log"
	"net/http"
	"os"
)

// service is used within main to store pointers to the database and muxer
type service struct {
	db  *data.DB
	mux *http.ServeMux
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// newService returns a new service, with DB connection and routes created
func newService() *service {
	s := &service{}
	s.setupRoutes()

	s.db = &data.DB{}
	s.db.Connect()
	s.db.Setup()

	env := os.Getenv("SCRIEVE_ENV")
	if env == "development" {
		log.SetFlags(log.Lshortfile)
	} else if env == "production" {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}
	return s
}

// run is called by main, allowing it to return an error
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
