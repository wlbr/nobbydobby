package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Guest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

func (g *Guest) String() string {
	return fmt.Sprintf(`"%s","%s","%s"`, g.Firstname, g.Lastname, g.Email)
}

type felixHandler struct {
	mux            *http.ServeMux
	guests         map[string]*Guest
	guestsfilename string
	m              sync.RWMutex
}

func NewFelixHandler() *felixHandler {
	return &felixHandler{mux: http.NewServeMux(), guests: make(map[string]*Guest), guestsfilename: "guests.db"}
}

func (h *felixHandler) ReadGuests(fname string) error {
	h.m.Lock()
	defer h.m.Unlock()
	f, err := os.Open(fname)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		log.Printf("Error during read Guests from file: %s", err)
		return fmt.Errorf("Error during read Guests from file: %s", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		j := scanner.Bytes()
		g := &Guest{}
		if err := json.Unmarshal(j, g); err != nil {
			log.Printf("Error unmarshalling guest from file: %s", err)
		} else {
			h.guests[g.Email] = g
		}
	}
	return nil
}

func (h *felixHandler) AddGuest(g *Guest) error {
	h.m.Lock()
	defer h.m.Unlock()
	if _, ok := h.guests[g.Email]; ok {
		log.Printf("Email %s already on guestlist, skipped", g.Email)
		return fmt.Errorf("Email %s already on guestlist, skipped", g.Email)
	}
	h.guests[g.Email] = g

	// Append the new guest to the file
	f, err := os.OpenFile(h.guestsfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Error opening guests file for append: %s", err)
	}
	defer f.Close()

	var j []byte
	if j, err = json.Marshal(g); err != nil {
		return fmt.Errorf("Error marshalling new guest for append: %s", err)

	}
	if _, err = f.Write(j); err != nil {
		return fmt.Errorf("Error writing new guest to file: %s", err)
	}
	if _, err = f.WriteString("\n"); err != nil {
		return fmt.Errorf("Error writing newline after new guest to file: %s", err)

	}
	return nil
}

func (h *felixHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got request from %s: %s %s", r.RemoteAddr, r.Method, r.URL)
	h.mux.ServeHTTP(w, r)
}

func (h *felixHandler) handleReadGuests(w http.ResponseWriter, r *http.Request) {
	h.m.RLock()
	defer h.m.RUnlock()

	w.Header().Add("Content-Type", "text/plain")
	//w.Header().Add("Content-Type", "text/csv")
	for _, g := range h.guests {
		fmt.Fprintln(w, g)
	}
}

func (h *felixHandler) handleRegisterGuest(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	g := &Guest{}
	if err := json.Unmarshal(b, g); err != nil {
		log.Printf("Error unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.AddGuest(g); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusAlreadyReported)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	h := NewFelixHandler()
	if err := h.ReadGuests(h.guestsfilename); err != nil {
		log.Fatalf("Error reading guests database.")
	}

	h.mux.HandleFunc("GET /all", h.handleReadGuests)
	h.mux.HandleFunc("POST /register", h.handleRegisterGuest)

	http.ListenAndServe(":8080", h)
}
