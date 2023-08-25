package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

type Message struct {
	Msg string
}

const addr = "0.0.0.0:7777"

func main() {
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		m := Message{}
		err := d.Decode(&m)
		if err != nil {
			json.NewEncoder(w).Encode(errors.New("unable to decode request body"))
			return
		}
		log.Printf("Received message: %v\n", m.Msg)

		json.NewEncoder(w).Encode(m)
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check called")
		w.Write([]byte("OK"))
	})
	r.Get("/healthfail", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check failed")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	})
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		log.Printf("Listening on %q", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	<- c

}
