package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type webserver struct {
	cfg *Config
	db  database
	r   *chi.Mux
}

func (s *webserver) GetRegistrations(w http.ResponseWriter, r *http.Request) {
	users, err := s.db.GetUserRegistrations()
	if err != nil {
		log.Printf("Could not get user registrations: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not get user registrations"))
		s.cfg.FatalExit()
	}
	for _, u := range users {
		w.Write([]byte(fmt.Sprintf("User: ID=%d, FirstName=%s, LastName=%s, Email=%s \n", u.ID, u.FirstName, u.LastName, u.Email)))
		w.WriteHeader(http.StatusOK)
	}
}

func (s *webserver) Register(w http.ResponseWriter, r *http.Request) {
	u := new(User)
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please enter correct registration data!!"))
		return
	}
	s.db.PutuserRegistration(u)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registered!!"))
}

func RunRestserver(cfg *Config, db database) {
	log.Println("Starting REST server")
	webserver := &webserver{cfg: cfg, db: db, r: chi.NewRouter()}

	r := webserver.r
	r.Use(middleware.Logger)

	r.Get("/", webserver.GetRegistrations)
	r.Post("/register", webserver.Register)

	http.ListenAndServe(":3000", r)
}
